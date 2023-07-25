package slos

import (
	"fmt"
	"github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3/common"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"regexp"
	"strconv"
	"strings"
)

type Service struct {
	Name       string
	SLIContent string
	SLOContent string
}

type Stage struct {
	Name       string
	Services   []Service
	SLIContent string
	SLOContent string
}

type Project struct {
	Name       string
	SLIContent string
	SLOContent string
	Stages     []Stage
}

type SLIConversion struct {
	Project  Project
	Provider string
}

type EvaluationDefinitionResult struct {
	EvaluationDefinition *v1alpha3.KeptnEvaluationDefinition
	Error                error
}

type MetricResult struct {
	Metric *unstructured.Unstructured
	Error  error
}

type ConversionResult struct {
	EvaluationDefinition EvaluationDefinitionResult
	Metrics              []MetricResult

	Name string
}

type ServiceConversionResult struct {
	Name string
	ConversionResult
}

type StageConversionResult struct {
	ConversionResult
	Name     string
	Services []ServiceConversionResult
}

type ProjectConversionResult struct {
	ConversionResult
	Name   string
	Stages []StageConversionResult
}

type SLIConversionResult struct {
	Name    string
	Project ProjectConversionResult
}

func TransformKeptnProject(conversion SLIConversion) SLIConversionResult {
	result := SLIConversionResult{
		Name: conversion.Project.Name,
		Project: ProjectConversionResult{
			Name:   conversion.Project.Name,
			Stages: []StageConversionResult{},
		},
	}

	if conversion.Project.SLOContent != "" {
		projectSLOName := common.TruncateString(fmt.Sprintf("%s-slo", conversion.Project.Name), 253)
		projectEvaluationDefinition, err := TransformSLOToEvaluationDefinition(projectSLOName, conversion.Project.SLOContent)

		result.Project.EvaluationDefinition = EvaluationDefinitionResult{
			EvaluationDefinition: projectEvaluationDefinition,
			Error:                err,
		}
	}

	for _, stage := range conversion.Project.Stages {

		stageResult := StageConversionResult{
			Name:     stage.Name,
			Services: []ServiceConversionResult{},
		}
		if stage.SLOContent != "" {
			stageSLOName := common.TruncateString(fmt.Sprintf("%s-%s-slo", conversion.Project, stage.Name), 253)
			stageEvaluationDefinition, err := TransformSLOToEvaluationDefinition(stageSLOName, stage.SLOContent)
			stageResult.EvaluationDefinition = EvaluationDefinitionResult{
				EvaluationDefinition: stageEvaluationDefinition,
				Error:                err,
			}
		}

		for _, service := range stage.Services {
			serviceResult := ServiceConversionResult{
				Name: service.Name,
			}

			if service.SLOContent != "" {
				serviceSLOName := common.TruncateString(fmt.Sprintf("%s-%s-%s-slo", conversion.Project, stage.Name, service.Name), 253)
				serviceEvaluationDefinition, err := TransformSLOToEvaluationDefinition(serviceSLOName, service.SLOContent)
				serviceResult.EvaluationDefinition = EvaluationDefinitionResult{
					EvaluationDefinition: serviceEvaluationDefinition,
					Error:                err,
				}
			}

			stageResult.Services = append(stageResult.Services, serviceResult)
		}

		result.Project.Stages = append(result.Project.Stages, stageResult)
	}
	return result
}

func TransformSLIToMetrics(conversion SLIConversion, name, query string) (*unstructured.Unstructured, error) {
	// we do not have placeholder support in SLIs, so here we would need to create a KeptnMetric for each combination of project/stage/service
	return nil, nil
}

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
				Pass: v1alpha3.CriteriaSet{
					AnyOf: []v1alpha3.Criteria{},
				},
				Warning: v1alpha3.CriteriaSet{
					AnyOf: []v1alpha3.Criteria{},
				},
			},
			Weight:       objective.Weight,
			KeyObjective: objective.KeySLI,
		}

		applyTargets(objective.Pass, &obj.SLOTargets.Pass)
		applyTargets(objective.Warning, &obj.SLOTargets.Warning)

		evaluationDef.Spec.Objectives = append(evaluationDef.Spec.Objectives, obj)
	}

	return evaluationDef, nil
}

func applyTargets(objectives []*SLOCriteria, dst *v1alpha3.CriteriaSet) {
	for _, passTarget := range objectives {
		kltPassCriteria := v1alpha3.Criteria{
			AllOf: []v1alpha3.Target{},
		}

		for _, criteriaStr := range passTarget.Criteria {
			kltTarget, err := parseCriteriaString(criteriaStr)
			if err != nil {
				// continue with the other criteria
				continue
			}

			kltPassCriteria.AllOf = append(kltPassCriteria.AllOf, *kltTarget)
		}

		dst.AnyOf = append(dst.AnyOf, kltPassCriteria)
	}
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
