package converter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsLessOrEqual(t *testing.T) {
	require.True(t, isLessOrEqual("<"))
	require.True(t, isLessOrEqual("<="))
	require.False(t, isLessOrEqual(">"))
	require.False(t, isLessOrEqual(">="))
}

func TestIsGreaterOrEqual(t *testing.T) {
	require.False(t, isGreaterOrEqual("<"))
	require.False(t, isGreaterOrEqual("<="))
	require.True(t, isGreaterOrEqual(">"))
	require.True(t, isGreaterOrEqual(">="))
}

func TestConvertResourceName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "contains invalid characters",
			args: args{
				name: "Invalid_resource",
			},
			want: "invalid-resource",
		},
		{
			name: "resource name too long",
			args: args{
				name: "abcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcd",
			},
			want: "abcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertResourceName(tt.args.name)
			require.Equal(t, tt.want, got)
		})
	}
}
