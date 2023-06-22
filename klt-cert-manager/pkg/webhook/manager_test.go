package webhook

import (
	"testing"

	"github.com/keptn/lifecycle-toolkit/klt-cert-manager/pkg/fake"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

func TestCreateOptions(t *testing.T) {

	t.Run("implements interface", func(t *testing.T) {
		var provider Provider = NewWebhookManagerProvider("certs-dir", "key-file", "cert-file")

		providerImpl := provider.(WebhookProvider)
		assert.Equal(t, "certs-dir", providerImpl.certificateDirectory)
		assert.Equal(t, "key-file", providerImpl.keyFileName)
		assert.Equal(t, "cert-file", providerImpl.certificateFileName)
	})
	t.Run("creates options", func(t *testing.T) {
		provider := WebhookProvider{}
		options := provider.createOptions(scheme.Scheme, "test-namespace")

		assert.NotNil(t, options)
		assert.Equal(t, "test-namespace", options.Namespace)
		assert.Equal(t, scheme.Scheme, options.Scheme)
		assert.Equal(t, metricsBindAddress, options.MetricsBindAddress)
		assert.Equal(t, port, options.Port)
	})
	t.Run("configures webhooks server", func(t *testing.T) {
		provider := NewWebhookManagerProvider("certs-dir", "key-file", "cert-file")
		expectedWebhookServer := webhook.NewServer(webhook.Options{})

		mgr := &fake.MockManager{
			GetWebhookServerFunc: func() webhook.Server {
				return expectedWebhookServer
			},
		}

		provider.SetupWebhookServer(mgr)

		assert.Equal(t, "certs-dir", expectedWebhookServer.(*webhook.DefaultServer).Options.CertDir)
		assert.Equal(t, "key-file", expectedWebhookServer.(*webhook.DefaultServer).Options.KeyName)
		assert.Equal(t, "cert-file", expectedWebhookServer.(*webhook.DefaultServer).Options.CertName)

		mgrWebhookServer := mgr.GetWebhookServer()
		assert.Equal(t, "certs-dir", mgrWebhookServer.(*webhook.DefaultServer).Options.CertDir)
		assert.Equal(t, "key-file", mgrWebhookServer.(*webhook.DefaultServer).Options.KeyName)
		assert.Equal(t, "cert-file", mgrWebhookServer.(*webhook.DefaultServer).Options.CertName)
	})
}
