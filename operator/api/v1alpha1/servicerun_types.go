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

// ServiceRunSpec defines the desired state of ServiceRun
type ServiceRunSpec struct {
	PreDeploymentCheck EventSpec `json:"preDeploymentCheck"`
	ApplicationName    string    `json:"application"`
	Version            string    `json:"version"`
	Owner              Owner     `json:"owner"`
}

type ServiceRunPhase string

const (
	// ServiceRunPending means the application has been accepted by the system, but one or more of its
	// serviceRuns has not been started.
	ServiceRunPending ServiceRunPhase = "Pending"
	// ServiceRunRunning means that all of the serviceRuns have been started.
	ServiceRunRunning ServiceRunPhase = "Running"
	// ServiceRunSucceeded means that all of the serviceRuns have been finished successfully.
	ServiceRunSucceeded ServiceRunPhase = "Succeeded"
	// ServiceRunFailed means that one or more pre-deployment checks was not successful and terminated.
	ServiceRunFailed ServiceRunPhase = "Failed"
	// ServiceRunUnknown means that for some reason the state of the application could not be obtained.
	ServiceRunUnknown ServiceRunPhase = "Unknown"
)

// ServiceRunStatus defines the observed state of ServiceRun
type ServiceRunStatus struct {
	PreDeploymentPhase     ServiceRunPhase `json:"preDeploymentPhase"`
	PreDeploymentTaskName  string          `json:"preDeploymentTaskName"`
	PostDeploymentTaskName string          `json:"postDeploymentTaskName"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ServiceRun is the Schema for the serviceruns API
type ServiceRun struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ServiceRunSpec   `json:"spec,omitempty"`
	Status ServiceRunStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ServiceRunList contains a list of ServiceRun
type ServiceRunList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ServiceRun `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ServiceRun{}, &ServiceRunList{})
}

func (s ServiceRun) IsCompleted() bool {
	if s.Status.PreDeploymentPhase == ServiceRunSucceeded || s.Status.PreDeploymentPhase == ServiceRunFailed || s.Status.PreDeploymentPhase == ServiceRunUnknown {
		return true
	}
	return false
}

func (s ServiceRun) IsDeploymentCheckNotCreated() bool {
	if s.Status.PreDeploymentPhase == ServiceRunPending || s.Status.PreDeploymentTaskName == "" {
		return true
	}
	return false
}
