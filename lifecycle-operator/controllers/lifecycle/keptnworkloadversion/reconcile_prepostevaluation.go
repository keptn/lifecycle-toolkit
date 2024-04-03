//nolint:dupl
package keptnworkloadversion

import (
	"context"
	"fmt"

	apilifecycle "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/evaluation"
)

func (r *KeptnWorkloadVersionReconciler) reconcilePrePostEvaluation(ctx context.Context, phaseCtx context.Context, workloadVersion *apilifecycle.KeptnWorkloadVersion, checkType apicommon.CheckType) (apicommon.KeptnState, error) {
	evaluationCreateAttributes := evaluation.CreateEvaluationAttributes{
		SpanName:  fmt.Sprintf(apicommon.CreateWorkloadEvalSpanName, checkType),
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

	newStatus, state, err := evaluationHandler.ReconcileEvaluations(ctx, phaseCtx, workloadVersion, evaluationCreateAttributes)
	if err != nil {
		return apicommon.StateUnknown, err
	}

	overallState := apicommon.GetOverallStateBlockedDeployment(state, r.Config.GetBlockDeployment())

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
