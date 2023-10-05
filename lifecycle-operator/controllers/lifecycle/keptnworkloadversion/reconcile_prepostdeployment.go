//nolint:dupl
package keptnworkloadversion

import (
	"context"
	"fmt"

	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	controllercommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common"
)

func (r *KeptnWorkloadVersionReconciler) reconcilePrePostDeployment(ctx context.Context, phaseCtx context.Context, workloadVersion *klcv1alpha3.KeptnWorkloadVersion, checkType apicommon.CheckType) (apicommon.KeptnState, error) {
	taskHandler := controllercommon.TaskHandler{
		Client:      r.Client,
		EventSender: r.EventSender,
		Log:         r.Log,
		Tracer:      r.getTracer(),
		Scheme:      r.Scheme,
		SpanHandler: r.SpanHandler,
	}

	taskCreateAttributes := controllercommon.CreateTaskAttributes{
		SpanName:  fmt.Sprintf(apicommon.CreateWorkloadTaskSpanName, checkType),
		CheckType: checkType,
	}

	newStatus, state, err := taskHandler.ReconcileTasks(ctx, phaseCtx, workloadVersion, taskCreateAttributes)
	if err != nil {
		return apicommon.StateUnknown, err
	}

	overallState := apicommon.GetOverallState(state)

	switch checkType {
	case apicommon.PreDeploymentCheckType:
		workloadVersion.Status.PreDeploymentStatus = overallState
		workloadVersion.Status.PreDeploymentTaskStatus = newStatus
	case apicommon.PostDeploymentCheckType:
		workloadVersion.Status.PostDeploymentStatus = overallState
		workloadVersion.Status.PostDeploymentTaskStatus = newStatus
	}

	// Write Status Field
	err = r.Client.Status().Update(ctx, workloadVersion)
	if err != nil {
		return apicommon.StateUnknown, err
	}
	return overallState, nil
}
