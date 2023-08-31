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

const InvalidOperatorErrMsg = "invalid operator: '%s'"
const UnableConvertValueErrMsg = "unable to convert value '%s' to decimal"

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

type Operator struct {
	Value    *inf.Dec
	Operator string
}

func (o *Objective) hasNotSupportedCriteria() bool {
	return len(o.Pass) > 1 || len(o.Warning) > 1
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
	// skip objective target conversion if it has criteria combined with logical OR -> not supported
	// this way the SLO will become "informative"
	// it will become informative as well when the pass criteria are not defined
	if o.hasNotSupportedCriteria() || len(o.Pass) == 0 {
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

	// if pass is superinterval of warn, the following logic is used:
	// !(pass criteria) -> fail criteria
	// warn criteria -> warn criteria
	if isSuperInterval(o.Pass[0].Operators, o.Warning[0].Operators) {
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

	// if warning is superinterval of pass OR we have a single rule criteria, the following logic is used:
	// !(warn criteria) -> fail criteria
	// !(pass criteria) -> warn criteria
	// TODO change if statement when isSuporInterval is implemented
	if (len(o.Pass[0].Operators) == 1 && len(o.Warning[0].Operators) == 1) || true /*isSuperInterval(o.Warning[0].Operators, o.Pass[0].Operators) */ {
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
	return target, nil
}

// TODO implement
func isSuperInterval(op1 []string, op2 []string) bool {
	// superOp1, superOpVal1, err := decodeOperatorAndValue(op1[0])
	// if err != nil {
	// 	return false, err
	// }
	// subOp1, subOpVal1, err := decodeOperatorAndValue(op2[0])
	// if err != nil {
	// 	return false, err
	// }
	return false
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

	return "", "", fmt.Errorf(InvalidOperatorErrMsg, op)
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

	return &metricsapi.Operator{}, fmt.Errorf("empty operators: '%v'", op)
}

// checks and negates the existing single operator
//
//nolint:dupl
func negateSingleOperator(op string, value string) (*metricsapi.Operator, error) {
	dec := inf.NewDec(1, 0)
	_, ok := dec.SetString(value)
	if !ok {
		return nil, fmt.Errorf(UnableConvertValueErrMsg, value)
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

	return &metricsapi.Operator{}, fmt.Errorf(InvalidOperatorErrMsg, op)
}

// checks and creates single operator
//
//nolint:dupl
func createSingleOperator(op string, value string) (*metricsapi.Operator, error) {
	dec := inf.NewDec(1, 0)
	_, ok := dec.SetString(value)
	if !ok {
		return nil, fmt.Errorf(UnableConvertValueErrMsg, value)
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

	return &metricsapi.Operator{}, fmt.Errorf(InvalidOperatorErrMsg, op)
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
	if (smallerOperator.Operator == ">" || smallerOperator.Operator == ">=") && (biggerOperator.Operator == "<" || biggerOperator.Operator == "<=") {
		return &metricsapi.Operator{
			InRange: r,
		}, nil
		// outOfRange interval
	} else if (smallerOperator.Operator == "<" || smallerOperator.Operator == "<=") && (biggerOperator.Operator == ">" || biggerOperator.Operator == ">=") {
		return &metricsapi.Operator{
			NotInRange: r,
		}, nil
	}

	return nil, fmt.Errorf("unconvertable combination of operators: '%s', '%s'", op1, op2)
}

// decides which of the values is smaller and binds operator to them
func decideIntervalBounds(op1 string, value1 string, op2 string, value2 string) (*Operator, *Operator, error) {
	dec1 := inf.NewDec(1, 0)
	_, ok := dec1.SetString(value1)
	if !ok {
		return nil, nil, fmt.Errorf(UnableConvertValueErrMsg, value1)
	}
	dec2 := inf.NewDec(1, 0)
	_, ok = dec2.SetString(value2)
	if !ok {
		return nil, nil, fmt.Errorf(UnableConvertValueErrMsg, value2)
	}

	operator1 := &Operator{
		Value:    dec1,
		Operator: op1,
	}

	operator2 := &Operator{
		Value:    dec2,
		Operator: op2,
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
