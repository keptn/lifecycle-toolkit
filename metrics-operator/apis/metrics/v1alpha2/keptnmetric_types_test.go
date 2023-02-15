package v1alpha2

import (
	"reflect"
	"testing"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func Test_IsProviderValid(t *testing.T) {

	tests := []struct {
		name         string
		providerName string
		fldPath      *field.Path
		want         bool
	}{
		{
			name:         "bad provider",
			providerName: "keptn-metric",
			want:         false,
		},
		{
			name:         "good provider",
			providerName: "prometheus",
			want:         true,
		},
	}
	for _, tt := range tests {

		s := &KeptnMetric{
			ObjectMeta: v1.ObjectMeta{Name: tt.name},
			Spec:       KeptnMetricSpec{Provider: ProviderRef{Name: tt.providerName}},
		}
		t.Run(tt.name, func(t *testing.T) {
			if got := s.IsProviderValid(tt.providerName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("validateProviderName() = %v, want %v", got, tt.want)
			}
		})
	}
}
