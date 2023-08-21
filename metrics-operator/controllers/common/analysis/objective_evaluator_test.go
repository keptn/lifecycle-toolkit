package analysis

import (
	"errors"
	"testing"

	"github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/analysis/fake"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/analysis/types"
	"github.com/stretchr/testify/require"
)

func TestObjectiveEvaluator_Evaluate(t *testing.T) {
	tests := []struct {
		name            string
		values          map[string]string
		o               v1alpha3.Objective
		want            types.ObjectiveResult
		mockedEvaluator ITargetEvaluator
	}{
		{
			name:   "no value in results map",
			values: map[string]string{},
			o: v1alpha3.Objective{
				AnalysisValueTemplateRef: v1alpha3.ObjectReference{
					Name: "name",
				},
			},
			mockedEvaluator: &fake.ITargetEvaluatorMock{},
			want: types.ObjectiveResult{
				Score: 0.0,
				Error: errors.New("required value 'name' not available"),
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
				Weight: 2,
			},
			mockedEvaluator: &fake.ITargetEvaluatorMock{
				EvaluateFunc: func(val float64, target *v1alpha3.Target) types.TargetResult {
					return types.TargetResult{
						Pass: true,
					}
				},
			},
			want: types.ObjectiveResult{
				Score: 2.0,
				Error: nil,
				Value: 20.0,
				Result: types.TargetResult{
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
				Weight: 2,
			},
			mockedEvaluator: &fake.ITargetEvaluatorMock{
				EvaluateFunc: func(val float64, target *v1alpha3.Target) types.TargetResult {
					return types.TargetResult{
						Warning: true,
						Pass:    false,
					}
				},
			},
			want: types.ObjectiveResult{
				Score: 1.0,
				Error: nil,
				Value: 20.0,
				Result: types.TargetResult{
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
				Weight: 2,
			},
			mockedEvaluator: &fake.ITargetEvaluatorMock{
				EvaluateFunc: func(val float64, target *v1alpha3.Target) types.TargetResult {
					return types.TargetResult{
						Warning: false,
						Pass:    false,
					}
				},
			},
			want: types.ObjectiveResult{
				Score: 0.0,
				Error: nil,
				Value: 20.0,
				Result: types.TargetResult{
					Pass:    false,
					Warning: false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oe := NewObjectiveEvaluator(tt.mockedEvaluator)
			require.Equal(t, tt.want, oe.Evaluate(tt.values, &tt.o))
		})
	}
}

func TestGetValueFromMap(t *testing.T) {
	tests := []struct {
		name    string
		values  map[string]string
		in      string
		val     float64
		wantErr bool
	}{
		{
			name: "happy path",
			values: map[string]string{
				"key": "7",
			},
			in:      "key",
			val:     7.0,
			wantErr: false,
		},
		{
			name: "key not found",
			values: map[string]string{
				"key1": "7",
			},
			in:      "key",
			val:     0.0,
			wantErr: true,
		},
		{
			name: "value not float",
			values: map[string]string{
				"key": "",
			},
			in:      "key",
			val:     0.0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := getValueFromMap(tt.values, tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
			require.Equal(t, tt.val, res)
		})
	}
}

func TestComputeKey(t *testing.T) {
	obj := v1alpha3.ObjectReference{
		Name: "key",
	}

	require.Equal(t, "key", computeKey(obj))

	obj.Namespace = "namespace"

	require.Equal(t, "key-namespace", computeKey(obj))
}
