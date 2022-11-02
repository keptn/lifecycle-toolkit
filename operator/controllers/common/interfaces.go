package common

import (
	"errors"

	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1/common"
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
	GetSpanName(phase string) string
	Complete()
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

func (pw *PhaseItemWrapper) Complete() {
	pw.Obj.Complete()
}

func (pw PhaseItemWrapper) GetVersion() string {
	return pw.Obj.GetVersion()
}

func (pw PhaseItemWrapper) GetSpanName(phase string) string {
	return pw.Obj.GetSpanName(phase)
}
