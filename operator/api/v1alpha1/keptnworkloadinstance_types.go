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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// KeptnWorkloadInstanceSpec defines the desired state of KeptnWorkloadInstance
type KeptnWorkloadInstanceSpec struct {
	KeptnWorkloadSpec `json:",inline"`
	WorkloadName      string `json:"workloadName"`
}

// KeptnWorkloadInstanceStatus defines the observed state of KeptnWorkloadInstance
type KeptnWorkloadInstanceStatus struct {
	// +kubebuilder:default:=Pending
	PreDeploymentStatus common.KeptnState `json:"preDeploymentStatus,omitempty"`
	// +kubebuilder:default:=Pending
	DeploymentStatus common.KeptnState `json:"deploymentStatus,omitempty"`
	// +kubebuilder:default:=Pending
	PostDeploymentStatus     common.KeptnState    `json:"postDeploymentStatus,omitempty"`
	PreDeploymentTaskStatus  []WorkloadTaskStatus `json:"preDeploymentTaskStatus,omitempty"`
	PostDeploymentTaskStatus []WorkloadTaskStatus `json:"postDeploymentTaskStatus,omitempty"`
	StartTime                metav1.Time          `json:"startTime,omitempty"`
	EndTime                  metav1.Time          `json:"endTime,omitempty"`
}

type WorkloadTaskStatus struct {
	TaskDefinitionName string            `json:"TaskDefinitionName,omitempty"`
	Status             common.KeptnState `json:"status,omitempty"`
	TaskName           string            `json:"taskName,omitempty"`
	StartTime          metav1.Time       `json:"startTime,omitempty"`
	EndTime            metav1.Time       `json:"endTime,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="AppName",type=string,JSONPath=`.spec.app`
// +kubebuilder:printcolumn:name="Workload",type=string,JSONPath=`.spec.workloadName`
// +kubebuilder:printcolumn:name="Version",type=string,JSONPath=`.spec.version`
// +kubebuilder:printcolumn:name="PreDeploymentStatus",type=string,JSONPath=`.status.preDeploymentStatus`
// +kubebuilder:printcolumn:name="DeploymentStatus",type=string,JSONPath=`.status.deploymentStatus`
// +kubebuilder:printcolumn:name="PostDeploymentStatus",type=string,JSONPath=`.status.postDeploymentStatus`

// KeptnWorkloadInstance is the Schema for the keptnworkloadinstances API
type KeptnWorkloadInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KeptnWorkloadInstanceSpec   `json:"spec,omitempty"`
	Status KeptnWorkloadInstanceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

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

func (i KeptnWorkloadInstance) IsPostDeploymentCompleted() bool {
	return i.Status.PostDeploymentStatus.IsCompleted()
}

func (i KeptnWorkloadInstance) IsDeploymentCompleted() bool {
	return i.Status.DeploymentStatus.IsCompleted()
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

func (i *WorkloadTaskStatus) SetStartTime() {
	if i.StartTime.IsZero() {
		i.StartTime = metav1.NewTime(time.Now().UTC())
	}
}

func (i *WorkloadTaskStatus) SetEndTime() {
	if i.EndTime.IsZero() {
		i.EndTime = metav1.NewTime(time.Now().UTC())
	}
}
