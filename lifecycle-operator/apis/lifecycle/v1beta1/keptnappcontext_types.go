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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DeploymentTaskSpec struct {
	// PreDeploymentTasks is a list of all tasks to be performed during the pre-deployment phase of the KeptnApp.
	// The items of this list refer to the names of KeptnTaskDefinitions
	// located in the same namespace as the KeptnApp, or in the Keptn namespace.
	PreDeploymentTasks []string `json:"preDeploymentTasks,omitempty"`
	// PostDeploymentTasks is a list of all tasks to be performed during the post-deployment phase of the KeptnApp.
	// The items of this list refer to the names of KeptnTaskDefinitions
	// located in the same namespace as the KeptnApp, or in the Keptn namespace.
	PostDeploymentTasks []string `json:"postDeploymentTasks,omitempty"`
	// PreDeploymentEvaluations is a list of all evaluations to be performed
	// during the pre-deployment phase of the KeptnApp.
	// The items of this list refer to the names of KeptnEvaluationDefinitions
	// located in the same namespace as the KeptnApp, or in the Keptn namespace.
	PreDeploymentEvaluations []string `json:"preDeploymentEvaluations,omitempty"`
	// PostDeploymentEvaluations is a list of all evaluations to be performed
	// during the post-deployment phase of the KeptnApp.
	// The items of this list refer to the names of KeptnEvaluationDefinitions
	// located in the same namespace as the KeptnApp, or in the Keptn namespace.
	PostDeploymentEvaluations []string `json:"postDeploymentEvaluations,omitempty"`
	// PromotionTasks is a list of all tasks to be performed during the promotion phase of the KeptnApp.
	// The items of this list refer to the names of KeptnTaskDefinitions
	// located in the same namespace as the KeptnApp, or in the Keptn namespace.
	PromotionTasks []string `json:"promotionTasks,omitempty"`
}

// KeptnAppContextSpec defines the desired state of KeptnAppContext
type KeptnAppContextSpec struct {
	DeploymentTaskSpec `json:",inline"`

	// +optional
	// Metadata contains additional key-value pairs for contextual information.
	Metadata map[string]string `json:"metadata,omitempty"`

	// +optional
	// SpanLinks are links to OpenTelemetry span IDs for tracking. These links establish relationships between spans across different services, enabling distributed tracing.
	// For more information on OpenTelemetry span links, refer to the documentation: https://opentelemetry.io/docs/concepts/signals/traces/#span-links
	SpanLinks []string `json:"spanLinks,omitempty"`
}

// KeptnAppContextStatus defines the observed state of KeptnAppContext
type KeptnAppContextStatus struct {
	// unused field
	// +optional
	Status string `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// KeptnAppContext is the Schema for the keptnappcontexts API
type KeptnAppContext struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KeptnAppContextSpec   `json:"spec,omitempty"`
	Status KeptnAppContextStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// KeptnAppContextList contains a list of KeptnAppContext
type KeptnAppContextList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KeptnAppContext `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KeptnAppContext{}, &KeptnAppContextList{})
}
