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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// KeptnAppCreationRequestSpec defines the desired state of KeptnAppCreationRequest
type KeptnAppCreationRequestSpec struct {
	// AppName is the name of the KeptnApp the KeptnAppCreationRequest should create if no user-defined object with that name is found.
	AppName string `json:"appName"`
}

// KeptnAppCreationRequestStatus defines the observed state of KeptnAppCreationRequest
type KeptnAppCreationRequestStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// KeptnAppCreationRequest is the Schema for the keptnappcreationrequests API
type KeptnAppCreationRequest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KeptnAppCreationRequestSpec   `json:"spec,omitempty"`
	Status KeptnAppCreationRequestStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// KeptnAppCreationRequestList contains a list of KeptnAppCreationRequest
type KeptnAppCreationRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KeptnAppCreationRequest `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KeptnAppCreationRequest{}, &KeptnAppCreationRequestList{})
}
