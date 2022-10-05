package keptnworkloadinstance

import (
	"context"

	klcv1alpha1 "github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1"
	"github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1/common"
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
