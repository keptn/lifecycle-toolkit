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

package v1beta1

import (
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1beta1/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KeptnTaskSpec defines the desired state of KeptnTask
type KeptnTaskSpec struct {
	// TaskDefinition refers to the name of the KeptnTaskDefinition
	// which includes the specification for the task to be performed.
	// The KeptnTaskDefinition can be
	// located in the same namespace as the KeptnTask, or in the Keptn namespace.
	TaskDefinition string `json:"taskDefinition"`
	// Context contains contextual information about the task execution.
	// +optional
	Context TaskContext `json:"context"`
	// Parameters contains parameters that will be passed to the job that executes the task.
	// +optional
	Parameters TaskParameters `json:"parameters,omitempty"`
	// SecureParameters contains secure parameters that will be passed to the job that executes the task.
	// These will be stored and accessed as secrets in the cluster.
	// +optional
	SecureParameters SecureParameters `json:"secureParameters,omitempty"`
	// Type indicates whether the KeptnTask is part of the pre- or postDeployment phase.
	// +optional
	Type common.CheckType `json:"checkType,omitempty"`
	// Retries indicates how many times the KeptnTask can be attempted in the case of an error
	// before considering the KeptnTask to be failed.
	// +kubebuilder:default:=10
	// +optional
	Retries *int32 `json:"retries,omitempty"`
	// Timeout specifies the maximum time to wait for the task to be completed successfully.
	// If the task does not complete successfully within this time frame, it will be
	// considered to be failed.
	// +kubebuilder:default:="5m"
	// +kubebuilder:validation:Pattern="^0|([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	// +kubebuilder:validation:Type:=string
	// +optional
	Timeout metav1.Duration `json:"timeout,omitempty"`
}

type TaskContext struct {
	// WorkloadName the name of the KeptnWorkload the KeptnTask is being executed for.
	// +optional
	WorkloadName string `json:"workloadName"`
	// AppName the name of the KeptnApp the KeptnTask is being executed for.
	// +optional
	AppName string `json:"appName"`
	// AppVersion the version of the KeptnApp the KeptnTask is being executed for.
	// +optional
	AppVersion string `json:"appVersion"`
	// WorkloadVersion the version of the KeptnWorkload the KeptnTask is being executed for.
	// +optional
	WorkloadVersion string `json:"workloadVersion"`
	// TaskType indicates whether the KeptnTask is part of the pre- or postDeployment phase.
	// +optional
	TaskType string `json:"taskType"`
	// ObjectType indicates whether the KeptnTask is being executed for a KeptnApp or KeptnWorkload.
	// +optional
	ObjectType string `json:"objectType"`
	// +optional
	// Metadata contains additional key-value pairs for contextual information.
	Metadata map[string]string `json:"metadata,omitempty"`
}

type TaskParameters struct {
	// Inline contains the parameters that will be made available to the job
	// executing the KeptnTask via the 'DATA' environment variable.
	// The 'DATA'  environment variable's content will be a json
	// encoded string containing all properties of the map provided.
	// +optional
	Inline map[string]string `json:"map,omitempty"`
}

type SecureParameters struct {
	// Secret contains the parameters that will be made available to the job
	// executing the KeptnTask via the 'SECRET_DATA' environment variable.
	// The 'SECRET_DATA'  environment variable's content will the same as value of the 'SECRET_DATA'
	// key of the referenced secret.
	// +optional
	Secret string `json:"secret,omitempty"`
}

// KeptnTaskStatus defines the observed state of KeptnTask
type KeptnTaskStatus struct {
	// JobName is the name of the Job executing the Task.
	// +optional
	JobName string `json:"jobName,omitempty"`
	// Status represents the overall state of the KeptnTask.
	// +kubebuilder:default:=Pending
	// +optional
	Status common.KeptnState `json:"status,omitempty"`
	// Message contains information about unexpected errors encountered during the execution of the KeptnTask.
	// +optional
	Message string `json:"message,omitempty"`
	// StartTime represents the time at which the KeptnTask started.
	// +optional
	StartTime metav1.Time `json:"startTime,omitempty"`
	// EndTime represents the time at which the KeptnTask finished.
	// +optional
	EndTime metav1.Time `json:"endTime,omitempty"`
	// Reason contains more information about the reason for the last transition of the Job executing the KeptnTask.
	// +optional
	Reason string `json:"reason,omitempty"`
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

	// Spec describes the desired state of the KeptnTask.
	// +optional
	Spec KeptnTaskSpec `json:"spec,omitempty"`
	// Status describes the current state of the KeptnTask.
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
