package analysis

import "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"

//go:generate moq -pkg fake -skip-ensure -out ./fake/analysis_evaluator_mock.go . IAnalysisEvaluator
type IAnalysisEvaluator interface {
	Evaluate(values map[string]string, ad v1alpha3.AnalysisDefinition) AnalysisResult
}

//go:generate moq -pkg fake -skip-ensure -out ./fake/objective_evaluator_mock.go . IObjectiveEvaluator
type IObjectiveEvaluator interface {
	Evaluate(values map[string]string, objective v1alpha3.Objective) ObjectiveResult
}

//go:generate moq -pkg fake -skip-ensure -out ./fake/target_evaluator_mock.go . ITargetEvaluator
type ITargetEvaluator interface {
	Evaluate(val float64, target v1alpha3.Target) TargetResult
}

//go:generate moq -pkg fake -skip-ensure -out ./fake/objective_evaluator_mock.go . IObjectiveEvaluator
type IOperatorEvaluator interface {
	Evaluate(val float64, criteria v1alpha3.Operator) OperatorResult
}
