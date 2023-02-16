package v1alpha2

import (
	"testing"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestKeptnMetric_IsStatusSet(t *testing.T) {
	type fields struct {
		TypeMeta   v1.TypeMeta
		ObjectMeta v1.ObjectMeta
		Spec       KeptnMetricSpec
		Status     KeptnMetricStatus
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "No value set",
			fields: fields{
				Status: KeptnMetricStatus{
					Value: "",
				},
			},
			want: false,
		},
		{
			name: "we have a value",
			fields: fields{
				Status: KeptnMetricStatus{
					Value: "1.0",
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &KeptnMetric{
				TypeMeta:   tt.fields.TypeMeta,
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
				Status:     tt.fields.Status,
			}
			if got := s.IsStatusSet(); got != tt.want {
				t.Errorf("IsStatusSet() = %v, want %v", got, tt.want)
			}
		})
	}
}
