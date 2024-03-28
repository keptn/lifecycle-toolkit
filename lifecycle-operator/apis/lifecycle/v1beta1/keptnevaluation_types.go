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

// KeptnEvaluationSpec defines the desired state of KeptnEvaluation
type KeptnEvaluationSpec struct {
	// Workload defines the KeptnWorkload for which the KeptnEvaluation is done.
	// +optional
	Workload string `json:"workload,omitempty"`
	// WorkloadVersion defines the version of the KeptnWorkload for which the KeptnEvaluation is done.
	WorkloadVersion string `json:"workloadVersion"`
	// AppName defines the KeptnApp for which the KeptnEvaluation is done.
	// +optional
	AppName string `json:"appName,omitempty"`
	// AppVersion defines the version of the KeptnApp for which the KeptnEvaluation is done.
	// +optional
	AppVersion string `json:"appVersion,omitempty"`
	// EvaluationDefinition refers to the name of the KeptnEvaluationDefinition
	// which includes the objectives for the KeptnEvaluation.
	// The KeptnEvaluationDefinition can be
	// located in the same namespace as the KeptnEvaluation, or in the Keptn namespace.
	EvaluationDefinition string `json:"evaluationDefinition"`
	// Type indicates whether the KeptnEvaluation is part of the pre- or postDeployment phase.
	// +optional
	Type common.CheckType `json:"checkType,omitempty"`
	// FailureConditions represent the failure conditions (number of retries and retry interval)
	// for the evaluation to be considered as failed
	FailureConditions `json:",inline"`
}

// KeptnEvaluationStatus defines the observed state of KeptnEvaluation
type KeptnEvaluationStatus struct {
	// RetryCount indicates how many times the KeptnEvaluation has been attempted already.
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
	// StartTime represents the time at which the KeptnEvaluation started.
	// +optional
	StartTime metav1.Time `json:"startTime,omitempty"`
	// EndTime represents the time at which the KeptnEvaluation finished.
	// +optional
	EndTime metav1.Time `json:"endTime,omitempty"`
}

type EvaluationStatusItem struct {
	// Value represents the value of the KeptnMetric being evaluated.
	Value string `json:"value"`
	// Status indicates the status of the objective being evaluated.
	Status common.KeptnState `json:"status"`
	// Message contains additional information about the evaluation of an objective.
	// This can include explanations about why an evaluation has failed (e.g. due to a missed objective),
	// or if there was any error during the evaluation of the objective.
	// +optional
	Message string `json:"message,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=keptnevaluations,shortName=ke
// +kubebuilder:printcolumn:name="AppName",type=string,JSONPath=`.spec.appName`
// +kubebuilder:printcolumn:name="AppVersion",type=string,JSONPath=`.spec.appVersion`
// +kubebuilder:printcolumn:name="WorkloadName",type=string,JSONPath=`.spec.workload`
// +kubebuilder:printcolumn:name="WorkloadVersion",type=string,JSONPath=`.spec.workloadVersion`
// +kubebuilder:printcolumn:name="RetryCount",type=string,JSONPath=`.status.retryCount`
// +kubebuilder:printcolumn:name="EvaluationStatus",type=string,JSONPath=`.status.evaluationStatus`
// +kubebuilder:printcolumn:name="OverallStatus",type=string,JSONPath=`.status.overallStatus`

// KeptnEvaluation is the Schema for the keptnevaluations API
type KeptnEvaluation struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec describes the desired state of the KeptnEvaluation.
	// +optional
	Spec KeptnEvaluationSpec `json:"spec,omitempty"`
	// Status describes the current state of the KeptnEvaluation.
	// +optional
	Status KeptnEvaluationStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KeptnEvaluationList contains a list of KeptnEvaluation
type KeptnEvaluationList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KeptnEvaluation `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KeptnEvaluation{}, &KeptnEvaluationList{})
}
