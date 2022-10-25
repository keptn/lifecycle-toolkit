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
	bindCRDSpan map[string]trace.Span
}

type PhaseResult struct {
	Continue bool
	ctrl.Result
}

func RecordEvent(recorder record.EventRecorder, phase common.KeptnPhaseType, eventType string, appVersion client.Object, shortReason string, longReason string, version string) {
	recorder.Event(appVersion, eventType, fmt.Sprintf("%s%s", phase.ShortName, shortReason), fmt.Sprintf("%s %s / Namespace: %s, Name: %s, Version: %s ", phase.LongName, longReason, appVersion.GetNamespace(), appVersion.GetName(), version))
}

func (r PhaseHandler) HandlePhase(ctx context.Context, ctxAppTrace context.Context, tracer trace.Tracer, appVersion client.Object, phase common.KeptnPhaseType, span trace.Span, reconcilePhase func() (common.KeptnState, error)) (*PhaseResult, error) {
	piWrapper, err := NewPhaseItemWrapperFromClientObject(appVersion)
	if err != nil {
		return &PhaseResult{Continue: false, Result: ctrl.Result{}}, err
	}
	oldStatus := piWrapper.GetState()
	oldPhase := piWrapper.GetCurrentPhase()
	piWrapper.SetCurrentPhase(phase.ShortName)

	r.Log.Info(phase.LongName + " not finished")
	ctxAppTrace, spanAppTrace := r.GetSpan(ctxAppTrace, tracer, appVersion, phase.ShortName)

	state, err := reconcilePhase()
	if err != nil {
		spanAppTrace.AddEvent(phase.LongName + " could not get reconciled")
		RecordEvent(r.Recorder, phase, "Warning", appVersion, "ReconcileErrored", "could not get reconciled", piWrapper.GetVersion())
		span.SetStatus(codes.Error, err.Error())
		return &PhaseResult{Continue: false, Result: ctrl.Result{}}, err
	}

	defer func(oldStatus common.KeptnState, oldPhase string, appVersion client.Object) {
		piWrapper, _ := NewPhaseItemWrapperFromClientObject(appVersion)
		if oldStatus != piWrapper.GetState() || oldPhase != piWrapper.GetCurrentPhase() {
			ctx, spanAppTrace = r.GetSpan(ctxAppTrace, tracer, appVersion, piWrapper.GetCurrentPhase())
			if err := r.Status().Update(ctx, appVersion); err != nil {
				r.Log.Error(err, "could not update status")
			}
		}
	}(oldStatus, oldPhase, appVersion)

	if state.IsCompleted() {
		if state.IsFailed() {
			piWrapper.Complete()
			piWrapper.SetState(common.StateFailed)
			spanAppTrace.AddEvent(phase.LongName + " has failed")
			spanAppTrace.SetStatus(codes.Error, "Failed")
			spanAppTrace.End()
			r.UnbindSpan(appVersion, phase.ShortName)
			RecordEvent(r.Recorder, phase, "Warning", appVersion, "Failed", "has failed", piWrapper.GetVersion())
			return &PhaseResult{Continue: false, Result: ctrl.Result{}}, nil
		}

		piWrapper.SetState(common.StateSucceeded)
		spanAppTrace.AddEvent(phase.LongName + " has succeeded")
		spanAppTrace.SetStatus(codes.Ok, "Succeeded")
		spanAppTrace.End()
		r.UnbindSpan(appVersion, phase.ShortName)
		RecordEvent(r.Recorder, phase, "Normal", appVersion, "Succeeded", "has succeeded", piWrapper.GetVersion())

		return &PhaseResult{Continue: true, Result: ctrl.Result{}}, nil
	}

	piWrapper.SetState(common.StateProgressing)
	RecordEvent(r.Recorder, phase, "Warning", appVersion, "NotFinished", "has not finished", piWrapper.GetVersion())

	return &PhaseResult{Continue: false, Result: ctrl.Result{Requeue: true, RequeueAfter: 5 * time.Second}}, nil
}

func (r PhaseHandler) GetSpan(ctx context.Context, tracer trace.Tracer, appv client.Object, phase string) (context.Context, trace.Span) {
	piWrapper, err := NewPhaseItemWrapperFromClientObject(appv)
	if err != nil {
		return nil, nil
	}
	appvName := piWrapper.GetSpanName(phase)
	if r.bindCRDSpan == nil {
		r.bindCRDSpan = make(map[string]trace.Span)
	}
	if span, ok := r.bindCRDSpan[appvName]; ok {
		return ctx, span
	}
	ctx, span := tracer.Start(ctx, phase, trace.WithSpanKind(trace.SpanKindConsumer))
	r.Log.Info("DEBUG: Created span " + appvName)
	r.bindCRDSpan[appvName] = span
	return ctx, span
}

func (r PhaseHandler) UnbindSpan(appv client.Object, phase string) {
	piWrapper, err := NewPhaseItemWrapperFromClientObject(appv)
	if err != nil {
		return
	}
	delete(r.bindCRDSpan, piWrapper.GetSpanName(phase))
}
