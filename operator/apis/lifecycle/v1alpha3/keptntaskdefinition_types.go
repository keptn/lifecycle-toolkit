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
	// Function contains the definition for the function that is to be executed in KeptnTasks based on
	// the KeptnTaskDefinitions.
	Function FunctionSpec `json:"function,omitempty"`
	// Retries specifies how many times a job executing the KeptnTaskDefinition should be restarted in the case
	// of an unsuccessful attempt.
	// +kubebuilder:default:=10
	Retries *int32 `json:"retries,omitempty"`
	// Timeout specifies the maximum time to wait for the task to be completed successfully.
	// If the task does not complete successfully within this time frame, it will be
	// considered to be failed.
	// +optional
	// +kubebuilder:default:="5m"
	// +kubebuilder:validation:Pattern="^0|([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	// +kubebuilder:validation:Type:=string
	// +optional
	Timeout metav1.Duration `json:"timeout,omitempty"`
}

type FunctionSpec struct {
	// FunctionReference allows to reference another KeptnTaskDefinition which contains the source code of the
	// function to be executes for KeptnTasks based on this KeptnTaskDefinition. This can be useful when you have
	// multiple KeptnTaskDefinitions that should execute the same logic, but each with different parameters.
	FunctionReference FunctionReference `json:"functionRef,omitempty"`
	// Inline allows to specify the code that should be executed directly in the KeptnTaskDefinition, as a multi-line
	// string.
	Inline Inline `json:"inline,omitempty"`
	// HttpReference allows to point to an HTTP URL containing the code of the function.
	HttpReference HttpReference `json:"httpRef,omitempty"`
	// ConfigMapReference allows to reference a ConfigMap containing the code of the function.
	// When referencing a ConfigMap, the code of the function must be available as a value of the 'code' key
	// of the referenced ConfigMap.
	ConfigMapReference ConfigMapReference `json:"configMapRef,omitempty"`
	// Parameters contains parameters that will be passed to the job that executes the task.
	Parameters TaskParameters `json:"parameters,omitempty"`
	// SecureParameters contains secure parameters that will be passed to the job that executes the task.
	// These will be stored and accessed as secrets in the cluster.
	SecureParameters SecureParameters `json:"secureParameters,omitempty"`
}

type ConfigMapReference struct {
	// Name is the name of the referenced ConfigMap.
	Name string `json:"name,omitempty"`
}

type FunctionReference struct {
	// Name is the name of the referenced KeptnTaksDefinition.
	Name string `json:"name,omitempty"`
}

type Inline struct {
	// Code contains the code of the function.
	Code string `json:"code,omitempty"`
}

type HttpReference struct {
	// Url is the URL containing the code of the function.
	Url string `json:"url,omitempty"`
}

type ContainerSpec struct {
}

// KeptnTaskDefinitionStatus defines the observed state of KeptnTaskDefinition
type KeptnTaskDefinitionStatus struct {
	// Function contains status information of the function definition for the task.
	Function FunctionStatus `json:"function,omitempty"`
}

type FunctionStatus struct {
	// ConfigMap indicates the ConfigMap in which the function code is stored.
	ConfigMap string `json:"configMap,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:storageversion
//+kubebuilder:subresource:status

// KeptnTaskDefinition is the Schema for the keptntaskdefinitions API
type KeptnTaskDefinition struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec describes the desired state of the KeptnTaskDefinition.
	Spec KeptnTaskDefinitionSpec `json:"spec,omitempty"`
	// Status describes the current state of the KeptnTaskDefinition.
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
