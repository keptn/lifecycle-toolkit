package types

import "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"

type TargetResult struct {
	FailureResult OperatorResult
	WarningResult OperatorResult
	Warning       bool
	Pass          bool
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
	Failed       bool
	Error        error
}

type AnalysisResult struct {
	ObjectiveResults []ObjectiveResult
	TotalScore       float64
	MaximumScore     float64
	Pass             bool
	Warning          bool
}
