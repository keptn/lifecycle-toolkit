package analysis

import "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"

//go:generate moq -pkg fake -skip-ensure -out ./fake/objective_evaluator_mock.go . IObjectiveEvaluator
type IObjectiveEvaluator interface {
	Evaluate(values map[string]float64, objective v1alpha3.Objective) v1alpha3.ObjectiveResult
}

//go:generate moq -pkg fake -skip-ensure -out ./fake/criteria_set_evaluator_mock.go . ICriteriaSetEvaluator
type ICriteriaSetEvaluator interface {
	Evaluate(val float64, criteriaSet v1alpha3.CriteriaSet) v1alpha3.CriteriaSetResult
}

//go:generate moq -pkg fake -skip-ensure -out ./fake/criteria_evaluator_mock.go . ICriteriaEvaluator
type ICriteriaEvaluator interface {
	Evaluate(val float64, criteria v1alpha3.Criteria) v1alpha3.CriteriaResult
}

//go:generate moq -pkg fake -skip-ensure -out ./fake/target_evaluator_mock.go . ITargetEvaluator
type ITargetEvaluator interface {
	Evaluate(val float64, target v1alpha3.Target) v1alpha3.TargetResult
}
