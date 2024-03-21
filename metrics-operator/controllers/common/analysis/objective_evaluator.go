package analysis

import (
	"fmt"
	"strconv"

	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1"
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

func (oe *ObjectiveEvaluator) Evaluate(values map[string]metricsapi.ProviderResult, obj *metricsapi.Objective) types.ObjectiveResult {
	result := types.ObjectiveResult{
		Score:     0.0,
		Value:     0.0,
		Objective: *obj,
	}

	// get the value
	floatVal, query, err := getResultFromMap(values, ComputeKey(obj.AnalysisValueTemplateRef))
	result.Query = query
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

	// if target fulfilled warning criteria, we return the half score
	if result.IsWarn() {
		result.Score = float64(obj.Weight) / 2
		return result
	}

	return result
}

func getResultFromMap(values map[string]metricsapi.ProviderResult, key string) (float64, string, error) {
	val, ok := values[key]
	if !ok {
		return 0.0, "", fmt.Errorf("required value '%s' not available", key)
	}
	floatVal, err := strconv.ParseFloat(val.Value, 64)
	if err != nil {
		return 0.0, val.Query, err
	}

	return floatVal, val.Query, nil
}

func ComputeKey(obj metricsapi.ObjectReference) string {
	if !obj.IsNamespaceSet() {
		return obj.Name
	}
	return obj.Name + "-" + obj.Namespace
}
