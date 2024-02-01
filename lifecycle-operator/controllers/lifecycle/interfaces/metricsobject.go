package interfaces

import (
	"time"

	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/errors"
	"go.opentelemetry.io/otel/attribute"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// MetricsObject represents an object whose metrics are stored
//
//go:generate moq -pkg fake --skip-ensure -out ./fake/metricsobject_mock.go . MetricsObject
type MetricsObject interface {
	GetDurationMetricsAttributes() []attribute.KeyValue
	GetMetricsAttributes() []attribute.KeyValue
	GetEndTime() time.Time
	GetStartTime() time.Time
	IsEndTimeSet() bool
	GetPreviousVersion() string
	GetParentName() string
	GetNamespace() string
}

type MetricsObjectWrapper struct {
	Obj MetricsObject
}

func NewMetricsObjectWrapperFromClientObject(object client.Object) (*MetricsObjectWrapper, error) {
	mo, ok := object.(MetricsObject)
	if !ok {
		return nil, errors.ErrCannotWrapToMetricsObject
	}
	return &MetricsObjectWrapper{Obj: mo}, nil
}

func (mo MetricsObjectWrapper) GetMetricsAttributes() []attribute.KeyValue {
	return mo.Obj.GetMetricsAttributes()
}

func (mo MetricsObjectWrapper) GetDurationMetricsAttributes() []attribute.KeyValue {
	return mo.Obj.GetDurationMetricsAttributes()
}

func (mo MetricsObjectWrapper) GetEndTime() time.Time {
	return mo.Obj.GetEndTime()
}

func (mo MetricsObjectWrapper) GetStartTime() time.Time {
	return mo.Obj.GetStartTime()
}

func (mo MetricsObjectWrapper) IsEndTimeSet() bool {
	return mo.Obj.IsEndTimeSet()
}

func (mo MetricsObjectWrapper) GetPreviousVersion() string {
	return mo.Obj.GetPreviousVersion()
}

func (mo MetricsObjectWrapper) GetParentName() string {
	return mo.Obj.GetParentName()
}

func (mo MetricsObjectWrapper) GetNamespace() string {
	return mo.Obj.GetNamespace()
}
