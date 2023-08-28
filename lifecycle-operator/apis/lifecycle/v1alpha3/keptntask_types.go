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

package v1alpha3

import (
	"time"

	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
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
	// TaskDefinition refers to the name of the KeptnTaskDefinition
	// which includes the specification for the task to be performed.
	// The KeptnTaskDefinition can be
	// located in the same namespace as the KeptnTask, or in the KLT namespace.
	TaskDefinition string `json:"taskDefinition"`
	// Context contains contextual information about the task execution.
	Context TaskContext `json:"context"`
	// Parameters contains parameters that will be passed to the job that executes the task.
	Parameters TaskParameters `json:"parameters,omitempty"`
	// SecureParameters contains secure parameters that will be passed to the job that executes the task.
	// These will be stored and accessed as secrets in the cluster.
	SecureParameters SecureParameters `json:"secureParameters,omitempty"`
	// Type indicates whether the KeptnTask is part of the pre- or postDeployment phase.
	Type common.CheckType `json:"checkType,omitempty"`
	// Retries indicates how many times the KeptnTask can be attempted in the case of an error
	// before considering the KeptnTask to be failed.
	// +kubebuilder:default:=10
	Retries *int32 `json:"retries,omitempty"`
	// Timeout specifies the maximum time to wait for the task to be completed successfully.
	// If the task does not complete successfully within this time frame, it will be
	// considered to be failed.
	// +optional
	// +kubebuilder:default:="5m"
	// +kubebuilder:validation:Pattern="^0|([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	// +kubebuilder:validation:Type:=string
	// +optional
	Timeout metav1.Duration `json:"timeout,omitempty"`
}

type TaskContext struct {
	// WorkloadName the name of the KeptnWorkload the KeptnTask is being executed for.
	WorkloadName string `json:"workloadName"`
	// AppName the name of the KeptnApp the KeptnTask is being executed for.
	AppName string `json:"appName"`
	// AppVersion the version of the KeptnApp the KeptnTask is being executed for.
	AppVersion string `json:"appVersion"`
	// WorkloadVersion the version of the KeptnWorkload the KeptnTask is being executed for.
	WorkloadVersion string `json:"workloadVersion"`
	// TaskType indicates whether the KeptnTask is part of the pre- or postDeployment phase.
	TaskType string `json:"taskType"`
	// ObjectType indicates whether the KeptnTask is being executed for a KeptnApp or KeptnWorkload.
	ObjectType string `json:"objectType"`
}

type TaskParameters struct {
	// Inline contains the parameters that will be made available to the job
	// executing the KeptnTask via the 'DATA' environment variable.
	// The 'DATA'  environment variable's content will be a json
	// encoded string containing all properties of the map provided.
	Inline map[string]string `json:"map,omitempty"`
}

type SecureParameters struct {
	// Secret contains the parameters that will be made available to the job
	// executing the KeptnTask via the 'SECRET_DATA' environment variable.
	// The 'SECRET_DATA'  environment variable's content will the same as value of the 'SECRET_DATA'
	// key of the referenced secret.
	Secret string `json:"secret,omitempty"`
}

// KeptnTaskStatus defines the observed state of KeptnTask
type KeptnTaskStatus struct {
	// JobName is the name of the Job executing the Task.
	JobName string `json:"jobName,omitempty"`
	// Status represents the overall state of the KeptnTask.
	// +kubebuilder:default:=Pending
	Status common.KeptnState `json:"status,omitempty"`
	// Message contains information about unexpected errors encountered during the execution of the KeptnTask.
	Message string `json:"message,omitempty"`
	// StartTime represents the time at which the KeptnTask started.
	StartTime metav1.Time `json:"startTime,omitempty"`
	// EndTime represents the time at which the KeptnTask finished.
	EndTime metav1.Time `json:"endTime,omitempty"`
	// Reason contains more information about the reason for the last transition of the Job executing the KeptnTask.
	Reason string `json:"reason,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
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

	// Spec describes the desired state of the KeptnTask.
	Spec KeptnTaskSpec `json:"spec,omitempty"`
	// Status describes the current state of the KeptnTask.
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
		common.AppName.String(t.Spec.Context.AppName),
		common.AppVersion.String(t.Spec.Context.AppVersion),
		common.WorkloadName.String(t.Spec.Context.WorkloadName),
		common.WorkloadVersion.String(t.Spec.Context.WorkloadVersion),
		common.TaskName.String(t.Name),
		common.TaskType.String(string(t.Spec.Type)),
	}
}

func (t KeptnTask) GetMetricsAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		common.AppName.String(t.Spec.Context.AppName),
		common.AppVersion.String(t.Spec.Context.AppVersion),
		common.WorkloadName.String(t.Spec.Context.WorkloadName),
		common.WorkloadVersion.String(t.Spec.Context.WorkloadVersion),
		common.TaskName.String(t.Name),
		common.TaskType.String(string(t.Spec.Type)),
		common.TaskStatus.String(string(t.Status.Status)),
	}
}

func (t KeptnTask) SetSpanAttributes(span trace.Span) {
	span.SetAttributes(t.GetSpanAttributes()...)
}

func (t KeptnTask) CreateKeptnAnnotations() map[string]string {
	if t.Spec.Context.WorkloadName != "" {
		return common.MergeMaps(t.Annotations, map[string]string{
			common.AppAnnotation:      t.Spec.Context.AppName,
			common.WorkloadAnnotation: t.Spec.Context.WorkloadName,
			common.VersionAnnotation:  t.Spec.Context.WorkloadVersion,
			common.TaskNameAnnotation: t.Name,
		})
	}
	return common.MergeMaps(t.Annotations, map[string]string{
		common.AppAnnotation:      t.Spec.Context.AppName,
		common.VersionAnnotation:  t.Spec.Context.AppVersion,
		common.TaskNameAnnotation: t.Name,
	})
}

func (t KeptnTask) GetSpanAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		common.AppName.String(t.Spec.Context.AppName),
		common.AppVersion.String(t.Spec.Context.AppVersion),
		common.WorkloadName.String(t.Spec.Context.WorkloadName),
		common.WorkloadVersion.String(t.Spec.Context.WorkloadVersion),
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

func (t KeptnTask) GetEventAnnotations() map[string]string {
	return map[string]string{
		"appName":            t.Spec.Context.AppName,
		"appVersion":         t.Spec.Context.AppVersion,
		"workloadName":       t.Spec.Context.WorkloadName,
		"workloadVersion":    t.Spec.Context.WorkloadVersion,
		"taskName":           t.Name,
		"taskDefinitionName": t.Spec.TaskDefinition,
	}
}

func (t KeptnTask) GetActiveDeadlineSeconds() *int64 {
	deadline, _ := time.ParseDuration(t.Spec.Timeout.Duration.String())
	seconds := int64(deadline.Seconds())
	return &seconds
}
