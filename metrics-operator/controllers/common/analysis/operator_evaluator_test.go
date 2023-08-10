package analysis

import (
	"testing"

	"github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestOperatorEvaluator_Evaluate(t *testing.T) {
	compValue := resource.NewQuantity(15, resource.DecimalSI)
	tests := []struct {
		name string
		val  float64
		o    v1alpha3.Operator
		want v1alpha3.OperatorResult
	}{
		{
			name: "less than - fulfilled",
			val:  10.0,
			o: v1alpha3.Operator{
				LessThan: &v1alpha3.OperatorValue{
					FixedValue: *compValue,
				},
			},
			want: v1alpha3.OperatorResult{
				Operator: v1alpha3.Operator{
					LessThan: &v1alpha3.OperatorValue{
						FixedValue: *compValue,
					},
				},
				Fulfilled: true,
			},
		},
		{
			name: "less than - not fulfilled",
			val:  16.0,
			o: v1alpha3.Operator{
				LessThan: &v1alpha3.OperatorValue{
					FixedValue: *compValue,
				},
			},
			want: v1alpha3.OperatorResult{
				Operator: v1alpha3.Operator{
					LessThan: &v1alpha3.OperatorValue{
						FixedValue: *compValue,
					},
				},
				Fulfilled: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oe := &OperatorEvaluator{}
			require.Equal(t, tt.want, oe.Evaluate(tt.val, tt.o))
		})
	}
}
