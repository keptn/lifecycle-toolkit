package handlers

import (
	"context"
	"testing"

	"github.com/go-logr/logr/testr"
	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	controllercommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/fake"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	k8sfake "sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
)

const namespace = "test-namespace"
const testAppWorkload = "my-workload-my-workload"

var errUpdate = errors.New("badupdate")
var errFetch = errors.New("bad")

func TestHandle(t *testing.T) {

	mockEventSender := controllercommon.NewK8sSender(record.NewFakeRecorder(100))
	log := testr.New(t)
	tr := &fake.ITracerMock{StartFunc: func(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
		return ctx, trace.SpanFromContext(ctx)
	}}
	workload := &klcv1alpha3.KeptnWorkload{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-workload-my-workload",
			Namespace: namespace,
		},
	}

	wantWorkload := &klcv1alpha3.KeptnWorkload{
		TypeMeta: metav1.TypeMeta{Kind: "KeptnWorkload", APIVersion: "lifecycle.keptn.sh/v1alpha3"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      testAppWorkload,
			Namespace: namespace,
			OwnerReferences: []metav1.OwnerReference{
				{},
			},
			ResourceVersion: "1",
		},
		Spec: klcv1alpha3.KeptnWorkloadSpec{
			AppName: TestWorkload,
			Version: "0.1",
		},
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example-pod",
			Namespace: namespace,
			Annotations: map[string]string{
				apicommon.WorkloadAnnotation: TestWorkload,
				apicommon.VersionAnnotation:  "0.1",
			},
		}}
	// Define test cases
	tests := []struct {
		name         string
		client       client.Client
		pod          *corev1.Pod
		wanterr      error
		wantWorkload *klcv1alpha3.KeptnWorkload
	}{
		{
			name:         "Create Workload",
			pod:          pod,
			client:       fake.NewClient(),
			wantWorkload: wantWorkload,
		},
		{
			name:         "Update Workload",
			pod:          pod,
			client:       fake.NewClient(wantWorkload),
			wantWorkload: wantWorkload,
		},
		{
			name: "Error Fetching Workload",
			pod:  &corev1.Pod{},
			client: k8sfake.NewClientBuilder().WithInterceptorFuncs(interceptor.Funcs{
				Get: func(ctx context.Context, client client.WithWatch, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
					return errFetch
				},
			}).Build(),
			wanterr: errFetch,
		},
		{
			name: "Error Creating Workload",
			pod:  pod,
			client: k8sfake.NewClientBuilder().WithInterceptorFuncs(interceptor.Funcs{
				Create: func(ctx context.Context, client client.WithWatch, obj client.Object, opts ...client.CreateOption) error {
					return errCreate
				},
			}).Build(),
			wanterr: errCreate,
		},
		{
			name: "Error Updating Workload",
			pod:  pod,
			client: k8sfake.NewClientBuilder().WithInterceptorFuncs(interceptor.Funcs{
				Update: func(ctx context.Context, client client.WithWatch, obj client.Object, opts ...client.UpdateOption) error {
					return errUpdate
				},
			}).WithObjects(workload).Build(),
			wanterr: errUpdate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workloadHandler := &WorkloadHandler{
				Client:      tt.client,
				Log:         log,
				EventSender: mockEventSender,
				Tracer:      tr,
			}
			err := workloadHandler.Handle(context.TODO(), tt.pod, "test-namespace")

			if tt.wanterr != nil {
				require.NotNil(t, err)
				require.ErrorIs(t, err, tt.wanterr)
			} else {
				require.Nil(t, err)
			}

			if tt.wantWorkload != nil {
				actualWorkload := &klcv1alpha3.KeptnWorkload{}
				err = tt.client.Get(context.TODO(), types.NamespacedName{Name: tt.wantWorkload.Name, Namespace: tt.wantWorkload.Namespace}, actualWorkload)
				require.Nil(t, err)
				require.Equal(t, tt.wantWorkload, actualWorkload)
			}

		})
	}
}

func TestUpdateWorkloadNoSpecChanges(t *testing.T) {
	mockEventSender := controllercommon.NewK8sSender(record.NewFakeRecorder(100))
	log := testr.New(t)
	tr := &fake.ITracerMock{StartFunc: func(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
		return ctx, trace.SpanFromContext(ctx)
	}}

	workload := &klcv1alpha3.KeptnWorkload{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testAppWorkload,
			Namespace: namespace,
		},
	}
	a := &WorkloadHandler{
		Client:      nil,
		Log:         log,
		Tracer:      tr,
		EventSender: mockEventSender,
	}
	err := a.updateWorkload(context.TODO(), workload, workload, nil)
	require.Nil(t, err)

}

func TestGenerateWorkload(t *testing.T) {
	testCases := []struct {
		name           string
		podAnnotations map[string]string
		expected       *klcv1alpha3.KeptnWorkload
	}{
		{
			name: "Pod with annotations",
			podAnnotations: map[string]string{
				apicommon.VersionAnnotation:                  "v1",
				apicommon.PreDeploymentTaskAnnotation:        "task1,task2",
				apicommon.PostDeploymentTaskAnnotation:       "task3,task4",
				apicommon.PreDeploymentEvaluationAnnotation:  "eval1,eval2",
				apicommon.PostDeploymentEvaluationAnnotation: "eval3,eval4",
				apicommon.K8sRecommendedAppAnnotations:       "my-app",
			},
			expected: &klcv1alpha3.KeptnWorkload{
				ObjectMeta: metav1.ObjectMeta{
					Name:        getWorkloadName(&metav1.ObjectMeta{}, "my-app"),
					Namespace:   "my-namespace",
					Annotations: map[string]string{},
					OwnerReferences: []metav1.OwnerReference{
						{
							UID:        "owner-uid",
							Kind:       "Deployment",
							Name:       "deployment-1",
							APIVersion: "apps/v1",
						},
					},
				},
				Spec: klcv1alpha3.KeptnWorkloadSpec{
					AppName:                   "my-app",
					Version:                   "v1",
					ResourceReference:         klcv1alpha3.ResourceReference{UID: "owner-uid", Kind: "Deployment", Name: "deployment-1"},
					PreDeploymentTasks:        []string{"task1", "task2"},
					PostDeploymentTasks:       []string{"task3", "task4"},
					PreDeploymentEvaluations:  []string{"eval1", "eval2"},
					PostDeploymentEvaluations: []string{"eval3", "eval4"},
				},
			},
		},
		{
			name:           "Pod with no annotations",
			podAnnotations: nil,
			expected: &klcv1alpha3.KeptnWorkload{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "-",
					Namespace:   "my-namespace",
					Annotations: map[string]string{},
					OwnerReferences: []metav1.OwnerReference{{
						APIVersion: "apps/v1",
						Kind:       "Deployment",
						Name:       "deployment-1",
						UID:        "owner-uid",
					},
					}},
				Spec: klcv1alpha3.KeptnWorkloadSpec{
					ResourceReference: klcv1alpha3.ResourceReference{
						UID:  "owner-uid",
						Kind: "Deployment",
						Name: "deployment-1"}}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a sample context and pod object
			ctx := context.TODO()
			pod := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: tc.podAnnotations,
					OwnerReferences: []metav1.OwnerReference{
						{
							UID:        "owner-uid",
							Kind:       "Deployment",
							Name:       "deployment-1",
							APIVersion: "apps/v1",
						},
					},
				},
			}

			result := generateWorkload(ctx, pod, "my-namespace")
			require.Equal(t, tc.expected, result)
		})
	}
}
