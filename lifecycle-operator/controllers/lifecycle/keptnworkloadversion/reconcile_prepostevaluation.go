//nolint:dupl
package keptnworkloadversion

import (
	"context"
	"fmt"

	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	klcv1alpha4 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha4"
	controllercommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common"
)

func (r *KeptnWorkloadVersionReconciler) reconcilePrePostEvaluation(ctx context.Context, phaseCtx context.Context, workloadVersion *klcv1alpha4.KeptnWorkloadVersion, checkType apicommon.CheckType) (apicommon.KeptnState, error) {
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

	newStatus, state, err := evaluationHandler.ReconcileEvaluations(ctx, phaseCtx, workloadVersion, evaluationCreateAttributes)
	if err != nil {
		return apicommon.StateUnknown, err
	}

	overallState := apicommon.GetOverallState(state)

	switch checkType {
	case apicommon.PreDeploymentEvaluationCheckType:
		workloadVersion.Status.PreDeploymentEvaluationStatus = overallState
		workloadVersion.Status.PreDeploymentEvaluationTaskStatus = newStatus
	case apicommon.PostDeploymentEvaluationCheckType:
		workloadVersion.Status.PostDeploymentEvaluationStatus = overallState
		workloadVersion.Status.PostDeploymentEvaluationTaskStatus = newStatus
	}

	// Write Status Field
	err = r.Client.Status().Update(ctx, workloadVersion)
	if err != nil {
		return apicommon.StateUnknown, err
	}
	return overallState, nil
}
