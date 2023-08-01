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

//////////////////

//go:generate moq -pkg fake -skip-ensure -out ./fake/objective_mock.go . IObjective
type IObjective interface {
	Evaluate(values map[string]float64) v1alpha3.ObjectiveResult
}

//go:generate moq -pkg fake -skip-ensure -out ./fake/criteria_set_mock.go . ICriteriaSet
type ICriteriaSet interface {
	Evaluate(val float64) v1alpha3.CriteriaSetResult
}

//go:generate moq -pkg fake -skip-ensure -out ./fake/criteria_mock.go . ICriteria
type ICriteria interface {
	Evaluate(val float64) v1alpha3.CriteriaResult
}

//go:generate moq -pkg fake -skip-ensure -out ./fake/target_mock.go . ITarget
type ITarget interface {
	Evaluate(val float64) v1alpha3.TargetResult
}

// Factory interfaces

//go:generate moq -pkg fake -skip-ensure -out ./fake/target_factory_mock.go . ITargetFactory
type ITargetFactory interface {
	GetTarget(target v1alpha3.Target) ITarget
}

//go:generate moq -pkg fake -skip-ensure -out ./fake/criteria_factory_mock.go . ICriteriaFactory
type ICriteriaFactory interface {
	GetCriteria(target v1alpha3.Target) ICriteria
}
