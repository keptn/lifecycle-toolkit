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
func (r *KeptnPrometheusProvider) EvaluateQuery(ctx context.Context, metric metricsapi.KeptnMetric, provider metricsapi.KeptnMetricsProvider) (string, []byte, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	client, err := promapi.NewClient(promapi.Config{Address: provider.Spec.TargetServer, Client: &r.HttpClient})
	if err != nil {
		return "", nil, err
	}

	api := prometheus.NewAPI(client)
	if metric.Spec.Range != nil {
		result, warnings, err := evaluateQueryWithRange(ctx, metric, r, api)
		if err != nil {
			return "", nil, err
		}
		if len(warnings) != 0 {
			r.Log.Info("Prometheus API returned warnings: " + warnings[0])
		}
		return getResultForMatrix(result, r)
	} else {
		result, warnings, err := evaluateQueryWithoutRange(ctx, metric, r, api)
		if err != nil {
			return "", nil, err
		}
		if len(warnings) != 0 {
			r.Log.Info("Prometheus API returned warnings: " + warnings[0])
		}
		return getResultForVector(result, r)
	}
}

// EvaluateQueryForStep fetches the SLI values from prometheus provider
func (r *KeptnPrometheusProvider) EvaluateQueryForStep(ctx context.Context, metric metricsapi.KeptnMetric, provider metricsapi.KeptnMetricsProvider) ([]string, []byte, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	client, err := promapi.NewClient(promapi.Config{Address: provider.Spec.TargetServer, Client: &r.HttpClient})
	if err != nil {
		return nil, nil, err
	}

	api := prometheus.NewAPI(client)
	result, warnings, err := evaluateQueryWithRange(ctx, metric, r, api)
	if err != nil {
		return nil, nil, err
	}
	if len(warnings) != 0 {
		r.Log.Info("Prometheus API returned warnings: " + warnings[0])
	}
	return getResultForStepMatrix(result, r)
}

func evaluateQueryWithRange(ctx context.Context, metric metricsapi.KeptnMetric, r *KeptnPrometheusProvider, api prometheus.API) (model.Value, prometheus.Warnings, error) {
	queryTime := time.Now().UTC()
	// Get the duration
	queryInterval, err := time.ParseDuration(metric.Spec.Range.Interval)
	if err != nil {
		return nil, nil, err
	}
	var stepInterval time.Duration
	if metric.Spec.Range.Step != "" {
		stepTime := metric.Spec.Range.Step
		stepInterval, err = time.ParseDuration(stepTime)
		if err != nil {
			return nil, nil, err
		}
	} else {
		stepInterval = queryInterval
	}
	// Convert type Duration to type Time
	startTime := queryTime.Add(-queryInterval).UTC()
	r.Log.Info(fmt.Sprintf(
		"Running query: /api/v1/query_range?query=%s&start=%d&end=%d&step=%v",
		metric.Spec.Query,
		startTime.Unix(), queryTime.Unix(),
		stepInterval,
	))
	queryRange := prometheus.Range{
		Start: startTime,
		End:   queryTime,
		Step:  stepInterval,
	}
	result, warnings, err := api.QueryRange(
		ctx,
		metric.Spec.Query,
		queryRange,
		[]prometheus.Option{}...,
	)
	if err != nil {
		return nil, nil, err
	}
	return result, warnings, nil
}

func evaluateQueryWithoutRange(ctx context.Context, metric metricsapi.KeptnMetric, r *KeptnPrometheusProvider, api prometheus.API) (model.Value, prometheus.Warnings, error) {
	queryTime := time.Now().UTC()
	r.Log.Info(fmt.Sprintf(
		"Running query: /api/v1/query?query=%s&time=%d",
		metric.Spec.Query,
		queryTime.Unix(),
	))
	result, warnings, err := api.Query(
		ctx,
		metric.Spec.Query,
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

func getResultForStepMatrix(result model.Value, r *KeptnPrometheusProvider) ([]string, []byte, error) {
	// check if we can cast the result to a matrix
	resultMatrix, ok := result.(model.Matrix)
	if !ok {
		return nil, nil, fmt.Errorf("could not cast result")
	}

	if len(resultMatrix) == 0 {
		r.Log.Info("No values in query result")
		return nil, nil, fmt.Errorf("no values in query result")
	} else if len(resultMatrix) > 1 {
		r.Log.Info("Too many values in the query result")
		return nil, nil, fmt.Errorf("too many values in the query result")
	}
	var resultSlice []string
	for i, value := range resultMatrix[0].Values {
		resultSlice[i] = value.Value.String()
	}	
	b, err := resultMatrix[0].MarshalJSON()
	if err != nil {
		return nil, nil, err
	}
	return resultSlice, b, nil
}
