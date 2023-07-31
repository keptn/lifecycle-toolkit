package analysis

import (
	"github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	"github.com/pkg/errors"
)

type Analysis struct {
	Definition v1alpha3.KeptnEvaluationDefinition
}

type MetricValue struct {
	Name  string
	Value float64
}

type IObjectiveEvaluator interface {
	Evaluate(values map[string]float32, objective v1alpha3.Objective) v1alpha3.ObjectiveResult
}

type ICriteriaSetEvaluator interface {
	Evaluate(val float64, criteriaSet v1alpha3.CriteriaSet) v1alpha3.CriteriaSetResult
}

type ICriteriaEvaluator interface {
	Evaluate(val float64, criteria v1alpha3.Criteria) v1alpha3.CriteriaResult
}

type ITargetEvaluator interface {
	Evaluate(val float64, target v1alpha3.Target) v1alpha3.TargetResult
}

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
		for _, target := range c.AllOf {
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
		for _, criteria := range cs.AllOf {
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

type ObjectiveEvaluator struct {
	CriteriaSetEvaluator ICriteriaSetEvaluator
}

func (oe *ObjectiveEvaluator) Evaluate(values map[string]float64, obj v1alpha3.Objective) v1alpha3.ObjectiveResult {
	result := v1alpha3.ObjectiveResult{
		KeyObjective: obj.KeyObjective,
		Score:        0.0,
	}
	// get the value
	val, ok := values[obj.KeptnMetricRef.Name]
	if !ok {
		result.Error = errors.New("required value not available")
		return result
	}

	result.Value = val

	// check 'pass'  criteria
	passEvaluation := oe.CriteriaSetEvaluator.Evaluate(val, obj.SLOTargets.Pass)

	// if pass criteria are successful, we can return without checking 'Warning criteria'
	if len(passEvaluation.ViolatedCriteria) == 0 {
		result.Score = float64(obj.Weight)
		return result
	}

	warnEvaluation := oe.CriteriaSetEvaluator.Evaluate(val, obj.SLOTargets.Warning)

	if len(warnEvaluation.ViolatedCriteria) == 0 {
		result.Score = float64(obj.Weight) / 2
	}

	return result
}

type AnalysisResult struct {
}

func EvaluateAnalysis(values map[string]float64, ed v1alpha3.KeptnEvaluationDefinition) (*AnalysisResult, error) {
	result := &AnalysisResult{}

	for _, objective := range ed.Spec.Objectives {
		objective.Evaluate(values)
	}

	return result, nil
}
