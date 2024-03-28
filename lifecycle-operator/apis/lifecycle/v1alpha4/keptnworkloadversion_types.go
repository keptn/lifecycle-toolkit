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

package v1alpha4

import (
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KeptnWorkloadVersionSpec defines the desired state of KeptnWorkloadVersion
type KeptnWorkloadVersionSpec struct {
	v1alpha3.KeptnWorkloadSpec `json:",inline"`
	// WorkloadName is the name of the KeptnWorkload.
	WorkloadName string `json:"workloadName"`
	// PreviousVersion is the version of the KeptnWorkload that has been deployed prior to this version.
	// +optional
	PreviousVersion string `json:"previousVersion,omitempty"`
	// TraceId contains the OpenTelemetry trace ID.
	// +optional
	TraceId map[string]string `json:"traceId,omitempty"`
}

// KeptnWorkloadVersionStatus defines the observed state of KeptnWorkloadVersion
type KeptnWorkloadVersionStatus struct {
	// PreDeploymentStatus indicates the current status of the KeptnWorkloadVersion's PreDeployment phase.
	// +kubebuilder:default:=Pending
	// +optional
	PreDeploymentStatus common.KeptnState `json:"preDeploymentStatus,omitempty"`
	// DeploymentStatus indicates the current status of the KeptnWorkloadVersion's Deployment phase.
	// +kubebuilder:default:=Pending
	// +optional
	DeploymentStatus common.KeptnState `json:"deploymentStatus,omitempty"`
	// PreDeploymentEvaluationStatus indicates the current status of the KeptnWorkloadVersion's PreDeploymentEvaluation phase.
	// +kubebuilder:default:=Pending
	// +optional
	PreDeploymentEvaluationStatus common.KeptnState `json:"preDeploymentEvaluationStatus,omitempty"`
	// PostDeploymentEvaluationStatus indicates the current status of the KeptnWorkloadVersion's PostDeploymentEvaluation phase.
	// +kubebuilder:default:=Pending
	// +optional
	PostDeploymentEvaluationStatus common.KeptnState `json:"postDeploymentEvaluationStatus,omitempty"`
	// PostDeploymentStatus indicates the current status of the KeptnWorkloadVersion's PostDeployment phase.
	// +kubebuilder:default:=Pending
	// +optional
	PostDeploymentStatus common.KeptnState `json:"postDeploymentStatus,omitempty"`
	// PreDeploymentTaskStatus indicates the current state of each preDeploymentTask of the KeptnWorkloadVersion.
	// +optional
	PreDeploymentTaskStatus []v1alpha3.ItemStatus `json:"preDeploymentTaskStatus,omitempty"`
	// PostDeploymentTaskStatus indicates the current state of each postDeploymentTask of the KeptnWorkloadVersion.
	// +optional
	PostDeploymentTaskStatus []v1alpha3.ItemStatus `json:"postDeploymentTaskStatus,omitempty"`
	// PreDeploymentEvaluationTaskStatus indicates the current state of each preDeploymentEvaluation of the KeptnWorkloadVersion.
	// +optional
	PreDeploymentEvaluationTaskStatus []v1alpha3.ItemStatus `json:"preDeploymentEvaluationTaskStatus,omitempty"`
	// PostDeploymentEvaluationTaskStatus indicates the current state of each postDeploymentEvaluation of the KeptnWorkloadVersion.
	// +optional
	PostDeploymentEvaluationTaskStatus []v1alpha3.ItemStatus `json:"postDeploymentEvaluationTaskStatus,omitempty"`
	// StartTime represents the time at which the deployment of the KeptnWorkloadVersion started.
	// +optional
	StartTime metav1.Time `json:"startTime,omitempty"`
	// EndTime represents the time at which the deployment of the KeptnWorkloadVersion finished.
	// +optional
	EndTime metav1.Time `json:"endTime,omitempty"`
	// CurrentPhase indicates the current phase of the KeptnWorkloadVersion. This can be:
	// - PreDeploymentTasks
	// - PreDeploymentEvaluations
	// - Deployment
	// - PostDeploymentTasks
	// - PostDeploymentEvaluations
	// +optional
	CurrentPhase string `json:"currentPhase,omitempty"`
	// PhaseTraceIDs contains the trace IDs of the OpenTelemetry spans of each phase of the KeptnWorkloadVersion
	// +optional
	PhaseTraceIDs common.PhaseTraceID `json:"phaseTraceIDs,omitempty"`
	// Status represents the overall status of the KeptnWorkloadVersion.
	// +kubebuilder:default:=Pending
	// +optional
	Status common.KeptnState `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=keptnworkloadversions,shortName=kwv
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="AppName",type=string,JSONPath=`.spec.app`
// +kubebuilder:printcolumn:name="WorkloadName",type=string,JSONPath=`.spec.workloadName`
// +kubebuilder:printcolumn:name="WorkloadVersion",type=string,JSONPath=`.spec.version`
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.currentPhase`
// +kubebuilder:printcolumn:name="PreDeploymentStatus",priority=1,type=string,JSONPath=`.status.preDeploymentStatus`
// +kubebuilder:printcolumn:name="PreDeploymentEvaluationStatus",priority=1,type=string,JSONPath=`.status.preDeploymentEvaluationStatus`
// +kubebuilder:printcolumn:name="DeploymentStatus",type=string,priority=1,JSONPath=`.status.deploymentStatus`
// +kubebuilder:printcolumn:name="PostDeploymentStatus",type=string,priority=1,JSONPath=`.status.postDeploymentStatus`
// +kubebuilder:printcolumn:name="PostDeploymentEvaluationStatus",priority=1,type=string,JSONPath=`.status.postDeploymentEvaluationStatus`

// KeptnWorkloadVersion is the Schema for the keptnworkloadversions API
type KeptnWorkloadVersion struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec describes the desired state of the KeptnWorkloadVersion.
	// +optional
	Spec KeptnWorkloadVersionSpec `json:"spec,omitempty"`
	// Status describes the current state of the KeptnWorkloadVersion.
	// +optional
	Status KeptnWorkloadVersionStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KeptnWorkloadVersionList contains a list of KeptnWorkloadVersion
type KeptnWorkloadVersionList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KeptnWorkloadVersion `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KeptnWorkloadVersion{}, &KeptnWorkloadVersionList{})
}
