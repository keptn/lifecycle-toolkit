package keptnwebhookcontroller

import (
	"time"
)

const (
	// DeploymentName is the name used for the Deployment of any webhooks and WebhookConfiguration objects.
	DeploymentName             = "klc-controller-manager"
	SuccessDuration            = 3 * time.Hour
	secretPostfix              = "-certs"
	certificatesSecretEmptyErr = "certificates secret is empty"
	couldNotUpdateCRDErr       = "could not update crd config"
)
