package v1alpha1

import (
	"reflect"
	"testing"

	"github.com/keptn/lifecycle-toolkit/operator/apis/metrics/v1alpha1/common"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestKeptnMetric_validateKeptnMetric(t *testing.T) {

	tests := []struct {
		name         string
		providerName string
		want         error
	}{
		{
			name:         "bad-provider",
			providerName: common.KeptnMetricProviderName,
			want:         common.ErrForbiddenProvider,
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

func TestKeptnMetric_validateProvider(t *testing.T) {

	tests := []struct {
		name string
		Spec KeptnMetricSpec
		want *field.Error
	}{
		{
			name: "bad-provider",
			Spec: KeptnMetricSpec{
				Provider: ProviderRef{Name: common.KeptnMetricProviderName},
			},
			want: field.Invalid(
				field.NewPath("spec").Child("provider").Child("name"),
				common.KeptnMetricProviderName,
				common.ErrForbiddenProvider.Error(),
			),
		},

		{
			name: "good-provider",
			Spec: KeptnMetricSpec{
				Provider: ProviderRef{Name: "prometheus"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &KeptnMetric{
				ObjectMeta: metav1.ObjectMeta{Name: tt.name},
				Spec:       KeptnMetricSpec{Provider: ProviderRef{Name: tt.Spec.Provider.Name}},
			}
			if got := s.validateProvider(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("validateProvider() = %v, want %v", got, tt.want)
			}
		})
	}
}
