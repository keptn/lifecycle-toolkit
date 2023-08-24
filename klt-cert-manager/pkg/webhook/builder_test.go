package webhook

import (
	"context"
	"testing"

	fake2 "github.com/keptn/lifecycle-toolkit/klt-cert-manager/pkg/certificates/fake"
	"github.com/keptn/lifecycle-toolkit/klt-cert-manager/pkg/fake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/controller-runtime/pkg/client"
	fakeClient "sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

func TestWebhookCommandBuilder(t *testing.T) {
	t.Run("set manager provider", func(t *testing.T) {
		expectedProvider := &fake.MockWebhookManager{}
		builder := NewWebhookBuilder().SetManagerProvider(expectedProvider)

		assert.Equal(t, expectedProvider, builder.managerProvider)
	})
	t.Run("set namespace", func(t *testing.T) {
		builder := NewWebhookBuilder().SetNamespace("namespace")

		assert.Equal(t, "namespace", builder.namespace)
	})
	t.Run("set certificate watcher", func(t *testing.T) {
		expectedCertificateWatcher := &fake2.MockCertificateWatcher{}
		builder := NewWebhookBuilder().SetCertificateWatcher(expectedCertificateWatcher)

		assert.Equal(t, expectedCertificateWatcher, builder.certificateWatcher)
	})
}

func TestBuilder_Run(t *testing.T) {
	mockProvider := &fake.MockWebhookManager{}

	mockProvider.SetupWebhookServerFunc = func(mgr manager.Manager) {}

	mockManager := &fake.MockManager{}

	mockManager.GetAPIReaderFunc = func() client.Reader {
		return newFakeClient()
	}
	webhookServer := &webhook.Server{}
	mockManager.GetWebhookServerFunc = func() *webhook.Server {
		return webhookServer
	}

	mockManager.StartFunc = func(ctx context.Context) error {
		return nil
	}

	mockCertificateWatcher := &fake2.MockCertificateWatcher{}
	mockCertificateWatcher.WaitForCertificatesFunc = func() {}

	builder := NewWebhookBuilder().SetManagerProvider(mockProvider).SetCertificateWatcher(mockCertificateWatcher)

	webhooks := map[string]*admission.Webhook{
		"/my-url":       {},
		"/my-other-url": {},
	}
	err := builder.Run(mockManager, webhooks)

	require.Nil(t, err)

	require.Len(t, mockProvider.SetupWebhookServerCalls(), 1)
	require.Len(t, mockManager.GetWebhookServerCalls(), 2)
	require.Len(t, mockCertificateWatcher.WaitForCertificatesCalls(), 1)
}

func newFakeClient(objs ...client.Object) client.Reader {
	return fakeClient.NewClientBuilder().WithObjects(objs...).Build()
}
