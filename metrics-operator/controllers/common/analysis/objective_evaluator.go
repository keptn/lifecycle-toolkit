package analysis

import (
	"errors"
	"strconv"

	"github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/analysis/types"
)

type ObjectiveEvaluator struct {
	TargetEvaluator ITargetEvaluator
}

func NewObjectiveEvaluator(t ITargetEvaluator) ObjectiveEvaluator {
	return ObjectiveEvaluator{
		TargetEvaluator: t,
	}
}

func (oe *ObjectiveEvaluator) Evaluate(values map[string]string, obj v1alpha3.Objective) types.ObjectiveResult {
	result := types.ObjectiveResult{
		KeyObjective: obj.KeyObjective,
		Score:        0.0,
		Value:        0.0,
	}

	// get the value
	floatVal, err := getValueFromMap(values, obj.AnalysisValueTemplateRef.Name)
	if err != nil {
		result.Error = err
		return result
	}

	result.Value = floatVal
	result.Result = oe.TargetEvaluator.Evaluate(floatVal, obj.Target)

	// if target passed, we return the full score
	if result.IsPass() {
		result.Score = float64(obj.Weight)
		return result
	}

	// if target fullfilled warning criteria, we return the half score
	if result.IsWarning() {
		result.Score = float64(obj.Weight) / 2
		return result
	}

	return result
}

func getValueFromMap(values map[string]string, name string) (float64, error) {
	val, ok := values[name]
	if !ok {
		return 0.0, errors.New("required value not available")
	}

	floatVal, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0.0, err
	}

	return floatVal, nil
}
