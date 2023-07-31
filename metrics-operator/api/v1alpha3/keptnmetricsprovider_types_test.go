package v1alpha3

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestKeptnMetricsProvider_GetType(t *testing.T) {
	type fields struct {
		TypeMeta   metav1.TypeMeta
		ObjectMeta metav1.ObjectMeta
		Spec       KeptnMetricsProviderSpec
		Status     metav1.Status
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "provider type set",
			fields: fields{
				ObjectMeta: metav1.ObjectMeta{
					Name: "provider1",
				},
				Spec: KeptnMetricsProviderSpec{
					Type:         "prometheus",
					TargetServer: "",
				},
			},
			want: "prometheus",
		},
		{
			name: "provider type not set, should return name",
			fields: fields{
				ObjectMeta: metav1.ObjectMeta{
					Name: "prometheus",
				},
				Spec: KeptnMetricsProviderSpec{
					Type:         "",
					TargetServer: "",
				},
			},
			want: "prometheus",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &KeptnMetricsProvider{
				TypeMeta:   tt.fields.TypeMeta,
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
				Status:     tt.fields.Status,
			}
			if got := p.GetType(); got != tt.want {
				t.Errorf("GetType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeptnMetricsProvider_HasSecretDefined(t *testing.T) {
	type fields struct {
		TypeMeta   metav1.TypeMeta
		ObjectMeta metav1.ObjectMeta
		Spec       KeptnMetricsProviderSpec
		Status     metav1.Status
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "secret key ref is defined and has a key defined",
			fields: fields{
				Spec: KeptnMetricsProviderSpec{
					SecretKeyRef: corev1.SecretKeySelector{
						Key: "some-secret",
					},
				},
			},
			want: false,
		},
		{
			name: "secret key ref is not defined",
			fields: fields{
				Spec: KeptnMetricsProviderSpec{},
			},
			want: false,
		},
		{
			name: "secret key ref key is empty",
			fields: fields{
				Spec: KeptnMetricsProviderSpec{
					SecretKeyRef: corev1.SecretKeySelector{
						Key: "",
					},
				},
			},
			want: false,
		},
		{
			name: "secret key ref name is empty",
			fields: fields{
				Spec: KeptnMetricsProviderSpec{
					SecretKeyRef: corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "",
						},
					},
				},
			},
			want: false,
		},
		{
			name: "secret key ref name is empty",
			fields: fields{
				Spec: KeptnMetricsProviderSpec{
					SecretKeyRef: corev1.SecretKeySelector{
						Key: "some-key",
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "some-name",
						},
					},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &KeptnMetricsProvider{
				TypeMeta:   tt.fields.TypeMeta,
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
				Status:     tt.fields.Status,
			}
			if got := p.HasSecretDefined(); got != tt.want {
				t.Errorf("HasSecretDefined() = %v, want %v", got, tt.want)
			}
		})
	}
}
