package analysis

import (
	"fmt"
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

func (oe *ObjectiveEvaluator) Evaluate(values map[string]v1alpha3.ProviderResult, obj *v1alpha3.Objective) types.ObjectiveResult {
	result := types.ObjectiveResult{
		Score:     0.0,
		Value:     0.0,
		Objective: obj,
	}

	// get the value
	floatVal, err := getValueFromMap(values, ComputeKey(obj.AnalysisValueTemplateRef))
	if err != nil {
		result.Error = err
		return result
	}

	result.Value = floatVal
	result.Result = oe.TargetEvaluator.Evaluate(floatVal, &obj.Target)

	// if target passed, we return the full score
	if result.IsPass() {
		result.Score = float64(obj.Weight)
		return result
	}

	// if target fullfilled warning criteria, we return the half score
	if result.IsWarn() {
		result.Score = float64(obj.Weight) / 2
		return result
	}

	return result
}

func getValueFromMap(values map[string]v1alpha3.ProviderResult, key string) (float64, error) {
	val, ok := values[key]
	if !ok {
		return 0.0, fmt.Errorf("required value '%s' not available", key)
	}
	floatVal, err := strconv.ParseFloat(val.Value, 64)
	if err != nil {
		return 0.0, err
	}

	return floatVal, nil
}

func ComputeKey(obj v1alpha3.ObjectReference) string {
	if !obj.IsNamespaceSet() {
		return obj.Name
	}
	return obj.Name + "-" + obj.Namespace
}
