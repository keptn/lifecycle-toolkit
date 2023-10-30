package pod_mutator

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/go-logr/logr/testr"
	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	controllercommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common"
	fakeclient "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/fake"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/webhooks/pod_mutator/handlers"
	fakehandler "github.com/keptn/lifecycle-toolkit/lifecycle-operator/webhooks/pod_mutator/handlers/fake"
	"github.com/stretchr/testify/require"
	admissionv1 "k8s.io/api/admission/v1"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	k8sfake "sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

const testWorkload = "my-workload"
const testPod = "example-pod"
const testNamespace = "default"
const testKeptnWorkload = "my-workload-my-workload"
const testDeployment = "my-deployment"

func TestPodMutatingWebhookHandleDisabledNamespace(t *testing.T) {
	fakeClient := fakeclient.NewClient(&corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: testNamespace,
		},
	})

	decoder := admission.NewDecoder(runtime.NewScheme())

	wh := &PodMutatingWebhook{
		Client:      fakeClient,
		Decoder:     decoder,
		EventSender: controllercommon.NewK8sSender(record.NewFakeRecorder(100)),
		Log:         testr.New(t),
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testPod,
			Namespace: testNamespace,
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

	request := generateRequest(pod, t)

	resp := wh.Handle(context.TODO(), admission.Request{
		AdmissionRequest: request,
	})

	require.NotNil(t, resp)
	require.True(t, resp.Allowed)
}

func TestPodMutatingWebhookHandleUnsupportedOwner(t *testing.T) {
	fakeClient := fakeclient.NewClient(&corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: testNamespace,
			Annotations: map[string]string{
				apicommon.NamespaceEnabledAnnotation: "enabled",
			},
		},
	})

	decoder := admission.NewDecoder(runtime.NewScheme())

	wh := &PodMutatingWebhook{
		Client:      fakeClient,
		Decoder:     decoder,
		EventSender: controllercommon.NewK8sSender(record.NewFakeRecorder(100)),
		Log:         testr.New(t),
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testPod,
			Namespace: testNamespace,
			Annotations: map[string]string{
				apicommon.WorkloadAnnotation: testWorkload,
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

	request := generateRequest(pod, t)

	resp := wh.Handle(context.TODO(), admission.Request{
		AdmissionRequest: request,
	})

	require.NotNil(t, resp)
	require.True(t, resp.Allowed)

	// if we get an unsupported owner for the pod, we expect not to have any KLT resources to have been created
	kacr := &klcv1alpha3.KeptnAppCreationRequest{}

	err := fakeClient.Get(context.Background(), types.NamespacedName{
		Namespace: testNamespace,
		Name:      testWorkload,
	}, kacr)

	require.NotNil(t, err)
	require.True(t, k8serrors.IsNotFound(err))

	workload := &klcv1alpha3.KeptnWorkload{}

	err = fakeClient.Get(context.TODO(), types.NamespacedName{
		Namespace: testNamespace,
		Name:      testKeptnWorkload,
	}, workload)

	require.NotNil(t, err)
	require.True(t, k8serrors.IsNotFound(err))
}

func TestPodMutatingWebhookHandleSingleService(t *testing.T) {
	fakeClient := fakeclient.NewClient(&corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: testNamespace,
			Annotations: map[string]string{
				apicommon.NamespaceEnabledAnnotation: "enabled",
			},
		},
	})

	decoder := admission.NewDecoder(runtime.NewScheme())
	log := testr.New(t)

	wh := NewPodMutator(fakeClient, decoder, controllercommon.NewK8sSender(record.NewFakeRecorder(100)), log, false)

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testPod,
			Namespace: testNamespace,
			Annotations: map[string]string{
				apicommon.WorkloadAnnotation: testWorkload,
				apicommon.VersionAnnotation:  "0.1",
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: "v1",
					Kind:       "Deployment",
					Name:       testDeployment,
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

	request := generateRequest(pod, t)

	resp := wh.Handle(context.TODO(), admission.Request{
		AdmissionRequest: request,
	})

	require.NotNil(t, resp)
	require.True(t, resp.Allowed)

	kacr := &klcv1alpha3.KeptnAppCreationRequest{}

	err := fakeClient.Get(context.Background(), types.NamespacedName{
		Namespace: testNamespace,
		Name:      testWorkload,
	}, kacr)

	require.Nil(t, err)

	require.Equal(t, testWorkload, kacr.Spec.AppName)
	require.Equal(t, string(apicommon.AppTypeSingleService), kacr.Annotations[apicommon.AppTypeAnnotation])

	workload := &klcv1alpha3.KeptnWorkload{}

	err = fakeClient.Get(context.TODO(), types.NamespacedName{
		Namespace: testNamespace,
		Name:      testKeptnWorkload,
	}, workload)

	require.Nil(t, err)

	require.Equal(t, klcv1alpha3.KeptnWorkloadSpec{
		AppName: kacr.Spec.AppName,
		Version: "0.1",
		ResourceReference: klcv1alpha3.ResourceReference{
			UID:  "1234",
			Kind: "Deployment",
			Name: testDeployment,
		},
	}, workload.Spec)
}

func TestPodMutatingWebhookHandleSchedulingGatesGateRemoved(t *testing.T) {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testPod,
			Namespace: testNamespace,
			Annotations: map[string]string{
				apicommon.WorkloadAnnotation:    testWorkload,
				apicommon.VersionAnnotation:     "0.1",
				apicommon.SchedulingGateRemoved: "true",
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: "v1",
					Kind:       "Deployment",
					Name:       testDeployment,
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
			Name: testNamespace,
			Annotations: map[string]string{
				apicommon.NamespaceEnabledAnnotation: "enabled",
			},
		},
	}
	fakeClient := fakeclient.NewClient(ns, pod)

	decoder := admission.NewDecoder(runtime.NewScheme())

	wh := &PodMutatingWebhook{
		SchedulingGatesEnabled: true,
		Client:                 fakeClient,
		Decoder:                decoder,
		EventSender:            controllercommon.NewK8sSender(record.NewFakeRecorder(100)),
		Log:                    testr.New(t),
	}

	request := generateRequest(pod, t)

	resp := wh.Handle(context.TODO(), admission.Request{
		AdmissionRequest: request,
	})

	require.NotNil(t, resp)
	require.True(t, resp.Allowed)

	// no changes to the pod are expected
	require.Len(t, resp.Patches, 0)
}

func TestPodMutatingWebhookHandleSchedulingGates(t *testing.T) {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testPod,
			Namespace: testNamespace,
			Annotations: map[string]string{
				apicommon.WorkloadAnnotation: testWorkload,
				apicommon.VersionAnnotation:  "0.1",
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: "v1",
					Kind:       "Deployment",
					Name:       testDeployment,
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
			Name: testNamespace,
			Annotations: map[string]string{
				apicommon.NamespaceEnabledAnnotation: "enabled",
			},
		},
	}
	fakeClient := fakeclient.NewClient(ns, pod)

	decoder := admission.NewDecoder(runtime.NewScheme())

	wh :=
		NewPodMutator(fakeClient, decoder, controllercommon.NewK8sSender(record.NewFakeRecorder(100)), testr.New(t), true)

	request := generateRequest(pod, t)

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

	err := fakeClient.Get(context.Background(), types.NamespacedName{
		Namespace: testNamespace,
		Name:      testWorkload,
	}, kacr)

	require.Nil(t, err)

	require.Equal(t, testWorkload, kacr.Spec.AppName)
	require.Equal(t, string(apicommon.AppTypeSingleService), kacr.Annotations[apicommon.AppTypeAnnotation])

	workload := &klcv1alpha3.KeptnWorkload{}

	err = fakeClient.Get(context.TODO(), types.NamespacedName{
		Namespace: testNamespace,
		Name:      testKeptnWorkload,
	}, workload)

	require.Nil(t, err)

	require.Equal(t, klcv1alpha3.KeptnWorkloadSpec{
		AppName: kacr.Spec.AppName,
		Version: "0.1",
		ResourceReference: klcv1alpha3.ResourceReference{
			UID:  "1234",
			Kind: "Deployment",
			Name: testDeployment,
		},
	}, workload.Spec)
}

func TestPodMutatingWebhookHandleSingleServiceAppCreationRequestAlreadyPresent(t *testing.T) {
	fakeClient := fakeclient.NewClient(&corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: testNamespace,
			Annotations: map[string]string{
				apicommon.NamespaceEnabledAnnotation: "enabled",
			},
		},
	}, &klcv1alpha3.KeptnAppCreationRequest{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      testWorkload,
			Namespace: testNamespace,
			Annotations: map[string]string{
				apicommon.AppTypeAnnotation: string(apicommon.AppTypeSingleService),
			},
			Labels: map[string]string{
				"donotchange": "true",
			},
		},
		Spec: klcv1alpha3.KeptnAppCreationRequestSpec{
			AppName: testWorkload,
		},
	})

	decoder := admission.NewDecoder(runtime.NewScheme())

	wh := NewPodMutator(fakeClient, decoder, controllercommon.NewK8sSender(record.NewFakeRecorder(100)), testr.New(t), false)

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testPod,
			Namespace: testNamespace,
			Annotations: map[string]string{
				apicommon.WorkloadAnnotation: testWorkload,
				apicommon.VersionAnnotation:  "0.1",
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: "v1",
					Kind:       "Deployment",
					Name:       testDeployment,
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

	request := generateRequest(pod, t)

	resp := wh.Handle(context.TODO(), admission.Request{
		AdmissionRequest: request,
	})

	require.NotNil(t, resp)
	require.True(t, resp.Allowed)

	kacr := &klcv1alpha3.KeptnAppCreationRequest{}

	err := fakeClient.Get(context.Background(), types.NamespacedName{
		Namespace: testNamespace,
		Name:      testWorkload,
	}, kacr)

	require.Nil(t, err)

	require.Equal(t, testWorkload, kacr.Spec.AppName)
	require.Equal(t, string(apicommon.AppTypeSingleService), kacr.Annotations[apicommon.AppTypeAnnotation])
	// verify that the previously created KACR has not been changed
	require.Equal(t, "true", kacr.Labels["donotchange"])

	workload := &klcv1alpha3.KeptnWorkload{}

	err = fakeClient.Get(context.TODO(), types.NamespacedName{
		Namespace: testNamespace,
		Name:      testKeptnWorkload,
	}, workload)

	require.Nil(t, err)

	require.Equal(t, klcv1alpha3.KeptnWorkloadSpec{
		AppName: kacr.Spec.AppName,
		Version: "0.1",
		ResourceReference: klcv1alpha3.ResourceReference{
			UID:  "1234",
			Kind: "Deployment",
			Name: testDeployment,
		},
	}, workload.Spec)
}

func TestPodMutatingWebhookHandleMultiService(t *testing.T) {
	fakeClient := fakeclient.NewClient(&corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: testNamespace,
			Annotations: map[string]string{
				apicommon.NamespaceEnabledAnnotation: "enabled",
			},
		},
	})

	pod, _, _, decoder := setupTestData()

	wh := NewPodMutator(fakeClient, decoder, controllercommon.NewK8sSender(record.NewFakeRecorder(100)), testr.New(t), false)

	request := generateRequest(pod, t)

	resp := wh.Handle(context.TODO(), admission.Request{
		AdmissionRequest: request,
	})

	require.NotNil(t, resp)
	require.True(t, resp.Allowed)

	kacr := &klcv1alpha3.KeptnAppCreationRequest{}

	err := fakeClient.Get(context.Background(), types.NamespacedName{
		Namespace: testNamespace,
		Name:      "my-app",
	}, kacr)

	require.Nil(t, err)

	require.Equal(t, "my-app", kacr.Spec.AppName) // this makes sure that everything is lowercase
	// here we do not want a single-service annotation
	require.Empty(t, kacr.Annotations[apicommon.AppTypeAnnotation])

	workload := &klcv1alpha3.KeptnWorkload{}

	err = fakeClient.Get(context.TODO(), types.NamespacedName{
		Namespace: testNamespace,
		Name:      "my-app-my-workload",
	}, workload)

	require.Nil(t, err)

	require.Equal(t, klcv1alpha3.KeptnWorkloadSpec{
		AppName: kacr.Spec.AppName,
		Version: "v0.1",
		ResourceReference: klcv1alpha3.ResourceReference{
			UID:  "1234",
			Kind: "Deployment",
			Name: testDeployment,
		},
	}, workload.Spec)
}

func TestPodMutatingWebhookHandleErrorPaths(t *testing.T) {

	pod, dp, ns, decoder := setupTestData()

	tests := []struct {
		name            string
		workloadHandler handlers.K8sHandler
		appHandler      handlers.K8sHandler
		client          client.Client
		decoder         handlers.Decoder
		message         string
		errorCode       int
	}{
		{
			name:    "DecoderError",
			message: "bad decode",
			decoder: &fakehandler.MockDecoder{DecodeFunc: func(req admission.Request, into runtime.Object) error {
				return k8serrors.NewResourceExpired("bad decode")
			}},
			errorCode: http.StatusBadRequest,
		},
		{
			name:    "NamespaceError",
			message: "could not get",
			client: k8sfake.NewClientBuilder().WithInterceptorFuncs(interceptor.Funcs{
				Get: func(ctx context.Context, client client.WithWatch, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
					return errors.New("could not get")
				},
			}).Build(),
			decoder: &fakehandler.MockDecoder{
				DecodeFunc: func(req admission.Request, into runtime.Object) error {
					return nil
				}},
			errorCode: http.StatusInternalServerError,
		},
		{
			name: "WorkloadError",
			workloadHandler: &fakehandler.MockHandler{
				HandleFunc: func(ctx context.Context, pod *corev1.Pod, namespace string) error {
					return errors.New("bad workload")
				},
			},
			message:   "bad workload",
			client:    fakeclient.NewClient(pod, dp, ns),
			decoder:   decoder,
			errorCode: http.StatusBadRequest,
		},
		{
			name: "AppError",
			workloadHandler: &fakehandler.MockHandler{
				HandleFunc: func(ctx context.Context, pod *corev1.Pod, namespace string) error {
					return nil
				},
			},
			appHandler: &fakehandler.MockHandler{
				HandleFunc: func(ctx context.Context, pod *corev1.Pod, namespace string) error {
					return errors.New("bad app")
				},
			},
			message:   "bad app",
			decoder:   decoder,
			client:    fakeclient.NewClient(pod, dp, ns),
			errorCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			wh := PodMutatingWebhook{
				Decoder:  tt.decoder,
				Log:      testr.New(t),
				Client:   tt.client,
				Workload: tt.workloadHandler,
				App:      tt.appHandler,
			}

			// Create an AdmissionRequest object
			request := generateRequest(pod, t)

			resp := wh.Handle(context.TODO(), admission.Request{
				AdmissionRequest: request,
			})

			require.NotNil(t, resp)
			require.False(t, resp.Allowed)
			require.Equal(t, tt.message, resp.Result.Message)
			require.Equal(t, tt.errorCode, int(resp.Result.Code))
		})
	}

}

func generateRequest(pod *corev1.Pod, t *testing.T) admissionv1.AdmissionRequest {
	// Convert the Pod object to a byte array
	podBytes, err := json.Marshal(pod)
	require.Nil(t, err)

	return admissionv1.AdmissionRequest{
		UID:       "12345",
		Kind:      metav1.GroupVersionKind{Group: "", Version: "v1", Kind: "Pod"},
		Operation: admissionv1.Create,
		Object: runtime.RawExtension{
			Raw: podBytes,
		},
		Namespace: testNamespace,
	}
}

func setupTestData() (*corev1.Pod, *v1.Deployment, *corev1.Namespace, *admission.Decoder) {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testPod,
			Namespace: testNamespace,
			Annotations: map[string]string{
				apicommon.WorkloadAnnotation: testWorkload,
				apicommon.VersionAnnotation:  "V0.1",
				apicommon.AppAnnotation:      "my-App",
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: "v1",
					Kind:       "Deployment",
					Name:       testDeployment,
					UID:        "1234",
				},
			},
		},
	}

	dp := &v1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind: "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: testDeployment, Namespace: testNamespace,
			UID: "1234"},
	}

	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: testNamespace,
			Annotations: map[string]string{
				apicommon.NamespaceEnabledAnnotation: "enabled",
			},
		},
	}

	decoder := admission.NewDecoder(runtime.NewScheme())
	return pod, dp, ns, decoder
}
