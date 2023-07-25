package slos

import (
	"fmt"
	"github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

const exampleSLO = `spec_version: '1.0'
filter:
  mz_id: "4711"
  svc_id: "a14b-cd87-0d51"
comparison:
  compare_with: "several_results"
  include_result_with_score: "pass"
  number_of_comparison_results: 3
  aggregate_function: avg
objectives:
- sli: response_time_p95
  displayName: "Response Time P95"
  pass:
  - criteria:
    - "<=+10%"
    - "<600"
  warning:
  - criteria:
    - "<=800"
total_score:
  pass: "90%"
  warning: "75%"`

func TestTransformSLOToEvaluationDefinition(t *testing.T) {
	type args struct {
		name         string
		sloConfigStr string
	}
	tests := []struct {
		name    string
		args    args
		want    *v1alpha3.KeptnEvaluationDefinition
		wantErr error
	}{
		{
			name: "convert SLO",
			args: args{
				name:         "my-slo",
				sloConfigStr: exampleSLO,
			},
			want: &v1alpha3.KeptnEvaluationDefinition{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-slo",
					Namespace: "keptn-lifecycle-toolkit-system",
				},
				Spec: v1alpha3.KeptnEvaluationDefinitionSpec{
					Objectives: []v1alpha3.Objective{
						{
							KeptnMetricRef: v1alpha3.KeptnMetricReference{
								Name:      "response_time_p95",
								Namespace: "keptn-lifecycle-toolkit-system",
							},
							EvaluationTarget: "",
							SLOTargets: &v1alpha3.SLOTarget{
								Pass: v1alpha3.CriteriaSet{
									AnyOf: []v1alpha3.Criteria{
										{
											AllOf: []v1alpha3.Target{
												{
													LessThanOrEqual: &v1alpha3.TargetValue{
														Comparison: &v1alpha3.ComparisonTarget{
															IncreaseByPercent: float32Ptr(10),
														},
													},
												},
												{
													LessThan: &v1alpha3.TargetValue{
														FixedValue: float32Ptr(600),
													},
												},
											},
										},
									},
								},
								Warning: v1alpha3.CriteriaSet{
									AnyOf: []v1alpha3.Criteria{
										{
											AllOf: []v1alpha3.Target{
												{
													LessThanOrEqual: &v1alpha3.TargetValue{
														FixedValue: float32Ptr(800),
													},
												},
											},
										},
									},
								},
							},
							Weight:       1,
							KeyObjective: false,
						},
					},
					TotalScore: &v1alpha3.PassThreshold{
						PassPercentage:    90,
						WarningPercentage: 75,
					},
					Comparison: &v1alpha3.ComparisonSpec{
						CompareWith:         "several_results",
						IncludeWarning:      false,
						NumberOfComparisons: 3,
						AggregationFunction: "avg",
					},
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TransformSLOToEvaluationDefinition(tt.args.name, tt.args.sloConfigStr)
			require.ErrorIs(t, err, tt.wantErr)

			require.Equal(t, tt.want.Spec.TotalScore, got.Spec.TotalScore)
			require.Equal(t, tt.want.Spec.Objectives, got.Spec.Objectives)
			require.Equal(t, tt.want.ObjectMeta, got.ObjectMeta)

			bytes, err := yaml.Marshal(got)

			fmt.Println(string(bytes))
		})
	}
}

func float32Ptr(value float32) *float32 {
	return &value
}
