package keptnappversion

import (
	"context"

	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3/common"
	controllercommon "github.com/keptn/lifecycle-toolkit/operator/controllers/common"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

func (r *KeptnAppVersionReconciler) reconcileWorkloads(ctx context.Context, appVersion *klcv1alpha3.KeptnAppVersion) (apicommon.KeptnState, error) {
	r.Log.Info("Reconciling Workloads")
	var summary apicommon.StatusSummary
	summary.Total = len(appVersion.Spec.Workloads)

	phase := apicommon.KeptnPhaseType{
		ShortName: "ReconcileWorkload",
		LongName:  "Reconcile Workloads",
	}

	var newStatus []klcv1alpha3.WorkloadStatus
	for _, w := range appVersion.Spec.Workloads {
		r.Log.Info("Reconciling workload " + w.Name)
		workload, err := r.getWorkloadInstance(ctx, getWorkloadInstanceName(appVersion.Namespace, appVersion.Spec.AppName, w.Name, w.Version))
		if err != nil && errors.IsNotFound(err) {
			controllercommon.RecordEvent(r.Recorder, phase, "Warning", appVersion, "NotFound", "workloadInstance not found", appVersion.GetVersion())
			workload.Status.Status = apicommon.StatePending
		} else if err != nil {
			r.Log.Error(err, "Could not get workload")
			workload.Status.Status = apicommon.StateUnknown
		}
		workloadStatus := workload.Status.Status

		newStatus = append(newStatus, klcv1alpha3.WorkloadStatus{
			Workload: w,
			Status:   workloadStatus,
		})
		summary = apicommon.UpdateStatusSummary(workloadStatus, summary)
	}

	overallState := apicommon.GetOverallState(summary)
	appVersion.Status.WorkloadOverallStatus = overallState
	r.Log.Info("Overall state of workloads", "state", appVersion.Status.WorkloadOverallStatus)

	appVersion.Status.WorkloadStatus = newStatus
	r.Log.Info("Workload status", "status", appVersion.Status.WorkloadStatus)

	// Write Status Field
	err := r.Client.Status().Update(ctx, appVersion)
	return overallState, err
}

func (r *KeptnAppVersionReconciler) getWorkloadInstance(ctx context.Context, workload types.NamespacedName) (klcv1alpha3.KeptnWorkloadInstance, error) {
	workloadInstance := &klcv1alpha3.KeptnWorkloadInstance{}
	err := r.Get(ctx, workload, workloadInstance)
	return *workloadInstance, err
}

func getWorkloadInstanceName(namespace string, appName string, workloadName string, version string) types.NamespacedName {
	return types.NamespacedName{Namespace: namespace, Name: appName + "-" + workloadName + "-" + version}
}
