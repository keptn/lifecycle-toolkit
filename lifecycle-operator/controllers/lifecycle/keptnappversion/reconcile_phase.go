//nolint:dupl
package keptnappversion

import (
	"context"
	"fmt"

	apilifecycle "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/task"
)

func (r *KeptnAppVersionReconciler) reconcilePhase(ctx context.Context, phaseCtx context.Context, appVersion *apilifecycle.KeptnAppVersion, checkType apicommon.CheckType) (apicommon.KeptnState, error) {
	taskHandler := task.Handler{
		Client:      r.Client,
		EventSender: r.EventSender,
		Log:         r.Log,
		Tracer:      r.getTracer(),
		Scheme:      r.Scheme,
		SpanHandler: r.SpanHandler,
	}

	taskCreateAttributes := task.CreateTaskAttributes{
		SpanName:  fmt.Sprintf(apicommon.CreateAppTaskSpanName, checkType),
		CheckType: checkType,
	}

	newStatus, state, err := taskHandler.ReconcileTasks(ctx, phaseCtx, appVersion, taskCreateAttributes)
	if err != nil {
		return apicommon.StateUnknown, err
	}
	overallState := apicommon.GetOverallStateBlockedDeployment(state, r.Config.GetBlockDeployment())

	switch checkType {
	case apicommon.PreDeploymentCheckType:
		appVersion.Status.PreDeploymentStatus = overallState
		appVersion.Status.PreDeploymentTaskStatus = newStatus
	case apicommon.PostDeploymentCheckType:
		appVersion.Status.PostDeploymentStatus = overallState
		appVersion.Status.PostDeploymentTaskStatus = newStatus
	case apicommon.PromotionCheckType:
		appVersion.Status.PromotionStatus = overallState
		appVersion.Status.PromotionTaskStatus = newStatus
	}

	// Write Status Field
	err = r.Client.Status().Update(ctx, appVersion)
	if err != nil {
		return apicommon.StateUnknown, err
	}
	return overallState, nil
}
