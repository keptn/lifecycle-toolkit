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

// KeptnConfigSpec defines the desired state of KeptnConfig
type KeptnConfigSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// OTelCollectorUrl can be used to set the Open Telemetry collector that the lifecycle operator should use
	// +optional
	OTelCollectorUrl string `json:"OTelCollectorUrl,omitempty"`
	// KeptnAppCreationRequestTimeoutSeconds is used to set the interval in which automatic app discovery
	// searches for workload to put into the same auto-generated KeptnApp
	// +kubebuilder:default:=30
	// +optional
	KeptnAppCreationRequestTimeoutSeconds uint `json:"keptnAppCreationRequestTimeoutSeconds,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// KeptnConfig is the Schema for the keptnconfigs API
type KeptnConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec KeptnConfigSpec `json:"spec,omitempty"`
	// unused field
	Status string `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KeptnConfigList contains a list of KeptnConfig
type KeptnConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KeptnConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KeptnConfig{}, &KeptnConfigList{})
}
