package analysis

import "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"

type CriteriaEvaluator struct {
	TargetEvaluator ITargetEvaluator
}

func (ce *CriteriaEvaluator) Evaluate(val float64, c v1alpha3.Criteria) v1alpha3.CriteriaResult {
	result := v1alpha3.CriteriaResult{
		ViolatedTargets: []v1alpha3.TargetResult{},
	}

	if c.AllOf != nil && len(c.AllOf) > 0 {
		result.Violated = false
		for _, target := range c.AllOf {
			targetResult := ce.TargetEvaluator.Evaluate(val, target)
			if targetResult.Violated {
				result.Violated = true
				result.ViolatedTargets = append(result.ViolatedTargets, targetResult)
			}
		}
	} else if c.AnyOf != nil && len(c.AnyOf) > 0 {
		result.Violated = true
		for _, target := range c.AnyOf {
			targetResult := ce.TargetEvaluator.Evaluate(val, target)
			if targetResult.Violated {
				result.ViolatedTargets = append(result.ViolatedTargets, targetResult)
			} else {
				result.Violated = false
			}
		}
	}
	return result
}
