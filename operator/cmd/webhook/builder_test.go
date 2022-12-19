package webhook

//func TestWebhookCommandBuilder(t *testing.T) {
//	t.Run("build command", func(t *testing.T) {
//		builder := NewWebhookCommandBuilder()
//		csiCommand := builder.Build()
//
//		assert.NotNil(t, csiCommand)
//		assert.Equal(t, use, csiCommand.Use)
//		assert.NotNil(t, csiCommand.RunE)
//	})
//	t.Run("set config provider", func(t *testing.T) {
//		builder := NewWebhookCommandBuilder()
//
//		assert.NotNil(t, builder)
//
//		expectedProvider := &config.MockProvider{}
//		builder = builder.SetConfigProvider(expectedProvider)
//
//		assert.Equal(t, expectedProvider, builder.configProvider)
//	})
//	t.Run("set manager provider", func(t *testing.T) {
//		expectedProvider := &cmdManager.MockProvider{}
//		builder := NewWebhookCommandBuilder().SetManagerProvider(expectedProvider)
//
//		assert.Equal(t, expectedProvider, builder.managerProvider)
//	})
//	t.Run("set namespace", func(t *testing.T) {
//		builder := NewWebhookCommandBuilder().SetNamespace("namespace")
//
//		assert.Equal(t, "namespace", builder.namespace)
//	})
//}
