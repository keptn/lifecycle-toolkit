package fake

import (
	lfcv1alpha1 "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

// NewClient returns a new controller-runtime fake Client configured with the Operator's scheme, and initialized with objs.
func NewClient(objs ...client.Object) (client.Client, error) {
	err := setupScheme()
	return fake.NewClientBuilder().WithScheme(scheme.Scheme).WithObjects(objs...).Build(), err
}

func setupScheme() error {
	return lfcv1alpha1.AddToScheme(scheme.Scheme)
}
