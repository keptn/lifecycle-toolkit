//nolint:dupl
package keptnappversion

import (
	"context"
	"fmt"

	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	controllercommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common"
)

func (r *KeptnAppVersionReconciler) reconcilePrePostDeployment(ctx context.Context, phaseCtx context.Context, appVersion *klcv1alpha3.KeptnAppVersion, checkType apicommon.CheckType) (apicommon.KeptnState, error) {
	taskHandler := controllercommon.TaskHandler{
		Client:      r.Client,
		EventSender: r.EventSender,
		Log:         r.Log,
		Tracer:      r.getTracer(),
		Scheme:      r.Scheme,
		SpanHandler: r.SpanHandler,
	}

	taskCreateAttributes := controllercommon.CreateTaskAttributes{
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
	err = r.Client.Update(ctx, appVersion)
	if err != nil {
		return apicommon.StateUnknown, err
	}
	return overallState, nil
}
