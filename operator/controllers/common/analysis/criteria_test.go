package analysis

import (
	"github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/common/analysis/fake"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCriteriaEvaluator_Evaluate(t *testing.T) {
	type fields struct {
		TargetEvaluator *fake.ITargetEvaluatorMock
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

			if tt.args.c.AnyOf != nil {
				require.Len(t, tt.fields.TargetEvaluator.EvaluateCalls(), len(tt.args.c.AnyOf))
			} else {
				require.Len(t, tt.fields.TargetEvaluator.EvaluateCalls(), len(tt.args.c.AllOf))
			}

		})
	}
}

func TestCriteria_Evaluate(t *testing.T) {
	type fields struct {
		Criteria v1alpha3.Criteria
	}
	type args struct {
		val float64
	}
	tests := []struct {
		name       string
		fields     fields
		TargetImpl *fake.ITargetMock
		args       args
		want       v1alpha3.CriteriaResult
	}{
		{
			name: "anyOf - one target passes",
			fields: fields{
				Criteria: v1alpha3.Criteria{
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
			args: args{
				val: 0,
			},
			TargetImpl: &fake.ITargetMock{
				EvaluateFunc: func(val float64) v1alpha3.TargetResult {
					// one should pass, another one should fail
					violated := true

					return v1alpha3.TargetResult{
						Violated: violated,
					}
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
				Criteria: v1alpha3.Criteria{
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
			args: args{
				val: 0,
			},
			TargetImpl: &fake.ITargetMock{
				EvaluateFunc: func(val float64) v1alpha3.TargetResult {
					return v1alpha3.TargetResult{
						Violated: true,
					}
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

			targetFactory := &fake.ITargetFactoryMock{GetTargetFunc: func(target v1alpha3.Target) ITarget {
				return tt.TargetImpl
			}}
			c := &Criteria{
				Criteria:      tt.fields.Criteria,
				TargetFactory: targetFactory,
			}
			got := c.Evaluate(tt.args.val)

			require.Equal(t, tt.want, got)

			require.Len(t, targetFactory.GetTargetCalls(), len(tt.fields.Criteria.AnyOf))
		})
	}
}
