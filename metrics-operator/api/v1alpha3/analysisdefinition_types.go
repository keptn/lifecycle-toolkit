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
	// Objectives defines a list of objectives to evaluate for an analysis
	Objectives []Objective `json:"objectives,omitempty" yaml:"objectives,omitempty"`
	// TotalScore defines the required score for an analysis to be successful
	TotalScore TotalScore `json:"totalScore" yaml:"totalScore"`
}

// TotalScore defines the required score for an analysis to be successful
type TotalScore struct {
	// PassPercentage defines the threshold to reach for an analysis to pass
	// +kubebuilder:validation:Minimum:=0
	// +kubebuilder:validation:Maximum:=100
	PassPercentage int `json:"passPercentage" yaml:"passPercentage"`
	// WarningPercentage defines the threshold to reach for an analysis to pass with a 'warning' status
	// +kubebuilder:validation:Minimum:=0
	// +kubebuilder:validation:Maximum:=100
	WarningPercentage int `json:"warningPercentage" yaml:"warningPercentage"`
}

// Objective defines an objective for analysis
type Objective struct {
	// AnalysisValueTemplateRef refers to the appropriate AnalysisValueTemplate
	AnalysisValueTemplateRef ObjectReference `json:"analysisValueTemplateRef" yaml:"analysisValueTemplateRef"`
	// Target defines failure or warning criteria
	Target Target `json:"target,omitempty" yaml:"target,omitempty"`
	// Weight can be used to emphasize the importance of one Objective over the others
	// +kubebuilder:default:=1
	Weight int `json:"weight,omitempty" yaml:"weight,omitempty"`
	// KeyObjective defines whether the whole analysis fails when this objective's target is not met
	// +kubebuilder:default:=false
	KeyObjective bool `json:"keyObjective,omitempty" yaml:"keyObjective,omitempty"`
}

// Target defines the failure and warning criteria
type Target struct {
	// Failure defines limits up to which an analysis fails
	Failure *Operator `json:"failure,omitempty" yaml:"failure,omitempty"`
	// Warning defines limits where the result does not pass or fail
	Warning *Operator `json:"warning,omitempty" yaml:"warning,omitempty"`
}

// OperatorValue represents the value to which the result is compared
type OperatorValue struct {
	// FixedValue defines the value for comparison
	FixedValue resource.Quantity `json:"fixedValue" yaml:"fixedValue"`
}

// RangeValue represents a range which the value should fit
type RangeValue struct {
	// LowBound defines the lower bound of the range
	LowBound resource.Quantity `json:"lowBound"`
	// HighBound defines the higher bound of the range
	HighBound resource.Quantity `json:"highBound"`
}

// Operator specifies the supported operators for value comparisons
type Operator struct {
	// LessThanOrEqual represents '<=' operator
	LessThanOrEqual *OperatorValue `json:"lessThanOrEqual,omitempty" yaml:"lessThanOrEqual,omitempty"`
	// LessThan represents '<' operator
	LessThan *OperatorValue `json:"lessThan,omitempty" yaml:"lessThan,omitempty"`
	// GreaterThan represents '>' operator
	GreaterThan *OperatorValue `json:"greaterThan,omitempty" yaml:"greaterThan,omitempty"`
	// GreaterThanOrEqual represents '>=' operator
	GreaterThanOrEqual *OperatorValue `json:"greaterThanOrEqual,omitempty" yaml:"greaterThanOrEqual,omitempty"`
	// EqualTo represents '==' operator
	EqualTo *OperatorValue `json:"equalTo,omitempty" yaml:"equalTo,omitempty"`
	// InRange represents operator checking the value is inclusively in the defined range
	InRange *RangeValue `json:"inRange,omitempty" yaml:"inRange,omitempty"`
	// NotInRange represents operator checking the value is exclusively out of the defined range
	NotInRange *RangeValue `json:"notInRange,omitempty" yaml:"notInRange,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// AnalysisDefinition is the Schema for the analysisdefinitions APIs
type AnalysisDefinition struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`

	Spec AnalysisDefinitionSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
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

func (o *OperatorValue) GetFloatValue() float64 {
	return o.FixedValue.AsApproximateFloat64()
}
