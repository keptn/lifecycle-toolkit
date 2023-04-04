package client

import "net/url"

// OAuthScope represents a scope provided for the registered OAuth client interacting with the DT API
type OAuthScope string

// These constants define the scopes that we currently need for the DQL metric functionality. This list might extend as new features will be added.
// For now, we keep this at the minimum set of scopes required, as these are currently likely to change
const (
	OAuthScopeStorageMetricsRead    OAuthScope = "storage:metrics:read"
	OAuthScopeEnvironmentRoleViewer OAuthScope = "environment:roles:viewer"
)

type oAuthCredentials struct {
	clientID     string
	clientSecret string
	scopes       []OAuthScope
	accessToken  string
}

func (oac oAuthCredentials) urlValues() url.Values {
	values := url.Values{}
	values.Add(oAuthGrantType, oAuthGrantTypeClientCredentials)
	values.Add(oAuthScope, oac.getScopesAsString())
	values.Add(oAuthClientID, oac.clientID)
	values.Add(oAuthClientSecret, oac.clientSecret)

	return values
}

func (oac oAuthCredentials) getScopesAsString() string {
	scopeStr := ""

	for i := 0; i < len(oac.scopes); i++ {
		if i == 0 {
			scopeStr += string(oac.scopes[i])
		} else {
			scopeStr += " " + string(oac.scopes[i])
		}
	}
	return scopeStr
}

type OAuthResponse struct {
	AccessToken string `json:"access_token"`
}
