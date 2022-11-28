package keptnappversion

import (
	"context"
	"fmt"

	klcv1alpha1 "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1/common"
	controllercommon "github.com/keptn/lifecycle-toolkit/operator/controllers/common"
)

func (r *KeptnAppVersionReconciler) reconcilePrePostDeployment(ctx context.Context, phaseCtx context.Context, appVersion *klcv1alpha1.KeptnAppVersion, checkType apicommon.CheckType) (apicommon.KeptnState, error) {
	taskHandler := controllercommon.TaskHandler{
		Client:      r.Client,
		Recorder:    r.Recorder,
		Log:         r.Log,
		Tracer:      r.Tracer,
		Scheme:      r.Scheme,
		SpanHandler: r.SpanHandler,
	}

	taskCreateAttributes := controllercommon.TaskCreateAttributes{
		SpanName:  fmt.Sprintf(apicommon.CreateAppTaskSpanName, checkType),
		CheckType: checkType,
	}

	newStatus, state, err := taskHandler.ReconcileTasks(ctx, phaseCtx, appVersion, taskCreateAttributes)
	if err != nil {
		return apicommon.StateUnknown, err
	}
	overallState := apicommon.GetOverallState(state)

	switch checkType {
	case apicommon.PreDeploymentCheckType:
		appVersion.Status.PreDeploymentStatus = overallState
		appVersion.Status.PreDeploymentTaskStatus = newStatus
	case apicommon.PostDeploymentCheckType:
		appVersion.Status.PostDeploymentStatus = overallState
		appVersion.Status.PostDeploymentTaskStatus = newStatus
	}

	// Write Status Field
	err = r.Client.Status().Update(ctx, appVersion)
	if err != nil {
		return apicommon.StateUnknown, err
	}
	return overallState, nil
}
