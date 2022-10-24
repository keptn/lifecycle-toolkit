package common

import "go.opentelemetry.io/otel/trace"

//go:generate moq -pkg fake -skip-ensure -out ./fake/tracer.go . ITracer
type ITracer interface {
	trace.Tracer
}
