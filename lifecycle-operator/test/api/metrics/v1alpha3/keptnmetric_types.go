package v1alpha3

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KeptnMetricSpec defines the desired state of KeptnMetric
type KeptnMetricSpec struct {
	// Provider represents the provider object
	Provider ProviderRef `json:"provider"`
	// Query represents the query to be run
	Query string `json:"query"`
	// FetchIntervalSeconds represents the update frequency in seconds that is used to update the metric
	FetchIntervalSeconds uint `json:"fetchIntervalSeconds"`
	// Range represents the time range for which data is to be queried
	Range *RangeSpec `json:"range,omitempty"`
}

// KeptnMetricStatus defines the observed state of KeptnMetric
type KeptnMetricStatus struct {
	// Value represents the resulting value
	Value string `json:"value"`
	// RawValue represents the resulting value in raw format
	RawValue []byte `json:"rawValue"`
	// LastUpdated represents the time when the status data was last updated
	LastUpdated metav1.Time `json:"lastUpdated"`
	// ErrMsg represents the error details when the query could not be evaluated
	ErrMsg string `json:"errMsg,omitempty"`
}

// ProviderRef represents the provider object
type ProviderRef struct {
	// Name of the provider
	Name string `json:"name"`
}

// RangeSpec defines the time range for which data is to be queried
type RangeSpec struct {
	// Interval specifies the duration of the time interval for the data query
	// +kubebuilder:default:="5m"
	Interval string `json:"interval,omitempty"`
}

// KeptnMetric is the Schema for the keptnmetrics API
type KeptnMetric struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KeptnMetricSpec   `json:"spec,omitempty"`
	Status KeptnMetricStatus `json:"status,omitempty"`
}

// KeptnMetricList contains a list of KeptnMetric
type KeptnMetricList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KeptnMetric `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KeptnMetric{}, &KeptnMetricList{})
}
