package interfaces

import (
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/errors"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// SpanItem represents an object whose metrics are stored
//
//go:generate moq -pkg fake --skip-ensure -out ./fake/spanitem_mock.go . SpanItem
type SpanItem interface {
	SetSpanAttributes(span trace.Span)
	SetPhaseTraceID(phase string, carrier propagation.MapCarrier)
	GetSpanKey(phase string) string
	GetSpanName(phase string) string
}

type SpanItemWrapper struct {
	Obj SpanItem
}

func NewSpanItemWrapperFromClientObject(object client.Object) (*SpanItemWrapper, error) {
	mo, ok := object.(SpanItem)
	if !ok {
		return nil, errors.ErrCannotWrapToSpanItem
	}
	return &SpanItemWrapper{Obj: mo}, nil
}

func (pw SpanItemWrapper) SetPhaseTraceID(phase string, carrier propagation.MapCarrier) {
	pw.Obj.SetPhaseTraceID(phase, carrier)
}
func (pw SpanItemWrapper) GetSpanKey(phase string) string {
	return pw.Obj.GetSpanKey(phase)
}

func (pw SpanItemWrapper) GetSpanName(phase string) string {
	return pw.Obj.GetSpanName(phase)
}

func (pw SpanItemWrapper) SetSpanAttributes(span trace.Span) {
	pw.Obj.SetSpanAttributes(span)
}
