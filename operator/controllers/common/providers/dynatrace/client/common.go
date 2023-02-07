package client

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

var ErrClientSecretInvalid = errors.New("the Dynatrace token has an invalid format")

const dtTokenPrefix = "dt0s08"

func validateOAuthSecret(token string) error {
	// must start with dt0s08
	// must have 2 dots
	// third part (split by dot) must be 64 chars
	if !strings.HasPrefix(token, dtTokenPrefix) {
		return fmt.Errorf("secret %s does not start with required prefix: %w", dtTokenPrefix, ErrClientSecretInvalid)
	}
	split := strings.Split(token, ".")
	if len(split) != 3 {
		return fmt.Errorf("secret does not contain three components: %w", ErrClientSecretInvalid)
	}
	secret := split[2]
	if secretLen := len(secret); secretLen != 64 {
		return fmt.Errorf("length of secret is not equal to 64: %w", ErrClientSecretInvalid)
	}
	return nil
}
