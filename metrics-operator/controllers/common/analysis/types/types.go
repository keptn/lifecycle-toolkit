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
	Result TargetResult
	Value  float64
	Score  float64
	Error  error
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

func (a *AnalysisResult) GetAchievedPercentage() float64 {
	if a.MaximumScore == 0.0 {
		return 0.0
	}
	if a.TotalScore == a.MaximumScore {
		return 100.0
	}
	return (a.TotalScore / a.MaximumScore) * 100.0
}
