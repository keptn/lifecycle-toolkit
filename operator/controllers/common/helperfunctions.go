package common

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3/common"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/lifecycle/interfaces"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const KLTNamespace = "keptn-lifecycle-toolkit-system"

// GetItemStatus retrieves the state of the task/evaluation, if it does not exists, it creates a default one
func GetItemStatus(name string, instanceStatus []klcv1alpha3.ItemStatus) klcv1alpha3.ItemStatus {
	for _, status := range instanceStatus {
		if status.DefinitionName == name {
			return status
		}
	}
	return klcv1alpha3.ItemStatus{
		DefinitionName: name,
		Status:         apicommon.StatePending,
		Name:           "",
	}
}

func GetAppVersionName(namespace string, appName string, version string) types.NamespacedName {
	return types.NamespacedName{Namespace: namespace, Name: appName + "-" + version}
}

// GetOldStatus retrieves the state of the task/evaluation
func GetOldStatus(name string, statuses []klcv1alpha3.ItemStatus) apicommon.KeptnState {
	var oldstatus apicommon.KeptnState
	for _, ts := range statuses {
		if ts.DefinitionName == name {
			oldstatus = ts.Status
		}
	}

	return oldstatus
}

// RecordEvent creates k8s Event and adds it to Eventqueue
func RecordEvent(recorder record.EventRecorder, phase apicommon.KeptnPhaseType, eventType string, reconcileObject client.Object, shortReason string, longReason string, version string) {
	msg := setEventMessage(phase, reconcileObject, longReason, version)
	annotations := setAnnotations(reconcileObject, phase)
	recorder.AnnotatedEventf(reconcileObject, annotations, eventType, fmt.Sprintf("%s%s", phase.ShortName, shortReason), msg)
}

func setEventMessage(phase apicommon.KeptnPhaseType, reconcileObject client.Object, longReason string, version string) string {
	if version == "" {
		return fmt.Sprintf("%s: %s / Namespace: %s, Name: %s", phase.LongName, longReason, reconcileObject.GetNamespace(), reconcileObject.GetName())
	}
	return fmt.Sprintf("%s: %s / Namespace: %s, Name: %s, Version: %s", phase.LongName, longReason, reconcileObject.GetNamespace(), reconcileObject.GetName(), version)
}

func setAnnotations(reconcileObject client.Object, phase apicommon.KeptnPhaseType) map[string]string {
	if reconcileObject == nil || reconcileObject.GetName() == "" || reconcileObject.GetNamespace() == "" {
		return nil
	}
	annotations := map[string]string{
		"namespace": reconcileObject.GetNamespace(),
		"name":      reconcileObject.GetName(),
		"phase":     phase.ShortName,
	}

	piWrapper, err := interfaces.NewEventObjectWrapperFromClientObject(reconcileObject)
	if err == nil {
		copyMap(annotations, piWrapper.GetEventAnnotations())
	}

	annotationsObject := reconcileObject.GetAnnotations()
	annotations["traceparent"] = annotationsObject["traceparent"]

	return annotations
}

func copyMap[M1 ~map[K]V, M2 ~map[K]V, K comparable, V any](dst M1, src M2) {
	for k, v := range src {
		dst[k] = v
	}
}

func GetTaskDefinition(k8sclient client.Client, log logr.Logger, ctx context.Context, definitionName string, namespace string) (*klcv1alpha3.KeptnTaskDefinition, error) {
	definition := &klcv1alpha3.KeptnTaskDefinition{}
	if err := getObject(k8sclient, log, ctx, definitionName, namespace, definition); err != nil {
		return nil, err
	}
	return definition, nil
}

func GetEvaluationDefinition(k8sclient client.Client, log logr.Logger, ctx context.Context, definitionName string, namespace string) (*klcv1alpha3.KeptnEvaluationDefinition, error) {
	definition := &klcv1alpha3.KeptnEvaluationDefinition{}
	if err := getObject(k8sclient, log, ctx, definitionName, namespace, definition); err != nil {
		return nil, err
	}
	return definition, nil
}

func getObject(k8sclient client.Client, log logr.Logger, ctx context.Context, definitionName string, namespace string, definition client.Object) error {
	err := k8sclient.Get(ctx, types.NamespacedName{Name: definitionName, Namespace: namespace}, definition)
	if err != nil {
		log.Info("Failed to get resource from application namespace", "resource type", fmt.Sprintf("%T", definition), "Definition name", definitionName, "namespace", namespace)
		if k8serrors.IsNotFound(err) {
			if err := k8sclient.Get(ctx, types.NamespacedName{Name: definitionName, Namespace: KLTNamespace}, definition); err != nil {
				log.Info("Failed to get KeptnEvaluationDefinition from default KLT namespace", "resource type", fmt.Sprintf("%T", definition), "definition name", definitionName)
				return err
			}
			return nil
		}
		return err
	}
	return nil
}
