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

// KeptnWorkloadInstanceSpec defines the desired state of KeptnWorkloadInstance
type KeptnWorkloadInstanceSpec struct {
	PreDeploymentCheck EventSpec         `json:"preDeploymentCheck"`
	AppName            string            `json:"app"`
	Version            string            `json:"version"`
	ResourceReference  ResourceReference `json:"resourceReference"`
}

// KeptnWorkloadInstanceStatus defines the observed state of KeptnWorkloadInstance
type KeptnWorkloadInstanceStatus struct {
	PreDeploymentPhase     WorkloadInstancePhase `json:"preDeploymentPhase"`
	PreDeploymentTaskName  string                `json:"preDeploymentTaskName"`
	PostDeploymentTaskName string                `json:"postDeploymentTaskName"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// KeptnWorkloadInstance is the Schema for the keptnworkloadinstances API
type KeptnWorkloadInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KeptnWorkloadInstanceSpec   `json:"spec,omitempty"`
	Status KeptnWorkloadInstanceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// KeptnWorkloadInstanceList contains a list of KeptnWorkloadInstance
type KeptnWorkloadInstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KeptnWorkloadInstance `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KeptnWorkloadInstance{}, &KeptnWorkloadInstanceList{})
}

type WorkloadInstancePhase string

const (
	// WorkloadInstancePhasePending means that none of the WorkloadInstances have been created.
	WorkloadInstancePhasePending WorkloadInstancePhase = "Pending"
	// WorkloadInstancePhaseRunning means that all of the WorkloadInstances have been started.
	WorkloadInstancePhaseRunning WorkloadInstancePhase = "Running"
	// WorkloadInstancePhaseSucceeded means that all of the WorkloadInstances have been finished successfully.
	WorkloadInstancePhaseSucceeded WorkloadInstancePhase = "Succeeded"
	// WorkloadInstancePhaseFailed means that one or more pre-deployment checks was not successful and terminated.
	WorkloadInstancePhaseFailed WorkloadInstancePhase = "Failed"
	// WorkloadInstancePhaseUnknown means that for some reason the state of the application could not be obtained.
	WorkloadInstancePhaseUnknown WorkloadInstancePhase = "Unknown"
)

func (i KeptnWorkloadInstance) IsCompleted() bool {
	if i.Status.PreDeploymentPhase == WorkloadInstancePhaseSucceeded || i.Status.PreDeploymentPhase == WorkloadInstancePhaseFailed || i.Status.PreDeploymentPhase == WorkloadInstancePhaseUnknown {
		return true
	}
	return false
}

func (i KeptnWorkloadInstance) IsDeploymentCheckNotCreated() bool {
	if i.Status.PreDeploymentPhase == WorkloadInstancePhasePending || i.Status.PreDeploymentTaskName == "" {
		return true
	}
	return false
}
