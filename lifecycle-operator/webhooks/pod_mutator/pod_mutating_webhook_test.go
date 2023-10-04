package pod_mutator

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/go-logr/logr/testr"
	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	controllercommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common"
	fakeclient "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/fake"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

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
