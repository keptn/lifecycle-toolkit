package analysis

import "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"

type AnalysisEvaluator struct {
	ObjectiveEvaluator IObjectiveEvaluator
}

type AnalysisResult struct {
	ObjectiveResults []v1alpha3.ObjectiveResult
	TotalScore       float64
	MaximumScore     float64
	Pass             bool
	Warning          bool
}

func (ae *AnalysisEvaluator) Evaluate(values map[string]float64, ed v1alpha3.KeptnEvaluationDefinition) (*AnalysisResult, error) {
	result := &AnalysisResult{
		ObjectiveResults: []v1alpha3.ObjectiveResult{},
	}

	keyObjectiveFailed := false
	for _, objective := range ed.Spec.Objectives {
		result.MaximumScore += float64(objective.Weight)
		objectiveResult := ae.ObjectiveEvaluator.Evaluate(values, objective)
		result.TotalScore += objectiveResult.Score
		if objectiveResult.Score == 0 && objectiveResult.KeyObjective {
			keyObjectiveFailed = true
		}
		result.ObjectiveResults = append(result.ObjectiveResults, objectiveResult)
	}

	achievedPercentage := (result.TotalScore / result.MaximumScore) * 100.0

	if achievedPercentage >= ed.Spec.TotalScore.PassPercentage {
		result.Pass = true
	} else if achievedPercentage >= ed.Spec.TotalScore.WarningPercentage {
		result.Warning = true
	}

	if keyObjectiveFailed {
		result.Pass = false
		result.Warning = false
	}

	return result, nil
}
