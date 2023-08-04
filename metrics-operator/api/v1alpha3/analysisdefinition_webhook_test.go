package v1alpha3

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/resource"
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
			wantErr: fmt.Errorf("Target: no operator set"),
		},
		{
			name: "multiple targets set",
			target: Target{
				LessThanOrEqual: &TargetValue{
					FixedValue: *resource.NewQuantity(5, resource.DecimalSI),
				},
				LessThan: &TargetValue{
					FixedValue: *resource.NewQuantity(5, resource.DecimalSI),
				},
			},
			wantErr: fmt.Errorf("Target: multiple operators can not be set within the same target"),
		},
		{
			name: "happy path",
			target: Target{
				LessThanOrEqual: &TargetValue{
					FixedValue: *resource.NewQuantity(5, resource.DecimalSI),
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

func TestScore_validate(t *testing.T) {
	tests := []struct {
		name    string
		score   Score
		wantErr error
	}{
		{
			name: "happy path",
			score: Score{
				PassPercentage:    90,
				WarningPercentage: 80,
			},
			wantErr: nil,
		},
		{
			name: "warn and pass equal",
			score: Score{
				PassPercentage:    90,
				WarningPercentage: 90,
			},
			wantErr: fmt.Errorf("Warn percentage score cannot be higher or equal than Pass percentage score"),
		},
		{
			name: "warn higher than pass",
			score: Score{
				PassPercentage:    90,
				WarningPercentage: 95,
			},
			wantErr: fmt.Errorf("Warn percentage score cannot be higher or equal than Pass percentage score"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.wantErr, tt.score.validate())
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
				SLOTargets: SLOTarget{},
			},
			wantErr: nil,
		},
		{
			name: "only pass set",
			obj: Objective{
				SLOTargets: SLOTarget{
					Pass: &CriteriaSet{
						AnyOf: []Criteria{
							{
								AnyOf: []Target{
									{
										EqualTo: &TargetValue{
											FixedValue: *resource.NewQuantity(5, resource.DecimalSI),
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
				SLOTargets: SLOTarget{
					Warning: &CriteriaSet{
						AnyOf: []Criteria{
							{
								AnyOf: []Target{
									{
										EqualTo: &TargetValue{
											FixedValue: *resource.NewQuantity(5, resource.DecimalSI),
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
				SLOTargets: SLOTarget{
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
											FixedValue: *resource.NewQuantity(5, resource.DecimalSI),
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: fmt.Errorf("Criteria: neither AllOf nor AnyOf set"),
		},
		{
			name: "warning and pass set properly",
			obj: Objective{
				SLOTargets: SLOTarget{
					Warning: &CriteriaSet{
						AnyOf: []Criteria{
							{
								AnyOf: []Target{
									{
										EqualTo: &TargetValue{
											FixedValue: *resource.NewQuantity(5, resource.DecimalSI),
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
											FixedValue: *resource.NewQuantity(5, resource.DecimalSI),
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
			name: "failure path - objective",
			obj: AnalysisDefinition{
				Spec: AnalysisDefinitionSpec{
					Objectives: []Objective{
						{
							SLOTargets: SLOTarget{
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
			wantErr: fmt.Errorf("Criteria: neither AllOf nor AnyOf set"),
		},
		{
			name: "failure path - score",
			obj: AnalysisDefinition{
				Spec: AnalysisDefinitionSpec{
					Objectives: []Objective{
						{
							SLOTargets: SLOTarget{
								Pass: &CriteriaSet{
									AnyOf: []Criteria{
										{
											AnyOf: []Target{
												{
													EqualTo: &TargetValue{
														FixedValue: *resource.NewQuantity(5, resource.DecimalSI),
													},
												},
											},
										},
									},
								},
							},
						},
					},
					TotalScore: Score{
						PassPercentage:    80,
						WarningPercentage: 90,
					},
				},
			},
			wantErr: fmt.Errorf("Warn percentage score cannot be higher or equal than Pass percentage score"),
		},
		{
			name: "happy path",
			obj: AnalysisDefinition{
				Spec: AnalysisDefinitionSpec{
					Objectives: []Objective{
						{
							SLOTargets: SLOTarget{
								Pass: &CriteriaSet{
									AnyOf: []Criteria{
										{
											AnyOf: []Target{
												{
													EqualTo: &TargetValue{
														FixedValue: *resource.NewQuantity(5, resource.DecimalSI),
													},
												},
											},
										},
									},
								},
							},
						},
					},
					TotalScore: Score{
						PassPercentage:    80,
						WarningPercentage: 70,
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
			name:     "neither AllOf nor AnyOf set",
			criteria: Criteria{},
			wantErr:  fmt.Errorf("Criteria: neither AllOf nor AnyOf set"),
		},
		{
			name: "AllOf and AnyOf set",
			criteria: Criteria{
				AnyOf: []Target{
					{
						EqualTo: &TargetValue{
							FixedValue: *resource.NewQuantity(5, resource.DecimalSI),
						},
					},
				},
				AllOf: []Target{
					{
						EqualTo: &TargetValue{
							FixedValue: *resource.NewQuantity(5, resource.DecimalSI),
						},
					},
				},
			},
			wantErr: fmt.Errorf("Criteria: AllOf and AnyOf are set simultaneously"),
		},
		{
			name: "AllOf validation fails",
			criteria: Criteria{
				AllOf: []Target{
					{},
				},
			},
			wantErr: fmt.Errorf("Target: no operator set"),
		},
		{
			name: "AnyOf validation fails",
			criteria: Criteria{
				AnyOf: []Target{
					{},
				},
			},
			wantErr: fmt.Errorf("Target: no operator set"),
		},
		{
			name: "happy path",
			criteria: Criteria{
				AnyOf: []Target{
					{
						EqualTo: &TargetValue{
							FixedValue: *resource.NewQuantity(5, resource.DecimalSI),
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
			wantErr: fmt.Errorf("CriteriaSet: AllOf and AnyOf are set simultaneously"),
		},
		{
			name: "AllOf validation fails",
			criteria: CriteriaSet{
				AllOf: []Criteria{
					{},
				},
			},
			wantErr: fmt.Errorf("Criteria: neither AllOf nor AnyOf set"),
		},
		{
			name: "AnyOf validation fails",
			criteria: CriteriaSet{
				AnyOf: []Criteria{
					{},
				},
			},
			wantErr: fmt.Errorf("Criteria: neither AllOf nor AnyOf set"),
		},
		{
			name: "happy path",
			criteria: CriteriaSet{
				AnyOf: []Criteria{
					{
						AnyOf: []Target{
							{
								EqualTo: &TargetValue{
									FixedValue: *resource.NewQuantity(5, resource.DecimalSI),
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
