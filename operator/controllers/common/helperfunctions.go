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

func RecordEvent(recorder record.EventRecorder, phase apicommon.KeptnPhaseType, eventType string, reconcileObject client.Object, shortReason string, longReason string, version string) {
	msg := fmt.Sprintf("%s %s / Namespace: %s, Name: %s, Version: %s ", phase.LongName, longReason, reconcileObject.GetNamespace(), reconcileObject.GetName(), version)
	if version == "" {
		msg = fmt.Sprintf("%s %s / Namespace: %s, Name: %s", phase.LongName, longReason, reconcileObject.GetNamespace(), reconcileObject.GetName())
	}

	annotations := map[string]string{
		"namespace": reconcileObject.GetNamespace(),
		"name":      reconcileObject.GetName(),
		"phase":     phase.ShortName,
	}
	if app, ok := reconcileObject.(*klcv1alpha2.KeptnApp); ok {
		annotations["appName"] = app.Name
		annotations["appVersion"] = app.Spec.Version
	} else if appVersion, ok := reconcileObject.(*klcv1alpha2.KeptnAppVersion); ok {
		annotations["appName"] = appVersion.Spec.AppName
		annotations["appVersion"] = app.Spec.Version
		annotations["appVersionName"] = appVersion.Name
	} else if workload, ok := reconcileObject.(*klcv1alpha2.KeptnWorkload); ok {
		annotations["appName"] = workload.Spec.AppName
		annotations["workloadName"] = workload.Spec.AppName
		annotations["workloadVersion"] = workload.Spec.Version
	} else if workloadInstance, ok := reconcileObject.(*klcv1alpha2.KeptnWorkloadInstance); ok {
		annotations["appName"] = workloadInstance.Spec.AppName
		annotations["workloadName"] = workloadInstance.Spec.WorkloadName
		annotations["workloadVersion"] = workloadInstance.Spec.Version
		annotations["workloadInstanceName"] = workloadInstance.Name
	} else if task, ok := reconcileObject.(*klcv1alpha2.KeptnTask); ok {
		annotations["appName"] = task.Spec.AppName
		annotations["appVersion"] = task.Spec.AppVersion
		annotations["workloadName"] = task.Spec.Workload
		annotations["workloadVersion"] = task.Spec.WorkloadVersion
		annotations["taskName"] = task.Name
		annotations["taskDefinitionName"] = task.Spec.TaskDefinition
	} else if evaluation, ok := reconcileObject.(*klcv1alpha2.KeptnEvaluation); ok {
		annotations["appName"] = evaluation.Spec.AppName
		annotations["appVersion"] = evaluation.Spec.AppVersion
		annotations["workloadName"] = evaluation.Spec.Workload
		annotations["workloadVersion"] = evaluation.Spec.WorkloadVersion
		annotations["evaluationName"] = evaluation.Name
		annotations["evaluationDefinitionName"] = evaluation.Spec.EvaluationDefinition
	}

	annotationsObject := reconcileObject.GetAnnotations()
	if val, ok := annotationsObject["traceparent"]; ok {
		annotations["traceId"] = val
	}

	recorder.AnnotatedEventf(reconcileObject, annotations, eventType, fmt.Sprintf("%s %s", phase.ShortName, shortReason), msg)
}
