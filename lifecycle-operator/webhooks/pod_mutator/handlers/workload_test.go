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
const myworkload = "my-workload"

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
			Name:      "my-workload-my-workload",
			Namespace: namespace,
			OwnerReferences: []metav1.OwnerReference{
				{},
			},
			ResourceVersion: "1",
		},
		Spec: klcv1alpha3.KeptnWorkloadSpec{
			AppName: myworkload,
			Version: "0.1",
		},
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example-pod",
			Namespace: namespace,
			Annotations: map[string]string{
				apicommon.WorkloadAnnotation: myworkload,
				apicommon.VersionAnnotation:  "0.1",
			},
		}}
	// Define test cases
	tests := []struct {
		name         string
		client       client.Client
		pod          *corev1.Pod
		wanterr      string
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
					return errors.New("bad")
				},
			}).Build(),
			wanterr: "could not fetch Workload: bad",
		},
		{
			name: "Error Creating Workload",
			pod:  pod,
			client: k8sfake.NewClientBuilder().WithInterceptorFuncs(interceptor.Funcs{
				Create: func(ctx context.Context, client client.WithWatch, obj client.Object, opts ...client.CreateOption) error {
					return errors.New("badcreate")
				},
			}).Build(),
			wanterr: "badcreate",
		},
		{
			name: "Error Updating Workload",
			pod:  pod,
			client: k8sfake.NewClientBuilder().WithInterceptorFuncs(interceptor.Funcs{
				Update: func(ctx context.Context, client client.WithWatch, obj client.Object, opts ...client.UpdateOption) error {
					return errors.New("badupdate")
				},
			}).WithObjects(workload).Build(),
			wanterr: "badupdate",
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

			if tt.wanterr != "" {
				require.NotNil(t, err)
				require.Contains(t, err.Error(), tt.wanterr)
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
			Name:      "my-workload-my-workload",
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
