package analysis

import (
	"testing"

	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/analysis/types"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestOperatorEvaluator_Evaluate(t *testing.T) {
	compValue := resource.NewQuantity(15, resource.DecimalSI)
	compValue2 := resource.NewQuantity(25, resource.DecimalSI)
	tests := []struct {
		name string
		val  float64
		o    metricsapi.Operator
		want types.OperatorResult
	}{
		{
			name: "less than - fulfilled",
			val:  10.0,
			o: metricsapi.Operator{
				LessThan: &metricsapi.OperatorValue{
					FixedValue: *compValue,
				},
			},
			want: types.OperatorResult{
				Operator: metricsapi.Operator{
					LessThan: &metricsapi.OperatorValue{
						FixedValue: *compValue,
					},
				},
				Fulfilled: true,
			},
		},
		{
			name: "less than - not fulfilled",
			val:  16.0,
			o: metricsapi.Operator{
				LessThan: &metricsapi.OperatorValue{
					FixedValue: *compValue,
				},
			},
			want: types.OperatorResult{
				Operator: metricsapi.Operator{
					LessThan: &metricsapi.OperatorValue{
						FixedValue: *compValue,
					},
				},
				Fulfilled: false,
			},
		},
		{
			name: "less than equal - fulfilled",
			val:  15.0,
			o: metricsapi.Operator{
				LessThanOrEqual: &metricsapi.OperatorValue{
					FixedValue: *compValue,
				},
			},
			want: types.OperatorResult{
				Operator: metricsapi.Operator{
					LessThanOrEqual: &metricsapi.OperatorValue{
						FixedValue: *compValue,
					},
				},
				Fulfilled: true,
			},
		},
		{
			name: "less than equal - not fulfilled",
			val:  16.0,
			o: metricsapi.Operator{
				LessThanOrEqual: &metricsapi.OperatorValue{
					FixedValue: *compValue,
				},
			},
			want: types.OperatorResult{
				Operator: metricsapi.Operator{
					LessThanOrEqual: &metricsapi.OperatorValue{
						FixedValue: *compValue,
					},
				},
				Fulfilled: false,
			},
		},
		{
			name: "equal - not fulfilled",
			val:  16.0,
			o: metricsapi.Operator{
				EqualTo: &metricsapi.OperatorValue{
					FixedValue: *compValue,
				},
			},
			want: types.OperatorResult{
				Operator: metricsapi.Operator{
					EqualTo: &metricsapi.OperatorValue{
						FixedValue: *compValue,
					},
				},
				Fulfilled: false,
			},
		},
		{
			name: "equal - fulfilled",
			val:  15.0,
			o: metricsapi.Operator{
				EqualTo: &metricsapi.OperatorValue{
					FixedValue: *compValue,
				},
			},
			want: types.OperatorResult{
				Operator: metricsapi.Operator{
					EqualTo: &metricsapi.OperatorValue{
						FixedValue: *compValue,
					},
				},
				Fulfilled: true,
			},
		},
		{
			name: "greater than - fulfilled",
			val:  16.0,
			o: metricsapi.Operator{
				GreaterThan: &metricsapi.OperatorValue{
					FixedValue: *compValue,
				},
			},
			want: types.OperatorResult{
				Operator: metricsapi.Operator{
					GreaterThan: &metricsapi.OperatorValue{
						FixedValue: *compValue,
					},
				},
				Fulfilled: true,
			},
		},
		{
			name: "greater than - not fulfilled",
			val:  10.0,
			o: metricsapi.Operator{
				GreaterThan: &metricsapi.OperatorValue{
					FixedValue: *compValue,
				},
			},
			want: types.OperatorResult{
				Operator: metricsapi.Operator{
					GreaterThan: &metricsapi.OperatorValue{
						FixedValue: *compValue,
					},
				},
				Fulfilled: false,
			},
		},
		{
			name: "greater than equal - fulfilled",
			val:  15.0,
			o: metricsapi.Operator{
				GreaterThanOrEqual: &metricsapi.OperatorValue{
					FixedValue: *compValue,
				},
			},
			want: types.OperatorResult{
				Operator: metricsapi.Operator{
					GreaterThanOrEqual: &metricsapi.OperatorValue{
						FixedValue: *compValue,
					},
				},
				Fulfilled: true,
			},
		},
		{
			name: "greater than equal - not fulfilled",
			val:  10.0,
			o: metricsapi.Operator{
				GreaterThanOrEqual: &metricsapi.OperatorValue{
					FixedValue: *compValue,
				},
			},
			want: types.OperatorResult{
				Operator: metricsapi.Operator{
					GreaterThanOrEqual: &metricsapi.OperatorValue{
						FixedValue: *compValue,
					},
				},
				Fulfilled: false,
			},
		},
		{
			name: "in range - fulfilled",
			val:  20.0,
			o: metricsapi.Operator{
				InRange: &metricsapi.RangeValue{
					LowBound:  *compValue,
					HighBound: *compValue2,
				},
			},
			want: types.OperatorResult{
				Operator: metricsapi.Operator{
					InRange: &metricsapi.RangeValue{
						LowBound:  *compValue,
						HighBound: *compValue2,
					},
				},
				Fulfilled: true,
			},
		},
		{
			name: "in range - not fulfilled",
			val:  30.0,
			o: metricsapi.Operator{
				InRange: &metricsapi.RangeValue{
					LowBound:  *compValue,
					HighBound: *compValue2,
				},
			},
			want: types.OperatorResult{
				Operator: metricsapi.Operator{
					InRange: &metricsapi.RangeValue{
						LowBound:  *compValue,
						HighBound: *compValue2,
					},
				},
				Fulfilled: false,
			},
		},
		{
			name: "not in range - fulfilled",
			val:  30.0,
			o: metricsapi.Operator{
				NotInRange: &metricsapi.RangeValue{
					LowBound:  *compValue,
					HighBound: *compValue2,
				},
			},
			want: types.OperatorResult{
				Operator: metricsapi.Operator{
					NotInRange: &metricsapi.RangeValue{
						LowBound:  *compValue,
						HighBound: *compValue2,
					},
				},
				Fulfilled: true,
			},
		},
		{
			name: "not in range - not fulfilled",
			val:  20.0,
			o: metricsapi.Operator{
				NotInRange: &metricsapi.RangeValue{
					LowBound:  *compValue,
					HighBound: *compValue2,
				},
			},
			want: types.OperatorResult{
				Operator: metricsapi.Operator{
					NotInRange: &metricsapi.RangeValue{
						LowBound:  *compValue,
						HighBound: *compValue2,
					},
				},
				Fulfilled: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oe := &OperatorEvaluator{}
			require.Equal(t, tt.want, oe.Evaluate(tt.val, &tt.o))
		})
	}
}
