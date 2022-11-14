package v1alpha1

import (
	"testing"
)

func TestKeptnAppVersion_GetWorkloadNameOfApp(t *testing.T) {
	type fields struct {
		Spec KeptnAppVersionSpec
	}
	type args struct {
		workloadName string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "",
			fields: fields{
				Spec: KeptnAppVersionSpec{AppName: "my-app"},
			},
			args: args{
				workloadName: "my-workload",
			},
			want: "my-app-my-workload",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := KeptnAppVersion{
				Spec: tt.fields.Spec,
			}
			if got := v.GetWorkloadNameOfApp(tt.args.workloadName); got != tt.want {
				t.Errorf("GetWorkloadNameOfApp() = %v, want %v", got, tt.want)
			}
		})
	}
}
