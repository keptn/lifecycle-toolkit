package common

import (
	"fmt"

	apicommon "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3/common"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type EventSender struct {
	recorder record.EventRecorder
}

func NewEventSender(recorder record.EventRecorder) EventSender {
	return EventSender{
		recorder: recorder,
	}
}

// SendK8sEvent creates k8s Event and adds it to Eventqueue
func (s *EventSender) SendK8sEvent(phase apicommon.KeptnPhaseType, eventType string, reconcileObject client.Object, shortReason string, longReason string, version string) {
	msg := setEventMessage(phase, reconcileObject, longReason, version)
	annotations := setAnnotations(reconcileObject, phase)
	s.recorder.AnnotatedEventf(reconcileObject, annotations, eventType, fmt.Sprintf("%s%s", phase.ShortName, shortReason), msg)
}
