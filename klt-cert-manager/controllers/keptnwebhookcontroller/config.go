package keptnwebhookcontroller

import (
	"time"
)

const (
	// DeploymentName is the name used for the Deployment of any webhooks and WebhookConfiguration objects.
	DeploymentName             = "klc-controller-manager"
	ServiceName                = "klc-webhook-service"
	SuccessDuration            = 3 * time.Hour
	MutatingWebhookconfig      = "klc-mutating-webhook-configuration"
	ValidatingWebhookconfig    = "klc-validating-webhook-configuration"
	secretPostfix              = "-certs"
	crdGroup                   = "lifecycle.keptn.sh"
	certificatesSecretEmptyErr = "certificates secret is empty"
	couldNotUpdateCRDErr       = "could not update crd config"
)
