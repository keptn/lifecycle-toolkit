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

// AnalysisValueTemplateSpec defines the desired state of AnalysisValueTemplate
type AnalysisValueTemplateSpec struct {
	// Provider represents the provider object
	Provider ProviderRef `json:"provider"`
	// Query represents the query to be run. It can include placeholders that are defined using the go template
	// syntax (https://pkg.go.dev/text/template).
	Query string `json:"query"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Provider",type=string,JSONPath=`.spec.provider.name`
// +kubebuilder:printcolumn:name="Query",type=string,JSONPath=`.spec.query`

// AnalysisValueTemplate is the Schema for the analysisvaluetemplates API
type AnalysisValueTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec contains the specification for the AnalysisValueTemplate
	Spec   AnalysisValueTemplateSpec `json:"spec,omitempty"`
	Status EmptyStatus               `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AnalysisValueTemplateList contains a list of AnalysisValueTemplate
type AnalysisValueTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AnalysisValueTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AnalysisValueTemplate{}, &AnalysisValueTemplateList{})
}
