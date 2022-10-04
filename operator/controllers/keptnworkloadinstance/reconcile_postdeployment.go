package keptnworkloadinstance

import (
	"context"
	klcv1alpha1 "github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1"
	"github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1/common"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *KeptnWorkloadInstanceReconciler) reconcilePostDeployment(ctx context.Context, req ctrl.Request, workloadInstance *klcv1alpha1.KeptnWorkloadInstance) error {

	newStatus, postDeploymentState, err := r.reconcileChecks(ctx, common.PostDeploymentCheckType, workloadInstance)
	if err != nil {
		return err
	}
	workloadInstance.Status.PostDeploymentStatus = getOverallState(postDeploymentState)
	workloadInstance.Status.PostDeploymentTaskStatus = newStatus

	// Write Status Field
	err = r.Client.Status().Update(ctx, workloadInstance)
	if err != nil {
		return err
	}
	return nil
}

func (r *KeptnWorkloadInstanceReconciler) GetDesiredReplicas(ctx context.Context, reference v1.OwnerReference, namespace string) (int32, error) {
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
