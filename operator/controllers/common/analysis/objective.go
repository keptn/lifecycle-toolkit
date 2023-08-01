package analysis

import (
	"errors"
	"github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
)

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
	result.PassResult = oe.CriteriaSetEvaluator.Evaluate(val, obj.SLOTargets.Pass)

	// if pass criteria are successful, we can return without checking 'Warning criteria'
	if len(result.PassResult.ViolatedCriteria) == 0 {
		result.Score = float64(obj.Weight)
		return result
	}

	result.WarningResult = oe.CriteriaSetEvaluator.Evaluate(val, obj.SLOTargets.Warning)

	if len(result.WarningResult.ViolatedCriteria) == 0 {
		result.Score = float64(obj.Weight) / 2
	}

	return result
}
