package analysis

import (
	"testing"

	"github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/analysis/fake"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/analysis/types"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestAnalysisEvaluator_Evaluate(t *testing.T) {
	tests := []struct {
		name            string
		values          map[string]string
		a               v1alpha3.AnalysisDefinition
		want            types.AnalysisResult
		mockedEvaluator IObjectiveEvaluator
	}{
		{
			name:   "no objectives",
			values: map[string]string{},
			a: v1alpha3.AnalysisDefinition{
				Spec: v1alpha3.AnalysisDefinitionSpec{
					Objectives: []v1alpha3.Objective{},
				},
			},
			want: types.AnalysisResult{
				TotalScore:       0.0,
				MaximumScore:     0.0,
				Pass:             true,
				Warning:          false,
				ObjectiveResults: []types.ObjectiveResult{},
			},
			mockedEvaluator: nil,
		},
		{
			name:   "pass scenario",
			values: map[string]string{},
			a: v1alpha3.AnalysisDefinition{
				Spec: v1alpha3.AnalysisDefinitionSpec{
					Objectives: []v1alpha3.Objective{
						{
							Weight: 10,
						},
					},
					TotalScore: v1alpha3.TotalScore{
						PassPercentage:    80,
						WarningPercentage: 50,
					},
				},
			},
			want: types.AnalysisResult{
				TotalScore:   10.0,
				MaximumScore: 10.0,
				Pass:         true,
				Warning:      false,
				ObjectiveResults: []types.ObjectiveResult{
					{
						Result: types.TargetResult{},
						Value:  5.0,
						Score:  10.0,
						Error:  nil,
					},
				},
			},
			mockedEvaluator: &fake.IObjectiveEvaluatorMock{
				EvaluateFunc: func(values map[string]string, objective v1alpha3.Objective) types.ObjectiveResult {
					return types.ObjectiveResult{
						Result: types.TargetResult{},
						Value:  5.0,
						Score:  10.0,
						Error:  nil,
					}
				},
			},
		},
		{
			name:   "pass scenario - multiple objectives",
			values: map[string]string{},
			a: v1alpha3.AnalysisDefinition{
				Spec: v1alpha3.AnalysisDefinitionSpec{
					Objectives: []v1alpha3.Objective{
						{
							Weight: 10,
						},
						{
							Weight: 10,
						},
					},
					TotalScore: v1alpha3.TotalScore{
						PassPercentage:    80,
						WarningPercentage: 50,
					},
				},
			},
			want: types.AnalysisResult{
				TotalScore:   20.0,
				MaximumScore: 20.0,
				Pass:         true,
				Warning:      false,
				ObjectiveResults: []types.ObjectiveResult{
					{
						Result: types.TargetResult{},
						Value:  5.0,
						Score:  10.0,
						Error:  nil,
					},
					{
						Result: types.TargetResult{},
						Value:  5.0,
						Score:  10.0,
						Error:  nil,
					},
				},
			},
			mockedEvaluator: &fake.IObjectiveEvaluatorMock{
				EvaluateFunc: func(values map[string]string, objective v1alpha3.Objective) types.ObjectiveResult {
					return types.ObjectiveResult{
						Result: types.TargetResult{},
						Value:  5.0,
						Score:  10.0,
						Error:  nil,
					}
				},
			},
		},
		{
			name:   "warning scenario",
			values: map[string]string{},
			a: v1alpha3.AnalysisDefinition{
				Spec: v1alpha3.AnalysisDefinitionSpec{
					Objectives: []v1alpha3.Objective{
						{
							Weight: 10,
						},
					},
					TotalScore: v1alpha3.TotalScore{
						PassPercentage:    80,
						WarningPercentage: 50,
					},
				},
			},
			want: types.AnalysisResult{
				TotalScore:   5.0,
				MaximumScore: 10.0,
				Pass:         false,
				Warning:      true,
				ObjectiveResults: []types.ObjectiveResult{
					{
						Result: types.TargetResult{},
						Value:  5.0,
						Score:  5.0,
						Error:  nil,
					},
				},
			},
			mockedEvaluator: &fake.IObjectiveEvaluatorMock{
				EvaluateFunc: func(values map[string]string, objective v1alpha3.Objective) types.ObjectiveResult {
					return types.ObjectiveResult{
						Result: types.TargetResult{},
						Value:  5.0,
						Score:  5.0,
						Error:  nil,
					}
				},
			},
		},
		{
			name:   "fail scenario",
			values: map[string]string{},
			a: v1alpha3.AnalysisDefinition{
				Spec: v1alpha3.AnalysisDefinitionSpec{
					Objectives: []v1alpha3.Objective{
						{
							Weight: 10,
						},
					},
					TotalScore: v1alpha3.TotalScore{
						PassPercentage:    80,
						WarningPercentage: 50,
					},
				},
			},
			want: types.AnalysisResult{
				TotalScore:   0.0,
				MaximumScore: 10.0,
				Pass:         false,
				Warning:      false,
				ObjectiveResults: []types.ObjectiveResult{
					{
						Result: types.TargetResult{},
						Value:  5.0,
						Score:  0.0,
						Error:  nil,
					},
				},
			},
			mockedEvaluator: &fake.IObjectiveEvaluatorMock{
				EvaluateFunc: func(values map[string]string, objective v1alpha3.Objective) types.ObjectiveResult {
					return types.ObjectiveResult{
						Result: types.TargetResult{},
						Value:  5.0,
						Score:  0.0,
						Error:  nil,
					}
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ae := NewAnalysisEvaluator(tt.mockedEvaluator)
			require.Equal(t, tt.want, ae.Evaluate(tt.values, tt.a))
		})
	}
}

// The scenario where the score is passing, but keyObjective fails
// we cannot mock, therefore we need to use the full slo scoring
// scenario to test this
func TestAnalysisEvaluator_Evaluate_failKeyObjective(t *testing.T) {
	te := NewTargetEvaluator(&OperatorEvaluator{})
	oe := NewObjectiveEvaluator(&te)
	ae := NewAnalysisEvaluator(&oe)

	values := map[string]string{
		"key1": "5",
		"key2": "0",
	}

	ad := v1alpha3.AnalysisDefinition{
		Spec: v1alpha3.AnalysisDefinitionSpec{
			Objectives: []v1alpha3.Objective{
				{
					AnalysisValueTemplateRef: v1alpha3.ObjectReference{
						Name: "key1",
					},
					Target: v1alpha3.Target{
						Failure: &v1alpha3.Operator{
							LessThan: &v1alpha3.OperatorValue{
								FixedValue: *resource.NewQuantity(3, resource.DecimalSI),
							},
						},
					},
					Weight:       10,
					KeyObjective: false,
				},
				{
					AnalysisValueTemplateRef: v1alpha3.ObjectReference{
						Name: "key2",
					},
					Target: v1alpha3.Target{
						Failure: &v1alpha3.Operator{
							LessThan: &v1alpha3.OperatorValue{
								FixedValue: *resource.NewQuantity(3, resource.DecimalSI),
							},
						},
					},
					Weight:       1,
					KeyObjective: true,
				},
			},
			TotalScore: v1alpha3.TotalScore{
				PassPercentage:    80,
				WarningPercentage: 50,
			},
		},
	}

	expectedRes := types.AnalysisResult{
		TotalScore:   10.0,
		MaximumScore: 11.0,
		Pass:         false,
		Warning:      false,
		ObjectiveResults: []types.ObjectiveResult{
			{
				Result: types.TargetResult{
					FailureResult: types.OperatorResult{
						Operator: v1alpha3.Operator{
							LessThan: &v1alpha3.OperatorValue{
								FixedValue: *resource.NewQuantity(3, resource.DecimalSI),
							},
						},
						Fulfilled: false,
					},
					Warning: false,
					Pass:    true,
				},
				Value: 5.0,
				Score: 10.0,
				Error: nil,
			},
			{
				Result: types.TargetResult{
					FailureResult: types.OperatorResult{
						Operator: v1alpha3.Operator{
							LessThan: &v1alpha3.OperatorValue{
								FixedValue: *resource.NewQuantity(3, resource.DecimalSI),
							},
						},
						Fulfilled: true,
					},
					Warning: false,
					Pass:    false,
				},
				Value: 0.0,
				Score: 0.0,
				Error: nil,
			},
		},
	}
	require.Equal(t, expectedRes, ae.Evaluate(values, ad))
}
