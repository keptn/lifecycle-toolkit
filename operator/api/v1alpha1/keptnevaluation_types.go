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

// KeptnEvaluationSpec defines the desired state of KeptnEvaluation
type KeptnEvaluationSpec struct {
	Workload             string `json:"workload,omitempty"`
	WorkloadVersion      string `json:"workloadVersion"`
	AppName              string `json:"appName,omitempty"`
	AppVersion           string `json:"appVersion,omitempty"`
	EvaluationDefinition string `json:"evaluationDefinition"`
	// +kubebuilder:default:=10
	Retries int `json:"retries,omitempty"`
	// +optional
	// +kubebuilder:default:="5s"
	// +kubebuilder:validation:Pattern="^0|([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$"
	// +kubebuilder:validation:Type:=string
	// +optional
	RetryInterval metav1.Duration  `json:"retryInterval,omitempty"`
	FailAction    string           `json:"failAction,omitempty"`
	Type          common.CheckType `json:"checkType,omitempty"`
}

// KeptnEvaluationStatus defines the observed state of KeptnEvaluation
type KeptnEvaluationStatus struct {
	// +kubebuilder:default:=0
	RetryCount       int                             `json:"retryCount"`
	EvaluationStatus map[string]EvaluationStatusItem `json:"evaluationStatus"`
	// +kubebuilder:default:=Pending
	OverallStatus common.KeptnState `json:"overallStatus"`
	StartTime     metav1.Time       `json:"startTime,omitempty"`
	EndTime       metav1.Time       `json:"endTime,omitempty"`
}

type EvaluationStatusItem struct {
	Value   string            `json:"value"`
	Status  common.KeptnState `json:"status"`
	Message string            `json:"message,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
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

	Spec   KeptnEvaluationSpec   `json:"spec,omitempty"`
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

func (i *KeptnEvaluation) SetStartTime() {
	if i.Status.StartTime.IsZero() {
		i.Status.StartTime = metav1.NewTime(time.Now().UTC())
	}
}

func (i *KeptnEvaluation) SetEndTime() {
	if i.Status.EndTime.IsZero() {
		i.Status.EndTime = metav1.NewTime(time.Now().UTC())
	}
}

func (i *KeptnEvaluation) IsStartTimeSet() bool {
	return !i.Status.StartTime.IsZero()
}

func (i *KeptnEvaluation) IsEndTimeSet() bool {
	return !i.Status.EndTime.IsZero()
}

func (i KeptnEvaluation) GetActiveMetricsAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		common.AppName.String(i.Spec.AppName),
		common.AppVersion.String(i.Spec.AppVersion),
		common.WorkloadName.String(i.Spec.Workload),
		common.WorkloadVersion.String(i.Spec.WorkloadVersion),
		common.EvaluationName.String(i.Name),
		common.EvaluationType.String(string(i.Spec.Type)),
	}
}

func (i KeptnEvaluation) GetMetricsAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		common.AppName.String(i.Spec.AppName),
		common.AppVersion.String(i.Spec.AppVersion),
		common.WorkloadName.String(i.Spec.Workload),
		common.WorkloadVersion.String(i.Spec.WorkloadVersion),
		common.EvaluationName.String(i.Name),
		common.EvaluationType.String(string(i.Spec.Type)),
		common.EvaluationStatus.String(string(i.Status.OverallStatus)),
	}
}

func (e *KeptnEvaluation) AddEvaluationStatus(objective Objective) {

	evaluationStatusItem := EvaluationStatusItem{
		Status: common.StatePending,
	}
	e.Status.EvaluationStatus[objective.Name] = evaluationStatusItem

}
