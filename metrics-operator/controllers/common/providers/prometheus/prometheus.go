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

	queryTime := time.Now().UTC()
	r.Log.Info("Running query: /api/v1/query?query=" + metric.Spec.Query + "&time=" + queryTime.String())
	client, err := promapi.NewClient(promapi.Config{Address: provider.Spec.TargetServer, Client: &r.HttpClient})
	if err != nil {
		return "", nil, err
	}
	api := prometheus.NewAPI(client)
	var result model.Value
	var warnings prometheus.Warnings
	
	if metric.Spec.Range != nil {

		queryInterval, err := time.ParseDuration(metric.Spec.Range.Interval)
		if err != nil {
			return "", nil, err
		}
		endTime := queryTime.Add(queryInterval)

		queryRange := prometheus.Range{
			Start: queryTime,
			End: endTime,
		}

		r, w, err := api.QueryRange(
			ctx,
			metric.Spec.Query,
			queryRange,
			[]prometheus.Option{}...,
		)

		if err != nil {
			return "", nil, err
		}
		result = r
		warnings = w

	} else {

		r, w, err := api.Query(
			ctx,
			metric.Spec.Query,
			queryTime,
			[]prometheus.Option{}...,
		)
		if err != nil {
			return "", nil, err
		}
		result = r
		warnings = w
	}
	
	if len(warnings) != 0 {
		r.Log.Info("Prometheus API returned warnings: " + warnings[0])
	}

	// check if we can cast the result to a vector, it might be another data struct which we can't process
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
