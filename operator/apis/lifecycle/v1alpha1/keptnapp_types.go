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
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// KeptnAppSpec defines the desired state of KeptnApp
type KeptnAppSpec struct {
	Version                   string             `json:"version"`
	Workloads                 []KeptnWorkloadRef `json:"workloads,omitempty"`
	PreDeploymentTasks        []string           `json:"preDeploymentTasks,omitempty"`
	PostDeploymentTasks       []string           `json:"postDeploymentTasks,omitempty"`
	PreDeploymentEvaluations  []string           `json:"preDeploymentEvaluations,omitempty"`
	PostDeploymentEvaluations []string           `json:"postDeploymentEvaluations,omitempty"`
}

// KeptnAppStatus defines the observed state of KeptnApp
type KeptnAppStatus struct {
	CurrentVersion string `json:"currentVersion,omitempty"`
}

type KeptnWorkloadRef struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// KeptnApp is the Schema for the keptnapps API
type KeptnApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KeptnAppSpec   `json:"spec,omitempty"`
	Status KeptnAppStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

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
	return strings.ToLower(a.Name + "-" + a.Spec.Version)
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
