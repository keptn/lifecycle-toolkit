package prometheus

import (
	"context"
	"errors"
	"fmt"
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

const apiKey = "ACCESS_TOKEN"

var ErrSecretKeyRefNotDefined = errors.New("the SecretKeyRef property with the Prometheus API Key is missing")

type transport struct {
	underlyingTransport http.RoundTripper
	apiToken            string
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	bearer := "Bearer " + t.apiToken
	req.Header.Add("Authorization", bearer)
	return t.underlyingTransport.RoundTrip(req)
}

func getPrometheusSecret(ctx context.Context, provider metricsapi.KeptnMetricsProvider, k8sClient client.Client) (string, error) {
	if !provider.HasSecretDefined() {
		return "", ErrSecretKeyRefNotDefined
	}
	secret := &corev1.Secret{}
	if err := k8sClient.Get(ctx, types.NamespacedName{Name: provider.Spec.SecretKeyRef.Name, Namespace: provider.Namespace}, secret); err != nil {
		return "", err
	}

	apiKeyVal := secret.Data[provider.Spec.SecretKeyRef.Key]
	if len(apiKeyVal) == 0 {
		return "", fmt.Errorf("secret does not contain %s", apiKey)
	}
	return strings.Trim(string(apiKeyVal), "\n"), nil
}
