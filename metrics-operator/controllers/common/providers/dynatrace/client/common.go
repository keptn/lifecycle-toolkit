package client

import (
	"github.com/pkg/errors"
)

var ErrRequestFailed = errors.New("the API returned a response with a status outside of the 2xx range")
var ErrAuthenticationFailed = errors.New("could not retrieve an OAuth token from the API")

const (
	defaultAuthURL                  = "https://dev.token.internal.dynatracelabs.com/sso/oauth2/token"
	oAuthGrantType                  = "grant_type"
	oAuthGrantTypeClientCredentials = "client_credentials"
	oAuthScope                      = "scope"
	oAuthClientID                   = "client_id"
	oAuthClientSecret               = "client_secret"
)

func isErrorStatus(statusCode int) bool {
	return statusCode < 200 || statusCode >= 300
}
