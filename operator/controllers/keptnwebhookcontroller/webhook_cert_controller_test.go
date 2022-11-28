package keptnwebhookcontroller

import (
	"context"
	"github.com/go-logr/logr/testr"
	"testing"
	"time"

	"github.com/keptn/lifecycle-toolkit/operator/controllers/common/fake"
	"github.com/keptn/lifecycle-toolkit/operator/webhooks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	testNamespace = "test-namespace"
	testDomain    = webhooks.DeploymentName + "." + testNamespace + ".svc"

	expectedSecretName = webhooks.DeploymentName + secretPostfix

	testBytes = 123
)

func TestReconcileCertificate_Create(t *testing.T) {
	clt := prepareFakeClient(false, false)
	controller, request := prepareController(t, clt)

	res, err := controller.Reconcile(context.TODO(), request)
	require.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, SuccessDuration, res.RequeueAfter)

	secret := &corev1.Secret{}
	err = clt.Get(context.TODO(), client.ObjectKey{Name: expectedSecretName, Namespace: testNamespace}, secret)
	require.NoError(t, err)

	assert.NotNil(t, secret.Data)
	assert.NotEmpty(t, secret.Data)
	assert.Contains(t, secret.Data, RootKey)
	assert.Contains(t, secret.Data, RootCert)
	assert.Contains(t, secret.Data, RootCertOld)
	assert.Contains(t, secret.Data, ServerKey)
	assert.Contains(t, secret.Data, ServerCert)
	assert.NotNil(t, secret.Data[RootCert])
	assert.NotEmpty(t, secret.Data[RootCert])
	assert.Empty(t, secret.Data[RootCertOld])

	verifyCertificates(t, controller, secret, clt, false)
}

func TestReconcileCertificate_Update(t *testing.T) {
	clt := prepareFakeClient(true, false)
	controller, request := prepareController(t, clt)

	res, err := controller.Reconcile(context.TODO(), request)
	require.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, SuccessDuration, res.RequeueAfter)

	secret := &corev1.Secret{}
	err = clt.Get(context.TODO(), client.ObjectKey{Name: expectedSecretName, Namespace: testNamespace}, secret)
	require.NoError(t, err)

	assert.NotNil(t, secret.Data)
	assert.NotEmpty(t, secret.Data)
	assert.Contains(t, secret.Data, RootKey)
	assert.Contains(t, secret.Data, RootCert)
	assert.Contains(t, secret.Data, RootCertOld)
	assert.Contains(t, secret.Data, ServerKey)
	assert.Contains(t, secret.Data, ServerCert)
	assert.NotNil(t, secret.Data[RootCert])
	assert.NotEmpty(t, secret.Data[RootCert])
	assert.Equal(t, []byte{testBytes}, secret.Data[RootCertOld])

	verifyCertificates(t, controller, secret, clt, true)
}

func TestReconcileCertificate_ExistingSecretWithValidCertificate(t *testing.T) {
	clt := prepareFakeClient(true, true)
	controller, request := prepareController(t, clt)

	res, err := controller.Reconcile(context.TODO(), request)
	require.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, SuccessDuration, res.RequeueAfter)

	secret := &corev1.Secret{}
	err = clt.Get(context.TODO(), client.ObjectKey{Name: expectedSecretName, Namespace: testNamespace}, secret)
	require.NoError(t, err)

	verifyCertificates(t, controller, secret, clt, false)
}

func TestReconcile(t *testing.T) {

	t.Run(`reconcile successfully without validatingwebhookconfiguration`, func(t *testing.T) {
		fakeClient := fake.NewClient(&admissionregistrationv1.MutatingWebhookConfiguration{
			ObjectMeta: metav1.ObjectMeta{
				Name: webhooks.DeploymentName,
			},
			Webhooks: []admissionregistrationv1.MutatingWebhook{
				{
					ClientConfig: admissionregistrationv1.WebhookClientConfig{},
				},
				{
					ClientConfig: admissionregistrationv1.WebhookClientConfig{},
				},
			},
		})

		controller, request := prepareController(t, fakeClient)
		result, err := controller.Reconcile(context.TODO(), request)

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	// Generation must not be skipped because webhook startup routine listens for the secret
	// See cmd/operator/manager.go and cmd/operator/watcher.go
	t.Run(`do not skip certificates generation if no configuration exists`, func(t *testing.T) {
		fakeClient := fake.NewClient()
		controller, request := prepareController(t, fakeClient)
		result, err := controller.Reconcile(context.TODO(), request)

		assert.NoError(t, err)
		assert.NotNil(t, result)

		secret := &corev1.Secret{}
		err = fakeClient.Get(context.TODO(), client.ObjectKey{Name: expectedSecretName, Namespace: testNamespace}, secret)
		assert.NoError(t, err)
	})
}

func prepareFakeClient(withSecret bool, generateValidSecret bool) client.Client {
	objs := []client.Object{
		&admissionregistrationv1.MutatingWebhookConfiguration{
			ObjectMeta: metav1.ObjectMeta{
				Name: webhooks.DeploymentName,
			},
			Webhooks: []admissionregistrationv1.MutatingWebhook{
				{
					ClientConfig: admissionregistrationv1.WebhookClientConfig{},
				},
				{
					ClientConfig: admissionregistrationv1.WebhookClientConfig{},
				},
			},
		},
	}
	if withSecret {
		certData := createInvalidTestCertData(nil)
		if generateValidSecret {
			certData = createValidTestCertData(nil)
		}

		objs = append(objs,
			createTestSecret(nil, certData),
		)
	}

	fake := fake.NewClient(objs...)
	return fake
}

func createInvalidTestCertData(_ *testing.T) map[string][]byte {
	return map[string][]byte{
		RootKey:    {testBytes},
		RootCert:   {testBytes},
		ServerKey:  {testBytes},
		ServerCert: {testBytes},
	}
}

func createValidTestCertData(_ *testing.T) map[string][]byte {
	cert := Certs{
		Domain: testDomain,
		Now:    time.Now(),
	}
	_ = cert.ValidateCerts()
	return cert.Data
}

func createTestSecret(_ *testing.T, certData map[string][]byte) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: testNamespace,
			Name:      expectedSecretName,
		},
		Data: certData,
	}
}

func prepareController(t *testing.T, clt client.Client) (*KeptnWebhookCertificateReconciler, reconcile.Request) {
	rec := &KeptnWebhookCertificateReconciler{
		ctx:       context.TODO(),
		Client:    clt,
		ApiReader: clt,
		namespace: testNamespace,
		Log:       testr.New(t),
	}

	request := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      webhooks.DeploymentName,
			Namespace: testNamespace,
		},
	}

	return rec, request
}

func testWebhookClientConfig(
	t *testing.T, webhookClientConfig *admissionregistrationv1.WebhookClientConfig,
	secretData map[string][]byte, isUpdate bool) {
	require.NotNil(t, webhookClientConfig)
	require.NotEmpty(t, webhookClientConfig.CABundle)

	expectedCert := secretData[RootCert]
	if isUpdate {
		expectedCert = append(expectedCert, []byte{123}...)
	}
	assert.Equal(t, expectedCert, webhookClientConfig.CABundle)
}

func verifyCertificates(t *testing.T, rec *KeptnWebhookCertificateReconciler, secret *corev1.Secret, clt client.Client, isUpdate bool) {
	cert := Certs{
		Domain:  getDomain(rec.namespace),
		Data:    secret.Data,
		SrcData: secret.Data,
		Now:     time.Now(),
	}

	// validateRootCerts and validateServerCerts return false if the certificates are valid
	assert.False(t, cert.validateRootCerts(time.Now()))
	assert.False(t, cert.validateServerCerts(time.Now()))

	mutatingWebhookConfig := &admissionregistrationv1.MutatingWebhookConfiguration{}
	err := clt.Get(context.TODO(), client.ObjectKey{
		Name: webhooks.DeploymentName,
	}, mutatingWebhookConfig)
	require.NoError(t, err)
	assert.Len(t, mutatingWebhookConfig.Webhooks, 2)
	testWebhookClientConfig(t, &mutatingWebhookConfig.Webhooks[0].ClientConfig, secret.Data, isUpdate)

}
