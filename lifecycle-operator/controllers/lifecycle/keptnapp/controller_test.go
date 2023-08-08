package keptnapp

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	lfcv1alpha3 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	controllercommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/fake"
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

	app := &lfcv1alpha3.KeptnApp{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:       "my-app",
			Namespace:  "default",
			Generation: 1,
		},
		Spec: lfcv1alpha3.KeptnAppSpec{
			Version: "1.0.0",
		},
		Status: lfcv1alpha3.KeptnAppStatus{},
	}
	r, _, _ := setupReconciler()

	appVersion, err := r.createAppVersion(context.TODO(), app)
	if err != nil {
		t.Errorf("Error Creating appVersion: %s", err.Error())
	}
	t.Log("Verifying created app")
	assert.Equal(t, appVersion.Namespace, app.Namespace)
	assert.Equal(t, appVersion.Name, fmt.Sprintf("%s-%s-%s", app.Name, app.Spec.Version, apicommon.Hash(app.Generation)))
}

func TestKeptnAppReconciler_createAppVersionWithLongName(t *testing.T) {
	//nolint:gci
	longName := `loremipsumissimplydummytextoftheprintingandtypesettingindustryloremipsumissimplydummytextoftheprintingandtypesettingindustryloremipsumissimplydummytextoftheprintingandtypesettingindustryloremipsumissimplydummytextoftheprintingandtypesettingindustryloremloremax`
	//nolint:gci
	trimmedName := `loremipsumissimplydummytextoftheprintingandtypesettingindustryloremipsumissimplydummytextoftheprintingandtypesettingindustryloremipsumissimplydummytextoftheprintingandtypesettingindustryloremipsumissimplydummytextoftheprintingandtypeset-version-5feceb66`

	app := &lfcv1alpha3.KeptnApp{
		ObjectMeta: metav1.ObjectMeta{
			Name: longName,
		},
		Spec: lfcv1alpha3.KeptnAppSpec{
			Version: "version",
		},
	}
	r, _, _ := setupReconciler()

	appVersion, err := r.createAppVersion(context.Background(), app)
	if err != nil {
		t.Errorf("Error creating app version: %s", err.Error())
	}
	t.Log("Verifying app name length is not greater than MaxK8sObjectLen")
	assert.Equal(t, appVersion.ObjectMeta.Name, trimmedName)
}

func TestKeptnAppReconciler_reconcile(t *testing.T) {

	r, _, tracer := setupReconciler()

	tests := []struct {
		name           string
		req            ctrl.Request
		wantErr        error
		appVersionName string
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

	// setting up fakeclient CRD data

	err := controllercommon.AddApp(r.Client, "myapp")
	require.Nil(t, err)
	err = controllercommon.AddApp(r.Client, "myfinishedapp")
	require.Nil(t, err)
	err = controllercommon.AddAppVersion(r.Client, "default", "myfinishedapp", "1.0.0-6b86b273", nil, lfcv1alpha3.KeptnAppVersionStatus{Status: apicommon.StateSucceeded})
	require.Nil(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			_, err := r.Reconcile(context.TODO(), tt.req)
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("Reconcile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.appVersionName != "" {
				keptnappversion := &lfcv1alpha3.KeptnAppVersion{}
				err = r.Client.Get(context.TODO(), types.NamespacedName{Namespace: "default", Name: "myapp-1.0.0-6b86b273"}, keptnappversion)
				require.Nil(t, err)
			}

		})
	}

	// check correct traces
	assert.Equal(t, len(tracer.StartCalls()), 4)
	// case 1 reconcile and create app ver
	assert.Equal(t, tracer.StartCalls()[0].SpanName, "reconcile_app")
	assert.Equal(t, tracer.StartCalls()[1].SpanName, "create_app_version")
	assert.Equal(t, tracer.StartCalls()[2].SpanName, "myapp-1.0.0-6b86b273")
	// case 2 creates no span because notfound
	// case 3 reconcile finished crd
	assert.Equal(t, tracer.StartCalls()[3].SpanName, "reconcile_app")
}

func TestKeptnAppReconciler_deprecateAppVersions(t *testing.T) {
	r, _, _ := setupReconciler()

	err := controllercommon.AddApp(r.Client, "myapp")
	require.Nil(t, err)

	_, err = r.Reconcile(context.TODO(), ctrl.Request{
		NamespacedName: types.NamespacedName{
			Namespace: "default",
			Name:      "myapp",
		},
	})

	require.Nil(t, err)

	keptnappversion := &lfcv1alpha3.KeptnAppVersion{}
	err = r.Client.Get(context.TODO(), types.NamespacedName{Namespace: "default", Name: "myapp-1.0.0-6b86b273"}, keptnappversion)
	require.Nil(t, err)

	err = controllercommon.UpdateAppRevision(r.Client, "myapp", 2)
	require.Nil(t, err)

	_, err = r.Reconcile(context.TODO(), ctrl.Request{
		NamespacedName: types.NamespacedName{
			Namespace: "default",
			Name:      "myapp",
		},
	})

	require.Nil(t, err)

	err = r.Client.Get(context.TODO(), types.NamespacedName{Namespace: "default", Name: "myapp-1.0.0-d4735e3a"}, keptnappversion)
	require.Nil(t, err)

	err = r.Client.Get(context.TODO(), types.NamespacedName{Namespace: "default", Name: "myapp-1.0.0-6b86b273"}, keptnappversion)
	require.Nil(t, err)
	require.Equal(t, apicommon.StateDeprecated, keptnappversion.Status.Status)
}

func setupReconciler() (*KeptnAppReconciler, chan string, *fake.ITracerMock) {
	// setup logger
	opts := zap.Options{
		Development: true,
	}
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	// fake a tracer
	tr := &fake.ITracerMock{StartFunc: func(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
		return ctx, trace.SpanFromContext(ctx)
	}}

	tf := &fake.TracerFactoryMock{GetTracerFunc: func(name string) trace.Tracer {
		return tr
	}}

	fakeClient := fake.NewClient()

	recorder := record.NewFakeRecorder(100)
	r := &KeptnAppReconciler{
		Client:        fakeClient,
		Scheme:        scheme.Scheme,
		EventSender:   controllercommon.NewEventSender(recorder),
		Log:           ctrl.Log.WithName("test-appController"),
		TracerFactory: tf,
	}
	return r, recorder.Events, tr
}
