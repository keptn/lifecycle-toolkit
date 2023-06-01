package v1alpha3

import (
	"reflect"
	"testing"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestKeptnMetric_validateRangeInterval(t *testing.T) {

	tests := []struct {
		name string
		Spec KeptnMetricSpec
		want *field.Error
	}{
		{
			name: "with-wrong-interval",
			Spec: KeptnMetricSpec{
				Range: &RangeSpec{Interval: "5mins"},
			},
			want: field.Invalid(
				field.NewPath("spec").Child("range").Child("interval"),
				"5mins",
				errors.New("Forbidden! The time interval cannot be parsed. Please check for suitable conventions").Error(),
			),
		},
		{
			name: "with-empty-interval",
			Spec: KeptnMetricSpec{
				Range: &RangeSpec{Interval: ""},
			},
			want: field.Invalid(
				field.NewPath("spec").Child("range").Child("interval"),
				"",
				errors.New("Forbidden! The time interval cannot be parsed. Please check for suitable conventions").Error(),
			),
		},
		{
			name: "with-right-interval",
			Spec: KeptnMetricSpec{
				Range: &RangeSpec{Interval: "5m"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &KeptnMetric{
				ObjectMeta: metav1.ObjectMeta{Name: tt.name},
				Spec:       KeptnMetricSpec{Range: &RangeSpec{Interval: tt.Spec.Range.Interval}},
			}
			if got := s.validateRangeInterval(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("validateRangeInterval() = %v, want %v", got, tt.want)
			}
		})
	}
}
