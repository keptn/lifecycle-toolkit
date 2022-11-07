package keptnappversion

import (
	"context"
	"fmt"

	klcv1alpha1 "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1"
	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1/common"
	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1/semconv"
	controllercommon "github.com/keptn/lifecycle-toolkit/operator/controllers/common"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *KeptnAppVersionReconciler) reconcilePrePostEvaluation(ctx context.Context, appVersion *klcv1alpha1.KeptnAppVersion, checkType common.CheckType) (common.KeptnState, error) {
	evaluationHandler := controllercommon.EvaluationHandler{
		Client:   r.Client,
		Recorder: r.Recorder,
		Log:      r.Log,
		Tracer:   r.Tracer,
		Scheme:   r.Scheme,
	}

	evaluationCreateAttributes := controllercommon.EvaluationCreateAttributes{
		SpanName:  fmt.Sprintf(common.CreateAppEvalSpanName, checkType),
		CheckType: checkType,
	}

	newStatus, state, err := evaluationHandler.ReconcileEvaluations(ctx, appVersion, evaluationCreateAttributes)
	if err != nil {
		return common.StateUnknown, err
	}

	overallState := common.GetOverallState(state)

	switch checkType {
	case common.PreDeploymentEvaluationCheckType:
		appVersion.Status.PreDeploymentEvaluationStatus = overallState
		appVersion.Status.PreDeploymentEvaluationTaskStatus = newStatus
	case common.PostDeploymentEvaluationCheckType:
		appVersion.Status.PostDeploymentEvaluationStatus = overallState
		appVersion.Status.PostDeploymentEvaluationTaskStatus = newStatus
	}

	// Write Status Field
	err = r.Client.Status().Update(ctx, appVersion)
	if err != nil {
		return common.StateUnknown, err
	}
	return overallState, nil
}
