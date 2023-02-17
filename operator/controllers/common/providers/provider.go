package providers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-logr/logr"
	klcv1alpha2 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/common/providers/datadog"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/common/providers/dynatrace"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/common/providers/keptnmetric"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/common/providers/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// KeptnSLIProvider is the interface that describes the operations that an SLI provider must implement
type KeptnSLIProvider interface {
	EvaluateQuery(ctx context.Context, objective klcv1alpha2.Objective, provider klcv1alpha2.KeptnEvaluationProvider) (string, []byte, error)
}

// NewProvider is a factory method that chooses the right implementation of KeptnSLIProvider
func NewProvider(provider string, log logr.Logger, k8sClient client.Client) (KeptnSLIProvider, error) {
	switch strings.ToLower(provider) {
	case PrometheusProviderName:
		return &prometheus.KeptnPrometheusProvider{
			HttpClient: http.Client{},
			Log:        log,
		}, nil
	case DynatraceProviderName:
		return &dynatrace.KeptnDynatraceProvider{
			HttpClient: http.Client{},
			Log:        log,
			K8sClient:  k8sClient,
		}, nil
	case DynatraceDQLProviderName:
		return dynatrace.NewKeptnDynatraceDQLProvider(
			k8sClient,
			dynatrace.WithLogger(log),
		), nil
	case KeptnMetricProviderName:
		return &keptnmetric.KeptnMetricProvider{
			Log:       log,
			K8sClient: k8sClient,
		}, nil
	case DataDogProviderName:
		return &datadog.KeptnDataDogProvider{
			Log:        log,
			HttpClient: http.Client{},
		}, nil
	default:
		return nil, fmt.Errorf("provider %s not supported", provider)
	}
}
