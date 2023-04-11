package providers

import (
	"testing"

	"github.com/go-logr/logr"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/providers/datadog"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/providers/dynatrace"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/providers/prometheus"
	"github.com/stretchr/testify/require"
)

func TestFactory(t *testing.T) {
	tests := []struct {
		name         string
		providerType string
		provider     interface{}
		err          bool
	}{
		{
			name:         PrometheusProviderType,
			providerType: "",
			provider:     &prometheus.KeptnPrometheusProvider{},
			err:          false,
		},
		{
			name:         DynatraceProviderType,
			providerType: "",
			provider:     &dynatrace.KeptnDynatraceProvider{},
			err:          false,
		},
		{
			name:         DataDogProviderType,
			providerType: "",
			provider:     &datadog.KeptnDataDogProvider{},
			err:          false,
		},
		{
			name:         PrometheusProviderType,
			providerType: PrometheusProviderType,
			provider:     &prometheus.KeptnPrometheusProvider{},
			err:          false,
		},
		{
			name:         DynatraceProviderType,
			providerType: DynatraceProviderType,
			provider:     &dynatrace.KeptnDynatraceProvider{},
			err:          false,
		},
		{
			name:         DataDogProviderType,
			providerType: DataDogProviderType,
			provider:     &datadog.KeptnDataDogProvider{},
			err:          false,
		},
		{
			name:         "invalid",
			providerType: "invalid",
			provider:     nil,
			err:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, e := NewProvider(tt.name, tt.providerType, logr.Logger{}, nil)
			require.IsType(t, tt.provider, p)
			if tt.err {
				require.NotNil(t, e)
			}
		})

	}
}
