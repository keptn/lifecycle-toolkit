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

func (ae *AnalysisEvaluator) Evaluate(values map[string]string, ad v1alpha3.AnalysisDefinition) (types.AnalysisResult, error) {
	result := types.AnalysisResult{
		ObjectiveResults: make([]types.ObjectiveResult, 0, len(ad.Spec.Objectives)),
	}

	keyObjectiveFailed := false
	for _, objective := range ad.Spec.Objectives {
		// evaluate a single objective and store it's result
		objectiveResult := ae.ObjectiveEvaluator.Evaluate(values, objective)
		result.ObjectiveResults = append(result.ObjectiveResults, objectiveResult)

		// count scores
		result.CountScores(objective, objectiveResult)

		//check if the objective was marked as 'key' and if it succeeded
		if objectiveResult.IsFailure() && objectiveResult.KeyObjective {
			keyObjectiveFailed = true
		}
	}

	result.HandlePercentageScore(ad)

	if keyObjectiveFailed {
		result.Pass = false
		result.Warning = false
	}

	return result, nil
}
