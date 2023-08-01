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

// Criteria implements the ICriteria interface. This approach has the following drawback:
// The v1alpha3.Criteria wrapped by this struct would need to not reference v1alpha3.Target structs,
// but a Target interface instead, if we would like to be able to switch to another implementation of the
// Target evaluation. However, in the structs defining the CRDs having interfaces will not work for kubebuilder.
// To circumvent this, a factory returning the implementation for the Target evaluation interface can be injected into
// this struct. This however increases the complexity of the unit tests in my opinion (see the related tests in criteria_test.go)
type Criteria struct {
	v1alpha3.Criteria
	TargetFactory ITargetFactory
}

func (c *Criteria) Evaluate(val float64) v1alpha3.CriteriaResult {
	result := v1alpha3.CriteriaResult{
		ViolatedTargets: []v1alpha3.TargetResult{},
	}

	if c.AllOf != nil && len(c.AllOf) > 0 {
		result.Violated = false
		for _, target := range c.AllOf {
			targetResult := c.TargetFactory.GetTarget(target).Evaluate(val)
			if targetResult.Violated {
				result.Violated = true
				result.ViolatedTargets = append(result.ViolatedTargets, targetResult)
			}
		}
	} else if c.AnyOf != nil && len(c.AnyOf) > 0 {
		result.Violated = true
		for _, target := range c.AnyOf {
			targetResult := c.TargetFactory.GetTarget(target).Evaluate(val)
			if targetResult.Violated {
				result.ViolatedTargets = append(result.ViolatedTargets, targetResult)
			} else {
				result.Violated = false
			}
		}
	}
	return result
}
