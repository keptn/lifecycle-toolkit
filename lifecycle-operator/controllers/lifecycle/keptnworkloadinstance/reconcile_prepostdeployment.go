//nolint:dupl
package keptnworkloadinstance

import (
	"context"
	"fmt"

	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	controllercommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common"
)

func (r *KeptnWorkloadInstanceReconciler) reconcilePrePostDeployment(ctx context.Context, phaseCtx context.Context, workloadInstance *klcv1alpha3.KeptnWorkloadInstance, checkType apicommon.CheckType) (apicommon.KeptnState, error) {
	taskHandler := controllercommon.TaskHandler{
		Client:           r.Client,
		EventSender:      r.EventSender,
		Log:              r.Log,
		Tracer:           r.getTracer(),
		Scheme:           r.Scheme,
		SpanHandler:      r.SpanHandler,
		DefaultNamespace: r.Namespace,
	}

	taskCreateAttributes := controllercommon.CreateTaskAttributes{
		SpanName:  fmt.Sprintf(apicommon.CreateWorkloadTaskSpanName, checkType),
		CheckType: checkType,
	}

	newStatus, state, err := taskHandler.ReconcileTasks(ctx, phaseCtx, workloadInstance, taskCreateAttributes)
	if err != nil {
		return apicommon.StateUnknown, err
	}

	overallState := apicommon.GetOverallState(state)

	switch checkType {
	case apicommon.PreDeploymentCheckType:
		workloadInstance.Status.PreDeploymentStatus = overallState
		workloadInstance.Status.PreDeploymentTaskStatus = newStatus
	case apicommon.PostDeploymentCheckType:
		workloadInstance.Status.PostDeploymentStatus = overallState
		workloadInstance.Status.PostDeploymentTaskStatus = newStatus
	}

	// Write Status Field
	err = r.Client.Status().Update(ctx, workloadInstance)
	if err != nil {
		return apicommon.StateUnknown, err
	}
	return overallState, nil
}
