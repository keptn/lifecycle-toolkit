package analysis

import "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"

type CriteriaSetEvaluator struct {
	CriteriaEvaluator ICriteriaEvaluator
}

func (cse *CriteriaSetEvaluator) Evaluate(val float64, cs v1alpha3.CriteriaSet) v1alpha3.CriteriaSetResult {
	result := v1alpha3.CriteriaSetResult{
		ViolatedCriteria: []v1alpha3.CriteriaResult{},
	}

	if cs.AllOf != nil && len(cs.AllOf) > 0 {
		result.Violated = false
		for _, criteria := range cs.AllOf {
			criteriaResult := cse.CriteriaEvaluator.Evaluate(val, criteria)
			if criteriaResult.Violated {
				result.Violated = true
				result.ViolatedCriteria = append(result.ViolatedCriteria, criteriaResult)
			}
		}
	} else if cs.AnyOf != nil && len(cs.AnyOf) > 0 {
		result.Violated = true
		for _, criteria := range cs.AnyOf {
			criteriaResult := cse.CriteriaEvaluator.Evaluate(val, criteria)
			if criteriaResult.Violated {
				result.ViolatedCriteria = append(result.ViolatedCriteria, criteriaResult)
			} else {
				result.Violated = false
			}
		}
	}
	return result
}
