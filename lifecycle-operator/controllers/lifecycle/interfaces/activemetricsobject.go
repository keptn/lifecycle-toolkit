package interfaces

import (
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/errors"
	"go.opentelemetry.io/otel/attribute"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ActiveMetricsObject represents an object whose active metrics are stored
//
//go:generate moq -pkg fake --skip-ensure -out ./fake/activemetricsobject_mock.go . ActiveMetricsObject
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
		return nil, errors.ErrCannotWrapToActiveMetricsObject
	}
	return &ActiveMetricsObjectWrapper{Obj: amo}, nil
}

func (amo ActiveMetricsObjectWrapper) GetActiveMetricsAttributes() []attribute.KeyValue {
	return amo.Obj.GetActiveMetricsAttributes()
}

func (amo ActiveMetricsObjectWrapper) IsEndTimeSet() bool {
	return amo.Obj.IsEndTimeSet()
}
