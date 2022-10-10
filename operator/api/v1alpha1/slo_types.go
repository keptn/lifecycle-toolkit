/*
Copyright 2022.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SloSpec defines the desired state of Slo
type SloSpec struct {
	Comparison SLOComparison `json:"comparison"`
	Filter
	Objectives
	TotalScore
}

// SloStatus defines the observed state of Slo
type SloStatus struct {
	// Slo does not have any status
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Slo is the Schema for the sloes API
type Slo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SloSpec   `json:"spec,omitempty"`
	Status SloStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SloList contains a list of Slo
type SloList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Slo `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Slo{}, &SloList{})
}
