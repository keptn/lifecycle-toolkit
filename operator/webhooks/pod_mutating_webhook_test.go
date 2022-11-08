package webhooks

import (
	"github.com/go-logr/logr"
	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1"
	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1/common"
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
							common.AppAnnotation: "SOME-APP-NAME",
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
							common.AppAnnotation: "SOME-APP-NAME",
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
							common.AppAnnotation: "SOME-APP-NAME-ANNOTATION",
						},
						Labels: map[string]string{
							common.AppAnnotation: "SOME-APP-NAME-LABEL",
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
							common.AppAnnotation:      "SOME-APP-NAME",
							common.WorkloadAnnotation: "SOME-WORKLOAD-NAME",
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
							common.AppAnnotation:      "SOME-APP-NAME",
							common.WorkloadAnnotation: "SOME-WORKLOAD-NAME",
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
							common.AppAnnotation:      "SOME-APP-NAME-ANNOTATION",
							common.WorkloadAnnotation: "SOME-WORKLOAD-NAME-ANNOTATION",
						},
						Labels: map[string]string{
							common.AppAnnotation:      "SOME-APP-NAME-LABEL",
							common.WorkloadAnnotation: "SOME-WORKLOAD-NAME-LABEL",
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

func Test_getLabelOrAnnotation(t *testing.T) {
	type args struct {
		resource            *metav1.ObjectMeta
		primaryAnnotation   string
		secondaryAnnotation string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		{
			name: "Test if primary annotation is returned from annotations",
			args: args{
				resource: &metav1.ObjectMeta{
					Annotations: map[string]string{
						common.AppAnnotation: "some-app-name",
					},
				},
				primaryAnnotation:   common.AppAnnotation,
				secondaryAnnotation: common.K8sRecommendedAppAnnotations,
			},
			want:  "some-app-name",
			want1: true,
		},
		{
			name: "Test if secondary annotation is returned from annotations",
			args: args{
				resource: &metav1.ObjectMeta{
					Annotations: map[string]string{
						common.K8sRecommendedAppAnnotations: "some-app-name",
					},
				},
				primaryAnnotation:   common.AppAnnotation,
				secondaryAnnotation: common.K8sRecommendedAppAnnotations,
			},
			want:  "some-app-name",
			want1: true,
		},
		{
			name: "Test if primary annotation is returned from labels",
			args: args{
				resource: &metav1.ObjectMeta{
					Labels: map[string]string{
						common.AppAnnotation: "some-app-name",
					},
				},
				primaryAnnotation:   common.AppAnnotation,
				secondaryAnnotation: common.K8sRecommendedAppAnnotations,
			},
			want:  "some-app-name",
			want1: true,
		},
		{
			name: "Test if secondary annotation is returned from labels",
			args: args{
				resource: &metav1.ObjectMeta{
					Labels: map[string]string{
						common.K8sRecommendedAppAnnotations: "some-app-name",
					},
				},
				primaryAnnotation:   common.AppAnnotation,
				secondaryAnnotation: common.K8sRecommendedAppAnnotations,
			},
			want:  "some-app-name",
			want1: true,
		},
		{
			name: "Test that empty string is returned when no annotations or labels are found",
			args: args{
				resource: &metav1.ObjectMeta{
					Annotations: map[string]string{
						"some-other-annotation": "some-app-name",
					},
				},
				primaryAnnotation:   common.AppAnnotation,
				secondaryAnnotation: common.K8sRecommendedAppAnnotations,
			},
			want:  "",
			want1: false,
		},
		{
			name: "Test that empty string is returned when primary annotation cannot be found and secondary annotation is empty",
			args: args{
				resource: &metav1.ObjectMeta{
					Annotations: map[string]string{
						"some-other-annotation": "some-app-name",
					},
				},
				primaryAnnotation:   common.AppAnnotation,
				secondaryAnnotation: "",
			},
			want:  "",
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := getLabelOrAnnotation(tt.args.resource, tt.args.primaryAnnotation, tt.args.secondaryAnnotation)
			if got != tt.want {
				t.Errorf("getLabelOrAnnotation() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("getLabelOrAnnotation() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
