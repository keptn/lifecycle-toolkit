package kubeutils

import (
	"context"
	"reflect"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type SecretQuery struct {
	kubeQuery
}

func NewSecretQuery(kubeClient client.Client, kubeReader client.Reader, log logr.Logger) SecretQuery {
	return SecretQuery{
		newKubeQuery(kubeClient, kubeReader, log),
	}
}

func (query SecretQuery) Get(ctx context.Context, objectKey client.ObjectKey) (corev1.Secret, error) {
	var secret corev1.Secret
	err := query.kubeReader.Get(ctx, objectKey, &secret)

	return secret, errors.WithStack(err)
}

func (query SecretQuery) Create(ctx context.Context, secret corev1.Secret) error {
	query.log.Info("creating secret", "name", secret.Name, "namespace", secret.Namespace)

	return errors.WithStack(query.kubeClient.Create(ctx, &secret))
}

func (query SecretQuery) Update(ctx context.Context, secret corev1.Secret) error {
	query.log.Info("updating secret", "name", secret.Name, "namespace", secret.Namespace)

	return errors.WithStack(query.kubeClient.Update(ctx, &secret))
}

func (query SecretQuery) CreateOrUpdate(ctx context.Context, secret corev1.Secret) error {
	currentSecret, err := query.Get(ctx, types.NamespacedName{Name: secret.Name, Namespace: secret.Namespace})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			err = query.Create(ctx, secret)
			if err != nil {
				return errors.WithStack(err)
			}
			return nil
		}
		return errors.WithStack(err)
	}

	if AreSecretsEqual(secret, currentSecret) {
		query.log.Info("secret unchanged", "name", secret.Name, "namespace", secret.Namespace)
		return nil
	}

	err = query.Update(ctx, secret)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func AreSecretsEqual(secret corev1.Secret, other corev1.Secret) bool {
	return reflect.DeepEqual(secret.Data, other.Data) && reflect.DeepEqual(secret.Labels, other.Labels)
}

func NewSecret(name string, namespace string, data map[string][]byte) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: data,
	}
}
