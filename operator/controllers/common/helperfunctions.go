package common

import (
	"fmt"

	klcv1alpha2 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2/common"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
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

func RecordEvent(recorder record.EventRecorder, phase apicommon.KeptnPhaseType, eventType string, reconcileObject client.Object, shortReason string, longReason string, version string) { //
	recorder.Event(reconcileObject, eventType, fmt.Sprintf("%s %s", phase.ShortName, shortReason), fmt.Sprintf("%s %s / Namespace: %s, Name: %s, Version: %s ", phase.LongName, longReason, reconcileObject.GetNamespace(), reconcileObject.GetName(), version))
}
