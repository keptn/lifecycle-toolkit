package common

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/config"
	fake "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/fake"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func TestEventSender_SendK8sEvent(t *testing.T) {
	fakeRecorder := record.NewFakeRecorder(100)
	eventSender := NewK8sSender(fakeRecorder)

	eventSender.Emit(common.PhaseAppDeployment, "pre-event", &v1alpha3.KeptnAppVersion{
		ObjectMeta: v1.ObjectMeta{
			Name:      "app",
			Namespace: "ns",
		},
	}, "reason-short", "reason-long", "ver1")

	event := <-fakeRecorder.Events

	require.Contains(t, event, fmt.Sprintf("%s: reason-long / Namespace: ns, Name: app, Version: ver1", common.PhaseAppDeployment.LongName))
}

func TestEventSender_SendCloudEvent(t *testing.T) {
	//config
	name := "app"
	ns := "my-ns"
	status := "my-status"
	eventType := "my-type"
	version := "v0.0.1-dev"
	msg := "my message"
	phase := common.PhaseAppDeployment
	waitToReceive := make(chan bool, 1)
	// when
	// we have a CloudEvent endpoint
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "POST", r.Method)
		require.Equal(t, "/", r.URL.Path)
		require.Equal(t, 1, len(r.Header["Ce-Id"]))
		require.Equal(t, 1, len(r.Header["Ce-Time"]))
		require.Equal(t, 1, len(r.Header["Ce-Type"]))
		data, err := io.ReadAll(r.Body)
		require.Nil(t, err)
		expected := fmt.Sprintf("{\"message\":\"%s: %s / Namespace: %s, Name: %s, Version: %s\",\"resource\":{\"group\":\"\",\"kind\":\"\",\"name\":\"%s\",\"namespace\":\"%s\",\"version\":\"\"},\"type\":\"%s\",\"version\":\"%s\"}",
			phase.LongName, msg, ns, name, version, name, ns, eventType, version)
		require.Equal(t, expected, string(data))

		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(""))
		require.Nil(t, err)
		waitToReceive <- true
	}))
	defer svr.Close()
	config.Instance().SetCloudEventsEndpoint(svr.URL)

	// then
	// we send a Cloud Event
	c, err := ce.NewClientHTTP()
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}
	ceSender := newCloudEventSender(ctrl.Log.WithName("testytest"), c)
	ceSender.Emit(phase, eventType, &v1alpha3.KeptnAppVersion{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
	}, status, msg, version)

	select {
	case <-waitToReceive:
		// we sent a Cloud Event
		return
	case <-time.After(5 * time.Second):
		t.Error("Didn't receive the cloud event")
	}
}

func TestEventSender_CloudEventNoFailure(t *testing.T) {

	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "no endpoint",
			input: "",
		},
		{
			name:  "invalid endpoint",
			input: "ftp://localhost:9080/",
		},
		{
			name:  "not existing endpoint",
			input: "http://127.0.0",
		},
	}

	for _, tt := range tests {
		// when
		// we don't have a CloudEvent endpoint
		config.Instance().SetCloudEventsEndpoint(tt.input)

		// then
		// we send a Cloud Event
		c, err := ce.NewClientHTTP()
		if err != nil {
			log.Fatalf("failed to create client, %v", err)
		}
		ceSender := newCloudEventSender(ctrl.Log.WithName("testytest"), c)
		ceSender.Emit(common.PhaseAppCompleted, "type", &v1alpha3.KeptnAppVersion{
			ObjectMeta: v1.ObjectMeta{
				Name:      "app",
				Namespace: "ns",
			},
		}, "status", tt.name, "version")
		// we don't fail
	}
}

func TestEventSender_Multiplexer_register(t *testing.T) {
	tests := []struct {
		input  IEvent
		expect int
	}{
		{
			input:  &fake.MockEvent{},
			expect: 1,
		},
		{
			input:  nil,
			expect: 0,
		},
	}
	for _, tt := range tests {
		em := EventMultiplexer{}
		em.register(tt.input)
		require.Equal(t, tt.expect, len(em.emitters))
	}
}

func TestEventSender_Multiplexer_new(t *testing.T) {
	// when
	// init the object
	em := NewEventMultiplexer(zap.New(), nil, nil)
	// then assert
	// k8s and ce are registered
	require.Equal(t, 2, len(em.emitters))
}

func TestEventSender_Multiplexer_emit(t *testing.T) {
	// when
	// init the object with two emitters

	recE1 := make(chan struct{})
	recE2 := make(chan struct{})

	em1 := &fake.MockEvent{}
	em1.EmitFunc = func(phase common.KeptnPhaseType, eventType string, reconcileObject client.Object, status string, message string, version string) {
		recE1 <- struct{}{}
	}

	em2 := &fake.MockEvent{}
	em2.EmitFunc = func(phase common.KeptnPhaseType, eventType string, reconcileObject client.Object, status string, message string, version string) {
		recE2 <- struct{}{}
	}
	emitter := EventMultiplexer{}
	emitter.register(em1)
	emitter.register(em2)
	// then
	// fire a new event
	msg := "my special message"
	emitter.Emit(common.PhaseAppDeployment, "", nil, "", msg, "")
	// assert we got one event
	// wait for the emitMocks to receive the events

	select {
	case <-recE1:
		break
	case <-time.After(3 * time.Second):
		t.Error("timed out waiting for the event emitter to be called")
	}
	select {
	case <-recE2:
		break
	case <-time.After(3 * time.Second):
		t.Error("timed out waiting for the event emitter to be called")
	}

	require.Equal(t, 2, len(emitter.emitters))
	require.Len(t, em1.EmitCalls(), 1)
	require.Len(t, em2.EmitCalls(), 1)
	require.Equal(t, msg, em2.EmitCalls()[0].Message)
	require.Equal(t, msg, em2.EmitCalls()[0].Message)
}
