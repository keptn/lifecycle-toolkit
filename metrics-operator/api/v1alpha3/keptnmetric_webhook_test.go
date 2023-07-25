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

func TestKeptnMetric_validateAggregation(t *testing.T) {

	tests := []struct {
		name    string
		verb    string
		Spec    KeptnMetricSpec
		want    error
		oldSpec runtime.Object
	}{
		{
			name: "create-with-wrong-aggregation",
			verb: "create",
			Spec: KeptnMetricSpec{
				Range: &RangeSpec{
					Interval:    "5m",
					Step:        "1m",
					Aggregation: "p91",
				},
			},
			want: apierrors.NewInvalid(
				schema.GroupKind{Group: "metrics.keptn.sh", Kind: "KeptnMetric"},
				"create-with-wrong-aggregation",
				field.ErrorList{
					field.NotSupported(
						field.NewPath("spec").Child("range").Child("aggregation"),
						"p91",
						[]string{"p90", "p95", "p99", "max", "min", "avg", "median"},
					),
				},
			),
		},
		{
			name: "create-with-right-aggregation",
			verb: "create",
			Spec: KeptnMetricSpec{
				Range: &RangeSpec{
					Interval:    "5m",
					Step:        "1m",
					Aggregation: "p90",
				},
			},
		},
		{
			name: "create-with-empty-aggregation-with-step",
			verb: "create",
			Spec: KeptnMetricSpec{
				Range: &RangeSpec{
					Interval:    "5m",
					Step:        "1m",
					Aggregation: "",
				},
			},
			want: apierrors.NewInvalid(
				schema.GroupKind{Group: "metrics.keptn.sh", Kind: "KeptnMetric"},
				"create-with-empty-aggregation-with-step",
				field.ErrorList{
					field.Required(
						field.NewPath("spec").Child("range").Child("aggregation"),
						"Aggregation field is required if defining the step field",
					),
				},
			),
		},
		{
			name: "create-with-aggregation-with-empty-step",
			verb: "create",
			Spec: KeptnMetricSpec{
				Range: &RangeSpec{
					Interval:    "5m",
					Step:        "",
					Aggregation: "p90",
				},
			},
			want: apierrors.NewInvalid(
				schema.GroupKind{Group: "metrics.keptn.sh", Kind: "KeptnMetric"},
				"create-with-aggregation-with-empty-step",
				field.ErrorList{
					field.Required(
						field.NewPath("spec").Child("range").Child("step"),
						"Forbidden! Step interval is required for the aggregation to work",
					),
				},
			),
		},
		{
			name: "update-with-wrong-aggregation",
			verb: "update",
			Spec: KeptnMetricSpec{
				Range: &RangeSpec{
					Interval:    "5m",
					Step:        "1m",
					Aggregation: "p91",
				},
			},
			want: apierrors.NewInvalid(
				schema.GroupKind{Group: "metrics.keptn.sh", Kind: "KeptnMetric"},
				"update-with-wrong-aggregation",
				field.ErrorList{
					field.NotSupported(
						field.NewPath("spec").Child("range").Child("aggregation"),
						"p91",
						[]string{"p90", "p95", "p99", "max", "min", "avg", "median"},
					),
				},
			),
			oldSpec: &KeptnMetric{
				Spec: KeptnMetricSpec{
					Range: &RangeSpec{
						Interval:    "5m",
						Step:        "1m",
						Aggregation: "p90",
					},
				},
			},
		},
		{
			name: "update-with-right-aggregation",
			verb: "update",
			Spec: KeptnMetricSpec{
				Range: &RangeSpec{
					Interval:    "5m",
					Step:        "1m",
					Aggregation: "p90",
				},
			},
			oldSpec: &KeptnMetric{
				Spec: KeptnMetricSpec{
					Range: &RangeSpec{
						Interval:    "5m",
						Step:        "1m",
						Aggregation: "p91",
					},
				},
			},
		},
		{
			name: "update-with-empty-aggregation-with-step",
			verb: "update",
			Spec: KeptnMetricSpec{
				Range: &RangeSpec{
					Interval:    "5m",
					Step:        "1m",
				},
			},
			want: apierrors.NewInvalid(
				schema.GroupKind{Group: "metrics.keptn.sh", Kind: "KeptnMetric"},
				"update-with-empty-aggregation-with-step",
				field.ErrorList{
					field.Required(
						field.NewPath("spec").Child("range").Child("aggregation"),
						"Aggregation field is required if defining the step field",
					),
				},
			),
			oldSpec: &KeptnMetric{
				Spec: KeptnMetricSpec{
					Range: &RangeSpec{
						Interval:    "5m",
						Step:        "1m",
						Aggregation: "p90",
					},
				},
			},
		},
		{
			name: "update-with-aggregation-with-empty-step",
			verb: "update",
			Spec: KeptnMetricSpec{
				Range: &RangeSpec{
					Interval:    "5m",
					Aggregation: "p90",
				},
			},
			want: apierrors.NewInvalid(
				schema.GroupKind{Group: "metrics.keptn.sh", Kind: "KeptnMetric"},
				"update-with-aggregation-with-empty-step",
				field.ErrorList{
					field.Required(
						field.NewPath("spec").Child("range").Child("step"),
						"Forbidden! Step interval is required for the aggregation to work",
					),
				},
			),
			oldSpec: &KeptnMetric{
				Spec: KeptnMetricSpec{
					Range: &RangeSpec{
						Interval:    "5m",
						Step:        "1m",
						Aggregation: "p90",
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
					Spec:       KeptnMetricSpec{Range: &RangeSpec{Interval: tt.Spec.Range.Interval, Step: tt.Spec.Range.Step, Aggregation: tt.Spec.Range.Aggregation}},
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
