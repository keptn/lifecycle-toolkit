package common

import (
	klcv1alpha2 "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha2"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha2/common"
	"k8s.io/apimachinery/pkg/types"
)

type CreateAttributes struct {
	SpanName   string
	Definition string
	CheckType  apicommon.CheckType
}

func GetTaskStatus(taskName string, instanceStatus []klcv1alpha2.ItemStatus) klcv1alpha2.ItemStatus {
	for _, status := range instanceStatus {
		if status.DefinitionName == taskName {
			return status
		}
	}
	return klcv1alpha2.ItemStatus{
		DefinitionName: taskName,
		Status:         apicommon.StatePending,
		Name:           "",
	}
}

func GetEvaluationStatus(evaluationName string, instanceStatus []klcv1alpha2.ItemStatus) klcv1alpha2.ItemStatus {
	for _, status := range instanceStatus {
		if status.DefinitionName == evaluationName {
			return status
		}
	}
	return klcv1alpha2.ItemStatus{
		DefinitionName: evaluationName,
		Status:         apicommon.StatePending,
		Name:           "",
	}
}

func GetAppVersionName(namespace string, appName string, version string) types.NamespacedName {
	return types.NamespacedName{Namespace: namespace, Name: appName + "-" + version}
}

func GetOldStatus(statuses []klcv1alpha2.ItemStatus, definitionName string) apicommon.KeptnState {
	var oldstatus apicommon.KeptnState
	for _, ts := range statuses {
		if ts.DefinitionName == definitionName {
			oldstatus = ts.Status
		}
	}

	return oldstatus
}
