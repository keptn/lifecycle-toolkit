package fake

import (
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha2"
	corev1 "k8s.io/api/core/v1"
	apiv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

// NewClient returns a new controller-runtime fake Client configured with the Operator's scheme, and initialized with objs.
func NewClient(objs ...client.Object) client.Client {
	setupSchemes()
	return fake.NewClientBuilder().WithScheme(scheme.Scheme).WithObjects(objs...).Build()
}

func setupSchemes() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme.Scheme))
	utilruntime.Must(corev1.AddToScheme(scheme.Scheme))
	utilruntime.Must(apiv1.AddToScheme(scheme.Scheme))
	utilruntime.Must(metricsapi.AddToScheme(scheme.Scheme))
}
