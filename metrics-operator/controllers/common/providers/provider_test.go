package providers

import (
	"testing"

	"github.com/go-logr/logr"
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/fake"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/providers/datadog"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/providers/dynatrace"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/providers/elastic"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/providers/prometheus"
	"github.com/stretchr/testify/require"
)

func TestFactory(t *testing.T) {
	tests := []struct {
		metricsProvider metricsapi.KeptnMetricsProvider
		provider        interface{}
		err             bool
	}{
		{
			metricsProvider: metricsapi.KeptnMetricsProvider{
				Spec: metricsapi.KeptnMetricsProviderSpec{
					Type: PrometheusProviderType,
				},
			},
			provider: &prometheus.KeptnPrometheusProvider{},
			err:      false,
		},
		{
			metricsProvider: metricsapi.KeptnMetricsProvider{
				Spec: metricsapi.KeptnMetricsProviderSpec{
					Type: ThanosProviderType,
				},
			},
			provider: &prometheus.KeptnPrometheusProvider{},
			err:      false,
		},
		{
			metricsProvider: metricsapi.KeptnMetricsProvider{
				Spec: metricsapi.KeptnMetricsProviderSpec{
					Type: CortexProviderType,
				},
			},
			provider: &prometheus.KeptnPrometheusProvider{},
			err:      false,
		},
		{
			metricsProvider: metricsapi.KeptnMetricsProvider{
				Spec: metricsapi.KeptnMetricsProviderSpec{
					Type: ElasticProviderType,
				},
			},
			provider: &elastic.KeptnElasticProvider{},
			err:      false,
		},
		{
			metricsProvider: metricsapi.KeptnMetricsProvider{
				Spec: metricsapi.KeptnMetricsProviderSpec{
					Type: DynatraceProviderType,
				},
			},
			provider: &dynatrace.KeptnDynatraceProvider{},
			err:      false,
		},
		{
			metricsProvider: metricsapi.KeptnMetricsProvider{
				Spec: metricsapi.KeptnMetricsProviderSpec{
					Type: DynatraceDQLProviderType,
				},
			},
			provider: dynatrace.NewKeptnDynatraceDQLProvider(fake.NewClient()),
			err:      false,
		},
		{
			metricsProvider: metricsapi.KeptnMetricsProvider{
				Spec: metricsapi.KeptnMetricsProviderSpec{
					Type: DataDogProviderType,
				},
			},
			provider: &datadog.KeptnDataDogProvider{},
			err:      false,
		},
		{
			metricsProvider: metricsapi.KeptnMetricsProvider{
				Spec: metricsapi.KeptnMetricsProviderSpec{
					Type: "invalid",
				},
			},
			provider: nil,
			err:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.metricsProvider.Spec.Type, func(t *testing.T) {
			p, e := NewProvider(&tt.metricsProvider, logr.Logger{}, nil)
			require.IsType(t, tt.provider, p)
			if tt.err {
				require.NotNil(t, e)
			}
		})

	}
}
