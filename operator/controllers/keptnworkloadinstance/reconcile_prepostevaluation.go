package keptnworkloadinstance

import (
	"context"
	"fmt"

	klcv1alpha1 "github.com/keptn/lifecycle-controller/operator/api/v1alpha1"
	"github.com/keptn/lifecycle-controller/operator/api/v1alpha1/common"
	controllercommon "github.com/keptn/lifecycle-controller/operator/controllers/common"
)

func (r *KeptnWorkloadInstanceReconciler) reconcilePrePostEvaluation(ctx context.Context, workloadInstance *klcv1alpha1.KeptnWorkloadInstance, checkType common.CheckType) (common.KeptnState, error) {
	evaluationHandler := controllercommon.EvaluationHandler{
		Client:   r.Client,
		Recorder: r.Recorder,
		Log:      r.Log,
		Tracer:   r.Tracer,
		Scheme:   r.Scheme,
	}

	evaluationCreateAttributes := controllercommon.EvaluationCreateAttributes{
		SpanName:  fmt.Sprintf(common.CreateWorkloadEvalSpanName, checkType),
		CheckType: checkType,
	}

	newStatus, state, err := evaluationHandler.ReconcileEvaluations(ctx, workloadInstance, evaluationCreateAttributes)
	if err != nil {
		return common.StateUnknown, err
	}

	overallState := common.GetOverallState(state)

	switch checkType {
	case common.PreDeploymentEvaluationCheckType:
		workloadInstance.Status.PreDeploymentEvaluationStatus = overallState
		workloadInstance.Status.PreDeploymentEvaluationTaskStatus = newStatus
	case common.PostDeploymentEvaluationCheckType:
		workloadInstance.Status.PostDeploymentEvaluationStatus = overallState
		workloadInstance.Status.PostDeploymentEvaluationTaskStatus = newStatus
	}

	// Write Status Field
	err = r.Client.Status().Update(ctx, workloadInstance)
	if err != nil {
		return common.StateUnknown, err
	}
	return overallState, nil
}
