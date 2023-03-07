package keptnevaluation

import (
	"testing"

	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	"github.com/stretchr/testify/require"
)

func TestCheckValue(t *testing.T) {
	tests := []struct {
		name   string
		obj    klcv1alpha3.Objective
		item   *klcv1alpha3.EvaluationStatusItem
		result bool
		err    bool
	}{
		{
			name:   "empty values",
			obj:    klcv1alpha3.Objective{},
			item:   &klcv1alpha3.EvaluationStatusItem{},
			result: false,
			err:    true,
		},
		{
			name: "garbage values",
			obj: klcv1alpha3.Objective{
				KeptnMetricRef: klcv1alpha3.KeptnMetricReference{
					Name:      "testytest",
					Namespace: "default",
				},
				EvaluationTarget: "testytest",
			},
			item: &klcv1alpha3.EvaluationStatusItem{
				Value:   "testytest",
				Status:  "testytest",
				Message: "testytest",
			},
			result: false,
			err:    true,
		},
		{
			name: "Item nan",
			obj: klcv1alpha3.Objective{
				KeptnMetricRef: klcv1alpha3.KeptnMetricReference{
					Name:      "testytest",
					Namespace: "default",
				},
				EvaluationTarget: "10",
			},
			item: &klcv1alpha3.EvaluationStatusItem{
				Value:   "nan",
				Status:  "all good",
				Message: "all good",
			},
			result: false,
			err:    false,
		},
		{
			name: "garbage comparison",
			obj: klcv1alpha3.Objective{
				KeptnMetricRef: klcv1alpha3.KeptnMetricReference{
					Name:      "testytest",
					Namespace: "default",
				},
				EvaluationTarget: "testytest",
			},
			item: &klcv1alpha3.EvaluationStatusItem{
				Value:   "10",
				Status:  "all good",
				Message: "all good",
			},
			result: false,
			err:    true,
		},
		{
			name: "objective nan",
			obj: klcv1alpha3.Objective{
				KeptnMetricRef: klcv1alpha3.KeptnMetricReference{
					Name:      "testytest",
					Namespace: "default",
				},
				EvaluationTarget: "nan",
			},
			item: &klcv1alpha3.EvaluationStatusItem{
				Value:   "10",
				Status:  "all good",
				Message: "all good",
			},
			result: false,
			err:    false,
		},
		{
			name: "10>10",
			obj: klcv1alpha3.Objective{
				KeptnMetricRef: klcv1alpha3.KeptnMetricReference{
					Name:      "testytest",
					Namespace: "default",
				},
				EvaluationTarget: ">10",
			},
			item: &klcv1alpha3.EvaluationStatusItem{
				Value:   "10",
				Status:  "all good",
				Message: "all good",
			},
			result: false,
			err:    false,
		},
		{
			name: "9>10",
			obj: klcv1alpha3.Objective{
				KeptnMetricRef: klcv1alpha3.KeptnMetricReference{
					Name:      "testytest",
					Namespace: "default",
				},
				EvaluationTarget: ">10",
			},
			item: &klcv1alpha3.EvaluationStatusItem{
				Value:   "9",
				Status:  "all good",
				Message: "all good",
			},
			result: false,
			err:    false,
		},
		{
			name: "11>10",
			obj: klcv1alpha3.Objective{
				KeptnMetricRef: klcv1alpha3.KeptnMetricReference{
					Name:      "testytest",
					Namespace: "default",
				},
				EvaluationTarget: ">10",
			},
			item: &klcv1alpha3.EvaluationStatusItem{
				Value:   "11",
				Status:  "all good",
				Message: "all good",
			},
			result: true,
			err:    false,
		},
		{
			name: "10<10",
			obj: klcv1alpha3.Objective{
				KeptnMetricRef: klcv1alpha3.KeptnMetricReference{
					Name:      "testytest",
					Namespace: "default",
				},
				EvaluationTarget: "<10",
			},
			item: &klcv1alpha3.EvaluationStatusItem{
				Value:   "10",
				Status:  "all good",
				Message: "all good",
			},
			result: false,
			err:    false,
		},
		{
			name: "9<10",
			obj: klcv1alpha3.Objective{
				KeptnMetricRef: klcv1alpha3.KeptnMetricReference{
					Name:      "testytest",
					Namespace: "default",
				},
				EvaluationTarget: "<10",
			},
			item: &klcv1alpha3.EvaluationStatusItem{
				Value:   "9",
				Status:  "all good",
				Message: "all good",
			},
			result: true,
			err:    false,
		},
		{
			name: "11<10",
			obj: klcv1alpha3.Objective{
				KeptnMetricRef: klcv1alpha3.KeptnMetricReference{
					Name:      "testytest",
					Namespace: "default",
				},
				EvaluationTarget: "<10",
			},
			item: &klcv1alpha3.EvaluationStatusItem{
				Value:   "11",
				Status:  "all good",
				Message: "all good",
			},
			result: false,
			err:    false,
		},
		{
			name: "invalid op",
			obj: klcv1alpha3.Objective{
				KeptnMetricRef: klcv1alpha3.KeptnMetricReference{
					Name:      "testytest",
					Namespace: "default",
				},
				EvaluationTarget: "-10",
			},
			item: &klcv1alpha3.EvaluationStatusItem{
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
