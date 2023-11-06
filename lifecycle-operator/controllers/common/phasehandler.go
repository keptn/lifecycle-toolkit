package common

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/telemetry"
	controllererrors "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/errors"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/lifecycle/interfaces"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PhaseHandler struct {
	client.Client
	EventSender IEvent
	Log         logr.Logger
	SpanHandler telemetry.ISpanHandler
}

type PhaseResult struct {
	Continue bool
	ctrl.Result
}

func (r PhaseHandler) HandlePhase(ctx context.Context, ctxTrace context.Context, tracer trace.Tracer, reconcileObject client.Object, phase apicommon.KeptnPhaseType, reconcilePhase func(phaseCtx context.Context) (apicommon.KeptnState, error)) (*PhaseResult, error) {
	requeueResult := ctrl.Result{Requeue: true, RequeueAfter: 5 * time.Second}
	piWrapper, err := interfaces.NewPhaseItemWrapperFromClientObject(reconcileObject)
	if err != nil {
		return &PhaseResult{Continue: false, Result: ctrl.Result{Requeue: true}}, err
	}
	oldStatus := piWrapper.GetState()
	oldPhase := piWrapper.GetCurrentPhase()
	// do not attempt to execute the current phase if the whole phase item is already in deprecated/failed state
	if shouldAbortPhase(oldStatus) {
		return &PhaseResult{Continue: false, Result: ctrl.Result{}}, nil
	}
	if oldPhase != phase.ShortName {
		r.EventSender.Emit(phase, "Normal", reconcileObject, apicommon.PhaseStateStarted, "has started", piWrapper.GetVersion())
		piWrapper.SetCurrentPhase(phase.ShortName)
	}

	spanPhaseCtx, spanPhaseTrace, err := r.SpanHandler.GetSpan(ctxTrace, tracer, reconcileObject, phase.ShortName)
	if err != nil {
		r.Log.Error(err, "could not get span")
	}

	state, err := reconcilePhase(spanPhaseCtx)
	if err != nil {
		spanPhaseTrace.AddEvent(phase.LongName + " could not get reconciled")
		r.EventSender.Emit(phase, "Warning", reconcileObject, apicommon.PhaseStateReconcileError, "could not get reconciled", piWrapper.GetVersion())
		return &PhaseResult{Continue: false, Result: requeueResult}, err
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

	return &PhaseResult{Continue: false, Result: requeueResult}, nil
}

func shouldAbortPhase(oldStatus apicommon.KeptnState) bool {
	return oldStatus.IsDeprecated() || oldStatus.IsFailed()
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
		r.EventSender.Emit(phase, "Warning", reconcileObject, apicommon.PhaseStateFailed, "has failed", piWrapper.GetVersion())
		piWrapper.DeprecateRemainingPhases(phase)
		return &PhaseResult{Continue: false, Result: ctrl.Result{}}, nil
	}

	// end the current phase do not set the overall state of the whole object to Succeeded here, as this can cause
	// premature progression of reconcile objects that depend on the completion of another
	spanPhaseTrace.AddEvent(phase.LongName + " has succeeded")
	spanPhaseTrace.SetStatus(codes.Ok, "Succeeded")
	spanPhaseTrace.End()
	if err := r.SpanHandler.UnbindSpan(reconcileObject, phase.ShortName); err != nil {
		r.Log.Error(err, controllererrors.ErrCouldNotUnbindSpan, reconcileObject.GetName())
	}
	r.EventSender.Emit(phase, "Normal", reconcileObject, apicommon.PhaseStateFinished, "has finished", piWrapper.GetVersion())

	return &PhaseResult{Continue: true, Result: ctrl.Result{Requeue: true, RequeueAfter: 5 * time.Second}}, nil
}
