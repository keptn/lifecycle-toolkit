package client

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewAPIConfig(t *testing.T) {
	config, err := NewAPIConfig(
		"my-url",
		mockSecret,
		WithScopes([]OAuthScope{OAuthScopeStorageMetricsRead, OAuthScopeEnvironmentRoleViewer}),
		WithAuthURL("my-url/auth"),
	)

	require.Nil(t, err)
	require.NotNil(t, config)

	expectedApiConfig := apiConfig{
		serverURL: "my-url",
		authURL:   "my-url/auth",
		oAuthCredentials: oAuthCredentials{
			clientID:     "dt0s08.XX",
			clientSecret: mockSecret,
			scopes:       []OAuthScope{OAuthScopeStorageMetricsRead, OAuthScopeEnvironmentRoleViewer},
		},
	}
	require.Equal(t, expectedApiConfig, *config)
}
