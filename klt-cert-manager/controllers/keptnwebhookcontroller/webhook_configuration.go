package keptnwebhookcontroller

import admissionregistrationv1 "k8s.io/api/admissionregistration/v1"

func getClientConfigsFromMutatingWebhook(mutatingWebhookConfigList *admissionregistrationv1.MutatingWebhookConfigurationList) []*admissionregistrationv1.WebhookClientConfig {
	if mutatingWebhookConfigList == nil {
		return nil
	}

	mutatingWebhookConfigs := []*admissionregistrationv1.WebhookClientConfig{}
	for i := range mutatingWebhookConfigList.Items {
		for j := range mutatingWebhookConfigList.Items[i].Webhooks {
			mutatingWebhookConfigs = append(mutatingWebhookConfigs, &mutatingWebhookConfigList.Items[i].Webhooks[j].ClientConfig)
		}
	}
	return mutatingWebhookConfigs
}

func getClientConfigsFromValidatingWebhook(validatingWebhookConfigList *admissionregistrationv1.ValidatingWebhookConfigurationList) []*admissionregistrationv1.WebhookClientConfig {
	if validatingWebhookConfigList == nil {
		return nil
	}

	validatingWebhookConfigs := []*admissionregistrationv1.WebhookClientConfig{}
	for i := range validatingWebhookConfigList.Items {
		for j := range validatingWebhookConfigList.Items[i].Webhooks {
			validatingWebhookConfigs = append(validatingWebhookConfigs, &validatingWebhookConfigList.Items[i].Webhooks[j].ClientConfig)
		}
	}
	return validatingWebhookConfigs
}
