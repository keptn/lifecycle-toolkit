package analysis

import (
	"github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/common/analysis/fake"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAnalysisEvaluator_Evaluate(t *testing.T) {
	type fields struct {
		ObjectiveEvaluator *fake.IObjectiveEvaluatorMock
	}
	type args struct {
		values map[string]float64
		ed     v1alpha3.KeptnEvaluationDefinition
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *AnalysisResult
		wantErr bool
	}{
		{
			name: "all objectives pass",
			fields: fields{
				ObjectiveEvaluator: &fake.IObjectiveEvaluatorMock{
					EvaluateFunc: func(values map[string]float64, objective v1alpha3.Objective) v1alpha3.ObjectiveResult {
						return v1alpha3.ObjectiveResult{
							PassResult:    v1alpha3.CriteriaSetResult{},
							WarningResult: v1alpha3.CriteriaSetResult{},
							Value:         10.0,
							Score:         float64(objective.Weight),
							KeyObjective:  false,
							Error:         nil,
						}
					},
				},
			},
			args: args{
				values: map[string]float64{},
				ed: v1alpha3.KeptnEvaluationDefinition{
					Spec: v1alpha3.KeptnEvaluationDefinitionSpec{
						Objectives: []v1alpha3.Objective{
							{
								Weight: 1,
							},
							{
								Weight: 2,
							},
						},
						TotalScore: &v1alpha3.PassThreshold{
							PassPercentage:    90,
							WarningPercentage: 80,
						},
					},
				},
			},
			want: &AnalysisResult{
				ObjectiveResults: []v1alpha3.ObjectiveResult{
					{
						PassResult: v1alpha3.CriteriaSetResult{
							Violated: false,
						},
						WarningResult: v1alpha3.CriteriaSetResult{
							Violated: false,
						},
						Value:        10.0,
						Score:        1,
						KeyObjective: false,
						Error:        nil,
					},
					{
						PassResult: v1alpha3.CriteriaSetResult{
							Violated: false,
						},
						WarningResult: v1alpha3.CriteriaSetResult{
							Violated: false,
						},
						Value:        10.0,
						Score:        2,
						KeyObjective: false,
						Error:        nil,
					},
				},
				TotalScore:   3,
				MaximumScore: 3,
				Pass:         true,
				Warning:      false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ae := &AnalysisEvaluator{
				ObjectiveEvaluator: tt.fields.ObjectiveEvaluator,
			}
			got, err := ae.Evaluate(tt.args.values, tt.args.ed)

			require.Nil(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
