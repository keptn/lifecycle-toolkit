package keptnwebhookcontroller

import (
	"context"
	"testing"

	"github.com/go-logr/logr/testr"
	"github.com/keptn/lifecycle-toolkit/keptn-cert-manager/fake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	testKey = "test-key"
)

var testValue1 = []byte{1, 2, 3, 4}
var testValue2 = []byte{5, 6, 7, 8}

func TestSetSecretFromReader(t *testing.T) {
	t.Run(`fill with empty secret if secret does not exist`, func(t *testing.T) {
		fakeclient := fake.NewClient()
		certSecret := newCertificateSecret(fakeclient)
		err := certSecret.setSecretFromReader(context.TODO(), testnamespace, testr.New(t))

		assert.NoError(t, err)
		assert.False(t, certSecret.existsInCluster)
		assert.NotNil(t, certSecret.secret)
	})
	t.Run(`find existing secret`, func(t *testing.T) {

		fakeclient := fake.NewClient(
			createTestSecret(t, createInvalidTestCertData(t)))
		certSecret := newCertificateSecret(fakeclient)
		err := certSecret.setSecretFromReader(context.TODO(), testnamespace, testr.New(t))

		assert.NoError(t, err)
		assert.True(t, certSecret.existsInCluster)
		assert.NotNil(t, certSecret.secret)
	})
}

func TestIsRecent(t *testing.T) {
	t.Run(`true if certs and secret are nil`, func(t *testing.T) {
		certSecret := newCertificateSecret(nil)

		assert.True(t, certSecret.isRecent())
	})
	t.Run(`false if only one is nil`, func(t *testing.T) {
		certSecret := newCertificateSecret(nil)
		certSecret.secret = &corev1.Secret{}

		assert.False(t, certSecret.isRecent())

		certSecret.secret = nil
		certSecret.certificates = &Certs{}

		assert.False(t, certSecret.isRecent())
	})
	t.Run(`true if data is equal, false otherwise`, func(t *testing.T) {
		certSecret := newCertificateSecret(nil)
		secret := corev1.Secret{
			Data: map[string][]byte{testKey: testValue1},
		}
		certs := Certs{
			Data: map[string][]byte{testKey: testValue1},
		}
		certSecret.secret = &secret
		certSecret.certificates = &certs

		assert.True(t, certSecret.isRecent())

		certSecret.secret.Data = map[string][]byte{testKey: testValue2}

		assert.False(t, certSecret.isRecent())
	})
}

func TestAreConfigsValid(t *testing.T) {
	t.Run(`true if no configs were given`, func(t *testing.T) {
		certSecret := newCertificateSecret(nil)

		assert.True(t, certSecret.areWebhookConfigsValid(nil))
		assert.True(t, certSecret.areWebhookConfigsValid(make([]*admissionregistrationv1.WebhookClientConfig, 0)))
	})
	t.Run(`true if all CABundle matches certificate data, false otherwise`, func(t *testing.T) {
		certSecret := newCertificateSecret(nil)
		certSecret.certificates = &Certs{
			Data: map[string][]byte{RootCert: testValue1},
		}
		webhookConfigs := make([]*admissionregistrationv1.WebhookClientConfig, 1)
		webhookConfigs = append(webhookConfigs, &admissionregistrationv1.WebhookClientConfig{
			CABundle: testValue1,
		})
		webhookConfigs = append(webhookConfigs, &admissionregistrationv1.WebhookClientConfig{
			CABundle: testValue1,
		})
		webhookConfigs = append(webhookConfigs, &admissionregistrationv1.WebhookClientConfig{
			CABundle: testValue1,
		})

		assert.True(t, certSecret.areWebhookConfigsValid(webhookConfigs))

		webhookConfigs = append(webhookConfigs, &admissionregistrationv1.WebhookClientConfig{
			CABundle: testValue2,
		})

		assert.False(t, certSecret.areWebhookConfigsValid(webhookConfigs))
	})
}

func TestCreateOrUpdateIfNecessary(t *testing.T) {
	t.Run(`do nothing if certificate is recent and exists`, func(t *testing.T) {
		fakeClient := fake.NewClient()
		certSecret := newCertificateSecret(fakeClient)
		certSecret.existsInCluster = true

		err := certSecret.createOrUpdateIfNecessary(context.TODO())

		assert.NoError(t, err)

		err = fakeClient.Get(context.TODO(), client.ObjectKey{Name: buildSecretName()}, &corev1.Secret{})

		assert.Error(t, err)
		assert.True(t, k8serrors.IsNotFound(err))
	})
	t.Run(`create if secret does not exist`, func(t *testing.T) {
		fakeClient := fake.NewClient()
		certSecret := newCertificateSecret(fakeClient)
		certSecret.secret = &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      buildSecretName(),
				Namespace: testnamespace,
			},
		}
		certSecret.certificates = &Certs{
			Data: map[string][]byte{testKey: testValue1},
		}

		err := certSecret.createOrUpdateIfNecessary(context.TODO())

		assert.NoError(t, err)

		newSecret := corev1.Secret{}
		err = fakeClient.Get(context.TODO(), client.ObjectKey{Name: buildSecretName(), Namespace: testnamespace}, &newSecret)

		assert.NoError(t, err)
		assert.NotNil(t, newSecret)
		assert.EqualValues(t, certSecret.certificates.Data, newSecret.Data)
	})
	t.Run(`update if secret exists`, func(t *testing.T) {
		fakeClient := fake.NewClient()
		certSecret := newCertificateSecret(fakeClient)
		certSecret.secret = &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      buildSecretName(),
				Namespace: testnamespace,
			},
		}
		certSecret.certificates = &Certs{
			Data: map[string][]byte{testKey: testValue1},
		}

		err := certSecret.createOrUpdateIfNecessary(context.TODO())

		require.NoError(t, err)

		newSecret := corev1.Secret{}
		err = fakeClient.Get(context.TODO(), client.ObjectKey{Name: buildSecretName(), Namespace: testnamespace}, &newSecret)

		require.NoError(t, err)
		require.NotNil(t, newSecret)
		require.EqualValues(t, certSecret.certificates.Data, newSecret.Data)

		certSecret.secret = &newSecret
		certSecret.certificates.Data = map[string][]byte{testKey: testValue2}
		certSecret.existsInCluster = true
		err = certSecret.createOrUpdateIfNecessary(context.TODO())

		assert.NoError(t, err)

		err = fakeClient.Get(context.TODO(), client.ObjectKey{Name: buildSecretName(), Namespace: testnamespace}, &newSecret)

		assert.NoError(t, err)
		assert.NotNil(t, newSecret)
		assert.EqualValues(t, certSecret.certificates.Data, newSecret.Data)
	})
}
