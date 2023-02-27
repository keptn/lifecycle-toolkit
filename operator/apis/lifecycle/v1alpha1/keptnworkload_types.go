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
	"strings"

	"github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha1/common"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// KeptnWorkloadSpec defines the desired state of KeptnWorkload
type KeptnWorkloadSpec struct {
	AppName                   string            `json:"app"`
	Version                   string            `json:"version"`
	PreDeploymentTasks        []string          `json:"preDeploymentTasks,omitempty"`
	PostDeploymentTasks       []string          `json:"postDeploymentTasks,omitempty"`
	PreDeploymentEvaluations  []string          `json:"preDeploymentEvaluations,omitempty"`
	PostDeploymentEvaluations []string          `json:"postDeploymentEvaluations,omitempty"`
	ResourceReference         ResourceReference `json:"resourceReference"`
}

// KeptnWorkloadStatus defines the observed state of KeptnWorkload
type KeptnWorkloadStatus struct {
	CurrentVersion string `json:"currentVersion,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="AppName",type=string,JSONPath=`.spec.app`
// +kubebuilder:printcolumn:name="Version",type=string,JSONPath=`.spec.version`

// KeptnWorkload is the Schema for the keptnworkloads API
type KeptnWorkload struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KeptnWorkloadSpec   `json:"spec,omitempty"`
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
	return strings.ToLower(w.Name + "-" + w.Spec.Version)
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

func (i KeptnWorkload) GetSpanAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		common.AppName.String(i.Spec.AppName),
		common.WorkloadName.String(i.Name),
		common.WorkloadVersion.String(i.Spec.Version),
	}
}
