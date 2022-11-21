package common

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	klcv1alpha1 "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1"
	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1/common"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PhaseHandler struct {
	client.Client
	Recorder    record.EventRecorder
	Log         logr.Logger
	SpanHandler ISpanHandler
}

type PhaseResult struct {
	Continue bool
	ctrl.Result
}

type failedCheckReason struct {
	Message string
	Time    time.Time
}

func RecordEvent(recorder record.EventRecorder, phase common.KeptnPhaseType, eventType string, reconcileObject client.Object, shortReason string, longReason string, version string) {
	recorder.Event(reconcileObject, eventType, fmt.Sprintf("%s%s", phase.ShortName, shortReason), fmt.Sprintf("%s %s / Namespace: %s, Name: %s, Version: %s ", phase.LongName, longReason, reconcileObject.GetNamespace(), reconcileObject.GetName(), version))
}

func (r PhaseHandler) HandlePhase(ctx context.Context, ctxTrace context.Context, tracer trace.Tracer, reconcileObject client.Object, phase common.KeptnPhaseType, span trace.Span, reconcilePhase func() (common.KeptnState, error)) (*PhaseResult, error) {
	requeueResult := ctrl.Result{Requeue: true, RequeueAfter: 5 * time.Second}
	piWrapper, err := NewPhaseItemWrapperFromClientObject(reconcileObject)
	if err != nil {
		return &PhaseResult{Continue: false, Result: ctrl.Result{Requeue: true}}, err
	}
	oldStatus := piWrapper.GetState()
	oldPhase := piWrapper.GetCurrentPhase()
	if oldStatus.IsCancelled() {
		return &PhaseResult{Continue: false, Result: ctrl.Result{}}, nil
	}
	piWrapper.SetCurrentPhase(phase.ShortName)

	r.Log.Info(phase.LongName + " not finished")
	_, spanAppTrace, err := r.SpanHandler.GetSpan(ctxTrace, tracer, reconcileObject, phase.ShortName)
	if err != nil {
		r.Log.Error(err, "could not get span")
	}

	state, err := reconcilePhase()
	if err != nil {
		spanAppTrace.AddEvent(phase.LongName + " could not get reconciled")
		RecordEvent(r.Recorder, phase, "Warning", reconcileObject, "ReconcileErrored", "could not get reconciled", piWrapper.GetVersion())
		span.SetStatus(codes.Error, err.Error())
		return &PhaseResult{Continue: false, Result: requeueResult}, err
	}

	if state.IsPending() {
		state = common.StateProgressing
	}

	defer func(oldStatus common.KeptnState, oldPhase string, reconcileObject client.Object) {
		piWrapper, _ := NewPhaseItemWrapperFromClientObject(reconcileObject)
		if oldStatus != piWrapper.GetState() || oldPhase != piWrapper.GetCurrentPhase() {
			ctx, spanAppTrace, err = r.SpanHandler.GetSpan(ctxTrace, tracer, reconcileObject, piWrapper.GetCurrentPhase())
			if err != nil {
				r.Log.Error(err, "could not get span")
			}
			if err := r.Status().Update(ctx, reconcileObject); err != nil {
				r.Log.Error(err, "could not update status")
			}
		}
	}(oldStatus, oldPhase, reconcileObject)

	if state.IsCompleted() {
		if state.IsFailed() {
			piWrapper.Complete()
			piWrapper.SetState(common.StateFailed)
			spanAppTrace.AddEvent(phase.LongName + " has failed")
			if err := r.createFailureReasonSpanEvents(ctx, phase, piWrapper, spanAppTrace); err != nil {
				r.Log.Error(err, "cannot add failure events to spans")
			}
			spanAppTrace.SetStatus(codes.Error, "Failed")
			spanAppTrace.End()
			if err := r.SpanHandler.UnbindSpan(reconcileObject, phase.ShortName); err != nil {
				r.Log.Error(err, "cannot unbind span")
			}
			RecordEvent(r.Recorder, phase, "Warning", reconcileObject, "Failed", "has failed", piWrapper.GetVersion())
			piWrapper.CancelRemainingPhases(phase)
			return &PhaseResult{Continue: false, Result: ctrl.Result{}}, nil
		}

		piWrapper.SetState(common.StateSucceeded)
		spanAppTrace.AddEvent(phase.LongName + " has succeeded")
		spanAppTrace.SetStatus(codes.Ok, "Succeeded")
		spanAppTrace.End()
		if err := r.SpanHandler.UnbindSpan(reconcileObject, phase.ShortName); err != nil {
			r.Log.Error(err, "cannot unbind span")
		}
		RecordEvent(r.Recorder, phase, "Normal", reconcileObject, "Succeeded", "has succeeded", piWrapper.GetVersion())

		return &PhaseResult{Continue: true, Result: requeueResult}, nil
	}

	piWrapper.SetState(common.StateProgressing)
	RecordEvent(r.Recorder, phase, "Warning", reconcileObject, "NotFinished", "has not finished", piWrapper.GetVersion())

	return &PhaseResult{Continue: false, Result: requeueResult}, nil
}

func (r PhaseHandler) createFailureReasonSpanEvents(ctx context.Context, phase common.KeptnPhaseType, object *PhaseItemWrapper, spanTrace trace.Span) error {
	var messageObjects []failedCheckReason
	var err error
	if phase.IsEvaluation() {
		messageObjects, err = r.GetEvaluationFailureReasons(ctx, phase, object)
	} else if phase.IsTask() {
		messageObjects, err = r.GetTaskFailureReasons(ctx, phase, object)
	}

	if err != nil {
		return err
	}

	for _, msgObject := range messageObjects {
		spanTrace.AddEvent(msgObject.Message, trace.WithTimestamp(msgObject.Time))
	}

	return nil
}

func (r PhaseHandler) GetEvaluationFailureReasons(ctx context.Context, phase common.KeptnPhaseType, object PhaseItem) ([]failedCheckReason, error) {
	resultEvents := []failedCheckReason{}
	var status []klcv1alpha1.EvaluationStatus
	if phase.IsPreEvaluation() {
		status = object.GetPreDeploymentEvaluationTaskStatus()
	} else {
		status = object.GetPostDeploymentEvaluationTaskStatus()
	}

	// there can be only one evaluation and in this section of the code, it can only be failed
	// checking length of the status only for safety reasons
	if len(status) != 1 {
		return nil, fmt.Errorf("evaluation status not found for %s/%s", object.GetAppName(), object.GetParentName())
	}

	evaluation := &klcv1alpha1.KeptnEvaluation{}
	if err := r.Client.Get(ctx, types.NamespacedName{Name: status[0].EvaluationName, Namespace: object.GetNamespace()}, evaluation); err != nil {
		return nil, fmt.Errorf("evaluation %s not found for %s/%s", status[0].EvaluationName, object.GetAppName(), object.GetParentName())
	}

	for k, v := range evaluation.Status.EvaluationStatus {
		if v.Status == common.StateFailed {
			obj := failedCheckReason{
				Time:    evaluation.Status.EndTime.Time,
				Message: fmt.Sprintf("evaluation of '%s' failed with value: '%s' and reason: '%s'", k, v.Value, v.Message),
			}
			resultEvents = append(resultEvents, obj)
		}
	}

	return resultEvents, nil
}

func (r PhaseHandler) GetTaskFailureReasons(ctx context.Context, phase common.KeptnPhaseType, object PhaseItem) ([]failedCheckReason, error) {
	resultEvents := []failedCheckReason{}
	var failedTasks []klcv1alpha1.KeptnTask
	var status []klcv1alpha1.TaskStatus
	if phase.IsPreTask() {
		status = object.GetPreDeploymentTaskStatus()
	} else {
		status = object.GetPostDeploymentTaskStatus()
	}

	for _, item := range status {
		if item.Status == common.StateFailed {
			task := &klcv1alpha1.KeptnTask{}
			if err := r.Client.Get(ctx, types.NamespacedName{Name: item.TaskName, Namespace: object.GetNamespace()}, task); err != nil {
				return nil, fmt.Errorf("task %s not found for %s/%s", item.TaskName, object.GetAppName(), object.GetParentName())
			}
			failedTasks = append(failedTasks, *task)
		}
	}

	for _, task := range failedTasks {
		obj := failedCheckReason{
			Time:    task.Status.EndTime.Time,
			Message: fmt.Sprintf("task '%s' failed with reason: '%s'", task.Name, task.Status.Message),
		}
		resultEvents = append(resultEvents, obj)
	}

	return resultEvents, nil
}
