package client

import (
	"net/http"
	"testing"
)

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
