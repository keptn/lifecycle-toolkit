package keptnworkloadinstance

import (
	"context"
	klcv1alpha1 "github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1"
	"github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1/common"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
		}
	}

	isReplicaRunning, err := r.isReplicaSetRunning(ctx, workloadInstance.Spec.ResourceReference, workloadInstance.Namespace)
	if err != nil {
		return common.StateUnknown, err
	}
	if isReplicaRunning {
		workloadInstance.Status.DeploymentStatus = common.StateSucceeded
	}

	err = r.Client.Status().Update(ctx, workloadInstance)
	if err != nil {
		return common.StateUnknown, err
	}
	return common.StateSucceeded, nil
}

func (r *KeptnWorkloadInstanceReconciler) isReplicaSetRunning(ctx context.Context, resource klcv1alpha1.ResourceReference, namespace string) (bool, error) {
	replica := &appsv1.ReplicaSetList{}
	if err := r.Client.List(ctx, replica, client.InNamespace(namespace)); err != nil {
		return false, err
	}
	for _, re := range replica.Items {
		if re.UID == resource.UID {
			replicas, err := r.getDesiredReplicas(ctx, re.OwnerReferences[0], namespace)
			if err != nil {
				return false, err
			}
			if re.Status.ReadyReplicas == replicas {
				return true, nil
			}
			return false, nil
		}
	}
	return false, nil

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

func (r *KeptnWorkloadInstanceReconciler) getDesiredReplicas(ctx context.Context, reference v1.OwnerReference, namespace string) (int32, error) {
	var replicas *int32
	switch reference.Kind {
	case "Deployment":
		dep := appsv1.Deployment{}
		err := r.Client.Get(ctx, types.NamespacedName{Name: reference.Name, Namespace: namespace}, &dep)
		if err != nil {
			return 0, err
		}
		replicas = dep.Spec.Replicas
	case "StatefulSet":
		sts := appsv1.StatefulSet{}
		err := r.Client.Get(ctx, types.NamespacedName{Name: reference.Name, Namespace: namespace}, &sts)
		if err != nil {
			return 0, err
		}
		replicas = sts.Spec.Replicas
	}

	return *replicas, nil

}
