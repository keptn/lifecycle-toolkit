package keptnworkloadinstance

import (
	"context"

	klcv1alpha1 "github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1"
	"github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1/common"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *KeptnWorkloadInstanceReconciler) reconcilePreDeployment(ctx context.Context, req ctrl.Request, workloadInstance *klcv1alpha1.KeptnWorkloadInstance) error {
	newStatus, preDeploymentState, err := r.reconcileChecks(ctx, common.PreDeploymentCheckType, workloadInstance)
	if err != nil {
		return err
	}
	r.Log.Info("Pre-Deployment Information", "Pre-Deployment State", preDeploymentState, "Pre-Deployment Workload State", newStatus)
	workloadInstance.Status.PreDeploymentStatus = common.GetOverallState(preDeploymentState)
	workloadInstance.Status.PreDeploymentTaskStatus = newStatus

	// Write Status Field
	err = r.Client.Status().Update(ctx, workloadInstance)
	if err != nil {
		return err
	}
	return nil
}
