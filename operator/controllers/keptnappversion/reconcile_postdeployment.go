package keptnappversion

import (
	"context"

	klcv1alpha1 "github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1"
	"github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1/common"
)

func (r *KeptnAppVersionReconciler) reconcilePostDeployment(ctx context.Context, appVersion *klcv1alpha1.KeptnAppVersion) (common.KeptnState, error) {

	newStatus, postDeploymentState, err := r.reconcileTasks(ctx, common.PostDeploymentCheckType, appVersion)
	if err != nil {
		return common.StateUnknown, err
	}
	overallState := common.GetOverallState(postDeploymentState)
	appVersion.Status.PostDeploymentStatus = overallState
	appVersion.Status.PostDeploymentTaskStatus = newStatus

	// Write Status Field
	err = r.Client.Status().Update(ctx, appVersion)
	if err != nil {
		return common.StateUnknown, err
	}
	return overallState, nil
}
