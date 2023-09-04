package types

import (
	"github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
)

type ProviderRequest struct {
	Objective *v1alpha3.Objective
	Query     string
	Provider  *v1alpha3.KeptnMetricsProvider
}

type TargetResult struct {
	FailResult OperatorResult `json:"failResult,omitempty"`
	WarnResult OperatorResult `json:"warnResult,omitempty"`
	Warning    bool           `json:"warning,omitempty"`
	Pass       bool           `json:"pass,omitempty"`
}

func (t *TargetResult) IsFail() bool {
	return t.FailResult.Fulfilled
}

func (t *TargetResult) IsWarn() bool {
	return t.WarnResult.Fulfilled
}

type OperatorResult struct {
	Operator  v1alpha3.Operator `json:"operator,omitempty"`
	Fulfilled bool              `json:"fulfilled,omitempty"`
}

type ObjectiveResult struct {
	Result TargetResult `json:"result,omitempty"`
	Value  float64      `json:"value,omitempty"`
	Score  float64      `json:"score,omitempty"`
	Error  error        `json:"error,omitempty"`
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
	ObjectiveResults []ObjectiveResult `json:"objectiveResults,omitempty"`
	TotalScore       float64           `json:"totalScore,omitempty"`
	MaximumScore     float64           `json:"maximumScore,omitempty"`
	Pass             bool              `json:"pass,omitempty"`
	Warning          bool              `json:"warning,omitempty"`
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
