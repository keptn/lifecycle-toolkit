package v1alpha3

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

func TestAnalysis_Validation(t *testing.T) {
	type fields struct {
		TypeMeta   v1.TypeMeta
		ObjectMeta v1.ObjectMeta
		Spec       AnalysisSpec
		Status     AnalysisStatus
	}
	tests := []struct {
		name    string
		verb    string
		fields  fields
		want    admission.Warnings
		wantErr bool
	}{
		// CREATE
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
			verb:    "create",
			want:    admission.Warnings{},
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
			verb:    "create",
			want:    admission.Warnings{},
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
			verb:    "create",
			want:    admission.Warnings{},
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
			verb:    "create",
			want:    admission.Warnings{},
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
			verb:    "create",
			want:    admission.Warnings{},
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
			verb:    "create",
			want:    admission.Warnings{},
			wantErr: true,
		},
		{
			name: "invalid Analysis with no timeframe info set",
			fields: fields{
				Spec: AnalysisSpec{
					Timeframe: Timeframe{},
				},
			},
			verb:    "create",
			want:    admission.Warnings{},
			wantErr: true,
		},
		// UPDATE
		{
			name: "valid Analysis with from/to timestamps - update",
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
			verb:    "update",
			want:    admission.Warnings{},
			wantErr: false,
		},
		{
			name: "valid Analysis with 'recent' being set' - update",
			fields: fields{
				Spec: AnalysisSpec{
					Timeframe: Timeframe{
						Recent: v1.Duration{
							Duration: 5 * time.Second,
						},
					},
				},
			},
			verb:    "update",
			want:    admission.Warnings{},
			wantErr: false,
		},
		{
			name: "invalid Analysis with from timestamp greater than to timestamps - update",
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
			verb:    "update",
			want:    admission.Warnings{},
			wantErr: true,
		},
		{
			name: "invalid Analysis with 'from' being nil - update",
			fields: fields{
				Spec: AnalysisSpec{
					Timeframe: Timeframe{
						To: v1.Time{
							Time: time.Now(),
						},
					},
				},
			},
			verb:    "update",
			want:    admission.Warnings{},
			wantErr: true,
		},
		{
			name: "invalid Analysis with 'to' being nil - update",
			fields: fields{
				Spec: AnalysisSpec{
					Timeframe: Timeframe{
						From: v1.Time{
							Time: time.Now(),
						},
					},
				},
			},
			verb:    "update",
			want:    admission.Warnings{},
			wantErr: true,
		},
		{
			name: "invalid Analysis with 'recent' ad 'from'/'to'  being set - update",
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
			verb:    "update",
			want:    admission.Warnings{},
			wantErr: true,
		},
		{
			name: "invalid Analysis with no timeframe info set - update",
			fields: fields{
				Spec: AnalysisSpec{
					Timeframe: Timeframe{},
				},
			},
			verb:    "update",
			want:    admission.Warnings{},
			wantErr: true,
		},
		// DELETE
		{
			name: "delete analysis",
			fields: fields{
				Spec: AnalysisSpec{
					Timeframe: Timeframe{},
				},
			},
			verb:    "delete",
			want:    admission.Warnings{},
			wantErr: false,
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
			var got []string
			var err error

			switch tt.verb {
			case "create":
				got, err = a.ValidateCreate()
			case "update":
				got, err = a.ValidateUpdate(&Analysis{})
			case "delete":
				got, err = a.ValidateDelete()
			default:
				got, err = a.ValidateCreate()
			}

			if !tt.wantErr {
				require.Nil(t, err)
			} else {
				require.NotNil(t, err)
			}
			require.EqualValues(t, tt.want, got)
		})
	}
}
