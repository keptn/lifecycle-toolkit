package v1alpha3

import (
	"testing"

	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestKeptnAppCreationRequest_IsSingleService(t *testing.T) {
	type fields struct {
		ObjectMeta v1.ObjectMeta
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "single-service application",
			fields: fields{ObjectMeta: v1.ObjectMeta{
				Annotations: map[string]string{
					common.AppTypeAnnotation: string(common.AppTypeSingleService),
				},
			}},
			want: true,
		},
		{
			name: "multi-service application",
			fields: fields{ObjectMeta: v1.ObjectMeta{
				Annotations: map[string]string{
					common.AppTypeAnnotation: string(common.AppTypeMultiService),
				},
			}},
			want: false,
		},
		{
			name: "anything else",
			fields: fields{ObjectMeta: v1.ObjectMeta{
				Annotations: map[string]string{
					common.AppTypeAnnotation: "",
				},
			}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kacr := &KeptnAppCreationRequest{
				ObjectMeta: tt.fields.ObjectMeta,
			}
			if got := kacr.IsSingleService(); got != tt.want {
				t.Errorf("IsSingleService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeptnAppCreationRequest_GetSpanAttributes(t *testing.T) {
	kacr := KeptnAppCreationRequest{
		ObjectMeta: v1.ObjectMeta{
			Name: "my-app",
		},
		Spec: KeptnAppCreationRequestSpec{},
	}

	spanAttrs := kacr.GetSpanAttributes()

	require.Equal(t, []attribute.KeyValue{
		common.AppName.String(kacr.Name),
	}, spanAttrs)
}
