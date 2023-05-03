package pod_mutator

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3/common"
	fakeclient "github.com/keptn/lifecycle-toolkit/operator/controllers/common/fake"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
	admissionv1 "k8s.io/api/admission/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

func TestPodMutatingWebhook_getOwnerReference(t *testing.T) {
	type fields struct {
		Client   client.Client
		Tracer   trace.Tracer
		decoder  *admission.Decoder
		Recorder record.EventRecorder
		Log      logr.Logger
	}
	type args struct {
		resource *metav1.ObjectMeta
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   metav1.OwnerReference
	}{
		{
			name: "Test simple return when UID and Kind is set",
			args: args{
				resource: &metav1.ObjectMeta{
					UID: "the-pod-uid",
					OwnerReferences: []metav1.OwnerReference{
						{
							Kind: "ReplicaSet",
							UID:  "the-replicaset-uid",
							Name: "some-name",
						},
					},
				},
			},
			want: metav1.OwnerReference{
				UID:  "the-replicaset-uid",
				Kind: "ReplicaSet",
				Name: "some-name",
			},
		},
		{
			name: "Test return is input argument if owner is not found",
			args: args{
				resource: &metav1.ObjectMeta{
					UID: "the-pod-uid",
					OwnerReferences: []metav1.OwnerReference{
						{
							Kind: "SomeNonExistentType",
							UID:  "the-replicaset-uid",
						},
					},
				},
			},
			want: metav1.OwnerReference{
				UID:  "",
				Kind: "",
				Name: "",
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
			if got := a.getOwnerReference(tt.args.resource); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getOwnerReference() = %v, want %v", got, tt.want)
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
							apicommon.AppAnnotation: "SOME-APP-NAME",
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
							apicommon.AppAnnotation: "SOME-APP-NAME",
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
							apicommon.AppAnnotation: "SOME-APP-NAME-ANNOTATION",
						},
						Labels: map[string]string{
							apicommon.AppAnnotation: "SOME-APP-NAME-LABEL",
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
							apicommon.AppAnnotation:      "SOME-APP-NAME",
							apicommon.WorkloadAnnotation: "SOME-WORKLOAD-NAME",
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
							apicommon.AppAnnotation:      "SOME-APP-NAME",
							apicommon.WorkloadAnnotation: "SOME-WORKLOAD-NAME",
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
							apicommon.AppAnnotation:      "SOME-APP-NAME-ANNOTATION",
							apicommon.WorkloadAnnotation: "SOME-WORKLOAD-NAME-ANNOTATION",
						},
						Labels: map[string]string{
							apicommon.AppAnnotation:      "SOME-APP-NAME-LABEL",
							apicommon.WorkloadAnnotation: "SOME-WORKLOAD-NAME-LABEL",
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
						apicommon.AppAnnotation: "some-app-name",
					},
				},
				primaryAnnotation:   apicommon.AppAnnotation,
				secondaryAnnotation: apicommon.K8sRecommendedAppAnnotations,
			},
			want:  "some-app-name",
			want1: true,
		},
		{
			name: "Test if secondary annotation is returned from annotations",
			args: args{
				resource: &metav1.ObjectMeta{
					Annotations: map[string]string{
						apicommon.K8sRecommendedAppAnnotations: "some-app-name",
					},
				},
				primaryAnnotation:   apicommon.AppAnnotation,
				secondaryAnnotation: apicommon.K8sRecommendedAppAnnotations,
			},
			want:  "some-app-name",
			want1: true,
		},
		{
			name: "Test if primary annotation is returned from labels",
			args: args{
				resource: &metav1.ObjectMeta{
					Labels: map[string]string{
						apicommon.AppAnnotation: "some-app-name",
					},
				},
				primaryAnnotation:   apicommon.AppAnnotation,
				secondaryAnnotation: apicommon.K8sRecommendedAppAnnotations,
			},
			want:  "some-app-name",
			want1: true,
		},
		{
			name: "Test if secondary annotation is returned from labels",
			args: args{
				resource: &metav1.ObjectMeta{
					Labels: map[string]string{
						apicommon.K8sRecommendedAppAnnotations: "some-app-name",
					},
				},
				primaryAnnotation:   apicommon.AppAnnotation,
				secondaryAnnotation: apicommon.K8sRecommendedAppAnnotations,
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
				primaryAnnotation:   apicommon.AppAnnotation,
				secondaryAnnotation: apicommon.K8sRecommendedAppAnnotations,
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
				primaryAnnotation:   apicommon.AppAnnotation,
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
							apicommon.AppAnnotation:      "SOME-APP-NAME-ANNOTATION",
							apicommon.WorkloadAnnotation: "workload-name-that-is-too-loooooooooooooooooooooooooooooooooooooooooooooooooong",
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
							apicommon.WorkloadAnnotation: "some-workload-name",
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
							apicommon.WorkloadAnnotation: "some-workload-name",
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
						apicommon.WorkloadAnnotation: "some-workload-name",
					},
					Annotations: map[string]string{
						apicommon.VersionAnnotation: "v1",
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
	rsUidWithNoOwner := types.UID("this-is-the-replicaset-with-no-owner")
	testStsUid := types.UID("this-is-the-stateful-set-uid")
	tstStsName := "test-stateful-set"
	testDsUid := types.UID("this-is-the-daemon-set-uid")
	testDsName := "test-daemon-set"

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
	// TODO: fix tests where an RS has a STS or DS as owner. they should not have a RS in between
	rsWithNoOwner := &appsv1.ReplicaSet{
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
			Name:      tstStsName,
			UID:       testStsUid,
			Namespace: testNamespace,
		},
	}
	testDs := &appsv1.DaemonSet{
		TypeMeta: metav1.TypeMeta{
			Kind: "DaemonSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      testDsName,
			UID:       testDsUid,
			Namespace: testNamespace,
		},
	}

	fakeClient := fakeclient.NewClient(rsWithDpOwner, rsWithNoOwner, testDp, testSts, testDs)

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
						UID: "some-uid",
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
								Name: rsWithDpOwner.Name,
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
			name: "Test fetching of statefulset owner of pod",
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
								Name: testSts.Name,
								UID:  testSts.UID,
								Kind: testSts.Kind,
							},
						},
					},
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Test fetching of daemonset owner of pod",
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
								Name: testDs.Name,
								UID:  testDs.UID,
								Kind: testDs.Kind,
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
								Name: rsWithNoOwner.Name,
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
						apicommon.WorkloadAnnotation:                 "some-workload-name",
						apicommon.AppAnnotation:                      "some-app-name",
						apicommon.VersionAnnotation:                  "v1.0.0",
						apicommon.PreDeploymentTaskAnnotation:        "some-pre-deployment-task",
						apicommon.PostDeploymentTaskAnnotation:       "some-post-deployment-task",
						apicommon.PreDeploymentEvaluationAnnotation:  "some-pre-deployment-evaluation",
						apicommon.PostDeploymentEvaluationAnnotation: "some-post-deployment-evaluation",
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
						apicommon.WorkloadAnnotation:                 "some-workload-name",
						apicommon.AppAnnotation:                      "some-app-name",
						apicommon.VersionAnnotation:                  "v1.0.0",
						apicommon.PreDeploymentTaskAnnotation:        "some-pre-deployment-task",
						apicommon.PostDeploymentTaskAnnotation:       "some-post-deployment-task",
						apicommon.PreDeploymentEvaluationAnnotation:  "some-pre-deployment-evaluation",
						apicommon.PostDeploymentEvaluationAnnotation: "some-post-deployment-evaluation",
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
						apicommon.WorkloadAnnotation:                 "some-workload-name",
						apicommon.AppAnnotation:                      "some-app-name",
						apicommon.VersionAnnotation:                  "v1.0.0",
						apicommon.PreDeploymentTaskAnnotation:        "some-pre-deployment-task",
						apicommon.PostDeploymentTaskAnnotation:       "some-post-deployment-task",
						apicommon.PreDeploymentEvaluationAnnotation:  "some-pre-deployment-evaluation",
						apicommon.PostDeploymentEvaluationAnnotation: "some-post-deployment-evaluation",
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
						apicommon.WorkloadAnnotation:                 "some-workload-name",
						apicommon.AppAnnotation:                      "some-app-name",
						apicommon.VersionAnnotation:                  "v1.0.0",
						apicommon.PreDeploymentTaskAnnotation:        "some-pre-deployment-task",
						apicommon.PostDeploymentTaskAnnotation:       "some-post-deployment-task",
						apicommon.PreDeploymentEvaluationAnnotation:  "some-pre-deployment-evaluation",
						apicommon.PostDeploymentEvaluationAnnotation: "some-post-deployment-evaluation",
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
						apicommon.WorkloadAnnotation:                 "some-workload-name",
						apicommon.AppAnnotation:                      "some-app-name",
						apicommon.PreDeploymentTaskAnnotation:        "some-pre-deployment-task",
						apicommon.PostDeploymentTaskAnnotation:       "some-post-deployment-task",
						apicommon.PreDeploymentEvaluationAnnotation:  "some-pre-deployment-evaluation",
						apicommon.PostDeploymentEvaluationAnnotation: "some-post-deployment-evaluation",
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
						apicommon.WorkloadAnnotation:                 "some-workload-name",
						apicommon.AppAnnotation:                      "some-app-name",
						apicommon.VersionAnnotation:                  "v1.0.0",
						apicommon.PreDeploymentTaskAnnotation:        "some-pre-deployment-task",
						apicommon.PostDeploymentTaskAnnotation:       "some-post-deployment-task",
						apicommon.PreDeploymentEvaluationAnnotation:  "some-pre-deployment-evaluation",
						apicommon.PostDeploymentEvaluationAnnotation: "some-post-deployment-evaluation",
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
						apicommon.WorkloadAnnotation: "some-workload-name-that-is-very-looooooooooooooooooooooong",
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

func TestPodMutatingWebhook_isAppAnnotationPresent(t *testing.T) {
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
			name: "Test return true when app annotation is present",
			args: args{
				pod: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							apicommon.AppAnnotation: "some-app-name",
						},
					},
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Test return false when app annotation is not present",
			args: args{
				pod: &corev1.Pod{},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Test return error when app annotation is too long",
			args: args{
				pod: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							apicommon.AppAnnotation: "some-app-annotation-that-is-very-looooooooooooooooooooong",
						},
					},
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Test that app name is copied when only workload name is present",
			args: args{
				pod: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							apicommon.WorkloadAnnotation: "some-workload-name",
						},
					},
				},
			},
			want:    false,
			wantErr: false,
			wantedPod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						apicommon.AppAnnotation:      "some-workload-name",
						apicommon.WorkloadAnnotation: "some-workload-name",
					},
				},
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
			got, err := a.isAppAnnotationPresent(tt.args.pod)
			if (err != nil) != tt.wantErr {
				t.Errorf("isAppAnnotationPresent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("isAppAnnotationPresent() got = %v, want %v", got, tt.want)
			}
			if tt.wantedPod != nil {
				require.Equal(t, tt.wantedPod, tt.args.pod)
			}
		})
	}
}

func TestPodMutatingWebhook_Handle_DisabledNamespace(t *testing.T) {
	fakeClient := fakeclient.NewClient(&corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "default",
		},
	})

	tr := &fakeclient.ITracerMock{StartFunc: func(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
		return ctx, trace.SpanFromContext(ctx)
	}}

	recorder := record.NewFakeRecorder(100)

	decoder, err := admission.NewDecoder(runtime.NewScheme())
	require.Nil(t, err)

	wh := &PodMutatingWebhook{
		Client:   fakeClient,
		Tracer:   tr,
		decoder:  decoder,
		Recorder: recorder,
		Log:      testr.New(t),
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example-pod",
			Namespace: "default",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "example-container",
					Image: "nginx",
				},
			},
		},
	}

	// Convert the Pod object to a byte array
	podBytes, err := json.Marshal(pod)
	require.Nil(t, err)

	// Create an AdmissionRequest object
	request := admissionv1.AdmissionRequest{
		UID:       "12345",
		Kind:      metav1.GroupVersionKind{Group: "", Version: "v1", Kind: "Pod"},
		Operation: admissionv1.Create,
		Object: runtime.RawExtension{
			Raw: podBytes,
		},
		Namespace: "default",
	}

	resp := wh.Handle(context.TODO(), admission.Request{
		AdmissionRequest: request,
	})

	require.NotNil(t, resp)
	require.True(t, resp.Allowed)
}

func TestPodMutatingWebhook_Handle_SingleService(t *testing.T) {
	fakeClient := fakeclient.NewClient(&corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "default",
			Annotations: map[string]string{
				apicommon.NamespaceEnabledAnnotation: "enabled",
			},
		},
	})

	tr := &fakeclient.ITracerMock{StartFunc: func(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
		return ctx, trace.SpanFromContext(ctx)
	}}

	recorder := record.NewFakeRecorder(100)

	decoder, err := admission.NewDecoder(runtime.NewScheme())
	require.Nil(t, err)

	wh := &PodMutatingWebhook{
		Client:   fakeClient,
		Tracer:   tr,
		decoder:  decoder,
		Recorder: recorder,
		Log:      testr.New(t),
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example-pod",
			Namespace: "default",
			Annotations: map[string]string{
				apicommon.WorkloadAnnotation: "my-workload",
				apicommon.VersionAnnotation:  "0.1",
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: "v1",
					Kind:       "Deployment",
					Name:       "my-deployment",
					UID:        "1234",
				},
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "example-container",
					Image: "nginx",
				},
			},
		},
	}

	// Convert the Pod object to a byte array
	podBytes, err := json.Marshal(pod)
	require.Nil(t, err)

	// Create an AdmissionRequest object
	request := admissionv1.AdmissionRequest{
		UID:       "12345",
		Kind:      metav1.GroupVersionKind{Group: "", Version: "v1", Kind: "Pod"},
		Operation: admissionv1.Create,
		Object: runtime.RawExtension{
			Raw: podBytes,
		},
		Namespace: "default",
	}

	resp := wh.Handle(context.TODO(), admission.Request{
		AdmissionRequest: request,
	})

	require.NotNil(t, resp)
	require.True(t, resp.Allowed)

	kacr := &klcv1alpha3.KeptnAppCreationRequest{}

	err = fakeClient.Get(context.Background(), types.NamespacedName{
		Namespace: "default",
		Name:      "my-workload",
	}, kacr)

	require.Nil(t, err)

	require.Equal(t, "my-workload", kacr.Spec.AppName)
	require.Equal(t, string(apicommon.AppTypeSingleService), kacr.Annotations[apicommon.AppTypeAnnotation])

	workload := &klcv1alpha3.KeptnWorkload{}

	err = fakeClient.Get(context.TODO(), types.NamespacedName{
		Namespace: "default",
		Name:      "my-workload-my-workload",
	}, workload)

	require.Nil(t, err)

	require.Equal(t, klcv1alpha3.KeptnWorkloadSpec{
		AppName: kacr.Spec.AppName,
		Version: "0.1",
		ResourceReference: klcv1alpha3.ResourceReference{
			UID:  "1234",
			Kind: "Deployment",
			Name: "my-deployment",
		},
	}, workload.Spec)
}

func TestPodMutatingWebhook_Handle_SingleService_AppCreationRequestAlreadyPresent(t *testing.T) {
	fakeClient := fakeclient.NewClient(&corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "default",
			Annotations: map[string]string{
				apicommon.NamespaceEnabledAnnotation: "enabled",
			},
		},
	}, &klcv1alpha3.KeptnAppCreationRequest{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-workload",
			Namespace: "default",
			Annotations: map[string]string{
				apicommon.AppTypeAnnotation: string(apicommon.AppTypeSingleService),
			},
			Labels: map[string]string{
				"donotchange": "true",
			},
		},
		Spec: klcv1alpha3.KeptnAppCreationRequestSpec{
			AppName: "my-workload",
		},
	})

	tr := &fakeclient.ITracerMock{StartFunc: func(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
		return ctx, trace.SpanFromContext(ctx)
	}}

	recorder := record.NewFakeRecorder(100)

	decoder, err := admission.NewDecoder(runtime.NewScheme())
	require.Nil(t, err)

	wh := &PodMutatingWebhook{
		Client:   fakeClient,
		Tracer:   tr,
		decoder:  decoder,
		Recorder: recorder,
		Log:      testr.New(t),
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example-pod",
			Namespace: "default",
			Annotations: map[string]string{
				apicommon.WorkloadAnnotation: "my-workload",
				apicommon.VersionAnnotation:  "0.1",
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: "v1",
					Kind:       "Deployment",
					Name:       "my-deployment",
					UID:        "1234",
				},
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "example-container",
					Image: "nginx",
				},
			},
		},
	}

	// Convert the Pod object to a byte array
	podBytes, err := json.Marshal(pod)
	require.Nil(t, err)

	// Create an AdmissionRequest object
	request := admissionv1.AdmissionRequest{
		UID:       "12345",
		Kind:      metav1.GroupVersionKind{Group: "", Version: "v1", Kind: "Pod"},
		Operation: admissionv1.Create,
		Object: runtime.RawExtension{
			Raw: podBytes,
		},
		Namespace: "default",
	}

	resp := wh.Handle(context.TODO(), admission.Request{
		AdmissionRequest: request,
	})

	require.NotNil(t, resp)
	require.True(t, resp.Allowed)

	kacr := &klcv1alpha3.KeptnAppCreationRequest{}

	err = fakeClient.Get(context.Background(), types.NamespacedName{
		Namespace: "default",
		Name:      "my-workload",
	}, kacr)

	require.Nil(t, err)

	require.Equal(t, "my-workload", kacr.Spec.AppName)
	require.Equal(t, string(apicommon.AppTypeSingleService), kacr.Annotations[apicommon.AppTypeAnnotation])
	// verify that the previously created KACR has not been changed
	require.Equal(t, "true", kacr.Labels["donotchange"])

	workload := &klcv1alpha3.KeptnWorkload{}

	err = fakeClient.Get(context.TODO(), types.NamespacedName{
		Namespace: "default",
		Name:      "my-workload-my-workload",
	}, workload)

	require.Nil(t, err)

	require.Equal(t, klcv1alpha3.KeptnWorkloadSpec{
		AppName: kacr.Spec.AppName,
		Version: "0.1",
		ResourceReference: klcv1alpha3.ResourceReference{
			UID:  "1234",
			Kind: "Deployment",
			Name: "my-deployment",
		},
	}, workload.Spec)
}

func TestPodMutatingWebhook_Handle_MultiService(t *testing.T) {
	fakeClient := fakeclient.NewClient(&corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "default",
			Annotations: map[string]string{
				apicommon.NamespaceEnabledAnnotation: "enabled",
			},
		},
	})

	tr := &fakeclient.ITracerMock{StartFunc: func(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
		return ctx, trace.SpanFromContext(ctx)
	}}

	recorder := record.NewFakeRecorder(100)

	decoder, err := admission.NewDecoder(runtime.NewScheme())
	require.Nil(t, err)

	wh := &PodMutatingWebhook{
		Client:   fakeClient,
		Tracer:   tr,
		decoder:  decoder,
		Recorder: recorder,
		Log:      testr.New(t),
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example-pod",
			Namespace: "default",
			Annotations: map[string]string{
				apicommon.WorkloadAnnotation: "my-workload",
				apicommon.VersionAnnotation:  "0.1",
				apicommon.AppAnnotation:      "my-app",
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: "v1",
					Kind:       "Deployment",
					Name:       "my-deployment",
					UID:        "1234",
				},
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "example-container",
					Image: "nginx",
				},
			},
		},
	}

	// Convert the Pod object to a byte array
	podBytes, err := json.Marshal(pod)
	require.Nil(t, err)

	// Create an AdmissionRequest object
	request := admissionv1.AdmissionRequest{
		UID:       "12345",
		Kind:      metav1.GroupVersionKind{Group: "", Version: "v1", Kind: "Pod"},
		Operation: admissionv1.Create,
		Object: runtime.RawExtension{
			Raw: podBytes,
		},
		Namespace: "default",
	}

	resp := wh.Handle(context.TODO(), admission.Request{
		AdmissionRequest: request,
	})

	require.NotNil(t, resp)
	require.True(t, resp.Allowed)

	kacr := &klcv1alpha3.KeptnAppCreationRequest{}

	err = fakeClient.Get(context.Background(), types.NamespacedName{
		Namespace: "default",
		Name:      "my-app",
	}, kacr)

	require.Nil(t, err)

	require.Equal(t, "my-app", kacr.Spec.AppName)
	// here we do not want a single-service annotation
	require.Empty(t, kacr.Annotations[apicommon.AppTypeAnnotation])

	workload := &klcv1alpha3.KeptnWorkload{}

	err = fakeClient.Get(context.TODO(), types.NamespacedName{
		Namespace: "default",
		Name:      "my-app-my-workload",
	}, workload)

	require.Nil(t, err)

	require.Equal(t, klcv1alpha3.KeptnWorkloadSpec{
		AppName: kacr.Spec.AppName,
		Version: "0.1",
		ResourceReference: klcv1alpha3.ResourceReference{
			UID:  "1234",
			Kind: "Deployment",
			Name: "my-deployment",
		},
	}, workload.Spec)
}
