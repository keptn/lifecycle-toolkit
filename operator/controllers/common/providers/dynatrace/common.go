package dynatrace

import (
	"context"
	"errors"
	"fmt"
	klcv1alpha2 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

var ErrSecretKeyRefNotDefined = errors.New("the SecretKeyRef property with the Dynatrace token is missing")
var ErrOAuthTokenInvalid = errors.New("the Dynatrace token has an invalid format")
var ErrInvalidResult = errors.New("the answer does not contain any data")

const dtTokenPrefix = "dts08"

func getDTSecret(ctx context.Context, provider klcv1alpha2.KeptnEvaluationProvider, k8sClient client.Client) (string, error) {
	if !provider.HasSecretDefined() {
		return "", ErrSecretKeyRefNotDefined
	}
	dtCredsSecret := &corev1.Secret{}
	if err := k8sClient.Get(ctx, types.NamespacedName{Name: provider.Spec.SecretKeyRef.Name, Namespace: provider.Namespace}, dtCredsSecret); err != nil {
		return "", err
	}

	token := dtCredsSecret.Data[provider.Spec.SecretKeyRef.Key]
	if len(token) == 0 {
		return "", fmt.Errorf("secret contains invalid key %s", provider.Spec.SecretKeyRef.Key)
	}
	return string(token), nil
}

func validateOAuthSecret(token string) error {
	// must start with dt0s08
	// must have 2 dots
	// third part (split by dot) must be 64 chars
	if !strings.HasPrefix(token, dtTokenPrefix) {
		return ErrOAuthTokenInvalid
	}
	split := strings.Split(token, ".")
	if len(split) != 3 {
		return ErrOAuthTokenInvalid
	}
	secret := split[2]
	if len(secret) != 64 {
		return ErrOAuthTokenInvalid
	}
	return nil
}

// consider a handle http client call func
