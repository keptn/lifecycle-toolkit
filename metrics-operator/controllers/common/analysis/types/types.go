package types

import (
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1"
)

type ProviderRequest struct {
	Objective metricsapi.Objective
	Query     string
	Provider  *metricsapi.KeptnMetricsProvider
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
	Operator  metricsapi.Operator `json:"operator"`
	Fulfilled bool                `json:"fulfilled"`
}

type ObjectiveResult struct {
	Result    TargetResult         `json:"result"`
	Objective metricsapi.Objective `json:"objective"`
	Value     float64              `json:"value"`
	Query     string               `json:"query"`
	Score     float64              `json:"score"`
	Error     error                `json:"error,omitempty"`
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

// AnalysisCompletion consolidates an analysis definition and its result into one struct, which is needed to communicate
// both objects via a channel
type AnalysisCompletion struct {
	Result   AnalysisResult
	Analysis metricsapi.Analysis
}
