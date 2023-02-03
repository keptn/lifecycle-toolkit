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
			input:  "",
			result: nil,
		},
		{
			name:   "wrong prefix",
			input:  "",
			result: nil,
		},
		{
			name:   "wrong format",
			input:  "",
			result: nil,
		},
		{
			name:   "wrong secret part",
			input:  "",
			result: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := validateOAuthSecret(tt.input)
			require.Equal(t, tt.result, e)
		})

	}
}
