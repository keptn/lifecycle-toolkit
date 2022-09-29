package keptnworkloadinstance

import (
	"context"

	klcv1alpha1 "github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1"
	"github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1/common"
	ctrl "sigs.k8s.io/controller-runtime"
)

type StatusSummary struct {
	failed    int
	succeeded int
	running   int
	pending   int
}

func (r *KeptnWorkloadInstanceReconciler) reconcilePreDeployment(ctx context.Context, req ctrl.Request, workloadInstance *klcv1alpha1.KeptnWorkloadInstance) error {
	newStatus, preDeploymentState, err := r.genericPrePost(ctx, common.PreDeploymentCheckType, workloadInstance)
	if err != nil {
		return err
	}
	workloadInstance.Status.PreDeploymentStatus = getOverallState(preDeploymentState)
	workloadInstance.Status.PreDeploymentTaskStatus = newStatus

	// Write Status Field
	err = r.Client.Status().Update(ctx, workloadInstance)
	if err != nil {
		return err
	}
	return nil
}
