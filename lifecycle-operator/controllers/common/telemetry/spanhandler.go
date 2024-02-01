package telemetry

import (
	"context"
	"sync"

	keptncontext "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/context"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/lifecycle/interfaces"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/spanhandler_mock.go . ISpanHandler
type ISpanHandler interface {
	GetSpan(ctx context.Context, tracer ITracer, reconcileObject client.Object, phase string, links ...trace.Link) (context.Context, trace.Span, error)
	UnbindSpan(reconcileObject client.Object, phase string) error
}

type keptnSpanCtx struct {
	Span trace.Span
	Ctx  context.Context //nolint:all
}

type Handler struct {
	bindCRDSpan map[string]keptnSpanCtx
	mtx         sync.Mutex
}

func (r *Handler) GetSpan(ctx context.Context, tracer ITracer, reconcileObject client.Object, phase string, links ...trace.Link) (context.Context, trace.Span, error) {
	piWrapper, err := interfaces.NewSpanItemWrapperFromClientObject(reconcileObject)
	if err != nil {
		return nil, nil, err
	}
	spanKey := piWrapper.GetSpanKey(phase)
	r.mtx.Lock()
	defer r.mtx.Unlock()
	if r.bindCRDSpan == nil {
		r.bindCRDSpan = make(map[string]keptnSpanCtx)
	}
	if span, ok := r.bindCRDSpan[spanKey]; ok {
		return span.Ctx, span.Span, nil
	}
	spanName := piWrapper.GetSpanName(phase)
	childCtx, span := tracer.Start(
		ctx,
		spanName,
		trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithLinks(links...),
	)
	piWrapper.SetSpanAttributes(span)

	// also get attributes from context
	if meta, ok := keptncontext.GetAppMetadataFromContext(ctx); ok {
		for key, value := range meta {
			span.SetAttributes(attribute.KeyValue{
				Key:   attribute.Key(key),
				Value: attribute.StringValue(value),
			})
		}
	}

	if phase != "" {
		traceContextCarrier := propagation.MapCarrier{}
		otel.GetTextMapPropagator().Inject(childCtx, traceContextCarrier)
		piWrapper.SetPhaseTraceID(phase, traceContextCarrier)
	}

	r.bindCRDSpan[spanKey] = keptnSpanCtx{
		Span: span,
		Ctx:  childCtx,
	}
	return childCtx, span, nil
}

func (r *Handler) UnbindSpan(reconcileObject client.Object, phase string) error {
	piWrapper, err := interfaces.NewSpanItemWrapperFromClientObject(reconcileObject)
	if err != nil {
		return err
	}
	r.mtx.Lock()
	defer r.mtx.Unlock()
	delete(r.bindCRDSpan, piWrapper.GetSpanKey(phase))
	return nil
}
