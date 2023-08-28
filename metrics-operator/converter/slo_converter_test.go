package converter

import (
	"testing"

	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/resource"
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

func TestConvertSLO(t *testing.T) {
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
					FixedValue: *resource.NewQuantity(1, resource.DecimalSI),
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
					FixedValue: *resource.NewQuantity(1, resource.DecimalSI),
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
					FixedValue: *resource.NewQuantity(1, resource.DecimalSI),
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
					FixedValue: *resource.NewQuantity(1, resource.DecimalSI),
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
					FixedValue: *resource.NewQuantity(1, resource.DecimalSI),
				},
			},
			wantErr: false,
		},
		{
			name: "less operator",
			op:   "<1",
			out: &metricsapi.Operator{
				GreaterThanOrEqual: &metricsapi.OperatorValue{
					FixedValue: *resource.NewQuantity(1, resource.DecimalSI),
				},
			},
			wantErr: false,
		},
		{
			name: "greaterEqual operator",
			op:   ">=1",
			out: &metricsapi.Operator{
				LessThan: &metricsapi.OperatorValue{
					FixedValue: *resource.NewQuantity(1, resource.DecimalSI),
				},
			},
			wantErr: false,
		},
		{
			name: "greater operator",
			op:   ">1",
			out: &metricsapi.Operator{
				LessThanOrEqual: &metricsapi.OperatorValue{
					FixedValue: *resource.NewQuantity(1, resource.DecimalSI),
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := setupOperator(tt.op)
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
			require.Equal(t, tt.want, shouldIgnoreObjective(tt.o))
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
