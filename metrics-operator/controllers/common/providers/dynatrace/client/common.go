package client

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

var ErrClientSecretInvalid = errors.New("the Dynatrace token has an invalid format")
var ErrRequestFailed = errors.New("the API returned a response with a status outside of the 2xx range")
var ErrAuthenticationFailed = errors.New("could not retrieve an OAuth token from the API")
var ErrClientUrlInvalid = errors.New("the Dynatrace authurl is not a valid url")

const (
	oAuthGrantType                  = "grant_type"
	oAuthGrantTypeClientCredentials = "client_credentials"
	oAuthScope                      = "scope"
	oAuthClientID                   = "client_id"
	oAuthClientSecret               = "client_secret"
)

const dtTokenPrefix = "dt0s08"

func validateOAuthSecret(token string, authurl string) error {
	// must start with dt0s08
	// must have 2 dots
	// third part (split by dot) must be 64 chars
	if !strings.HasPrefix(token, dtTokenPrefix) {
		return fmt.Errorf("secret does not start with required prefix %s: %w", dtTokenPrefix, ErrClientSecretInvalid)
	}
	split := strings.Split(token, ".")
	if len(split) != 3 {
		return fmt.Errorf("secret does not contain three components: %w", ErrClientSecretInvalid)
	}
	secret := split[2]
	if secretLen := len(secret); secretLen != 64 {
		return fmt.Errorf("length of secret is %d, which is not equal to 64: %w", secretLen, ErrClientSecretInvalid)
	}
	_, err := url.ParseRequestURI(authurl)
	if err != nil {
		return fmt.Errorf("authurl is not a valid url: %w", ErrClientUrlInvalid)
	}
	return nil
}

func isErrorStatus(statusCode int) bool {
	return statusCode < 200 || statusCode >= 300
}
