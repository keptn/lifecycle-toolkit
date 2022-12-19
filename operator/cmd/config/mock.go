package config

import (
	"github.com/stretchr/testify/mock"
	"k8s.io/client-go/rest"
)

type MockProvider struct {
	mock.Mock
}

func (provider *MockProvider) GetConfig() (*rest.Config, error) {
	args := provider.Called()
	return args.Get(0).(*rest.Config), args.Error(1)
}
