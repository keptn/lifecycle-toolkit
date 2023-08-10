package analysis

import (
	"github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
)

type OperatorEvaluator struct{}

func (te *OperatorEvaluator) Evaluate(val float64, t v1alpha3.Operator) v1alpha3.OperatorResult {
	result := v1alpha3.OperatorResult{
		Operator:  t,
		Fulfilled: false,
	}

	if t.EqualTo != nil {
		result.Fulfilled = (val == t.EqualTo.FixedValue.AsApproximateFloat64())
	} else if t.LessThanOrEqual != nil {
		result.Fulfilled = (val <= t.LessThanOrEqual.FixedValue.AsApproximateFloat64())
	} else if t.LessThan != nil {
		result.Fulfilled = (val < t.LessThan.FixedValue.AsApproximateFloat64())
	} else if t.GreaterThan != nil {
		result.Fulfilled = (val > t.GreaterThan.FixedValue.AsApproximateFloat64())
	} else if t.GreaterThanOrEqual != nil {
		result.Fulfilled = (val >= t.GreaterThanOrEqual.FixedValue.AsApproximateFloat64())
	}

	return result
}
