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

const resourceNameWith253Characters = "my-sample-app-pod-name-is-a-very-long-example-name-012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890253-my-sample-app-pod-name-is-a-very-long-example-name-my-sample-app-pod-name-is-a-very-lon"

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
			name: "starts with '-'",
			args: args{
				name: "-my-resource",
			},
			want: "my-resource",
		},
		{
			name: "ends with '-'",
			args: args{
				name: "my-resource-",
			},
			want: "my-resource",
		},
		{
			name: "empty string",
			args: args{
				name: "",
			},
			want: "",
		},
		{
			name: "name is 253 characters long",
			args: args{
				name: resourceNameWith253Characters,
			},
			want: resourceNameWith253Characters,
		},
		{
			name: "name is 254 characters long",
			args: args{
				name: resourceNameWith253Characters + "x",
			},
			want: resourceNameWith253Characters,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertResourceName(tt.args.name)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestValidateResourceName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid resource name",
			args: args{
				name: "my-valid-resource-name-1",
			},
			wantErr: false,
		},
		{
			name: "invalid resource name",
			args: args{
				name: "My-invalid-resource-name",
			},
			wantErr: true,
		},
		{
			name: "invalid resource name containing '_'",
			args: args{
				name: "my_invalid-resource-name",
			},
			wantErr: true,
		},
		{
			name: "253 characters long",
			args: args{
				name: resourceNameWith253Characters,
			},
			wantErr: false,
		},
		{
			name: "254 characters long",
			args: args{
				name: resourceNameWith253Characters + "x",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateResourceName(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("ValidateResourceName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
