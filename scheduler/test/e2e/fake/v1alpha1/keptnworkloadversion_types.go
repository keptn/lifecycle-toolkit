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

	"github.com/keptn/lifecycle-toolkit/scheduler/test/e2e/fake/v1alpha1/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// KeptnWorkloadVersionSpec defines the desired state of KeptnWorkloadVersion
type KeptnWorkloadVersionSpec struct {
	KeptnWorkloadSpec `json:",inline"`
	WorkloadName      string            `json:"workloadName"`
	PreviousVersion   string            `json:"previousVersion,omitempty"`
	TraceId           map[string]string `json:"traceId,omitempty"`
}

// KeptnWorkloadVersionStatus defines the observed state of KeptnWorkloadVersion
type KeptnWorkloadVersionStatus struct {
	// +kubebuilder:default:=Pending
	PreDeploymentStatus common.KeptnState `json:"preDeploymentStatus,omitempty"`
	// +kubebuilder:default:=Pending
	DeploymentStatus common.KeptnState `json:"deploymentStatus,omitempty"`
	// +kubebuilder:default:=Pending
	PreDeploymentEvaluationStatus common.KeptnState `json:"preDeploymentEvaluationStatus,omitempty"`
	// +kubebuilder:default:=Pending
	PostDeploymentEvaluationStatus common.KeptnState `json:"postDeploymentEvaluationStatus,omitempty"`
	// +kubebuilder:default:=Pending
	PostDeploymentStatus               common.KeptnState  `json:"postDeploymentStatus,omitempty"`
	PreDeploymentTaskStatus            []TaskStatus       `json:"preDeploymentTaskStatus,omitempty"`
	PostDeploymentTaskStatus           []TaskStatus       `json:"postDeploymentTaskStatus,omitempty"`
	PreDeploymentEvaluationTaskStatus  []EvaluationStatus `json:"preDeploymentEvaluationTaskStatus,omitempty"`
	PostDeploymentEvaluationTaskStatus []EvaluationStatus `json:"postDeploymentEvaluationTaskStatus,omitempty"`
	StartTime                          metav1.Time        `json:"startTime,omitempty"`
	EndTime                            metav1.Time        `json:"endTime,omitempty"`
	CurrentPhase                       string             `json:"currentPhase,omitempty"`
	// +kubebuilder:default:=Pending
	Status common.KeptnState `json:"status,omitempty"`
}

type TaskStatus struct {
	TaskDefinitionName string `json:"taskDefinitionName,omitempty"`
	// +kubebuilder:default:=Pending
	Status    common.KeptnState `json:"status,omitempty"`
	TaskName  string            `json:"taskName,omitempty"`
	StartTime metav1.Time       `json:"startTime,omitempty"`
	EndTime   metav1.Time       `json:"endTime,omitempty"`
}

type EvaluationStatus struct {
	EvaluationDefinitionName string `json:"evaluationDefinitionName,omitempty"`
	// +kubebuilder:default:=Pending
	Status         common.KeptnState `json:"status,omitempty"`
	EvaluationName string            `json:"evaluationName,omitempty"`
	StartTime      metav1.Time       `json:"startTime,omitempty"`
	EndTime        metav1.Time       `json:"endTime,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=keptnworkloadversions,shortName=kwi
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
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KeptnWorkloadVersionSpec   `json:"spec,omitempty"`
	Status KeptnWorkloadVersionStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KeptnWorkloadVersionList contains a list of KeptnWorkloadVersion
type KeptnWorkloadVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KeptnWorkloadVersion `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KeptnWorkloadVersion{}, &KeptnWorkloadVersionList{})
}

func (i KeptnWorkloadVersion) IsPreDeploymentCompleted() bool {
	return i.Status.PreDeploymentStatus.IsCompleted()
}

func (v KeptnWorkloadVersion) IsPreDeploymentEvaluationCompleted() bool {
	return v.Status.PreDeploymentEvaluationStatus.IsCompleted()
}

func (i KeptnWorkloadVersion) IsPreDeploymentSucceeded() bool {
	return i.Status.PreDeploymentStatus.IsSucceeded()
}

func (i KeptnWorkloadVersion) IsPreDeploymentFailed() bool {
	return i.Status.PreDeploymentStatus.IsFailed()
}

func (v KeptnWorkloadVersion) IsPreDeploymentEvaluationSucceeded() bool {
	return v.Status.PreDeploymentEvaluationStatus.IsSucceeded()
}

func (v KeptnWorkloadVersion) IsPreDeploymentEvaluationFailed() bool {
	return v.Status.PreDeploymentEvaluationStatus.IsFailed()
}

func (i KeptnWorkloadVersion) IsPostDeploymentCompleted() bool {
	return i.Status.PostDeploymentStatus.IsCompleted()
}

func (v KeptnWorkloadVersion) IsPostDeploymentEvaluationCompleted() bool {
	return v.Status.PostDeploymentEvaluationStatus.IsCompleted()
}

func (i KeptnWorkloadVersion) IsPostDeploymentSucceeded() bool {
	return i.Status.PostDeploymentStatus.IsSucceeded()
}

func (i KeptnWorkloadVersion) IsPostDeploymentFailed() bool {
	return i.Status.PostDeploymentStatus.IsFailed()
}

func (v KeptnWorkloadVersion) IsPostDeploymentEvaluationSucceeded() bool {
	return v.Status.PostDeploymentEvaluationStatus.IsSucceeded()
}

func (v KeptnWorkloadVersion) IsPostDeploymentEvaluationFailed() bool {
	return v.Status.PostDeploymentEvaluationStatus.IsFailed()
}

func (i KeptnWorkloadVersion) IsDeploymentCompleted() bool {
	return i.Status.DeploymentStatus.IsCompleted()
}

func (i KeptnWorkloadVersion) IsDeploymentSucceeded() bool {
	return i.Status.DeploymentStatus.IsSucceeded()
}

func (i KeptnWorkloadVersion) IsDeploymentFailed() bool {
	return i.Status.DeploymentStatus.IsFailed()
}

func (i *KeptnWorkloadVersion) SetStartTime() {
	if i.Status.StartTime.IsZero() {
		i.Status.StartTime = metav1.NewTime(time.Now().UTC())
	}
}

func (i *KeptnWorkloadVersion) SetEndTime() {
	if i.Status.EndTime.IsZero() {
		i.Status.EndTime = metav1.NewTime(time.Now().UTC())
	}
}

func (i *KeptnWorkloadVersion) IsStartTimeSet() bool {
	return !i.Status.StartTime.IsZero()
}

func (i *KeptnWorkloadVersion) IsEndTimeSet() bool {
	return !i.Status.EndTime.IsZero()
}

func (i *TaskStatus) SetStartTime() {
	if i.StartTime.IsZero() {
		i.StartTime = metav1.NewTime(time.Now().UTC())
	}
}

func (i *TaskStatus) SetEndTime() {
	if i.EndTime.IsZero() {
		i.EndTime = metav1.NewTime(time.Now().UTC())
	}
}

func (i *EvaluationStatus) SetStartTime() {
	if i.StartTime.IsZero() {
		i.StartTime = metav1.NewTime(time.Now().UTC())
	}
}

func (i *EvaluationStatus) SetEndTime() {
	if i.EndTime.IsZero() {
		i.EndTime = metav1.NewTime(time.Now().UTC())
	}
}
