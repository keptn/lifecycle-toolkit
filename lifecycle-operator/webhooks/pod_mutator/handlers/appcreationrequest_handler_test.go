package handlers

import (
	"context"
	"errors"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/fake"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	k8sfake "sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
)

func TestAppHandlerHandle(t *testing.T) {

	mockEventSender := common.NewK8sSender(record.NewFakeRecorder(100))
	log := testr.New(t)
	tr := &fake.ITracerMock{StartFunc: func(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
		return ctx, trace.SpanFromContext(ctx)
	}}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example-pod",
			Namespace: namespace,
			Annotations: map[string]string{
				apicommon.WorkloadAnnotation: "my-workload",
				apicommon.VersionAnnotation:  "0.1",
			},
		}}

	singleServiceCreationReq := &klcv1alpha3.KeptnAppCreationRequest{
		TypeMeta: metav1.TypeMeta{
			Kind:       "KeptnAppCreationRequest",
			APIVersion: "lifecycle.keptn.sh/v1alpha3",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:            "my-workload",
			Namespace:       namespace,
			ResourceVersion: "1",
			Annotations: map[string]string{
				"keptn.sh/app-type": "single-service",
			},
		},
		Spec: klcv1alpha3.KeptnAppCreationRequestSpec{AppName: "my-workload"},
	}

	tests := []struct {
		name    string
		client  client.Client
		pod     *corev1.Pod
		wanterr string
		wantReq *klcv1alpha3.KeptnAppCreationRequest
	}{
		{
			name:    "Create AppCreationRequest inherit from workload",
			pod:     pod,
			client:  fake.NewClient(),
			wantReq: singleServiceCreationReq,
		},
		{
			name:    "AppCreationRequest already exists",
			pod:     pod,
			client:  fake.NewClient(singleServiceCreationReq),
			wantReq: singleServiceCreationReq,
		},
		{
			name: "Create AppCreationRequest",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "example-pod",
					Namespace: namespace,
					Annotations: map[string]string{
						apicommon.AppAnnotation:      "my-app",
						apicommon.WorkloadAnnotation: "my-workload",
						apicommon.VersionAnnotation:  "0.1",
					},
				}},
			client: fake.NewClient(),
			wantReq: &klcv1alpha3.KeptnAppCreationRequest{
				TypeMeta: metav1.TypeMeta{
					Kind:       "KeptnAppCreationRequest",
					APIVersion: "lifecycle.keptn.sh/v1alpha3",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:            "my-app",
					Namespace:       namespace,
					ResourceVersion: "1",
				},
				Spec: klcv1alpha3.KeptnAppCreationRequestSpec{AppName: "my-app"},
			},
		},
		{
			name: "Error Fetching AppCreationRequest",
			pod:  &corev1.Pod{},
			client: k8sfake.NewClientBuilder().WithInterceptorFuncs(interceptor.Funcs{
				Get: func(ctx context.Context, client client.WithWatch, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
					return errors.New("bad")
				},
			}).Build(),
			wanterr: "could not fetch AppCreationRequest: bad",
		},
		{
			name: "Error Creating AppCreationRequest",
			pod:  pod,
			client: k8sfake.NewClientBuilder().WithInterceptorFuncs(interceptor.Funcs{
				Create: func(ctx context.Context, client client.WithWatch, obj client.Object, opts ...client.CreateOption) error {
					return errors.New("badcreate")
				},
			}).Build(),
			wanterr: "badcreate",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appHandler := &AppCreationRequestHandler{
				Client:      tt.client,
				Log:         log,
				EventSender: mockEventSender,
				Tracer:      tr,
			}
			err := appHandler.Handle(context.TODO(), tt.pod, namespace)

			if tt.wanterr != "" {
				require.NotNil(t, err)
				require.Contains(t, err.Error(), tt.wanterr)
			} else {
				require.Nil(t, err)
			}

			if tt.wantReq != nil {
				creationReq := &klcv1alpha3.KeptnAppCreationRequest{}
				err = tt.client.Get(context.TODO(), types.NamespacedName{Name: tt.wantReq.Name, Namespace: tt.wantReq.Namespace}, creationReq)
				require.Nil(t, err)
				require.Equal(t, tt.wantReq, creationReq)
			}

		})
	}
}

func TestAppHandlerCreateAppSucceeds(t *testing.T) {
	fakeClient := fake.NewClient()
	logger := logr.Discard()
	eventSender := common.NewK8sSender(record.NewFakeRecorder(100))
	tracer := &fake.ITracerMock{StartFunc: func(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
		return ctx, trace.SpanFromContext(ctx)
	}}
	appHandler := &AppCreationRequestHandler{
		Client:      fakeClient,
		Log:         logger,
		Tracer:      tracer,
		EventSender: eventSender,
	}

	ctx := context.TODO()
	name := "myappcreationreq"
	newAppCreationRequest := &klcv1alpha3.KeptnAppCreationRequest{
		ObjectMeta: metav1.ObjectMeta{Name: name},
	}
	err := appHandler.createApp(ctx, newAppCreationRequest, trace.SpanFromContext(ctx))

	require.Nil(t, err)
	creationReq := &klcv1alpha3.KeptnAppCreationRequest{}
	err = fakeClient.Get(ctx, types.NamespacedName{Name: name}, creationReq)
	require.Nil(t, err)

}

func TestAppHandlerCreateAppFails(t *testing.T) {
	fakeClient := fake.NewClient()
	logger := logr.Discard()
	eventSender := common.NewK8sSender(record.NewFakeRecorder(100))
	tracer := &fake.ITracerMock{StartFunc: func(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
		return ctx, trace.SpanFromContext(ctx)
	}}
	appHandler := &AppCreationRequestHandler{
		Client:      fakeClient,
		Log:         logger,
		Tracer:      tracer,
		EventSender: eventSender,
	}

	ctx := context.TODO()
	newAppCreationRequest := &klcv1alpha3.KeptnAppCreationRequest{
		ObjectMeta: metav1.ObjectMeta{},
	}
	err := appHandler.createApp(ctx, newAppCreationRequest, trace.SpanFromContext(ctx))
	require.Error(t, err)

}
