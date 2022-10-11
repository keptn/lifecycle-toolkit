package keptnworkloadinstance

import (
	"context"
	klcv1alpha1 "github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1"
	"github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1/common"
)

func (r *KeptnWorkloadInstanceReconciler) reconcilePostDeployment(ctx context.Context, workloadInstance *klcv1alpha1.KeptnWorkloadInstance) (common.KeptnState, error) {
	newStatus, postDeploymentState, err := r.reconcileChecks(ctx, common.PostDeploymentCheckType, workloadInstance)
	if err != nil {
		return common.StateUnknown, err
	}
	overallState := common.GetOverallState(postDeploymentState)
	workloadInstance.Status.PostDeploymentStatus = overallState
	workloadInstance.Status.PostDeploymentTaskStatus = newStatus

	// Write Status Field
	err = r.Client.Status().Update(ctx, workloadInstance)
	if err != nil {
		return common.StateUnknown, err
	}
	return overallState, nil
}

func (r *KeptnWorkloadInstanceReconciler) reconcilePrePostDeployment(ctx context.Context, workloadInstance *klcv1alpha1.KeptnWorkloadInstance, checkType common.CheckType) (common.KeptnState, error) {
	newStatus, state, err := r.reconcileTasks(ctx, checkType, workloadInstance)
	if err != nil {
		return common.StateUnknown, err
	}
	overallState := common.GetOverallState(state)

	switch checkType {
	case common.PreDeploymentCheckType:
		workloadInstance.Status.PreDeploymentStatus = overallState
		workloadInstance.Status.PreDeploymentTaskStatus = newStatus
	case common.PostDeploymentCheckType:
		workloadInstance.Status.PostDeploymentStatus = overallState
		workloadInstance.Status.PostDeploymentTaskStatus = newStatus
	}

	// Write Status Field
	err = r.Client.Status().Update(ctx, workloadInstance)
	if err != nil {
		return common.StateUnknown, err
	}
	return overallState, nil
}
