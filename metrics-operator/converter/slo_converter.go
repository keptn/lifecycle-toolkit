package converter

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	"gopkg.in/inf.v0"
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

type SLOConverter struct {
}

func NewSLOConverter() *SLOConverter {
	return &SLOConverter{}
}

type SLO struct {
	Objectives []*Objective `yaml:"objectives" json:"objectives"`
	TotalScore Score        `yaml:"total_score" json:"total_score"`
}

type Score struct {
	Pass    string `yaml:"pass" json:"pass"`
	Warning string `yaml:"warning" json:"warning"`
}

type Objective struct {
	Name    string     `yaml:"sli" json:"sli"`
	KeySLI  bool       `yaml:"key_sli,omitempty" json:"key_sli,omitempty"`
	Weight  int        `yaml:"weight,omitempty" json:"weight,omitempty"`
	Warning []Criteria `yaml:"warning,omitempty" json:"warning,omitempty"`
	Pass    []Criteria `yaml:"pass,omitempty" json:"pass,omitempty"`
}

type Criteria struct {
	Operators []string `yaml:"criteria,omitempty" json:"criteria,omitempty"`
}

func (o *Objective) hasNotSupportedCriteria() bool {
	// no pass criteria -> informative
	if len(o.Pass) == 0 {
		return true
	}
	// support only warning criteria with a single criteria element
	if len(o.Warning) > 1 {
		return true
	}
	// warning criteria == 1, pass can be only 1
	if len(o.Warning) == 1 {
		return len(o.Pass) > 1
	}

	// warn criteria == 0, pass can be anything
	return false
}

func (c *SLOConverter) Convert(fileContent []byte, analysisDef string, namespace string) (string, error) {
	//check that provider and namespace is set
	if analysisDef == "" || namespace == "" {
		return "", fmt.Errorf("missing arguments: 'definition' and 'namespace' needs to be set for conversion")
	}

	// unmarshall content
	content := &SLO{}
	err := yaml.Unmarshal(fileContent, content)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling file content: %s", err.Error())
	}

	// convert
	analysisDefinition, err := c.convertSLO(content, analysisDef, namespace)
	if err != nil {
		return "", err
	}

	// marshal AnalysisDefinition to Yaml
	yamlData, err := yaml.Marshal(analysisDefinition)
	if err != nil {
		return "", fmt.Errorf("error marshalling data: %s", err.Error())
	}

	return string(yamlData), nil
}

func (c *SLOConverter) convertSLO(sloContent *SLO, name string, namespace string) (*metricsapi.AnalysisDefinition, error) {
	// define resulting AnalysisDefinition with easy conversions
	passPercentage, err := removePercentage(sloContent.TotalScore.Pass)
	if err != nil {
		return nil, err
	}
	warnPercentage, err := removePercentage(sloContent.TotalScore.Warning)
	if err != nil {
		return nil, err
	}
	definition := &metricsapi.AnalysisDefinition{
		TypeMeta: v1.TypeMeta{
			Kind:       "AnalysisDefinition",
			APIVersion: "metrics.keptn.sh/v1alpha3",
		},
		ObjectMeta: v1.ObjectMeta{
			Name: name,
		},
		Spec: metricsapi.AnalysisDefinitionSpec{
			TotalScore: metricsapi.TotalScore{
				PassPercentage:    passPercentage,
				WarningPercentage: warnPercentage,
			},
			// create a slice of size of len(sloContent.Objectives), but reserve capacity for
			// double the size, as some objectives may be twice there (conversion of criteria with logical AND)
			Objectives: make([]metricsapi.Objective, len(sloContent.Objectives)),
		},
	}

	// convert objectives one after another
	indexObjectives := 0
	for _, o := range sloContent.Objectives {
		// remove criteria, which contain % in their operators
		o = cleanupObjective(o)
		// set up target
		target, err := setupTarget(o)
		if err != nil {
			return nil, err
		}
		objective := metricsapi.Objective{
			AnalysisValueTemplateRef: metricsapi.ObjectReference{
				Name:      o.Name,
				Namespace: namespace,
			},
			KeyObjective: o.KeySLI,
			Weight:       o.Weight,
			Target:       *target,
		}
		definition.Spec.Objectives[indexObjectives] = objective
		indexObjectives++
	}
	return definition, nil
}

// removes % symbol from the scoring values and converts to numeric value
func removePercentage(str string) (int, error) {
	t := strings.ReplaceAll(str, "%", "")
	f, err := strconv.ParseFloat(t, 64)
	if err != nil {
		return 0, err
	}
	return int(math.Round(f)), nil
}

// creates and sets up the target struct from objective
// nolint:gocognit,gocyclo
func setupTarget(o *Objective) (*metricsapi.Target, error) {
	target := &metricsapi.Target{}
	// skip unsupported combination of criteria and informative objectives
	if o.hasNotSupportedCriteria() {
		return target, nil
	}

	// multiple criteria combined with logical OR operator
	if len(o.Pass) > 1 {
		ops := []string{o.Pass[0].Operators[0], o.Pass[1].Operators[0]}
		op, err := newOperator(ops, true)
		if err != nil {
			return target, err
		}
		target.Failure = op
		return target, nil
	}

	// if warning criteria are not defined, negate the pass criteria to create fail criteria
	if len(o.Warning) == 0 {
		if len(o.Pass) > 0 {
			if len(o.Pass[0].Operators) > 0 {
				op, err := newOperator(o.Pass[0].Operators, true)
				if err != nil {
					return target, err
				}
				target.Failure = op
				return target, nil
			}
		}
	}

	// if warning is superinterval of pass OR we have a single rule criteria, the following logic is used:
	// !(warn criteria) -> fail criteria
	// !(pass criteria) -> warn criteria
	isWarningSuperInterval, err := isSuperInterval(o.Warning[0].Operators, o.Pass[0].Operators)
	if err != nil {
		return target, err
	}
	if (len(o.Pass[0].Operators) == 1 && len(o.Warning[0].Operators) == 1) || isWarningSuperInterval {
		op1, err := newOperator(o.Pass[0].Operators, true)
		if err != nil {
			return target, err
		}
		op2, err := newOperator(o.Warning[0].Operators, true)
		if err != nil {
			return target, err
		}
		target.Failure = op2
		target.Warning = op1
		return target, nil
	}

	// if pass is superinterval of warn, the following logic is used:
	// !(pass criteria) -> fail criteria
	// warn criteria -> warn criteria
	isPassSuperInterval, err := isSuperInterval(o.Pass[0].Operators, o.Warning[0].Operators)
	if err != nil {
		return target, err
	}
	if isPassSuperInterval {
		op1, err := newOperator(o.Warning[0].Operators, false)
		if err != nil {
			return target, err
		}
		op2, err := newOperator(o.Pass[0].Operators, true)
		if err != nil {
			return target, err
		}
		target.Failure = op2
		target.Warning = op1
		return target, nil
	}

	return target, nil
}

// checks if interval is valid and if the first set of operators defines interval
// which is superset of the interval defined by second set of operators
func isSuperInterval(op1 []string, op2 []string) (bool, error) {
	superInterval, err := createInterval(op1)
	if err != nil {
		return false, err
	}
	subInterval, err := createInterval(op2)
	if err != nil {
		return false, err
	}

	return superInterval.Start.Cmp(subInterval.Start) < 1 && superInterval.End.Cmp(subInterval.End) >= 0, nil
}

// creates interval from set of operators
func createInterval(op []string) (*Interval, error) {
	// if it's unbounded interval, we have only one operator
	if len(op) == 1 {
		return createUnboundedInterval(op[0])
	}

	//bounded interval
	return createBoundedInterval(op)
}

func createBoundedInterval(op []string) (*Interval, error) {
	if len(op) < 2 {
		return nil, NewUnsupportedIntervalCombinationErr(op)
	}
	//fetch operators and values
	operator1, value1, err := decodeOperatorAndValue(op[0])
	if err != nil {
		return nil, err
	}
	operator2, value2, err := decodeOperatorAndValue(op[1])
	if err != nil {
		return nil, err
	}
	// determine lower and higher bouds
	smallerOperator, biggerOperator, err := decideIntervalBounds(operator1, value1, operator2, value2)
	if err != nil {
		return nil, err
	}
	//check if the interval makes logical sense for conversions, e.g. 5 < x < 10; unsupported: x < 5 && x > 10
	if isGreaterOrEqual(smallerOperator.Operation) && isLessOrEqual(biggerOperator.Operation) {
		return &Interval{
			Start: smallerOperator.Value,
			End:   biggerOperator.Value,
		}, nil
	}

	return nil, NewUnsupportedIntervalCombinationErr(op)
}

func createUnboundedInterval(op string) (*Interval, error) {
	//fetch operator and value
	operator, value, err := decodeOperatorAndValue(op)
	if err != nil {
		return nil, err
	}
	dec := inf.NewDec(1, 0)
	_, ok := dec.SetString(value)
	if !ok {
		return nil, NewUnconvertableValueErr(value)
	}
	// interval of (val, Inf)
	if isGreaterOrEqual(operator) {
		return &Interval{
			Start: dec,
			End:   inf.NewDec(int64(MaxInt), 0),
		}, nil
		// interval of (-Inf, val)
	}

	return &Interval{
		Start: inf.NewDec(int64(MinInt), 0),
		End:   dec,
	}, nil
}

func cleanupObjective(o *Objective) *Objective {
	o.Pass = cleanupCriteria(o.Pass)
	o.Warning = cleanupCriteria(o.Warning)
	return o
}

// remove % operators from criterium structure
// if criteria did have only % operators, remove it from strucutre
func cleanupCriteria(criteria []Criteria) []Criteria {
	newCriteria := make([]Criteria, 0, len(criteria))
	for _, c := range criteria {
		operators := make([]string, 0, len(c.Operators))
		for _, op := range c.Operators {
			// keep only criteria with real values, not percentage
			if !strings.Contains(op, "%") {
				// remove unneeded whitespaces from criteria string
				operators = append(operators, strings.Replace(op, " ", "", -1))
			}
		}
		// if criterium does have operator, store it
		if len(operators) > 0 {
			newCriteria = append(newCriteria, Criteria{Operators: operators})
		}
	}

	return newCriteria
}

// check if operator is valid and split it to operator and value
func decodeOperatorAndValue(op string) (string, string, error) {
	operators := []string{"<=", "<", ">=", ">"}
	for _, operator := range operators {
		if strings.HasPrefix(op, operator) {
			return operator, strings.TrimPrefix(op, operator), nil
		}
	}

	return "", "", NewInvalidOperatorErr(op)
}

// create operator for Target
func newOperator(op []string, negate bool) (*metricsapi.Operator, error) {
	// convert single operator
	if len(op) == 1 {
		operator, value, err := decodeOperatorAndValue(op[0])
		if err != nil {
			return nil, err
		}
		if negate {
			return negateSingleOperator(operator, value)
		} else {
			return createSingleOperator(operator, value)
		}
	} else if len(op) >= 2 { // convert operators representing range
		operator1, value1, err := decodeOperatorAndValue(op[0])
		if err != nil {
			return nil, err
		}
		operator2, value2, err := decodeOperatorAndValue(op[1])
		if err != nil {
			return nil, err
		}
		if negate {
			return negateDoubleOperator(operator1, value1, operator2, value2)
		} else {
			return createDoubleOperator(operator1, value1, operator2, value2)
		}
	}

	return &metricsapi.Operator{}, NewEmptyOperatorErr(op)
}

// checks and negates the existing single operator
//
//nolint:dupl
func negateSingleOperator(op string, value string) (*metricsapi.Operator, error) {
	dec := inf.NewDec(1, 0)
	_, ok := dec.SetString(value)
	if !ok {
		return nil, NewUnconvertableValueErr(value)
	}
	if op == "<=" {
		return &metricsapi.Operator{
			GreaterThan: &metricsapi.OperatorValue{
				FixedValue: *resource.NewDecimalQuantity(*dec, resource.DecimalSI),
			},
		}, nil
	} else if op == "<" {
		return &metricsapi.Operator{
			GreaterThanOrEqual: &metricsapi.OperatorValue{
				FixedValue: *resource.NewDecimalQuantity(*dec, resource.DecimalSI),
			},
		}, nil
	} else if op == ">=" {
		return &metricsapi.Operator{
			LessThan: &metricsapi.OperatorValue{
				FixedValue: *resource.NewDecimalQuantity(*dec, resource.DecimalSI),
			},
		}, nil
	} else if op == ">" {
		return &metricsapi.Operator{
			LessThanOrEqual: &metricsapi.OperatorValue{
				FixedValue: *resource.NewDecimalQuantity(*dec, resource.DecimalSI),
			},
		}, nil
	}

	return &metricsapi.Operator{}, NewInvalidOperatorErr(op)
}

// checks and creates single operator
//
//nolint:dupl
func createSingleOperator(op string, value string) (*metricsapi.Operator, error) {
	dec := inf.NewDec(1, 0)
	_, ok := dec.SetString(value)
	if !ok {
		return nil, NewUnconvertableValueErr(value)
	}
	if op == "<=" {
		return &metricsapi.Operator{
			LessThanOrEqual: &metricsapi.OperatorValue{
				FixedValue: *resource.NewDecimalQuantity(*dec, resource.DecimalSI),
			},
		}, nil
	} else if op == "<" {
		return &metricsapi.Operator{
			LessThan: &metricsapi.OperatorValue{
				FixedValue: *resource.NewDecimalQuantity(*dec, resource.DecimalSI),
			},
		}, nil
	} else if op == ">=" {
		return &metricsapi.Operator{
			GreaterThanOrEqual: &metricsapi.OperatorValue{
				FixedValue: *resource.NewDecimalQuantity(*dec, resource.DecimalSI),
			},
		}, nil
	} else if op == ">" {
		return &metricsapi.Operator{
			GreaterThan: &metricsapi.OperatorValue{
				FixedValue: *resource.NewDecimalQuantity(*dec, resource.DecimalSI),
			},
		}, nil
	}

	return &metricsapi.Operator{}, NewInvalidOperatorErr(op)
}

// checks and creates double operator
func createDoubleOperator(op1 string, value1 string, op2 string, value2 string) (*metricsapi.Operator, error) {
	smallerOperator, biggerOperator, err := decideIntervalBounds(op1, value1, op2, value2)
	if err != nil {
		return nil, err
	}

	// create range
	r := &metricsapi.RangeValue{
		LowBound:  *resource.NewDecimalQuantity(*smallerOperator.Value, resource.DecimalSI),
		HighBound: *resource.NewDecimalQuantity(*biggerOperator.Value, resource.DecimalSI),
	}

	// inRange interval
	if isGreaterOrEqual(smallerOperator.Operation) && isLessOrEqual(biggerOperator.Operation) {
		return &metricsapi.Operator{
			InRange: r,
		}, nil
		// outOfRange interval
	} else if isLessOrEqual(smallerOperator.Operation) && isGreaterOrEqual(biggerOperator.Operation) {
		return &metricsapi.Operator{
			NotInRange: r,
		}, nil
	}

	return nil, NewUnconvertableOperatorCombinationErr(op1, op2)
}

// decides which of the values is smaller and binds operator to them
func decideIntervalBounds(op1 string, value1 string, op2 string, value2 string) (*Operator, *Operator, error) {
	dec1 := inf.NewDec(1, 0)
	_, ok := dec1.SetString(value1)
	if !ok {
		return nil, nil, NewUnconvertableValueErr(value1)
	}
	dec2 := inf.NewDec(1, 0)
	_, ok = dec2.SetString(value2)
	if !ok {
		return nil, nil, NewUnconvertableValueErr(value2)
	}

	operator1 := &Operator{
		Value:     dec1,
		Operation: op1,
	}

	operator2 := &Operator{
		Value:     dec2,
		Operation: op2,
	}

	if dec1.Cmp(dec2) == -1 {
		return operator1, operator2, nil
	}

	return operator2, operator1, nil
}

// checks and negates double operator
func negateDoubleOperator(op1 string, value1 string, op2 string, value2 string) (*metricsapi.Operator, error) {
	// create range operator
	operator, err := createDoubleOperator(op1, value1, op2, value2)
	if err != nil {
		return operator, err
	}

	// negate it
	if operator.NotInRange != nil {
		return &metricsapi.Operator{
			InRange: operator.NotInRange,
		}, nil
	}

	return &metricsapi.Operator{
		NotInRange: operator.InRange,
	}, nil
}
