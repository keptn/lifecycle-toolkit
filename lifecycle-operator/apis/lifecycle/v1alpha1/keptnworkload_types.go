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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// KeptnWorkloadSpec defines the desired state of KeptnWorkload
type KeptnWorkloadSpec struct {
	AppName string `json:"app"`
	Version string `json:"version"`
	// +optional
	PreDeploymentTasks []string `json:"preDeploymentTasks,omitempty"`
	// +optional
	PostDeploymentTasks []string `json:"postDeploymentTasks,omitempty"`
	// +optional
	PreDeploymentEvaluations []string `json:"preDeploymentEvaluations,omitempty"`
	// +optional
	PostDeploymentEvaluations []string          `json:"postDeploymentEvaluations,omitempty"`
	ResourceReference         ResourceReference `json:"resourceReference"`
}

// KeptnWorkloadStatus defines the observed state of KeptnWorkload
type KeptnWorkloadStatus struct {
	// +optional
	CurrentVersion string `json:"currentVersion,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="AppName",type=string,JSONPath=`.spec.app`
// +kubebuilder:printcolumn:name="Version",type=string,JSONPath=`.spec.version`

// KeptnWorkload is the Schema for the keptnworkloads API
type KeptnWorkload struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +optional
	Spec KeptnWorkloadSpec `json:"spec,omitempty"`
	// +optional
	Status KeptnWorkloadStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KeptnWorkloadList contains a list of KeptnWorkload
type KeptnWorkloadList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KeptnWorkload `json:"items"`
}

// ResourceReference represents the parent resource of Workload
type ResourceReference struct {
	UID  types.UID `json:"uid"`
	Kind string    `json:"kind"`
	Name string    `json:"name"`
}

func init() {
	SchemeBuilder.Register(&KeptnWorkload{}, &KeptnWorkloadList{})
}
