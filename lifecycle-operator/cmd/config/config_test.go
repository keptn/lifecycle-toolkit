package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigProvider(t *testing.T) {
	t.Run("NewKubeConfigProvider", func(t *testing.T) {
		provider := NewKubeConfigProvider()
		assert.NotNil(t, provider)
	})

	// t.Run("GetConfigSuccess", func(t *testing.T) {
	// 	provider := NewKubeConfigProvider()

	// 	config, err := provider.GetConfig()
	// 	assert.NoError(t, err)
	// 	assert.NotNil(t, config)
	// })
}
