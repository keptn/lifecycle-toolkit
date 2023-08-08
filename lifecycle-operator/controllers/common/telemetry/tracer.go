package telemetry

import "go.opentelemetry.io/otel/trace"

//go:generate moq -pkg fake -skip-ensure -out ./fake/tracer_mock.go . ITracer
type ITracer = trace.Tracer

//go:generate moq -pkg fake -skip-ensure -out ./fake/tracerfactory_mock.go . TracerFactory
type TracerFactory interface {
	GetTracer(name string) ITracer
}
