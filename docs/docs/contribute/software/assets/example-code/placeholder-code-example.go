package placeholder

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-logr/logr"
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1"
)

type KeptnPlaceholderProvider struct {
	Log        logr.Logger
	HttpClient http.Client
}

func (d *KeptnPlaceholderProvider) FetchAnalysisValue(ctx context.Context, query string, analysis metricsapi.Analysis, provider *metricsapi.KeptnMetricsProvider) (string, error) {
	return fmt.Sprintf("placeholder provider FetchAnalysisValue was called with query %s from %d to %d", query, analysis.GetFrom().Unix(), analysis.GetTo().Unix()), nil
}

func (d *KeptnPlaceholderProvider) EvaluateQuery(ctx context.Context, metric metricsapi.KeptnMetric, provider metricsapi.KeptnMetricsProvider) (string, []byte, error) {
	return fmt.Sprintf("placeholder provider EvaluateQuery was called with query %s", metric.Spec.Query), nil, nil
}

func (d *KeptnPlaceholderProvider) EvaluateQueryForStep(ctx context.Context, metric metricsapi.KeptnMetric, provider metricsapi.KeptnMetricsProvider) ([]string, []byte, error) {
	fromTime, toTime, stepInterval, err := getTimeRangeForStep(metric)
	if err != nil {
		return nil, nil, err
	}
	result := fmt.Sprintf("placeholder provider EvaluateQueryForStep was called with query %s from %d to %d at an interval %d", metric.Spec.Query, fromTime, toTime, stepInterval)
	return []string{result}, nil, nil
}

func getTimeRangeForStep(metric metricsapi.KeptnMetric) (int64, int64, int64, error) {
	intervalDuration, err := time.ParseDuration(metric.Spec.Range.Interval)
	if err != nil {
		return 0, 0, 0, err
	}
	stepDuration, err := time.ParseDuration(metric.Spec.Range.Step)
	if err != nil {
		return 0, 0, 0, err
	}
	return time.Now().Add(-intervalDuration).Unix(), time.Now().Unix(), stepDuration.Milliseconds(), nil
}
