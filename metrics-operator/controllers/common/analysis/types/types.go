package types

import "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"

type TargetResult struct {
	FailureResult OperatorResult
	WarningResult OperatorResult
	Warning       bool
	Pass          bool
}

func (t *TargetResult) IsFailure() bool {
	return t.FailureResult.Fulfilled
}

func (t *TargetResult) IsWarning() bool {
	return t.WarningResult.Fulfilled
}

type OperatorResult struct {
	Operator  v1alpha3.Operator
	Fulfilled bool
}

type ObjectiveResult struct {
	Result       TargetResult
	Value        float64
	Score        float64
	KeyObjective bool
	Error        error
}

func (o *ObjectiveResult) IsFailure() bool {
	return o.Score == 0.0
}

func (o *ObjectiveResult) IsPass() bool {
	return o.Result.Pass
}

func (o *ObjectiveResult) IsWarning() bool {
	return o.Result.Warning
}

type AnalysisResult struct {
	ObjectiveResults []ObjectiveResult
	TotalScore       float64
	MaximumScore     float64
	Pass             bool
	Warning          bool
}

func (a *AnalysisResult) CountScores(obj v1alpha3.Objective, result ObjectiveResult) {
	a.MaximumScore += float64(obj.Weight)
	a.TotalScore += result.Score
}

func (a *AnalysisResult) GetAchievedPercentage() float64 {
	return (a.TotalScore / a.MaximumScore) * 100.0
}

func (a *AnalysisResult) HandlePercentageScore(ad v1alpha3.AnalysisDefinition) {
	achievedPercentage := a.GetAchievedPercentage()

	if achievedPercentage >= float64(ad.Spec.TotalScore.PassPercentage) {
		a.Pass = true
	} else if achievedPercentage >= float64(ad.Spec.TotalScore.WarningPercentage) {
		a.Warning = true
	}
}
