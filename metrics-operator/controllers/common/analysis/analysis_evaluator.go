package analysis

import (
	"github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
)

type AnalysisEvaluator struct {
	ObjectiveEvaluator IObjectiveEvaluator
}

func NewAnalysisEvaluator(o IObjectiveEvaluator) AnalysisEvaluator {
	return AnalysisEvaluator{
		ObjectiveEvaluator: o,
	}
}

func (ae *AnalysisEvaluator) Evaluate(values map[string]string, ad v1alpha3.AnalysisDefinition) (v1alpha3.AnalysisResult, error) {
	result := v1alpha3.AnalysisResult{
		ObjectiveResults: make([]v1alpha3.ObjectiveResult, 0, len(ad.Spec.Objectives)),
	}

	keyObjectiveFailed := false
	for _, objective := range ad.Spec.Objectives {
		// evaluate a single objective and store it's result
		objectiveResult := ae.ObjectiveEvaluator.Evaluate(values, objective)
		result.ObjectiveResults = append(result.ObjectiveResults, objectiveResult)

		// count scores
		result.MaximumScore += float64(objective.Weight)
		result.TotalScore += objectiveResult.Score

		//check if the objective was marked as 'key' and if it succeeded
		if objectiveResult.Failed && objectiveResult.KeyObjective {
			keyObjectiveFailed = true
		}
	}

	achievedPercentage := (result.TotalScore / result.MaximumScore) * 100.0

	if achievedPercentage >= float64(ad.Spec.TotalScore.PassPercentage) {
		result.Pass = true
	} else if achievedPercentage >= float64(ad.Spec.TotalScore.WarningPercentage) {
		result.Warning = true
	}

	if keyObjectiveFailed {
		result.Pass = false
		result.Warning = false
	}

	return result, nil
}
