package webhooks

import (
	"github.com/go-logr/logr"
	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1"
	"go.opentelemetry.io/otel/trace"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"testing"
)

func TestPodMutatingWebhook_getOwnerOfReplicaSet(t *testing.T) {
	type fields struct {
		Client   client.Client
		Tracer   trace.Tracer
		decoder  *admission.Decoder
		Recorder record.EventRecorder
		Log      logr.Logger
	}
	type args struct {
		rs *appsv1.ReplicaSet
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   v1alpha1.ResourceReference
	}{
		{
			name:   "Test simple return when UID and Kind is set",
			fields: fields{},
			args: args{
				rs: &appsv1.ReplicaSet{
					ObjectMeta: metav1.ObjectMeta{
						OwnerReferences: []metav1.OwnerReference{
							{
								Kind: "Deployment",
								UID:  "someUID-123456",
							},
						},
					},
				},
			},
			want: v1alpha1.ResourceReference{
				UID:  "someUID-123456",
				Kind: "Deployment",
			},
		},
		{
			name:   "Test return is input argument if owner is not found",
			fields: fields{},
			args: args{
				rs: &appsv1.ReplicaSet{
					TypeMeta: metav1.TypeMeta{
						Kind: "ReplicaSet",
					},
					ObjectMeta: metav1.ObjectMeta{
						UID: "replicaset-UID-abc123",
						OwnerReferences: []metav1.OwnerReference{
							{
								Kind: "SomeNonExistentType",
								UID:  "someUID-123456",
							},
						},
					},
				},
			},
			want: v1alpha1.ResourceReference{
				Kind: "ReplicaSet",
				UID:  "replicaset-UID-abc123",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &PodMutatingWebhook{
				Client:   tt.fields.Client,
				Tracer:   tt.fields.Tracer,
				decoder:  tt.fields.decoder,
				Recorder: tt.fields.Recorder,
				Log:      tt.fields.Log,
			}
			if got := a.getOwnerOfReplicaSet(tt.args.rs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getOwnerOfReplicaSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPodMutatingWebhook_getReplicaSetOfPod(t *testing.T) {
	type fields struct {
		Client   client.Client
		Tracer   trace.Tracer
		decoder  *admission.Decoder
		Recorder record.EventRecorder
		Log      logr.Logger
	}
	type args struct {
		pod *corev1.Pod
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   v1alpha1.ResourceReference
	}{
		{
			name:   "Test simple return when UID and Kind is set",
			fields: fields{},
			args: args{
				pod: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						UID: "the-pod-uid",
						OwnerReferences: []metav1.OwnerReference{
							{
								Kind: "ReplicaSet",
								UID:  "the-replicaset-uid",
							},
						},
					},
				},
			},
			want: v1alpha1.ResourceReference{
				UID:  "the-replicaset-uid",
				Kind: "ReplicaSet",
			},
		},
		{
			name:   "Test return is input argument if owner is not found",
			fields: fields{},
			args: args{
				pod: &corev1.Pod{
					TypeMeta: metav1.TypeMeta{
						Kind: "Pod",
					},
					ObjectMeta: metav1.ObjectMeta{
						UID: "the-pod-uid",
						OwnerReferences: []metav1.OwnerReference{
							{
								Kind: "SomeNonExistentType",
								UID:  "the-replicaset-uid",
							},
						},
					},
				},
			},
			want: v1alpha1.ResourceReference{
				UID:  "the-pod-uid",
				Kind: "Pod",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &PodMutatingWebhook{
				Client:   tt.fields.Client,
				Tracer:   tt.fields.Tracer,
				decoder:  tt.fields.decoder,
				Recorder: tt.fields.Recorder,
				Log:      tt.fields.Log,
			}
			if got := a.getReplicaSetOfPod(tt.args.pod); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getReplicaSetOfPod() = %v, want %v", got, tt.want)
			}
		})
	}
}
