package common

import (
	"fmt"

	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/event_mock.go . IEvent:MockEvent
type IEvent interface {
	SendEvent(phase apicommon.KeptnPhaseType, eventType string, reconcileObject client.Object, status string, message string, version string)
}

// ===== Main =====

func NewEventSender(recorder record.EventRecorder) IEvent {
	return newK8sSender(recorder)
}

// ===== Cloud Event Sender =====
// TODO: implement Cloud Event logic

// ===== K8s Event Sender =====

type k8sEvent struct {
	recorder record.EventRecorder
}

func newK8sSender(recorder record.EventRecorder) IEvent {
	return &k8sEvent{
		recorder: recorder,
	}
}

// SendEvent creates k8s Event and adds it to Eventqueue
func (s *k8sEvent) SendEvent(phase apicommon.KeptnPhaseType, eventType string, reconcileObject client.Object, status string, message string, version string) {
	msg := setEventMessage(phase, reconcileObject, message, version)
	annotations := setAnnotations(reconcileObject, phase)
	s.recorder.AnnotatedEventf(reconcileObject, annotations, eventType, fmt.Sprintf("%s%s", phase.ShortName, status), msg)
}
