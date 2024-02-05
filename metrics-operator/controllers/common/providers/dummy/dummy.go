package dummy

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-logr/logr"
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1beta1"
)

type KeptnDummyProvider struct {
	Log        logr.Logger
	HttpClient http.Client
}

func (d *KeptnDummyProvider) FetchAnalysisValue(ctx context.Context, query string, analysis metricsapi.Analysis, provider *metricsapi.KeptnMetricsProvider) (string, error) {
	return fmt.Sprintf("dummy provider FetchAnalysisValue was called with query %s", query), nil
}

func (d *KeptnDummyProvider) EvaluateQuery(ctx context.Context, metric metricsapi.KeptnMetric, provider metricsapi.KeptnMetricsProvider) (string, []byte, error) {
	return fmt.Sprintf("dummy provider EvaluateQuery was called with query %s", metric.Spec.Query), nil, nil
}

func (d *KeptnDummyProvider) EvaluateQueryForStep(ctx context.Context, metric metricsapi.KeptnMetric, provider metricsapi.KeptnMetricsProvider) ([]string, []byte, error) {
	result := fmt.Sprintf("dummy provider EvaluateQueryForStep was called with query %s", metric.Spec.Query)
	return []string{result}, nil, nil
}
