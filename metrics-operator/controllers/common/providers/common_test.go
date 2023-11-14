package providers

import "testing"

func TestIsProviderSupported(t *testing.T) {
	type args struct {
		providerName string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "supported provider",
			args: args{
				providerName: "dynatrace",
			},
			want: true,
		},
		{
			name: "unsupported provider",
			args: args{
				providerName: "foo",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsProviderSupported(tt.args.providerName); got != tt.want {
				t.Errorf("IsProviderSupported() = %v, want %v", got, tt.want)
			}
		})
	}
}
