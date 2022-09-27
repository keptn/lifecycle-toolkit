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

// KeptnAppSpec defines the desired state of KeptnApp
type KeptnAppSpec struct {
	Workloads           []KeptnWorkload `json:"workloads,omitempty"`
	PreDeploymentTasks  []string        `json:"preDeploymentTasks,omitempty"`
	PostDeploymentTasks []string        `json:"postDeploymentTasks,omitempty"`
}

// KeptnAppStatus defines the observed state of KeptnApp
type KeptnAppStatus struct {
	Status ApplicationPhase `json:"status"`
}

type ApplicationPhase string

const (
	// ApplicationPending means the application has been accepted by the system, but one or more of its
	// workloads has not been started.
	ApplicationPending ApplicationPhase = "Pending"
	// ApplicationRunning means that all of the workloads have been started.
	ApplicationRunning ApplicationPhase = "Running"
	// ApplicationFailed means that one or more pre-deployment checks was not successful and terminated.
	ApplicationFailed ApplicationPhase = "Failed"
	// ApplicationUnknown means that for some reason the state of the application could not be obtained.
	ApplicationUnknown ApplicationPhase = "Unknown"
)

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// KeptnApp is the Schema for the keptnapps API
type KeptnApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KeptnAppSpec   `json:"spec,omitempty"`
	Status KeptnAppStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// KeptnAppList contains a list of KeptnApp
type KeptnAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KeptnApp `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KeptnApp{}, &KeptnAppList{})
}
