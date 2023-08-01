package analysis

import "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"

type TargetEvaluator struct{}

func (te *TargetEvaluator) Evaluate(val float64, t v1alpha3.Target) v1alpha3.TargetResult {
	result := v1alpha3.TargetResult{
		Target:   t,
		Violated: false,
	}

	if t.EqualTo != nil && t.EqualTo.FixedValue != nil {
		result.Violated = !(*t.EqualTo.FixedValue == val)
	} else if t.LessThanOrEqual != nil && t.LessThanOrEqual.FixedValue != nil {
		result.Violated = !(val <= *t.LessThanOrEqual.FixedValue)
	} else if t.LessThan != nil && t.LessThan.FixedValue != nil {
		result.Violated = !(val < *t.LessThan.FixedValue)
	} else if t.GreaterThan != nil && t.GreaterThan.FixedValue != nil {
		result.Violated = !(val > *t.GreaterThan.FixedValue)
	} else if t.GreaterThanOrEqual != nil && t.GreaterThanOrEqual.FixedValue != nil {
		result.Violated = !(val >= *t.GreaterThanOrEqual.FixedValue)
	}

	return result
}
