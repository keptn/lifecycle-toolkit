package handlers

import (
	"context"
	"reflect"
	"testing"

	"github.com/go-logr/logr"
	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common"
	"go.opentelemetry.io/otel/trace"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	}{
		// TODO: Add test cases.
	}
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

func TestAppHandler_createApp(t *testing.T) {
	type fields struct {
		Client      client.Client
		Log         logr.Logger
		Tracer      trace.Tracer
		EventSender common.IEvent
	}
	type args struct {
		ctx                   context.Context
		newAppCreationRequest *klcv1alpha3.KeptnAppCreationRequest
		span                  trace.Span
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AppHandler{
				Client:      tt.fields.Client,
				Log:         tt.fields.Log,
				Tracer:      tt.fields.Tracer,
				EventSender: tt.fields.EventSender,
			}
			if err := a.createApp(tt.args.ctx, tt.args.newAppCreationRequest, tt.args.span); (err != nil) != tt.wantErr {
				t.Errorf("createApp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_generateAppCreationRequest(t *testing.T) {
	type args struct {
		ctx       context.Context
		pod       *corev1.Pod
		namespace string
	}
	tests := []struct {
		name string
		args args
		want *klcv1alpha3.KeptnAppCreationRequest
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateAppCreationRequest(tt.args.ctx, tt.args.pod, tt.args.namespace); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateAppCreationRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_inheritWorkloadAnnotation(t *testing.T) {
	type args struct {
		meta *metav1.ObjectMeta
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inheritWorkloadAnnotation(tt.args.meta)
		})
	}
}

func Test_isAppAnnotationPresent(t *testing.T) {
	type args struct {
		meta *metav1.ObjectMeta
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isAppAnnotationPresent(tt.args.meta); got != tt.want {
				t.Errorf("isAppAnnotationPresent() = %v, want %v", got, tt.want)
			}
		})
	}
}
