/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha3

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AnalysisValueSpec defines the desired state of AnalysisValue
type AnalysisValueSpec struct {
	AnalysisTemplate AnalysisTemplateRef  `json:"analysisTemplate"`
	Timeframe        *TimeframeDefinition `json:"timeframe,omitempty"`
	Selectors        map[string]string    `json:"selectors,omitempty"`
}

type TimeframeDefinition struct {
	From metav1.Time `json:"from"`
	To   metav1.Time `json:"to"`
}

type AnalysisTemplateRef struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
}

// AnalysisValueStatus defines the observed state of AnalysisValue
type AnalysisValueStatus struct {
	Query    string `json:"query,omitempty"`
	Value    string `json:"value,omitempty"`
	RawValue []byte `json:"rawValue"`
	ErrMsg   string `json:"errMsg,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// AnalysisValue is the Schema for the analysisvalues API
type AnalysisValue struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AnalysisValueSpec   `json:"spec,omitempty"`
	Status AnalysisValueStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AnalysisValueList contains a list of AnalysisValue
type AnalysisValueList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AnalysisValue `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AnalysisValue{}, &AnalysisValueList{})
}
