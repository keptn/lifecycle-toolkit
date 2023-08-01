package analysis

import (
	"github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestTargetEvaluator_Evaluate(t *testing.T) {

	compValue := 15.0
	type args struct {
		val float64
		t   v1alpha3.Target
	}
	tests := []struct {
		name string
		args args
		want v1alpha3.TargetResult
	}{
		{
			name: "less than",
			args: args{
				val: 10.0,
				t: v1alpha3.Target{
					LessThan: &v1alpha3.TargetValue{
						FixedValue: &compValue,
					},
				},
			},
			want: v1alpha3.TargetResult{
				Target: v1alpha3.Target{
					LessThan: &v1alpha3.TargetValue{
						FixedValue: &compValue,
					},
				},
				Violated: false,
			},
		},
		{
			name: "less than - violated",
			args: args{
				val: 16.0,
				t: v1alpha3.Target{
					LessThan: &v1alpha3.TargetValue{
						FixedValue: &compValue,
					},
				},
			},
			want: v1alpha3.TargetResult{
				Target: v1alpha3.Target{
					LessThan: &v1alpha3.TargetValue{
						FixedValue: &compValue,
					},
				},
				Violated: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			te := &TargetEvaluator{}
			got := te.Evaluate(tt.args.val, tt.args.t)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestTarget_Evaluate(t1 *testing.T) {

	compValue := 15.0

	type fields struct {
		Target v1alpha3.Target
	}
	type args struct {
		val float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   v1alpha3.TargetResult
	}{
		{
			name: "less than",
			fields: fields{
				Target: v1alpha3.Target{
					LessThan: &v1alpha3.TargetValue{
						FixedValue: &compValue,
					},
				},
			},
			args: args{
				val: 10.0,
			},
			want: v1alpha3.TargetResult{
				Target: v1alpha3.Target{
					LessThan: &v1alpha3.TargetValue{
						FixedValue: &compValue,
					},
				},
				Violated: false,
			},
		},
		{
			name: "less than - violated",
			fields: fields{
				Target: v1alpha3.Target{
					LessThan: &v1alpha3.TargetValue{
						FixedValue: &compValue,
					},
				},
			},
			args: args{
				val: 16.0,
			},
			want: v1alpha3.TargetResult{
				Target: v1alpha3.Target{
					LessThan: &v1alpha3.TargetValue{
						FixedValue: &compValue,
					},
				},
				Violated: true,
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Target{
				Target: tt.fields.Target,
			}
			if got := t.Evaluate(tt.args.val); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}
