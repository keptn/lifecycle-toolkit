package interfaces

import "go.opentelemetry.io/otel/trace"

//go:generate moq -pkg fake -skip-ensure -out ./fake/tracer_mock.go . ITracer
type ITracer = trace.Tracer
