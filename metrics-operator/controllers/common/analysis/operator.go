package analysis

import (
	"github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
)

type OperatorEvaluator struct{}

type OperatorResult struct {
	Operator   v1alpha3.Operator
	Fullfilled bool
}

func (te *OperatorEvaluator) Evaluate(val float64, t v1alpha3.Operator) OperatorResult {
	result := OperatorResult{
		Operator:   t,
		Fullfilled: false,
	}

	if t.EqualTo != nil {
		result.Fullfilled = (val == t.EqualTo.FixedValue.AsApproximateFloat64())
	} else if t.LessThanOrEqual != nil {
		result.Fullfilled = (val <= t.LessThanOrEqual.FixedValue.AsApproximateFloat64())
	} else if t.LessThan != nil {
		result.Fullfilled = (val < t.LessThan.FixedValue.AsApproximateFloat64())
	} else if t.GreaterThan != nil {
		result.Fullfilled = (val > t.GreaterThan.FixedValue.AsApproximateFloat64())
	} else if t.GreaterThanOrEqual != nil {
		result.Fullfilled = (val >= t.GreaterThanOrEqual.FixedValue.AsApproximateFloat64())
	}

	return result
}
