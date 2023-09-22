package prometheus

import (
	"context"
	"errors"
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	promapi "github.com/prometheus/client_golang/api"
	"github.com/prometheus/common/config"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const userName = "user"
const password = "password"

var ErrSecretKeyRefNotDefined = errors.New("the SecretKeyRef property with the Prometheus API Key is missing")
var ErrInvalidSecretFormat = errors.New("secret key does not contain user and password")

type SecretData struct {
	User     string        `json:"user"`
	Password config.Secret `json:"password"`
}

func getRoundtripper(ctx context.Context, provider metricsapi.KeptnMetricsProvider, k8sClient client.Client) (http.RoundTripper, error) {
	secret, err := getPrometheusSecret(ctx, provider, k8sClient)
	if err != nil {
		if errors.Is(err, ErrSecretKeyRefNotDefined) {
			return promapi.DefaultRoundTripper, nil
		}
		return nil, err
	}
	return config.NewBasicAuthRoundTripper(secret.User, secret.Password, "", promapi.DefaultRoundTripper), nil

}

func getPrometheusSecret(ctx context.Context, provider metricsapi.KeptnMetricsProvider, k8sClient client.Client) (*SecretData, error) {
	if !provider.HasSecretDefined() {
		return nil, ErrSecretKeyRefNotDefined
	}
	secret := &corev1.Secret{}
	if err := k8sClient.Get(ctx, types.NamespacedName{Name: provider.Spec.SecretKeyRef.Name, Namespace: provider.Namespace}, secret); err != nil {
		return nil, err
	}

	var secretData SecretData
	if data, ok := secret.Data[provider.Spec.SecretKeyRef.Key]; ok {
		// Unmarshal the JSON data into the SecretData struct
		if err := json.Unmarshal(data, &secretData); err != nil {
			return nil, ErrInvalidSecretFormat
		}
	} else if _, ok := secret.Data[userName]; !ok {
		return nil, ErrInvalidSecretFormat
	}
	secretData.User = string(secret.Data[userName])
	secretData.Password = config.Secret(secret.Data[password])
	return &secretData, nil
}
