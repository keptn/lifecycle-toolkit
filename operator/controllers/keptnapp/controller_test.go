package keptnapp

import (
	"context"
	"reflect"
	"testing"

	lfcv1alpha2 "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha2"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha2/common"
	controllercommon "github.com/keptn/lifecycle-toolkit/operator/controllers/common"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/common/fake"
	interfacesfake "github.com/keptn/lifecycle-toolkit/operator/controllers/interfaces/fake"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

// Example Unit test on help function
func TestKeptnAppReconciler_createAppVersionSuccess(t *testing.T) {

	app := &lfcv1alpha2.KeptnApp{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-app",
			Namespace: "default",
		},
		Spec: lfcv1alpha2.KeptnAppSpec{
			Version: "1.0.0",
		},
		Status: lfcv1alpha2.KeptnAppStatus{},
	}
	r, _, _ := setupReconciler(t)

	appVersion, err := r.createAppVersion(context.TODO(), app)
	if err != nil {
		t.Errorf("Error Creating appVersion: %s", err.Error())
	}
	t.Log("Verifying created app")
	assert.Equal(t, appVersion.Namespace, app.Namespace)
	assert.Equal(t, appVersion.Name, app.Name+"-"+app.Spec.Version)

}

func TestKeptnAppReconciler_reconcile(t *testing.T) {

	r, eventChannel, tracer := setupReconciler(t)

	tests := []struct {
		name    string
		req     ctrl.Request
		wantErr error
		event   string //check correct events are generated
	}{
		{
			name: "test simple create appVersion",
			req: ctrl.Request{
				NamespacedName: types.NamespacedName{
					Namespace: "default",
					Name:      "myapp",
				},
			},
			wantErr: nil,
			event:   `Normal AppVersionCreated Created KeptnAppVersion / Namespace: default, Name: myapp-1.0.0`,
		},
		{
			name: "test simple notfound should not return error nor event",
			req: ctrl.Request{
				NamespacedName: types.NamespacedName{
					Namespace: "default",
					Name:      "mynotthereapp",
				},
			},
			wantErr: nil,
		},
		{
			name: "test existing appVersion nothing done",
			req: ctrl.Request{
				NamespacedName: types.NamespacedName{
					Namespace: "default",
					Name:      "myfinishedapp",
				},
			},
			wantErr: nil,
		},
	}

	//setting up fakeclient CRD data

	err := controllercommon.AddApp(r.Client, "myapp")
	require.Nil(t, err)
	err = controllercommon.AddApp(r.Client, "myfinishedapp")
	require.Nil(t, err)
	err = controllercommon.AddAppVersion(r.Client, "default", "myfinishedapp", "1.0.0", nil, lfcv1alpha2.KeptnAppVersionStatus{Status: apicommon.StateSucceeded})
	require.Nil(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			_, err := r.Reconcile(context.TODO(), tt.req)
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("Reconcile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.event != "" {
				event := <-eventChannel
				assert.Matches(t, event, tt.event)
			}

		})

	}

	// check correct traces
	assert.Equal(t, len(tracer.StartCalls()), 4)
	// case 1 reconcile and create app ver
	assert.Equal(t, tracer.StartCalls()[0].SpanName, "reconcile_app")
	assert.Equal(t, tracer.StartCalls()[1].SpanName, "create_app_version")
	assert.Equal(t, tracer.StartCalls()[2].SpanName, "myapp-1.0.0")
	//case 2 creates no span because notfound
	//case 3 reconcile finished crd
	assert.Equal(t, tracer.StartCalls()[3].SpanName, "reconcile_app")
}

func setupReconciler(t *testing.T) (*KeptnAppReconciler, chan string, *interfacesfake.ITracerMock) {
	//setup logger
	opts := zap.Options{
		Development: true,
	}
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	//fake a tracer
	tr := &interfacesfake.ITracerMock{StartFunc: func(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
		return ctx, trace.SpanFromContext(ctx)
	}}

	fakeClient, err := fake.NewClient()
	if err != nil {
		t.Errorf("Reconcile() error when setting up fake client %v", err)
	}
	recorder := record.NewFakeRecorder(100)
	r := &KeptnAppReconciler{
		Client:   fakeClient,
		Scheme:   scheme.Scheme,
		Recorder: recorder,
		Log:      ctrl.Log.WithName("test-appController"),
		Tracer:   tr,
	}
	return r, recorder.Events, tr
}
