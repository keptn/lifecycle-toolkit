package common

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha2/common"
	controllererrors "github.com/keptn/lifecycle-toolkit/operator/controllers/errors"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/interfaces"
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
	SpanHandler ISpanHandler
}

type PhaseResult struct {
	Continue bool
	ctrl.Result
}

func RecordEvent(recorder record.EventRecorder, phase apicommon.KeptnPhaseType, eventType string, reconcileObject client.Object, shortReason string, longReason string, version string) {
	recorder.Event(reconcileObject, eventType, fmt.Sprintf("%s%s", phase.ShortName, shortReason), fmt.Sprintf("%s %s / Namespace: %s, Name: %s, Version: %s ", phase.LongName, longReason, reconcileObject.GetNamespace(), reconcileObject.GetName(), version))
}

func (r PhaseHandler) HandlePhase(ctx context.Context, ctxTrace context.Context, tracer trace.Tracer, reconcileObject client.Object, phase apicommon.KeptnPhaseType, span trace.Span, reconcilePhase func(phaseCtx context.Context) (apicommon.KeptnState, error)) (*PhaseResult, error) {
	requeueResult := ctrl.Result{Requeue: true, RequeueAfter: 5 * time.Second}
	piWrapper, err := interfaces.NewPhaseItemWrapperFromClientObject(reconcileObject)
	if err != nil {
		return &PhaseResult{Continue: false, Result: ctrl.Result{Requeue: true}}, err
	}
	oldStatus := piWrapper.GetState()
	oldPhase := piWrapper.GetCurrentPhase()
	if oldStatus.IsDeprecated() {
		return &PhaseResult{Continue: false, Result: ctrl.Result{}}, nil
	}
	piWrapper.SetCurrentPhase(phase.ShortName)

	r.Log.Info(phase.LongName + " not finished")
	spanPhaseCtx, spanPhaseTrace, err := r.SpanHandler.GetSpan(ctxTrace, tracer, reconcileObject, phase.ShortName)
	if err != nil {
		r.Log.Error(err, "could not get span")
	}

	state, err := reconcilePhase(spanPhaseCtx)
	if err != nil {
		spanPhaseTrace.AddEvent(phase.LongName + " could not get reconciled")
		RecordEvent(r.Recorder, phase, "Warning", reconcileObject, "ReconcileErrored", "could not get reconciled", piWrapper.GetVersion())
		span.SetStatus(codes.Error, err.Error())
		return &PhaseResult{Continue: false, Result: requeueResult}, err
	}

	if state.IsPending() {
		state = apicommon.StateProgressing
	}

	defer func(ctx context.Context, oldStatus apicommon.KeptnState, oldPhase string, reconcileObject client.Object) {
		piWrapper, _ := interfaces.NewPhaseItemWrapperFromClientObject(reconcileObject)
		if oldStatus != piWrapper.GetState() || oldPhase != piWrapper.GetCurrentPhase() {
			if err := r.Status().Update(ctx, reconcileObject); err != nil {
				r.Log.Error(err, "could not update status")
			}
		}
	}(ctx, oldStatus, oldPhase, reconcileObject)

	if state.IsCompleted() {
		return r.handleCompletedPhase(state, piWrapper, phase, reconcileObject, spanPhaseTrace)
	}

	piWrapper.SetState(apicommon.StateProgressing)
	RecordEvent(r.Recorder, phase, "Warning", reconcileObject, "NotFinished", "has not finished", piWrapper.GetVersion())

	return &PhaseResult{Continue: false, Result: requeueResult}, nil
}

func (r PhaseHandler) handleCompletedPhase(state apicommon.KeptnState, piWrapper *interfaces.PhaseItemWrapper, phase apicommon.KeptnPhaseType, reconcileObject client.Object, spanPhaseTrace trace.Span) (*PhaseResult, error) {
	if state.IsFailed() {
		piWrapper.Complete()
		piWrapper.SetState(apicommon.StateFailed)
		spanPhaseTrace.AddEvent(phase.LongName + " has failed")
		spanPhaseTrace.SetStatus(codes.Error, "Failed")
		spanPhaseTrace.End()
		if err := r.SpanHandler.UnbindSpan(reconcileObject, phase.ShortName); err != nil {
			r.Log.Error(err, controllererrors.ErrCouldNotUnbindSpan, reconcileObject.GetName())
		}
		RecordEvent(r.Recorder, phase, "Warning", reconcileObject, "Failed", "has failed", piWrapper.GetVersion())
		piWrapper.DeprecateRemainingPhases(phase)
		return &PhaseResult{Continue: false, Result: ctrl.Result{}}, nil
	}

	piWrapper.SetState(apicommon.StateSucceeded)
	spanPhaseTrace.AddEvent(phase.LongName + " has succeeded")
	spanPhaseTrace.SetStatus(codes.Ok, "Succeeded")
	spanPhaseTrace.End()
	if err := r.SpanHandler.UnbindSpan(reconcileObject, phase.ShortName); err != nil {
		r.Log.Error(err, controllererrors.ErrCouldNotUnbindSpan, reconcileObject.GetName())
	}
	RecordEvent(r.Recorder, phase, "Normal", reconcileObject, "Succeeded", "has succeeded", piWrapper.GetVersion())

	return &PhaseResult{Continue: true, Result: ctrl.Result{Requeue: true, RequeueAfter: 5 * time.Second}}, nil
}
