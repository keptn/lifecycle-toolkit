package prometheus

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	promapi "github.com/prometheus/client_golang/api"
	prometheus "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var errCouldNotCast = fmt.Errorf("could not cast result")
var errNoValues = fmt.Errorf("no values in query result")
var errTooManyValues = fmt.Errorf("too many values in query result")

type KeptnPrometheusProvider struct {
	Log       logr.Logger
	K8sClient client.Client
}

func (r *KeptnPrometheusProvider) FetchAnalysisValue(ctx context.Context, query string, analysis metricsapi.Analysis, provider *metricsapi.KeptnMetricsProvider) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	api, err := r.setupApi(ctx, *provider)
	if err != nil {
		return "", err
	}

	r.Log.Info(fmt.Sprintf(
		"Running query: /api/v1/query_range?query=%s&start=%d&end=%d",
		query,
		analysis.GetFrom().Unix(), analysis.GetTo().Unix(),
	))
	queryRange := prometheus.Range{
		Start: analysis.GetFrom(),
		End:   analysis.GetTo(),
		Step:  time.Minute,
	}
	result, warnings, err := api.QueryRange(
		ctx,
		query,
		queryRange,
		[]prometheus.Option{}...,
	)

	if err != nil {
		return "", err
	}
	if len(warnings) != 0 {
		r.Log.Info("Prometheus API returned warnings: " + warnings[0])
	}
	res, _, err := getResultForMatrix(result)
	return res, err
}

// EvaluateQuery fetches the SLI values from prometheus provider
func (r *KeptnPrometheusProvider) EvaluateQuery(ctx context.Context, metric metricsapi.KeptnMetric, provider metricsapi.KeptnMetricsProvider) (string, []byte, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	api, err := r.setupApi(ctx, provider)

	if err != nil {
		return "", nil, err
	}

	if metric.Spec.Range != nil {
		result, warnings, err := evaluateQueryWithRange(ctx, metric, r, api)
		if err != nil {
			return "", nil, err
		}
		if len(warnings) != 0 {
			r.Log.Info("Prometheus API returned warnings: " + warnings[0])
		}
		return getResultForMatrix(result)
	} else {
		result, warnings, err := evaluateQueryWithoutRange(ctx, metric, r, api)
		if err != nil {
			return "", nil, err
		}
		if len(warnings) != 0 {
			r.Log.Info("Prometheus API returned warnings: " + warnings[0])
		}
		return getResultForVector(result)
	}
}

// EvaluateQueryForStep fetches the metric values from prometheus provider
func (r *KeptnPrometheusProvider) EvaluateQueryForStep(ctx context.Context, metric metricsapi.KeptnMetric, provider metricsapi.KeptnMetricsProvider) ([]string, []byte, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	api, err := r.setupApi(ctx, provider)
	if err != nil {
		return nil, nil, err
	}

	result, warnings, err := evaluateQueryWithRange(ctx, metric, r, api)
	if err != nil {
		return nil, nil, err
	}
	if len(warnings) != 0 {
		r.Log.Info("Prometheus API returned warnings: " + warnings[0])
	}
	return getResultForStepMatrix(result)
}

func (r *KeptnPrometheusProvider) setupApi(ctx context.Context, provider metricsapi.KeptnMetricsProvider) (prometheus.API, error) {
	rt, err := getRoundtripper(ctx, provider, r.K8sClient)
	if err != nil {
		return nil, err
	}

	pClient, err := promapi.NewClient(promapi.Config{Address: provider.Spec.TargetServer, RoundTripper: rt})
	if err != nil {
		return nil, err
	}
	return prometheus.NewAPI(pClient), nil
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

func getResultForMatrix(result model.Value) (string, []byte, error) {
	// check if we can cast the result to a matrix
	resultMatrix, ok := result.(model.Matrix)
	if !ok {
		return "", nil, errCouldNotCast
	}
	// We are only allowed to return one value, if not the query may be malformed
	// we are using two different errors to give the user more information about the result
	// There can be more than 1 values in the matrixResults but we are defining the step
	// parameter as the interval itself, hence there can only be one value.
	// This logic should be changed, once we work onto the aggregation functions.
	if len(resultMatrix) == 0 {
		return "", nil, errNoValues
	} else if len(resultMatrix) > 1 {
		return "", nil, errTooManyValues
	}
	value := resultMatrix[0].Values[0].Value.String()
	b, err := resultMatrix[0].Values[0].Value.MarshalJSON()
	if err != nil {
		return "", nil, err
	}
	return value, b, nil
}

func getResultForVector(result model.Value) (string, []byte, error) {
	// check if we can cast the result to a vector
	resultVector, ok := result.(model.Vector)
	if !ok {
		return "", nil, errCouldNotCast
	}
	// We are only allowed to return one value, if not the query may be malformed
	// we are using two different errors to give the user more information about the result
	if len(resultVector) == 0 {
		return "", nil, errNoValues
	} else if len(resultVector) > 1 {
		return "", nil, errTooManyValues
	}
	value := resultVector[0].Value.String()
	b, err := resultVector[0].Value.MarshalJSON()
	if err != nil {
		return "", nil, err
	}
	return value, b, nil
}

func getResultForStepMatrix(result model.Value) ([]string, []byte, error) {
	// check if we can cast the result to a matrix
	resultMatrix, ok := result.(model.Matrix)
	if !ok {
		return nil, nil, errCouldNotCast
	}

	if len(resultMatrix) == 0 {
		return nil, nil, errNoValues
	} else if len(resultMatrix) > 1 {
		return nil, nil, errTooManyValues
	}

	resultSlice := make([]string, len(resultMatrix[0].Values))
	for i, value := range resultMatrix[0].Values {
		resultSlice[i] = value.Value.String()
	}

	b, err := json.Marshal(resultSlice)
	if err != nil {
		return nil, nil, err
	}

	return resultSlice, b, nil
}
