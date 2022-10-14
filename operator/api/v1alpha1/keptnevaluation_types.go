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
	Workload             string           `json:"workload,omitempty"`
	WorkloadVersion      string           `json:"workloadVersion"`
	AppName              string           `json:"app,omitempty"`
	AppVersion           string           `json:"appVersion,omitempty"`
	EvaluationDefinition string           `json:"evaluationDefinition"`
	Source               string           `json:"source"`
	Timeframe            time.Duration    `json:"timeframe,omitempty"`
	Retries              int              `json:"retries,omitempty"`
	RetryInterval        time.Duration    `json:"retryInterval,omitempty"`
	FailAction           string           `json:"failAction,omitempty"`
	Type                 common.CheckType `json:"checkType,omitempty"`
}

// KeptnEvaluationStatus defines the observed state of KeptnEvaluation
type KeptnEvaluationStatus struct {
	RetryCount       int                    `json:"retryCount"`
	EvaluationStatus []EvaluationStatusItem `json:"evaluationStatus"`
	OverallStatus    common.KeptnState      `json:"overallStatus"`
	StartTime        metav1.Time            `json:"startTime,omitempty"`
	EndTime          metav1.Time            `json:"endTime,omitempty"`
}

type EvaluationStatusItem struct {
	Name   string            `json:"name"`
	Value  string            `json:"value"`
	Status common.KeptnState `json:"status"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:path=keptnevaluations,shortName=ke

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

func (e *KeptnEvaluation) InitializeEvaluationStatuses(definition KeptnEvaluationDefinition) {
	for _, query := range definition.Spec.Objectives {
		evaluationStatusItem := EvaluationStatusItem{
			Name:   query.Name,
			Status: common.StatePending,
		}
		e.Status.EvaluationStatus = append(e.Status.EvaluationStatus, evaluationStatusItem)
	}
}
