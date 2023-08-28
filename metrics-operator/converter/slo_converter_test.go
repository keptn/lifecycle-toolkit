package converter

import (
	"testing"

	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	"github.com/stretchr/testify/require"
	"gopkg.in/inf.v0"
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const SLOContent = `---
spec_version: "0.1.1"
comparison:
  aggregate_function: "avg"
  compare_with: "single_result"
  include_result_with_score: "pass"
  number_of_comparison_results: 1
filter:
objectives:
  - sli: "response_time_p95"
    key_sli: false
    pass:
      - criteria:
          - "<=+75%"
          - "<800"
    warning:
      - criteria:
          - "<=1000"
          - "<=+100%"
    weight: 1
  - sli: "cpu"
    pass:
      - criteria:
          - "<=+100%"
          - ">=80"
      - criteria:
          - "<=+100%"
          - ">=80"
  - sli: "throughput"
    pass:
      - criteria:
          - "<=+100%"
          - ">=-80%"
  - sli: "error_rate"
total_score:
  pass: "100%"
  warning: "65%"`

const expectedOutput = `apiVersion: metrics.keptn.sh/v1alpha3
kind: AnalysisDefinition
metadata:
  creationTimestamp: null
  name: defname
spec:
  objectives:
  - analysisValueTemplateRef:
      name: response_time_p95
      namespace: default
    target:
      failure:
        greaterThan:
          fixedValue: 1k
      warning:
        greaterThanOrEqual:
          fixedValue: "800"
    weight: 1
  - analysisValueTemplateRef:
      name: cpu
      namespace: default
    target: {}
  - analysisValueTemplateRef:
      name: throughput
      namespace: default
    target: {}
  - analysisValueTemplateRef:
      name: error_rate
      namespace: default
    target: {}
  totalScore:
    passPercentage: 100
    warningPercentage: 65
`

func TestConvert(t *testing.T) {
	c := NewSLOConverter()
	// no provider nor namespace
	res, err := c.Convert([]byte(SLOContent), "", "")
	require.NotNil(t, err)
	require.Equal(t, "", res)

	// invalid file content
	res, err = c.Convert([]byte("invalid"), "dynatrace", "keptn")
	require.NotNil(t, err)
	require.Equal(t, "", res)

	// happy path
	res, err = c.Convert([]byte(SLOContent), "defname", "default")
	require.Nil(t, err)
	require.Equal(t, expectedOutput, res)
}

func TestCreateOperator(t *testing.T) {
	dec := inf.NewDec(1, 0)
	_, ok := dec.SetString("1")
	require.True(t, ok)

	tests := []struct {
		name    string
		op      string
		value   string
		out     *metricsapi.Operator
		wantErr bool
	}{
		{
			name:    "invalid int value",
			op:      "",
			value:   "val",
			out:     nil,
			wantErr: true,
		},
		{
			name:    "unsupported operator",
			op:      "",
			value:   "1",
			out:     nil,
			wantErr: true,
		},
		{
			name:  "lessEqual operator",
			op:    "<=",
			value: "1",
			out: &metricsapi.Operator{
				GreaterThan: &metricsapi.OperatorValue{
					FixedValue: *resource.NewDecimalQuantity(*dec, resource.DecimalSI),
				},
			},
			wantErr: false,
		},
		{
			name:  "less operator",
			op:    "<",
			value: "1",
			out: &metricsapi.Operator{
				GreaterThanOrEqual: &metricsapi.OperatorValue{
					FixedValue: *resource.NewDecimalQuantity(*dec, resource.DecimalSI),
				},
			},
			wantErr: false,
		},
		{
			name:  "greaterEqual operator",
			op:    ">=",
			value: "1",
			out: &metricsapi.Operator{
				LessThan: &metricsapi.OperatorValue{
					FixedValue: *resource.NewDecimalQuantity(*dec, resource.DecimalSI),
				},
			},
			wantErr: false,
		},
		{
			name:  "greater operator",
			op:    ">",
			value: "1",
			out: &metricsapi.Operator{
				LessThanOrEqual: &metricsapi.OperatorValue{
					FixedValue: *resource.NewDecimalQuantity(*dec, resource.DecimalSI),
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := createOperator(tt.op, tt.value)
			if tt.wantErr {
				require.NotNil(t, err)
			} else {
				require.Equal(t, tt.out, res)
			}
		})

	}
}

func TestSetupOperator(t *testing.T) {
	dec := inf.NewDec(1, 0)
	_, ok := dec.SetString("1")
	require.True(t, ok)

	tests := []struct {
		name    string
		op      string
		out     *metricsapi.Operator
		wantErr bool
	}{
		{
			name:    "unsupported operator",
			op:      "",
			out:     nil,
			wantErr: true,
		},
		{
			name: "lessEqual operator",
			op:   "<=1",
			out: &metricsapi.Operator{
				GreaterThan: &metricsapi.OperatorValue{
					FixedValue: *resource.NewDecimalQuantity(*dec, resource.DecimalSI),
				},
			},
			wantErr: false,
		},
		{
			name: "less operator",
			op:   "<1",
			out: &metricsapi.Operator{
				GreaterThanOrEqual: &metricsapi.OperatorValue{
					FixedValue: *resource.NewDecimalQuantity(*dec, resource.DecimalSI),
				},
			},
			wantErr: false,
		},
		{
			name: "greaterEqual operator",
			op:   ">=1",
			out: &metricsapi.Operator{
				LessThan: &metricsapi.OperatorValue{
					FixedValue: *resource.NewDecimalQuantity(*dec, resource.DecimalSI),
				},
			},
			wantErr: false,
		},
		{
			name: "greater operator",
			op:   ">1",
			out: &metricsapi.Operator{
				LessThanOrEqual: &metricsapi.OperatorValue{
					FixedValue: *resource.NewDecimalQuantity(*dec, resource.DecimalSI),
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := newOperator(tt.op)
			if tt.wantErr {
				require.NotNil(t, err)
			} else {
				require.Equal(t, tt.out, res)
			}
		})

	}
}

func TestShouldIgnoreObjective(t *testing.T) {
	tests := []struct {
		name string
		o    *Objective
		want bool
	}{
		{
			name: "empty criteria",
			o: &Objective{
				Pass:    []Criteria{},
				Warning: []Criteria{},
			},
			want: false,
		},
		{
			name: "valid criteria",
			o: &Objective{
				Pass: []Criteria{
					{
						Operators: []string{},
					},
				},
				Warning: []Criteria{
					{
						Operators: []string{},
					},
				},
			},
			want: false,
		},
		{
			name: "OR criteria",
			o: &Objective{
				Pass: []Criteria{
					{
						Operators: []string{},
					},
					{
						Operators: []string{},
					},
				},
				Warning: []Criteria{
					{
						Operators: []string{},
					},
				},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.o.hasSupportedCriteria())
		})

	}
}

func TestRemovePercentage(t *testing.T) {
	tests := []struct {
		name    string
		val     string
		want    int
		wantErr bool
	}{
		{
			name:    "no percentage",
			val:     "1",
			want:    1,
			wantErr: false,
		},
		{
			name:    "percentage",
			val:     "1%",
			want:    1,
			wantErr: false,
		},
		{
			name:    "only percentage",
			val:     "%",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := removePercentage(tt.val)
			if tt.wantErr {
				require.NotNil(t, err)
			} else {
				require.Equal(t, tt.want, res)
			}
		})
	}
}

func TestCleanupCriteria(t *testing.T) {
	tests := []struct {
		name string
		in   []Criteria
		out  []Criteria
	}{
		{
			name: "empty criteria",
			in:   []Criteria{},
			out:  []Criteria{},
		},
		{
			name: "no criteria to clean up",
			in: []Criteria{
				{
					Operators: []string{"<100"},
				},
			},
			out: []Criteria{
				{
					Operators: []string{"<100"},
				},
			},
		},
		{
			name: "criteria to clean up",
			in: []Criteria{
				{
					Operators: []string{"<100", "<10%"},
				},
			},
			out: []Criteria{
				{
					Operators: []string{"<100"},
				},
			},
		},
		{
			name: "multiple criteria to clean up",
			in: []Criteria{
				{
					Operators: []string{"<100", "<10%"},
				},
				{
					Operators: []string{"<10%"},
				},
			},
			out: []Criteria{
				{
					Operators: []string{"<100"},
				},
			},
		},
		{
			name: "all criteria to clean up",
			in: []Criteria{
				{
					Operators: []string{"<10%"},
				},
				{
					Operators: []string{"<10%"},
				},
			},
			out: []Criteria{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.out, cleanupCriteria(tt.in))
		})
	}
}

func TestCleanupObjective(t *testing.T) {
	tests := []struct {
		name string
		in   *Objective
		out  *Objective
	}{
		{
			name: "pass criteria to clean up",
			in: &Objective{
				Pass: []Criteria{
					{
						Operators: []string{"<100", "<10%"},
					},
					{
						Operators: []string{"<10%"},
					},
				},
			},
			out: &Objective{
				Pass: []Criteria{
					{
						Operators: []string{"<100"},
					},
				},
				Warning: []Criteria{},
			},
		},
		{
			name: "warning criteria to clean up",
			in: &Objective{
				Warning: []Criteria{
					{
						Operators: []string{"<100", "<10%"},
					},
					{
						Operators: []string{"<10%"},
					},
				},
			},
			out: &Objective{
				Warning: []Criteria{
					{
						Operators: []string{"<100"},
					},
				},
				Pass: []Criteria{},
			},
		},
		{
			name: "no criteria to clean up",
			in: &Objective{
				Warning: []Criteria{
					{
						Operators: []string{"<100"},
					},
				},
			},
			out: &Objective{
				Warning: []Criteria{
					{
						Operators: []string{"<100"},
					},
				},
				Pass: []Criteria{},
			},
		},
		{
			name: "no criteria",
			in:   &Objective{},
			out: &Objective{
				Warning: []Criteria{},
				Pass:    []Criteria{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.out, cleanupObjective(tt.in))
		})
	}
}

func TestSetupTarget(t *testing.T) {
	dec10 := inf.NewDec(1, 0)
	_, ok := dec10.SetString("10")
	require.True(t, ok)

	dec15 := inf.NewDec(1, 0)
	_, ok = dec15.SetString("15")
	require.True(t, ok)

	tests := []struct {
		name    string
		o       *Objective
		want    *metricsapi.Target
		wantErr bool
	}{
		{
			name: "informative slo",
			o: &Objective{
				Name: "informative",
			},
			want:    &metricsapi.Target{},
			wantErr: false,
		},
		{
			name: "logical OR criteria -> informative slo",
			o: &Objective{
				Name: "criteria",
				Pass: []Criteria{
					{
						Operators: []string{"<10"},
					},
					{
						Operators: []string{">1"},
					},
				},
			},
			want:    &metricsapi.Target{},
			wantErr: false,
		},
		{
			name: "criteria with % -> informative slo",
			o: &Objective{
				Name: "criteria",
				Pass: []Criteria{
					{
						Operators: []string{"<10%"},
					},
				},
			},
			want:    &metricsapi.Target{},
			wantErr: false,
		},
		{
			name: "no warn criteria - conversion",
			o: &Objective{
				Name: "criteria",
				Pass: []Criteria{
					{
						Operators: []string{"<10"},
					},
				},
			},
			want: &metricsapi.Target{
				Failure: &metricsapi.Operator{
					GreaterThanOrEqual: &metricsapi.OperatorValue{
						FixedValue: *resource.NewDecimalQuantity(*dec10, resource.DecimalSI),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "no warn criteria - error",
			o: &Objective{
				Name: "criteria",
				Pass: []Criteria{
					{
						Operators: []string{"10"},
					},
				},
			},
			want:    &metricsapi.Target{},
			wantErr: true,
		},
		{
			name: "with warn criteria - conversion",
			o: &Objective{
				Name: "criteria",
				Pass: []Criteria{
					{
						Operators: []string{"<10"},
					},
				},
				Warning: []Criteria{
					{
						Operators: []string{"<=15"},
					},
				},
			},
			want: &metricsapi.Target{
				Failure: &metricsapi.Operator{
					GreaterThan: &metricsapi.OperatorValue{
						FixedValue: *resource.NewDecimalQuantity(*dec15, resource.DecimalSI),
					},
				},
				Warning: &metricsapi.Operator{
					GreaterThanOrEqual: &metricsapi.OperatorValue{
						FixedValue: *resource.NewDecimalQuantity(*dec10, resource.DecimalSI),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "with warn criteria - error pass",
			o: &Objective{
				Name: "criteria",
				Pass: []Criteria{
					{
						Operators: []string{"10"},
					},
				},
				Warning: []Criteria{
					{
						Operators: []string{"<=15"},
					},
				},
			},
			want:    &metricsapi.Target{},
			wantErr: true,
		},
		{
			name: "with warn criteria - error warn",
			o: &Objective{
				Name: "criteria",
				Pass: []Criteria{
					{
						Operators: []string{"<10"},
					},
				},
				Warning: []Criteria{
					{
						Operators: []string{"15"},
					},
				},
			},
			want:    &metricsapi.Target{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := setupTarget(tt.o)
			if tt.wantErr {
				require.NotNil(t, err)
			} else {
				require.Equal(t, tt.want, res)
			}
		})
	}
}

func TestConvertSLO(t *testing.T) {
	dec10 := inf.NewDec(1, 0)
	_, ok := dec10.SetString("10")
	require.True(t, ok)

	dec15 := inf.NewDec(1, 0)
	_, ok = dec15.SetString("15")
	require.True(t, ok)

	c := NewSLOConverter()

	tests := []struct {
		name      string
		slo       *SLO
		defName   string
		namespace string
		out       *metricsapi.AnalysisDefinition
		wantErr   bool
	}{
		{
			name: "removePercentage pass - error",
			slo: &SLO{
				TotalScore: Score{
					Pass: "hhh",
				},
			},
			defName:   "defName",
			namespace: "default",
			out:       &metricsapi.AnalysisDefinition{},
			wantErr:   true,
		},
		{
			name: "removePercentage warn - error",
			slo: &SLO{
				TotalScore: Score{
					Pass:    "10",
					Warning: "hhh",
				},
			},
			defName:   "defName",
			namespace: "default",
			out:       &metricsapi.AnalysisDefinition{},
			wantErr:   true,
		},
		{
			name: "no objectives",
			slo: &SLO{
				TotalScore: Score{
					Pass:    "50",
					Warning: "20",
				},
			},
			defName:   "defName",
			namespace: "default",
			out: &metricsapi.AnalysisDefinition{
				TypeMeta: v1.TypeMeta{
					Kind:       "AnalysisDefinition",
					APIVersion: "metrics.keptn.sh/v1alpha3",
				},
				ObjectMeta: v1.ObjectMeta{
					Name: "defName",
				},
				Spec: metricsapi.AnalysisDefinitionSpec{
					TotalScore: metricsapi.TotalScore{
						PassPercentage:    50,
						WarningPercentage: 20,
					},
					Objectives: []metricsapi.Objective{},
				},
			},
			wantErr: false,
		},
		{
			name: "objectives conversion",
			slo: &SLO{
				TotalScore: Score{
					Pass:    "50",
					Warning: "20",
				},
				Objectives: []*Objective{
					{
						Name: "criteria",
						Pass: []Criteria{
							{
								Operators: []string{"<10"},
							},
						},
						Warning: []Criteria{
							{
								Operators: []string{"<=15"},
							},
						},
						KeySLI: true,
						Weight: 10,
					},
					{
						Name: "criteria2",
						Pass: []Criteria{
							{
								Operators: []string{"<10"},
							},
						},
						Weight: 5,
					},
				},
			},
			defName:   "defName",
			namespace: "default",
			out: &metricsapi.AnalysisDefinition{
				TypeMeta: v1.TypeMeta{
					Kind:       "AnalysisDefinition",
					APIVersion: "metrics.keptn.sh/v1alpha3",
				},
				ObjectMeta: v1.ObjectMeta{
					Name: "defName",
				},
				Spec: metricsapi.AnalysisDefinitionSpec{
					TotalScore: metricsapi.TotalScore{
						PassPercentage:    50,
						WarningPercentage: 20,
					},
					Objectives: []metricsapi.Objective{
						{
							Target: metricsapi.Target{
								Failure: &metricsapi.Operator{
									GreaterThan: &metricsapi.OperatorValue{
										FixedValue: *resource.NewDecimalQuantity(*dec15, resource.DecimalSI),
									},
								},
								Warning: &metricsapi.Operator{
									GreaterThanOrEqual: &metricsapi.OperatorValue{
										FixedValue: *resource.NewDecimalQuantity(*dec10, resource.DecimalSI),
									},
								},
							},
							Weight:       10,
							KeyObjective: true,
							AnalysisValueTemplateRef: metricsapi.ObjectReference{
								Name:      "criteria",
								Namespace: "default",
							},
						},
						{
							Target: metricsapi.Target{
								Failure: &metricsapi.Operator{
									GreaterThanOrEqual: &metricsapi.OperatorValue{
										FixedValue: *resource.NewDecimalQuantity(*dec10, resource.DecimalSI),
									},
								},
							},
							Weight: 5,
							AnalysisValueTemplateRef: metricsapi.ObjectReference{
								Name:      "criteria2",
								Namespace: "default",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "objectives conversion - setupTarget error",
			slo: &SLO{
				TotalScore: Score{
					Pass:    "50",
					Warning: "20",
				},
				Objectives: []*Objective{
					{
						Name: "criteria",
						Pass: []Criteria{
							{
								Operators: []string{"<10"},
							},
						},
						Warning: []Criteria{
							{
								Operators: []string{"15"},
							},
						},
						KeySLI: true,
						Weight: 10,
					},
					{
						Name: "criteria2",
						Pass: []Criteria{
							{
								Operators: []string{"<10"},
							},
						},
						Weight: 5,
					},
				},
			},
			defName:   "defName",
			namespace: "default",
			out:       &metricsapi.AnalysisDefinition{},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := c.convertSLO(tt.slo, tt.defName, tt.namespace)
			if tt.wantErr {
				require.NotNil(t, err)
			} else {
				require.Equal(t, tt.out, res)
			}
		})
	}

}
