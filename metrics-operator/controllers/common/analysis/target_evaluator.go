package analysis

import (
	"github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/analysis/types"
)

type TargetEvaluator struct {
	OperatorEvaluator IOperatorEvaluator
}

func NewTargetEvaluator(o IOperatorEvaluator) TargetEvaluator {
	return TargetEvaluator{
		OperatorEvaluator: o,
	}
}

func (te *TargetEvaluator) Evaluate(val float64, t v1alpha3.Target) types.TargetResult {
	result := types.TargetResult{
		Warning: false,
		Pass:    false,
	}

	// check 'failure'  criteria
	if t.Failure != nil {
		result.FailResult = te.OperatorEvaluator.Evaluate(val, t.Failure)

		// if failure criteria are met, we can return without checking warning criteria
		if result.IsFail() {
			return result
		}
	}

	// check 'warning' criteria
	if t.Warning != nil {
		result.WarnResult = te.OperatorEvaluator.Evaluate(val, t.Warning)

		// if warning criteria are met, we can return warning
		if result.IsWarn() {
			result.Warning = true
			return result
		}
	}

	// if failure and warning criteria are not met, we pass
	result.Pass = true
	return result
}
