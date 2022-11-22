package common

import (
	"go.opentelemetry.io/otel/attribute"
	"sigs.k8s.io/controller-runtime/pkg/client"
)


// ActiveMetricsObject represents an object whose active metrics are stored
type ActiveMetricsObject interface {
	GetActiveMetricsAttributes() []attribute.KeyValue
	IsEndTimeSet() bool
}

type ActiveMetricsObjectWrapper struct {
	Obj ActiveMetricsObject
}

func NewActiveMetricsObjectWrapperFromClientObject(object client.Object) (*ActiveMetricsObjectWrapper, error) {
	amo, ok := object.(ActiveMetricsObject)
	if !ok {
		return nil, ErrCannotWrapToActiveMetricsObject
	}
	return &ActiveMetricsObjectWrapper{Obj: amo}, nil
}

func (amo ActiveMetricsObjectWrapper) GetActiveMetricsAttributes() []attribute.KeyValue {
	return amo.Obj.GetActiveMetricsAttributes()
}

func (amo ActiveMetricsObjectWrapper) IsEndTimeSet() bool {
	return amo.Obj.IsEndTimeSet()
}
