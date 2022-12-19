package manager

import (
	"context"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/common/fake"
	"k8s.io/apimachinery/pkg/runtime"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/config/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type TestManager struct {
	manager.Manager
}

func (mgr *TestManager) GetClient() client.Client {
	return fake.NewClient()

}

func (mgr *TestManager) GetAPIReader() client.Reader {
	return struct{ client.Reader }{}
}

func (mgr *TestManager) GetControllerOptions() v1alpha1.ControllerConfigurationSpec {
	return v1alpha1.ControllerConfigurationSpec{}
}

func (mgr *TestManager) GetScheme() *runtime.Scheme {
	return fake.NewClient().Scheme()
}

func (mgr *TestManager) GetLogger() logr.Logger {
	return testr.New(&testing.T{})
}

func (mgr *TestManager) SetFields(interface{}) error {
	return nil
}

func (mgr *TestManager) Add(manager.Runnable) error {
	return nil
}

func (mgr *TestManager) Start(_ context.Context) error {
	return nil
}
