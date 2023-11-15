package telemetry

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/tracer_mock.go . ITracer
type ITracer interface {
	Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span)
}

//go:generate moq -pkg fake -skip-ensure -out ./fake/tracerfactory_mock.go . TracerFactory
type TracerFactory interface {
	GetTracer(name string) ITracer
}
