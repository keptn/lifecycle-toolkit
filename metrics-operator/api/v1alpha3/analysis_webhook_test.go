package v1alpha3

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"testing"
	"time"
)

func TestAnalysis_ValidateCreate(t *testing.T) {
	type fields struct {
		TypeMeta   v1.TypeMeta
		ObjectMeta v1.ObjectMeta
		Spec       AnalysisSpec
		Status     AnalysisStatus
	}
	tests := []struct {
		name    string
		fields  fields
		want    admission.Warnings
		wantErr bool
	}{
		{
			name: "valid Analysis with from/to timestamps",
			fields: fields{
				Spec: AnalysisSpec{
					Timeframe: Timeframe{
						From: v1.Time{
							Time: time.Now(),
						},
						To: v1.Time{
							Time: time.Now().Add(1 * time.Second),
						},
					},
				},
			},
			want:    []string{},
			wantErr: false,
		},
		{
			name: "valid Analysis with 'recent' being set'",
			fields: fields{
				Spec: AnalysisSpec{
					Timeframe: Timeframe{
						Recent: v1.Duration{
							Duration: 5 * time.Second,
						},
					},
				},
			},
			want:    []string{},
			wantErr: false,
		},
		{
			name: "invalid Analysis with from timestamp greater than to timestamps",
			fields: fields{
				Spec: AnalysisSpec{
					Timeframe: Timeframe{
						From: v1.Time{
							Time: time.Now().Add(1 * time.Second),
						},
						To: v1.Time{
							Time: time.Now(),
						},
					},
				},
			},
			want:    []string{},
			wantErr: true,
		},
		{
			name: "invalid Analysis with 'from' being nil",
			fields: fields{
				Spec: AnalysisSpec{
					Timeframe: Timeframe{
						To: v1.Time{
							Time: time.Now(),
						},
					},
				},
			},
			want:    []string{},
			wantErr: true,
		},
		{
			name: "invalid Analysis with 'to' being nil",
			fields: fields{
				Spec: AnalysisSpec{
					Timeframe: Timeframe{
						From: v1.Time{
							Time: time.Now(),
						},
					},
				},
			},
			want:    []string{},
			wantErr: true,
		},
		{
			name: "invalid Analysis with 'recent' ad 'from'/'to'  being set",
			fields: fields{
				Spec: AnalysisSpec{
					Timeframe: Timeframe{
						From: v1.Time{
							Time: time.Now(),
						},
						Recent: v1.Duration{
							Duration: 1 * time.Second,
						},
					},
				},
			},
			want:    []string{},
			wantErr: true,
		},
		{
			name: "invalid Analysis with no timeframe info set",
			fields: fields{
				Spec: AnalysisSpec{
					Timeframe: Timeframe{},
				},
			},
			want:    []string{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Analysis{
				TypeMeta:   tt.fields.TypeMeta,
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
				Status:     tt.fields.Status,
			}
			got, err := a.ValidateCreate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCreate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateCreate() got = %v, want %v", got, tt.want)
			}
		})
	}
}
