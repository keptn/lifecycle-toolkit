package keptnappversion

import (
	"context"

	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3/common"
	controllercommon "github.com/keptn/lifecycle-toolkit/operator/controllers/common"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
	workloadInstanceList, err := r.getWorkloadInstanceList(ctx, appVersion.Namespace, appVersion.Name)
	if err != nil {
		r.Log.Error(err, "Could not get workloads")
		return apicommon.StatePending, err
	}

	if len(workloadInstanceList.Items) == 0 {
		r.Log.Info("No WorkloadInstances found")
		controllercommon.RecordEvent(r.Recorder, phase, "Warning", appVersion, "NotFound", "workloadInstances not found", appVersion.GetVersion())
		return apicommon.StatePending, nil
	}

	for _, w := range appVersion.Spec.Workloads {
		r.Log.Info("Reconciling workload " + w.Name)
		workloadStatus := apicommon.StatePending
		for _, i := range workloadInstanceList.Items {
			if w.Name == i.Name && w.Version == i.Spec.Version {
				workloadStatus = i.Status.Status
			}
		}

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
	err = r.Client.Status().Update(ctx, appVersion)
	return overallState, err
}

func (r *KeptnAppVersionReconciler) getWorkloadInstanceList(ctx context.Context, namespace string, appName string) (*klcv1alpha3.KeptnWorkloadInstanceList, error) {
	workloadInstanceList := &klcv1alpha3.KeptnWorkloadInstanceList{}
	err := r.Client.List(ctx, workloadInstanceList, client.InNamespace(namespace), client.MatchingFields{
		"spec.app": appName,
	})
	return workloadInstanceList, err
}

func getWorkloadInstanceName(namespace string, appName string, workloadName string, version string) types.NamespacedName {
	return types.NamespacedName{Namespace: namespace, Name: appName + "-" + workloadName + "-" + version}
}
