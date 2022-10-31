package common

import (
	klcv1alpha1 "github.com/keptn/lifecycle-controller/operator/api/v1alpha1"
	apicommon "github.com/keptn/lifecycle-controller/operator/api/v1alpha1/common"
)

func GetTaskStatus(taskName string, instanceStatus []klcv1alpha1.TaskStatus) klcv1alpha1.TaskStatus {
	for _, status := range instanceStatus {
		if status.TaskDefinitionName == taskName {
			return status
		}
	}
	return klcv1alpha1.TaskStatus{
		TaskDefinitionName: taskName,
		Status:             apicommon.StatePending,
		TaskName:           "",
	}
}

func GetEvaluationStatus(evaluationName string, instanceStatus []klcv1alpha1.EvaluationStatus) klcv1alpha1.EvaluationStatus {
	for _, status := range instanceStatus {
		if status.EvaluationDefinitionName == evaluationName {
			return status
		}
	}
	return klcv1alpha1.EvaluationStatus{
		EvaluationDefinitionName: evaluationName,
		Status:                   apicommon.StatePending,
		EvaluationName:           "",
	}
}
