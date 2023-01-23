package keptnwebhookcontroller

import admissionregistrationv1 "k8s.io/api/admissionregistration/v1"

func getClientConfigsFromMutatingWebhook(mutatingWebhookConfig *admissionregistrationv1.MutatingWebhookConfiguration) []*admissionregistrationv1.WebhookClientConfig {
	if mutatingWebhookConfig == nil {
		return nil
	}

	mutatingWebhookConfigs := make([]*admissionregistrationv1.WebhookClientConfig, len(mutatingWebhookConfig.Webhooks))
	for i := range mutatingWebhookConfig.Webhooks {
		mutatingWebhookConfigs[i] = &mutatingWebhookConfig.Webhooks[i].ClientConfig
	}
	return mutatingWebhookConfigs
}

func getClientConfigsFromValidatingWebhook(validatingWebhookConfig *admissionregistrationv1.ValidatingWebhookConfiguration) []*admissionregistrationv1.WebhookClientConfig {
	if validatingWebhookConfig == nil {
		return nil
	}

	mutatingWebhookConfigs := make([]*admissionregistrationv1.WebhookClientConfig, len(validatingWebhookConfig.Webhooks))
	for i := range validatingWebhookConfig.Webhooks {
		mutatingWebhookConfigs[i] = &validatingWebhookConfig.Webhooks[i].ClientConfig
	}
	return mutatingWebhookConfigs
}
