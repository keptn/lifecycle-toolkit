package v1alpha3

import (
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
