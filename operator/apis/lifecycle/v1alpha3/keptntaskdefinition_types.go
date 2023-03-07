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

package v1alpha3

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// KeptnTaskDefinitionSpec defines the desired state of KeptnTaskDefinition
type KeptnTaskDefinitionSpec struct {
	Function FunctionSpec `json:"function,omitempty"`
}

type FunctionSpec struct {
	FunctionReference  FunctionReference  `json:"functionRef,omitempty"`
	Inline             Inline             `json:"inline,omitempty"`
	HttpReference      HttpReference      `json:"httpRef,omitempty"`
	ConfigMapReference ConfigMapReference `json:"configMapRef,omitempty"`
	Parameters         TaskParameters     `json:"parameters,omitempty"`
	SecureParameters   SecureParameters   `json:"secureParameters,omitempty"`
}

type ConfigMapReference struct {
	Name string `json:"name,omitempty"`
}

type FunctionReference struct {
	Name string `json:"name,omitempty"`
}

type Inline struct {
	Code string `json:"code,omitempty"`
}

type HttpReference struct {
	Url string `json:"url,omitempty"`
}

type ContainerSpec struct {
}

// KeptnTaskDefinitionStatus defines the observed state of KeptnTaskDefinition
type KeptnTaskDefinitionStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Function FunctionStatus `json:"function,omitempty"`
}

type FunctionStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	ConfigMap string `json:"configMap,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:storageversion
//+kubebuilder:subresource:status

// KeptnTaskDefinition is the Schema for the keptntaskdefinitions API
type KeptnTaskDefinition struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KeptnTaskDefinitionSpec   `json:"spec,omitempty"`
	Status KeptnTaskDefinitionStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// KeptnTaskDefinitionList contains a list of KeptnTaskDefinition
type KeptnTaskDefinitionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KeptnTaskDefinition `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KeptnTaskDefinition{}, &KeptnTaskDefinitionList{})
}
