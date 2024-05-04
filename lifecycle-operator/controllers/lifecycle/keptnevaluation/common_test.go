package keptnevaluation

import (
	"testing"

	apilifecycle "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1"
	"github.com/stretchr/testify/require"
)

func TestCheckValue(t *testing.T) {
	tests := []struct {
		name   string
		obj    apilifecycle.Objective
		item   *apilifecycle.EvaluationStatusItem
		result bool
		err    bool
	}{
		{
			name:   "empty values",
			obj:    apilifecycle.Objective{},
			item:   &apilifecycle.EvaluationStatusItem{},
			result: false,
			err:    true,
		},
		{
			name: "garbage values",
			obj: apilifecycle.Objective{
				KeptnMetricRef: apilifecycle.KeptnMetricReference{
					Name:      "testytest",
					Namespace: "default",
				},
				EvaluationTarget: "testytest",
			},
			item: &apilifecycle.EvaluationStatusItem{
				Value:   "testytest",
				Status:  "testytest",
				Message: "testytest",
			},
			result: false,
			err:    true,
		},
		{
			name: "Item nan",
			obj: apilifecycle.Objective{
				KeptnMetricRef: apilifecycle.KeptnMetricReference{
					Name:      "testytest",
					Namespace: "default",
				},
				EvaluationTarget: "10",
			},
			item: &apilifecycle.EvaluationStatusItem{
				Value:   "nan",
				Status:  "all good",
				Message: "all good",
			},
			result: false,
			err:    false,
		},
		{
			name: "garbage comparison",
			obj: apilifecycle.Objective{
				KeptnMetricRef: apilifecycle.KeptnMetricReference{
					Name:      "testytest",
					Namespace: "default",
				},
				EvaluationTarget: "testytest",
			},
			item: &apilifecycle.EvaluationStatusItem{
				Value:   "10",
				Status:  "all good",
				Message: "all good",
			},
			result: false,
			err:    true,
		},
		{
			name: "objective nan",
			obj: apilifecycle.Objective{
				KeptnMetricRef: apilifecycle.KeptnMetricReference{
					Name:      "testytest",
					Namespace: "default",
				},
				EvaluationTarget: "nan",
			},
			item: &apilifecycle.EvaluationStatusItem{
				Value:   "10",
				Status:  "all good",
				Message: "all good",
			},
			result: false,
			err:    false,
		},
		{
			name: "10>10",
			obj: apilifecycle.Objective{
				KeptnMetricRef: apilifecycle.KeptnMetricReference{
					Name:      "testytest",
					Namespace: "default",
				},
				EvaluationTarget: ">10",
			},
			item: &apilifecycle.EvaluationStatusItem{
				Value:   "10",
				Status:  "all good",
				Message: "all good",
			},
			result: false,
			err:    false,
		},
		{
			name: "9>10",
			obj: apilifecycle.Objective{
				KeptnMetricRef: apilifecycle.KeptnMetricReference{
					Name:      "testytest",
					Namespace: "default",
				},
				EvaluationTarget: ">10",
			},
			item: &apilifecycle.EvaluationStatusItem{
				Value:   "9",
				Status:  "all good",
				Message: "all good",
			},
			result: false,
			err:    false,
		},
		{
			name: "11>10",
			obj: apilifecycle.Objective{
				KeptnMetricRef: apilifecycle.KeptnMetricReference{
					Name:      "testytest",
					Namespace: "default",
				},
				EvaluationTarget: ">10",
			},
			item: &apilifecycle.EvaluationStatusItem{
				Value:   "11",
				Status:  "all good",
				Message: "all good",
			},
			result: true,
			err:    false,
		},
		{
			name: "10<10",
			obj: apilifecycle.Objective{
				KeptnMetricRef: apilifecycle.KeptnMetricReference{
					Name:      "testytest",
					Namespace: "default",
				},
				EvaluationTarget: "<10",
			},
			item: &apilifecycle.EvaluationStatusItem{
				Value:   "10",
				Status:  "all good",
				Message: "all good",
			},
			result: false,
			err:    false,
		},
		{
			name: "9<10",
			obj: apilifecycle.Objective{
				KeptnMetricRef: apilifecycle.KeptnMetricReference{
					Name:      "testytest",
					Namespace: "default",
				},
				EvaluationTarget: "<10",
			},
			item: &apilifecycle.EvaluationStatusItem{
				Value:   "9",
				Status:  "all good",
				Message: "all good",
			},
			result: true,
			err:    false,
		},
		{
			name: "11<10",
			obj: apilifecycle.Objective{
				KeptnMetricRef: apilifecycle.KeptnMetricReference{
					Name:      "testytest",
					Namespace: "default",
				},
				EvaluationTarget: "<10",
			},
			item: &apilifecycle.EvaluationStatusItem{
				Value:   "11",
				Status:  "all good",
				Message: "all good",
			},
			result: false,
			err:    false,
		},
		{
			name: "invalid op",
			obj: apilifecycle.Objective{
				KeptnMetricRef: apilifecycle.KeptnMetricReference{
					Name:      "testytest",
					Namespace: "default",
				},
				EvaluationTarget: "-10",
			},
			item: &apilifecycle.EvaluationStatusItem{
				Value:   "11",
				Status:  "all good",
				Message: "all good",
			},
			result: false,
			err:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, e := checkValue(tt.obj, tt.item)
			require.Equal(t, tt.result, r)
			if tt.err {
				require.NotNil(t, e)
			}
		})

	}
}
