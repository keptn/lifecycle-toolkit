package v1alpha3

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTarget_validate(t *testing.T) {
	tests := []struct {
		name    string
		target  Target
		wantErr error
	}{
		{
			name:    "no target set",
			target:  Target{},
			wantErr: fmt.Errorf("Target: not set"),
		},
		{
			name: "multiple targets set",
			target: Target{
				LessThanOrEqual: &TargetValue{
					FixedValue: 5,
				},
				LessThan: &TargetValue{
					FixedValue: 5,
				},
			},
			wantErr: fmt.Errorf("Target: multiple targets set anot allowed per Analysis"),
		},
		{
			name: "happy path",
			target: Target{
				LessThanOrEqual: &TargetValue{
					FixedValue: 5,
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.wantErr, tt.target.validate())
		})
	}
}

func TestObjective_validate(t *testing.T) {
	tests := []struct {
		name    string
		obj     Objective
		wantErr error
	}{
		{
			name:    "no SLOTarget set",
			obj:     Objective{},
			wantErr: nil,
		},
		{
			name: "neither pass nor warning set",
			obj: Objective{
				SLOTargets: &SLOTarget{},
			},
			wantErr: nil,
		},
		{
			name: "only pass set",
			obj: Objective{
				SLOTargets: &SLOTarget{
					Pass: &CriteriaSet{
						AnyOf: []Criteria{
							{
								AnyOf: []Target{
									{
										EqualTo: &TargetValue{
											FixedValue: 5,
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "only warning set",
			obj: Objective{
				SLOTargets: &SLOTarget{
					Warning: &CriteriaSet{
						AnyOf: []Criteria{
							{
								AnyOf: []Target{
									{
										EqualTo: &TargetValue{
											FixedValue: 5,
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: fmt.Errorf("Warning criteria cannot be set without Pass criteria"),
		},
		{
			name: "warning not set properly",
			obj: Objective{
				SLOTargets: &SLOTarget{
					Warning: &CriteriaSet{
						AnyOf: []Criteria{
							{},
						},
					},
					Pass: &CriteriaSet{
						AnyOf: []Criteria{
							{
								AnyOf: []Target{
									{
										EqualTo: &TargetValue{
											FixedValue: 5,
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: fmt.Errorf("Criteria: AllOf nor Anyof set"),
		},
		{
			name: "warning and pass set properly",
			obj: Objective{
				SLOTargets: &SLOTarget{
					Warning: &CriteriaSet{
						AnyOf: []Criteria{
							{
								AnyOf: []Target{
									{
										EqualTo: &TargetValue{
											FixedValue: 5,
										},
									},
								},
							},
						},
					},
					Pass: &CriteriaSet{
						AnyOf: []Criteria{
							{
								AnyOf: []Target{
									{
										EqualTo: &TargetValue{
											FixedValue: 5,
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.wantErr, tt.obj.validate())
		})
	}
}

func TestAnalysisDefinition_validateCreateUpdate(t *testing.T) {
	tests := []struct {
		name    string
		obj     AnalysisDefinition
		wantErr error
	}{
		{
			name: "no spec set",
			obj: AnalysisDefinition{
				Spec: AnalysisDefinitionSpec{},
			},
			wantErr: nil,
		},
		{
			name: "failure path",
			obj: AnalysisDefinition{
				Spec: AnalysisDefinitionSpec{
					Objectives: []Objective{
						{
							SLOTargets: &SLOTarget{
								Pass: &CriteriaSet{
									AnyOf: []Criteria{
										{},
									},
								},
							},
						},
					},
				},
			},
			wantErr: fmt.Errorf("Criteria: AllOf nor Anyof set"),
		},
		{
			name: "happy path",
			obj: AnalysisDefinition{
				Spec: AnalysisDefinitionSpec{
					Objectives: []Objective{
						{
							SLOTargets: &SLOTarget{
								Pass: &CriteriaSet{
									AnyOf: []Criteria{
										{
											AnyOf: []Target{
												{
													EqualTo: &TargetValue{
														FixedValue: 5,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.wantErr, tt.obj.ValidateCreate())
			require.Equal(t, tt.wantErr, tt.obj.ValidateUpdate(&AnalysisDefinition{}))
		})
	}
}

func TestCriteria_validate(t *testing.T) {
	tests := []struct {
		name     string
		criteria Criteria
		wantErr  error
	}{
		{
			name:     "either AllOf nor AnyOf set",
			criteria: Criteria{},
			wantErr:  fmt.Errorf("Criteria: AllOf nor Anyof set"),
		},
		{
			name: "AllOf and AnyOf set",
			criteria: Criteria{
				AnyOf: []Target{
					{
						EqualTo: &TargetValue{
							FixedValue: 5,
						},
					},
				},
				AllOf: []Target{
					{
						EqualTo: &TargetValue{
							FixedValue: 5,
						},
					},
				},
			},
			wantErr: fmt.Errorf("Criteria: AllOf and Anyof are set simultaneusly"),
		},
		{
			name: "AllOf validation fails",
			criteria: Criteria{
				AllOf: []Target{
					{},
				},
			},
			wantErr: fmt.Errorf("Target: not set"),
		},
		{
			name: "AnyOf validation fails",
			criteria: Criteria{
				AnyOf: []Target{
					{},
				},
			},
			wantErr: fmt.Errorf("Target: not set"),
		},
		{
			name: "happy path",
			criteria: Criteria{
				AnyOf: []Target{
					{
						EqualTo: &TargetValue{
							FixedValue: 5,
						},
					},
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.wantErr, tt.criteria.validate())
		})
	}
}

func TestCriteriaSet_validate(t *testing.T) {
	tests := []struct {
		name     string
		criteria CriteriaSet
		wantErr  error
	}{
		{
			name:     "neither AllOf nor AnyOf set",
			criteria: CriteriaSet{},
			wantErr:  nil,
		},
		{
			name: "AllOf and AnyOf set",
			criteria: CriteriaSet{
				AnyOf: []Criteria{
					{},
				},
				AllOf: []Criteria{
					{},
				},
			},
			wantErr: fmt.Errorf("CriteriaSet: AllOf and Anyof are set simultaneusly"),
		},
		{
			name: "AllOf validation fails",
			criteria: CriteriaSet{
				AllOf: []Criteria{
					{},
				},
			},
			wantErr: fmt.Errorf("Criteria: AllOf nor Anyof set"),
		},
		{
			name: "AnyOf validation fails",
			criteria: CriteriaSet{
				AnyOf: []Criteria{
					{},
				},
			},
			wantErr: fmt.Errorf("Criteria: AllOf nor Anyof set"),
		},
		{
			name: "happy path",
			criteria: CriteriaSet{
				AnyOf: []Criteria{
					{
						AnyOf: []Target{
							{
								EqualTo: &TargetValue{
									FixedValue: 5,
								},
							},
						},
					},
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.wantErr, tt.criteria.validate())
		})
	}
}
