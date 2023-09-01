//nolint:dupl
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
  - sli: "response_time_p90"
    key_sli: false
    pass:
    - criteria:
        - ">600"
        - "<800"
    warning:
    - criteria:
        - "<=1000"
        - ">500"
    weight: 2
  - sli: "response_time_p80"
    key_sli: false
    pass:
      - criteria:
          - ">600"
          - "<800"
    warning:
      - criteria:
          - "<=1000"
    weight: 2
  - sli: "response_time_p70"
    key_sli: false
    warning:
      - criteria:
          - ">600"
          - "<800"
    pass:
      - criteria:
          - "<=1000"
    weight: 2
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
      name: response_time_p90
      namespace: default
    target:
      failure:
        notInRange:
          highBound: 1k
          lowBound: "500"
      warning:
        notInRange:
          highBound: "800"
          lowBound: "600"
    weight: 2
  - analysisValueTemplateRef:
      name: response_time_p80
      namespace: default
    target:
      failure:
        greaterThan:
          fixedValue: 1k
      warning:
        notInRange:
          highBound: "800"
          lowBound: "600"
    weight: 2
  - analysisValueTemplateRef:
      name: response_time_p70
      namespace: default
    target:
      failure:
        greaterThan:
          fixedValue: 1k
      warning:
        inRange:
          highBound: "800"
          lowBound: "600"
    weight: 2
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

func TestNegateSingleOperator(t *testing.T) {
	dec := inf.NewDec(1, 0)

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
			res, err := negateSingleOperator(tt.op, tt.value)
			if tt.wantErr {
				require.NotNil(t, err)
			} else {
				require.Equal(t, tt.out, res)
			}
		})

	}
}

func TestCreateSingleOperator(t *testing.T) {
	dec := inf.NewDec(1, 0)

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
				LessThanOrEqual: &metricsapi.OperatorValue{
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
				LessThan: &metricsapi.OperatorValue{
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
				GreaterThanOrEqual: &metricsapi.OperatorValue{
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
				GreaterThan: &metricsapi.OperatorValue{
					FixedValue: *resource.NewDecimalQuantity(*dec, resource.DecimalSI),
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := createSingleOperator(tt.op, tt.value)
			if tt.wantErr {
				require.NotNil(t, err)
			} else {
				require.Equal(t, tt.out, res)
			}
		})

	}
}

func TestCreateDoubleOperator(t *testing.T) {
	dec := inf.NewDec(1, 0)
	dec5 := inf.NewDec(5, 0)

	tests := []struct {
		name    string
		op1     string
		value1  string
		op2     string
		value2  string
		out     *metricsapi.Operator
		wantErr bool
	}{
		{
			name:    "invalid int value",
			op1:     "",
			value1:  "val",
			out:     nil,
			wantErr: true,
		},
		{
			name:    "unsupported operator",
			op1:     "",
			value1:  "1",
			out:     nil,
			wantErr: true,
		},
		{
			name:   "inRange operator",
			op1:    "<=",
			value1: "5",
			op2:    ">=",
			value2: "1",
			out: &metricsapi.Operator{
				InRange: &metricsapi.RangeValue{
					LowBound:  *resource.NewDecimalQuantity(*dec, resource.DecimalSI),
					HighBound: *resource.NewDecimalQuantity(*dec5, resource.DecimalSI),
				},
			},
			wantErr: false,
		},
		{
			name:   "inRange operator",
			op1:    "<",
			value1: "5",
			op2:    ">",
			value2: "1",
			out: &metricsapi.Operator{
				InRange: &metricsapi.RangeValue{
					LowBound:  *resource.NewDecimalQuantity(*dec, resource.DecimalSI),
					HighBound: *resource.NewDecimalQuantity(*dec5, resource.DecimalSI),
				},
			},
			wantErr: false,
		},
		{
			name:   "inRange operator",
			op1:    "<=",
			value1: "5",
			op2:    ">",
			value2: "1",
			out: &metricsapi.Operator{
				InRange: &metricsapi.RangeValue{
					LowBound:  *resource.NewDecimalQuantity(*dec, resource.DecimalSI),
					HighBound: *resource.NewDecimalQuantity(*dec5, resource.DecimalSI),
				},
			},
			wantErr: false,
		},
		{
			name:   "inRange operator",
			op1:    "<",
			value1: "5",
			op2:    ">=",
			value2: "1",
			out: &metricsapi.Operator{
				InRange: &metricsapi.RangeValue{
					LowBound:  *resource.NewDecimalQuantity(*dec, resource.DecimalSI),
					HighBound: *resource.NewDecimalQuantity(*dec5, resource.DecimalSI),
				},
			},
			wantErr: false,
		},
		{
			name:   "notinRange operator",
			op1:    ">=",
			value1: "5",
			op2:    "<=",
			value2: "1",
			out: &metricsapi.Operator{
				NotInRange: &metricsapi.RangeValue{
					LowBound:  *resource.NewDecimalQuantity(*dec, resource.DecimalSI),
					HighBound: *resource.NewDecimalQuantity(*dec5, resource.DecimalSI),
				},
			},
			wantErr: false,
		},
		{
			name:   "notinRange operator",
			op1:    ">",
			value1: "5",
			op2:    "<",
			value2: "1",
			out: &metricsapi.Operator{
				NotInRange: &metricsapi.RangeValue{
					LowBound:  *resource.NewDecimalQuantity(*dec, resource.DecimalSI),
					HighBound: *resource.NewDecimalQuantity(*dec5, resource.DecimalSI),
				},
			},
			wantErr: false,
		},
		{
			name:   "notinRange operator",
			op1:    ">=",
			value1: "5",
			op2:    "<",
			value2: "1",
			out: &metricsapi.Operator{
				NotInRange: &metricsapi.RangeValue{
					LowBound:  *resource.NewDecimalQuantity(*dec, resource.DecimalSI),
					HighBound: *resource.NewDecimalQuantity(*dec5, resource.DecimalSI),
				},
			},
			wantErr: false,
		},
		{
			name:   "notinRange operator",
			op1:    ">",
			value1: "5",
			op2:    "<=",
			value2: "1",
			out: &metricsapi.Operator{
				NotInRange: &metricsapi.RangeValue{
					LowBound:  *resource.NewDecimalQuantity(*dec, resource.DecimalSI),
					HighBound: *resource.NewDecimalQuantity(*dec5, resource.DecimalSI),
				},
			},
			wantErr: false,
		},
		{
			name:    "unsupported combination",
			op1:     ">",
			value1:  "5",
			op2:     ">",
			value2:  "1",
			out:     &metricsapi.Operator{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := createDoubleOperator(tt.op1, tt.value1, tt.op2, tt.value2)
			if tt.wantErr {
				require.NotNil(t, err)
			} else {
				require.Equal(t, tt.out, res)
			}
		})

	}
}

func TestNegateDoubleOperator(t *testing.T) {
	dec := inf.NewDec(1, 0)
	dec5 := inf.NewDec(5, 0)

	tests := []struct {
		name    string
		op1     string
		value1  string
		op2     string
		value2  string
		out     *metricsapi.Operator
		wantErr bool
	}{
		{
			name:    "invalid int value",
			op1:     "",
			value1:  "val",
			out:     nil,
			wantErr: true,
		},
		{
			name:    "unsupported operator",
			op1:     "",
			value1:  "1",
			out:     nil,
			wantErr: true,
		},
		{
			name:   "Not inRange operator",
			op1:    "<=",
			value1: "5",
			op2:    ">=",
			value2: "1",
			out: &metricsapi.Operator{
				NotInRange: &metricsapi.RangeValue{
					LowBound:  *resource.NewDecimalQuantity(*dec, resource.DecimalSI),
					HighBound: *resource.NewDecimalQuantity(*dec5, resource.DecimalSI),
				},
			},
			wantErr: false,
		},
		{
			name:   "Not inRange operator",
			op1:    "<",
			value1: "5",
			op2:    ">",
			value2: "1",
			out: &metricsapi.Operator{
				NotInRange: &metricsapi.RangeValue{
					LowBound:  *resource.NewDecimalQuantity(*dec, resource.DecimalSI),
					HighBound: *resource.NewDecimalQuantity(*dec5, resource.DecimalSI),
				},
			},
			wantErr: false,
		},
		{
			name:   "Not inRange operator",
			op1:    "<=",
			value1: "5",
			op2:    ">",
			value2: "1",
			out: &metricsapi.Operator{
				NotInRange: &metricsapi.RangeValue{
					LowBound:  *resource.NewDecimalQuantity(*dec, resource.DecimalSI),
					HighBound: *resource.NewDecimalQuantity(*dec5, resource.DecimalSI),
				},
			},
			wantErr: false,
		},
		{
			name:   "Not inRange operator",
			op1:    "<",
			value1: "5",
			op2:    ">=",
			value2: "1",
			out: &metricsapi.Operator{
				NotInRange: &metricsapi.RangeValue{
					LowBound:  *resource.NewDecimalQuantity(*dec, resource.DecimalSI),
					HighBound: *resource.NewDecimalQuantity(*dec5, resource.DecimalSI),
				},
			},
			wantErr: false,
		},
		{
			name:   "inRange operator",
			op1:    ">=",
			value1: "5",
			op2:    "<=",
			value2: "1",
			out: &metricsapi.Operator{
				InRange: &metricsapi.RangeValue{
					LowBound:  *resource.NewDecimalQuantity(*dec, resource.DecimalSI),
					HighBound: *resource.NewDecimalQuantity(*dec5, resource.DecimalSI),
				},
			},
			wantErr: false,
		},
		{
			name:   "inRange operator",
			op1:    ">",
			value1: "5",
			op2:    "<",
			value2: "1",
			out: &metricsapi.Operator{
				InRange: &metricsapi.RangeValue{
					LowBound:  *resource.NewDecimalQuantity(*dec, resource.DecimalSI),
					HighBound: *resource.NewDecimalQuantity(*dec5, resource.DecimalSI),
				},
			},
			wantErr: false,
		},
		{
			name:   "inRange operator",
			op1:    ">=",
			value1: "5",
			op2:    "<",
			value2: "1",
			out: &metricsapi.Operator{
				InRange: &metricsapi.RangeValue{
					LowBound:  *resource.NewDecimalQuantity(*dec, resource.DecimalSI),
					HighBound: *resource.NewDecimalQuantity(*dec5, resource.DecimalSI),
				},
			},
			wantErr: false,
		},
		{
			name:   "inRange operator",
			op1:    ">",
			value1: "5",
			op2:    "<=",
			value2: "1",
			out: &metricsapi.Operator{
				InRange: &metricsapi.RangeValue{
					LowBound:  *resource.NewDecimalQuantity(*dec, resource.DecimalSI),
					HighBound: *resource.NewDecimalQuantity(*dec5, resource.DecimalSI),
				},
			},
			wantErr: false,
		},
		{
			name:    "unsupported combination",
			op1:     ">",
			value1:  "5",
			op2:     ">",
			value2:  "1",
			out:     &metricsapi.Operator{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := negateDoubleOperator(tt.op1, tt.value1, tt.op2, tt.value2)
			if tt.wantErr {
				require.NotNil(t, err)
			} else {
				require.Equal(t, tt.out, res)
			}
		})

	}
}

func TestNewOperator(t *testing.T) {
	dec := inf.NewDec(1, 0)
	dec5 := inf.NewDec(5, 0)

	tests := []struct {
		name    string
		op      []string
		negate  bool
		out     *metricsapi.Operator
		wantErr bool
	}{
		{
			name:    "empty operator",
			op:      []string{},
			negate:  true,
			out:     nil,
			wantErr: true,
		},
		{
			name:    "unsupported operator",
			op:      []string{""},
			negate:  true,
			out:     nil,
			wantErr: true,
		},
		{
			name:    "unsupported operator double - first",
			op:      []string{"", ">5"},
			negate:  true,
			out:     nil,
			wantErr: true,
		},
		{
			name:    "unsupported operator double - second",
			op:      []string{"<5", "5"},
			negate:  true,
			out:     nil,
			wantErr: true,
		},
		{
			name:   "lessEqual operator - negate",
			op:     []string{"<=1"},
			negate: true,
			out: &metricsapi.Operator{
				GreaterThan: &metricsapi.OperatorValue{
					FixedValue: *resource.NewDecimalQuantity(*dec, resource.DecimalSI),
				},
			},
			wantErr: false,
		},
		{
			name:   "less operator - negate",
			op:     []string{"<1"},
			negate: true,
			out: &metricsapi.Operator{
				GreaterThanOrEqual: &metricsapi.OperatorValue{
					FixedValue: *resource.NewDecimalQuantity(*dec, resource.DecimalSI),
				},
			},
			wantErr: false,
		},
		{
			name:   "greaterEqual operator - negate",
			op:     []string{">=1"},
			negate: true,
			out: &metricsapi.Operator{
				LessThan: &metricsapi.OperatorValue{
					FixedValue: *resource.NewDecimalQuantity(*dec, resource.DecimalSI),
				},
			},
			wantErr: false,
		},
		{
			name:   "greater operator - negate",
			op:     []string{">1"},
			negate: true,
			out: &metricsapi.Operator{
				LessThanOrEqual: &metricsapi.OperatorValue{
					FixedValue: *resource.NewDecimalQuantity(*dec, resource.DecimalSI),
				},
			},
			wantErr: false,
		},
		{
			name:   "lessEqual operator",
			op:     []string{"<=1"},
			negate: false,
			out: &metricsapi.Operator{
				LessThanOrEqual: &metricsapi.OperatorValue{
					FixedValue: *resource.NewDecimalQuantity(*dec, resource.DecimalSI),
				},
			},
			wantErr: false,
		},
		{
			name:   "less operator",
			op:     []string{"<1"},
			negate: false,
			out: &metricsapi.Operator{
				LessThan: &metricsapi.OperatorValue{
					FixedValue: *resource.NewDecimalQuantity(*dec, resource.DecimalSI),
				},
			},
			wantErr: false,
		},
		{
			name:   "greaterEqual operator",
			op:     []string{">=1"},
			negate: false,
			out: &metricsapi.Operator{
				GreaterThanOrEqual: &metricsapi.OperatorValue{
					FixedValue: *resource.NewDecimalQuantity(*dec, resource.DecimalSI),
				},
			},
			wantErr: false,
		},
		{
			name:   "greater operator",
			op:     []string{">1"},
			negate: false,
			out: &metricsapi.Operator{
				GreaterThan: &metricsapi.OperatorValue{
					FixedValue: *resource.NewDecimalQuantity(*dec, resource.DecimalSI),
				},
			},
			wantErr: false,
		},
		{
			name:   "double operator - negate",
			op:     []string{">1", "<5"},
			negate: true,
			out: &metricsapi.Operator{
				NotInRange: &metricsapi.RangeValue{
					LowBound:  *resource.NewDecimalQuantity(*dec, resource.DecimalSI),
					HighBound: *resource.NewDecimalQuantity(*dec5, resource.DecimalSI),
				},
			},
			wantErr: false,
		},
		{
			name:   "double operator",
			op:     []string{">1", "<5"},
			negate: false,
			out: &metricsapi.Operator{
				InRange: &metricsapi.RangeValue{
					LowBound:  *resource.NewDecimalQuantity(*dec, resource.DecimalSI),
					HighBound: *resource.NewDecimalQuantity(*dec5, resource.DecimalSI),
				},
			},
			wantErr: false,
		},
		{
			name:   "double operator - third one ignored",
			op:     []string{">1", "<5", ">8"},
			negate: false,
			out: &metricsapi.Operator{
				InRange: &metricsapi.RangeValue{
					LowBound:  *resource.NewDecimalQuantity(*dec, resource.DecimalSI),
					HighBound: *resource.NewDecimalQuantity(*dec5, resource.DecimalSI),
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := newOperator(tt.op, tt.negate)
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
			require.Equal(t, tt.want, tt.o.hasNotSupportedCriteria())
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
			name:    "percentage with float - round down",
			val:     "1.333333%",
			want:    1,
			wantErr: false,
		},
		{
			name:    "percentage with float - round up",
			val:     "1.833333%",
			want:    2,
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
			name: "criteria with whitespaces",
			in: []Criteria{
				{
					Operators: []string{"   <  1   0  0   "},
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
	dec5 := inf.NewDec(5, 0)
	dec10 := inf.NewDec(10, 0)
	dec15 := inf.NewDec(15, 0)
	dec20 := inf.NewDec(20, 0)

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
			name: "bogus operator",
			o: &Objective{
				Name: "criteria",
				Pass: []Criteria{
					{
						Operators: []string{"<<<<10"},
					},
				},
			},
			want:    &metricsapi.Target{},
			wantErr: true,
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
			name: "no warn criteria single - conversion",
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
			name: "no warn criteria single - error",
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
			name: "no warn criteria double - conversion",
			o: &Objective{
				Name: "criteria",
				Pass: []Criteria{
					{
						Operators: []string{">10", "<15"},
					},
				},
			},
			want: &metricsapi.Target{
				Failure: &metricsapi.Operator{
					NotInRange: &metricsapi.RangeValue{
						LowBound:  *resource.NewDecimalQuantity(*dec10, resource.DecimalSI),
						HighBound: *resource.NewDecimalQuantity(*dec15, resource.DecimalSI),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "no warn criteria double - error",
			o: &Objective{
				Name: "criteria",
				Pass: []Criteria{
					{
						Operators: []string{"10", "<15"},
					},
				},
			},
			want:    &metricsapi.Target{},
			wantErr: true,
		},
		{
			name: "with warn criteria single pass criteria single - conversion",
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
			name: "with warn criteria double pass criteria double - warn in superset - conversion",
			o: &Objective{
				Name: "criteria",
				Pass: []Criteria{
					{
						Operators: []string{">10", "<15"},
					},
				},
				Warning: []Criteria{
					{
						Operators: []string{"<=20", ">5"},
					},
				},
			},
			want: &metricsapi.Target{
				Failure: &metricsapi.Operator{
					NotInRange: &metricsapi.RangeValue{
						LowBound:  *resource.NewDecimalQuantity(*dec5, resource.DecimalSI),
						HighBound: *resource.NewDecimalQuantity(*dec20, resource.DecimalSI),
					},
				},
				Warning: &metricsapi.Operator{
					NotInRange: &metricsapi.RangeValue{
						LowBound:  *resource.NewDecimalQuantity(*dec10, resource.DecimalSI),
						HighBound: *resource.NewDecimalQuantity(*dec15, resource.DecimalSI),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "with warn criteria single pass criteria double - warn is superset - conversion",
			o: &Objective{
				Name: "criteria",
				Pass: []Criteria{
					{
						Operators: []string{">10", "<15"},
					},
				},
				Warning: []Criteria{
					{
						Operators: []string{"<=20"},
					},
				},
			},
			want: &metricsapi.Target{
				Failure: &metricsapi.Operator{
					GreaterThan: &metricsapi.OperatorValue{
						FixedValue: *resource.NewDecimalQuantity(*dec20, resource.DecimalSI),
					},
				},
				Warning: &metricsapi.Operator{
					NotInRange: &metricsapi.RangeValue{
						LowBound:  *resource.NewDecimalQuantity(*dec10, resource.DecimalSI),
						HighBound: *resource.NewDecimalQuantity(*dec15, resource.DecimalSI),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "with warn criteria double pass criteria single - pass is superset - conversion",
			o: &Objective{
				Name: "criteria",
				Warning: []Criteria{
					{
						Operators: []string{">10", "<15"},
					},
				},
				Pass: []Criteria{
					{
						Operators: []string{"<=20"},
					},
				},
			},
			want: &metricsapi.Target{
				Failure: &metricsapi.Operator{
					GreaterThan: &metricsapi.OperatorValue{
						FixedValue: *resource.NewDecimalQuantity(*dec20, resource.DecimalSI),
					},
				},
				Warning: &metricsapi.Operator{
					InRange: &metricsapi.RangeValue{
						LowBound:  *resource.NewDecimalQuantity(*dec10, resource.DecimalSI),
						HighBound: *resource.NewDecimalQuantity(*dec15, resource.DecimalSI),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "with warn criteria double pass criteria double - pass is superset - conversion",
			o: &Objective{
				Name: "criteria",
				Warning: []Criteria{
					{
						Operators: []string{">10", "<15"},
					},
				},
				Pass: []Criteria{
					{
						Operators: []string{"<=20", ">5"},
					},
				},
			},
			want: &metricsapi.Target{
				Failure: &metricsapi.Operator{
					NotInRange: &metricsapi.RangeValue{
						LowBound:  *resource.NewDecimalQuantity(*dec5, resource.DecimalSI),
						HighBound: *resource.NewDecimalQuantity(*dec20, resource.DecimalSI),
					},
				},
				Warning: &metricsapi.Operator{
					InRange: &metricsapi.RangeValue{
						LowBound:  *resource.NewDecimalQuantity(*dec10, resource.DecimalSI),
						HighBound: *resource.NewDecimalQuantity(*dec15, resource.DecimalSI),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "with warn criteria double pass criteria double - no intersection - conversion",
			o: &Objective{
				Name: "criteria",
				Pass: []Criteria{
					{
						Operators: []string{">15", "<20"},
					},
				},
				Warning: []Criteria{
					{
						Operators: []string{"<=10", ">5"},
					},
				},
			},
			want:    &metricsapi.Target{},
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
	dec10 := inf.NewDec(10, 0)
	dec15 := inf.NewDec(15, 0)

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
			name: "objectives with percentage only - informative",
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
								Operators: []string{"<10%"},
							},
						},
						Warning: []Criteria{
							{
								Operators: []string{"<=15%"},
							},
						},
						KeySLI: true,
						Weight: 10,
					},
					{
						Name: "criteria2",
						Pass: []Criteria{
							{
								Operators: []string{"<10%"},
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
							Target:       metricsapi.Target{},
							Weight:       10,
							KeyObjective: true,
							AnalysisValueTemplateRef: metricsapi.ObjectReference{
								Name:      "criteria",
								Namespace: "default",
							},
						},
						{
							Target: metricsapi.Target{},
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

func TestDecodeOperatorAndValue(t *testing.T) {
	tests := []struct {
		name    string
		o       string
		opOut   string
		opVal   string
		wantErr bool
	}{
		{
			name:    "unsupported operator",
			o:       "--",
			wantErr: true,
		},
		{
			name:    "happy path - less",
			o:       "<5",
			opOut:   "<",
			opVal:   "5",
			wantErr: false,
		},
		{
			name:    "happy path - lessEqual",
			o:       "<=5",
			opOut:   "<=",
			opVal:   "5",
			wantErr: false,
		},
		{
			name:    "happy path - greater",
			o:       ">5",
			opOut:   ">",
			opVal:   "5",
			wantErr: false,
		},
		{
			name:    "happy path - greaterEqual",
			o:       ">=5",
			opOut:   ">=",
			opVal:   "5",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opOut, opVal, err := decodeOperatorAndValue(tt.o)
			if tt.wantErr {
				require.NotNil(t, err)
			} else {
				require.Equal(t, tt.opOut, opOut)
				require.Equal(t, tt.opVal, opVal)
			}
		})

	}
}

func TestDecideIntervalBounds(t *testing.T) {
	dec10 := inf.NewDec(10, 0)
	dec15 := inf.NewDec(15, 0)

	tests := []struct {
		name      string
		op1       string
		val1      string
		op2       string
		val2      string
		smallerOp *Operator
		biggerOp  *Operator
		wantErr   bool
	}{
		{
			name:    "error converting first value",
			op1:     "<",
			val1:    "-",
			op2:     "<",
			val2:    "5",
			wantErr: true,
		},
		{
			name:    "error converting second value",
			op1:     "<",
			val1:    "5",
			op2:     "<",
			val2:    "-",
			wantErr: true,
		},
		{
			name: "fist value smaller",
			op1:  ">",
			val1: "10",
			op2:  "<",
			val2: "15",
			smallerOp: &Operator{
				Value:     dec10,
				Operation: ">",
			},
			biggerOp: &Operator{
				Value:     dec15,
				Operation: "<",
			},
			wantErr: false,
		},
		{
			name: "second value smaller",
			op1:  ">",
			val1: "15",
			op2:  "<",
			val2: "10",
			smallerOp: &Operator{
				Value:     dec10,
				Operation: "<",
			},
			biggerOp: &Operator{
				Value:     dec15,
				Operation: ">",
			},
			wantErr: false,
		},
		{
			name: "equal values",
			op1:  ">",
			val1: "15",
			op2:  "<",
			val2: "15",
			smallerOp: &Operator{
				Value:     dec15,
				Operation: "<",
			},
			biggerOp: &Operator{
				Value:     dec15,
				Operation: ">",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			smallerOperator, biggerOperator, err := decideIntervalBounds(tt.op1, tt.val1, tt.op2, tt.val2)
			if tt.wantErr {
				require.NotNil(t, err)
			} else {
				require.Equal(t, tt.smallerOp, smallerOperator)
				require.Equal(t, tt.biggerOp, biggerOperator)
			}
		})

	}
}

func TestCreateUnboundedInterval(t *testing.T) {
	dec10 := inf.NewDec(10, 0)
	max := inf.NewDec(int64(MaxInt), 0)
	min := inf.NewDec(int64(MinInt), 0)

	tests := []struct {
		name    string
		op      string
		i       *Interval
		wantErr bool
	}{
		{
			name:    "unable to decode operator",
			op:      "--",
			wantErr: true,
		},
		{
			name:    "unable to decode dec number",
			op:      "<--",
			wantErr: true,
		},
		{
			name:    "unsupported operator",
			op:      "=5",
			wantErr: true,
		},
		{
			name:    "inf interval greater",
			op:      ">10",
			wantErr: false,
			i: &Interval{
				Start: dec10,
				End:   max,
			},
		},
		{
			name:    "inf interval greater equal",
			op:      ">=10",
			wantErr: false,
			i: &Interval{
				Start: dec10,
				End:   max,
			},
		},
		{
			name:    "inf interval less",
			op:      "<10",
			wantErr: false,
			i: &Interval{
				Start: min,
				End:   dec10,
			},
		},
		{
			name:    "inf interval less equal",
			op:      "<10",
			wantErr: false,
			i: &Interval{
				Start: min,
				End:   dec10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i, err := createUnboundedInterval(tt.op)
			if tt.wantErr {
				require.NotNil(t, err)
			} else {
				require.Equal(t, tt.i, i)
			}
		})

	}
}

func TestCreateBoundedInterval(t *testing.T) {
	dec10 := inf.NewDec(10, 0)
	dec15 := inf.NewDec(15, 0)

	tests := []struct {
		name    string
		op      []string
		i       *Interval
		wantErr bool
	}{
		{
			name:    "empty array",
			op:      []string{},
			wantErr: true,
		},
		{
			name:    "unable to decode operator1",
			op:      []string{"--", "<5"},
			wantErr: true,
		},
		{
			name:    "unable to decode operator2",
			op:      []string{"<5", "-"},
			wantErr: true,
		},
		{
			name:    "unable to decode inteval bounds",
			op:      []string{"<-", ">5"},
			wantErr: true,
		},
		{
			name:    "unsupported interval",
			op:      []string{"<5", ">10"},
			wantErr: true,
		},
		{
			name:    "happy path",
			op:      []string{">10", "<15"},
			wantErr: false,
			i: &Interval{
				Start: dec10,
				End:   dec15,
			},
		},
		{
			name:    "happy path",
			op:      []string{">=10", "<=15"},
			wantErr: false,
			i: &Interval{
				Start: dec10,
				End:   dec15,
			},
		},
		{
			name:    "happy path",
			op:      []string{">10", "<=15"},
			wantErr: false,
			i: &Interval{
				Start: dec10,
				End:   dec15,
			},
		},
		{
			name:    "happy path",
			op:      []string{">=10", "<15"},
			wantErr: false,
			i: &Interval{
				Start: dec10,
				End:   dec15,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i, err := createBoundedInterval(tt.op)
			if tt.wantErr {
				require.NotNil(t, err)
			} else {
				require.Equal(t, tt.i, i)
			}
		})

	}
}

func TestCreateInterval(t *testing.T) {
	// unbounded interval
	dec10 := inf.NewDec(10, 0)
	dec15 := inf.NewDec(15, 0)
	min := inf.NewDec(int64(MinInt), 0)

	i, err := createInterval([]string{"<10"})
	require.Nil(t, err)
	require.Equal(t, &Interval{
		Start: min,
		End:   dec10,
	}, i)

	// bounded interval
	i, err = createInterval([]string{">10", "<15"})
	require.Nil(t, err)
	require.Equal(t, &Interval{
		Start: dec10,
		End:   dec15,
	}, i)
}

func TestIsSuperInterval(t *testing.T) {
	tests := []struct {
		name    string
		op1     []string
		op2     []string
		want    bool
		wantErr bool
	}{
		{
			name:    "error creating super interval",
			op1:     []string{"--"},
			want:    false,
			wantErr: true,
		},
		{
			name:    "error creating sub interval",
			op1:     []string{"<5"},
			op2:     []string{"--"},
			want:    false,
			wantErr: true,
		},
		{
			name:    "intervals do not intercept",
			op1:     []string{"<5"},
			op2:     []string{">10"},
			want:    false,
			wantErr: false,
		},
		{
			name:    "intervals intercept partially",
			op1:     []string{"<10"},
			op2:     []string{">5"},
			want:    false,
			wantErr: false,
		},
		{
			name:    "subinterval is superinterval",
			op1:     []string{">5", "<7"},
			op2:     []string{"<10"},
			want:    false,
			wantErr: false,
		},
		{
			name:    "equal intervals",
			op1:     []string{"<10"},
			op2:     []string{"<10"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "superinterval unbounded",
			op1:     []string{"<10"},
			op2:     []string{">5", "<7"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "superinterval bounded",
			op1:     []string{">5", "<10"},
			op2:     []string{">5", "<7"},
			want:    true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := isSuperInterval(tt.op1, tt.op2)
			if tt.wantErr {
				require.NotNil(t, err)
			} else {
				require.Equal(t, tt.want, res)
			}
		})

	}
}
