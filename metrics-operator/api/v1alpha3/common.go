package v1alpha3

type ObjectReference struct {
	// Name defines the name of the referenced object
	Name string `json:"name"`
	// Namespace defines the namespace of the referenced object
	// +optional
	Namespace string `json:"namespace,omitempty"`
}

// AnalysisState represents the state of the analysis
type AnalysisState string

const (
	StatePending     AnalysisState = "Pending"
	StateProgressing AnalysisState = "Progressing"
	StateCompleted   AnalysisState = "Completed"
)
