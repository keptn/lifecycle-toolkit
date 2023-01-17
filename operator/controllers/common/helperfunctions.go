package common

import (
	"context"

	"github.com/go-logr/logr"
	klcv1alpha2 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2/common"
	controllererrors "github.com/keptn/lifecycle-toolkit/operator/controllers/errors"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type CreateAttributes struct {
	SpanName   string
	Definition string
	CheckType  apicommon.CheckType
}

// GetItemStatus retrieves the state of the task/evaluation, if it does not exists, it creates a default one
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

// GetOldStatus retrieves the state of the task/evaluation
func GetOldStatus(name string, statuses []klcv1alpha2.ItemStatus) apicommon.KeptnState {
	var oldstatus apicommon.KeptnState
	for _, ts := range statuses {
		if ts.DefinitionName == name {
			oldstatus = ts.Status
		}
	}

	return oldstatus
}

func RemoveGates(ctx context.Context, c client.Client, log logr.Logger, workloadInstance *klcv1alpha2.KeptnWorkloadInstance) error {
	switch workloadInstance.Spec.ResourceReference.Kind {
	case "Pod":
		return removePodGates(ctx, c, log, workloadInstance.Spec.ResourceReference.Name, workloadInstance.Namespace)
	case "ReplicaSet", "StatefulSet", "DaemonSet":
		podList, err := getPodsOfOwner(ctx, c, log, workloadInstance.Spec.ResourceReference.UID, workloadInstance.Spec.ResourceReference.Kind, workloadInstance.Namespace)
		if err != nil {
			log.Error(err, "cannot get pods")
			return err
		}
		for _, pod := range podList {
			err := removePodGates(ctx, c, log, pod, workloadInstance.Namespace)
			if err != nil {
				log.Error(err, "cannot remove gates from pod")
				return err
			}
		}
	default:
		return controllererrors.ErrUnsupportedWorkloadInstanceResourceReference
	}

	return nil
}

func removePodGates(ctx context.Context, c client.Client, log logr.Logger, podName string, podNamespace string) error {
	pod := &v1.Pod{}
	err := c.Get(ctx, types.NamespacedName{Namespace: podNamespace, Name: podName}, pod)
	if err != nil {
		log.Error(err, "cannot remove gates from pod - inner")
		return err
	}
	if len(pod.Annotations) == 0 {
		pod.Annotations = make(map[string]string)
	}
	pod.Annotations[apicommon.SchedullingGateRemoved] = "true"
	pod.Spec.SchedulingGates = nil
	return c.Update(ctx, pod)
}

func getPodsOfOwner(ctx context.Context, c client.Client, log logr.Logger, ownerUID types.UID, ownerKind string, namespace string) ([]string, error) {
	pods := &v1.PodList{}
	err := c.List(ctx, pods, client.InNamespace(namespace))
	if err != nil {
		log.Error(err, "cannot list pods - inner")
		return nil, err
	}

	var resultPods []string

	for _, pod := range pods.Items {
		for _, owner := range pod.OwnerReferences {
			if owner.Kind == ownerKind && owner.UID == ownerUID {
				resultPods = append(resultPods, pod.Name)
				break
			}
		}
	}

	return resultPods, nil
}
