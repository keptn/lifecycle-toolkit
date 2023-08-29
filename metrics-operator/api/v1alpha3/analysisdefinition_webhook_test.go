package v1alpha3

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestOperator_validate(t *testing.T) {
	tests := []struct {
		name     string
		operator Operator
		wantErr  error
	}{
		{
			name:     "no operator set",
			operator: Operator{},
			wantErr:  fmt.Errorf("Operator: no operator set"),
		},
		{
			name: "multiple operators set",
			operator: Operator{
				LessThanOrEqual: &OperatorValue{
					FixedValue: *resource.NewQuantity(5, resource.DecimalSI),
				},
				LessThan: &OperatorValue{
					FixedValue: *resource.NewQuantity(5, resource.DecimalSI),
				},
				InRange: &RangeValue{
					LowBound:  *resource.NewQuantity(5, resource.DecimalSI),
					HighBound: *resource.NewQuantity(15, resource.DecimalSI),
				},
			},
			wantErr: fmt.Errorf("Operator: multiple operators can not be set"),
		},
		{
			name: "in range - fail validation",
			operator: Operator{
				InRange: &RangeValue{
					LowBound:  *resource.NewQuantity(25, resource.DecimalSI),
					HighBound: *resource.NewQuantity(15, resource.DecimalSI),
				},
			},
			wantErr: fmt.Errorf("RangeValue: lower bound of the range needs to be smaller than higher bound"),
		},
		{
			name: "not in range - fail validation",
			operator: Operator{
				NotInRange: &RangeValue{
					LowBound:  *resource.NewQuantity(25, resource.DecimalSI),
					HighBound: *resource.NewQuantity(15, resource.DecimalSI),
				},
			},
			wantErr: fmt.Errorf("RangeValue: lower bound of the range needs to be smaller than higher bound"),
		},
		{
			name: "happy path",
			operator: Operator{
				LessThanOrEqual: &OperatorValue{
					FixedValue: *resource.NewQuantity(5, resource.DecimalSI),
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.wantErr, tt.operator.validate())
		})
	}
}

func TestRangeValue_validate(t *testing.T) {
	tests := []struct {
		name    string
		r       RangeValue
		wantErr error
	}{
		{
			name: "happy path",
			r: RangeValue{
				LowBound:  *resource.NewQuantity(5, resource.DecimalSI),
				HighBound: *resource.NewQuantity(15, resource.DecimalSI),
			},
			wantErr: nil,
		},
		{
			name: "equal bounds",
			r: RangeValue{
				LowBound:  *resource.NewQuantity(5, resource.DecimalSI),
				HighBound: *resource.NewQuantity(5, resource.DecimalSI),
			},
			wantErr: fmt.Errorf("RangeValue: lower bound of the range needs to be smaller than higher bound"),
		},
		{
			name: "lower greater that higher bound",
			r: RangeValue{
				LowBound:  *resource.NewQuantity(15, resource.DecimalSI),
				HighBound: *resource.NewQuantity(5, resource.DecimalSI),
			},
			wantErr: fmt.Errorf("RangeValue: lower bound of the range needs to be smaller than higher bound"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.wantErr, tt.r.validate())
		})
	}
}

func TestScore_validate(t *testing.T) {
	tests := []struct {
		name    string
		score   TotalScore
		wantErr error
	}{
		{
			name: "happy path",
			score: TotalScore{
				PassPercentage:    90,
				WarningPercentage: 80,
			},
			wantErr: nil,
		},
		{
			name: "warn and pass equal",
			score: TotalScore{
				PassPercentage:    90,
				WarningPercentage: 90,
			},
			wantErr: fmt.Errorf("Warn percentage score cannot be higher or equal than Pass percentage score"),
		},
		{
			name: "warn higher than pass",
			score: TotalScore{
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
			name:    "no Target set",
			obj:     Objective{},
			wantErr: nil,
		},
		{
			name: "only warning set",
			obj: Objective{
				Target: Target{
					Warning: &Operator{
						EqualTo: &OperatorValue{
							FixedValue: *resource.NewQuantity(5, resource.DecimalSI),
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "only failure set",
			obj: Objective{
				Target: Target{
					Failure: &Operator{
						EqualTo: &OperatorValue{
							FixedValue: *resource.NewQuantity(5, resource.DecimalSI),
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "warning and failure set properly",
			obj: Objective{
				Target: Target{
					Warning: &Operator{
						EqualTo: &OperatorValue{
							FixedValue: *resource.NewQuantity(5, resource.DecimalSI),
						},
					},
					Failure: &Operator{
						EqualTo: &OperatorValue{
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
							Target: Target{
								Failure: &Operator{},
							},
						},
					},
				},
			},
			wantErr: fmt.Errorf("Operator: no operator set"),
		},
		{
			name: "failure path - score",
			obj: AnalysisDefinition{
				Spec: AnalysisDefinitionSpec{
					Objectives: []Objective{
						{
							Target: Target{
								Failure: &Operator{
									EqualTo: &OperatorValue{
										FixedValue: *resource.NewQuantity(5, resource.DecimalSI),
									},
								},
							},
						},
					},
					TotalScore: TotalScore{
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
							Target: Target{
								Failure: &Operator{
									EqualTo: &OperatorValue{
										FixedValue: *resource.NewQuantity(5, resource.DecimalSI),
									},
								},
							},
						},
					},
					TotalScore: TotalScore{
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

func TestTarget_validate(t *testing.T) {
	tests := []struct {
		name    string
		target  Target
		wantErr error
	}{
		{
			name:    "neither Failure and Warning set",
			target:  Target{},
			wantErr: nil,
		},
		{
			name: "Failure set",
			target: Target{
				Failure: &Operator{
					LessThanOrEqual: &OperatorValue{
						FixedValue: *resource.NewQuantity(5, resource.DecimalSI),
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "only warning set",
			target: Target{
				Warning: &Operator{
					LessThanOrEqual: &OperatorValue{
						FixedValue: *resource.NewQuantity(5, resource.DecimalSI),
					},
				},
			},
			wantErr: nil,
		},
		{
			name:    "neither failure nor warning set",
			target:  Target{},
			wantErr: nil,
		},
		{
			name: "only failure set",
			target: Target{
				Failure: &Operator{
					EqualTo: &OperatorValue{
						FixedValue: *resource.NewQuantity(5, resource.DecimalSI),
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "only warning set",
			target: Target{
				Warning: &Operator{
					EqualTo: &OperatorValue{
						FixedValue: *resource.NewQuantity(5, resource.DecimalSI),
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "warning not set properly",
			target: Target{
				Warning: &Operator{},
				Failure: &Operator{
					EqualTo: &OperatorValue{
						FixedValue: *resource.NewQuantity(5, resource.DecimalSI),
					},
				},
			},
			wantErr: fmt.Errorf("Operator: no operator set"),
		},
		{
			name: "warning and failure set properly",
			target: Target{
				Warning: &Operator{
					EqualTo: &OperatorValue{
						FixedValue: *resource.NewQuantity(5, resource.DecimalSI),
					},
				},
				Failure: &Operator{
					EqualTo: &OperatorValue{
						FixedValue: *resource.NewQuantity(5, resource.DecimalSI),
					},
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
