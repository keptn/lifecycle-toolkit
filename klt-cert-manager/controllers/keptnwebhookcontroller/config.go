package keptnwebhookcontroller

import (
	"time"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// DeploymentName is the name used for the Deployment of any webhooks and WebhookConfiguration objects.
	DeploymentName             = "klc-controller-manager"
	ServiceName                = "klc-webhook-service"
	SuccessDuration            = 3 * time.Hour
	Webhookconfig              = "klc-mutating-webhook-configuration"
	secretPostfix              = "-certs"
	crdGroup                   = "lifecycle.keptn.sh"
	certificatesSecretEmptyErr = "certificates secret is empty"
	couldNotUpdateCRDErr       = "could not update crd config"
	namespace                  = "keptn-lifecycle-toolkit-system"
)

var keptnCRD = schema.GroupVersionResource{Group: "lifecycle.keptn.sh"}
var crdversion = schema.GroupVersion{
	Group:   "apiextensions.k8s.io",
	Version: "v1",
}
