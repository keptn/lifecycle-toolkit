package v1alpha3

type ObjectReference struct {
	// Name defines the name of the referenced object
	Name string `json:"name"`
	// Namespace defines the namespace of the referenced object
	Namespace string `json:"namespace,omitempty"`
}

type TargetResult struct {
	FailureResult OperatorResult
	WarningResult OperatorResult
	Warning       bool
	Pass          bool
}

type OperatorResult struct {
	Operator  Operator
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
