//nolint:dupl
package keptnappversion

import (
	"context"
	"fmt"

	apilifecycle "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/evaluation"
)

func (r *KeptnAppVersionReconciler) reconcilePrePostEvaluation(ctx context.Context, phaseCtx context.Context, appVersion *apilifecycle.KeptnAppVersion, checkType apicommon.CheckType) (apicommon.KeptnState, error) {
	evaluationCreateAttributes := evaluation.CreateEvaluationAttributes{
		SpanName:  fmt.Sprintf(apicommon.CreateAppEvalSpanName, checkType),
		CheckType: checkType,
	}
	evaluationHandler := evaluation.NewHandler(
		r.Client,
		r.EventSender,
		r.Log,
		r.getTracer(),
		r.Client.Scheme(),
		r.SpanHandler,
	)

	newStatus, state, err := evaluationHandler.ReconcileEvaluations(ctx, phaseCtx, appVersion, evaluationCreateAttributes)
	if err != nil {
		return apicommon.StateUnknown, err
	}

	overallState := apicommon.GetOverallStateBlockedDeployment(state, r.Config.GetBlockDeployment())

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
