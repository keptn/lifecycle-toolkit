package keptnworkloadinstance

import (
	"context"

	klcv1alpha1 "github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1"
	"github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1/common"
)

func (r *KeptnWorkloadInstanceReconciler) reconcilePreDeployment(ctx context.Context, workloadInstance *klcv1alpha1.KeptnWorkloadInstance) (common.KeptnState, error) {
	newStatus, preDeploymentState, err := r.reconcileChecks(ctx, common.PreDeploymentCheckType, workloadInstance)
	if err != nil {
		return common.StateUnknown, err
	}
	overallState := common.GetOverallState(preDeploymentState)
	workloadInstance.Status.PreDeploymentStatus = overallState
	workloadInstance.Status.PreDeploymentTaskStatus = newStatus

	// Write Status Field
	err = r.Client.Status().Update(ctx, workloadInstance)
	if err != nil {
		return common.StateUnknown, err
	}
	return overallState, nil
}
