/*
Copyright 2023.

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
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AnalysisDefinitionSpec defines the desired state of AnalysisDefinition
type AnalysisDefinitionSpec struct {
	Objectives []Objective `json:"objectives,omitempty"`
	TotalScore Score       `json:"totalScore"`
}

// Score defines the required score for an evaluation to be successful
type Score struct {
	// PassPercentage defines the threshold which needs to be reached for an evaluation to pass.
	PassPercentage int `json:"passPercentage"`
	// WarningPercentage defines the threshold which needs to be reached for an evaluation to pass with a 'warning' status.
	WarningPercentage int `json:"warningPercentage"`
}

// Objective defines a list of objectives
type Objective struct {
	// AnalysisValueTemplateRef defines a reference to the used AnalysisValueTemplate
	AnalysisValueTemplateRef ObjectReference `json:"analysisValueTemplateRef"`
	// SLOTargets defines a list of SLOTargests
	SLOTargets *SLOTarget `json:"sloTargets,omitempty"`
	// Weigeht defines the importance of one SLI over the others
	// +kubebuilder:default:=1
	Weight int `json:"weight,omitempty"`
	// KeyObjective defines the meaning that the analysis fails if the objective is not met
	// +kubebuilder:default:=false
	KeyObjective bool `json:"keyObjective,omitempty"`
}

// SLOTarget defines the Criteria
type SLOTarget struct {
	// Pass defines limit up to which an evaluation is successful
	Pass *CriteriaSet `json:"pass,omitempty"`
	// Warning defines the border where the result is not pass and not fail
	Warning *CriteriaSet `json:"warning,omitempty"`
}

type TargetValue struct {
	FixedValue resource.Quantity `json:"fixedValue"`
}

type CriteriaSet struct {
	// AnyOf contains a list of targets where any of them needs to be successful for the Criteria to pass
	AnyOf []Criteria `json:"anyOf,omitempty"`
	// AllOf contains a list of targets where all of them need to be successful for the Criteria to pass
	AllOf []Criteria `json:"allOf,omitempty"`
}

type Criteria struct {
	// AnyOf contains a list of criteria where any of them needs to be successful for the CriteriaSet to pass
	AnyOf []Target `json:"anyOf,omitempty"`
	// AllOf contains a list of criteria where all of them need to be successful for the CriteriaSet to pass
	AllOf []Target `json:"allOf,omitempty"`
}

// Target specifies the supported operators for value comparisons
type Target struct {
	LessThanOrEqual    *TargetValue `json:"lessThanOrEqual,omitempty"`
	LessThan           *TargetValue `json:"lessThan,omitempty"`
	GreaterThan        *TargetValue `json:"greaterThan,omitempty"`
	GreaterThanOrEqual *TargetValue `json:"greaterThanOrEqual,omitempty"`
	EqualTo            *TargetValue `json:"equalTo,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// AnalysisDefinition is the Schema for the analysisdefinitions API
type AnalysisDefinition struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec AnalysisDefinitionSpec `json:"spec,omitempty"`
	// unused field
	Status string `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AnalysisDefinitionList contains a list of AnalysisDefinition
type AnalysisDefinitionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AnalysisDefinition `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AnalysisDefinition{}, &AnalysisDefinitionList{})
}
