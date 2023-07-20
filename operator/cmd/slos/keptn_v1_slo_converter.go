package slos

import (
	"fmt"
	"github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"regexp"
	"strconv"
	"strings"
)

func TransformSLOToEvaluationDefinition(name, sloConfigStr string) (*v1alpha3.KeptnEvaluationDefinition, error) {
	sloConfig, err := parseSLO([]byte(sloConfigStr))
	if err != nil {
		return nil, err
	}
	kltNamespace := "keptn-lifecycle-toolkit-system"

	evaluationDef := &v1alpha3.KeptnEvaluationDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: kltNamespace,
		},
		Spec: v1alpha3.KeptnEvaluationDefinitionSpec{
			Objectives: []v1alpha3.Objective{},
			TotalScore: &v1alpha3.PassThreshold{
				PassPercentage:    toFloat(sloConfig.TotalScore.Pass),
				WarningPercentage: toFloat(sloConfig.TotalScore.Warning),
			},
			Comparison: &v1alpha3.ComparisonSpec{
				CompareWith:         sloConfig.Comparison.CompareWith,
				IncludeWarning:      false,
				NumberOfComparisons: sloConfig.Comparison.NumberOfComparisonResults,
				AggregationFunction: sloConfig.Comparison.AggregateFunction,
			},
		},
	}

	if sloConfig.Comparison.CompareWith == "pass_or_warn" {
		evaluationDef.Spec.Comparison.IncludeWarning = true
	}

	for _, objective := range sloConfig.Objectives {

		obj := v1alpha3.Objective{
			KeptnMetricRef: v1alpha3.KeptnMetricReference{
				Name:      objective.SLI,
				Namespace: kltNamespace,
			},
			SLOTargets: &v1alpha3.SLOTarget{
				Pass: v1alpha3.OrCombinedCriteriaSet{
					AnyOf: []v1alpha3.Criteria{},
				},
				Warning: v1alpha3.OrCombinedCriteriaSet{
					AnyOf: []v1alpha3.Criteria{},
				},
			},
			Weight:       objective.Weight,
			KeyObjective: objective.KeySLI,
		}

		for _, passTarget := range objective.Pass {

			kltPassCriteria := v1alpha3.Criteria{
				Targets: []v1alpha3.Target{},
			}

			for _, criteriaStr := range passTarget.Criteria {
				kltTarget, err := parseCriteriaString(criteriaStr)
				if err != nil {
					// continue with the other criteria
					continue
				}

				kltPassCriteria.Targets = append(kltPassCriteria.Targets, *kltTarget)
			}

			obj.SLOTargets.Pass.AnyOf = append(obj.SLOTargets.Pass.AnyOf, kltPassCriteria)

		}

		evaluationDef.Spec.Objectives = append(evaluationDef.Spec.Objectives, obj)
	}

	return evaluationDef, nil
}

func parseCriteriaString(criteria string) (*v1alpha3.Target, error) {
	// example values: <+15%, <500, >-8%, =0
	// possible operators: <, <=, =, >, >=
	// regex: ^([<|<=|=|>|>=]{1,2})([+|-]{0,1}\\d*\.?\d*)([%]{0,1})
	regex := `^([<|<=|=|>|>=]{1,2})([+|-]{0,1}\d*\.?\d*)([%]{0,1})`
	var re *regexp.Regexp
	re = regexp.MustCompile(regex)

	// remove whitespaces
	criteria = strings.Replace(criteria, " ", "", -1)

	if !re.MatchString(criteria) {
		return nil, errors.New("invalid criteria string")
	}

	kltTarget := v1alpha3.Target{
		LessThanOrEqual: nil,
	}

	targetValue := v1alpha3.TargetValue{}

	operators := []string{"<=", "<", "=", ">=", ">"}
	operatorToUse := ""

	for _, operator := range operators {
		if strings.HasPrefix(criteria, operator) {
			operatorToUse = operator
			criteria = strings.TrimPrefix(criteria, operator)
			break
		}
	}

	if strings.HasSuffix(criteria, "%") {
		targetValue.Comparison = &v1alpha3.ComparisonTarget{}
		criteria = strings.TrimSuffix(criteria, "%")

		if strings.HasPrefix(criteria, "-") {
			criteria = strings.TrimPrefix(criteria, "-")
			floatValue := toFloat(criteria)
			targetValue.Comparison.DecreaseByPercent = &floatValue
		} else if strings.HasPrefix(criteria, "+") {
			criteria = strings.TrimPrefix(criteria, "+")
			criteria = strings.TrimPrefix(criteria, "-")
			floatValue := toFloat(criteria)
			targetValue.Comparison.IncreaseByPercent = &floatValue
		}
	} else {
		floatValue := toFloat(criteria)
		targetValue.FixedValue = &floatValue
	}

	switch operatorToUse {
	case "<=":
		kltTarget.LessThanOrEqual = &targetValue
	case "<":
		kltTarget.LessThan = &targetValue
	case "=":
		kltTarget.EqualTo = &targetValue
	case ">":
		kltTarget.GreaterThan = &targetValue
	case ">=":
		kltTarget.GreaterThanOrEqual = &targetValue
	}
	return &kltTarget, nil
}

func toFloat(strVal string) float32 {
	val, err := strconv.ParseFloat(strings.TrimSuffix(strVal, "%"), 32)
	if err != nil {
		return 0.0
	}
	return float32(val)
}

func parseSLO(input []byte) (*ServiceLevelObjectives, error) {
	slo := &ServiceLevelObjectives{}
	err := yaml.Unmarshal(input, &slo)

	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	if slo.Comparison == nil {
		slo.Comparison = &SLOComparison{
			CompareWith:               "single_result",
			IncludeResultWithScore:    "all",
			NumberOfComparisonResults: 1,
			AggregateFunction:         "avg",
		}
	}

	if slo.Comparison != nil {
		if slo.Comparison.IncludeResultWithScore == "" {
			slo.Comparison.IncludeResultWithScore = "all"
		}
		if slo.Comparison.NumberOfComparisonResults == 0 {
			slo.Comparison.NumberOfComparisonResults = 3
		}
		if slo.Comparison.AggregateFunction == "" {
			slo.Comparison.AggregateFunction = "avg"
		}
	}

	objectives := []*SLO{}
	for _, objective := range slo.Objectives {
		if objective == nil {
			continue
		}
		if objective.Weight == 0 {
			objective.Weight = 1
		}
		objectives = append(objectives, objective)
	}
	slo.Objectives = objectives

	return slo, nil
}
