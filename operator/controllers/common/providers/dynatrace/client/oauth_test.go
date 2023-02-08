package client

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_oAuthCredentials_getScopesAsString(t *testing.T) {
	oAuth := oAuthCredentials{
		scopes: []OAuthScope{OAuthScopeStorageMetricsRead, OAuthScopeEnvironmentRoleViewer},
	}

	require.Equal(t, "storage:metrics:read environment:roles:viewer", oAuth.getScopesAsString())
}

func Test_oAuthCredentials_getScopesAsStringEmptyScopes(t *testing.T) {
	oAuth := oAuthCredentials{}

	require.Equal(t, "", oAuth.getScopesAsString())
}

func Test_oAuthCredentials_urlValues(t *testing.T) {
	oAuth := oAuthCredentials{
		clientID:     "client-id",
		clientSecret: "client-secret",
		scopes:       []OAuthScope{OAuthScopeStorageMetricsRead, OAuthScopeEnvironmentRoleViewer},
	}

	urlValues := oAuth.urlValues()

	require.Equal(t, "client-id", urlValues.Get(oAuthClientID))
	require.Equal(t, "client-secret", urlValues.Get(oAuthClientSecret))
	require.Equal(t, oAuthGrantTypeClientCredentials, urlValues.Get(oAuthGrantType))
	require.Equal(t, "storage:metrics:read environment:roles:viewer", urlValues.Get(oAuthScope))
}
