package analysis

import (
	"github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/common/analysis/fake"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCriteriaEvaluator_Evaluate(t *testing.T) {
	type fields struct {
		TargetEvaluator ITargetEvaluator
	}
	type args struct {
		val float64
		c   v1alpha3.Criteria
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   v1alpha3.CriteriaResult
	}{
		{
			name: "anyOf - one target passes",
			fields: fields{
				TargetEvaluator: &fake.ITargetEvaluatorMock{EvaluateFunc: func(val float64, target v1alpha3.Target) v1alpha3.TargetResult {
					// one should pass, another one should fail
					violated := true

					if target.EqualTo != nil {
						violated = false
					}

					return v1alpha3.TargetResult{
						Target:   target,
						Violated: violated,
					}
				}},
			},
			args: args{
				val: 0,
				c: v1alpha3.Criteria{
					AnyOf: []v1alpha3.Target{
						{
							EqualTo: &v1alpha3.TargetValue{},
						},
						{
							LessThan: &v1alpha3.TargetValue{},
						},
					},
				},
			},
			want: v1alpha3.CriteriaResult{
				ViolatedTargets: []v1alpha3.TargetResult{
					{
						Target: v1alpha3.Target{
							LessThan: &v1alpha3.TargetValue{},
						},
						Violated: true,
					},
				},
				Violated: false,
			},
		},
		{
			name: "anyOf - no target passes",
			fields: fields{
				TargetEvaluator: &fake.ITargetEvaluatorMock{EvaluateFunc: func(val float64, target v1alpha3.Target) v1alpha3.TargetResult {

					return v1alpha3.TargetResult{
						Target:   target,
						Violated: true,
					}
				}},
			},
			args: args{
				val: 0,
				c: v1alpha3.Criteria{
					AnyOf: []v1alpha3.Target{
						{
							EqualTo: &v1alpha3.TargetValue{},
						},
						{
							LessThan: &v1alpha3.TargetValue{},
						},
					},
				},
			},
			want: v1alpha3.CriteriaResult{
				ViolatedTargets: []v1alpha3.TargetResult{
					{
						Target: v1alpha3.Target{
							EqualTo: &v1alpha3.TargetValue{},
						},
						Violated: true,
					},
					{
						Target: v1alpha3.Target{
							LessThan: &v1alpha3.TargetValue{},
						},
						Violated: true,
					},
				},
				Violated: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ce := &CriteriaEvaluator{
				TargetEvaluator: tt.fields.TargetEvaluator,
			}
			got := ce.Evaluate(tt.args.val, tt.args.c)

			require.Equal(t, tt.want, got)
		})
	}
}
