package keptnappversion

import (
	"context"
	"fmt"

	controllercommon "github.com/keptn/lifecycle-toolkit/operator/controllers/common"

	klcv1alpha1 "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1"
	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1/common"
)

func (r *KeptnAppVersionReconciler) reconcilePrePostDeployment(ctx context.Context, appVersion *klcv1alpha1.KeptnAppVersion, checkType common.CheckType) (common.KeptnState, error) {
	taskHandler := controllercommon.TaskHandler{
		Client:   r.Client,
		Recorder: r.Recorder,
		Log:      r.Log,
		Tracer:   r.Tracer,
		Scheme:   r.Scheme,
	}

	taskCreateAttributes := controllercommon.TaskCreateAttributes{
		SpanName:  fmt.Sprintf(common.CreateAppTaskSpanName, checkType),
		CheckType: checkType,
	}

	newStatus, state, err := taskHandler.ReconcileTasks(ctx, appVersion, taskCreateAttributes)
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
