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
	"github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1/common"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// EventSpec defines the desired state of Event
type EventSpec struct {
	Service     string          `json:"service,omitempty"`
	Application string          `json:"application,omitempty"`
	JobSpec     batchv1.JobSpec `json:"job,omitempty"`
}

// EventStatus defines the observed state of Event
type EventStatus struct {
	Phase   common.KeptnState `json:"phase"`
	JobName string            `json:"jobName"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Event is the Schema for the events API
type Event struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EventSpec   `json:"spec,omitempty"`
	Status EventStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// EventList contains a list of Event
type EventList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Event `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Event{}, &EventList{})
}

func (e Event) IsCompleted() bool {
	if e.Status.Phase == common.StateSucceeded || e.Status.Phase == common.StateFailed || e.Status.Phase == common.StateUnknown {
		return true
	}
	return false
}

func (e Event) IsJobNotCreated() bool {
	if e.Status.Phase == common.StatePending || e.Status.JobName == "" {
		return true
	}
	return false
}
