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
	cfg "sigs.k8s.io/controller-runtime/pkg/config/v1alpha1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// OptionsSpec defines the desired state of Options
type OptionsSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of OptionsNew. Edit optionsnew_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// OptionsStatus defines the observed state of Options
type OptionsStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Options is the Schema for the options API
type Options struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OptionsSpec   `json:"spec,omitempty"`
	Status OptionsStatus `json:"status,omitempty"`

	// ControllerManagerConfigurationSpec returns the configurations for controllers
	cfg.ControllerManagerConfigurationSpec `json:",inline"`

	// OTelCollectorUrl can be used to send Open Telemetry metrics to an external collector
	OTelCollectorUrl string `json:"otelCollectorUrl,omitempty"`

	// DisableWebhook determines whether the pod mutating webhook should be set up to enable all features powered
	// by KLT
	DisableWebhook bool `json:"disableWebhook,omitEmpty"`

	// FunctionsRunnerImage can be used to customize the runner image and version that is used to run
	// Pre- and Post-Deployment Tasks
	FunctionsRunnerImage string `json:"functionsRunnerImage,omitempty"`
}

// +kubebuilder:object:root=true

// OptionsList contains a list of Options
type OptionsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Options `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Options{}, &OptionsList{})
}
