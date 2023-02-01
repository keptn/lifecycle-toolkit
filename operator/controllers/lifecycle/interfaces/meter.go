package interfaces

import (
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument/asyncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/asyncint64"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/meter_mock.go . IMeter
type IMeter = metric.Meter

//go:generate moq -pkg fake -skip-ensure -out ./fake/async/tracer_provider_int_mock.go . ITracerProviderAsyncInt64
type ITracerProviderAsyncInt64 = asyncint64.InstrumentProvider

//go:generate moq -pkg fake -skip-ensure -out ./fake/async/tracer_provider_float_mock.go . ITracerProviderAsyncFloat64
type ITracerProviderAsyncFloat64 = asyncfloat64.InstrumentProvider

//go:generate moq -pkg fake -skip-ensure -out ./fake/sync/tracer_provider_int_mock.go . ITracerProviderSyncInt64
type ITracerProviderSyncInt64 = syncint64.InstrumentProvider

//go:generate moq -pkg fake -skip-ensure -out ./fake/sync/tracer_provider_float_mock.go . ITracerProviderSyncFloat64
type ITracerProviderSyncFloat64 = syncfloat64.InstrumentProvider
