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
	"go.opentelemetry.io/otel/attribute"
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
	PreDeploymentStatus common.KeptnState `json:"preDeploymentStatus,omitempty"`
	// +kubebuilder:default:=Pending
	PostDeploymentStatus common.KeptnState `json:"postDeploymentStatus,omitempty"`
	// +kubebuilder:default:=Pending
	PreDeploymentEvaluationStatus common.KeptnState `json:"preDeploymentEvaluationStatus,omitempty"`
	// +kubebuilder:default:=Pending
	PostDeploymentEvaluationStatus common.KeptnState `json:"postDeploymentEvaluationStatus,omitempty"`
	// +kubebuilder:default:=Pending
	WorkloadOverallStatus common.KeptnState `json:"workloadOverallStatus,omitempty"`
	WorkloadStatus        []WorkloadStatus  `json:"workloadStatus,omitempty"`

	PreDeploymentTaskStatus            []TaskStatus       `json:"preDeploymentTaskStatus,omitempty"`
	PostDeploymentTaskStatus           []TaskStatus       `json:"postDeploymentTaskStatus,omitempty"`
	PreDeploymentEvaluationTaskStatus  []EvaluationStatus `json:"preDeploymentEvaluationTaskStatus,omitempty"`
	PostDeploymentEvaluationTaskStatus []EvaluationStatus `json:"postDeploymentEvaluationTaskStatus,omitempty"`

	StartTime metav1.Time `json:"startTime,omitempty"`
	EndTime   metav1.Time `json:"endTime,omitempty"`
}

type WorkloadStatus struct {
	Workload KeptnWorkloadRef `json:"workload,omitempty"`
	// +kubebuilder:default:=Pending
	Status common.KeptnState `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:path=keptnappversions,shortName=kav
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="AppName",type=string,JSONPath=`.spec.appName`
// +kubebuilder:printcolumn:name="Version",type=string,JSONPath=`.spec.version`
// +kubebuilder:printcolumn:name="PreDeploymentStatus",type=string,JSONPath=`.status.preDeploymentStatus`
// +kubebuilder:printcolumn:name="WorkloadOverallStatus",type=string,JSONPath=`.status.workloadOverallStatus`
// +kubebuilder:printcolumn:name="PostDeploymentStatus",type=string,JSONPath=`.status.postDeploymentStatus`

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

func (v KeptnAppVersion) IsPreDeploymentCompleted() bool {
	return v.Status.PreDeploymentStatus.IsCompleted()
}

func (v KeptnAppVersion) IsPreDeploymentEvaluationCompleted() bool {
	return v.Status.PreDeploymentEvaluationStatus.IsCompleted()
}

func (v KeptnAppVersion) IsPreDeploymentSucceeded() bool {
	return v.Status.PreDeploymentStatus.IsSucceeded()
}

func (v KeptnAppVersion) IsPreDeploymentFailed() bool {
	return v.Status.PreDeploymentStatus.IsFailed()
}

func (v KeptnAppVersion) IsPreDeploymentEvaluationSucceeded() bool {
	return v.Status.PreDeploymentEvaluationStatus.IsSucceeded()
}

func (v KeptnAppVersion) IsPreDeploymentEvaluationFailed() bool {
	return v.Status.PreDeploymentEvaluationStatus.IsFailed()
}

func (v KeptnAppVersion) IsPostDeploymentCompleted() bool {
	return v.Status.PostDeploymentStatus.IsCompleted()
}

func (v KeptnAppVersion) IsPostDeploymentEvaluationCompleted() bool {
	return v.Status.PostDeploymentEvaluationStatus.IsCompleted()
}

func (v KeptnAppVersion) IsPostDeploymentFailed() bool {
	return v.Status.PostDeploymentStatus.IsFailed()
}

func (v KeptnAppVersion) IsPostDeploymentEvaluationSucceeded() bool {
	return v.Status.PostDeploymentEvaluationStatus.IsSucceeded()
}

func (v KeptnAppVersion) IsPostDeploymentEvaluationFailed() bool {
	return v.Status.PostDeploymentEvaluationStatus.IsFailed()
}

func (v KeptnAppVersion) IsPostDeploymentSucceeded() bool {
	return v.Status.PostDeploymentStatus.IsSucceeded()
}

func (v KeptnAppVersion) AreWorkloadsCompleted() bool {
	return v.Status.WorkloadOverallStatus.IsCompleted()
}

func (v KeptnAppVersion) AreWorkloadsSucceeded() bool {
	return v.Status.WorkloadOverallStatus.IsSucceeded()
}

func (v KeptnAppVersion) AreWorkloadsFailed() bool {
	return v.Status.WorkloadOverallStatus.IsFailed()
}

func (v *KeptnAppVersion) SetStartTime() {
	if v.Status.StartTime.IsZero() {
		v.Status.StartTime = metav1.NewTime(time.Now().UTC())
	}
}

func (v *KeptnAppVersion) SetEndTime() {
	if v.Status.EndTime.IsZero() {
		v.Status.EndTime = metav1.NewTime(time.Now().UTC())
	}
}

func (v *KeptnAppVersion) IsStartTimeSet() bool {
	return !v.Status.StartTime.IsZero()
}

func (v *KeptnAppVersion) IsEndTimeSet() bool {
	return !v.Status.EndTime.IsZero()
}

func (v KeptnAppVersion) GetActiveMetricsAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		common.AppName.String(v.Spec.AppName),
		common.AppVersion.String(v.Spec.Version),
		common.AppNamespace.String(v.Namespace),
	}
}

func (v KeptnAppVersion) GetMetricsAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		common.AppName.String(v.Spec.AppName),
		common.AppVersion.String(v.Spec.Version),
		common.AppNamespace.String(v.Namespace),
		common.AppStatus.String(string(v.Status.PostDeploymentStatus)),
	}
}
