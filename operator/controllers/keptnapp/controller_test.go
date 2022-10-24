package keptnapp

import (
	"context"
	lifecyclev1alpha1 "github.com/keptn/lifecycle-controller/operator/api/v1alpha1"
	"github.com/keptn/lifecycle-controller/operator/controllers/common/fake"
	"github.com/magiconair/properties/assert"
	"go.opentelemetry.io/otel/trace"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
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
	tr := fake.ITracerMock{
		StartFunc: func(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
			return ctx, trace.SpanFromContext(ctx)
		},
	}

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
