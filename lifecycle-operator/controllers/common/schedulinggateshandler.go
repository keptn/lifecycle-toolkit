package common

import (
	"context"

	"github.com/go-logr/logr"
	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	controllererrors "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/errors"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/schedulinggateshandler_mock.go . ISchedulingGatesHandler
type ISchedulingGatesHandler interface {
	RemoveGates(ctx context.Context, workloadVersion *klcv1alpha3.KeptnWorkloadVersion) error
	Enabled() bool
}

type RemoveGatesFunc func(ctx context.Context, c client.Client, podName string, podNamespace string) error
type GetPodsFunc func(ctx context.Context, c client.Client, ownerUID types.UID, ownerKind string, namespace string) ([]string, error)

type SchedulingGatesHandler struct {
	client.Client
	logr.Logger
	enabled     bool
	removeGates RemoveGatesFunc
	getPods     GetPodsFunc
}

func NewSchedulingGatesHandler(c client.Client, l logr.Logger, enabled bool) *SchedulingGatesHandler {
	return &SchedulingGatesHandler{
		Client:      c,
		Logger:      l,
		enabled:     enabled,
		removeGates: removePodGates,
		getPods:     getPodsOfOwner,
	}
}

func (h *SchedulingGatesHandler) RemoveGates(ctx context.Context, workloadVersion *klcv1alpha3.KeptnWorkloadVersion) error {
	switch workloadVersion.Spec.ResourceReference.Kind {
	case "Pod":
		return h.removeGates(ctx, h.Client, workloadVersion.Spec.ResourceReference.Name, workloadVersion.Namespace)
	case "ReplicaSet", "StatefulSet", "DaemonSet":
		podList, err := h.getPods(ctx, h.Client, workloadVersion.Spec.ResourceReference.UID, workloadVersion.Spec.ResourceReference.Kind, workloadVersion.Namespace)
		if err != nil {
			h.Logger.Error(err, "cannot get pods")
			return err
		}
		for _, pod := range podList {
			err := h.removeGates(ctx, h.Client, pod, workloadVersion.Namespace)
			if err != nil {
				h.Logger.Error(err, "cannot remove gates from pod")
				return err
			}
		}
	default:
		return controllererrors.ErrUnsupportedWorkloadVersionResourceReference
	}

	return nil
}

func (h *SchedulingGatesHandler) Enabled() bool {
	return h.enabled
}

func removePodGates(ctx context.Context, c client.Client, podName string, podNamespace string) error {
	pod := &v1.Pod{}
	err := c.Get(ctx, types.NamespacedName{Namespace: podNamespace, Name: podName}, pod)
	if err != nil {
		return err
	}

	if pod.Annotations[apicommon.SchedulingGateRemoved] != "" {
		return nil
	}

	if len(pod.Annotations) == 0 {
		pod.Annotations = make(map[string]string, 1)
	}
	pod.Annotations[apicommon.SchedulingGateRemoved] = "true"
	pod.Spec.SchedulingGates = nil
	return c.Update(ctx, pod)
}

func getPodsOfOwner(ctx context.Context, c client.Client, ownerUID types.UID, ownerKind string, namespace string) ([]string, error) {
	pods := &v1.PodList{}
	err := c.List(ctx, pods, client.InNamespace(namespace))
	if err != nil {
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
