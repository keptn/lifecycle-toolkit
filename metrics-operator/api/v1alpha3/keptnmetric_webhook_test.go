package v1alpha3

import (
	"reflect"
	"testing"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestKeptnMetric_validateRangeInterval(t *testing.T) {

	tests := []struct {
		name string
		Spec KeptnMetricSpec
		want *apierrors.StatusError
	}{
		{
			name: "with-nil-range",
			Spec: KeptnMetricSpec{
				Range: nil,
			},
		},
		{
			name: "with-wrong-interval",
			Spec: KeptnMetricSpec{
				Range: &RangeSpec{Interval: "5mins"},
			},
			want: apierrors.NewInvalid(
				schema.GroupKind{Group: "metrics.keptn.sh", Kind: "KeptnMetric"},
				"with-wrong-interval",
				field.ErrorList{
					field.Invalid(
						field.NewPath("spec").Child("range").Child("interval"),
						"5mins",
						"Forbidden! The time interval cannot be parsed. Please check for suitable conventions",
					),
				},	
			),
		},
		{
			name: "with-empty-interval",
			Spec: KeptnMetricSpec{
				Range: &RangeSpec{Interval: ""},
			},
			want: apierrors.NewInvalid(
				schema.GroupKind{Group: "metrics.keptn.sh", Kind: "KeptnMetric"},
				"with-empty-interval",
				field.ErrorList{
					field.Invalid(
						field.NewPath("spec").Child("range").Child("interval"),
						"",
						"Forbidden! The time interval cannot be parsed. Please check for suitable conventions",
					),
				},	
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
			if tt.Spec.Range == nil {
				s := &KeptnMetric{
					ObjectMeta: metav1.ObjectMeta{Name: tt.name},
					Spec:       KeptnMetricSpec{Range: tt.Spec.Range},
				}
				if got := s.validateKeptnMetric(); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("validateKeptnMetric() = %v, want %v", got, tt.want)
				}
			} else {
				s := &KeptnMetric{
					ObjectMeta: metav1.ObjectMeta{Name: tt.name},
					Spec:       KeptnMetricSpec{Range: &RangeSpec{Interval: tt.Spec.Range.Interval}},
				}
				if got := s.validateKeptnMetric(); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("validateKeptnMetric() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
