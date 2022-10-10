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

// KeptnAppVersionSpec defines the desired state of KeptnAppVersion
type KeptnAppVersionSpec struct {
	KeptnAppSpec `json:",inline"`
	AppName      string `json:"appName"`
}

// KeptnAppVersionStatus defines the observed state of KeptnAppVersion
type KeptnAppVersionStatus struct {
	// +kubebuilder:default:=Pending
	PreDeploymentStatus common.KeptnState `json:"preDeploymentStatus,omitempty"` // WLV watches this, Scheduler watches this?
	// +kubebuilder:default:=Pending
	PostDeploymentStatus common.KeptnState `json:"postDeploymentStatus,omitempty"`
	// +kubebuilder:default:=Pending
	WorkloadOverallStatus common.KeptnState `json:"workloadOverallStatus,omitempty"`
	WorkloadStatus        []WorkloadStatus  `json:"workloadStatus,omitempty"` //Workload Instance post dep check

	PreDeploymentTaskStatus  []TaskStatus `json:"preDeploymentTaskStatus,omitempty"`
	PostDeploymentTaskStatus []TaskStatus `json:"postDeploymentTaskStatus,omitempty"`

	StartTime metav1.Time `json:"startTime,omitempty"`
	EndTime   metav1.Time `json:"endTime,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

type WorkloadStatus struct {
	Workload KeptnWorkloadRef
	status   common.KeptnState
}

// KeptnAppVersion is the Schema for the keptnappversions API
type KeptnAppVersion struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KeptnAppVersionSpec   `json:"spec,omitempty"`
	Status KeptnAppVersionStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// KeptnAppVersionList contains a list of KeptnAppVersion
type KeptnAppVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KeptnAppVersion `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KeptnAppVersion{}, &KeptnAppVersionList{})
}
