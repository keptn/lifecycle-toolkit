package common

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/config"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/lifecycle/interfaces"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

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
		log.Info("Could not find resource in application namespace", "resource type", fmt.Sprintf("%T", definition), "Definition name", definitionName, "namespace", namespace)
		if k8serrors.IsNotFound(err) {
			if err := k8sclient.Get(ctx, types.NamespacedName{Name: definitionName, Namespace: config.Instance().GetDefaultNamespace()}, definition); err != nil {
				log.Info("Could not find resource in default Keptn namespace", "resource type", fmt.Sprintf("%T", definition), "definition name", definitionName)
				return err
			}
			return nil
		}
		return err
	}
	return nil
}

// GetRequestInfo extracts name and namespace from a controller request.
func GetRequestInfo(req ctrl.Request) map[string]string {
	return map[string]string{
		"name":      req.Name,
		"namespace": req.Namespace,
	}
}
