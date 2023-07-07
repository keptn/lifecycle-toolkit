package common

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3/common"
	controllererrors "github.com/keptn/lifecycle-toolkit/operator/controllers/errors"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/lifecycle/interfaces"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PhaseHandler struct {
	client.Client
	EventSender EventSender
	Log         logr.Logger
	SpanHandler ISpanHandler
}

type PhaseResult struct {
	Continue bool
	ctrl.Result
}

func (r PhaseHandler) HandlePhase(ctx context.Context, ctxTrace context.Context, tracer trace.Tracer, reconcileObject client.Object, phase apicommon.KeptnPhaseType, span trace.Span, reconcilePhase func(phaseCtx context.Context) (apicommon.KeptnState, error)) (*PhaseResult, error) {
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
	piWrapper.SetCurrentPhase(phase.ShortName)

	r.Log.Info(phase.LongName + " not finished")
	spanPhaseCtx, spanPhaseTrace, err := r.SpanHandler.GetSpan(ctxTrace, tracer, reconcileObject, phase.ShortName)
	if err != nil {
		r.Log.Error(err, "could not get span")
	}

	state, err := reconcilePhase(spanPhaseCtx)
	if err != nil {
		spanPhaseTrace.AddEvent(phase.LongName + " could not get reconciled")
		r.EventSender.SendK8sEvent(phase, "Warning", reconcileObject, "ReconcileErrored", "could not get reconciled", piWrapper.GetVersion())
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
	r.EventSender.SendK8sEvent(phase, "Warning", reconcileObject, "NotFinished", "has not finished", piWrapper.GetVersion())

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
		r.EventSender.SendK8sEvent(phase, "Warning", reconcileObject, "Failed", "has failed", piWrapper.GetVersion())
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
	r.EventSender.SendK8sEvent(phase, "Normal", reconcileObject, "Succeeded", "has succeeded", piWrapper.GetVersion())

	return &PhaseResult{Continue: true, Result: ctrl.Result{Requeue: true, RequeueAfter: 5 * time.Second}}, nil
}
