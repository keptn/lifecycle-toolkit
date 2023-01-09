package common

import (
	klcv1alpha2 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2/common"
	"k8s.io/apimachinery/pkg/types"
)

func GetItemStatus(name string, instanceStatus []klcv1alpha2.ItemStatus) klcv1alpha2.ItemStatus {
	for _, status := range instanceStatus {
		if status.DefinitionName == name {
			return status
		}
	}
	return klcv1alpha2.ItemStatus{
		DefinitionName: name,
		Status:         apicommon.StatePending,
		Name:           "",
	}
}

func GetAppVersionName(namespace string, appName string, version string) types.NamespacedName {
	return types.NamespacedName{Namespace: namespace, Name: appName + "-" + version}
}
