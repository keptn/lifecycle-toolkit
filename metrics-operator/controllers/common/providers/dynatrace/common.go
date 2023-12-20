package dynatrace

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1beta1"
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

type SecretValues struct {
	Token   string `json:"token"`
	AuthUrl string `json:"authurl"`
}

func getDTSecret(ctx context.Context, provider metricsapi.KeptnMetricsProvider, k8sClient client.Client) ([]byte, error) {
	if !provider.HasSecretDefined() || !provider.HasSecretKeyDefined() {
		return []byte{}, ErrSecretKeyRefNotDefined
	}
	dtCredsSecret := &corev1.Secret{}
	if err := k8sClient.Get(ctx, types.NamespacedName{Name: provider.Spec.SecretKeyRef.Name, Namespace: provider.Namespace}, dtCredsSecret); err != nil {
		return []byte{}, err
	}

	token := string(dtCredsSecret.Data[provider.Spec.SecretKeyRef.Key])
	secretValues := SecretValues{
		Token:   token,
		AuthUrl: fmt.Sprintf("https://dev.%s.internal.dynatracelabs.com/sso/oauth2/%s", token, token),
	}

	// Use json.Marshal to encode the Person struct to JSON
	jsonData, err := json.Marshal(secretValues)
	if err != nil {
		fmt.Println("Error encoding to JSON:", err)
	}

	if len(token) == 0 {
		return []byte{}, fmt.Errorf("secret contains invalid key %s", provider.Spec.SecretKeyRef.Key)
	}
	// return strings.Trim(string(encoder.Bytes()), "\n"), nil
	return jsonData, nil
}

func urlEncodeQuery(query string) string {
	params := strings.Split(query, "&")

	result := ""

	for i, param := range params {
		keyAndValue := strings.Split(param, "=")
		if len(keyAndValue) == 2 {
			encodedKeyAndValue := keyAndValue[0] + "=" + url.QueryEscape(keyAndValue[1])
			if i > 0 {
				result += "&"
			}
			result += encodedKeyAndValue
		}
	}

	return result
}
