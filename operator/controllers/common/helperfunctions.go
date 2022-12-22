package common

import (
	klcv1alpha2 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2/common"
	"k8s.io/apimachinery/pkg/types"
)

func GetTaskStatus(taskName string, instanceStatus []klcv1alpha2.TaskStatus) klcv1alpha2.TaskStatus {
	for _, status := range instanceStatus {
		if status.TaskDefinitionName == taskName {
			return status
		}
	}
	return klcv1alpha2.TaskStatus{
		TaskDefinitionName: taskName,
		Status:             apicommon.StatePending,
		TaskName:           "",
	}
}

func GetEvaluationStatus(evaluationName string, instanceStatus []klcv1alpha2.EvaluationStatus) klcv1alpha2.EvaluationStatus {
	for _, status := range instanceStatus {
		if status.EvaluationDefinitionName == evaluationName {
			return status
		}
	}
	return klcv1alpha2.EvaluationStatus{
		EvaluationDefinitionName: evaluationName,
		Status:                   apicommon.StatePending,
		EvaluationName:           "",
	}
}

func GetAppVersionName(namespace string, appName string, version string) types.NamespacedName {
	return types.NamespacedName{Namespace: namespace, Name: appName + "-" + version}
}
