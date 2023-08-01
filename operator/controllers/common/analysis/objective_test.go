package analysis

import (
	"github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/common/analysis/fake"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestObjectiveEvaluator_Evaluate(t *testing.T) {
	type fields struct {
		CriteriaSetEvaluator *fake.ICriteriaSetEvaluatorMock
	}
	type args struct {
		values map[string]float64
		obj    v1alpha3.Objective
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   v1alpha3.ObjectiveResult
	}{
		{
			name: "pass objective reached",
			fields: fields{
				CriteriaSetEvaluator: &fake.ICriteriaSetEvaluatorMock{
					EvaluateFunc: func(val float64, criteriaSet v1alpha3.CriteriaSet) v1alpha3.CriteriaSetResult {
						return v1alpha3.CriteriaSetResult{
							ViolatedCriteria: nil,
							Violated:         false,
						}
					},
				},
			},
			args: args{
				values: map[string]float64{
					"my-value": 10.0,
				},
				obj: v1alpha3.Objective{
					KeptnMetricRef: v1alpha3.KeptnMetricReference{
						Name: "my-value",
					},
					SLOTargets: &v1alpha3.SLOTarget{
						Pass: v1alpha3.CriteriaSet{},
					},
					Weight:       2,
					KeyObjective: false,
				},
			},
			want: v1alpha3.ObjectiveResult{
				PassResult: v1alpha3.CriteriaSetResult{
					ViolatedCriteria: nil,
					Violated:         false,
				},
				Value:        10.0,
				Score:        2,
				KeyObjective: false,
				Error:        nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oe := &ObjectiveEvaluator{
				CriteriaSetEvaluator: tt.fields.CriteriaSetEvaluator,
			}
			got := oe.Evaluate(tt.args.values, tt.args.obj)

			require.Equal(t, tt.want, got)
		})
	}
}
