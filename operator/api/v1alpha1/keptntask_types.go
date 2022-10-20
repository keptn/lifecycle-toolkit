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
	"time"

	"github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1/common"
	"go.opentelemetry.io/otel/attribute"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// KeptnTaskSpec defines the desired state of KeptnTask
type KeptnTaskSpec struct {
	Workload         string           `json:"workload"`
	WorkloadVersion  string           `json:"workloadVersion"`
	AppName          string           `json:"app"`
	AppVersion       string           `json:"appVersion"`
	TaskDefinition   string           `json:"taskDefinition"`
	Context          TaskContext      `json:"context"`
	Parameters       TaskParameters   `json:"parameters,omitempty"`
	SecureParameters SecureParameters `json:"secureParameters,omitempty"`
	Type             common.CheckType `json:"checkType,omitempty"`
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
	Inline map[string]string `json:"map,omitempty"`
}

type SecureParameters struct {
	Secret string `json:"secret,omitempty"`
}

// KeptnTaskStatus defines the observed state of KeptnTask
type KeptnTaskStatus struct {
	JobName string `json:"jobName,omitempty"`
	// +kubebuilder:default:=Pending
	Status    common.KeptnState `json:"status,omitempty"`
	StartTime metav1.Time       `json:"startTime,omitempty"`
	EndTime   metav1.Time       `json:"endTime,omitempty"`
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="AppName",type=string,JSONPath=`.spec.app`
// +kubebuilder:printcolumn:name="AppVersion",type=string,JSONPath=`.spec.appVersion`
// +kubebuilder:printcolumn:name="WorkloadName",type=string,JSONPath=`.spec.workload`
// +kubebuilder:printcolumn:name="WorkloadVersion",type=string,JSONPath=`.spec.workloadVersion`
// +kubebuilder:printcolumn:name="Job Name",type=string,JSONPath=`.status.jobName`
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.status`

// KeptnTask is the Schema for the keptntasks API
type KeptnTask struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KeptnTaskSpec   `json:"spec,omitempty"`
	Status KeptnTaskStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// KeptnTaskList contains a list of KeptnTask
type KeptnTaskList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KeptnTask `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KeptnTask{}, &KeptnTaskList{})
}

func (i *KeptnTask) SetStartTime() {
	if i.Status.StartTime.IsZero() {
		i.Status.StartTime = metav1.NewTime(time.Now().UTC())
	}
}

func (i *KeptnTask) SetEndTime() {
	if i.Status.EndTime.IsZero() {
		i.Status.EndTime = metav1.NewTime(time.Now().UTC())
	}
}

func (i *KeptnTask) IsStartTimeSet() bool {
	return !i.Status.StartTime.IsZero()
}

func (i *KeptnTask) IsEndTimeSet() bool {
	return !i.Status.EndTime.IsZero()
}

func (i KeptnTask) GetActiveMetricsAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		common.AppName.String(i.Spec.AppName),
		common.AppVersion.String(i.Spec.AppVersion),
		common.WorkloadName.String(i.Spec.Workload),
		common.WorkloadVersion.String(i.Spec.WorkloadVersion),
		common.TaskName.String(i.Name),
		common.TaskType.String(string(i.Spec.Type)),
	}
}

func (i KeptnTask) GetMetricsAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		common.AppName.String(i.Spec.AppName),
		common.AppVersion.String(i.Spec.AppVersion),
		common.WorkloadName.String(i.Spec.Workload),
		common.WorkloadVersion.String(i.Spec.WorkloadVersion),
		common.TaskName.String(i.Name),
		common.TaskType.String(string(i.Spec.Type)),
		common.TaskStatus.String(string(i.Status.Status)),
	}
}
