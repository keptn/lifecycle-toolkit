package v1alpha3

import (
	"testing"

	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestKeptnWorkload_GetNameWithoutAppPrefix(t *testing.T) {
	type fields struct {
		ObjectMeta v1.ObjectMeta
		Spec       KeptnWorkloadSpec
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "remove app prefix",
			fields: fields{
				ObjectMeta: v1.ObjectMeta{
					Name: "my-app-my-workload",
				},
				Spec: KeptnWorkloadSpec{
					AppName: "my-app",
				},
			},
			want: "my-workload",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := KeptnWorkload{
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
			}
			got := w.GetNameWithoutAppPrefix()

			require.Equal(t, tt.want, got)
		})
	}
}
