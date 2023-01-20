package v1alpha1

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func Test_validateProviderName(t *testing.T) {

	tests := []struct {
		name         string
		providerName string
		fldPath      *field.Path
		want         *field.Error
	}{
		{
			name:         "bad provider",
			providerName: "keptn-metric",
			fldPath:      nil,
			want:         field.Invalid(nil, "keptn-metric", ErrForbiddenProvider.Error()),
		},
		{
			name:         "good provider",
			providerName: "prometheus",
			fldPath:      nil,
			want:         nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateProviderName(tt.providerName, tt.fldPath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("validateProviderName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeptnMetric_validateKeptnMetric(t *testing.T) {

	tests := []struct {
		name         string
		providerName string
		want         error
	}{
		{
			name:         "bad-provider",
			providerName: KeptnMetricProviderName,
			want:         ErrForbiddenProvider,
		},

		{
			name:         "good-provider",
			providerName: "prometheus",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &KeptnMetric{
				ObjectMeta: metav1.ObjectMeta{Name: tt.name},
				Spec:       KeptnMetricSpec{Provider: ProviderRef{Name: tt.providerName}},
			}
			got := r.validateKeptnMetric()
			if tt.want != nil {
				require.NotNil(t, got)
				require.Contains(t, got.Error(), tt.want.Error())
			} else {
				require.Nil(t, got)
			}
		})
	}
}
