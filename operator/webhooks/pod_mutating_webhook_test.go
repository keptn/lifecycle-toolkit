package webhooks

import (
	"context"
	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1"
	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1/common"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/common/fake"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
	admissionv1 "k8s.io/api/admission/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
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
				Kind: "",
				UID:  "",
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

func TestPodMutatingWebhook_isPodAnnotated(t *testing.T) {
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
		name      string
		fields    fields
		args      args
		want      bool
		wantErr   bool
		wantedPod *corev1.Pod
	}{
		{
			name: "Test error when workload name is too long",
			args: args{
				pod: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							common.AppAnnotation:      "SOME-APP-NAME-ANNOTATION",
							common.WorkloadAnnotation: "workload-name-that-is-too-loooooooooooooooooooooooooooooooooooooooooooooooooong",
						},
					},
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Test return true when pod has workload annotation",
			args: args{
				pod: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							common.WorkloadAnnotation: "some-workload-name",
						},
					},
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Test return true and initialize annotations when labels are set",
			args: args{
				pod: &corev1.Pod{
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Image: "some-image:v1",
							},
						},
					},
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							common.WorkloadAnnotation: "some-workload-name",
						},
					},
				},
			},
			want:    true,
			wantErr: false,
			wantedPod: &corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Image: "some-image:v1",
						},
					},
				},
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						common.WorkloadAnnotation: "some-workload-name",
					},
					Annotations: map[string]string{
						common.VersionAnnotation: "v1",
					},
				},
			},
		},
		{
			name: "Test return false when annotations and labels are not set",
			args: args{
				pod: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"some-other-label": "some-value",
						},
					},
				},
			},
			want:    false,
			wantErr: false,
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
			got, err := a.isPodAnnotated(tt.args.pod)
			if (err != nil) != tt.wantErr {
				t.Errorf("isPodAnnotated() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("isPodAnnotated() got = %v, want %v", got, tt.want)
			}
			if tt.wantedPod != nil {
				require.Equal(t, tt.wantedPod, tt.args.pod)
			}
		})
	}
}

func TestPodMutatingWebhook_copyAnnotationsIfParentAnnotated(t *testing.T) {
	testNamespace := "test-namespace"
	rsUidWithDpOwner := types.UID("this-is-the-replicaset-with-dp-owner")
	rsUidWithStsOwner := types.UID("this-is-the-replicaset-with-sts-owner")
	rsUidWithDsOwner := types.UID("this-is-the-replicaset-with-ds-owner")
	rsUidWithNoOwner := types.UID("this-is-the-replicaset-with-no-owner")

	fakeClient, err := fake.NewClient()

	if err != nil {
		t.Errorf("Error when setting up fake client %v", err)
	}

	rsWithDpOwner := &appsv1.ReplicaSet{
		TypeMeta: metav1.TypeMeta{
			Kind: "ReplicaSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-replicaset1",
			UID:       rsUidWithDpOwner,
			Namespace: testNamespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					Kind: "Deployment",
					Name: "this-is-the-deployment",
					UID:  "this-is-the-deployment-uid",
				},
			},
		},
	}
	rsWithStsOwner := &appsv1.ReplicaSet{
		TypeMeta: metav1.TypeMeta{
			Kind: "ReplicaSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-replicaset2",
			UID:       rsUidWithStsOwner,
			Namespace: testNamespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					Kind: "Deployment",
					Name: "this-is-the-deployment",
					UID:  "this-is-the-stateful-set-uid",
				},
			},
		},
	}
	rsWithDsOwner := &appsv1.ReplicaSet{
		TypeMeta: metav1.TypeMeta{
			Kind: "ReplicaSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-replicaset3",
			UID:       rsUidWithDsOwner,
			Namespace: testNamespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					Kind: "Deployment",
					Name: "this-is-the-deployment",
					UID:  "this-is-the-daemonset-uid",
				},
			},
		},
	}
	testRsWithNoOwner := &appsv1.ReplicaSet{
		TypeMeta: metav1.TypeMeta{
			Kind: "ReplicaSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-replicaset4",
			UID:       rsUidWithNoOwner,
			Namespace: testNamespace,
		},
	}
	testDp := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind: "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-deployment",
			UID:       "this-is-the-deployment-uid",
			Namespace: testNamespace,
		},
	}
	testSts := &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			Kind: "StatefulSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-stateful-set",
			UID:       "this-is-the-stateful-set-uid",
			Namespace: testNamespace,
		},
	}
	testDs := &appsv1.DaemonSet{
		TypeMeta: metav1.TypeMeta{
			Kind: "DaemonSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-daemonset",
			UID:       "this-is-the-daemonset-uid",
			Namespace: testNamespace,
		},
	}

	err = fakeClient.Create(context.TODO(), rsWithDpOwner)
	err = fakeClient.Create(context.TODO(), rsWithStsOwner)
	err = fakeClient.Create(context.TODO(), rsWithDsOwner)
	err = fakeClient.Create(context.TODO(), testRsWithNoOwner)
	err = fakeClient.Create(context.TODO(), testDp)
	err = fakeClient.Create(context.TODO(), testSts)
	err = fakeClient.Create(context.TODO(), testDs)

	if err != nil {
		t.Errorf("Error when creating objects in fake client %v", err)
	}

	type fields struct {
		Client   client.Client
		Tracer   trace.Tracer
		decoder  *admission.Decoder
		Recorder record.EventRecorder
		Log      logr.Logger
	}
	type args struct {
		ctx context.Context
		req *admission.Request
		pod *corev1.Pod
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Test that nothing happens if owner UID is pod UID",
			fields: fields{
				Log: testr.New(t),
			},
			args: args{
				ctx: context.TODO(),
				req: nil,
				pod: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						UID: "some-uid",
						OwnerReferences: []metav1.OwnerReference{
							{
								UID:  "some-uid",
								Kind: "ReplicaSet",
							},
						},
					},
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Test fetching of replicaset owner of pod and deployment owner of replicaset",
			fields: fields{
				Log:    testr.New(t),
				Client: fakeClient,
			},
			args: args{
				ctx: context.TODO(),
				req: &admission.Request{
					AdmissionRequest: admissionv1.AdmissionRequest{
						Namespace: testNamespace,
					},
				},
				pod: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						UID: "this-is-the-pod-uid",
						OwnerReferences: []metav1.OwnerReference{
							{
								UID:  rsUidWithDpOwner,
								Kind: "ReplicaSet",
							},
						},
					},
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Test fetching of replicaset owner of pod and statefulset owner of replicaset",
			fields: fields{
				Log:    testr.New(t),
				Client: fakeClient,
			},
			args: args{
				ctx: context.TODO(),
				req: &admission.Request{
					AdmissionRequest: admissionv1.AdmissionRequest{
						Namespace: testNamespace,
					},
				},
				pod: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						UID: "this-is-the-pod-uid",
						OwnerReferences: []metav1.OwnerReference{
							{
								UID:  rsUidWithStsOwner,
								Kind: "ReplicaSet",
							},
						},
					},
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Test fetching of replicaset owner of pod and daemonset owner of replicaset",
			fields: fields{
				Log:    testr.New(t),
				Client: fakeClient,
			},
			args: args{
				ctx: context.TODO(),
				req: &admission.Request{
					AdmissionRequest: admissionv1.AdmissionRequest{
						Namespace: testNamespace,
					},
				},
				pod: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						UID: "this-is-the-pod-uid",
						OwnerReferences: []metav1.OwnerReference{
							{
								UID:  rsUidWithDsOwner,
								Kind: "ReplicaSet",
							},
						},
					},
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Test that method returns without doing anything when we get a pod with replicaset without owner",
			fields: fields{
				Log:    testr.New(t),
				Client: fakeClient,
			},
			args: args{
				ctx: context.TODO(),
				req: &admission.Request{
					AdmissionRequest: admissionv1.AdmissionRequest{
						Namespace: testNamespace,
					},
				},
				pod: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						UID: "this-is-the-pod-uid",
						OwnerReferences: []metav1.OwnerReference{
							{
								UID:  rsUidWithNoOwner,
								Kind: "ReplicaSet",
							},
						},
					},
				},
			},
			want:    false,
			wantErr: false,
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
			got, err := a.copyAnnotationsIfParentAnnotated(tt.args.ctx, tt.args.req, tt.args.pod)
			if (err != nil) != tt.wantErr {
				t.Errorf("copyAnnotationsIfParentAnnotated() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("copyAnnotationsIfParentAnnotated() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPodMutatingWebhook_copyResourceLabelsIfPresent(t *testing.T) {
	type fields struct {
		Client   client.Client
		Tracer   trace.Tracer
		decoder  *admission.Decoder
		Recorder record.EventRecorder
		Log      logr.Logger
	}
	type args struct {
		sourceResource *metav1.ObjectMeta
		targetPod      *corev1.Pod
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      bool
		wantErr   bool
		wantedPod *corev1.Pod
	}{
		{
			name: "Test that annotations get copied from source to target",
			args: args{
				sourceResource: &metav1.ObjectMeta{
					Name: "testSourceObject",
					Annotations: map[string]string{
						common.WorkloadAnnotation:                 "some-workload-name",
						common.AppAnnotation:                      "some-app-name",
						common.VersionAnnotation:                  "v1.0.0",
						common.PreDeploymentTaskAnnotation:        "some-pre-deployment-task",
						common.PostDeploymentTaskAnnotation:       "some-post-deployment-task",
						common.PreDeploymentEvaluationAnnotation:  "some-pre-deployment-evaluation",
						common.PostDeploymentEvaluationAnnotation: "some-post-deployment-evaluation",
					},
				},
				targetPod: &corev1.Pod{
					TypeMeta:   metav1.TypeMeta{},
					ObjectMeta: metav1.ObjectMeta{},
					Spec:       corev1.PodSpec{},
					Status:     corev1.PodStatus{},
				},
			},
			want:    true,
			wantErr: false,
			wantedPod: &corev1.Pod{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						common.WorkloadAnnotation:                 "some-workload-name",
						common.AppAnnotation:                      "some-app-name",
						common.VersionAnnotation:                  "v1.0.0",
						common.PreDeploymentTaskAnnotation:        "some-pre-deployment-task",
						common.PostDeploymentTaskAnnotation:       "some-post-deployment-task",
						common.PreDeploymentEvaluationAnnotation:  "some-pre-deployment-evaluation",
						common.PostDeploymentEvaluationAnnotation: "some-post-deployment-evaluation",
					},
				},
			},
		},
		{
			name: "Test that source labels get copied to target annotations",
			args: args{
				sourceResource: &metav1.ObjectMeta{
					Name: "testSourceObject",
					Labels: map[string]string{
						common.WorkloadAnnotation:                 "some-workload-name",
						common.AppAnnotation:                      "some-app-name",
						common.VersionAnnotation:                  "v1.0.0",
						common.PreDeploymentTaskAnnotation:        "some-pre-deployment-task",
						common.PostDeploymentTaskAnnotation:       "some-post-deployment-task",
						common.PreDeploymentEvaluationAnnotation:  "some-pre-deployment-evaluation",
						common.PostDeploymentEvaluationAnnotation: "some-post-deployment-evaluation",
					},
				},
				targetPod: &corev1.Pod{},
			},
			want:    true,
			wantErr: false,
			wantedPod: &corev1.Pod{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						common.WorkloadAnnotation:                 "some-workload-name",
						common.AppAnnotation:                      "some-app-name",
						common.VersionAnnotation:                  "v1.0.0",
						common.PreDeploymentTaskAnnotation:        "some-pre-deployment-task",
						common.PostDeploymentTaskAnnotation:       "some-post-deployment-task",
						common.PreDeploymentEvaluationAnnotation:  "some-pre-deployment-evaluation",
						common.PostDeploymentEvaluationAnnotation: "some-post-deployment-evaluation",
					},
				},
			},
		},
		{
			name: "Test that version label is generated correctly and rest is copied",
			args: args{
				sourceResource: &metav1.ObjectMeta{
					Name: "testSourceObject",
					Labels: map[string]string{
						common.WorkloadAnnotation:                 "some-workload-name",
						common.AppAnnotation:                      "some-app-name",
						common.PreDeploymentTaskAnnotation:        "some-pre-deployment-task",
						common.PostDeploymentTaskAnnotation:       "some-post-deployment-task",
						common.PreDeploymentEvaluationAnnotation:  "some-pre-deployment-evaluation",
						common.PostDeploymentEvaluationAnnotation: "some-post-deployment-evaluation",
					},
				},
				targetPod: &corev1.Pod{
					TypeMeta:   metav1.TypeMeta{},
					ObjectMeta: metav1.ObjectMeta{},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Image: "some-image:v1.0.0",
							},
						},
					},
					Status: corev1.PodStatus{},
				},
			},
			want:    true,
			wantErr: false,
			wantedPod: &corev1.Pod{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						common.WorkloadAnnotation:                 "some-workload-name",
						common.AppAnnotation:                      "some-app-name",
						common.VersionAnnotation:                  "v1.0.0",
						common.PreDeploymentTaskAnnotation:        "some-pre-deployment-task",
						common.PostDeploymentTaskAnnotation:       "some-post-deployment-task",
						common.PreDeploymentEvaluationAnnotation:  "some-pre-deployment-evaluation",
						common.PostDeploymentEvaluationAnnotation: "some-post-deployment-evaluation",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Image: "some-image:v1.0.0",
						},
					},
				},
				Status: corev1.PodStatus{},
			},
		},
		{
			name: "Test that error is return with too long workload name",
			args: args{
				sourceResource: &metav1.ObjectMeta{
					Name: "testSourceObject",
					Labels: map[string]string{
						common.WorkloadAnnotation: "some-workload-name-that-is-very-looooooooooooooooooooooong",
					},
				},
				targetPod: &corev1.Pod{},
			},
			want:    false,
			wantErr: true,
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
			got, err := a.copyResourceLabelsIfPresent(tt.args.sourceResource, tt.args.targetPod)
			if (err != nil) != tt.wantErr {
				t.Errorf("copyResourceLabelsIfPresent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("copyResourceLabelsIfPresent() got = %v, want %v", got, tt.want)
			}
			if tt.wantedPod != nil {
				require.Equal(t, tt.wantedPod, tt.args.targetPod)
			}
		})
	}
}
