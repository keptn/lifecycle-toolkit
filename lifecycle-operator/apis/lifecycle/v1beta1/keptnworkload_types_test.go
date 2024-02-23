package v1beta1

import (
	"testing"

	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1beta1/common"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestKeptnWorkload(t *testing.T) {
	workload := &KeptnWorkload{
		ObjectMeta: v1.ObjectMeta{
			Name:      "workload",
			Namespace: "namespace",
		},
		Spec: KeptnWorkloadSpec{
			Version: "version",
			AppName: "app",
		},
	}

	require.Equal(t, []attribute.KeyValue{
		common.AppName.String("app"),
		common.WorkloadName.String("workload"),
		common.WorkloadVersion.String("version"),
	}, workload.GetSpanAttributes())

	require.Equal(t, map[string]string{
		"appName":         "app",
		"workloadName":    "workload",
		"workloadVersion": "version",
	}, workload.GetEventAnnotations())
}

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
