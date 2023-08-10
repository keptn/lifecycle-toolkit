package analysis

import (
	"errors"
	"testing"

	"github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/analysis/fake"
	"github.com/stretchr/testify/require"
)

func TestObjectiveEvaluator_Evaluate(t *testing.T) {
	tests := []struct {
		name            string
		values          map[string]string
		o               v1alpha3.Objective
		want            v1alpha3.ObjectiveResult
		mockedEvaluator ITargetEvaluator
	}{
		{
			name:   "no value in results map",
			values: map[string]string{},
			o: v1alpha3.Objective{
				AnalysisValueTemplateRef: v1alpha3.ObjectReference{
					Name: "name",
				},
				KeyObjective: true,
			},
			mockedEvaluator: &fake.ITargetEvaluatorMock{},
			want: v1alpha3.ObjectiveResult{
				KeyObjective: true,
				Score:        0.0,
				Failed:       true,
				Error:        errors.New("required value not available"),
			},
		},
		{
			name: "evaluation passed",
			values: map[string]string{
				"name": "20",
			},
			o: v1alpha3.Objective{
				AnalysisValueTemplateRef: v1alpha3.ObjectReference{
					Name: "name",
				},
				KeyObjective: true,
				Weight:       2,
			},
			mockedEvaluator: &fake.ITargetEvaluatorMock{
				EvaluateFunc: func(val float64, target v1alpha3.Target) v1alpha3.TargetResult {
					return v1alpha3.TargetResult{
						Pass: true,
					}
				},
			},
			want: v1alpha3.ObjectiveResult{
				KeyObjective: true,
				Score:        2.0,
				Failed:       false,
				Error:        nil,
				Value:        20.0,
				Result: v1alpha3.TargetResult{
					Pass: true,
				},
			},
		},
		{
			name: "evaluation finished with warning",
			values: map[string]string{
				"name": "20",
			},
			o: v1alpha3.Objective{
				AnalysisValueTemplateRef: v1alpha3.ObjectReference{
					Name: "name",
				},
				KeyObjective: true,
				Weight:       2,
			},
			mockedEvaluator: &fake.ITargetEvaluatorMock{
				EvaluateFunc: func(val float64, target v1alpha3.Target) v1alpha3.TargetResult {
					return v1alpha3.TargetResult{
						Warning: true,
						Pass:    false,
					}
				},
			},
			want: v1alpha3.ObjectiveResult{
				KeyObjective: true,
				Score:        1.0,
				Failed:       false,
				Error:        nil,
				Value:        20.0,
				Result: v1alpha3.TargetResult{
					Pass:    false,
					Warning: true,
				},
			},
		},
		{
			name: "evaluation failed",
			values: map[string]string{
				"name": "20",
			},
			o: v1alpha3.Objective{
				AnalysisValueTemplateRef: v1alpha3.ObjectReference{
					Name: "name",
				},
				KeyObjective: true,
				Weight:       2,
			},
			mockedEvaluator: &fake.ITargetEvaluatorMock{
				EvaluateFunc: func(val float64, target v1alpha3.Target) v1alpha3.TargetResult {
					return v1alpha3.TargetResult{
						Warning: false,
						Pass:    false,
					}
				},
			},
			want: v1alpha3.ObjectiveResult{
				KeyObjective: true,
				Score:        0.0,
				Failed:       true,
				Error:        nil,
				Value:        20.0,
				Result: v1alpha3.TargetResult{
					Pass:    false,
					Warning: false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oe := NewObjectiveEvaluator(tt.mockedEvaluator)
			require.Equal(t, tt.want, oe.Evaluate(tt.values, tt.o))
		})
	}
}
