package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const configContent = `apiVersion: v1
clusters:
- cluster:
    server: https://example.com
  name: test-cluster
contexts:
- context:
    cluster: test-cluster
    user: test-user
  name: test-context
current-context: test-context
kind: Config
preferences: {}
users:
- name: test-user
  user:
    password: test-password
    username: test-username
`

func TestConfigProvider_NewKubeConfigProvider(t *testing.T) {
	provider := NewKubeConfigProvider()
	require.NotNil(t, provider)
}

func TestConfigProvider_GetConfigSuccess(t *testing.T) {
	tempDir := setupMockedKubeConfig(t)
	defer os.RemoveAll(tempDir)

	provider := NewKubeConfigProvider()

	config, err := provider.GetConfig()
	require.Nil(t, err)
	require.NotNil(t, config)
}

func setupMockedKubeConfig(t *testing.T) string {
	tempDir := t.TempDir()

	kubeConfigFile := filepath.Join(tempDir, ".kube", "config")
	err := os.MkdirAll(filepath.Dir(kubeConfigFile), 0755)
	require.Nil(t, err)

	err = os.WriteFile(kubeConfigFile, []byte(configContent), 0644)
	require.Nil(t, err)

	t.Setenv("HOME", tempDir)

	t.Setenv("KUBECONFIG", kubeConfigFile)

	return tempDir
}
