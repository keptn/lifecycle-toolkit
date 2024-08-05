package kubeutils

import (
	"context"
	"reflect"
	"testing"

	"github.com/go-logr/logr/testr"
	"github.com/keptn/lifecycle-toolkit/keptn-cert-manager/fake"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestSecretQuery(t *testing.T) {
	t.Run(`Get secret`, testGetSecret)
	t.Run(`Create secret`, testCreateSecret)
	t.Run(`Update secret`, testUpdateSecret)
	t.Run(`Create or update secret`, testCreateOrUpdateSecret)
	t.Run(`Identical secret is not updated`, testIdenticalSecretIsNotUpdated)
	t.Run(`Update secret when data has changed`, testUpdateSecretWhenDataChanged)
	t.Run(`Update secret when labels have changed`, testUpdateSecretWhenLabelsChanged)
	t.Run(`Create secret in target namespace`, testCreateSecretInTargetNamespace)
	t.Run(`New Secret`, Secretwithmutiplekeys)
}

func testGetSecret(t *testing.T) {
	secret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testSecretName,
			Namespace: testNamespace,
		},
		Data: map[string][]byte{testKey1: []byte(testSecretValue)},
	}
	fakeClient := fake.NewClient(&secret)

	secretQuery := NewSecretQuery(fakeClient, fakeClient, testr.New(t))

	foundSecret, err := secretQuery.Get(context.TODO(), client.ObjectKey{Name: testSecretName, Namespace: testNamespace})

	assert.NoError(t, err)
	assert.True(t, AreSecretsEqual(secret, foundSecret))
}

func testCreateSecret(t *testing.T) {
	fakeClient := fake.NewClient()

	secretQuery := NewSecretQuery(fakeClient, fakeClient, testr.New(t))
	secret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testSecretName,
			Namespace: testNamespace,
		},
		Data: map[string][]byte{testKey1: []byte(testSecretValue)},
	}

	err := secretQuery.Create(context.TODO(), secret)

	assert.NoError(t, err)

	var actualSecret corev1.Secret
	err = fakeClient.Get(context.TODO(), client.ObjectKey{Name: testSecretName, Namespace: testNamespace}, &actualSecret)

	assert.NoError(t, err)
	assert.True(t, AreSecretsEqual(secret, actualSecret))
}

func testUpdateSecret(t *testing.T) {
	secret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testSecretName,
			Namespace: testNamespace,
		},
		Data: map[string][]byte{testKey1: []byte(testSecretValue)},
	}
	fakeClient := fake.NewClient()

	secretQuery := NewSecretQuery(fakeClient, fakeClient, testr.New(t))

	err := secretQuery.Update(context.TODO(), secret)

	assert.Error(t, err)

	secret.Data = nil
	fakeClient = fake.NewClient(&secret)

	secretQuery.kubeClient = fakeClient

	err = secretQuery.Update(context.TODO(), secret)

	assert.NoError(t, err)

	var updatedSecret corev1.Secret
	err = fakeClient.Get(context.TODO(), client.ObjectKey{Name: secret.Name, Namespace: secret.Namespace}, &updatedSecret)

	assert.NoError(t, err)
	assert.True(t, AreSecretsEqual(secret, updatedSecret))
}

func testCreateOrUpdateSecret(t *testing.T) {
	secret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testSecretName,
			Namespace: testNamespace,
		},
		Data: map[string][]byte{testKey1: []byte(testSecretValue)},
	}
	fakeClient := fake.NewClient()
	secretQuery := NewSecretQuery(fakeClient, fakeClient, testr.New(t))

	err := secretQuery.CreateOrUpdate(context.TODO(), secret)
	assert.NoError(t, err)

	var createdSecret corev1.Secret
	err = fakeClient.Get(context.TODO(), client.ObjectKey{Name: secret.Name, Namespace: secret.Namespace}, &createdSecret)

	assert.NoError(t, err)
	assert.True(t, AreSecretsEqual(secret, createdSecret))

	fakeClient = fake.NewClient(&secret)

	secret = corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testSecretName,
			Namespace: testNamespace,
		},
		Data: nil,
	}
	secretQuery.kubeClient = fakeClient

	err = secretQuery.CreateOrUpdate(context.TODO(), secret)

	assert.NoError(t, err)

	var updatedSecret corev1.Secret
	err = fakeClient.Get(context.TODO(), client.ObjectKey{Name: secret.Name, Namespace: secret.Namespace}, &updatedSecret)

	assert.NoError(t, err)
	assert.True(t, AreSecretsEqual(secret, updatedSecret))
}

func testIdenticalSecretIsNotUpdated(t *testing.T) {
	data := map[string][]byte{testKey1: []byte(testValue1)}
	labels := map[string]string{
		"label": "test",
	}
	fakeClient := fake.NewClient(&corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testSecretName,
			Namespace: testNamespace,
			Labels:    labels,
		},
		Data: data,
	})

	secret := createTestSecret(labels, data)
	secretQuery := NewSecretQuery(fakeClient, fakeClient, testr.New(t))

	err := secretQuery.CreateOrUpdate(context.TODO(), *secret)
	assert.NoError(t, err)
}

func testUpdateSecretWhenDataChanged(t *testing.T) {
	data := map[string][]byte{testKey1: []byte(testValue1)}
	labels := map[string]string{
		"label": "test",
	}
	fakeClient := fake.NewClient(&corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testSecretName,
			Namespace: testNamespace,
			Labels:    labels,
		},
		Data: map[string][]byte{},
	})

	secret := createTestSecret(labels, data)
	secretQuery := NewSecretQuery(fakeClient, fakeClient, testr.New(t))

	err := secretQuery.CreateOrUpdate(context.TODO(), *secret)
	assert.NoError(t, err)

	var updatedSecret corev1.Secret
	err = fakeClient.Get(context.TODO(), types.NamespacedName{Name: testSecretName, Namespace: testNamespace}, &updatedSecret)

	assert.NoError(t, err)
	assert.True(t, reflect.DeepEqual(data, updatedSecret.Data))
}

func testUpdateSecretWhenLabelsChanged(t *testing.T) {
	data := map[string][]byte{testKey1: []byte(testValue1)}
	labels := map[string]string{
		"label": "test",
	}
	fakeClient := fake.NewClient(&corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testSecretName,
			Namespace: testNamespace,
			Labels:    map[string]string{},
		},
		Data: data,
	})

	secret := createTestSecret(labels, data)
	secretQuery := NewSecretQuery(fakeClient, fakeClient, testr.New(t))

	err := secretQuery.CreateOrUpdate(context.TODO(), *secret)
	assert.NoError(t, err)

	var updatedSecret corev1.Secret
	err = fakeClient.Get(context.TODO(), types.NamespacedName{Name: testSecretName, Namespace: testNamespace}, &updatedSecret)

	assert.NoError(t, err)
	assert.True(t, reflect.DeepEqual(labels, updatedSecret.Labels))
}

func testCreateSecretInTargetNamespace(t *testing.T) {
	data := map[string][]byte{testKey1: []byte(testValue1)}
	labels := map[string]string{
		"label": "test",
	}
	fakeClient := fake.NewClient(&corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testSecretName,
			Namespace: "other",
		},
		Data: map[string][]byte{},
	})

	secret := createTestSecret(labels, data)
	secretQuery := NewSecretQuery(fakeClient, fakeClient, testr.New(t))

	err := secretQuery.CreateOrUpdate(context.TODO(), *secret)

	assert.NoError(t, err)

	var newSecret corev1.Secret
	err = fakeClient.Get(context.TODO(), types.NamespacedName{Name: testSecretName, Namespace: testNamespace}, &newSecret)

	assert.NoError(t, err)
	assert.True(t, reflect.DeepEqual(data, newSecret.Data))
	assert.True(t, reflect.DeepEqual(labels, newSecret.Labels))
	assert.Equal(t, testSecretName, newSecret.Name)
	assert.Equal(t, testNamespace, newSecret.Namespace)
	assert.Equal(t, corev1.SecretTypeOpaque, newSecret.Type)
}

func createTestSecret(labels map[string]string, data map[string][]byte) *corev1.Secret {
	secret := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      testSecretName,
			Namespace: testNamespace,
			Labels:    labels,
		},
		Data: data,
		Type: corev1.SecretTypeOpaque,
	}
	return secret
}

func Secretwithmutiplekeys(t *testing.T) {

	testCases := []struct {
		name      string
		namespace string
		data      map[string][]byte
	}{
		{
			name:      "test-secret-1",
			namespace: "test-namespace-1",
			data: map[string][]byte{
				"key1": []byte("value1"),
				"key2": []byte("value2"),
			},
		},
		{
			name:      "test-secret-2",
			namespace: "test-namespace-2",
			data: map[string][]byte{
				"key1": []byte("value1"),
			},
		},
		{
			name:      "test-secret-3",
			namespace: "test-namespace-3",
			data:      map[string][]byte{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			secret := NewSecret(tc.name, tc.namespace, tc.data)

			assert.Equal(t, tc.name, secret.Name)
			assert.Equal(t, tc.namespace, secret.Namespace)
			assert.Equal(t, tc.data, secret.Data)
		})
	}
}
