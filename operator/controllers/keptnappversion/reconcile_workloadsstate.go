package keptnappversion

import (
	"context"
	"fmt"

	klcv1alpha1 "github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1"
	"github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1/common"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

func (r *KeptnAppVersionReconciler) reconcileWorkloads(ctx context.Context, appVersion *klcv1alpha1.KeptnAppVersion) (common.KeptnState, error) {
	r.Log.Info("Reconciling Workloads")
	var summary common.StatusSummary
	summary.Total = len(appVersion.Spec.Workloads)

	var newStatus []klcv1alpha1.WorkloadStatus
	for _, w := range appVersion.Spec.Workloads {
		r.Log.Info("Reconciling workload " + w.Name)
		workload, err := r.getWorkloadInstance(ctx, getWorkloadInstanceName(appVersion.Namespace, appVersion.Spec.AppName, w.Name, w.Version))
		if err != nil && errors.IsNotFound(err) {
			r.Recorder.Event(appVersion, "Warning", "WorkloadNotFound", fmt.Sprintf("Could not find KeptnWorkloadInstance / Namespace: %s, Name: %s ", appVersion.Namespace, w.Name))
			workload.Status.Status = common.StatePending
		} else if err != nil {
			r.Log.Error(err, "Could not get workload")
			workload.Status.Status = common.StateUnknown
		}
		workloadStatus := workload.Status.Status

		newStatus = append(newStatus, klcv1alpha1.WorkloadStatus{
			Workload: w,
			Status:   workloadStatus,
		})
		summary = common.UpdateStatusSummary(workloadStatus, summary)
	}

	overallState := common.GetOverallState(summary)
	appVersion.Status.WorkloadOverallStatus = overallState
	r.Log.Info("Overall state of workloads", "state", appVersion.Status.WorkloadOverallStatus)

	appVersion.Status.WorkloadStatus = newStatus
	r.Log.Info("Workload status", "status", appVersion.Status.WorkloadStatus)

	// Write Status Field
	err := r.Client.Status().Update(ctx, appVersion)
	return overallState, err
}

func (r *KeptnAppVersionReconciler) getWorkloadInstance(ctx context.Context, workload types.NamespacedName) (klcv1alpha1.KeptnWorkloadInstance, error) {
	workloadInstance := &klcv1alpha1.KeptnWorkloadInstance{}
	err := r.Get(ctx, workload, workloadInstance)
	return *workloadInstance, err
}

func getWorkloadInstanceName(namespace string, appName string, workloadName string, version string) types.NamespacedName {
	return types.NamespacedName{Namespace: namespace, Name: appName + "-" + workloadName + "-" + version}
}
