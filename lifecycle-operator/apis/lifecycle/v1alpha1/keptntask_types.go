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
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha1/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KeptnTaskSpec defines the desired state of KeptnTask
type KeptnTaskSpec struct {
	Workload        string      `json:"workload"`
	WorkloadVersion string      `json:"workloadVersion"`
	AppName         string      `json:"app"`
	AppVersion      string      `json:"appVersion"`
	TaskDefinition  string      `json:"taskDefinition"`
	Context         TaskContext `json:"context"`
	// +optional
	Parameters TaskParameters `json:"parameters,omitempty"`
	// +optional
	SecureParameters SecureParameters `json:"secureParameters,omitempty"`
	// +optional
	Type common.CheckType `json:"checkType,omitempty"`
}

type TaskContext struct {
	WorkloadName    string `json:"workloadName"`
	AppName         string `json:"appName"`
	AppVersion      string `json:"appVersion"`
	WorkloadVersion string `json:"workloadVersion"`
	TaskType        string `json:"taskType"`
	ObjectType      string `json:"objectType"`
}

type TaskParameters struct {
	// +optional
	Inline map[string]string `json:"map,omitempty"`
}

type SecureParameters struct {
	// +optional
	Secret string `json:"secret,omitempty"`
}

// KeptnTaskStatus defines the observed state of KeptnTask
type KeptnTaskStatus struct {
	// +optional
	JobName string `json:"jobName,omitempty"`
	// +kubebuilder:default:=Pending
	// +optional
	Status common.KeptnState `json:"status,omitempty"`
	// +optional
	Message string `json:"message,omitempty"`
	// +optional
	StartTime metav1.Time `json:"startTime,omitempty"`
	// +optional
	EndTime metav1.Time `json:"endTime,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="AppName",type=string,JSONPath=`.spec.app`
// +kubebuilder:printcolumn:name="AppVersion",type=string,JSONPath=`.spec.appVersion`
// +kubebuilder:printcolumn:name="WorkloadName",type=string,JSONPath=`.spec.workload`
// +kubebuilder:printcolumn:name="WorkloadVersion",type=string,JSONPath=`.spec.workloadVersion`
// +kubebuilder:printcolumn:name="Job Name",type=string,JSONPath=`.status.jobName`
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.status`

// KeptnTask is the Schema for the keptntasks API
type KeptnTask struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +optional
	Spec KeptnTaskSpec `json:"spec,omitempty"`
	// +optional
	Status KeptnTaskStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KeptnTaskList contains a list of KeptnTask
type KeptnTaskList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KeptnTask `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KeptnTask{}, &KeptnTaskList{})
}
