package pod_mutator

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	controllercommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common"
	fakeclient "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/fake"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
	admissionv1 "k8s.io/api/admission/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

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
		Client      client.Client
		Tracer      trace.Tracer
		Decoder     *admission.Decoder
		EventSender controllercommon.IEvent
		Log         logr.Logger
	}
	type args struct {
		ctx context.Context
		req *admission.Request
		pod *corev1.Pod
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
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
			want: false,
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
			want: false,
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
			want: false,
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
			want: false,
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
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &PodMutatingWebhook{
				Client:      tt.fields.Client,
				Tracer:      tt.fields.Tracer,
				Decoder:     tt.fields.Decoder,
				EventSender: tt.fields.EventSender,
				Log:         tt.fields.Log,
			}
			got := a.copyAnnotationsIfParentAnnotated(tt.args.ctx, tt.args.req, tt.args.pod)
			if got != tt.want {
				t.Errorf("copyAnnotationsIfParentAnnotated() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPodMutatingWebhook_copyResourceLabelsIfPresent(t *testing.T) {

	type args struct {
		sourceResource *metav1.ObjectMeta
		targetPod      *corev1.Pod
	}
	tests := []struct {
		name      string
		args      args
		want      bool
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
			want: true,
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
			want: true,
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
			want: true,
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := copyResourceLabelsIfPresent(tt.args.sourceResource, tt.args.targetPod)
			if got != tt.want {
				t.Errorf("copyResourceLabelsIfPresent() got = %v, want %v", got, tt.want)
			}
			if tt.wantedPod != nil {
				require.Equal(t, tt.wantedPod, tt.args.targetPod)
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

	decoder := admission.NewDecoder(runtime.NewScheme())

	wh := &PodMutatingWebhook{
		Client:      fakeClient,
		Tracer:      tr,
		Decoder:     decoder,
		EventSender: controllercommon.NewK8sSender(record.NewFakeRecorder(100)),
		Log:         testr.New(t),
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

func TestPodMutatingWebhook_Handle_UnsupportedOwner(t *testing.T) {
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

	decoder := admission.NewDecoder(runtime.NewScheme())

	wh := &PodMutatingWebhook{
		Client:      fakeClient,
		Tracer:      tr,
		Decoder:     decoder,
		EventSender: controllercommon.NewK8sSender(record.NewFakeRecorder(100)),
		Log:         testr.New(t),
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
					APIVersion: "batchv1",
					Kind:       "Job",
					Name:       "my-job",
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

	// if we get an unsupported owner for the pod, we expect not to have any KLT resources to have been created
	kacr := &klcv1alpha3.KeptnAppCreationRequest{}

	err = fakeClient.Get(context.Background(), types.NamespacedName{
		Namespace: "default",
		Name:      "my-workload",
	}, kacr)

	require.NotNil(t, err)
	require.True(t, errors.IsNotFound(err))

	workload := &klcv1alpha3.KeptnWorkload{}

	err = fakeClient.Get(context.TODO(), types.NamespacedName{
		Namespace: "default",
		Name:      "my-workload-my-workload",
	}, workload)

	require.NotNil(t, err)
	require.True(t, errors.IsNotFound(err))
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

	decoder := admission.NewDecoder(runtime.NewScheme())

	wh := &PodMutatingWebhook{
		Client:      fakeClient,
		Tracer:      tr,
		Decoder:     decoder,
		EventSender: controllercommon.NewK8sSender(record.NewFakeRecorder(100)),
		Log:         testr.New(t),
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

func TestPodMutatingWebhook_Handle_SchedulingGates_GateRemoved(t *testing.T) {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example-pod",
			Namespace: "default",
			Annotations: map[string]string{
				apicommon.WorkloadAnnotation:    "my-workload",
				apicommon.VersionAnnotation:     "0.1",
				apicommon.SchedulingGateRemoved: "true",
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
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "default",
			Annotations: map[string]string{
				apicommon.NamespaceEnabledAnnotation: "enabled",
			},
		},
	}
	fakeClient := fakeclient.NewClient(ns, pod)

	tr := &fakeclient.ITracerMock{StartFunc: func(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
		return ctx, trace.SpanFromContext(ctx)
	}}

	decoder := admission.NewDecoder(runtime.NewScheme())

	wh := &PodMutatingWebhook{
		SchedulingGatesEnabled: true,
		Client:                 fakeClient,
		Tracer:                 tr,
		Decoder:                decoder,
		EventSender:            controllercommon.NewK8sSender(record.NewFakeRecorder(100)),
		Log:                    testr.New(t),
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

	// no changes to the pod are expected
	require.Len(t, resp.Patches, 0)
}

func TestPodMutatingWebhook_Handle_SchedulingGates(t *testing.T) {
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
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "default",
			Annotations: map[string]string{
				apicommon.NamespaceEnabledAnnotation: "enabled",
			},
		},
	}
	fakeClient := fakeclient.NewClient(ns, pod)

	tr := &fakeclient.ITracerMock{StartFunc: func(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
		return ctx, trace.SpanFromContext(ctx)
	}}

	decoder := admission.NewDecoder(runtime.NewScheme())

	wh := &PodMutatingWebhook{
		SchedulingGatesEnabled: true,
		Client:                 fakeClient,
		Tracer:                 tr,
		Decoder:                decoder,
		EventSender:            controllercommon.NewK8sSender(record.NewFakeRecorder(100)),
		Log:                    testr.New(t),
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

	expectedValue := []interface{}{map[string]interface{}{"name": apicommon.KeptnGate}}
	require.Len(t, resp.Patches, 2)
	if resp.Patches[0].Path == "/spec/schedulingGates" {
		require.Equal(t, expectedValue, resp.Patches[0].Value)
	} else {
		require.Equal(t, expectedValue, resp.Patches[1].Value)
	}

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

	decoder := admission.NewDecoder(runtime.NewScheme())

	wh := &PodMutatingWebhook{
		Client:      fakeClient,
		Tracer:      tr,
		Decoder:     decoder,
		EventSender: controllercommon.NewK8sSender(record.NewFakeRecorder(100)),
		Log:         testr.New(t),
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

	decoder := admission.NewDecoder(runtime.NewScheme())

	wh := &PodMutatingWebhook{
		Client:      fakeClient,
		Tracer:      tr,
		Decoder:     decoder,
		EventSender: controllercommon.NewK8sSender(record.NewFakeRecorder(100)),
		Log:         testr.New(t),
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example-pod",
			Namespace: "default",
			Annotations: map[string]string{
				apicommon.WorkloadAnnotation: "my-workload",
				apicommon.VersionAnnotation:  "V0.1",
				apicommon.AppAnnotation:      "my-App",
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

	require.Equal(t, "my-app", kacr.Spec.AppName) // this makes sure that everything is lowercase
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
		Version: "v0.1",
		ResourceReference: klcv1alpha3.ResourceReference{
			UID:  "1234",
			Kind: "Deployment",
			Name: "my-deployment",
		},
	}, workload.Spec)
}
