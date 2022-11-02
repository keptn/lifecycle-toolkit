package keptnappversion

import (
	"context"

	controllercommon "github.com/keptn/lifecycle-controller/operator/controllers/common"

	klcv1alpha1 "github.com/keptn/lifecycle-controller/operator/api/v1alpha1"
	"github.com/keptn/lifecycle-controller/operator/api/v1alpha1/common"
)

func (r *KeptnAppVersionReconciler) reconcilePrePostDeployment(ctx context.Context, appVersion *klcv1alpha1.KeptnAppVersion, checkType common.CheckType) (common.KeptnState, error) {
	taskHandler := controllercommon.TaskHandler{
		Client:      r.Client,
		Recorder:    r.Recorder,
		Log:         r.Log,
		SpanHandler: r.SpanHandler,
		Tracer:      r.Tracer,
		Scheme:      r.Scheme,
	}
	newStatus, state, err := taskHandler.ReconcileTasks(ctx, checkType, appVersion, true)
	if err != nil {
		return common.StateUnknown, err
	}
	overallState := common.GetOverallState(state)

	switch checkType {
	case common.PreDeploymentCheckType:
		appVersion.Status.PreDeploymentStatus = overallState
		appVersion.Status.PreDeploymentTaskStatus = newStatus
	case common.PostDeploymentCheckType:
		appVersion.Status.PostDeploymentStatus = overallState
		appVersion.Status.PostDeploymentTaskStatus = newStatus
	}

	// Write Status Field
	err = r.Client.Status().Update(ctx, appVersion)
	if err != nil {
		return common.StateUnknown, err
	}
	return overallState, nil
}
