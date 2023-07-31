package analysis

import "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"

type Analysis struct {
	Definition v1alpha3.KeptnEvaluationDefinition
}

type MetricValue struct {
	Name  string
	Value float64
}

type ObjectiveEvaluator interface {
	Evaluate(values map[string]float32) v1alpha3.ObjectiveResult
}

type CriteriaSetEvaluator interface {
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
