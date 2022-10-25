package keptnapp

import (
	"context"
	lifecyclev1alpha1 "github.com/keptn/lifecycle-controller/operator/api/v1alpha1"
	keptncommon "github.com/keptn/lifecycle-controller/operator/api/v1alpha1/common"
	"github.com/keptn/lifecycle-controller/operator/controllers/common/fake"
	"github.com/magiconair/properties/assert"
	"go.opentelemetry.io/otel/trace"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"testing"
)

//EXample Unit test on help function
func TestKeptnAppReconciler_createAppVersionSuccess(t *testing.T) {

	app := &lifecyclev1alpha1.KeptnApp{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-app",
			Namespace: "default",
		},
		Spec: lifecyclev1alpha1.KeptnAppSpec{
			Version: "1.0.0",
		},
		Status: lifecyclev1alpha1.KeptnAppStatus{},
	}
	//setup logger
	opts := zap.Options{
		Development: true,
	}
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	//fake a tracer
	tr := &fake.ITracerMock{StartFunc: func(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
		return ctx, trace.SpanFromContext(ctx)
	}}

	//add keptn lfc scheme to k8sdefault
	setupScheme(t)

	r := &KeptnAppReconciler{
		Log:    ctrl.Log.WithName("test-appController"),
		Tracer: tr,
		Scheme: scheme.Scheme,
	}

	appVersion, err := r.createAppVersion(context.TODO(), app)
	if err != nil {
		t.Errorf("Error Creating appVersion: %s", err.Error())
	}
	t.Log("Verifying created app")
	assert.Equal(t, appVersion.Namespace, app.Namespace)
	assert.Equal(t, appVersion.Name, app.Name+"-"+app.Spec.Version)

}

func setupScheme(t *testing.T) {
	err := lifecyclev1alpha1.AddToScheme(scheme.Scheme)
	if err != nil {
		t.Fatalf("Could not set scheme %s", err.Error())
	}
}

func TestKeptnAppReconciler_Reconcile(t *testing.T) {

	r, eventChannel := setupReconciler(t)

	tests := []struct {
		name    string
		req     ctrl.Request
		want    ctrl.Result
		wantErr error
		event   string
	}{
		{
			name: "test simple create appVersion",
			req: ctrl.Request{
				NamespacedName: types.NamespacedName{
					Namespace: "default",
					Name:      "myapp",
				},
			},
			want:    ctrl.Result{},
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
			want:    ctrl.Result{},
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
			want:    ctrl.Result{},
			wantErr: nil,
		},
	}

	//setting up fakeclient CRD data

	addApp(r, "myapp")
	addApp(r, "myfinishedapp")
	addAppVersion(r, "myfinishedapp-1.0.0", keptncommon.StateSucceeded)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := r.Reconcile(context.TODO(), tt.req)
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("Reconcile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reconcile() got = %v, want %v", got, tt.want)
			}
			if tt.event != "" {
				event := <-eventChannel
				assert.Matches(t, event, tt.event)
			}
		})
	}
}

func setupReconciler(t *testing.T) (*KeptnAppReconciler, chan string) {
	//setup logger
	opts := zap.Options{
		Development: true,
	}
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	//fake a tracer
	tr := &fake.ITracerMock{StartFunc: func(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
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
	return r, recorder.Events
}

func addApp(r *KeptnAppReconciler, name string) error {
	app := &lifecyclev1alpha1.KeptnApp{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
		},
		Spec: lifecyclev1alpha1.KeptnAppSpec{
			Version: "1.0.0",
		},
		Status: lifecyclev1alpha1.KeptnAppStatus{},
	}
	return r.Client.Create(context.TODO(), app)

}

func addAppVersion(r *KeptnAppReconciler, name string, status keptncommon.KeptnState) error {
	app := &lifecyclev1alpha1.KeptnAppVersion{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
		},
		Spec: lifecyclev1alpha1.KeptnAppVersionSpec{},
		Status: lifecyclev1alpha1.KeptnAppVersionStatus{
			Status: status,
		},
	}
	return r.Client.Create(context.TODO(), app)

}
