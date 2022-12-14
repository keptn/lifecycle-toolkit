package keptnevaluation

import (
	"fmt"
	"math"
	"strconv"

	klcv1alpha2 "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha2"
)

func checkValue(objective klcv1alpha2.Objective, item *klcv1alpha2.EvaluationStatusItem) (bool, error) {

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
