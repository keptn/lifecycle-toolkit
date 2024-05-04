package handlers

import (
	"context"
	"errors"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	apilifecycle "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/eventsender"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/testcommon"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	k8sfake "sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
)

const testApp = "my-app"
const TestWorkload = "my-workload"

var errCreate = errors.New("badcreate")
var errAppCreate = errors.New("bad")

func TestAppHandlerHandle(t *testing.T) {

	mockEventSender := eventsender.NewK8sSender(record.NewFakeRecorder(100))
	log := testr.New(t)

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example-pod",
			Namespace: namespace,
			Annotations: map[string]string{
				apicommon.WorkloadAnnotation: TestWorkload,
				apicommon.VersionAnnotation:  "0.1",
			},
		}}

	singleServiceCreationReq := &apilifecycle.KeptnAppCreationRequest{
		TypeMeta: metav1.TypeMeta{
			Kind:       "KeptnAppCreationRequest",
			APIVersion: "lifecycle.keptn.sh/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:            TestWorkload,
			Namespace:       namespace,
			ResourceVersion: "1",
			Annotations: map[string]string{
				"keptn.sh/app-type": "single-service",
			},
		},
		Spec: apilifecycle.KeptnAppCreationRequestSpec{AppName: TestWorkload},
	}

	tests := []struct {
		name    string
		client  client.Client
		pod     *corev1.Pod
		wanterr error
		wantReq *apilifecycle.KeptnAppCreationRequest
	}{
		{
			name:    "Create AppCreationRequest inherit from workload",
			pod:     pod,
			client:  testcommon.NewTestClient(),
			wantReq: singleServiceCreationReq,
		},
		{
			name:    "AppCreationRequest already exists",
			pod:     pod,
			client:  testcommon.NewTestClient(singleServiceCreationReq),
			wantReq: singleServiceCreationReq,
		},
		{
			name: "Create AppCreationRequest",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "example-pod",
					Namespace: namespace,
					Annotations: map[string]string{
						apicommon.AppAnnotation:      testApp,
						apicommon.WorkloadAnnotation: TestWorkload,
						apicommon.VersionAnnotation:  "0.1",
					},
				}},
			client: testcommon.NewTestClient(),
			wantReq: &apilifecycle.KeptnAppCreationRequest{
				TypeMeta: metav1.TypeMeta{
					Kind:       "KeptnAppCreationRequest",
					APIVersion: "lifecycle.keptn.sh/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:            testApp,
					Namespace:       namespace,
					ResourceVersion: "1",
				},
				Spec: apilifecycle.KeptnAppCreationRequestSpec{AppName: testApp},
			},
		},
		{
			name: "Error Fetching AppCreationRequest",
			pod:  &corev1.Pod{},
			client: k8sfake.NewClientBuilder().WithInterceptorFuncs(interceptor.Funcs{
				Get: func(ctx context.Context, client client.WithWatch, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
					return errAppCreate
				},
			}).Build(),
			wanterr: errAppCreate,
		},
		{
			name: "Error Creating AppCreationRequest",
			pod:  pod,
			client: k8sfake.NewClientBuilder().WithInterceptorFuncs(interceptor.Funcs{
				Create: func(ctx context.Context, client client.WithWatch, obj client.Object, opts ...client.CreateOption) error {
					return errCreate
				},
			}).Build(),
			wanterr: errCreate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appHandler := &AppCreationRequestHandler{
				Client:      tt.client,
				Log:         log,
				EventSender: mockEventSender,
			}
			err := appHandler.Handle(context.TODO(), tt.pod, namespace)

			if tt.wanterr != nil {
				require.NotNil(t, err)
				require.ErrorIs(t, err, tt.wanterr)
			} else {
				require.Nil(t, err)
			}

			if tt.wantReq != nil {
				creationReq := &apilifecycle.KeptnAppCreationRequest{}
				err = tt.client.Get(context.TODO(), types.NamespacedName{Name: tt.wantReq.Name, Namespace: tt.wantReq.Namespace}, creationReq)
				require.Nil(t, err)
				require.Equal(t, tt.wantReq, creationReq)
			}

		})
	}
}

func TestAppHandlerCreateAppSucceeds(t *testing.T) {
	fakeClient := testcommon.NewTestClient()
	logger := logr.Discard()
	eventSender := eventsender.NewK8sSender(record.NewFakeRecorder(100))

	appHandler := &AppCreationRequestHandler{
		Client:      fakeClient,
		Log:         logger,
		EventSender: eventSender,
	}

	ctx := context.TODO()
	name := "myappcreationreq"
	newAppCreationRequest := &apilifecycle.KeptnAppCreationRequest{
		ObjectMeta: metav1.ObjectMeta{Name: name},
	}
	err := appHandler.createResource(ctx, newAppCreationRequest)

	require.Nil(t, err)
	creationReq := &apilifecycle.KeptnAppCreationRequest{}
	err = fakeClient.Get(ctx, types.NamespacedName{Name: name}, creationReq)
	require.Nil(t, err)

}

func TestAppHandlerCreateAppFails(t *testing.T) {
	fakeClient := testcommon.NewTestClient()
	logger := logr.Discard()
	eventSender := eventsender.NewK8sSender(record.NewFakeRecorder(100))

	appHandler := &AppCreationRequestHandler{
		Client:      fakeClient,
		Log:         logger,
		EventSender: eventSender,
	}

	ctx := context.TODO()
	newAppCreationRequest := &apilifecycle.KeptnAppCreationRequest{
		ObjectMeta: metav1.ObjectMeta{},
	}
	err := appHandler.createResource(ctx, newAppCreationRequest)
	require.Error(t, err)

}

func TestGenerateAppCreationRequest(t *testing.T) {
	// Mock a context with OpenTelemetry tracer enabled
	ctx := context.Background()

	// Create a sample Pod
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: namespace,
			Annotations: map[string]string{
				apicommon.WorkloadAnnotation: workloadName,
			},
		},
	}

	// Test case 1: Pod does not have app annotation
	t.Run("PodWithoutAppAnnotation", func(t *testing.T) {
		kacr := generateResource(ctx, pod, namespace)

		require.Equal(t, namespace, kacr.Namespace)
		require.Equal(t, string(apicommon.AppTypeSingleService), kacr.Annotations[apicommon.AppTypeAnnotation])
		require.Equal(t, workloadName, kacr.Name)
		require.Equal(t, workloadName, kacr.Spec.AppName)
	})

	// Test case 2: Pod has app annotation
	t.Run("PodWithAppAnnotation", func(t *testing.T) {
		// Add app annotation to the Pod
		pod.ObjectMeta.Annotations[apicommon.AppAnnotation] = lowerAppName
		kacr := generateResource(ctx, pod, namespace)

		require.Equal(t, namespace, kacr.Namespace)
		require.Empty(t, kacr.Annotations[apicommon.AppTypeAnnotation]) // No app type annotation
		require.Equal(t, lowerAppName, kacr.Name)
		require.Equal(t, lowerAppName, kacr.Spec.AppName)
	})
}
