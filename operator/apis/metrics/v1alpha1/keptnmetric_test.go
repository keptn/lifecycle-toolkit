package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKeptnMetric_IsStatusSet(t *testing.T) {
	tests := []struct {
		name   string
		metric KeptnMetric
		want   bool
	}{
		{
			name: "results set",
			metric: KeptnMetric{
				Status: KeptnMetricStatus{
					Value: "val",
				},
			},
			want: true,
		},
		{
			name: "value empty",
			metric: KeptnMetric{
				Status: KeptnMetricStatus{
					Value: "",
				},
			},
			want: false,
		},
		{
			name: "value not set",
			metric: KeptnMetric{
				Status: KeptnMetricStatus{},
			},
			want: false,
		},
		{
			name:   "status not set",
			metric: KeptnMetric{},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.metric.IsStatusSet())
		})
	}
}
