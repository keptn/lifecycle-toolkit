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
	"github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3/common"
	operatorcommon "github.com/keptn/lifecycle-toolkit/operator/common"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// KeptnAppSpec defines the desired state of KeptnApp
type KeptnAppSpec struct {
	// Version defines the version of the application. For automatically created KeptnApps,
	// the version is a function of all KeptnWorkloads that are part of the KeptnApp.
	Version string `json:"version"`
	// Revision can be modified to trigger another deployment of a KeptnApp of the same version.
	// This can be used for restarting a KeptnApp which failed to deploy,
	// e.g. due to a failed preDeploymentEvaluation/preDeploymentTask.
	// +kubebuilder:default:=1
	Revision uint `json:"revision,omitempty"`
	// Workloads is a list of all KeptnWorkloads that are part of the KeptnApp.
	Workloads []KeptnWorkloadRef `json:"workloads,omitempty"`
	// PreDeploymentTasks is a list of all tasks to be performed during the pre-deployment phase of the KeptnApp.
	// The items of this list refer to the names of KeptnTaskDefinitions
	// located in the same namespace as the KeptnApp, or in the KLT namespace.
	PreDeploymentTasks []string `json:"preDeploymentTasks,omitempty"`
	// PostDeploymentTasks is a list of all tasks to be performed during the post-deployment phase of the KeptnApp.
	// The items of this list refer to the names of KeptnTaskDefinitions
	// located in the same namespace as the KeptnApp, or in the KLT namespace.
	PostDeploymentTasks []string `json:"postDeploymentTasks,omitempty"`
	// PreDeploymentEvaluations is a list of all evaluations to be performed
	// during the pre-deployment phase of the KeptnApp.
	// The items of this list refer to the names of KeptnEvaluationDefinitions
	// located in the same namespace as the KeptnApp, or in the KLT namespace.
	PreDeploymentEvaluations []string `json:"preDeploymentEvaluations,omitempty"`
	// PostDeploymentEvaluations is a list of all evaluations to be performed
	// during the post-deployment phase of the KeptnApp.
	// The items of this list refer to the names of KeptnEvaluationDefinitions
	// located in the same namespace as the KeptnApp, or in the KLT namespace.
	PostDeploymentEvaluations []string `json:"postDeploymentEvaluations,omitempty"`
}

// KeptnAppStatus defines the observed state of KeptnApp
type KeptnAppStatus struct {
	// CurrentVersion indicates the version that is currently deployed or being reconciled.
	CurrentVersion string `json:"currentVersion,omitempty"`
}

// KeptnWorkloadRef refers to a KeptnWorkload that is part of a KeptnApp
type KeptnWorkloadRef struct {
	// Name is the name of the KeptnWorkload.
	Name string `json:"name"`
	// Version is the version of the KeptnWorkload.
	Version string `json:"version"`
}

//+kubebuilder:object:root=true
//+kubebuilder:storageversion
//+kubebuilder:subresource:status

// KeptnApp is the Schema for the keptnapps API
type KeptnApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec describes the desired state of the KeptnApp.
	Spec KeptnAppSpec `json:"spec,omitempty"`
	// Status describes the current state of the KeptnApp.
	Status KeptnAppStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// KeptnAppList contains a list of KeptnApp
type KeptnAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KeptnApp `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KeptnApp{}, &KeptnAppList{})
}

func (a KeptnApp) GetAppVersionName() string {
	return operatorcommon.CreateResourceName(common.MaxK8sObjectLength, common.MinKLTNameLen, a.Name, a.Spec.Version, common.Hash(a.Generation))
}

func (a KeptnApp) SetSpanAttributes(span trace.Span) {
	span.SetAttributes(a.GetSpanAttributes()...)
}

func (a KeptnApp) GenerateAppVersion(previousVersion string, traceContextCarrier map[string]string) KeptnAppVersion {
	return KeptnAppVersion{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: traceContextCarrier,
			Name:        a.GetAppVersionName(),
			Namespace:   a.Namespace,
		},
		Spec: KeptnAppVersionSpec{
			KeptnAppSpec:    a.Spec,
			AppName:         a.Name,
			PreviousVersion: previousVersion,
		},
	}
}

func (a KeptnApp) GetSpanAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		common.AppName.String(a.Name),
		common.AppVersion.String(a.Spec.Version),
	}
}

func (a KeptnApp) GetEventAnnotations() map[string]string {
	return map[string]string{
		"appName":     a.Name,
		"appVersion":  a.Spec.Version,
		"appRevision": common.Hash(a.Generation),
	}
}
