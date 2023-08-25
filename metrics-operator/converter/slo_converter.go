package converter

import (
	"fmt"
	"strconv"
	"strings"

	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
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
	Objectives []Objective `yaml:"objectives" json:"objectives"`
	TotalScore Score       `yaml:"total_score" json:"total_score"`
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

func (c *SLOConverter) Convert(fileContent []byte, analysisDef string, namespace string) (string, error) {
	//check that provider and namespace is set
	if analysisDef == "" || namespace == "" {
		return "", fmt.Errorf("--definition and --slo-namespace needs to be set for conversion")
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
				PassPercentage:    removePercentage(sloContent.TotalScore.Pass),
				WarningPercentage: removePercentage(sloContent.TotalScore.Warning),
			},
			Objectives: make([]metricsapi.Objective, len(sloContent.Objectives), len(sloContent.Objectives)*2),
		},
	}

	for i, o := range sloContent.Objectives {
		objective := metricsapi.Objective{
			AnalysisValueTemplateRef: metricsapi.ObjectReference{
				Name:      o.Name,
				Namespace: namespace,
			},
			KeyObjective: o.KeySLI,
			Weight:       o.Weight,
			Target:       setupTarget(o),
		}
		definition.Spec.Objectives[i] = objective
	}
	return definition, nil
}

func removePercentage(str string) int {
	t := strings.ReplaceAll(str, "%", "")
	y, _ := strconv.Atoi(t)
	return y
}

func setupTarget(o Objective) metricsapi.Target {
	target := metricsapi.Target{}
	o = cleanupObjective(o)
	if shouldIgnoreObjective(o) {
		return target
	}

	if len(o.Warning) == 0 {
		if len(o.Pass) > 0 {
			if len(o.Pass[0].Operators) > 0 {
				op, _ := setupOperator(o.Pass[0].Operators[0])
				target.Failure = op
				return target
			}
		}
	}

	if len(o.Pass) > 0 {
		if len(o.Pass[0].Operators) > 0 {
			op, _ := setupOperator(o.Pass[0].Operators[0])
			target.Warning = op
		}
		if len(o.Warning[0].Operators) > 0 {
			op, _ := setupOperator(o.Warning[0].Operators[0])
			target.Failure = op
		}
	}

	return target
}

func cleanupObjective(o Objective) Objective {
	o.Pass = cleanupCriteria(o.Pass)
	o.Warning = cleanupCriteria(o.Warning)
	return o
}

func cleanupCriteria(criteria []Criteria) []Criteria {
	newCriteria := make([]Criteria, 0, len(criteria))
	for _, c := range criteria {
		operators := make([]string, 0, len(c.Operators))
		for _, op := range c.Operators {
			// keep only criteria with real values, not percentage
			if !strings.Contains(op, "%") {
				operators = append(operators, op)
			}
		}
		// if criterium does have operator, store it
		if len(operators) > 0 {
			newCriteria = append(newCriteria, Criteria{Operators: operators})
		}
	}

	return newCriteria
}

func shouldIgnoreObjective(o Objective) bool {
	return len(o.Pass) > 1 || len(o.Warning) > 1
}

func setupOperator(op string) (*metricsapi.Operator, error) {
	// remove whitespaces
	op = strings.Replace(op, " ", "", -1)

	operators := []string{"<=", "<", ">=", ">"}
	for _, operator := range operators {
		if strings.HasPrefix(op, operator) {
			return createOperator(operator, strings.TrimPrefix(op, operator))
		}
	}

	return &metricsapi.Operator{}, nil
}

func createOperator(o string, value string) (*metricsapi.Operator, error) {
	v, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return nil, err
	}
	if o == "<=" {
		return &metricsapi.Operator{
			GreaterThan: &metricsapi.OperatorValue{
				FixedValue: *resource.NewQuantity(v, resource.DecimalSI),
			},
		}, nil
	} else if o == "<" {
		return &metricsapi.Operator{
			GreaterThanOrEqual: &metricsapi.OperatorValue{
				FixedValue: *resource.NewQuantity(v, resource.DecimalSI),
			},
		}, nil
	} else if o == ">=" {
		return &metricsapi.Operator{
			LessThan: &metricsapi.OperatorValue{
				FixedValue: *resource.NewQuantity(v, resource.DecimalSI),
			},
		}, nil
	} else if o == ">" {
		return &metricsapi.Operator{
			LessThanOrEqual: &metricsapi.OperatorValue{
				FixedValue: *resource.NewQuantity(v, resource.DecimalSI),
			},
		}, nil
	}

	return &metricsapi.Operator{}, fmt.Errorf("invalid operator")
}
