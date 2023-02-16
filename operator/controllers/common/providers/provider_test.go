package providers

import (
	"testing"

	"github.com/go-logr/logr"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/common/providers/keptnmetric"
	"github.com/stretchr/testify/require"
)

func TestFactory(t *testing.T) {
	tests := []struct {
		name     string
		provider interface{}
		err      bool
	}{
		{
			name:     KeptnMetricProviderName,
			provider: &keptnmetric.KeptnMetricProvider{},
			err:      false,
		},
		{
			name:     "invalid",
			provider: nil,
			err:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, e := NewProvider(tt.name, logr.Logger{}, nil)
			require.IsType(t, tt.provider, p)
			if tt.err {
				require.NotNil(t, e)
			}
		})

	}
}
