package manager

import (
	"context"

	"github.com/stretchr/testify/mock"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

type MockManager struct {
	TestManager
	mock.Mock
}

func (mgr *MockManager) Start(ctx context.Context) error {
	args := mgr.Called(ctx)
	return args.Error(0)
}

func (mgr *MockManager) AddHealthzCheck(name string, check healthz.Checker) error {
	args := mgr.Called(name, check)
	return args.Error(0)
}

func (mgr *MockManager) AddReadyzCheck(name string, check healthz.Checker) error {
	args := mgr.Called(name, check)
	return args.Error(0)
}

func (mgr *MockManager) GetWebhookServer() *webhook.Server {
	args := mgr.Called()
	return args.Get(0).(*webhook.Server)
}

type MockProvider struct {
	mock.Mock
}

func (provider *MockProvider) SetupWebhookServer(mgr manager.Manager) {
	provider.Called(mgr)
}
