package keptnworkloadinstance

import (
	"context"
	klcv1alpha1 "github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1"
	"github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1/common"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

var PostDeploymentState StatusSummary

func (r *KeptnWorkloadInstanceReconciler) reconcilePostDeployment(ctx context.Context, req ctrl.Request, workloadInstance *klcv1alpha1.KeptnWorkloadInstance) error {

	// Check if pre-deployment is already done
	if workloadInstance.IsPostDeploymentCompleted() {
		return nil
	}

	// Check current state of the PostDeploymentTasks
	var newStatus []klcv1alpha1.WorkloadTaskStatus
	for _, taskDefinitionName := range workloadInstance.Spec.PostDeploymentTasks {
		taskStatus := r.getTaskStatus(taskDefinitionName, workloadInstance.Status.PostDeploymentTaskStatus)
		task := &klcv1alpha1.KeptnTask{}
		taskExists := false

		// Create new state entry for the pre-deployment Task if it does not exist
		if taskStatus == (klcv1alpha1.WorkloadTaskStatus{}) {
			taskStatus = klcv1alpha1.WorkloadTaskStatus{
				TaskDefinitionName: taskDefinitionName,
				Status:             common.StatePending,
				TaskName:           "",
			}
		}
		// Check if task has already succeeded or failed
		if taskStatus.Status == common.StateSucceeded || taskStatus.Status == common.StateFailed {
			continue
		}

		// Check if Task is already created
		if taskStatus.TaskName != "" {
			err := r.Client.Get(ctx, types.NamespacedName{Name: taskStatus.TaskName, Namespace: workloadInstance.Namespace}, task)
			if err != nil && errors.IsNotFound(err) {
				taskStatus.TaskName = ""
			} else {
				return err
			}
			taskExists = true
		}

		// Create new Task if it does not exist
		if !taskExists {
			taskName, err := r.createKeptnTask(ctx, req.Namespace, workloadInstance, taskDefinitionName)
			if err != nil {
				return err
			}
			taskStatus.TaskName = taskName
		} else {
			// Update state of Task if it is already created
			taskStatus.Status = task.Status.Status
		}
		// Update state of the Pre-Deployment Task
		newStatus = append(newStatus, taskStatus)

		// Update overall state for Pre-Deployment
		updateStatusSummary(taskStatus.Status, PostDeploymentState)
	}

	workloadInstance.Status.PostDeploymentStatus = getOverallState(PostDeploymentState)
	workloadInstance.Status.PostDeploymentTaskStatus = newStatus

	// Write Status Field
	err := r.Client.Status().Update(ctx, workloadInstance)
	if err != nil {
		return err
	}
	return nil
}
