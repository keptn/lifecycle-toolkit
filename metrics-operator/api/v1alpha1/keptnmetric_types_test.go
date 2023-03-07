package v1alpha1

import (
	"testing"

	"github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha2"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

func TestKeptnMetric_ConvertTo(t *testing.T) {
	now := v1.Now()
	type fields struct {
		TypeMeta   v1.TypeMeta
		ObjectMeta v1.ObjectMeta
		Spec       KeptnMetricSpec
		Status     KeptnMetricStatus
	}
	tests := []struct {
		name       string
		fields     fields
		wantErr    bool
		wantResult *v1alpha2.KeptnMetric
	}{
		{
			name: "convert to hub",
			fields: fields{
				ObjectMeta: v1.ObjectMeta{
					Name:      "my-metric",
					Namespace: "my-namespace",
				},
				Spec: KeptnMetricSpec{
					Provider: ProviderRef{
						Name: "my-provider",
					},
					Query:                "my-query",
					FetchIntervalSeconds: 10,
				},
				Status: KeptnMetricStatus{
					Value:       "10.0",
					RawValue:    []byte("10.0"),
					LastUpdated: now,
				},
			},
			wantErr: false,

			wantResult: &v1alpha2.KeptnMetric{
				ObjectMeta: v1.ObjectMeta{
					Name:      "my-metric",
					Namespace: "my-namespace",
				},
				Spec: v1alpha2.KeptnMetricSpec{
					Provider: v1alpha2.ProviderRef{
						Name: "my-provider",
					},
					Query:                "my-query",
					FetchIntervalSeconds: 10,
				},
				Status: v1alpha2.KeptnMetricStatus{
					Value:       "10.0",
					RawValue:    []byte("10.0"),
					LastUpdated: now,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := &KeptnMetric{
				TypeMeta:   tt.fields.TypeMeta,
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
				Status:     tt.fields.Status,
			}
			dst := &v1alpha2.KeptnMetric{}
			err := src.ConvertTo(dst)

			if tt.wantErr {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
				require.EqualValues(t, tt.wantResult, dst)
			}
		})
	}
}

func TestKeptnMetric_ConvertFrom(t *testing.T) {
	now := v1.Now()

	type args struct {
		srcRaw conversion.Hub
	}
	tests := []struct {
		name       string
		dst        *KeptnMetric
		args       args
		wantErr    bool
		wantResult *KeptnMetric
	}{
		{
			name: "convert from hub",
			args: args{
				srcRaw: &v1alpha2.KeptnMetric{
					ObjectMeta: v1.ObjectMeta{
						Name:      "my-metric",
						Namespace: "my-namespace",
					},
					Spec: v1alpha2.KeptnMetricSpec{
						Provider: v1alpha2.ProviderRef{
							Name: "my-provider",
						},
						Query:                "my-query",
						FetchIntervalSeconds: 10,
					},
					Status: v1alpha2.KeptnMetricStatus{
						Value:       "10.0",
						RawValue:    []byte("10.0"),
						LastUpdated: now,
					},
				},
			},
			wantErr: false,

			wantResult: &KeptnMetric{
				ObjectMeta: v1.ObjectMeta{
					Name:      "my-metric",
					Namespace: "my-namespace",
				},
				Spec: KeptnMetricSpec{
					Provider: ProviderRef{
						Name: "my-provider",
					},
					Query:                "my-query",
					FetchIntervalSeconds: 10,
				},
				Status: KeptnMetricStatus{
					Value:       "10.0",
					RawValue:    []byte("10.0"),
					LastUpdated: now,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := &KeptnMetric{}
			err := dst.ConvertFrom(tt.args.srcRaw)
			if tt.wantErr {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
				require.EqualValues(t, tt.wantResult, dst)
			}
		})
	}
}
