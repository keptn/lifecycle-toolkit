package analysis

import (
	"testing"

	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/analysis/types"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestTargetEvaluator_Evaluate(t *testing.T) {
	compValue15 := resource.NewQuantity(15, resource.DecimalSI)
	compValue20 := resource.NewQuantity(20, resource.DecimalSI)
	tests := []struct {
		name string
		val  float64
		t    metricsapi.Target
		want types.TargetResult
	}{
		{
			name: "failure nor warning target set",
			val:  10.0,
			t:    metricsapi.Target{},
			want: types.TargetResult{
				Warning: false,
				Pass:    true,
			},
		},
		{
			name: "failure scenario",
			val:  10.0,
			t: metricsapi.Target{
				Failure: &metricsapi.Operator{
					LessThan: &metricsapi.OperatorValue{
						FixedValue: *compValue15,
					},
				},
			},
			want: types.TargetResult{
				FailResult: types.OperatorResult{
					Operator: metricsapi.Operator{
						LessThan: &metricsapi.OperatorValue{
							FixedValue: *compValue15,
						},
					},
					Fulfilled: true,
				},
				Warning: false,
				Pass:    false,
			},
		},
		{
			name: "warning scenario",
			val:  17.0,
			t: metricsapi.Target{
				Failure: &metricsapi.Operator{
					LessThan: &metricsapi.OperatorValue{
						FixedValue: *compValue15,
					},
				},
				Warning: &metricsapi.Operator{
					LessThan: &metricsapi.OperatorValue{
						FixedValue: *compValue20,
					},
				},
			},
			want: types.TargetResult{
				FailResult: types.OperatorResult{
					Operator: metricsapi.Operator{
						LessThan: &metricsapi.OperatorValue{
							FixedValue: *compValue15,
						},
					},
					Fulfilled: false,
				},
				WarnResult: types.OperatorResult{
					Operator: metricsapi.Operator{
						LessThan: &metricsapi.OperatorValue{
							FixedValue: *compValue20,
						},
					},
					Fulfilled: true,
				},
				Warning: true,
				Pass:    false,
			},
		},
		{
			name: "pass scenario",
			val:  27.0,
			t: metricsapi.Target{
				Failure: &metricsapi.Operator{
					LessThan: &metricsapi.OperatorValue{
						FixedValue: *compValue15,
					},
				},
				Warning: &metricsapi.Operator{
					LessThan: &metricsapi.OperatorValue{
						FixedValue: *compValue20,
					},
				},
			},
			want: types.TargetResult{
				FailResult: types.OperatorResult{
					Operator: metricsapi.Operator{
						LessThan: &metricsapi.OperatorValue{
							FixedValue: *compValue15,
						},
					},
					Fulfilled: false,
				},
				WarnResult: types.OperatorResult{
					Operator: metricsapi.Operator{
						LessThan: &metricsapi.OperatorValue{
							FixedValue: *compValue20,
						},
					},
					Fulfilled: false,
				},
				Warning: false,
				Pass:    true,
			},
		},
		{
			name: "pass scenario - only failed defined",
			val:  17.0,
			t: metricsapi.Target{
				Failure: &metricsapi.Operator{
					LessThan: &metricsapi.OperatorValue{
						FixedValue: *compValue15,
					},
				},
			},
			want: types.TargetResult{
				FailResult: types.OperatorResult{
					Operator: metricsapi.Operator{
						LessThan: &metricsapi.OperatorValue{
							FixedValue: *compValue15,
						},
					},
					Fulfilled: false,
				},
				Warning: false,
				Pass:    true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			te := NewTargetEvaluator(&OperatorEvaluator{})
			require.Equal(t, tt.want, te.Evaluate(tt.val, &tt.t))
		})
	}
}
