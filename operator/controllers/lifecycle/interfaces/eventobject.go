package interfaces

import (
	"github.com/keptn/lifecycle-toolkit/operator/controllers/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// EventObject represents an object who can send k8s Events with annotations
//
//go:generate moq -pkg fake --skip-ensure -out ./fake/eventobject_mock.go . EventObject
type EventObject interface {
	GetEventAnnotations() map[string]string
}

type EventObjectWrapper struct {
	Obj EventObject
}

func NewEventObjectWrapperFromClientObject(object client.Object) (*EventObjectWrapper, error) {
	amo, ok := object.(EventObject)
	if !ok {
		return nil, errors.ErrCannotWrapToEventObject
	}
	return &EventObjectWrapper{Obj: amo}, nil
}

func (amo EventObjectWrapper) GetEventAnnotations() map[string]string {
	return amo.Obj.GetEventAnnotations()
}
