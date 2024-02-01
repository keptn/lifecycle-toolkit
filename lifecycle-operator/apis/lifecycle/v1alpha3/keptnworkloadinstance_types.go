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
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// KeptnWorkloadInstanceSpec defines the desired state of KeptnWorkloadInstance
type KeptnWorkloadInstanceSpec struct {
	KeptnWorkloadSpec `json:",inline"`
	// WorkloadName is the name of the KeptnWorkload.
	WorkloadName string `json:"workloadName"`
	// PreviousVersion is the version of the KeptnWorkload that has been deployed prior to this version.
	// +optional
	PreviousVersion string `json:"previousVersion,omitempty"`
	// TraceId contains the OpenTelemetry trace ID.
	// +optional
	TraceId map[string]string `json:"traceId,omitempty"`
}

// KeptnWorkloadInstanceStatus defines the observed state of KeptnWorkloadInstance
type KeptnWorkloadInstanceStatus struct {
	// PreDeploymentStatus indicates the current status of the KeptnWorkloadInstance's PreDeployment phase.
	// +kubebuilder:default:=Pending
	// +optional
	PreDeploymentStatus common.KeptnState `json:"preDeploymentStatus,omitempty"`
	// DeploymentStatus indicates the current status of the KeptnWorkloadInstance's Deployment phase.
	// +kubebuilder:default:=Pending
	// +optional
	DeploymentStatus common.KeptnState `json:"deploymentStatus,omitempty"`
	// PreDeploymentEvaluationStatus indicates the current status of the KeptnWorkloadInstance's PreDeploymentEvaluation phase.
	// +kubebuilder:default:=Pending
	// +optional
	PreDeploymentEvaluationStatus common.KeptnState `json:"preDeploymentEvaluationStatus,omitempty"`
	// PostDeploymentEvaluationStatus indicates the current status of the KeptnWorkloadInstance's PostDeploymentEvaluation phase.
	// +kubebuilder:default:=Pending
	// +optional
	PostDeploymentEvaluationStatus common.KeptnState `json:"postDeploymentEvaluationStatus,omitempty"`
	// PostDeploymentStatus indicates the current status of the KeptnWorkloadInstance's PostDeployment phase.
	// +kubebuilder:default:=Pending
	// +optional
	PostDeploymentStatus common.KeptnState `json:"postDeploymentStatus,omitempty"`
	// PreDeploymentTaskStatus indicates the current state of each preDeploymentTask of the KeptnWorkloadInstance.
	// +optional
	PreDeploymentTaskStatus []ItemStatus `json:"preDeploymentTaskStatus,omitempty"`
	// PostDeploymentTaskStatus indicates the current state of each postDeploymentTask of the KeptnWorkloadInstance.
	// +optional
	PostDeploymentTaskStatus []ItemStatus `json:"postDeploymentTaskStatus,omitempty"`
	// PreDeploymentEvaluationTaskStatus indicates the current state of each preDeploymentEvaluation of the KeptnWorkloadInstance.
	// +optional
	PreDeploymentEvaluationTaskStatus []ItemStatus `json:"preDeploymentEvaluationTaskStatus,omitempty"`
	// PostDeploymentEvaluationTaskStatus indicates the current state of each postDeploymentEvaluation of the KeptnWorkloadInstance.
	// +optional
	PostDeploymentEvaluationTaskStatus []ItemStatus `json:"postDeploymentEvaluationTaskStatus,omitempty"`
	// StartTime represents the time at which the deployment of the KeptnWorkloadInstance started.
	// +optional
	StartTime metav1.Time `json:"startTime,omitempty"`
	// EndTime represents the time at which the deployment of the KeptnWorkloadInstance finished.
	// +optional
	EndTime metav1.Time `json:"endTime,omitempty"`
	// CurrentPhase indicates the current phase of the KeptnWorkloadInstance. This can be:
	// - PreDeploymentTasks
	// - PreDeploymentEvaluations
	// - Deployment
	// - PostDeploymentTasks
	// - PostDeploymentEvaluations
	// +optional
	CurrentPhase string `json:"currentPhase,omitempty"`
	// PhaseTraceIDs contains the trace IDs of the OpenTelemetry spans of each phase of the KeptnWorkloadInstance
	// +optional
	PhaseTraceIDs common.PhaseTraceID `json:"phaseTraceIDs,omitempty"`
	// Status represents the overall status of the KeptnWorkloadInstance.
	// +kubebuilder:default:=Pending
	// +optional
	Status common.KeptnState `json:"status,omitempty"`
}

type ItemStatus struct {
	// DefinitionName is the name of the EvaluationDefinition/TaskDefinition
	// +optional
	DefinitionName string `json:"definitionName,omitempty"`
	// +kubebuilder:default:=Pending
	// +optional
	Status common.KeptnState `json:"status,omitempty"`
	// Name is the name of the Evaluation/Task
	// +optional
	Name string `json:"name,omitempty"`
	// StartTime represents the time at which the Item (Evaluation/Task) started.
	// +optional
	StartTime metav1.Time `json:"startTime,omitempty"`
	// EndTime represents the time at which the Item (Evaluation/Task) started.
	// +optional
	EndTime metav1.Time `json:"endTime,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=keptnworkloadinstances,shortName=kwi
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:printcolumn:name="AppName",type=string,JSONPath=`.spec.app`
// +kubebuilder:printcolumn:name="WorkloadName",type=string,JSONPath=`.spec.workloadName`
// +kubebuilder:printcolumn:name="WorkloadVersion",type=string,JSONPath=`.spec.version`
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.currentPhase`
// +kubebuilder:printcolumn:name="PreDeploymentStatus",priority=1,type=string,JSONPath=`.status.preDeploymentStatus`
// +kubebuilder:printcolumn:name="PreDeploymentEvaluationStatus",priority=1,type=string,JSONPath=`.status.preDeploymentEvaluationStatus`
// +kubebuilder:printcolumn:name="DeploymentStatus",type=string,priority=1,JSONPath=`.status.deploymentStatus`
// +kubebuilder:printcolumn:name="PostDeploymentStatus",type=string,priority=1,JSONPath=`.status.postDeploymentStatus`
// +kubebuilder:printcolumn:name="PostDeploymentEvaluationStatus",priority=1,type=string,JSONPath=`.status.postDeploymentEvaluationStatus`

// KeptnWorkloadInstance is the Schema for the keptnworkloadinstances API
type KeptnWorkloadInstance struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec describes the desired state of the KeptnWorkloadInstance.
	// +optional
	Spec KeptnWorkloadInstanceSpec `json:"spec,omitempty"`
	// Status describes the current state of the KeptnWorkloadInstance.
	// +optional
	Status KeptnWorkloadInstanceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KeptnWorkloadInstanceList contains a list of KeptnWorkloadInstance
type KeptnWorkloadInstanceList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KeptnWorkloadInstance `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KeptnWorkloadInstance{}, &KeptnWorkloadInstanceList{})
}
