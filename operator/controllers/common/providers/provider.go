package providers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-logr/logr"
	klcv1alpha2 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const DynatraceProviderName = "dynatrace"
const PrometheusProviderName = "prometheus"
const KeptnMetricProviderName = "keptn-metric"

// KeptnSLIProvider is the interface that describes the operations that an SLI provider must implement
type KeptnSLIProvider interface {
	EvaluateQuery(ctx context.Context, objective klcv1alpha2.Objective, provider klcv1alpha2.KeptnEvaluationProvider) (string, []byte, error)
}

// NewProvider is a factory method that chooses the right implementation of KeptnSLIProvider
func NewProvider(provider string, log logr.Logger, k8sClient client.Client, opts ...func(p KeptnSLIProvider)) (KeptnSLIProvider, error) {
	switch strings.ToLower(provider) {
	case PrometheusProviderName:
		return &KeptnPrometheusProvider{
			httpClient: http.Client{},
			Log:        log,
		}, nil
	case DynatraceProviderName:
		return &KeptnDynatraceProvider{
			httpClient: http.Client{},
			Log:        log,
			k8sClient:  k8sClient,
		}, nil
	case KeptnMetricProviderName:
		return &KeptnMetricProvider{
			Log:       log,
			k8sClient: k8sClient,
		}, nil
	default:
		return nil, fmt.Errorf("provider %s not supported", provider)
	}
}
