package client

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewAPIConfig(t *testing.T) {
	myData := SecretValues{"dt0s08.XX.XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX", "https://dev.token.internal.my-auth-url.com/sso/oauth2/token"}
	config, err := NewAPIConfig(
		"my-url",
		myData,
		WithScopes([]OAuthScope{OAuthScopeStorageMetricsRead, OAuthScopeEnvironmentRoleViewer}),
		WithAuthURL("https://dev.token.internal.my-auth-url.com/sso/oauth2/token"),
	)

	require.Nil(t, err)
	require.NotNil(t, config)

	expectedApiConfig := apiConfig{
		serverURL: "my-url",
		authURL:   "https://dev.token.internal.my-auth-url.com/sso/oauth2/token",
		oAuthCredentials: oAuthCredentials{
			clientID:     "dt0s08.XX",
			clientSecret: mockSecret,
			scopes:       []OAuthScope{OAuthScopeStorageMetricsRead, OAuthScopeEnvironmentRoleViewer},
		},
	}
	require.Equal(t, expectedApiConfig, *config)
}
