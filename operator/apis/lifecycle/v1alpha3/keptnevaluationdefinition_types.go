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

// PassThreshold defines the required score for an evaluation to be successful
type PassThreshold struct {
	// PassPercentage defines the threshold which needs to be reached for an evaluation to pass.
	PassPercentage float64 `json:"passPercentage"`
	// WarningPercentage defines the threshold which needs to be reached for an evaluation to pass with a 'warning'  status.
	WarningPercentage float64 `json:"passPercentage"`
}

// ComparisonTarget defines how the target value for a comparison based evaluation is calculated. Only one of its field can be set
type ComparisonTarget struct {
	IncreaseByPercent *float32 `json:"increaseByPercent,omitempty"`
	DecreaseByPercent *float32 `json:"decreaseByPercent,omitempty"`
	IncreaseBy        *float32 `json:"increaseBy,omitempty"`
	DecreaseBy        *float32 `json:"decreaseBy,omitempty"`
	Equal             *bool    `json:"equal,omitempty"`
}

type TargetValue struct {
	FixedValue *float64          `json:"fixedValue,omitempty"`
	Comparison *ComparisonTarget `json:"compareValue,omitempty"`
}

type Target struct {
	LessThanOrEqual    *TargetValue `json:"lessThanOrEqual,omitempty"`
	LessThan           *TargetValue `json:"lessThan,omitempty"`
	GreaterThan        *TargetValue `json:"greaterThan,omitempty"`
	GreaterThanOrEqual *TargetValue `json:"greaterThanOrEqual,omitempty"`
	EqualTo            *TargetValue `json:"equalTo,omitempty"`
}

type Criteria struct {
	// AnyOf contains a list of targets where any of them needs to be successful for the Criteria to pass
	AnyOf []Target `json:"anyOf"`
	// AllOf contains a list of targets where all of them need to be successful for the Criteria to pass
	AllOf []Target `json:"allOf"`
}

type CriteriaSet struct {
	// AnyOf contains a list of criteria where any of them needs to be successful for the CriteriaSet to pass
	AnyOf []Criteria `json:"anyOf"`
	// AllOf contains a list of criteria where all of them need to be successful for the CriteriaSet to pass
	AllOf []Criteria `json:"allOf"`
}

type TargetResult struct {
	Target
	Violated bool
}

type CriteriaResult struct {
	ViolatedTargets []TargetResult
	Violated        bool
}

type CriteriaSetResult struct {
	ViolatedCriteria []CriteriaResult
	Violated         bool
}

type SLOTarget struct {
	Pass    CriteriaSet `json:"pass"`
	Warning CriteriaSet `json:"warning"`
}

type ComparisonSpec struct {
	CompareWith         string `json:"compareWith"`
	IncludeWarning      bool   `json:"includeWarning"`
	NumberOfComparisons int    `json:"numberOfComparisons"`
	AggregationFunction string `json:"aggregationFunction"`
}

// KeptnEvaluationDefinitionSpec defines the desired state of KeptnEvaluationDefinition
type KeptnEvaluationDefinitionSpec struct {
	// Objectives is a list of objectives that have to be met for a KeptnEvaluation referencing this
	// KeptnEvaluationDefinition to be successful.
	Objectives []Objective `json:"objectives"`
	// TotalScore allows to define a minimum required score for passing an evaluation.
	// If this is not defined, all objectives have to be successful for the evaluation to pass.
	TotalScore *PassThreshold `json:"totalScore"`
	// Comparison defines which previous KeptnEvaluations should be taken into consideration for comparison based targets
	Comparison *ComparisonSpec `json:"comparison,omitempty"`
}

type Objective struct {
	// KeptnMetricRef references the KeptnMetric that should be evaluated.
	KeptnMetricRef KeptnMetricReference `json:"keptnMetricRef"`
	// EvaluationTarget specifies the target value for the references KeptnMetric.
	// Needs to start with either '<' or '>', followed by the target value (e.g. '<10').
	// Likely to be deprecated and replaced by SLOTargets
	EvaluationTarget string `json:"evaluationTarget"`
	// SLOTargets provide a more flexible way to define targets for a metric.
	SLOTargets *SLOTarget `json:"slo_targets"`
	// Weight defines how much the Objective affects the overall score. Defaults to 1.
	Weight int `json:"weight"`
	// KeyObjective defines if the Objective is mandatory for an KeptnEvaluation to pass
	KeyObjective bool `json:"keyObjective"`
}

type ObjectiveResult struct {
	PassResult    CriteriaSetResult
	WarningResult CriteriaSetResult
	Value         float64
	Score         float64
	KeyObjective  bool
	Error         error
}

type KeptnMetricReference struct {
	// Name is the name of the referenced KeptnMetric.
	Name string `json:"name"`
	// Namespace is the namespace where the referenced KeptnMetric is located.
	Namespace string `json:"namespace,omitempty"`
}

// KeptnEvaluationDefinitionStatus defines the observed state of KeptnEvaluationDefinition.
type KeptnEvaluationDefinitionStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:storageversion
//+kubebuilder:resource:path=keptnevaluationdefinitions,shortName=ked

// KeptnEvaluationDefinition is the Schema for the keptnevaluationdefinitions API
type KeptnEvaluationDefinition struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec describes the desired state of the KeptnEvaluationDefinition.
	Spec KeptnEvaluationDefinitionSpec `json:"spec,omitempty"`
	// Status describes the current state of the KeptnEvaluationDefinition.
	Status KeptnEvaluationDefinitionStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// KeptnEvaluationDefinitionList contains a list of KeptnEvaluationDefinition
type KeptnEvaluationDefinitionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KeptnEvaluationDefinition `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KeptnEvaluationDefinition{}, &KeptnEvaluationDefinitionList{})
}
