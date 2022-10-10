package keptnappversion

import (
	"context"
	klcv1alpha1 "github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1"
	"github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1/common"
	"k8s.io/apimachinery/pkg/types"
)

func (r *KeptnAppVersionReconciler) reconcileWorkloads(ctx context.Context, appVersion *klcv1alpha1.KeptnAppVersion) error {

	var summary common.StatusSummary

	var newStatus []klcv1alpha1.WorkloadStatus
	for _, w := range appVersion.Spec.Workloads {
		workload, err := r.getWorkload(ctx, getWorkloadInstanceName(appVersion.Namespace, appVersion.Spec.AppName, w.Name, w.Version))
		if err != nil {
			return err
		}
		workloadStatus := workload.Status.PostDeploymentStatus

		newStatus = append(newStatus, klcv1alpha1.WorkloadStatus{
			Workload: w,
			Status:   workloadStatus,
		})
		summary.UpdateStatusSummary(workloadStatus)
	}

	appVersion.Status.WorkloadOverallStatus = common.GetOverallState(summary)
	appVersion.Status.WorkloadStatus = newStatus

	// Write Status Field
	err := r.Client.Status().Update(ctx, appVersion)
	return err
}

func (r *KeptnAppVersionReconciler) getWorkload(ctx context.Context, workload types.NamespacedName) (klcv1alpha1.KeptnWorkloadInstance, error) {
	workloadInstance := &klcv1alpha1.KeptnWorkloadInstance{}
	err := r.Get(ctx, workload, workloadInstance)
	return *workloadInstance, err
}

func getWorkloadInstanceName(namespace string, appName string, workloadName string, version string) types.NamespacedName {
	return types.NamespacedName{Namespace: namespace, Name: appName + "-" + workloadName + "-" + version}
}
