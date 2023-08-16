package dynatrace

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var ErrSecretKeyRefNotDefined = errors.New("the SecretKeyRef property with the Dynatrace token is missing")
var ErrInvalidResult = errors.New("the answer does not contain any data")
var ErrDQLQueryTimeout = errors.New("timed out waiting for result of DQL query")

const ErrAPIMsg = "provider api response: %s"

type Error struct {
	Code    int    `json:"-"` // optional
	Message string `json:"message"`
}

func getDTSecret(ctx context.Context, provider metricsapi.KeptnMetricsProvider, k8sClient client.Client) (string, error) {
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

func urlEncodeQuery(query string) string {
	params := strings.Split(query, "&")

	result := ""

	for i, param := range params {
		keyAndValue := strings.Split(param, "=")
		if len(keyAndValue) == 2 {
			encodedKeyAndValue := keyAndValue[0] + "=" + url.QueryEscape(keyAndValue[1])
			if i == 0 {
				result += encodedKeyAndValue
			} else {
				result += "&" + encodedKeyAndValue
			}
		}
	}

	return result
}
