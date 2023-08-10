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
	"fmt"
	"strings"

	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	operatorcommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/common"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// KeptnWorkloadSpec defines the desired state of KeptnWorkload
type KeptnWorkloadSpec struct {
	// AppName is the name of the KeptnApp containing the KeptnWorkload.
	AppName string `json:"app"`
	// Version defines the version of the KeptnWorkload.
	Version string `json:"version"`
	// PreDeploymentTasks is a list of all tasks to be performed during the pre-deployment phase of the KeptnWorkload.
	// The items of this list refer to the names of KeptnTaskDefinitions
	// located in the same namespace as the KeptnApp, or in the KLT namespace.
	PreDeploymentTasks []string `json:"preDeploymentTasks,omitempty"`
	// PostDeploymentTasks is a list of all tasks to be performed during the post-deployment phase of the KeptnWorkload.
	// The items of this list refer to the names of KeptnTaskDefinitions
	// located in the same namespace as the KeptnWorkload, or in the KLT namespace.
	PostDeploymentTasks []string `json:"postDeploymentTasks,omitempty"`
	// PreDeploymentEvaluations is a list of all evaluations to be performed
	// during the pre-deployment phase of the KeptnWorkload.
	// The items of this list refer to the names of KeptnEvaluationDefinitions
	// located in the same namespace as the KeptnWorkload, or in the KLT namespace.
	PreDeploymentEvaluations []string `json:"preDeploymentEvaluations,omitempty"`
	// PostDeploymentEvaluations is a list of all evaluations to be performed
	// during the post-deployment phase of the KeptnWorkload.
	// The items of this list refer to the names of KeptnEvaluationDefinitions
	// located in the same namespace as the KeptnWorkload, or in the KLT namespace.
	PostDeploymentEvaluations []string `json:"postDeploymentEvaluations,omitempty"`
	// ResourceReference is a reference to the Kubernetes resource
	// (Deployment, DaemonSet, StatefulSet or ReplicaSet) the KeptnWorkload is representing.
	ResourceReference ResourceReference `json:"resourceReference"`
}

// KeptnWorkloadStatus defines the observed state of KeptnWorkload
type KeptnWorkloadStatus struct {
	// CurrentVersion indicates the version that is currently deployed or being reconciled.
	CurrentVersion string `json:"currentVersion,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:printcolumn:name="AppName",type=string,JSONPath=`.spec.app`
// +kubebuilder:printcolumn:name="Version",type=string,JSONPath=`.spec.version`

// KeptnWorkload is the Schema for the keptnworkloads API
type KeptnWorkload struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec describes the desired state of the KeptnWorkload.
	Spec KeptnWorkloadSpec `json:"spec,omitempty"`
	// Status describes the current state of the KeptnWorkload.
	Status KeptnWorkloadStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KeptnWorkloadList contains a list of KeptnWorkload
type KeptnWorkloadList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KeptnWorkload `json:"items"`
}

type ResourceReference struct {
	UID  types.UID `json:"uid"`
	Kind string    `json:"kind"`
	Name string    `json:"name"`
}

func init() {
	SchemeBuilder.Register(&KeptnWorkload{}, &KeptnWorkloadList{})
}

func (w KeptnWorkload) GetWorkloadInstanceName() string {
	return operatorcommon.CreateResourceName(common.MaxK8sObjectLength, common.MinKLTNameLen, w.Name, w.Spec.Version)
}

func (w KeptnWorkload) SetSpanAttributes(span trace.Span) {
	span.SetAttributes(w.GetSpanAttributes()...)
}

func (w KeptnWorkload) GenerateWorkloadInstance(previousVersion string, traceContextCarrier map[string]string) KeptnWorkloadInstance {
	return KeptnWorkloadInstance{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: traceContextCarrier,
			Name:        w.GetWorkloadInstanceName(),
			Namespace:   w.Namespace,
		},
		Spec: KeptnWorkloadInstanceSpec{
			KeptnWorkloadSpec: w.Spec,
			WorkloadName:      w.Name,
			PreviousVersion:   previousVersion,
		},
	}
}

func (w KeptnWorkload) GetSpanAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		common.AppName.String(w.Spec.AppName),
		common.WorkloadName.String(w.Name),
		common.WorkloadVersion.String(w.Spec.Version),
	}
}

func (w KeptnWorkload) GetEventAnnotations() map[string]string {
	return map[string]string{
		"appName":         w.Spec.AppName,
		"workloadName":    w.Name,
		"workloadVersion": w.Spec.Version,
	}
}

func (w KeptnWorkload) GetNameWithoutAppPrefix() string {
	return strings.TrimPrefix(w.Name, fmt.Sprintf("%s-", w.Spec.AppName))
}
