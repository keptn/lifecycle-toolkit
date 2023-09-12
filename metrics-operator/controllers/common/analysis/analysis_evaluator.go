package analysis

import (
	"github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/analysis/types"
)

type AnalysisEvaluator struct {
	ObjectiveEvaluator IObjectiveEvaluator
}

func NewAnalysisEvaluator(o IObjectiveEvaluator) AnalysisEvaluator {
	return AnalysisEvaluator{
		ObjectiveEvaluator: o,
	}
}

func (ae *AnalysisEvaluator) Evaluate(values map[string]v1alpha3.ProviderResult, ad *v1alpha3.AnalysisDefinition) types.AnalysisResult {
	result := types.AnalysisResult{
		ObjectiveResults: make([]types.ObjectiveResult, len(ad.Spec.Objectives)),
	}

	keyObjectiveFailed := false
	for index, objective := range ad.Spec.Objectives {
		// evaluate a single objective and store it's result
		objectiveResult := ae.ObjectiveEvaluator.Evaluate(values, &objective)
		result.ObjectiveResults[index] = objectiveResult

		// count scores
		result.MaximumScore += float64(objective.Weight)
		result.TotalScore += objectiveResult.Score

		//check if the objective was marked as 'key' and if it succeeded
		if objectiveResult.IsFail() && objective.KeyObjective {
			keyObjectiveFailed = true
		}
	}

	achievedPercentage := result.GetAchievedPercentage()

	if achievedPercentage >= float64(ad.Spec.TotalScore.PassPercentage) {
		result.Pass = true
	} else if achievedPercentage >= float64(ad.Spec.TotalScore.WarningPercentage) {
		result.Warning = true
	}

	if keyObjectiveFailed {
		result.Pass = false
		result.Warning = false
	}

	return result
}
