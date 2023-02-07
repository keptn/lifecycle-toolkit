package client

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_validateOAuthSecret(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		result error
	}{
		{
			name:   "good token",
			input:  "dt0s08.XX.XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
			result: nil,
		},
		{
			name:   "wrong prefix",
			input:  "",
			result: ErrClientSecretInvalid,
		},
		{
			name:   "wrong format",
			input:  "",
			result: ErrClientSecretInvalid,
		},
		{
			name:   "wrong secret part",
			input:  "",
			result: ErrClientSecretInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateOAuthSecret(tt.input)

			require.ErrorIs(t, err, tt.result)
		})

	}
}
