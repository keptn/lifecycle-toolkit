package providers

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-logr/logr"
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/providers/datadog"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/providers/dynatrace"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/providers/elastic"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/providers/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// KeptnSLIProvider is the interface that describes the operations that an SLI provider must implement
//
//go:generate moq -pkg fake -skip-ensure -out ./fake/provider_mock.go . KeptnSLIProvider
type KeptnSLIProvider interface {
	EvaluateQuery(ctx context.Context, metric metricsapi.KeptnMetric, provider metricsapi.KeptnMetricsProvider) (string, []byte, error)
	EvaluateQueryForStep(ctx context.Context, metric metricsapi.KeptnMetric, provider metricsapi.KeptnMetricsProvider) ([]string, []byte, error)
	FetchAnalysisValue(ctx context.Context, query string, spec metricsapi.Analysis, provider *metricsapi.KeptnMetricsProvider) (string, error)
}

type ProviderFactory func(provider *metricsapi.KeptnMetricsProvider, log logr.Logger, k8sClient client.Client) (KeptnSLIProvider, error)

// NewProvider is a factory method that chooses the right implementation of KeptnSLIProvider
func NewProvider(provider *metricsapi.KeptnMetricsProvider, log logr.Logger, k8sClient client.Client) (KeptnSLIProvider, error) {

	switch strings.ToLower(provider.Spec.Type) {
	case PrometheusProviderType, ThanosProviderType, CortexProviderType:
		return prometheus.NewPrometheusProvider(log, k8sClient), nil
	case DynatraceProviderType:
		return &dynatrace.KeptnDynatraceProvider{
			HttpClient: http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: provider.Spec.InsecureSkipTlsVerify,
					},
				},
			},
			Log:       log,
			K8sClient: k8sClient,
		}, nil
	case DynatraceDQLProviderType:
		return dynatrace.NewKeptnDynatraceDQLProvider(
			k8sClient,
			dynatrace.WithLogger(log),
			dynatrace.WithInsecureSkipTlsVerify(provider.Spec.InsecureSkipTlsVerify),
		), nil
	case DataDogProviderType:
		return &datadog.KeptnDataDogProvider{
			Log: log,
			HttpClient: http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: provider.Spec.InsecureSkipTlsVerify,
					},
				},
			},
			K8sClient: k8sClient,
		}, nil
	case ElasticProviderType:
		es, err := elastic.GetElasticClient(*provider)
		if err != nil {
			return nil, err
		}
		return &elastic.KeptnElasticProvider{
			Log:       log,
			K8sClient: k8sClient,
			Elastic:   es,
		}, nil
	default:
		return nil, fmt.Errorf("provider %s not supported", provider.Spec.Type)
	}
}
