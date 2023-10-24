//nolint:dupl
package keptnappversion

import (
	"context"
	"fmt"

	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/evaluation"
)

func (r *KeptnAppVersionReconciler) reconcilePrePostEvaluation(ctx context.Context, phaseCtx context.Context, appVersion *klcv1alpha3.KeptnAppVersion, checkType apicommon.CheckType) (apicommon.KeptnState, error) {
	evaluationHandler := evaluation.EvaluationHandler{
		Client:      r.Client,
		EventSender: r.EventSender,
		Log:         r.Log,
		Tracer:      r.getTracer(),
		Scheme:      r.Scheme,
		SpanHandler: r.SpanHandler,
	}

	evaluationCreateAttributes := evaluation.CreateEvaluationAttributes{
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
