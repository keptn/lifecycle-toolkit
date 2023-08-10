package analysis

import (
	"errors"
	"strconv"

	"github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
)

type ObjectiveEvaluator struct {
	TargetEvaluator ITargetEvaluator
}

type ObjectiveResult struct {
	Result       TargetResult
	Value        float64
	Score        float64
	KeyObjective bool
	Failed       bool
	Error        error
}

func NewObjectiveEvaluator(t ITargetEvaluator) ObjectiveEvaluator {
	return ObjectiveEvaluator{
		TargetEvaluator: t,
	}
}

func (oe *ObjectiveEvaluator) Evaluate(values map[string]string, obj v1alpha3.Objective) ObjectiveResult {
	result := ObjectiveResult{
		KeyObjective: obj.KeyObjective,
		Score:        0.0,
		Failed:       false,
	}

	// get the value
	floatVal, err := oe.getValue(values, obj.AnalysisValueTemplateRef.Name)
	if err != nil {
		result.Error = err
		result.Failed = true
		return result
	}

	result.Value = floatVal
	result.Result = oe.TargetEvaluator.Evaluate(floatVal, obj.Target)

	// if target passed, we return the full score
	if result.Result.Pass {
		result.Score = float64(obj.Weight)
		return result
	}

	// if target fullfilled warning criteria, we return the half score
	if result.Result.Warning {
		result.Score = float64(obj.Weight) / 2
		return result
	}

	result.Failed = true
	return result
}

func (oe *ObjectiveEvaluator) getValue(values map[string]string, name string) (float64, error) {
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
