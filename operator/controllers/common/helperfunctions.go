package common

import (
	"context"

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

func RemoveGates(ctx context.Context, c client.Client, workloadInstance *klcv1alpha2.KeptnWorkloadInstance) error {
	switch workloadInstance.Spec.ResourceReference.Kind {
	case "Pod":
		pod := &v1.Pod{}
		err := c.Get(ctx, types.NamespacedName{Namespace: workloadInstance.Namespace, Name: workloadInstance.Spec.ResourceReference.Name}, pod)
		if err != nil {
			return err
		}
		return removePodGates(ctx, c, pod)
	case "ReplicaSet", "StatefulSet", "DaemonSet":
		podList, err := getPodsOfOwner(ctx, c, workloadInstance.Spec.ResourceReference.UID, workloadInstance.Spec.ResourceReference.Kind, workloadInstance.Namespace)
		if err != nil {
			return err
		}
		for _, pod := range podList {
			err := removePodGates(ctx, c, &pod)
			if err != nil {
				return err
			}
		}
	default:
		return controllererrors.ErrUnsupportedWorkloadInstanceResourceReference
	}

	return nil
}

func removePodGates(ctx context.Context, c client.Client, pod *v1.Pod) error {
	pod.Spec.SchedulingGates = nil
	return c.Update(ctx, pod)
}

func getPodsOfOwner(ctx context.Context, c client.Client, ownerUID types.UID, ownerKind string, namespace string) ([]v1.Pod, error) {
	pods := &v1.PodList{}
	err := c.List(ctx, pods, client.InNamespace(namespace))
	if err != nil {
		return nil, err
	}

	var resultPods []v1.Pod

	for _, pod := range pods.Items {
		for _, owner := range pod.OwnerReferences {
			if owner.Kind == ownerKind && owner.UID == ownerUID {
				resultPods = append(resultPods, pod)
				break
			}
		}
	}

	return resultPods, nil
}
