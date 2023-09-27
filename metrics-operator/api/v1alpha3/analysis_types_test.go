package v1alpha3

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestAnalysis_SetAndRetrieveTime(t *testing.T) {
	now := time.Now().UTC()
	before := now.Add(-5 * time.Minute)
	type fields struct {
		TypeMeta   v1.TypeMeta
		ObjectMeta v1.ObjectMeta
		Spec       AnalysisSpec
		Status     AnalysisStatus
	}
	tests := []struct {
		name     string
		fields   fields
		wantFrom time.Time
		wantTo   time.Time
	}{
		{
			name: "from and to timestamps set",
			fields: fields{
				Spec: AnalysisSpec{
					Timeframe: Timeframe{
						From: v1.Time{
							Time: before,
						},
						To: v1.Time{
							Time: now,
						},
					},
				},
			},
			wantFrom: before,
			wantTo:   now,
		},
		{
			name: "'recent' set",
			fields: fields{
				Spec: AnalysisSpec{
					Timeframe: Timeframe{
						Recent: v1.Duration{
							Duration: 3 * time.Minute,
						},
					},
				},
			},
			wantFrom: now.Add(-3 * time.Minute),
			wantTo:   now,
		},
		{
			name: "nothing set",
			fields: fields{
				Spec: AnalysisSpec{},
			},
			wantFrom: time.Time{},
			wantTo:   time.Time{},
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
			a.EnsureTimeframeIsSet()

			require.WithinDuration(t, tt.wantFrom, a.GetFrom(), 1*time.Minute)
			require.WithinDuration(t, tt.wantTo, a.GetTo(), 1*time.Minute)
		})
	}
}
