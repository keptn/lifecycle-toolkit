package keptnworkloadinstance

import (
	"context"
	"fmt"
	"github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1/common"
	"math/rand"
	"time"

	klcv1alpha1 "github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type StatusSummary struct {
	failed    int
	succeeded int
	running   int
	pending   int
}

var preDeploymentState StatusSummary

func (r *KeptnWorkloadInstanceReconciler) reconcilePreDeployment(ctx context.Context, req ctrl.Request, workloadInstance *klcv1alpha1.KeptnWorkloadInstance) (ctrl.Result, error) {

	if workloadInstance.IsPreDeploymentCompleted() {
		return ctrl.Result{Requeue: true, RequeueAfter: 10 * time.Second}, nil
	}

	if workloadInstance.Status.PreDeploymentStatus == common.StatePending || workloadInstance.Status.PreDeploymentStatus == "" {
		var newStatus []klcv1alpha1.WorkloadTaskStatus
		// tasks not created yet, need to create them
		for _, taskDefinition := range workloadInstance.Spec.PreDeploymentTasks {
			taskName, err := r.createKeptnTask(ctx, req, workloadInstance, taskDefinition)
			if err != nil {
				return reconcile.Result{}, err
			}
			newStatus = append(newStatus, klcv1alpha1.WorkloadTaskStatus{
				TaskDefinitionName: taskDefinition,
				Status:             common.StatePending,
				TaskName:           taskName,
			})
		}
		workloadInstance.Status.PreDeploymentTaskStatus = newStatus
		workloadInstance.Status.PreDeploymentStatus = common.StateRunning
		err := r.Client.Status().Update(ctx, workloadInstance)
		if err != nil {
			return ctrl.Result{}, err
		}

		return ctrl.Result{Requeue: true, RequeueAfter: 10 * time.Second}, nil
	}
	// tasks exist, check status
	summary := StatusSummary{0, 0, 0, 0}
	for _, taskStatus := range workloadInstance.Status.PreDeploymentTaskStatus {
		if taskStatus.Status != common.StateFailed && taskStatus.Status != common.StateSucceeded {
			task, err := r.getKeptnTask(ctx, taskStatus.TaskName, workloadInstance.Namespace)
			if err != nil {
				return ctrl.Result{}, err
			}
			taskStatus.Status = task.Status.Status
		}
		updateStatusSummary(taskStatus.Status, summary)
	}

	workloadInstance.Status.PreDeploymentStatus = getOverallState(summary)
	err := r.Client.Status().Update(ctx, workloadInstance)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{Requeue: true, RequeueAfter: 10 * time.Second}, nil
}

func (r *KeptnWorkloadInstanceReconciler) getTaskStatus(taskName string, instanceStatus []klcv1alpha1.WorkloadTaskStatus) klcv1alpha1.WorkloadTaskStatus {
	for _, status := range instanceStatus {
		if status.TaskDefinitionName == taskName {
			return status
		}
	}
	return klcv1alpha1.WorkloadTaskStatus{}
}

func (r *KeptnWorkloadInstanceReconciler) getKeptnTask(ctx context.Context, taskName string, namespace string) (*klcv1alpha1.KeptnTask, error) {
	task := &klcv1alpha1.KeptnTask{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: taskName, Namespace: namespace}, task)
	if err != nil {
		return task, err
	}
	return task, nil
}

func updateStatusSummary(status common.KeptnState, summary StatusSummary) {
	switch status {
	case common.StateFailed:
		summary.failed++
	case common.StateSucceeded:
		summary.succeeded++
	case common.StateRunning:
		summary.running++
	case common.StatePending:
		summary.pending++
	}
}

func getOverallState(summary StatusSummary) common.KeptnState {
	if summary.failed > 0 {
		return common.StateFailed
	}
	if summary.running > 0 {
		return common.StateRunning
	}
	if summary.pending > 0 {
		return common.StatePending
	}
	return common.StateSucceeded
}

func generateTaskName(instance klcv1alpha1.KeptnWorkloadInstance, taskName string) string {
	randomId := rand.Intn(99999-10000) + 10000
	return fmt.Sprintf("%s-%s-%d", instance.Name, taskName, randomId)
}

func (r *KeptnWorkloadInstanceReconciler) createKeptnTask(ctx context.Context, req ctrl.Request, workloadInstance *klcv1alpha1.KeptnWorkloadInstance, taskDefinition string) (string, error) {
	newTask := &klcv1alpha1.KeptnTask{
		ObjectMeta: metav1.ObjectMeta{
			Name:      generateTaskName(*workloadInstance, taskDefinition),
			Namespace: req.Namespace,
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
