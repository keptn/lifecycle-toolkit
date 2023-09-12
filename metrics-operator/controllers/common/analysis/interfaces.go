package analysis

import (
	"github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/analysis/types"
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/analysis_evaluator_mock.go . IAnalysisEvaluator
type IAnalysisEvaluator interface {
	Evaluate(values map[string]v1alpha3.ProviderResult, ad *v1alpha3.AnalysisDefinition) types.AnalysisResult
}

//go:generate moq -pkg fake -skip-ensure -out ./fake/objective_evaluator_mock.go . IObjectiveEvaluator
type IObjectiveEvaluator interface {
	Evaluate(values map[string]v1alpha3.ProviderResult, objective *v1alpha3.Objective) types.ObjectiveResult
}

//go:generate moq -pkg fake -skip-ensure -out ./fake/target_evaluator_mock.go . ITargetEvaluator
type ITargetEvaluator interface {
	Evaluate(val float64, target *v1alpha3.Target) types.TargetResult
}

//go:generate moq -pkg fake -skip-ensure -out ./fake/operator_evaluator_mock.go . IOperatorEvaluator
type IOperatorEvaluator interface {
	Evaluate(val float64, criteria *v1alpha3.Operator) types.OperatorResult
}
