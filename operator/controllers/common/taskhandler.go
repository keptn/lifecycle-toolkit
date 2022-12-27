package common

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	klcv1alpha2 "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha2"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha2/common"
	controllererrors "github.com/keptn/lifecycle-toolkit/operator/controllers/errors"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/interfaces"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type TaskHandler struct {
	client.Client
	Recorder    record.EventRecorder
	Log         logr.Logger
	Tracer      trace.Tracer
	Scheme      *runtime.Scheme
	SpanHandler ISpanHandler
}

type TaskCreateAttributes struct {
	SpanName       string
	TaskDefinition string
	CheckType      apicommon.CheckType
}

func (r TaskHandler) ReconcileTasks(ctx context.Context, phaseCtx context.Context, reconcileObject client.Object, taskCreateAttributes TaskCreateAttributes) ([]klcv1alpha2.TaskStatus, apicommon.StatusSummary, error) {
	piWrapper, err := interfaces.NewPhaseItemWrapperFromClientObject(reconcileObject)
	if err != nil {
		return nil, apicommon.StatusSummary{}, err
	}

	phase := apicommon.KeptnPhaseType{
		ShortName: "ReconcileTasks",
		LongName:  "Reconcile Tasks",
	}

	tasks, statuses := r.setupTasks(taskCreateAttributes, piWrapper)

	var summary apicommon.StatusSummary
	summary.Total = len(tasks)
	// Check current state of the PrePostDeploymentTasks
	var newStatus []klcv1alpha2.TaskStatus
	for _, taskDefinitionName := range tasks {
		oldstatus := getOldStatusOfTasks(statuses, taskDefinitionName)

		taskStatus := GetTaskStatus(taskDefinitionName, statuses)
		task := &klcv1alpha2.KeptnTask{}
		taskExists := false

		if oldstatus != taskStatus.Status {
			RecordEvent(r.Recorder, phase, "Normal", reconcileObject, "TaskStatusChanged", fmt.Sprintf("task status changed from %s to %s", oldstatus, taskStatus.Status), piWrapper.GetVersion())
		}

		// Check if task has already succeeded or failed
		if taskStatus.Status == apicommon.StateSucceeded || taskStatus.Status == apicommon.StateFailed {
			newStatus = append(newStatus, taskStatus)
			continue
		}

		// Check if Task is already created
		if taskStatus.TaskName != "" {
			err := r.Client.Get(ctx, types.NamespacedName{Name: taskStatus.TaskName, Namespace: piWrapper.GetNamespace()}, task)
			if err != nil && errors.IsNotFound(err) {
				taskStatus.TaskName = ""
			} else if err != nil {
				return nil, summary, err
			}
			taskExists = true
		}

		// Create new Task if it does not exist
		if !taskExists {
			taskCreateAttributes.TaskDefinition = taskDefinitionName
			taskName, err := r.CreateKeptnTask(ctx, piWrapper.GetNamespace(), reconcileObject, taskCreateAttributes)
			if err != nil {
				return nil, summary, err
			}
			taskStatus.TaskName = taskName
			taskStatus.SetStartTime()
			_, _, err = r.SpanHandler.GetSpan(phaseCtx, r.Tracer, task, "")
			if err != nil {
				r.Log.Error(err, "could not get span")
			}
		} else {
			_, spanTaskTrace, err := r.SpanHandler.GetSpan(phaseCtx, r.Tracer, task, "")
			if err != nil {
				r.Log.Error(err, "could not get span")
			}
			// Update state of Task if it is already created
			taskStatus.Status = task.Status.Status
			if taskStatus.Status.IsCompleted() {
				if taskStatus.Status.IsSucceeded() {
					spanTaskTrace.AddEvent(task.Name + " has finished")
					spanTaskTrace.SetStatus(codes.Ok, "Finished")
				} else {
					spanTaskTrace.AddEvent(task.Name + " has failed")
					r.setTaskFailureEvents(task, spanTaskTrace)
					spanTaskTrace.SetStatus(codes.Error, "Failed")
				}
				spanTaskTrace.End()
				if err := r.SpanHandler.UnbindSpan(task, ""); err != nil {
					r.Log.Error(err, controllererrors.ErrCouldNotUnbindSpan, task.Name)
				}
				taskStatus.SetEndTime()
			}
		}
		// Update state of the Check
		newStatus = append(newStatus, taskStatus)
	}

	for _, ns := range newStatus {
		summary = apicommon.UpdateStatusSummary(ns.Status, summary)
	}
	if apicommon.GetOverallState(summary) != apicommon.StateSucceeded {
		RecordEvent(r.Recorder, phase, "Warning", reconcileObject, "NotFinished", "has not finished", piWrapper.GetVersion())
	}
	return newStatus, summary, nil
}

func (r TaskHandler) CreateKeptnTask(ctx context.Context, namespace string, reconcileObject client.Object, taskCreateAttributes TaskCreateAttributes) (string, error) {
	piWrapper, err := interfaces.NewPhaseItemWrapperFromClientObject(reconcileObject)
	if err != nil {
		return "", err
	}

	phase := apicommon.KeptnPhaseType{
		ShortName: "KeptnTaskCreate",
		LongName:  "Keptn Task Create",
	}

	newTask := piWrapper.GenerateTask(taskCreateAttributes.TaskDefinition, taskCreateAttributes.CheckType)
	err = controllerutil.SetControllerReference(reconcileObject, &newTask, r.Scheme)
	if err != nil {
		r.Log.Error(err, "could not set controller reference:")
	}
	err = r.Client.Create(ctx, &newTask)
	if err != nil {
		r.Log.Error(err, "could not create KeptnTask")
		RecordEvent(r.Recorder, phase, "Warning", reconcileObject, "CreateFailed", "could not create KeptnTask", piWrapper.GetVersion())
		return "", err
	}
	RecordEvent(r.Recorder, phase, "Normal", reconcileObject, "Created", "created", piWrapper.GetVersion())

	return newTask.Name, nil
}

func (r TaskHandler) setTaskFailureEvents(task *klcv1alpha2.KeptnTask, spanTrace trace.Span) {
	spanTrace.AddEvent(fmt.Sprintf("task '%s' failed with reason: '%s'", task.Name, task.Status.Message), trace.WithTimestamp(time.Now().UTC()))
}

func (r TaskHandler) setupTasks(taskCreateAttributes TaskCreateAttributes, piWrapper *interfaces.PhaseItemWrapper) ([]string, []klcv1alpha2.TaskStatus) {
	var tasks []string
	var statuses []klcv1alpha2.TaskStatus

	switch taskCreateAttributes.CheckType {
	case apicommon.PreDeploymentCheckType:
		tasks = piWrapper.GetPreDeploymentTasks()
		statuses = piWrapper.GetPreDeploymentTaskStatus()
	case apicommon.PostDeploymentCheckType:
		tasks = piWrapper.GetPostDeploymentTasks()
		statuses = piWrapper.GetPostDeploymentTaskStatus()
	}
	return tasks, statuses
}

func getOldStatusOfTasks(statuses []klcv1alpha2.TaskStatus, taskDefinitionName string) apicommon.KeptnState {
	var oldstatus apicommon.KeptnState
	for _, ts := range statuses {
		if ts.TaskDefinitionName == taskDefinitionName {
			oldstatus = ts.Status
		}
	}

	return oldstatus
}
