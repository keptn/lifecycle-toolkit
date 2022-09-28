package keptnworkloadinstance

import (
	"context"
	"fmt"
	klcv1alpha1 "github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1"
	"github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1/common"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"math/rand"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type StatusSummary struct {
	failed    int
	succeeded int
	running   int
	pending   int
}

func (r *KeptnWorkloadInstanceReconciler) reconcilePreDeployment(ctx context.Context, req ctrl.Request, workloadInstance *klcv1alpha1.KeptnWorkloadInstance) error {
	var preDeploymentState StatusSummary
	// Check if pre-deployment is already done
	if workloadInstance.IsPreDeploymentCompleted() {
		return nil
	}

	// Check current state of the PreDeploymentTasks
	var newStatus []klcv1alpha1.WorkloadTaskStatus
	for _, taskDefinitionName := range workloadInstance.Spec.PreDeploymentTasks {
		taskStatus := r.getTaskStatus(taskDefinitionName, workloadInstance.Status.PreDeploymentTaskStatus)
		task := &klcv1alpha1.KeptnTask{}
		taskExists := false

		// Check if task has already succeeded or failed
		if taskStatus.Status == common.StateSucceeded || taskStatus.Status == common.StateFailed {
			newStatus = append(newStatus, taskStatus)
			continue
		}

		// Check if Task is already created
		if taskStatus.TaskName != "" {
			err := r.Client.Get(ctx, types.NamespacedName{Name: taskStatus.TaskName, Namespace: workloadInstance.Namespace}, task)
			if err != nil && errors.IsNotFound(err) {
				taskStatus.TaskName = ""
			} else if err != nil {
				return err
			}
			taskExists = true
		}

		// Create new Task if it does not exist
		if !taskExists {
			taskName, err := r.createKeptnTask(ctx, req.Namespace, workloadInstance, taskDefinitionName, "pre")
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
	}

	for _, ns := range newStatus {
		preDeploymentState = updateStatusSummary(ns.Status, preDeploymentState)
	}
	workloadInstance.Status.PreDeploymentStatus = getOverallState(preDeploymentState)
	workloadInstance.Status.PreDeploymentTaskStatus = newStatus

	// Write Status Field
	err := r.Client.Status().Update(ctx, workloadInstance)
	if err != nil {
		return err
	}
	return nil
}

func (r *KeptnWorkloadInstanceReconciler) getTaskStatus(taskName string, instanceStatus []klcv1alpha1.WorkloadTaskStatus) klcv1alpha1.WorkloadTaskStatus {
	for _, status := range instanceStatus {
		if status.TaskDefinitionName == taskName {
			return status
		}
	}
	return klcv1alpha1.WorkloadTaskStatus{
		TaskDefinitionName: taskName,
		Status:             common.StatePending,
		TaskName:           "",
	}
}

func (r *KeptnWorkloadInstanceReconciler) getKeptnTask(ctx context.Context, taskName string, namespace string) (*klcv1alpha1.KeptnTask, error) {
	task := &klcv1alpha1.KeptnTask{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: taskName, Namespace: namespace}, task)
	if err != nil {
		return task, err
	}
	return task, nil
}

func updateStatusSummary(status common.KeptnState, summary StatusSummary) StatusSummary {
	switch status {
	case common.StateFailed:
		summary.failed++
	case common.StateSucceeded:
		summary.succeeded++
	case common.StateRunning:
		summary.running++
	case common.StatePending, "":
		summary.pending++
	}
	return summary
}

func getOverallState(summary StatusSummary) common.KeptnState {
	if summary.failed > 0 {
		return common.StateFailed
	}
	if summary.pending > 0 {
		return common.StatePending
	}
	if summary.running > 0 {
		return common.StateRunning
	}
	return common.StateSucceeded
}

func generateTaskName(checkType string, taskName string) string {
	randomId := rand.Intn(99999-10000) + 10000
	return fmt.Sprintf("%s-%s-%d", checkType, taskName, randomId)
}

func (r *KeptnWorkloadInstanceReconciler) createKeptnTask(ctx context.Context, namespace string, workloadInstance *klcv1alpha1.KeptnWorkloadInstance, taskDefinition string, checkType string) (string, error) {
	newTask := &klcv1alpha1.KeptnTask{
		ObjectMeta: metav1.ObjectMeta{
			Name:      generateTaskName(checkType, taskDefinition),
			Namespace: namespace,
		},
		Spec: klcv1alpha1.KeptnTaskSpec{
			Workload:         workloadInstance.Spec.WorkloadName,
			WorkloadVersion:  workloadInstance.Spec.Version,
			AppName:          workloadInstance.Spec.AppName,
			TaskDefinition:   taskDefinition,
			Parameters:       klcv1alpha1.TaskParameters{},
			SecureParameters: klcv1alpha1.SecureParameters{},
		},
	}
	err := controllerutil.SetControllerReference(workloadInstance, newTask, r.Scheme)
	if err != nil {
		r.Log.Error(err, "could not set controller reference:")
	}
	err = r.Client.Create(ctx, newTask)
	if err != nil {
		r.Log.Error(err, "could not create KeptnTask")
		r.Recorder.Event(workloadInstance, "Warning", "KeptnTaskNotCreated", fmt.Sprintf("Could not create KeptnTask / Namespace: %s, Name: %s ", newTask.Namespace, newTask.Name))
		return "", err
	}
	r.Recorder.Event(workloadInstance, "Normal", "KeptnTaskCreated", fmt.Sprintf("Created KeptnTask / Namespace: %s, Name: %s ", newTask.Namespace, newTask.Name))
	return newTask.Name, nil
}
