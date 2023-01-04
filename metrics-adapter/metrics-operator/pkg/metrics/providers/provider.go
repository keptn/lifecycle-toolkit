package providers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-logr/logr"
	metricsv1alpha1 "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// KeptnSLIProvider is the interface that describes the operations that an SLI provider must implement
type MetricsProvider interface {
	EvaluateQuery(ctx context.Context, metric metricsv1alpha1.Metric, provider metricsv1alpha1.Provider) (string, error)
}

// NewProvider is a factory method that chooses the right implementation of MetricsProvider
func NewProvider(provider string, log logr.Logger, k8sClient client.Client) (MetricsProvider, error) {
	switch strings.ToLower(provider) {
	case "prometheus":
		return &KeptnPrometheusProvider{
			httpClient: http.Client{},
			Log:        log,
		}, nil
	case "dynatrace":
		return &KeptnDynatraceProvider{
			httpClient: http.Client{},
			Log:        log,
			k8sClient:  k8sClient,
		}, nil
	default:
		return nil, fmt.Errorf("provider %s not supported", provider)
	}
}
