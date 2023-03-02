package client

import (
	"net/http"
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

func Test_isErrorStatus(t *testing.T) {
	type args struct {
		statusCode int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "is not an error",
			args: args{
				statusCode: http.StatusOK,
			},
			want: false,
		},
		{
			name: "is an error",
			args: args{
				statusCode: http.StatusInternalServerError,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isErrorStatus(tt.args.statusCode); got != tt.want {
				t.Errorf("isErrorStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
