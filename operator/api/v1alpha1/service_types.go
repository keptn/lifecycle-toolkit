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

// ServiceSpec defines the desired state of Service
type ServiceSpec struct {
	ApplicationName   string    `json:"application,omitempty"`
	PreDeplymentCheck EventSpec `json:"preDeploymentChecks"`
}

// ServiceStatus defines the observed state of Service
type ServiceStatus struct {
	Phase                  ServicePhase `json:"phase"`
	PreDeploymentCheckName string       `json:"preDeploymentChecksName"`
}

type ServicePhase string

const (
	// ServicePending means the application has been accepted by the system, but one or more of its
	// services has not been started.
	ServicePending ServicePhase = "Pending"
	// ServiceRunning means that all of the services have been started.
	ServiceRunning ServicePhase = "Running"
	// ServiceSucceeded means that all of the services have been finished successfully.
	ServiceSucceeded ServicePhase = "Succeeded"
	// ServiceFailed means that one or more pre-deployment checks was not successful and terminated.
	ServiceFailed ServicePhase = "Failed"
	// ServiceUnknown means that for some reason the state of the application could not be obtained.
	ServiceUnknown ServicePhase = "Unknown"
)

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Service is the Schema for the services API
type Service struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ServiceSpec   `json:"spec,omitempty"`
	Status ServiceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ServiceList contains a list of Service
type ServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Service `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Service{}, &ServiceList{})
}

func (s Service) IsCompleted() bool {
	if s.Status.Phase == ServiceSucceeded || s.Status.Phase == ServiceFailed || s.Status.Phase == ServiceUnknown {
		return true
	}
	return false
}

func (s Service) IsDeploymentCheckNotCreated() bool {
	if s.Status.Phase == ServicePending || s.Status.PreDeploymentCheckName == "" {
		return true
	}
	return false
}
