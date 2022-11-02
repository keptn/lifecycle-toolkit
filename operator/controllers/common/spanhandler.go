package common

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/trace"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type SpanHandler struct {
	bindCRDSpan map[string]trace.Span
	mtx         sync.Mutex
}

func (r *SpanHandler) GetSpan(ctx context.Context, tracer trace.Tracer, reconcileObject client.Object, phase string) (context.Context, trace.Span, error) {
	piWrapper, err := NewPhaseItemWrapperFromClientObject(reconcileObject)
	if err != nil {
		return nil, nil, err
	}
	appvName := piWrapper.GetSpanName(phase)
	r.mtx.Lock()
	defer r.mtx.Unlock()
	if r.bindCRDSpan == nil {
		r.bindCRDSpan = make(map[string]trace.Span)
	}
	if span, ok := r.bindCRDSpan[appvName]; ok {
		return ctx, span, nil
	}
	ctx, span := tracer.Start(ctx, phase, trace.WithSpanKind(trace.SpanKindConsumer))
	r.bindCRDSpan[appvName] = span
	return ctx, span, nil
}

func (r *SpanHandler) UnbindSpan(reconcileObject client.Object, phase string) error {
	piWrapper, err := NewPhaseItemWrapperFromClientObject(reconcileObject)
	if err != nil {
		return err
	}
	r.mtx.Lock()
	defer r.mtx.Unlock()
	delete(r.bindCRDSpan, piWrapper.GetSpanName(phase))
	return nil
}
