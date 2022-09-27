package keptnworkloadinstance

import (
	"context"

	klcv1alpha1 "github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *KeptnWorkloadInstanceReconciler) reconcilePostDeployment(ctx context.Context, req ctrl.Request, workloadInstance *klcv1alpha1.KeptnWorkloadInstance) (ctrl.Result, error) {

	return ctrl.Result{}, nil
}
