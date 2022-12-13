package keptnwebhookcontroller

const (

	// SecretCertsName is the name of the secret where the webhook certificates are stored.
	SecretCertsName = "klc-controller-manager-certs"

	// DeploymentName is the name used for the Deployment of any webhooks and WebhookConfiguration objects.
	DeploymentName = "klc-mutating-webhook-configuration"
)
