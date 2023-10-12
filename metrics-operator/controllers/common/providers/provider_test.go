package providers

import (
	"testing"

	"github.com/go-logr/logr"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/fake"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/providers/datadog"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/providers/dynatrace"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/providers/prometheus"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
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

func TestGetRequestInfo(t *testing.T) {
	req := ctrl.Request{
		NamespacedName: types.NamespacedName{
			Name:      "example",
			Namespace: "test-namespace",
		}}

	info := GetRequestInfo(req)
	expected := map[string]string{
		"name":      "example",
		"namespace": "test-namespace",
	}
	require.Equal(t, expected, info)
}
