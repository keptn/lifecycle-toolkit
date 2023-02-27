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
	"time"

	"github.com/keptn/lifecycle-toolkit/scheduler/test/e2e/fake/v1alpha2/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// KeptnWorkloadInstanceSpec defines the desired state of KeptnWorkloadInstance
type KeptnWorkloadInstanceSpec struct {
	KeptnWorkloadSpec `json:",inline"`
	WorkloadName      string            `json:"workloadName"`
	PreviousVersion   string            `json:"previousVersion,omitempty"`
	TraceId           map[string]string `json:"traceId,omitempty"`
}

// KeptnWorkloadInstanceStatus defines the observed state of KeptnWorkloadInstance
type KeptnWorkloadInstanceStatus struct {
	// +kubebuilder:default:=Pending
	PreDeploymentStatus common.KeptnState `json:"preDeploymentStatus,omitempty"`
	// +kubebuilder:default:=Pending
	DeploymentStatus common.KeptnState `json:"deploymentStatus,omitempty"`
	// +kubebuilder:default:=Pending
	PreDeploymentEvaluationStatus common.KeptnState `json:"preDeploymentEvaluationStatus,omitempty"`
	// +kubebuilder:default:=Pending
	PostDeploymentEvaluationStatus common.KeptnState `json:"postDeploymentEvaluationStatus,omitempty"`
	// +kubebuilder:default:=Pending
	PostDeploymentStatus               common.KeptnState `json:"postDeploymentStatus,omitempty"`
	PreDeploymentTaskStatus            []ItemStatus      `json:"preDeploymentTaskStatus,omitempty"`
	PostDeploymentTaskStatus           []ItemStatus      `json:"postDeploymentTaskStatus,omitempty"`
	PreDeploymentEvaluationTaskStatus  []ItemStatus      `json:"preDeploymentEvaluationTaskStatus,omitempty"`
	PostDeploymentEvaluationTaskStatus []ItemStatus      `json:"postDeploymentEvaluationTaskStatus,omitempty"`
	StartTime                          metav1.Time       `json:"startTime,omitempty"`
	EndTime                            metav1.Time       `json:"endTime,omitempty"`
	CurrentPhase                       string            `json:"currentPhase,omitempty"`
	// +kubebuilder:default:=Pending
	Status common.KeptnState `json:"status,omitempty"`
}

type ItemStatus struct {
	// name of EvaluationDefinition/TaskDefiniton
	DefinitionName string `json:"definitionName,omitempty"`
	// +kubebuilder:default:=Pending
	Status common.KeptnState `json:"status,omitempty"`
	// name of Evaluation/Task
	Name      string      `json:"name,omitempty"`
	StartTime metav1.Time `json:"startTime,omitempty"`
	EndTime   metav1.Time `json:"endTime,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=keptnworkloadinstances,shortName=kwi
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

// KeptnWorkloadInstance is the Schema for the keptnworkloadinstances API
type KeptnWorkloadInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KeptnWorkloadInstanceSpec   `json:"spec,omitempty"`
	Status KeptnWorkloadInstanceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KeptnWorkloadInstanceList contains a list of KeptnWorkloadInstance
type KeptnWorkloadInstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KeptnWorkloadInstance `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KeptnWorkloadInstance{}, &KeptnWorkloadInstanceList{})
}

func (i KeptnWorkloadInstance) IsPreDeploymentCompleted() bool {
	return i.Status.PreDeploymentStatus.IsCompleted()
}

func (v KeptnWorkloadInstance) IsPreDeploymentEvaluationCompleted() bool {
	return v.Status.PreDeploymentEvaluationStatus.IsCompleted()
}

func (i KeptnWorkloadInstance) IsPreDeploymentSucceeded() bool {
	return i.Status.PreDeploymentStatus.IsSucceeded()
}

func (i KeptnWorkloadInstance) IsPreDeploymentFailed() bool {
	return i.Status.PreDeploymentStatus.IsFailed()
}

func (v KeptnWorkloadInstance) IsPreDeploymentEvaluationSucceeded() bool {
	return v.Status.PreDeploymentEvaluationStatus.IsSucceeded()
}

func (v KeptnWorkloadInstance) IsPreDeploymentEvaluationFailed() bool {
	return v.Status.PreDeploymentEvaluationStatus.IsFailed()
}

func (i KeptnWorkloadInstance) IsPostDeploymentCompleted() bool {
	return i.Status.PostDeploymentStatus.IsCompleted()
}

func (v KeptnWorkloadInstance) IsPostDeploymentEvaluationCompleted() bool {
	return v.Status.PostDeploymentEvaluationStatus.IsCompleted()
}

func (i KeptnWorkloadInstance) IsPostDeploymentSucceeded() bool {
	return i.Status.PostDeploymentStatus.IsSucceeded()
}

func (i KeptnWorkloadInstance) IsPostDeploymentFailed() bool {
	return i.Status.PostDeploymentStatus.IsFailed()
}

func (v KeptnWorkloadInstance) IsPostDeploymentEvaluationSucceeded() bool {
	return v.Status.PostDeploymentEvaluationStatus.IsSucceeded()
}

func (v KeptnWorkloadInstance) IsPostDeploymentEvaluationFailed() bool {
	return v.Status.PostDeploymentEvaluationStatus.IsFailed()
}

func (i KeptnWorkloadInstance) IsDeploymentCompleted() bool {
	return i.Status.DeploymentStatus.IsCompleted()
}

func (i KeptnWorkloadInstance) IsDeploymentSucceeded() bool {
	return i.Status.DeploymentStatus.IsSucceeded()
}

func (i KeptnWorkloadInstance) IsDeploymentFailed() bool {
	return i.Status.DeploymentStatus.IsFailed()
}

func (i *KeptnWorkloadInstance) SetStartTime() {
	if i.Status.StartTime.IsZero() {
		i.Status.StartTime = metav1.NewTime(time.Now().UTC())
	}
}

func (i *KeptnWorkloadInstance) SetEndTime() {
	if i.Status.EndTime.IsZero() {
		i.Status.EndTime = metav1.NewTime(time.Now().UTC())
	}
}

func (i *KeptnWorkloadInstance) IsStartTimeSet() bool {
	return !i.Status.StartTime.IsZero()
}

func (i *KeptnWorkloadInstance) IsEndTimeSet() bool {
	return !i.Status.EndTime.IsZero()
}

func (i *ItemStatus) SetStartTime() {
	if i.StartTime.IsZero() {
		i.StartTime = metav1.NewTime(time.Now().UTC())
	}
}

func (i *ItemStatus) SetEndTime() {
	if i.EndTime.IsZero() {
		i.EndTime = metav1.NewTime(time.Now().UTC())
	}
}
