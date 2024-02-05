package dummy

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-logr/logr"
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1beta1"
)

type KeptnDummyProvider struct {
	Log        logr.Logger
	HttpClient http.Client
}

func (d *KeptnDummyProvider) FetchAnalysisValue(ctx context.Context, query string, analysis metricsapi.Analysis, provider *metricsapi.KeptnMetricsProvider) (string, error) {
	fromTime := strconv.Itoa(int(analysis.GetFrom().Unix()))
	toTime := strconv.Itoa(int(analysis.GetFrom().Unix()))
	return fmt.Sprintf("dummy provider FetchAnalysisValue was called with query %s from %q to %q", query, fromTime, toTime), nil
}

func (d *KeptnDummyProvider) EvaluateQuery(ctx context.Context, metric metricsapi.KeptnMetric, provider metricsapi.KeptnMetricsProvider) (string, []byte, error) {
	return fmt.Sprintf("dummy provider EvaluateQuery was called with query %s", metric.Spec.Query), nil, nil
}

func (d *KeptnDummyProvider) EvaluateQueryForStep(ctx context.Context, metric metricsapi.KeptnMetric, provider metricsapi.KeptnMetricsProvider) ([]string, []byte, error) {
	fromTime, toTime, stepInterval, err := getTimeRangeForStep(metric)
	if err != nil {
		return nil, nil, err
	}
	result := fmt.Sprintf("dummy provider EvaluateQueryForStep was called with query %s from %q to %q at an interval %q", metric.Spec.Query, fromTime, toTime, stepInterval)
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
