package webhook

import (
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

const (
	metricsBindAddress = ":8383"
	port               = 8443
)

type WebhookProvider struct {
	certificateDirectory string
	certificateFileName  string
	keyFileName          string
}

func NewWebhookManagerProvider(certificateDirectory string, keyFileName string, certificateFileName string) WebhookProvider {
	return WebhookProvider{
		certificateDirectory: certificateDirectory,
		certificateFileName:  certificateFileName,
		keyFileName:          keyFileName,
	}
}

func (provider WebhookProvider) CreateManager(namespace string, scheme *runtime.Scheme, config *rest.Config) (manager.Manager, error) {
	mgr, err := ctrl.NewManager(config, provider.createOptions(scheme, namespace))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return provider.setupWebhookServer(mgr), nil
}

func (provider WebhookProvider) createOptions(scheme *runtime.Scheme, namespace string) ctrl.Options {
	return ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsBindAddress,
		Port:               port,
		Namespace:          namespace,
	}
}

func (provider WebhookProvider) setupWebhookServer(mgr manager.Manager) manager.Manager {
	webhookServer := mgr.GetWebhookServer()
	webhookServer.CertDir = provider.certificateDirectory
	webhookServer.KeyName = provider.keyFileName
	webhookServer.CertName = provider.certificateFileName

	return mgr
}
