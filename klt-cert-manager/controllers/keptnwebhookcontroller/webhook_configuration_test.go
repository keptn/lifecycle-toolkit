package keptnwebhookcontroller

import (
	"testing"

	"github.com/stretchr/testify/assert"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
)

func createTestMutatingWebhookConfigList(_ *testing.T) *admissionregistrationv1.MutatingWebhookConfigurationList {
	return &admissionregistrationv1.MutatingWebhookConfigurationList{
		Items: []admissionregistrationv1.MutatingWebhookConfiguration{
			{
				Webhooks: []admissionregistrationv1.MutatingWebhook{
					{},
					{ClientConfig: admissionregistrationv1.WebhookClientConfig{}},
					{
						ClientConfig: admissionregistrationv1.WebhookClientConfig{
							CABundle: []byte{0, 1, 2, 3, 4},
						},
					},
				},
			},
		},
	}
}

func createTestValidatingWebhookConfigList(_ *testing.T) *admissionregistrationv1.ValidatingWebhookConfigurationList {
	return &admissionregistrationv1.ValidatingWebhookConfigurationList{
		Items: []admissionregistrationv1.ValidatingWebhookConfiguration{
			{
				Webhooks: []admissionregistrationv1.ValidatingWebhook{
					{},
					{ClientConfig: admissionregistrationv1.WebhookClientConfig{}},
					{
						ClientConfig: admissionregistrationv1.WebhookClientConfig{
							CABundle: []byte{0, 1, 2, 3, 4},
						},
					},
				},
			},
		},
	}
}

func TestGetClientConfigsFromMutatingWebhook(t *testing.T) {
	t.Run(`returns nil when config is nil`, func(t *testing.T) {
		clientConfigs := getClientConfigsFromMutatingWebhook(nil)
		assert.Nil(t, clientConfigs)
	})
	t.Run(`returns client configs of all configured webhooks`, func(t *testing.T) {
		const expectedClientConfigs = 3
		clientConfigs := getClientConfigsFromMutatingWebhook(createTestMutatingWebhookConfigList(t))

		assert.NotNil(t, clientConfigs)
		assert.Equal(t, expectedClientConfigs, len(clientConfigs))
	})
}

func TestGetClientConfigsFromValidatingWebhook(t *testing.T) {
	t.Run(`returns nil when config is nil`, func(t *testing.T) {
		clientConfigs := getClientConfigsFromValidatingWebhook(nil)
		assert.Nil(t, clientConfigs)
	})
	t.Run(`returns client configs of all configured webhooks`, func(t *testing.T) {
		const expectedClientConfigs = 3
		clientConfigs := getClientConfigsFromValidatingWebhook(createTestValidatingWebhookConfigList(t))

		assert.NotNil(t, clientConfigs)
		assert.Equal(t, expectedClientConfigs, len(clientConfigs))
	})
}
