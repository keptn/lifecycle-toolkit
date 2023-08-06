package common

import (
	"context"
	"fmt"
	"strings"

	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/go-logr/logr"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/config"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/event_mock.go . IEvent:MockEvent
type IEvent interface {
	Emit(phase apicommon.KeptnPhaseType, eventType string, reconcileObject client.Object, status string, message string, version string)
}

// ===== Main =====

type EventMultiplexer struct {
	logger   logr.Logger
	emitters []IEvent
}

func NewEventMultiplexer(logger logr.Logger, recorder record.EventRecorder, client ce.Client) *EventMultiplexer {
	multiplexer := &EventMultiplexer{
		logger: logger,
	}
	multiplexer.register(newCloudEventSender(logger, client))
	multiplexer.register(NewK8sSender(recorder))
	return multiplexer
}

func (e *EventMultiplexer) register(emitter IEvent) {
	if emitter != nil {
		e.emitters = append(e.emitters, emitter)
	}
}

func (e *EventMultiplexer) Emit(phase apicommon.KeptnPhaseType, eventType string, reconcileObject client.Object, status string, message string, version string) {
	for _, emitter := range e.emitters {
		e.logger.Info(fmt.Sprintf("Emitting event using %T", emitter))
		emitter.Emit(phase, eventType, reconcileObject, status, message, version)
	}
}

// ===== Cloud Event Sender =====

type cloudEvent struct {
	client ce.Client
	logger logr.Logger
}

func newCloudEventSender(logger logr.Logger, client ce.Client) *cloudEvent {
	return &cloudEvent{
		client: client,
		logger: logger,
	}
}

// Emit creates a Cloud Event and send it to the endpoint
func (e *cloudEvent) Emit(phase apicommon.KeptnPhaseType, eventType string, reconcileObject client.Object, status string, message string, version string) {
	endpoint := config.Instance().GetCloudEventsEndpoint()
	if endpoint == "" {
		// if no endpoint is configured we don't emit any event
		if !strings.HasPrefix(endpoint, "http") {
			e.logger.V(5).Info(fmt.Sprintf("CloudEvent endpoint configured but it does not start with http: %s", endpoint))
		}
		return
	}
	event := ce.NewEvent()
	event.SetSource("keptn.sh")
	event.SetType(fmt.Sprintf("%s.%s", phase.LongName, status))

	msg := setEventMessage(phase, reconcileObject, message, version)
	err := event.SetData(ce.ApplicationJSON, map[string]string{
		"message": msg,
		"type":    eventType,
		"version": version,
	})
	if err != nil {
		e.logger.V(5).Info(fmt.Sprintf("Failed to set data for CloudEvent: %v", err))
		return
	}

	ctx := ce.ContextWithTarget(context.TODO(), endpoint)
	if result := e.client.Send(ctx, event); ce.IsUndelivered(result) {
		e.logger.V(5).Info(fmt.Sprintf("Failed to send CloudEvent: %v", event))
	}
}

// ===== K8s Event Sender =====

type k8sEvent struct {
	recorder record.EventRecorder
}

func NewK8sSender(recorder record.EventRecorder) IEvent {
	return &k8sEvent{
		recorder: recorder,
	}
}

// Emit creates k8s Event and adds it to Eventqueue
func (e *k8sEvent) Emit(phase apicommon.KeptnPhaseType, eventType string, reconcileObject client.Object, status string, message string, version string) {
	msg := setEventMessage(phase, reconcileObject, message, version)
	annotations := setAnnotations(reconcileObject, phase)
	e.recorder.AnnotatedEventf(reconcileObject, annotations, eventType, fmt.Sprintf("%s%s", phase.ShortName, status), msg)
}
