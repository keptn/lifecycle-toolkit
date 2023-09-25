package v1alpha3

import (
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reflect"
	"testing"
	"time"
)

func TestAnalysis_GetFromTime(t *testing.T) {
	type fields struct {
		TypeMeta   v1.TypeMeta
		ObjectMeta v1.ObjectMeta
		Spec       AnalysisSpec
		Status     AnalysisStatus
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Time
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Analysis{
				TypeMeta:   tt.fields.TypeMeta,
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
				Status:     tt.fields.Status,
			}
			if got := a.GetFrom(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFrom() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnalysis_GetToTime(t *testing.T) {
	type fields struct {
		TypeMeta   v1.TypeMeta
		ObjectMeta v1.ObjectMeta
		Spec       AnalysisSpec
		Status     AnalysisStatus
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Time
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Analysis{
				TypeMeta:   tt.fields.TypeMeta,
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
				Status:     tt.fields.Status,
			}
			if got := a.GetTo(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeframe_GetFrom(t1 *testing.T) {
	type fields struct {
		From   v1.Time
		To     v1.Time
		Recent v1.Duration
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Time
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Timeframe{
				From:   tt.fields.From,
				To:     tt.fields.To,
				Recent: tt.fields.Recent,
			}
			if got := t.GetFrom(); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("GetFrom() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeframe_GetTo(t1 *testing.T) {
	type fields struct {
		From   v1.Time
		To     v1.Time
		Recent v1.Duration
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Time
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Timeframe{
				From:   tt.fields.From,
				To:     tt.fields.To,
				Recent: tt.fields.Recent,
			}
			if got := t.GetTo(); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("GetTo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnalysis_EnsureTimeframeIsSet(t *testing.T) {
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
