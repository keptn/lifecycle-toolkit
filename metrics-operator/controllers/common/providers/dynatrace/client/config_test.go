package client

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewAPIConfig(t *testing.T) {
	myData := secretValues{"dt0s08.XX.XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX", "my-authurl"}
	jsonData, _ := json.Marshal(myData)
	config, err := NewAPIConfig(
		"my-url",
		jsonData,
		WithScopes([]OAuthScope{OAuthScopeStorageMetricsRead, OAuthScopeEnvironmentRoleViewer}),
		WithAuthURL("my-authurl"),
	)

	require.Nil(t, err)
	require.NotNil(t, config)

	expectedApiConfig := apiConfig{
		serverURL: "my-url",
		authURL:   "my-authurl",
		oAuthCredentials: oAuthCredentials{
			clientID:     "dt0s08.XX",
			clientSecret: mockSecret,
			scopes:       []OAuthScope{OAuthScopeStorageMetricsRead, OAuthScopeEnvironmentRoleViewer},
		},
	}
	require.Equal(t, expectedApiConfig, *config)
}
