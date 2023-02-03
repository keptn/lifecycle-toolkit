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
			input:  "dts08.XX.XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
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
			e := validateOAuthSecret(tt.input)
			if tt.result != nil {
				require.ErrorIs(t, e, tt.result)
			} else {
				require.Nil(t, tt.result)
			}
		})

	}
}
