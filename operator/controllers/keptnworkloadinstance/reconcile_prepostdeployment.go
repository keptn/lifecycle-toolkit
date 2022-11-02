package keptnworkloadinstance

import (
	"context"
	"fmt"

	klcv1alpha1 "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1"
	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1/common"
	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1/semconv"
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
	newStatus, state, err := r.reconcileTasks(ctx, checkType, workloadInstance)
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

func (r *KeptnWorkloadInstanceReconciler) createKeptnTask(ctx context.Context, namespace string, workloadInstance *klcv1alpha1.KeptnWorkloadInstance, taskDefinition string, checkType common.CheckType) (string, error) {
	ctx, span := r.Tracer.Start(ctx, fmt.Sprintf("create_%s_deployment_task", checkType), trace.WithSpanKind(trace.SpanKindProducer))
	defer span.End()

	semconv.AddAttributeFromWorkloadInstance(span, *workloadInstance)

	// create TraceContext
	// follow up with a Keptn propagator that JSON-encoded the OTel map into our own key
	traceContextCarrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, traceContextCarrier)

	phase := common.KeptnPhaseType{
		ShortName: "KeptnTaskCreate",
		LongName:  "Keptn Task Create",
	}

	newTask := &klcv1alpha1.KeptnTask{
		ObjectMeta: metav1.ObjectMeta{
			Name:        common.GenerateTaskName(checkType, taskDefinition),
			Namespace:   namespace,
			Annotations: traceContextCarrier,
		},
		Spec: klcv1alpha1.KeptnTaskSpec{
			AppName:          workloadInstance.Spec.AppName,
			WorkloadVersion:  workloadInstance.Spec.Version,
			Workload:         workloadInstance.Spec.WorkloadName,
			TaskDefinition:   taskDefinition,
			Parameters:       klcv1alpha1.TaskParameters{},
			SecureParameters: klcv1alpha1.SecureParameters{},
			Type:             checkType,
		},
	}
	err := controllerutil.SetControllerReference(workloadInstance, newTask, r.Scheme)
	if err != nil {
		r.Log.Error(err, "could not set controller reference:")
	}
	err = r.Client.Create(ctx, newTask)
	if err != nil {
		r.Log.Error(err, "could not create KeptnTask")
		controllercommon.RecordEvent(r.Recorder, phase, "Warning", workloadInstance, "CreateFailed", "could not create KeptnTask", workloadInstance.GetVersion())
		return "", err
	}
	controllercommon.RecordEvent(r.Recorder, phase, "Normal", workloadInstance, "Created", "created", workloadInstance.GetVersion())

	return newTask.Name, nil
}

func (r *KeptnWorkloadInstanceReconciler) reconcileTasks(ctx context.Context, checkType common.CheckType, workloadInstance *klcv1alpha1.KeptnWorkloadInstance) ([]klcv1alpha1.TaskStatus, common.StatusSummary, error) {
	var tasks []string
	var statuses []klcv1alpha1.TaskStatus

	phase := common.KeptnPhaseType{
		ShortName: "ReconcileTasks",
		LongName:  "Reconcile Tasks",
	}

	switch checkType {
	case common.PreDeploymentCheckType:
		tasks = workloadInstance.Spec.PreDeploymentTasks
		statuses = workloadInstance.Status.PreDeploymentTaskStatus
	case common.PostDeploymentCheckType:
		tasks = workloadInstance.Spec.PostDeploymentTasks
		statuses = workloadInstance.Status.PostDeploymentTaskStatus
	}

	var summary common.StatusSummary
	summary.Total = len(tasks)
	// Check current state of the PrePostDeploymentTasks
	var newStatus []klcv1alpha1.TaskStatus
	for _, taskDefinitionName := range tasks {
		var oldstatus common.KeptnState

		for _, ts := range statuses {
			if ts.TaskDefinitionName == taskDefinitionName {
				oldstatus = ts.Status
			}
		}

		taskStatus := controllercommon.GetTaskStatus(taskDefinitionName, statuses)
		task := &klcv1alpha1.KeptnTask{}
		taskExists := false

		if oldstatus != taskStatus.Status {
			controllercommon.RecordEvent(r.Recorder, phase, "Normal", workloadInstance, "TaskStatusChanged", fmt.Sprintf("task status changed from %s to %s", oldstatus, taskStatus.Status), workloadInstance.GetVersion())
		}

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
				return nil, summary, err
			}
			taskExists = true
		}

		// Create new Task if it does not exist
		if !taskExists {
			taskName, err := r.createKeptnTask(ctx, workloadInstance.Namespace, workloadInstance, taskDefinitionName, checkType)
			if err != nil {
				return nil, summary, err
			}
			taskStatus.TaskName = taskName
			taskStatus.SetStartTime()
		} else {
			// Update state of Task if it is already created
			taskStatus.Status = task.Status.Status
			if taskStatus.Status.IsCompleted() {
				taskStatus.SetEndTime()
			}
		}
		// Update state of the Check
		newStatus = append(newStatus, taskStatus)
	}

	for _, ns := range newStatus {
		summary = common.UpdateStatusSummary(ns.Status, summary)
	}
	if common.GetOverallState(summary) != common.StateSucceeded {
		controllercommon.RecordEvent(r.Recorder, phase, "Warning", workloadInstance, "NotFinished", "tasks have not finished", workloadInstance.GetVersion())
	}
	return newStatus, summary, nil
}
