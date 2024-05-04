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

package v1alpha2

import (
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha2/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KeptnAppVersionSpec defines the desired state of KeptnAppVersion
type KeptnAppVersionSpec struct {
	KeptnAppSpec `json:",inline"`
	AppName      string `json:"appName"`
	// +optional
	PreviousVersion string `json:"previousVersion,omitempty"`

	// +optional
	TraceId map[string]string `json:"traceId,omitempty"`
}

// KeptnAppVersionStatus defines the observed state of KeptnAppVersion
type KeptnAppVersionStatus struct {
	// +kubebuilder:default:=Pending
	// +optional
	PreDeploymentStatus common.KeptnState `json:"preDeploymentStatus,omitempty"`
	// +kubebuilder:default:=Pending
	// +optional
	PostDeploymentStatus common.KeptnState `json:"postDeploymentStatus,omitempty"`
	// +kubebuilder:default:=Pending
	// +optional
	PreDeploymentEvaluationStatus common.KeptnState `json:"preDeploymentEvaluationStatus,omitempty"`
	// +kubebuilder:default:=Pending
	// +optional
	PostDeploymentEvaluationStatus common.KeptnState `json:"postDeploymentEvaluationStatus,omitempty"`
	// +kubebuilder:default:=Pending
	// +optional
	WorkloadOverallStatus common.KeptnState `json:"workloadOverallStatus,omitempty"`
	// +optional
	WorkloadStatus []WorkloadStatus `json:"workloadStatus,omitempty"`
	// +optional
	CurrentPhase string `json:"currentPhase,omitempty"`
	// +optional
	PreDeploymentTaskStatus []ItemStatus `json:"preDeploymentTaskStatus,omitempty"`
	// +optional
	PostDeploymentTaskStatus []ItemStatus `json:"postDeploymentTaskStatus,omitempty"`
	// +optional
	PreDeploymentEvaluationTaskStatus []ItemStatus `json:"preDeploymentEvaluationTaskStatus,omitempty"`
	// +optional
	PostDeploymentEvaluationTaskStatus []ItemStatus `json:"postDeploymentEvaluationTaskStatus,omitempty"`
	// +optional
	PhaseTraceIDs common.PhaseTraceID `json:"phaseTraceIDs,omitempty"`
	// +kubebuilder:default:=Pending
	// +optional
	Status common.KeptnState `json:"status,omitempty"`

	// +optional
	StartTime metav1.Time `json:"startTime,omitempty"`
	// +optional
	EndTime metav1.Time `json:"endTime,omitempty"`
}

type WorkloadStatus struct {
	// +optional
	Workload KeptnWorkloadRef `json:"workload,omitempty"`
	// +kubebuilder:default:=Pending
	// +optional
	Status common.KeptnState `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=keptnappversions,shortName=kav
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="AppName",type=string,JSONPath=`.spec.appName`
// +kubebuilder:printcolumn:name="Version",type=string,JSONPath=`.spec.version`
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.currentPhase`
// +kubebuilder:printcolumn:name="PreDeploymentStatus",priority=1,type=string,JSONPath=`.status.preDeploymentStatus`
// +kubebuilder:printcolumn:name="PreDeploymentEvaluationStatus",priority=1,type=string,JSONPath=`.status.preDeploymentEvaluationStatus`
// +kubebuilder:printcolumn:name="WorkloadOverallStatus",priority=1,type=string,JSONPath=`.status.workloadOverallStatus`
// +kubebuilder:printcolumn:name="PostDeploymentStatus",priority=1,type=string,JSONPath=`.status.postDeploymentStatus`
// +kubebuilder:printcolumn:name="PostDeploymentEvaluationStatus",priority=1,type=string,JSONPath=`.status.postDeploymentEvaluationStatus`

// KeptnAppVersion is the Schema for the keptnappversions API
type KeptnAppVersion struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +optional
	Spec KeptnAppVersionSpec `json:"spec,omitempty"`
	// +optional
	Status KeptnAppVersionStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KeptnAppVersionList contains a list of KeptnAppVersion
type KeptnAppVersionList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KeptnAppVersion `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KeptnAppVersion{}, &KeptnAppVersionList{})
}
