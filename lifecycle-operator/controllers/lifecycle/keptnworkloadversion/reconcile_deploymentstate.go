package keptnworkloadversion

import (
	"context"

	argov1alpha1 "github.com/argoproj/argo-rollouts/pkg/apis/rollouts/v1alpha1"
	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	klcv1alpha4 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha4"
	controllererrors "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *KeptnWorkloadVersionReconciler) reconcileDeployment(ctx context.Context, workloadVersion *klcv1alpha4.KeptnWorkloadVersion) (apicommon.KeptnState, error) {
	var isRunning bool
	var err error

	switch workloadVersion.Spec.ResourceReference.Kind {
	case "Pod":
		isRunning, err = r.isPodRunning(ctx, workloadVersion.Spec.ResourceReference, workloadVersion.Namespace)
	case "ReplicaSet":
		isRunning, err = r.isReplicaSetRunning(ctx, workloadVersion.Spec.ResourceReference, workloadVersion.Namespace)
	case "StatefulSet":
		isRunning, err = r.isStatefulSetRunning(ctx, workloadVersion.Spec.ResourceReference, workloadVersion.Namespace)
	case "DaemonSet":
		isRunning, err = r.isDaemonSetRunning(ctx, workloadVersion.Spec.ResourceReference, workloadVersion.Namespace)
	default:
		isRunning, err = false, controllererrors.ErrUnsupportedWorkloadVersionResourceReference
	}

	if err != nil {
		return apicommon.StateUnknown, err
	}
	if isRunning {
		workloadVersion.Status.DeploymentStatus = apicommon.StateSucceeded
	} else {
		workloadVersion.Status.DeploymentStatus = apicommon.StateProgressing
	}

	err = r.Client.Status().Update(ctx, workloadVersion)
	if err != nil {
		return apicommon.StateUnknown, err
	}
	return workloadVersion.Status.DeploymentStatus, nil
}

func (r *KeptnWorkloadVersionReconciler) isReplicaSetRunning(ctx context.Context, resource klcv1alpha3.ResourceReference, namespace string) (bool, error) {
	rep := appsv1.ReplicaSet{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: resource.Name, Namespace: namespace}, &rep)
	if err != nil {
		return false, err
	}

	for _, ownerRef := range rep.OwnerReferences {
		if ownerRef.Kind == "Rollout" {
			return r.isRolloutRunning(ctx, klcv1alpha3.ResourceReference{Name: ownerRef.Name, UID: ownerRef.UID}, namespace)
		}
	}

	return *rep.Spec.Replicas == rep.Status.AvailableReplicas, nil
}

func (r *KeptnWorkloadVersionReconciler) isDaemonSetRunning(ctx context.Context, resource klcv1alpha3.ResourceReference, namespace string) (bool, error) {
	daemonSet := &appsv1.DaemonSet{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: resource.Name, Namespace: namespace}, daemonSet)
	if err != nil {
		return false, err
	}
	return daemonSet.Status.DesiredNumberScheduled == daemonSet.Status.NumberReady, nil
}

func (r *KeptnWorkloadVersionReconciler) isPodRunning(ctx context.Context, resource klcv1alpha3.ResourceReference, namespace string) (bool, error) {
	podList := &corev1.PodList{}
	if err := r.Client.List(ctx, podList, client.InNamespace(namespace)); err != nil {
		return false, err
	}
	for _, p := range podList.Items {
		if p.UID == resource.UID {
			if p.Status.Phase == corev1.PodRunning {
				return true, nil
			}
			return false, nil
		}
	}
	return false, nil
}

func (r *KeptnWorkloadVersionReconciler) isStatefulSetRunning(ctx context.Context, resource klcv1alpha3.ResourceReference, namespace string) (bool, error) {
	sts := appsv1.StatefulSet{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: resource.Name, Namespace: namespace}, &sts)
	if err != nil {
		return false, err
	}
	return *sts.Spec.Replicas == sts.Status.AvailableReplicas, nil
}

func (r *KeptnWorkloadVersionReconciler) isRolloutRunning(ctx context.Context, resource klcv1alpha3.ResourceReference, namespace string) (bool, error) {
	rollout := argov1alpha1.Rollout{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: resource.Name, Namespace: namespace}, &rollout)
	if err != nil {
		return false, err
	}
	return rollout.Status.Replicas == rollout.Status.UpdatedReplicas && rollout.Status.Phase == argov1alpha1.RolloutPhaseHealthy, nil
}
