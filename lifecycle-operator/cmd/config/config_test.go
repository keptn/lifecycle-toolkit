package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigProvider(t *testing.T) {
	t.Run("NewKubeConfigProvider", func(t *testing.T) {
		provider := NewKubeConfigProvider()
		assert.NotNil(t, provider)
	})

	t.Run("GetConfigSuccess", func(t *testing.T) {

		tempDir := setupMockedKubeConfig(t)
		defer os.RemoveAll(tempDir)

		provider := NewKubeConfigProvider()

		config, err := provider.GetConfig()
		assert.NoError(t, err)
		assert.NotNil(t, config)
	})
}

func setupMockedKubeConfig(t *testing.T) string {
	tempDir := t.TempDir()

	kubeConfigFile := filepath.Join(tempDir, ".kube", "config")
	err := os.MkdirAll(filepath.Dir(kubeConfigFile), 0755)
	assert.NoError(t, err)
	err = os.WriteFile(kubeConfigFile, []byte("mocked kubeconfig content"), 0644)
	assert.NoError(t, err)

	os.Setenv("HOME", tempDir)

	os.Setenv("KUBECONFIG", kubeConfigFile)

	return tempDir
}
