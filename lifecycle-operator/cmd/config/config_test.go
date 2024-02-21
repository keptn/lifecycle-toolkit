package config

import (
	"os"
	"path/filepath"
	"strings"
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

	configContent := genKubeconfig("test-context")
	err = os.WriteFile(kubeConfigFile, []byte(configContent), 0644)
	assert.NoError(t, err)

	t.Setenv("HOME", tempDir)

	t.Setenv("KUBECONFIG", kubeConfigFile)

	return tempDir
}

func genKubeconfig(contexts ...string) string {
	var sb strings.Builder
	sb.WriteString(`---
apiVersion: v1
kind: Config
clusters:
`)
	for _, ctx := range contexts {
		sb.WriteString(`- cluster:
    server: ` + ctx + `
  name: ` + ctx + `
`)
	}
	sb.WriteString("contexts:\n")
	for _, ctx := range contexts {
		sb.WriteString(`- context:
    cluster: ` + ctx + `
    user: ` + ctx + `
  name: ` + ctx + `
`)
	}

	sb.WriteString("users:\n")
	for _, ctx := range contexts {
		sb.WriteString(`- name: ` + ctx + `
`)
	}
	sb.WriteString("preferences: {}\n")
	if len(contexts) > 0 {
		sb.WriteString("current-context: " + contexts[0] + "\n")
	}

	return sb.String()
}
