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

	"github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha1/common"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
	Message   string            `json:"message,omitempty"`
	StartTime metav1.Time       `json:"startTime,omitempty"`
	EndTime   metav1.Time       `json:"endTime,omitempty"`
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
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
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KeptnTaskSpec   `json:"spec,omitempty"`
	Status KeptnTaskStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KeptnTaskList contains a list of KeptnTask
type KeptnTaskList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KeptnTask `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KeptnTask{}, &KeptnTaskList{})
}

func (t KeptnTaskList) GetItems() []client.Object {
	var b []client.Object
	for _, i := range t.Items {
		b = append(b, &i)
	}
	return b
}

func (t *KeptnTask) SetStartTime() {
	if t.Status.StartTime.IsZero() {
		t.Status.StartTime = metav1.NewTime(time.Now().UTC())
	}
}

func (t *KeptnTask) SetEndTime() {
	if t.Status.EndTime.IsZero() {
		t.Status.EndTime = metav1.NewTime(time.Now().UTC())
	}
}

func (t *KeptnTask) IsStartTimeSet() bool {
	return !t.Status.StartTime.IsZero()
}

func (t *KeptnTask) IsEndTimeSet() bool {
	return !t.Status.EndTime.IsZero()
}

func (t KeptnTask) GetActiveMetricsAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		common.AppName.String(t.Spec.AppName),
		common.AppVersion.String(t.Spec.AppVersion),
		common.WorkloadName.String(t.Spec.Workload),
		common.WorkloadVersion.String(t.Spec.WorkloadVersion),
		common.TaskName.String(t.Name),
		common.TaskType.String(string(t.Spec.Type)),
	}
}

func (t KeptnTask) GetMetricsAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		common.AppName.String(t.Spec.AppName),
		common.AppVersion.String(t.Spec.AppVersion),
		common.WorkloadName.String(t.Spec.Workload),
		common.WorkloadVersion.String(t.Spec.WorkloadVersion),
		common.TaskName.String(t.Name),
		common.TaskType.String(string(t.Spec.Type)),
		common.TaskStatus.String(string(t.Status.Status)),
	}
}

func (t KeptnTask) SetSpanAttributes(span trace.Span) {
	span.SetAttributes(t.GetSpanAttributes()...)
}

func (t KeptnTask) CreateKeptnLabels() map[string]string {
	if t.Spec.Workload != "" {
		return map[string]string{
			common.AppAnnotation:      t.Spec.AppName,
			common.WorkloadAnnotation: t.Spec.Workload,
			common.VersionAnnotation:  t.Spec.WorkloadVersion,
			common.TaskNameAnnotation: t.Name,
		}
	}
	return map[string]string{
		common.AppAnnotation:      t.Spec.AppName,
		common.VersionAnnotation:  t.Spec.AppVersion,
		common.TaskNameAnnotation: t.Name,
	}
}

func (t KeptnTask) GetSpanAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		common.AppName.String(t.Spec.AppName),
		common.AppVersion.String(t.Spec.AppVersion),
		common.WorkloadName.String(t.Spec.Workload),
		common.WorkloadVersion.String(t.Spec.WorkloadVersion),
		common.TaskName.String(t.Name),
		common.TaskType.String(string(t.Spec.Type)),
	}
}

func (t *KeptnTask) SetPhaseTraceID(phase string, carrier propagation.MapCarrier) {
	// present due to SpanItem interface
}

func (t KeptnTask) GetSpanKey(phase string) string {
	return t.Name
}

func (t KeptnTask) GetSpanName(phase string) string {
	return t.Name
}
