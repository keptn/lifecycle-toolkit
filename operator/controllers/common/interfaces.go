package common

import (
	"errors"
	"time"

	"github.com/keptn/lifecycle-controller/operator/api/v1alpha1/common"
	"go.opentelemetry.io/otel/attribute"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//go:generate moq -pkg common_mock --skip-ensure -out ./fake/phaseitem_mock.go . PhaseItem
type PhaseItem interface {
	GetState() common.KeptnState
	SetState(common.KeptnState)
	GetCurrentPhase() string
	SetCurrentPhase(string)
	GetVersion() string
	GetMetricsAttributes() []attribute.KeyValue
	GetActiveMetricsAttributes() []attribute.KeyValue
	GetSpanName(phase string) string
	Complete()
	IsEndTimeSet() bool
	GetDurationMetricsAttributes() []attribute.KeyValue
	GetEndTime() time.Time
	GetStartTime() time.Time
	GetPreviousVersion() string
	GetParentName() string
	GetNamespace() string
}

type PhaseItemWrapper struct {
	Obj PhaseItem
}

func NewPhaseItemWrapperFromClientObject(object client.Object) (*PhaseItemWrapper, error) {
	pi, ok := object.(PhaseItem)
	if !ok {
		return nil, errors.New("provided object does not implement PhaseItem interface")
	}
	return &PhaseItemWrapper{Obj: pi}, nil
}

func (pw PhaseItemWrapper) GetState() common.KeptnState {
	return pw.Obj.GetState()
}

func (pw *PhaseItemWrapper) SetState(state common.KeptnState) {
	pw.Obj.SetState(state)
}

func (pw PhaseItemWrapper) GetCurrentPhase() string {
	return pw.Obj.GetCurrentPhase()
}

func (pw *PhaseItemWrapper) SetCurrentPhase(phase string) {
	pw.Obj.SetCurrentPhase(phase)
}

func (pw PhaseItemWrapper) GetMetricsAttributes() []attribute.KeyValue {
	return pw.Obj.GetMetricsAttributes()
}

func (pw PhaseItemWrapper) GetDurationMetricsAttributes() []attribute.KeyValue {
	return pw.Obj.GetDurationMetricsAttributes()
}

func (pw PhaseItemWrapper) GetEndTime() time.Time {
	return pw.Obj.GetEndTime()
}

func (pw PhaseItemWrapper) GetStartTime() time.Time {
	return pw.Obj.GetStartTime()
}

func (pw *PhaseItemWrapper) Complete() {
	pw.Obj.Complete()
}

func (pw PhaseItemWrapper) GetVersion() string {
	return pw.Obj.GetVersion()
}

func (pw PhaseItemWrapper) GetSpanName(phase string) string {
	return pw.Obj.GetSpanName(phase)
}

func (pw PhaseItemWrapper) IsEndTimeSet() bool {
	return pw.Obj.IsEndTimeSet()
}

func (pw PhaseItemWrapper) GetPreviousVersion() string {
	return pw.Obj.GetPreviousVersion()
}

func (pw PhaseItemWrapper) GetParentName() string {
	return pw.Obj.GetParentName()
}

func (pw PhaseItemWrapper) GetNamespace() string {
	return pw.Obj.GetNamespace()
}

func (pw PhaseItemWrapper) GetActiveMetricsAttributes() []attribute.KeyValue {
	return pw.Obj.GetActiveMetricsAttributes()
}

func NewListItemWrapperFromClientObjectList(object client.ObjectList) (*ListItemWrapper, error) {
	pi, ok := object.(ListItem)
	if !ok {
		return nil, errors.New("provided object does not implement ListItem interface")
	}
	return &ListItemWrapper{Obj: pi}, nil
}

type ListItem interface {
	GetItems() []PhaseItem
}

type ListItemWrapper struct {
	Obj ListItem
}

func (pw ListItemWrapper) GetItems() []PhaseItem {
	return pw.Obj.GetItems()
}
