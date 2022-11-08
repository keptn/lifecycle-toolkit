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
			name: "Test simple return when UID and Kind is set",
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
			name: "Test return is input argument if owner is not found",
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
			name: "Test simple return when UID and Kind is set",
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
			name: "Test return is input argument if owner is not found",
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

func TestPodMutatingWebhook_getAppName(t *testing.T) {
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
		want   string
	}{
		{
			name: "Return keptn app name in lower case when annotation is set",
			args: args{
				pod: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							"keptn.sh/app": "SOME-APP-NAME",
						},
					},
				},
			},
			want: "some-app-name",
		},
		{
			name: "Return keptn app name in lower case when label is set",
			args: args{
				pod: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"keptn.sh/app": "SOME-APP-NAME",
						},
					},
				},
			},
			want: "some-app-name",
		},
		{
			name: "Return keptn app name from annotation in lower case when annotation and label is set",
			args: args{
				pod: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							"keptn.sh/app": "SOME-APP-NAME-ANNOTATION",
						},
						Labels: map[string]string{
							"keptn.sh/app": "SOME-APP-NAME-LABEL",
						},
					},
				},
			},
			want: "some-app-name-annotation",
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
			if got := a.getAppName(tt.args.pod); got != tt.want {
				t.Errorf("getAppName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPodMutatingWebhook_getWorkloadName(t *testing.T) {
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
		want   string
	}{
		{
			name: "Return concatenated app name and workload name in lower case when annotations are set",
			args: args{
				pod: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							"keptn.sh/app":      "SOME-APP-NAME",
							"keptn.sh/workload": "SOME-WORKLOAD-NAME",
						},
					},
				},
			},
			want: "some-app-name-some-workload-name",
		},
		{
			name: "Return concatenated app name and workload name in lower case when labels are set",
			args: args{
				pod: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"keptn.sh/app":      "SOME-APP-NAME",
							"keptn.sh/workload": "SOME-WORKLOAD-NAME",
						},
					},
				},
			},
			want: "some-app-name-some-workload-name",
		},
		{
			name: "Return concatenated keptn app name and workload name from annotation in lower case when annotations and labels are set",
			args: args{
				pod: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							"keptn.sh/app":      "SOME-APP-NAME-ANNOTATION",
							"keptn.sh/workload": "SOME-WORKLOAD-NAME-ANNOTATION",
						},
						Labels: map[string]string{
							"keptn.sh/app":      "SOME-APP-NAME-LABEL",
							"keptn.sh/workload": "SOME-WORKLOAD-NAME-LABEL",
						},
					},
				},
			},
			want: "some-app-name-annotation-some-workload-name-annotation",
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
			if got := a.getWorkloadName(tt.args.pod); got != tt.want {
				t.Errorf("getWorkloadName() = %v, want %v", got, tt.want)
			}
		})
	}
}
