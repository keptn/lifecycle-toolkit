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
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

func TestWebhookCommandBuilder(t *testing.T) {

	t.Run("set namespace", func(t *testing.T) {
		builder := NewWebhookServerBuilder().SetNamespace("namespace")

		assert.Equal(t, "namespace", builder.namespace)
	})
}

func TestBuilder_Run(t *testing.T) {

	mockManager := &fake.MockManager{}

	mockManager.GetAPIReaderFunc = func() client.Reader {
		return newFakeClient()
	}

	mockManager.GetWebhookServerFunc = func() webhook.Server {
		return webhook.NewServer(webhook.Options{})
	}

	mockManager.StartFunc = func(ctx context.Context) error {
		return nil
	}

	mockCertificateWatcher := &fake2.MockCertificateWatcher{}
	mockCertificateWatcher.WaitForCertificatesFunc = func() {}

	builder := NewWebhookServerBuilder().SetCertificateWatcher(mockCertificateWatcher)

	webhooks := map[string]*admission.Webhook{
		"/my-url":       {},
		"/my-other-url": {},
	}
	builder.Register(mockManager, webhooks)

	require.Len(t, mockManager.GetWebhookServerCalls(), 2)
	require.Len(t, mockCertificateWatcher.WaitForCertificatesCalls(), 1)
}

func newFakeClient(objs ...client.Object) client.Reader {
	return fakeClient.NewClientBuilder().WithObjects(objs...).Build()
}
