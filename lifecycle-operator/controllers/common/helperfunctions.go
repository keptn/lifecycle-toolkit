package common

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	apilifecycle "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/config"
	"golang.org/x/exp/maps"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GetItemStatus retrieves the state of the task/evaluation, if it does not exists, it creates a default one
func GetItemStatus(name string, instanceStatus []apilifecycle.ItemStatus) apilifecycle.ItemStatus {
	for _, status := range instanceStatus {
		if status.DefinitionName == name {
			return status
		}
	}
	return apilifecycle.ItemStatus{
		DefinitionName: name,
		Status:         apicommon.StatePending,
		Name:           "",
	}
}

// GetOldStatus retrieves the state of the task/evaluation
func GetOldStatus(name string, statuses []apilifecycle.ItemStatus) apicommon.KeptnState {
	var oldstatus apicommon.KeptnState
	for _, ts := range statuses {
		if ts.DefinitionName == name {
			oldstatus = ts.Status
		}
	}

	return oldstatus
}

func MergeMaps[M1 ~map[K]V, K comparable, V any](map1 M1, map2 M1) M1 {
	merged := make(M1, len(map1)+len(map2))
	// we copy the map1 first, so the values in the overlapping
	// properties are set from map2 in the resulting map
	maps.Copy[M1](merged, map1)
	maps.Copy[M1](merged, map2)
	return merged
}

func GetTaskDefinition(k8sclient client.Client, log logr.Logger, ctx context.Context, definitionName string, namespace string) (*apilifecycle.KeptnTaskDefinition, error) {
	definition := &apilifecycle.KeptnTaskDefinition{}
	if err := getObject(k8sclient, log, ctx, definitionName, namespace, definition); err != nil {
		return nil, err
	}
	return definition, nil
}

func GetEvaluationDefinition(k8sclient client.Client, log logr.Logger, ctx context.Context, definitionName string, namespace string) (*apilifecycle.KeptnEvaluationDefinition, error) {
	definition := &apilifecycle.KeptnEvaluationDefinition{}
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

func KeptnWorkloadVersionResourceRefUIDIndexFunc(rawObj client.Object) []string {
	// Extract the ResourceReference UID name from the KeptnWorkloadVersion Spec, if one is provided
	workloadVersion, ok := rawObj.(*apilifecycle.KeptnWorkloadVersion)
	if !ok {
		return nil
	}
	if workloadVersion.Spec.ResourceReference.UID == "" {
		return nil
	}
	return []string{string(workloadVersion.Spec.ResourceReference.UID)}
}
