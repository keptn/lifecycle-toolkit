package handlers

import (
	"context"
	"testing"

	"github.com/go-logr/logr"
	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/fake"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestAppHandler_Handle(t *testing.T) {
	type fields struct {
		Client      client.Client
		Log         logr.Logger
		Tracer      trace.Tracer
		EventSender common.IEvent
	}
	type args struct {
		ctx       context.Context
		pod       *corev1.Pod
		namespace string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AppHandler{
				Client:      tt.fields.Client,
				Log:         tt.fields.Log,
				Tracer:      tt.fields.Tracer,
				EventSender: tt.fields.EventSender,
			}
			if err := a.Handle(tt.args.ctx, tt.args.pod, tt.args.namespace); (err != nil) != tt.wantErr {
				t.Errorf("Handle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAppHandler_createApp_succeeds(t *testing.T) {
	fakeClient := fake.NewClient()
	logger := logr.Discard()
	eventSender := common.NewK8sSender(record.NewFakeRecorder(100))
	tracer := &fake.ITracerMock{StartFunc: func(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
		return ctx, trace.SpanFromContext(ctx)
	}}
	appHandler := &AppHandler{
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

func TestAppHandler_createApp_fails(t *testing.T) {
	fakeClient := fake.NewClient()
	logger := logr.Discard()
	eventSender := common.NewK8sSender(record.NewFakeRecorder(100))
	tracer := &fake.ITracerMock{StartFunc: func(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
		return ctx, trace.SpanFromContext(ctx)
	}}
	appHandler := &AppHandler{
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
