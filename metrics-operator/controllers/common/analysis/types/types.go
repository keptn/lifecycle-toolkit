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
	FailResult OperatorResult `json:"failResult"`
	WarnResult OperatorResult `json:"warnResult"`
	Warning    bool           `json:"warning"`
	Pass       bool           `json:"pass"`
}

func (t *TargetResult) IsFail() bool {
	return t.FailResult.Fulfilled
}

func (t *TargetResult) IsWarn() bool {
	return t.WarnResult.Fulfilled
}

type OperatorResult struct {
	Operator  v1alpha3.Operator `json:"operator"`
	Fulfilled bool              `json:"fulfilled"`
}

type ObjectiveResult struct {
	Result TargetResult `json:"result"`
	Value  float64      `json:"value"`
	Score  float64      `json:"score"`
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
	ObjectiveResults []ObjectiveResult `json:"objectiveResults"`
	TotalScore       float64           `json:"totalScore"`
	MaximumScore     float64           `json:"maximumScore"`
	Pass             bool              `json:"pass"`
	Warning          bool              `json:"warning"`
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
