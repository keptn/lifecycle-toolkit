package common

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/keptn/lifecycle-controller/operator/api/v1alpha1/common"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PhaseHandler struct {
	client.Client
	Recorder    record.EventRecorder
	Log         logr.Logger
	SpanHandler SpanHandler
}

type PhaseResult struct {
	Continue bool
	ctrl.Result
}

func RecordEvent(recorder record.EventRecorder, phase common.KeptnPhaseType, eventType string, reconcileObject client.Object, shortReason string, longReason string, version string) {
	recorder.Event(reconcileObject, eventType, fmt.Sprintf("%s%s", phase.ShortName, shortReason), fmt.Sprintf("%s %s / Namespace: %s, Name: %s, Version: %s ", phase.LongName, longReason, reconcileObject.GetNamespace(), reconcileObject.GetName(), version))
}

func (r PhaseHandler) HandlePhase(ctx context.Context, ctxAppTrace context.Context, tracer trace.Tracer, reconcileObject client.Object, phase common.KeptnPhaseType, span trace.Span, reconcilePhase func() (common.KeptnState, error)) (*PhaseResult, error) {
	requeueResult := ctrl.Result{Requeue: true, RequeueAfter: 5 * time.Second}
	piWrapper, err := NewPhaseItemWrapperFromClientObject(reconcileObject)
	if err != nil {
		return &PhaseResult{Continue: false, Result: requeueResult}, err
	}
	oldStatus := piWrapper.GetState()
	oldPhase := piWrapper.GetCurrentPhase()
	piWrapper.SetCurrentPhase(phase.ShortName)

	r.Log.Info(phase.LongName + " not finished")
	ctxAppTrace, spanAppTrace, err := r.SpanHandler.GetSpan(ctxAppTrace, tracer, reconcileObject, phase.ShortName)
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

	defer func(oldStatus common.KeptnState, oldPhase string, reconcileObject client.Object) {
		piWrapper, _ := NewPhaseItemWrapperFromClientObject(reconcileObject)
		if oldStatus != piWrapper.GetState() || oldPhase != piWrapper.GetCurrentPhase() {
			ctx, spanAppTrace, err = r.SpanHandler.GetSpan(ctxAppTrace, tracer, reconcileObject, piWrapper.GetCurrentPhase())
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
			spanAppTrace.SetStatus(codes.Error, "Failed")
			spanAppTrace.End()
			if err := r.SpanHandler.UnbindSpan(reconcileObject, phase.ShortName); err != nil {
				r.Log.Error(err, "cannot unbind span")
			}
			RecordEvent(r.Recorder, phase, "Warning", reconcileObject, "Failed", "has failed", piWrapper.GetVersion())
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
