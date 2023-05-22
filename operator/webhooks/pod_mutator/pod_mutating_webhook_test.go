package pod_mutator

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3/common"
	fakeclient "github.com/keptn/lifecycle-toolkit/operator/controllers/common/fake"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

func TestPodMutatingWebhook_getAppName(t *testing.T) {
	type fields struct {
		Client   client.Client
		Tracer   trace.Tracer
		Decoder  *admission.Decoder
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
				Decoder:  tt.fields.Decoder,
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
		Decoder  *admission.Decoder
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
				Decoder:  tt.fields.Decoder,
				Recorder: tt.fields.Recorder,
				Log:      tt.fields.Log,
			}
			if got := a.getWorkloadName(tt.args.pod); got != tt.want {
				t.Errorf("getWorkloadName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPodMutatingWebhook_isAppAnnotationPresent(t *testing.T) {
	type fields struct {
		Client   client.Client
		Tracer   trace.Tracer
		Decoder  *admission.Decoder
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
				Decoder:  tt.fields.Decoder,
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

	Decoder := admission.NewDecoder(runtime.NewScheme())

	wh := &PodMutatingWebhook{
		Client:   fakeClient,
		Tracer:   tr,
		Decoder:  Decoder,
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

	Decoder := admission.NewDecoder(runtime.NewScheme())

	wh := &PodMutatingWebhook{
		Client:   fakeClient,
		Tracer:   tr,
		Decoder:  Decoder,
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

	Decoder := admission.NewDecoder(runtime.NewScheme())

	wh := &PodMutatingWebhook{
		Client:   fakeClient,
		Tracer:   tr,
		Decoder:  Decoder,
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

	Decoder := admission.NewDecoder(runtime.NewScheme())

	wh := &PodMutatingWebhook{
		Client:   fakeClient,
		Tracer:   tr,
		Decoder:  Decoder,
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
