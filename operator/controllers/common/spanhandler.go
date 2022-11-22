package common

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"sync"

	"go.opentelemetry.io/otel/trace"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/spanhandler_mock.go . ISpanHandler
type ISpanHandler interface {
	GetSpan(ctx context.Context, tracer trace.Tracer, reconcileObject client.Object, phase string) (context.Context, trace.Span, error)
	UnbindSpan(reconcileObject client.Object, phase string) error
}

type spanContext struct {
	Span trace.Span
	Ctx  context.Context
}

type SpanHandler struct {
	bindCRDSpan map[string]spanContext
	mtx         sync.Mutex
}

func (r *SpanHandler) GetSpan(ctx context.Context, tracer trace.Tracer, reconcileObject client.Object, phase string) (context.Context, trace.Span, error) {
	piWrapper, err := NewPhaseItemWrapperFromClientObject(reconcileObject)
	if err != nil {
		return nil, nil, err
	}
	spanKey := piWrapper.GetSpanKey(phase)
	r.mtx.Lock()
	defer r.mtx.Unlock()
	if r.bindCRDSpan == nil {
		r.bindCRDSpan = make(map[string]spanContext)
	}
	if span, ok := r.bindCRDSpan[spanKey]; ok {
		return span.Ctx, span.Span, nil
	}
	spanName := piWrapper.GetSpanName(phase)
	childCtx, span := tracer.Start(ctx, spanName, trace.WithSpanKind(trace.SpanKindConsumer))
	piWrapper.SetSpanAttributes(span)

	traceContextCarrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(childCtx, traceContextCarrier)
	piWrapper.SetPhaseTraceID(phase, traceContextCarrier)

	r.bindCRDSpan[spanKey] = spanContext{
		Span: span,
		Ctx:  childCtx,
	}
	return childCtx, span, nil
}

func (r *SpanHandler) UnbindSpan(reconcileObject client.Object, phase string) error {
	piWrapper, err := NewPhaseItemWrapperFromClientObject(reconcileObject)
	if err != nil {
		return err
	}
	r.mtx.Lock()
	defer r.mtx.Unlock()
	delete(r.bindCRDSpan, piWrapper.GetSpanKey(phase))
	return nil
}
