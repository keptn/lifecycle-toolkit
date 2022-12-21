package webhook

import (
	"testing"

	"github.com/keptn/lifecycle-toolkit/operator/cmd/config"
	cmdManager "github.com/keptn/lifecycle-toolkit/operator/cmd/manager"
	"github.com/stretchr/testify/assert"
)

func TestWebhookCommandBuilder(t *testing.T) {

	t.Run("set config provider", func(t *testing.T) {
		builder := NewWebhookBuilder()

		assert.NotNil(t, builder)

		expectedProvider := &config.MockProvider{}
		builder = builder.SetConfigProvider(expectedProvider)

		assert.Equal(t, expectedProvider, builder.configProvider)
	})
	t.Run("set manager provider", func(t *testing.T) {
		expectedProvider := &cmdManager.MockProvider{}
		builder := NewWebhookBuilder().SetManagerProvider(expectedProvider)

		assert.Equal(t, expectedProvider, builder.managerProvider)
	})
	t.Run("set namespace", func(t *testing.T) {
		builder := NewWebhookBuilder().SetNamespace("namespace")

		assert.Equal(t, "namespace", builder.namespace)
	})
}
