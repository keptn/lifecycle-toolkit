package analysis

import (
	"github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/analysis/types"
)

type OperatorEvaluator struct{}

func (te *OperatorEvaluator) Evaluate(val float64, t v1alpha3.Operator) types.OperatorResult {
	result := types.OperatorResult{
		Operator:  t,
		Fulfilled: false,
	}

	if t.EqualTo != nil {
		result.Fulfilled = (val == t.EqualTo.GetFloatValue())
	} else if t.LessThanOrEqual != nil {
		result.Fulfilled = (val <= t.LessThanOrEqual.GetFloatValue())
	} else if t.LessThan != nil {
		result.Fulfilled = (val < t.LessThan.GetFloatValue())
	} else if t.GreaterThanOrEqual != nil {
		result.Fulfilled = (val >= t.GreaterThanOrEqual.GetFloatValue())
	} else if t.GreaterThan != nil {
		result.Fulfilled = (val > t.GreaterThan.GetFloatValue())
	}

	return result
}
