package client

import (
	"fmt"
	"strings"
)

type apiConfig struct {
	serverURL        string
	authURL          string
	oAuthCredentials oAuthCredentials
}

type SecretValues struct {
	Token   string `json:"token"`
	AuthUrl string `json:"authurl"`
}

type APIConfigOption func(config *apiConfig)

func WithAuthURL(authURL string) APIConfigOption {
	return func(config *apiConfig) {
		config.authURL = authURL
	}
}

// WithScopes passes the given scopes to the client config
func WithScopes(scopes []OAuthScope) APIConfigOption {
	return func(config *apiConfig) {
		config.oAuthCredentials.scopes = scopes
	}
}

// NewAPIConfig returns a new apiConfig that can be used for initializing a DTAPIClient with the NewAPIClient function
func NewAPIConfig(serverURL string, secret SecretValues, opts ...APIConfigOption) (*apiConfig, error) {

	if err := validateOAuthSecret(secret.Token, secret.AuthUrl); err != nil {
		return nil, err
	}

	secretParts := strings.Split(secret.Token, ".")
	clientId := fmt.Sprintf("%s.%s", secretParts[0], secretParts[1])
	clientSecret := fmt.Sprintf("%s.%s", clientId, secretParts[2])

	cfg := &apiConfig{
		serverURL: serverURL,
		authURL:   secret.AuthUrl,
		oAuthCredentials: oAuthCredentials{
			clientID:     clientId,
			clientSecret: clientSecret,
			scopes:       []OAuthScope{OAuthScopeStorageMetricsRead, OAuthScopeEnvironmentRoleViewer},
		},
	}

	for _, o := range opts {
		o(cfg)
	}

	return cfg, nil
}
