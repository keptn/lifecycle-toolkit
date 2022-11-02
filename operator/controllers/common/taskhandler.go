package common

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/keptn/lifecycle-controller/operator/api/v1alpha1/common"
	apicommon "github.com/keptn/lifecycle-controller/operator/api/v1alpha1/common"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	klcv1alpha1 "github.com/keptn/lifecycle-controller/operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type TaskHandler struct {
	client.Client
	Recorder    record.EventRecorder
	Log         logr.Logger
	SpanHandler SpanHandler
	Tracer      trace.Tracer
	Scheme      *runtime.Scheme
}

type CreateTaskFunc func(ctx context.Context, namespace string, appVersion client.Object, taskDefinition string, checkType common.CheckType) (string, error)

func (r TaskHandler) ReconcileTasks(ctx context.Context, checkType apicommon.CheckType, reconcileObject client.Object, app bool) ([]klcv1alpha1.TaskStatus, apicommon.StatusSummary, error) {
	piWrapper, err := NewPhaseItemWrapperFromClientObject(reconcileObject)
	if err != nil {
		return nil, apicommon.StatusSummary{}, err
	}

	phase := apicommon.KeptnPhaseType{
		ShortName: "ReconcileTasks",
		LongName:  "Reconcile Tasks",
	}

	var tasks []string
	var statuses []klcv1alpha1.TaskStatus

	switch checkType {
	case apicommon.PreDeploymentCheckType:
		tasks = piWrapper.GetPreDeploymentTasks()
		statuses = piWrapper.GetPreDeploymentTaskStatus()
	case apicommon.PostDeploymentCheckType:
		tasks = piWrapper.GetPreDeploymentTasks()
		statuses = piWrapper.GetPreDeploymentTaskStatus()
	}

	var summary apicommon.StatusSummary
	summary.Total = len(tasks)
	// Check current state of the PrePostDeploymentTasks
	var newStatus []klcv1alpha1.TaskStatus
	for _, taskDefinitionName := range tasks {
		var oldstatus apicommon.KeptnState
		for _, ts := range statuses {
			if ts.TaskDefinitionName == taskDefinitionName {
				oldstatus = ts.Status
			}
		}

		taskStatus := GetTaskStatus(taskDefinitionName, statuses)
		task := &klcv1alpha1.KeptnTask{}
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
			var taskName string
			if app {
				taskName, err = r.CreateKeptnTaskFromApp(ctx, piWrapper.GetNamespace(), reconcileObject, taskDefinitionName, checkType)
				if err != nil {
					return nil, summary, err
				}
			} else {
				taskName, err = r.CreateKeptnTaskFromWorkload(ctx, piWrapper.GetNamespace(), reconcileObject, taskDefinitionName, checkType)
				if err != nil {
					return nil, summary, err
				}
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
		summary = apicommon.UpdateStatusSummary(ns.Status, summary)
	}
	if apicommon.GetOverallState(summary) != apicommon.StateSucceeded {
		RecordEvent(r.Recorder, phase, "Warning", reconcileObject, "NotFinished", "has not finished", piWrapper.GetVersion())
	}
	return newStatus, summary, nil
}

func (r TaskHandler) CreateKeptnTaskFromApp(ctx context.Context, namespace string, appVersion client.Object, taskDefinition string, checkType common.CheckType) (string, error) {
	piWrapper, err := NewPhaseItemWrapperFromClientObject(appVersion)
	if err != nil {
		return "", err
	}

	ctx, span := r.Tracer.Start(ctx, "create_app_task", trace.WithSpanKind(trace.SpanKindProducer))
	defer span.End()

	//semconv.AddAttributeFromAppVersion(span, piWrapper)
	span.SetAttributes(common.AppName.String(piWrapper.GetParentName()))
	span.SetAttributes(common.AppVersion.String(piWrapper.GetVersion()))
	span.SetAttributes(common.WorkloadVersion.String(piWrapper.GetVersion()))

	// create TraceContext
	// follow up with a Keptn propagator that JSON-encoded the OTel map into our own key
	traceContextCarrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, traceContextCarrier)

	phase := apicommon.KeptnPhaseType{
		ShortName: "KeptnTaskCreate",
		LongName:  "Keptn Task Create",
	}

	newTask := &klcv1alpha1.KeptnTask{
		ObjectMeta: metav1.ObjectMeta{
			Name:        apicommon.GenerateTaskName(checkType, taskDefinition),
			Namespace:   namespace,
			Annotations: traceContextCarrier,
		},
		Spec: klcv1alpha1.KeptnTaskSpec{
			AppVersion:       piWrapper.GetVersion(),
			AppName:          piWrapper.GetParentName(),
			TaskDefinition:   taskDefinition,
			Parameters:       klcv1alpha1.TaskParameters{},
			SecureParameters: klcv1alpha1.SecureParameters{},
			Type:             checkType,
		},
	}
	err = controllerutil.SetControllerReference(appVersion, newTask, r.Scheme)
	if err != nil {
		r.Log.Error(err, "could not set controller reference:")
	}
	err = r.Client.Create(ctx, newTask)
	if err != nil {
		r.Log.Error(err, "could not create KeptnTask")
		RecordEvent(r.Recorder, phase, "Warning", appVersion, "CreateFailed", "could not create KeptnTask", piWrapper.GetVersion())
		return "", err
	}
	RecordEvent(r.Recorder, phase, "Normal", appVersion, "Created", "created", piWrapper.GetVersion())

	return newTask.Name, nil
}

func (r TaskHandler) CreateKeptnTaskFromWorkload(ctx context.Context, namespace string, workloadInstance client.Object, taskDefinition string, checkType common.CheckType) (string, error) {
	piWrapper, err := NewPhaseItemWrapperFromClientObject(workloadInstance)
	if err != nil {
		return "", err
	}

	ctx, span := r.Tracer.Start(ctx, fmt.Sprintf("create_%s_deployment_task", checkType), trace.WithSpanKind(trace.SpanKindProducer))
	defer span.End()

	//semconv.AddAttributeFromWorkloadInstance(span, *workloadInstance)
	span.SetAttributes(common.AppName.String(piWrapper.GetAppName()))
	span.SetAttributes(common.WorkloadName.String(piWrapper.GetParentName()))
	span.SetAttributes(common.WorkloadVersion.String(piWrapper.GetVersion()))

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
			AppName:          piWrapper.GetAppName(),
			WorkloadVersion:  piWrapper.GetParentName(),
			Workload:         piWrapper.GetVersion(),
			TaskDefinition:   taskDefinition,
			Parameters:       klcv1alpha1.TaskParameters{},
			SecureParameters: klcv1alpha1.SecureParameters{},
			Type:             checkType,
		},
	}
	err = controllerutil.SetControllerReference(workloadInstance, newTask, r.Scheme)
	if err != nil {
		r.Log.Error(err, "could not set controller reference:")
	}
	err = r.Client.Create(ctx, newTask)
	if err != nil {
		r.Log.Error(err, "could not create KeptnTask")
		RecordEvent(r.Recorder, phase, "Warning", workloadInstance, "CreateFailed", "could not create KeptnTask", piWrapper.GetVersion())
		return "", err
	}
	RecordEvent(r.Recorder, phase, "Normal", workloadInstance, "Created", "created", piWrapper.GetVersion())

	return newTask.Name, nil
}
