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

// KeptnComponentSpec defines the desired state of the KeptnComponent
type KeptnComponentSpec struct {
	ApplicationName    string         `json:"application,omitempty"`
	PreDeploymentCheck KeptnEventSpec `json:"preDeploymentCheck"`
}

// KeptnComponentStatus defines the observed state of the KeptnComponent
type KeptnComponentStatus struct {
	Phase          ServiceRunPhase `json:"phase"`
	ServiceRunName string          `json:"serviceRunName"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// KeptnComponent is the Schema for the keptncomponent API
type KeptnComponent struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KeptnComponentSpec   `json:"spec,omitempty"`
	Status KeptnComponentStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// KeptnComponentList contains a list of KeptnComponents
type KeptnComponentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KeptnComponent `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KeptnComponent{}, &KeptnComponentList{})
}

func (s KeptnComponent) IsCompleted() bool {
	if s.Status.Phase == ServiceRunSucceeded || s.Status.Phase == ServiceRunFailed || s.Status.Phase == ServiceRunUnknown {
		return true
	}
	return false
}

func (s KeptnComponent) IsServiceRunNotCreated() bool {
	if s.Status.Phase == "" {
		return true
	}
	return false
}
