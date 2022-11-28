package manager

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type Provider interface {
	CreateManager(namespace string, scheme *runtime.Scheme, config *rest.Config) (manager.Manager, error)
}
