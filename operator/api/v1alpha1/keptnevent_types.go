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
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// KeptnEventSpec defines the desired state of KeptnEvent
type KeptnEventSpec struct {
	Component   string          `json:"component,omitempty"`
	Application string          `json:"application,omitempty"`
	JobSpec     batchv1.JobSpec `json:"job,omitempty"`
}

// KeptnEventStatus defines the observed state of KeptnEvent
type KeptnEventStatus struct {
	Phase   KeptnEventPhase `json:"phase"`
	JobName string          `json:"jobName"`
}

type KeptnEventPhase string

const (
	// EventPending means the application has been accepted by the system, but one or more of its
	// services has not been started.
	EventPending KeptnEventPhase = "Pending"
	// EventRunning means that all of the services have been started.
	EventRunning KeptnEventPhase = "Running"
	// EventRunning means that all of the services have been finished successfully.
	EventSucceeded KeptnEventPhase = "Succeeded"
	// EventFailed means that one or more pre-deployment checks was not successful and terminated.
	EventFailed KeptnEventPhase = "Failed"
	// EventUnknown means that for some reason the state of the application could not be obtained.
	EventUnknown KeptnEventPhase = "Unknown"
)

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// KeptnEvent is the Schema for the keptnevents API
type KeptnEvent struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KeptnEventSpec   `json:"spec,omitempty"`
	Status KeptnEventStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// KeptnEventList contains a list of Event
type KeptnEventList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KeptnEvent `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KeptnEvent{}, &KeptnEventList{})
}

func (e KeptnEvent) IsCompleted() bool {
	if e.Status.Phase == EventSucceeded || e.Status.Phase == EventFailed || e.Status.Phase == EventUnknown {
		return true
	}
	return false
}

func (e KeptnEvent) IsJobNotCreated() bool {
	if e.Status.Phase == EventPending || e.Status.JobName == "" {
		return true
	}
	return false
}
