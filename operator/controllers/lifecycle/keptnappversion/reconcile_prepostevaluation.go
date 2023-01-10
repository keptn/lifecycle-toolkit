package keptnappversion

import (
	"context"
	"fmt"

	klcv1alpha2 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2/common"
	controllercommon "github.com/keptn/lifecycle-toolkit/operator/controllers/common"
)

func (r *KeptnAppVersionReconciler) reconcilePrePostEvaluation(ctx context.Context, phaseCtx context.Context, appVersion *klcv1alpha2.KeptnAppVersion, checkType apicommon.CheckType) (apicommon.KeptnState, error) {
	evaluationHandler := controllercommon.EvaluationHandler{
		Client:      r.Client,
		Recorder:    r.Recorder,
		Log:         r.Log,
		Tracer:      r.Tracer,
		Scheme:      r.Scheme,
		SpanHandler: r.SpanHandler,
	}

	evaluationCreateAttributes := controllercommon.CreateAttributes{
		SpanName:  fmt.Sprintf(apicommon.CreateAppEvalSpanName, checkType),
		CheckType: checkType,
	}

	newStatus, state, err := evaluationHandler.ReconcileEvaluations(ctx, phaseCtx, appVersion, evaluationCreateAttributes)
	if err != nil {
		return apicommon.StateUnknown, err
	}

	overallState := apicommon.GetOverallState(state)

	switch checkType {
	case apicommon.PreDeploymentEvaluationCheckType:
		appVersion.Status.PreDeploymentEvaluationStatus = overallState
		appVersion.Status.PreDeploymentEvaluationTaskStatus = newStatus
	case apicommon.PostDeploymentEvaluationCheckType:
		appVersion.Status.PostDeploymentEvaluationStatus = overallState
		appVersion.Status.PostDeploymentEvaluationTaskStatus = newStatus
	}

	// Write Status Field
	err = r.Client.Status().Update(ctx, appVersion)
	if err != nil {
		return apicommon.StateUnknown, err
	}
	return overallState, nil
}
