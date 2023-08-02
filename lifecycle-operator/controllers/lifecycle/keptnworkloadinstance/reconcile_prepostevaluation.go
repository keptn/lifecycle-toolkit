//nolint:dupl
package keptnworkloadinstance

import (
	"context"
	"fmt"

	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	controllercommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common"
)

func (r *KeptnWorkloadInstanceReconciler) reconcilePrePostEvaluation(ctx context.Context, phaseCtx context.Context, workloadInstance *klcv1alpha3.KeptnWorkloadInstance, checkType apicommon.CheckType) (apicommon.KeptnState, error) {
	evaluationHandler := controllercommon.EvaluationHandler{
		Client:      r.Client,
		EventSender: r.EventSender,
		Log:         r.Log,
		Tracer:      r.getTracer(),
		Scheme:      r.Scheme,
		SpanHandler: r.SpanHandler,
	}

	evaluationCreateAttributes := controllercommon.CreateEvaluationAttributes{
		SpanName:  fmt.Sprintf(apicommon.CreateWorkloadEvalSpanName, checkType),
		CheckType: checkType,
	}

	newStatus, state, err := evaluationHandler.ReconcileEvaluations(ctx, phaseCtx, workloadInstance, evaluationCreateAttributes)
	if err != nil {
		return apicommon.StateUnknown, err
	}

	overallState := apicommon.GetOverallState(state)

	switch checkType {
	case apicommon.PreDeploymentEvaluationCheckType:
		workloadInstance.Status.PreDeploymentEvaluationStatus = overallState
		workloadInstance.Status.PreDeploymentEvaluationTaskStatus = newStatus
	case apicommon.PostDeploymentEvaluationCheckType:
		workloadInstance.Status.PostDeploymentEvaluationStatus = overallState
		workloadInstance.Status.PostDeploymentEvaluationTaskStatus = newStatus
	}

	// Write Status Field
	err = r.Client.Status().Update(ctx, workloadInstance)
	if err != nil {
		return apicommon.StateUnknown, err
	}
	return overallState, nil
}
