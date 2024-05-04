package eventsender

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	ce "github.com/cloudevents/sdk-go/v2"
	apilifecycle "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/config"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/eventsender/fake"
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

	eventSender.Emit(apicommon.PhaseAppDeployment, "pre-event", &apilifecycle.KeptnAppVersion{
		ObjectMeta: v1.ObjectMeta{
			Name:      "app",
			Namespace: "ns",
		},
	}, "reason-short", "reason-long", "ver1")

	event := <-fakeRecorder.Events

	require.Contains(t, event, fmt.Sprintf("%s: reason-long / Namespace: ns, Name: app, Version: ver1", apicommon.PhaseAppDeployment.LongName))
}

func TestEventSender_SendCloudEvent(t *testing.T) {
	//config
	name := "app"
	ns := "my-ns"
	status := "my-status"
	eventType := "my-type"
	version := "v0.0.1-dev"
	msg := "my message"
	phase := apicommon.PhaseAppDeployment
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
	ceSender.Emit(phase, eventType, &apilifecycle.KeptnAppVersion{
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
		ceSender.Emit(apicommon.PhaseAppCompleted, "type", &apilifecycle.KeptnAppVersion{
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
	em1.EmitFunc = func(phase apicommon.KeptnPhaseType, eventType string, reconcileObject client.Object, status string, message string, version string) {
		recE1 <- struct{}{}
	}

	em2 := &fake.MockEvent{}
	em2.EmitFunc = func(phase apicommon.KeptnPhaseType, eventType string, reconcileObject client.Object, status string, message string, version string) {
		recE2 <- struct{}{}
	}
	emitter := EventMultiplexer{}
	emitter.register(em1)
	emitter.register(em2)
	// then
	// fire a new event
	msg := "my special message"
	emitter.Emit(apicommon.PhaseAppDeployment, "", nil, "", msg, "")
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

func Test_setEventMessage(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    string
	}{
		{
			name:    "version empty",
			version: "",
			want:    "App Deployment: longReason / Namespace: namespace, Name: app",
		},
		{
			name:    "version set",
			version: "1.0.0",
			want:    "App Deployment: longReason / Namespace: namespace, Name: app, Version: 1.0.0",
		},
	}

	appVersion := &apilifecycle.KeptnAppVersion{
		ObjectMeta: v1.ObjectMeta{
			Name:      "app",
			Namespace: "namespace",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, setEventMessage(apicommon.PhaseAppDeployment, appVersion, "longReason", tt.version), tt.want)
		})
	}
}

func Test_setAnnotations(t *testing.T) {
	tests := []struct {
		name   string
		object client.Object
		want   map[string]string
	}{
		{
			name:   "nil object",
			object: nil,
			want:   nil,
		},
		{
			name:   "empty object",
			object: &apilifecycle.KeptnEvaluationDefinition{},
			want:   nil,
		},
		{
			name: "unknown object",
			object: &apilifecycle.KeptnEvaluationDefinition{
				ObjectMeta: v1.ObjectMeta{
					Name:      "def",
					Namespace: "namespace",
				},
			},
			want: map[string]string{
				"namespace":   "namespace",
				"name":        "def",
				"phase":       "AppDeploy",
				"traceparent": "",
			},
		},
		{
			name: "object with traceparent",
			object: &apilifecycle.KeptnEvaluationDefinition{
				ObjectMeta: v1.ObjectMeta{
					Name:      "def",
					Namespace: "namespace",
					Annotations: map[string]string{
						"traceparent": "23232333",
					},
				},
			},
			want: map[string]string{
				"namespace":   "namespace",
				"name":        "def",
				"phase":       "AppDeploy",
				"traceparent": "23232333",
			},
		},
		{
			name: "KeptnApp",
			object: &apilifecycle.KeptnApp{
				ObjectMeta: v1.ObjectMeta{
					Name:       "app",
					Namespace:  "namespace",
					Generation: 1,
				},
				Spec: apilifecycle.KeptnAppSpec{
					Version: "1.0.0",
				},
			},
			want: map[string]string{
				"namespace":   "namespace",
				"name":        "app",
				"phase":       "AppDeploy",
				"appName":     "app",
				"appVersion":  "1.0.0",
				"appRevision": "6b86b273",
				"traceparent": "",
			},
		},
		{
			name: "KeptnAppVersion",
			object: &apilifecycle.KeptnAppVersion{
				ObjectMeta: v1.ObjectMeta{
					Name:      "appVersion",
					Namespace: "namespace",
				},
				Spec: apilifecycle.KeptnAppVersionSpec{
					AppName: "app",
					KeptnAppSpec: apilifecycle.KeptnAppSpec{
						Version: "1.0.0",
					},
				},
			},
			want: map[string]string{
				"namespace":      "namespace",
				"name":           "appVersion",
				"phase":          "AppDeploy",
				"appName":        "app",
				"appVersion":     "1.0.0",
				"appVersionName": "appVersion",
				"traceparent":    "",
			},
		},
		{
			name: "KeptnWorkload",
			object: &apilifecycle.KeptnWorkload{
				ObjectMeta: v1.ObjectMeta{
					Name:      "workload",
					Namespace: "namespace",
				},
				Spec: apilifecycle.KeptnWorkloadSpec{
					AppName: "app",
					Version: "1.0.0",
				},
			},
			want: map[string]string{
				"namespace":       "namespace",
				"name":            "workload",
				"phase":           "AppDeploy",
				"appName":         "app",
				"workloadVersion": "1.0.0",
				"workloadName":    "workload",
				"traceparent":     "",
			},
		},
		{
			name: "KeptnWorkloadVersion",
			object: &apilifecycle.KeptnWorkloadVersion{
				ObjectMeta: v1.ObjectMeta{
					Name:      "workloadVersion",
					Namespace: "namespace",
				},
				Spec: apilifecycle.KeptnWorkloadVersionSpec{
					KeptnWorkloadSpec: apilifecycle.KeptnWorkloadSpec{
						AppName: "app",
						Version: "1.0.0",
					},
					WorkloadName: "workload",
				},
			},
			want: map[string]string{
				"namespace":           "namespace",
				"name":                "workloadVersion",
				"phase":               "AppDeploy",
				"appName":             "app",
				"workloadVersion":     "1.0.0",
				"workloadName":        "workload",
				"workloadVersionName": "workloadVersion",
				"traceparent":         "",
			},
		},
		{
			name: "KeptnTask",
			object: &apilifecycle.KeptnTask{
				ObjectMeta: v1.ObjectMeta{
					Name:      "task",
					Namespace: "namespace",
				},
				Spec: apilifecycle.KeptnTaskSpec{
					TaskDefinition: "def",
					Context: apilifecycle.TaskContext{
						WorkloadName:    "workload",
						AppName:         "app",
						AppVersion:      "1.0.0",
						WorkloadVersion: "2.0.0",
					},
				},
			},
			want: map[string]string{
				"namespace":          "namespace",
				"name":               "task",
				"phase":              "AppDeploy",
				"appName":            "app",
				"appVersion":         "1.0.0",
				"workloadName":       "workload",
				"workloadVersion":    "2.0.0",
				"taskDefinitionName": "def",
				"taskName":           "task",
				"traceparent":        "",
			},
		},
		{
			name: "KeptnEvaluation",
			object: &apilifecycle.KeptnEvaluation{
				ObjectMeta: v1.ObjectMeta{
					Name:      "eval",
					Namespace: "namespace",
				},
				Spec: apilifecycle.KeptnEvaluationSpec{
					AppName:              "app",
					AppVersion:           "1.0.0",
					Workload:             "workload",
					WorkloadVersion:      "2.0.0",
					EvaluationDefinition: "def",
				},
			},
			want: map[string]string{
				"namespace":                "namespace",
				"name":                     "eval",
				"phase":                    "AppDeploy",
				"appName":                  "app",
				"appVersion":               "1.0.0",
				"workloadName":             "workload",
				"workloadVersion":          "2.0.0",
				"evaluationDefinitionName": "def",
				"evaluationName":           "eval",
				"traceparent":              "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, setAnnotations(tt.object, apicommon.PhaseAppDeployment), tt.want)
		})
	}
}
