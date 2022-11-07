package keptnworkloadinstance

import (
	"context"
	"fmt"

	klcv1alpha1 "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1"
	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1/common"
	controllercommon "github.com/keptn/lifecycle-toolkit/operator/controllers/common"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *KeptnWorkloadInstanceReconciler) reconcilePrePostDeployment(ctx context.Context, workloadInstance *klcv1alpha1.KeptnWorkloadInstance, checkType common.CheckType) (common.KeptnState, error) {
	taskHandler := controllercommon.TaskHandler{
		Client:   r.Client,
		Recorder: r.Recorder,
		Log:      r.Log,
		Tracer:   r.Tracer,
		Scheme:   r.Scheme,
	}

	taskCreateAttributes := controllercommon.TaskCreateAttributes{
		SpanName:  fmt.Sprintf(common.CreateWorkloadTaskSpanName, checkType),
		CheckType: checkType,
	}

	newStatus, state, err := taskHandler.ReconcileTasks(ctx, workloadInstance, taskCreateAttributes)
	if err != nil {
		return common.StateUnknown, err
	}

	overallState := common.GetOverallState(state)

	switch checkType {
	case common.PreDeploymentCheckType:
		workloadInstance.Status.PreDeploymentStatus = overallState
		workloadInstance.Status.PreDeploymentTaskStatus = newStatus
	case common.PostDeploymentCheckType:
		workloadInstance.Status.PostDeploymentStatus = overallState
		workloadInstance.Status.PostDeploymentTaskStatus = newStatus
	}

	// Write Status Field
	err = r.Client.Status().Update(ctx, workloadInstance)
	if err != nil {
		return common.StateUnknown, err
	}
	return overallState, nil
}
