package keptnappversion

import (
	"context"

	klcv1alpha1 "github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1"
	"github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1/common"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *KeptnAppVersionReconciler) reconcilePostDeployment(ctx context.Context, req ctrl.Request, appVersion *klcv1alpha1.KeptnAppVersion) error {

	newStatus, postDeploymentState, err := r.reconcileChecks(ctx, common.PostDeploymentCheckType, appVersion)
	if err != nil {
		return err
	}
	appVersion.Status.PostDeploymentStatus = common.GetOverallState(postDeploymentState)
	appVersion.Status.PostDeploymentTaskStatus = newStatus

	// Write Status Field
	err = r.Client.Status().Update(ctx, appVersion)
	if err != nil {
		return err
	}
	return nil
}
