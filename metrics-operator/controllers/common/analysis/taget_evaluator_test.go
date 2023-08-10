package analysis

import (
	"testing"

	"github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestTargetEvaluator_Evaluate(t *testing.T) {
	compValue15 := resource.NewQuantity(15, resource.DecimalSI)
	compValue20 := resource.NewQuantity(20, resource.DecimalSI)
	tests := []struct {
		name string
		val  float64
		t    v1alpha3.Target
		want v1alpha3.TargetResult
	}{
		{
			name: "failure scenario",
			val:  10.0,
			t: v1alpha3.Target{
				Failure: &v1alpha3.Operator{
					LessThan: &v1alpha3.OperatorValue{
						FixedValue: *compValue15,
					},
				},
			},
			want: v1alpha3.TargetResult{
				FailureResult: v1alpha3.OperatorResult{
					Operator: v1alpha3.Operator{
						LessThan: &v1alpha3.OperatorValue{
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
			t: v1alpha3.Target{
				Failure: &v1alpha3.Operator{
					LessThan: &v1alpha3.OperatorValue{
						FixedValue: *compValue15,
					},
				},
				Warning: &v1alpha3.Operator{
					LessThan: &v1alpha3.OperatorValue{
						FixedValue: *compValue20,
					},
				},
			},
			want: v1alpha3.TargetResult{
				FailureResult: v1alpha3.OperatorResult{
					Operator: v1alpha3.Operator{
						LessThan: &v1alpha3.OperatorValue{
							FixedValue: *compValue15,
						},
					},
					Fulfilled: false,
				},
				WarningResult: v1alpha3.OperatorResult{
					Operator: v1alpha3.Operator{
						LessThan: &v1alpha3.OperatorValue{
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
			t: v1alpha3.Target{
				Failure: &v1alpha3.Operator{
					LessThan: &v1alpha3.OperatorValue{
						FixedValue: *compValue15,
					},
				},
				Warning: &v1alpha3.Operator{
					LessThan: &v1alpha3.OperatorValue{
						FixedValue: *compValue20,
					},
				},
			},
			want: v1alpha3.TargetResult{
				FailureResult: v1alpha3.OperatorResult{
					Operator: v1alpha3.Operator{
						LessThan: &v1alpha3.OperatorValue{
							FixedValue: *compValue15,
						},
					},
					Fulfilled: false,
				},
				WarningResult: v1alpha3.OperatorResult{
					Operator: v1alpha3.Operator{
						LessThan: &v1alpha3.OperatorValue{
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
			t: v1alpha3.Target{
				Failure: &v1alpha3.Operator{
					LessThan: &v1alpha3.OperatorValue{
						FixedValue: *compValue15,
					},
				},
			},
			want: v1alpha3.TargetResult{
				FailureResult: v1alpha3.OperatorResult{
					Operator: v1alpha3.Operator{
						LessThan: &v1alpha3.OperatorValue{
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
			require.Equal(t, tt.want, te.Evaluate(tt.val, tt.t))
		})
	}
}
