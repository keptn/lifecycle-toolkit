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

// KeptnTaskDefinitionSpec defines the desired state of KeptnTaskDefinition
type KeptnTaskDefinitionSpec struct {
	// +optional
	Function FunctionSpec `json:"function,omitempty"`
}

type FunctionSpec struct {
	// +optional
	FunctionReference FunctionReference `json:"functionRef,omitempty"`
	// +optional
	Inline Inline `json:"inline,omitempty"`
	// +optional
	HttpReference HttpReference `json:"httpRef,omitempty"`
	// +optional
	ConfigMapReference ConfigMapReference `json:"configMapRef,omitempty"`
	// +optional
	Parameters TaskParameters `json:"parameters,omitempty"`
	// +optional
	SecureParameters SecureParameters `json:"secureParameters,omitempty"`
}

type ConfigMapReference struct {
	// +optional
	Name string `json:"name,omitempty"`
}

type FunctionReference struct {
	// +optional
	Name string `json:"name,omitempty"`
}

type Inline struct {
	// +optional
	Code string `json:"code,omitempty"`
}

type HttpReference struct {
	// +optional
	Url string `json:"url,omitempty"`
}

type ContainerSpec struct {
}

// KeptnTaskDefinitionStatus defines the observed state of KeptnTaskDefinition
type KeptnTaskDefinitionStatus struct {
	// +optional
	Function FunctionStatus `json:"function,omitempty"`
}

type FunctionStatus struct {
	// +optional
	ConfigMap string `json:"configMap,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// KeptnTaskDefinition is the Schema for the keptntaskdefinitions API
type KeptnTaskDefinition struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +optional
	Spec KeptnTaskDefinitionSpec `json:"spec,omitempty"`
	// +optional
	Status KeptnTaskDefinitionStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KeptnTaskDefinitionList contains a list of KeptnTaskDefinition
type KeptnTaskDefinitionList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KeptnTaskDefinition `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KeptnTaskDefinition{}, &KeptnTaskDefinitionList{})
}
