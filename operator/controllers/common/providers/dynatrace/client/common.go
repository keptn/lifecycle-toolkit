package client

import (
	"strings"

	"github.com/pkg/errors"
)

var ErrOAuthTokenInvalid = errors.New("the Dynatrace token has an invalid format")

const dtTokenPrefix = "dts08"

func validateOAuthSecret(token string) error {
	// must start with dt0s08
	// must have 2 dots
	// third part (split by dot) must be 64 chars
	if !strings.HasPrefix(token, dtTokenPrefix) {
		return ErrOAuthTokenInvalid
	}
	split := strings.Split(token, ".")
	if len(split) != 3 {
		return ErrOAuthTokenInvalid
	}
	secret := split[2]
	if len(secret) != 64 {
		return ErrOAuthTokenInvalid
	}
	return nil
}
