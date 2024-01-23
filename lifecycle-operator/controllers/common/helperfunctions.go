package common

import (
	"context"
	"fmt"
	"maps"

	"github.com/go-logr/logr"
	klcv1beta1 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1beta1"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1beta1/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/config"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GetItemStatus retrieves the state of the task/evaluation, if it does not exists, it creates a default one
func GetItemStatus(name string, instanceStatus []klcv1beta1.ItemStatus) klcv1beta1.ItemStatus {
	for _, status := range instanceStatus {
		if status.DefinitionName == name {
			return status
		}
	}
	return klcv1beta1.ItemStatus{
		DefinitionName: name,
		Status:         apicommon.StatePending,
		Name:           "",
	}
}

// GetOldStatus retrieves the state of the task/evaluation
func GetOldStatus(name string, statuses []klcv1beta1.ItemStatus) apicommon.KeptnState {
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
	// we copy the map2 first, so the values in the overlapping
	// properties are set from map1 in the resulting map
	maps.Copy[M1](merged, map2)
	maps.Copy[M1](merged, map1)
	return merged
}

func GetTaskDefinition(k8sclient client.Client, log logr.Logger, ctx context.Context, definitionName string, namespace string) (*klcv1beta1.KeptnTaskDefinition, error) {
	definition := &klcv1beta1.KeptnTaskDefinition{}
	if err := getObject(k8sclient, log, ctx, definitionName, namespace, definition); err != nil {
		return nil, err
	}
	return definition, nil
}

func GetEvaluationDefinition(k8sclient client.Client, log logr.Logger, ctx context.Context, definitionName string, namespace string) (*klcv1beta1.KeptnEvaluationDefinition, error) {
	definition := &klcv1beta1.KeptnEvaluationDefinition{}
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
