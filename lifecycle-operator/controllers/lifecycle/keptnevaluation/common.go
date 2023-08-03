package keptnevaluation

import (
	"fmt"
	"math"
	"strconv"

	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3"
)

func checkValue(objective klcv1alpha3.Objective, item *klcv1alpha3.EvaluationStatusItem) (bool, error) {

	if len(item.Value) == 0 || len(objective.EvaluationTarget) == 0 {
		return false, fmt.Errorf("no values")
	}

	eval := objective.EvaluationTarget[1:]
	sign := objective.EvaluationTarget[:1]

	resultValue, err := strconv.ParseFloat(item.Value, 64)
	if err != nil || math.IsNaN(resultValue) {
		return false, err
	}

	compareValue, err := strconv.ParseFloat(eval, 64)
	if err != nil || math.IsNaN(compareValue) {
		return false, err
	}

	// choose comparator
	switch sign {
	case ">":
		return resultValue > compareValue, nil
	case "<":
		return resultValue < compareValue, nil
	default:
		return false, fmt.Errorf("invalid operator")
	}
}
