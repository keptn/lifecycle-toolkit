package v1alpha3

import (
	"testing"

	"github.com/stretchr/testify/require"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestKeptnMetric_validateRangeInterval(t *testing.T) {

	tests := []struct {
		name    string
		verb    string
		Spec    KeptnMetricSpec
		want    error
		oldSpec runtime.Object
	}{
		{
			name: "create-with-nil-range",
			verb: "create",
			Spec: KeptnMetricSpec{
				Range: nil,
			},
		},
		{
			name: "create-with-wrong-interval",
			verb: "create",
			Spec: KeptnMetricSpec{
				Range: &RangeSpec{Interval: "5mins"},
			},
			want: apierrors.NewInvalid(
				schema.GroupKind{Group: "metrics.keptn.sh", Kind: "KeptnMetric"},
				"create-with-wrong-interval",
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
			name: "create-with-empty-interval",
			verb: "create",
			Spec: KeptnMetricSpec{
				Range: &RangeSpec{Interval: ""},
			},
			want: apierrors.NewInvalid(
				schema.GroupKind{Group: "metrics.keptn.sh", Kind: "KeptnMetric"},
				"create-with-empty-interval",
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
			name: "create-with-right-interval",
			verb: "create",
			Spec: KeptnMetricSpec{
				Range: &RangeSpec{Interval: "5m"},
			},
		},
		{
			name: "update-with-right-interval",
			verb: "update",
			Spec: KeptnMetricSpec{
				Range: &RangeSpec{Interval: "5m"},
			},
			oldSpec: &KeptnMetric{
				Spec: KeptnMetricSpec{
					Range: &RangeSpec{Interval: "5mins"},
				},
			},
		},
		{
			name: "update-with-wrong-interval",
			verb: "update",
			Spec: KeptnMetricSpec{
				Range: &RangeSpec{Interval: "5mins"},
			},
			want: apierrors.NewInvalid(
				schema.GroupKind{Group: "metrics.keptn.sh", Kind: "KeptnMetric"},
				"update-with-wrong-interval",
				field.ErrorList{
					field.Invalid(
						field.NewPath("spec").Child("range").Child("interval"),
						"5mins",
						"Forbidden! The time interval cannot be parsed. Please check for suitable conventions",
					),
				},
			),
			oldSpec: &KeptnMetric{
				Spec: KeptnMetricSpec{
					Range: &RangeSpec{Interval: "5m"},
				},
			},
		},
		{
			name: "delete",
			verb: "delete",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s *KeptnMetric
			if tt.Spec.Range == nil {
				s = &KeptnMetric{
					ObjectMeta: metav1.ObjectMeta{Name: tt.name},
					Spec:       KeptnMetricSpec{Range: tt.Spec.Range},
				}
			} else {
				s = &KeptnMetric{
					ObjectMeta: metav1.ObjectMeta{Name: tt.name},
					Spec:       KeptnMetricSpec{Range: &RangeSpec{Interval: tt.Spec.Range.Interval}},
				}
			}
			var err error
			switch tt.verb {
			case "create":
				err = s.ValidateCreate()
			case "update":
				err = s.ValidateUpdate(tt.oldSpec)
			case "delete":
				err = s.ValidateDelete()
			}
			if tt.want == nil {
				require.Nil(t, err)
			} else {
				require.Equal(t, tt.want, err)
			}
		})
	}
}

func TestKeptnMetric_validateRangeStep(t *testing.T) {

	tests := []struct {
		name    string
		verb    string
		Spec    KeptnMetricSpec
		want    error
		oldSpec runtime.Object
	}{
		{
			name: "create-with-wrong-step",
			verb: "create",
			Spec: KeptnMetricSpec{
				Range: &RangeSpec{
					Interval: "5m",
					Step:     "1mins",
				},
			},
			want: apierrors.NewInvalid(
				schema.GroupKind{Group: "metrics.keptn.sh", Kind: "KeptnMetric"},
				"create-with-wrong-step",
				field.ErrorList{
					field.Invalid(
						field.NewPath("spec").Child("range").Child("step"),
						"1mins",
						"Forbidden! The time interval cannot be parsed. Please check for suitable conventions",
					),
				},
			),
		},
		{
			name: "create-with-empty-step",
			verb: "create",
			Spec: KeptnMetricSpec{
				Range: &RangeSpec{
					Interval: "5m",
					Step:     "",
				},
			},
			want: apierrors.NewInvalid(
				schema.GroupKind{Group: "metrics.keptn.sh", Kind: "KeptnMetric"},
				"create-with-empty-step",
				field.ErrorList{
					field.Invalid(
						field.NewPath("spec").Child("range").Child("step"),
						"",
						"Forbidden! The time interval cannot be parsed. Please check for suitable conventions",
					),
				},
			),
		},
		{
			name: "create-with-right-step",
			verb: "create",
			Spec: KeptnMetricSpec{
				Range: &RangeSpec{
					Interval: "5m",
					Step:     "1m",
				},
			},
		},
		{
			name: "create-with-wrong-interval-right-step",
			verb: "update",
			Spec: KeptnMetricSpec{
				Range: &RangeSpec{
					Interval: "5mins",
					Step:     "1m",
				},
			},
			want: apierrors.NewInvalid(
				schema.GroupKind{Group: "metrics.keptn.sh", Kind: "KeptnMetric"},
				"create-with-wrong-interval-right-step",
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
			name: "create-with-wrong-interval-wrong-step",
			verb: "update",
			Spec: KeptnMetricSpec{
				Range: &RangeSpec{
					Interval: "5mins",
					Step:     "1mins",
				},
			},
			want: apierrors.NewInvalid(
				schema.GroupKind{Group: "metrics.keptn.sh", Kind: "KeptnMetric"},
				"create-with-wrong-interval-wrong-step",
				field.ErrorList{
					field.Invalid(
						field.NewPath("spec").Child("range").Child("interval"),
						"5mins",
						"Forbidden! The time interval cannot be parsed. Please check for suitable conventions",
					),
					field.Invalid(
						field.NewPath("spec").Child("range").Child("step"),
						"1mins",
						"Forbidden! The time interval cannot be parsed. Please check for suitable conventions",
					),
				},
			),
		},
		{
			name: "update-with-right-step",
			verb: "update",
			Spec: KeptnMetricSpec{
				Range: &RangeSpec{
					Interval: "5m",
					Step:     "1m",
				},
			},
			oldSpec: &KeptnMetric{
				Spec: KeptnMetricSpec{
					Range: &RangeSpec{
						Interval: "5m",
						Step:     "1mins",
					},
				},
			},
		},
		{
			name: "update-with-wrong-step",
			verb: "update",
			Spec: KeptnMetricSpec{
				Range: &RangeSpec{
					Interval: "5m",
					Step:     "1mins",
				},
			},
			want: apierrors.NewInvalid(
				schema.GroupKind{Group: "metrics.keptn.sh", Kind: "KeptnMetric"},
				"update-with-wrong-step",
				field.ErrorList{
					field.Invalid(
						field.NewPath("spec").Child("range").Child("step"),
						"1mins",
						"Forbidden! The time interval cannot be parsed. Please check for suitable conventions",
					),
				},
			),
			oldSpec: &KeptnMetric{
				Spec: KeptnMetricSpec{
					Range: &RangeSpec{
						Interval: "5m",
						Step:     "1m",
					},
				},
			},
		},
		{
			name: "delete",
			verb: "delete",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s *KeptnMetric
			if tt.Spec.Range == nil {
				s = &KeptnMetric{
					ObjectMeta: metav1.ObjectMeta{Name: tt.name},
					Spec:       KeptnMetricSpec{Range: tt.Spec.Range},
				}
			} else {
				s = &KeptnMetric{
					ObjectMeta: metav1.ObjectMeta{Name: tt.name},
					Spec:       KeptnMetricSpec{Range: &RangeSpec{Interval: tt.Spec.Range.Interval, Step: tt.Spec.Range.Step}},
				}
			}
			var err error
			switch tt.verb {
			case "create":
				err = s.ValidateCreate()
			case "update":
				err = s.ValidateUpdate(tt.oldSpec)
			case "delete":
				err = s.ValidateDelete()
			}
			if tt.want == nil {
				require.Nil(t, err)
			} else {
				require.Equal(t, tt.want, err)
			}
		})
	}
}
