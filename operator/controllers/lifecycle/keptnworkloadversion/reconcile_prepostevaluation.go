//nolint:dupl
package keptnworkloadversion

import (
	"context"
	"fmt"

	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3/common"
	controllercommon "github.com/keptn/lifecycle-toolkit/operator/controllers/common"
)

func (r *KeptnWorkloadVersionReconciler) reconcilePrePostEvaluation(ctx context.Context, phaseCtx context.Context, workloadVersion *klcv1alpha3.KeptnWorkloadVersion, checkType apicommon.CheckType) (apicommon.KeptnState, error) {
	evaluationHandler := controllercommon.EvaluationHandler{
		Client:      r.Client,
		Recorder:    r.Recorder,
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
