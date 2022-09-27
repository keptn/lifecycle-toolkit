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
	PreDeploymentStatus      common.KeptnState    `json:"preDeploymentStatus"`
	PostDeploymentStatus     common.KeptnState    `json:"postDeploymentStatus"`
	PreDeploymentTaskStatus  []WorkloadTaskStatus `json:"preDeploymentTaskStatus"`
	PostDeploymentTaskStatus []WorkloadTaskStatus `json:"postDeploymentTaskStatus"`
}

type WorkloadTaskStatus struct {
	TaskDefinitionName string            `json:"TaskDefinitionName"`
	Status             common.KeptnState `json:"status,omitempty"`
	TaskName           string            `json:"taskName,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

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
	if i.Status.PreDeploymentStatus == common.StateSucceeded || i.Status.PreDeploymentStatus == common.StateFailed || i.Status.PreDeploymentStatus == common.StateUnknown {
		return true
	}
	return false
}

func (i KeptnWorkloadInstance) IsPostDeploymentCompleted() bool {
	if i.Status.PostDeploymentStatus == common.StateSucceeded || i.Status.PostDeploymentStatus == common.StateFailed || i.Status.PostDeploymentStatus == common.StateUnknown {
		return true
	}
	return false
}
