package interfaces

import (
	"go.opentelemetry.io/otel/metric"
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/meter_mock.go . IMeter
type IMeter interface {
	Int64Counter(name string, options ...metric.Int64CounterOption) (metric.Int64Counter, error)
	Int64Histogram(name string, options ...metric.Int64HistogramOption) (metric.Int64Histogram, error)
	Float64Counter(name string, options ...metric.Float64CounterOption) (metric.Float64Counter, error)
	Float64Histogram(name string, options ...metric.Float64HistogramOption) (metric.Float64Histogram, error)
	RegisterCallback(f metric.Callback, instruments ...metric.Observable) (metric.Registration, error)
	Int64ObservableGauge(name string, options ...metric.Int64ObservableGaugeOption) (metric.Int64ObservableGauge, error)
	Float64ObservableGauge(name string, options ...metric.Float64ObservableGaugeOption) (metric.Float64ObservableGauge, error)
}
