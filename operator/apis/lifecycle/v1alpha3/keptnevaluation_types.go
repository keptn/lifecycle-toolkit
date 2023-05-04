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

	"github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3/common"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// KeptnEvaluationSpec defines the desired state of KeptnEvaluation
type KeptnEvaluationSpec struct {
	// Workload defines the KeptnWorkload for which the KeptnEvaluation is done.
	Workload string `json:"workload,omitempty"`
	// WorkloadVersion defines the version of the KeptnWorkload for which the KeptnEvaluation is done.
	WorkloadVersion string `json:"workloadVersion"`
	// AppName defines the KeptnApp for which the KeptnEvaluation is done.
	AppName string `json:"appName,omitempty"`
	// AppVersion defines the version of the KeptnApp for which the KeptnEvaluation is done.
	AppVersion string `json:"appVersion,omitempty"`
	// EvaluationDefinition refers to the name of the KeptnEvaluationDefinition
	// which includes the objectives for the KeptnEvaluation.
	// The KeptnEvaluationDefinition can be
	// located in the same namespace as the KeptnEvaluation, or in the KLT namespace
	EvaluationDefinition string `json:"evaluationDefinition"`
	// Retries indicates how many times the KeptnEvaluation can be attempted in the case of an error or
	// missed evaluation objective, before considering the KeptnEvaluation to be failed.
	// +kubebuilder:default:=10
	Retries int `json:"retries,omitempty"`
	// RetryInterval specifies the interval at which the KeptnEvaluation is retried in the case of an error
	// or a missed objective.
	// +optional
	// +kubebuilder:default:="5s"
	// +kubebuilder:validation:Pattern="^0|([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	// +kubebuilder:validation:Type:=string
	// +optional
	RetryInterval metav1.Duration `json:"retryInterval,omitempty"`
	FailAction    string          `json:"failAction,omitempty"`
	// Type indicates whether the KeptnEvaluation is part of the pre- or postDeployment phase
	Type common.CheckType `json:"checkType,omitempty"`
}

// KeptnEvaluationStatus defines the observed state of KeptnEvaluation
type KeptnEvaluationStatus struct {
	// RetryCount indicates how many times the KeptnEvaluation has been attempted already
	// +kubebuilder:default:=0
	RetryCount int `json:"retryCount"`
	// EvaluationStatus describes the status of each objective of the KeptnEvaluationDefinition
	// referenced by the KeptnEvaluation.
	EvaluationStatus map[string]EvaluationStatusItem `json:"evaluationStatus"`
	// OverallStatus describes the overall status of the KeptnEvaluation. The Overall status is derived
	// from the status of the individual objectives of the KeptnEvaluationDefinition
	// referenced by the KeptnEvaluation.
	// +kubebuilder:default:=Pending
	OverallStatus common.KeptnState `json:"overallStatus"`
	// StartTime represents the time at which the KeptnEvaluation started
	StartTime metav1.Time `json:"startTime,omitempty"`
	// EndTime represents the time at which the KeptnEvaluation finished
	EndTime metav1.Time `json:"endTime,omitempty"`
}

type EvaluationStatusItem struct {
	// Value represents the value of the KeptnMetric being evaluated
	Value string `json:"value"`
	// Status indicates the status of the objective being evaluated
	Status common.KeptnState `json:"status"`
	// Message contains additional information about the evaluation of an objective.
	// This can include explanations about why an evaluation has failed (e.g. due to a missed objective),
	// or if there was any error during the evaluation of the objective.
	Message string `json:"message,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:storageversion
//+kubebuilder:resource:path=keptnevaluations,shortName=ke
//+kubebuilder:printcolumn:name="AppName",type=string,JSONPath=`.spec.appName`
//+kubebuilder:printcolumn:name="AppVersion",type=string,JSONPath=`.spec.appVersion`
//+kubebuilder:printcolumn:name="WorkloadName",type=string,JSONPath=`.spec.workload`
//+kubebuilder:printcolumn:name="WorkloadVersion",type=string,JSONPath=`.spec.workloadVersion`
//+kubebuilder:printcolumn:name="RetryCount",type=string,JSONPath=`.status.retryCount`
//+kubebuilder:printcolumn:name="EvaluationStatus",type=string,JSONPath=`.status.evaluationStatus`
//+kubebuilder:printcolumn:name="OverallStatus",type=string,JSONPath=`.status.overallStatus`

// KeptnEvaluation is the Schema for the keptnevaluations API
type KeptnEvaluation struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec describes the desired state of the KeptnEvaluation
	Spec KeptnEvaluationSpec `json:"spec,omitempty"`
	// Status describes the current state of the KeptnEvaluation
	Status KeptnEvaluationStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// KeptnEvaluationList contains a list of KeptnEvaluation
type KeptnEvaluationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KeptnEvaluation `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KeptnEvaluation{}, &KeptnEvaluationList{})
}

func (e KeptnEvaluationList) GetItems() []client.Object {
	var b []client.Object
	for _, i := range e.Items {
		b = append(b, &i)
	}
	return b
}

func (e *KeptnEvaluation) SetStartTime() {
	if e.Status.StartTime.IsZero() {
		e.Status.StartTime = metav1.NewTime(time.Now().UTC())
	}
}

func (e *KeptnEvaluation) SetEndTime() {
	if e.Status.EndTime.IsZero() {
		e.Status.EndTime = metav1.NewTime(time.Now().UTC())
	}
}

func (e *KeptnEvaluation) IsStartTimeSet() bool {
	return !e.Status.StartTime.IsZero()
}

func (e *KeptnEvaluation) IsEndTimeSet() bool {
	return !e.Status.EndTime.IsZero()
}

func (e KeptnEvaluation) GetActiveMetricsAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		common.AppName.String(e.Spec.AppName),
		common.AppVersion.String(e.Spec.AppVersion),
		common.WorkloadName.String(e.Spec.Workload),
		common.WorkloadVersion.String(e.Spec.WorkloadVersion),
		common.EvaluationName.String(e.Name),
		common.EvaluationType.String(string(e.Spec.Type)),
	}
}

func (e KeptnEvaluation) GetMetricsAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		common.AppName.String(e.Spec.AppName),
		common.AppVersion.String(e.Spec.AppVersion),
		common.WorkloadName.String(e.Spec.Workload),
		common.WorkloadVersion.String(e.Spec.WorkloadVersion),
		common.EvaluationName.String(e.Name),
		common.EvaluationType.String(string(e.Spec.Type)),
		common.EvaluationStatus.String(string(e.Status.OverallStatus)),
	}
}

func (e *KeptnEvaluation) AddEvaluationStatus(objective Objective) {

	evaluationStatusItem := EvaluationStatusItem{
		Status: common.StatePending,
	}
	if e.Status.EvaluationStatus == nil {
		e.Status.EvaluationStatus = make(map[string]EvaluationStatusItem)
	}
	e.Status.EvaluationStatus[objective.KeptnMetricRef.Name] = evaluationStatusItem

}

func (e KeptnEvaluation) SetSpanAttributes(span trace.Span) {
	span.SetAttributes(e.GetSpanAttributes()...)
}

func (e KeptnEvaluation) GetSpanAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		common.AppName.String(e.Spec.AppName),
		common.AppVersion.String(e.Spec.AppVersion),
		common.WorkloadName.String(e.Spec.Workload),
		common.WorkloadVersion.String(e.Spec.WorkloadVersion),
		common.EvaluationName.String(e.Name),
		common.EvaluationType.String(string(e.Spec.Type)),
	}
}

func (e *KeptnEvaluation) SetPhaseTraceID(phase string, carrier propagation.MapCarrier) {
	// present due to SpanItem interface
}

func (e KeptnEvaluation) GetSpanKey(phase string) string {
	return e.Name
}

func (e KeptnEvaluation) GetSpanName(phase string) string {
	return e.Name
}

func (e KeptnEvaluation) GetEventAnnotations() map[string]string {
	return map[string]string{
		"appName":                  e.Spec.AppName,
		"appVersion":               e.Spec.AppVersion,
		"workloadName":             e.Spec.Workload,
		"workloadVersion":          e.Spec.WorkloadVersion,
		"evaluationName":           e.Name,
		"evaluationDefinitionName": e.Spec.EvaluationDefinition,
	}
}
