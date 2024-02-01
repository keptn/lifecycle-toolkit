package providers

import (
	"testing"

	"github.com/go-logr/logr"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/fake"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/providers/datadog"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/providers/dynatrace"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/providers/prometheus"
	"github.com/stretchr/testify/require"
)

func TestFactory(t *testing.T) {
	tests := []struct {
		providerType string
		provider     interface{}
		err          bool
	}{
		{
			providerType: PrometheusProviderType,
			provider:     &prometheus.KeptnPrometheusProvider{},
			err:          false,
		},
		{
			providerType: DynatraceProviderType,
			provider:     &dynatrace.KeptnDynatraceProvider{},
			err:          false,
		},
		{
			providerType: DynatraceDQLProviderType,
			provider:     dynatrace.NewKeptnDynatraceDQLProvider(fake.NewClient()),
			err:          false,
		},
		{
			providerType: DataDogProviderType,
			provider:     &datadog.KeptnDataDogProvider{},
			err:          false,
		},
		{
			providerType: "invalid",
			provider:     nil,
			err:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.providerType, func(t *testing.T) {
			p, e := NewProvider(tt.providerType, logr.Logger{}, nil)
			require.IsType(t, tt.provider, p)
			if tt.err {
				require.NotNil(t, e)
			}
		})

	}
}
