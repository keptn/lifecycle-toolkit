package types

import (
	"github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
)

type ProviderRequest struct {
	Objective *v1alpha3.Objective
	Query     string
	Provider  *v1alpha3.KeptnMetricsProvider
}

type ProviderResult struct {
	Objective v1alpha3.ObjectReference
	Value     string
	Err       error
}

type TargetResult struct {
	FailResult OperatorResult
	WarnResult OperatorResult
	Warning    bool
	Pass       bool
}

func (t *TargetResult) IsFail() bool {
	return t.FailResult.Fulfilled
}

func (t *TargetResult) IsWarn() bool {
	return t.WarnResult.Fulfilled
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

func (o *ObjectiveResult) IsFail() bool {
	return o.Score == 0.0
}

func (o *ObjectiveResult) IsPass() bool {
	return o.Result.Pass
}

func (o *ObjectiveResult) IsWarn() bool {
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
