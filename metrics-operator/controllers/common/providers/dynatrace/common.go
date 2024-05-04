package dynatrace

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var ErrSecretKeyRefNotDefined = errors.New("the SecretKeyRef property with the Dynatrace token is missing")
var ErrInvalidResult = errors.New("the answer does not contain any data")
var ErrDQLQueryTimeout = errors.New("timed out waiting for result of DQL query")
var ErrInvalidAuthURL = errors.New("the Dynatrace auth URL is not a valid URL")
var ErrInvalidToken = errors.New("the Dynatrace token has an invalid format")

const (
	ErrAPIMsg     = "provider api response: %s"
	dtTokenPrefix = "dt0s08"
)

type Error struct {
	Code    int    `json:"-"` // optional
	Message string `json:"message"`
}

type DQLSecret struct {
	Token   string `json:"token"`
	AuthUrl string `json:"authUrl"`
}

func (s DQLSecret) validate() error {
	// must start with dt0s08
	// must have 2 dots
	// third part (split by dot) must be 64 chars
	if !strings.HasPrefix(s.Token, dtTokenPrefix) {
		return fmt.Errorf("secret does not start with required prefix %s: %w", dtTokenPrefix, ErrInvalidToken)
	}
	split := strings.Split(s.Token, ".")
	if len(split) != 3 {
		return fmt.Errorf("secret does not contain three components: %w", ErrInvalidToken)
	}
	secret := split[2]
	if secretLen := len(secret); secretLen != 64 {
		return fmt.Errorf("length of secret is %d, which is not equal to 64: %w", secretLen, ErrInvalidToken)
	}

	if s.AuthUrl != "" {
		if _, err := url.ParseRequestURI(s.AuthUrl); err != nil {
			return fmt.Errorf("authurl is not a valid url: %w", ErrInvalidAuthURL)
		}
	}

	return nil
}

func getDTSecret(ctx context.Context, provider metricsapi.KeptnMetricsProvider, k8sClient client.Client) (string, error) {
	if !provider.HasSecretDefined() || !provider.HasSecretKeyDefined() {
		return "", ErrSecretKeyRefNotDefined
	}
	dtCredsSecret := &corev1.Secret{}
	if err := k8sClient.Get(ctx, types.NamespacedName{Name: provider.Spec.SecretKeyRef.Name, Namespace: provider.Namespace}, dtCredsSecret); err != nil {
		return "", err
	}

	token := string(dtCredsSecret.Data[provider.Spec.SecretKeyRef.Key])
	if len(token) == 0 {
		return "", fmt.Errorf("secret contains invalid key %s", provider.Spec.SecretKeyRef.Key)
	}
	return strings.Trim(token, "\n"), nil
}

func getDQLSecret(ctx context.Context, provider metricsapi.KeptnMetricsProvider, k8sClient client.Client) (*DQLSecret, error) {
	if !provider.HasSecretDefined() || !provider.HasSecretKeyDefined() {
		return nil, ErrSecretKeyRefNotDefined
	}
	dtCredsSecret := &corev1.Secret{}
	if err := k8sClient.Get(ctx, types.NamespacedName{Name: provider.Spec.SecretKeyRef.Name, Namespace: provider.Namespace}, dtCredsSecret); err != nil {
		return nil, err
	}

	credentialsStr := string(dtCredsSecret.Data[provider.Spec.SecretKeyRef.Key])

	credentialsObj := &DQLSecret{}

	if err := json.Unmarshal([]byte(credentialsStr), credentialsObj); err != nil {
		return nil, fmt.Errorf("could not unmarshal secret containing access credentials: %w", err)
	}

	if err := credentialsObj.validate(); err != nil {
		return nil, fmt.Errorf("secret contains invalid credentials: %w", err)
	}

	return credentialsObj, nil
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
