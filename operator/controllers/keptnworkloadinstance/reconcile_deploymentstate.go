package keptnworkloadinstance

import (
	"context"
	klcv1alpha1 "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1"
	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1/common"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *KeptnWorkloadInstanceReconciler) reconcileDeployment(ctx context.Context, workloadInstance *klcv1alpha1.KeptnWorkloadInstance) (common.KeptnState, error) {
	if workloadInstance.Spec.ResourceReference.Kind == "Pod" {
		isPodRunning, err := r.isPodRunning(ctx, workloadInstance.Spec.ResourceReference, workloadInstance.Namespace)
		if err != nil {
			return common.StateUnknown, err
		}
		if isPodRunning {
			workloadInstance.Status.DeploymentStatus = common.StateSucceeded
		} else {
			workloadInstance.Status.DeploymentStatus = common.StateProgressing
		}
	} else {
		isReplicaRunning, err := r.isReferencedWorkloadRunning(ctx, workloadInstance.Spec.ResourceReference, workloadInstance.Namespace)
		if err != nil {
			return common.StateUnknown, err
		}
		if isReplicaRunning {
			workloadInstance.Status.DeploymentStatus = common.StateSucceeded
		} else {
			workloadInstance.Status.DeploymentStatus = common.StateProgressing
		}
	}

	err := r.Client.Status().Update(ctx, workloadInstance)
	if err != nil {
		return common.StateUnknown, err
	}
	return workloadInstance.Status.DeploymentStatus, nil
}

func (r *KeptnWorkloadInstanceReconciler) isReferencedWorkloadRunning(ctx context.Context, resource klcv1alpha1.ResourceReference, namespace string) (bool, error) {

	var replicas *int32
	var desiredReplicas int32
	switch resource.Kind {
	case "ReplicaSet":
		rep := appsv1.ReplicaSet{}
		err := r.Client.Get(ctx, types.NamespacedName{Name: resource.Name, Namespace: namespace}, &rep)
		if err != nil {
			return false, err
		}
		replicas = rep.Spec.Replicas
		desiredReplicas = rep.Status.AvailableReplicas
	case "StatefulSet":
		sts := appsv1.StatefulSet{}
		err := r.Client.Get(ctx, types.NamespacedName{Name: resource.Name, Namespace: namespace}, &sts)
		if err != nil {
			return false, err
		}
		replicas = sts.Spec.Replicas
		desiredReplicas = sts.Status.AvailableReplicas
	}

	return *replicas == desiredReplicas, nil

}

func (r *KeptnWorkloadInstanceReconciler) isPodRunning(ctx context.Context, resource klcv1alpha1.ResourceReference, namespace string) (bool, error) {
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
