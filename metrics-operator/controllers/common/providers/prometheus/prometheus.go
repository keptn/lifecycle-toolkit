package prometheus

import (
	"context"
	"fmt"      //nolint:gci
	"net/http" //nolint:gci
	"time"

	"github.com/go-logr/logr"
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	promapi "github.com/prometheus/client_golang/api"
	prometheus "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

type KeptnPrometheusProvider struct {
	Log        logr.Logger
	HttpClient http.Client
}

// EvaluateQuery fetches the SLI values from prometheus provider
func (r *KeptnPrometheusProvider) EvaluateQuery(ctx context.Context, analysisValue metricsapi.AnalysisValue, provider metricsapi.KeptnMetricsProvider) (string, []byte, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	client, err := promapi.NewClient(promapi.Config{Address: provider.Spec.TargetServer, Client: &r.HttpClient})
	if err != nil {
		return "", nil, err
	}

	api := prometheus.NewAPI(client)
	if analysisValue.Spec.Timeframe != nil {
		result, warnings, err := evaluateQueryWithRange(ctx, analysisValue, r, api)
		if err != nil {
			return "", nil, err
		}
		if len(warnings) != 0 {
			r.Log.Info("Prometheus API returned warnings: " + warnings[0])
		}
		return getResultForMatrix(result, r)
	} else {
		result, warnings, err := evaluateQueryWithoutRange(ctx, analysisValue, r, api)
		if err != nil {
			return "", nil, err
		}
		if len(warnings) != 0 {
			r.Log.Info("Prometheus API returned warnings: " + warnings[0])
		}
		return getResultForVector(result, r)
	}
}

func evaluateQueryWithRange(ctx context.Context, analysisValue metricsapi.AnalysisValue, r *KeptnPrometheusProvider, api prometheus.API) (model.Value, prometheus.Warnings, error) {
	queryInterval := analysisValue.Spec.Timeframe.To.Time.Sub(analysisValue.Spec.Timeframe.From.Time)
	r.Log.Info(fmt.Sprintf(
		"Running query: /api/v1/query_range?query=%s&start=%d&end=%d&step=%v",
		analysisValue.Status.Query,
		analysisValue.Spec.Timeframe.From.Time.Unix(), analysisValue.Spec.Timeframe.To.Time.Unix(),
		queryInterval,
	))
	queryRange := prometheus.Range{
		Start: analysisValue.Spec.Timeframe.From.Time.UTC(),
		End:   analysisValue.Spec.Timeframe.To.Time.UTC(),
		Step:  queryInterval,
	}
	result, warnings, err := api.QueryRange(
		ctx,
		analysisValue.Status.Query,
		queryRange,
		[]prometheus.Option{}...,
	)
	if err != nil {
		return nil, nil, err
	}
	return result, warnings, nil
}

func evaluateQueryWithoutRange(ctx context.Context, analysisValue metricsapi.AnalysisValue, r *KeptnPrometheusProvider, api prometheus.API) (model.Value, prometheus.Warnings, error) {
	queryTime := time.Now().UTC()
	r.Log.Info(fmt.Sprintf(
		"Running query: /api/v1/query?query=%s&time=%d",
		analysisValue.Status.Query,
		queryTime.Unix(),
	))
	result, warnings, err := api.Query(
		ctx,
		analysisValue.Status.Query,
		queryTime,
		[]prometheus.Option{}...,
	)
	if err != nil {
		return nil, nil, err
	}
	return result, warnings, nil
}

func getResultForMatrix(result model.Value, r *KeptnPrometheusProvider) (string, []byte, error) {
	// check if we can cast the result to a matrix
	resultMatrix, ok := result.(model.Matrix)
	if !ok {
		return "", nil, fmt.Errorf("could not cast result")
	}
	// We are only allowed to return one value, if not the query may be malformed
	// we are using two different errors to give the user more information about the result
	// There can be more than 1 values in the matrixResults but we are defining the step
	// parameter as the interval itself, hence there can only be one value.
	// This logic should be changed, once we work onto the aggregation functions.
	if len(resultMatrix) == 0 {
		r.Log.Info("No values in query result")
		return "", nil, fmt.Errorf("no values in query result")
	} else if len(resultMatrix) > 1 {
		r.Log.Info("Too many values in the query result")
		return "", nil, fmt.Errorf("too many values in the query result")
	}
	value := resultMatrix[0].Values[0].Value.String()
	b, err := resultMatrix[0].Values[0].Value.MarshalJSON()
	if err != nil {
		return "", nil, err
	}
	return value, b, nil
}

func getResultForVector(result model.Value, r *KeptnPrometheusProvider) (string, []byte, error) {
	// check if we can cast the result to a vector
	resultVector, ok := result.(model.Vector)
	if !ok {
		return "", nil, fmt.Errorf("could not cast result")
	}
	// We are only allowed to return one value, if not the query may be malformed
	// we are using two different errors to give the user more information about the result
	if len(resultVector) == 0 {
		r.Log.Info("No values in query result")
		return "", nil, fmt.Errorf("no values in query result")
	} else if len(resultVector) > 1 {
		r.Log.Info("Too many values in the query result")
		return "", nil, fmt.Errorf("too many values in the query result")
	}
	value := resultVector[0].Value.String()
	b, err := resultVector[0].Value.MarshalJSON()
	if err != nil {
		return "", nil, err
	}
	return value, b, nil
}
