package config

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/rest"
)

func TestMockConfigProvider(t *testing.T) {
	mockProvider := &MockProvider{}
	var _ Provider = mockProvider

	expectedCfg := &rest.Config{}
	expectedErr := errors.New("config provider error")
	mockProvider.On("GetConfig").Return(expectedCfg, expectedErr)

	cfg, err := mockProvider.GetConfig()

	mockProvider.AssertCalled(t, "GetConfig")
	assert.NotNil(t, cfg)
	assert.NotNil(t, err)
	assert.Equal(t, expectedCfg, cfg)
	assert.Equal(t, expectedErr, err)
}
