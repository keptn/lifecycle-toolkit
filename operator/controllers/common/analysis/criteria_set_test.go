package analysis

import (
	"github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/common/analysis/fake"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCriteriaSetEvaluator_Evaluate(t *testing.T) {
	type fields struct {
		CriteriaEvaluator *fake.ICriteriaEvaluatorMock
	}
	type args struct {
		val float64
		cs  v1alpha3.CriteriaSet
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   v1alpha3.CriteriaSetResult
	}{
		{
			name: "allOf - all pass",
			fields: fields{
				CriteriaEvaluator: &fake.ICriteriaEvaluatorMock{EvaluateFunc: func(val float64, criteria v1alpha3.Criteria) v1alpha3.CriteriaResult {
					// all criteria are passing
					return v1alpha3.CriteriaResult{
						Violated: false,
					}
				}},
			},
			args: args{
				val: 0,
				cs: v1alpha3.CriteriaSet{
					AllOf: []v1alpha3.Criteria{
						{},
						{},
					},
				},
			},
			want: v1alpha3.CriteriaSetResult{
				ViolatedCriteria: []v1alpha3.CriteriaResult{},
				Violated:         false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cse := &CriteriaSetEvaluator{
				CriteriaEvaluator: tt.fields.CriteriaEvaluator,
			}
			got := cse.Evaluate(tt.args.val, tt.args.cs)

			require.Equal(t, tt.want, got)

			if tt.args.cs.AllOf != nil {
				require.Len(t, tt.fields.CriteriaEvaluator.EvaluateCalls(), len(tt.args.cs.AllOf))
			} else {
				require.Len(t, tt.fields.CriteriaEvaluator.EvaluateCalls(), len(tt.args.cs.AnyOf))
			}
		})
	}
}
