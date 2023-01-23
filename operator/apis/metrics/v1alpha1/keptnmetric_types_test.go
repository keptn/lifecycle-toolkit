package v1alpha1

import (
	"reflect"
	"testing"

	"github.com/keptn/lifecycle-toolkit/operator/apis/metrics/v1alpha1/common"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestKeptnMetric_checkAllowedProvider(t *testing.T) {
	type fields struct {
		TypeMeta   v1.TypeMeta
		ObjectMeta v1.ObjectMeta
		Spec       KeptnMetricSpec
		Status     KeptnMetricStatus
	}
	type args struct {
		provider string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &KeptnMetric{
				TypeMeta:   tt.fields.TypeMeta,
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
				Status:     tt.fields.Status,
			}
			if err := s.checkAllowedProvider(tt.args.provider); (err != nil) != tt.wantErr {
				t.Errorf("checkAllowedProvider() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateProviderName(t *testing.T) {

	tests := []struct {
		name         string
		providerName string
		fldPath      *field.Path
		want         error
	}{
		{
			name:         "bad provider",
			providerName: "keptn-metric",
			want:         common.ErrForbiddenProvider,
		},
		{
			name:         "good provider",
			providerName: "prometheus",
			want:         nil,
		},
	}
	for _, tt := range tests {

		s := &KeptnMetric{
			ObjectMeta: v1.ObjectMeta{Name: tt.name},
			Spec:       KeptnMetricSpec{Provider: ProviderRef{Name: tt.providerName}},
		}
		t.Run(tt.name, func(t *testing.T) {
			if got := s.checkAllowedProvider(tt.providerName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("validateProviderName() = %v, want %v", got, tt.want)
			}
		})
	}
}
