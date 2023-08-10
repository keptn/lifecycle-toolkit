package analysis

import "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"

type TargetEvaluator struct {
	OperatorEvaluator IOperatorEvaluator
}

type TargetResult struct {
	FailureResult OperatorResult
	WarningResult OperatorResult
	Warning       bool
	Pass          bool
}

func NewTargetEvaluator(o IOperatorEvaluator) TargetEvaluator {
	return TargetEvaluator{
		OperatorEvaluator: o,
	}
}

func (te *TargetEvaluator) Evaluate(val float64, t v1alpha3.Target) TargetResult {
	result := TargetResult{
		Warning: false,
		Pass:    false,
	}

	// check 'failure'  criteria
	result.FailureResult = te.OperatorEvaluator.Evaluate(val, *t.Failure)

	// if failure criteria are met, we can return without checking warning criteria
	if result.FailureResult.Fullfilled {
		return result
	}

	// check 'warning'  criteria
	result.WarningResult = te.OperatorEvaluator.Evaluate(val, *t.Warning)

	// if warning criteria are met, we can return warning
	if result.WarningResult.Fullfilled {
		result.Warning = true
		return result
	}

	// if failure and warning criteria are not met, we pass
	result.Pass = true
	return result
}
