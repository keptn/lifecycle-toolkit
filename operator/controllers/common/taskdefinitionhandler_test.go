package common

import (
	"reflect"
	"testing"

	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
)

func TestGetRuntimeImage(t *testing.T) {

	t.Setenv(FunctionRuntimeImageKey, "js")
	t.Setenv(PythonRuntimeImageKey, "python")
	tests := []struct {
		name string
		def  *klcv1alpha3.KeptnTaskDefinition
		want string
	}{
		{
			name: "python",
			def: &klcv1alpha3.KeptnTaskDefinition{
				Spec: klcv1alpha3.KeptnTaskDefinitionSpec{
					Python: &klcv1alpha3.RuntimeSpec{
						HttpReference: klcv1alpha3.HttpReference{
							Url: "testy.com",
						},
					},
				},
			},
			want: "python",
		}, {
			name: "js",
			def: &klcv1alpha3.KeptnTaskDefinition{
				Spec: klcv1alpha3.KeptnTaskDefinitionSpec{
					Deno: &klcv1alpha3.RuntimeSpec{
						HttpReference: klcv1alpha3.HttpReference{
							Url: "testy.com",
						},
					},
				},
			},
			want: "js",
		}, {
			name: "default function",
			def: &klcv1alpha3.KeptnTaskDefinition{
				Spec: klcv1alpha3.KeptnTaskDefinitionSpec{
					Function: &klcv1alpha3.RuntimeSpec{
						HttpReference: klcv1alpha3.HttpReference{
							Url: "testy.com",
						},
					},
				},
			},
			want: "js",
		},
		{
			name: "default and python defined, default wins",
			def: &klcv1alpha3.KeptnTaskDefinition{
				Spec: klcv1alpha3.KeptnTaskDefinitionSpec{
					Function: &klcv1alpha3.RuntimeSpec{
						HttpReference: klcv1alpha3.HttpReference{
							Url: "testy.com",
						},
					},
					Python: &klcv1alpha3.RuntimeSpec{
						HttpReference: klcv1alpha3.HttpReference{
							Url: "testy.com",
						},
					},
				},
			},
			want: "js",
		},
		{
			name: "deno and python defined, deno wins",
			def: &klcv1alpha3.KeptnTaskDefinition{
				Spec: klcv1alpha3.KeptnTaskDefinitionSpec{
					Deno: &klcv1alpha3.RuntimeSpec{
						HttpReference: klcv1alpha3.HttpReference{
							Url: "testy.com",
						},
					},
					Python: &klcv1alpha3.RuntimeSpec{
						HttpReference: klcv1alpha3.HttpReference{
							Url: "testy.com",
						},
					},
				},
			},
			want: "js",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetRuntimeImage(tt.def); got != tt.want {
				t.Errorf("GetRuntimeImage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetRuntimeSpec(t *testing.T) {
	tests := []struct {
		name string
		def  *klcv1alpha3.KeptnTaskDefinition
		want *klcv1alpha3.RuntimeSpec
	}{
		{
			name: "python",
			def: &klcv1alpha3.KeptnTaskDefinition{
				Spec: klcv1alpha3.KeptnTaskDefinitionSpec{
					Python: &klcv1alpha3.RuntimeSpec{
						HttpReference: klcv1alpha3.HttpReference{
							Url: "testy.com",
						},
					},
				},
			},
			want: &klcv1alpha3.RuntimeSpec{
				HttpReference: klcv1alpha3.HttpReference{
					Url: "testy.com",
				},
			},
		},
		{
			name: "js",
			def: &klcv1alpha3.KeptnTaskDefinition{
				Spec: klcv1alpha3.KeptnTaskDefinitionSpec{
					Deno: &klcv1alpha3.RuntimeSpec{
						HttpReference: klcv1alpha3.HttpReference{
							Url: "testy.com",
						},
					},
				},
			},
			want: &klcv1alpha3.RuntimeSpec{
				HttpReference: klcv1alpha3.HttpReference{
					Url: "testy.com",
				},
			},
		},
		{
			name: "default function",
			def: &klcv1alpha3.KeptnTaskDefinition{
				Spec: klcv1alpha3.KeptnTaskDefinitionSpec{
					Function: &klcv1alpha3.RuntimeSpec{
						HttpReference: klcv1alpha3.HttpReference{
							Url: "testy.com",
						},
					},
				},
			},
			want: &klcv1alpha3.RuntimeSpec{
				HttpReference: klcv1alpha3.HttpReference{
					Url: "testy.com",
				},
			},
		},
		{
			name: "default function & python exist",
			def: &klcv1alpha3.KeptnTaskDefinition{
				Spec: klcv1alpha3.KeptnTaskDefinitionSpec{
					Function: &klcv1alpha3.RuntimeSpec{
						HttpReference: klcv1alpha3.HttpReference{
							Url: "testy.com",
						},
					},
					Python: &klcv1alpha3.RuntimeSpec{
						HttpReference: klcv1alpha3.HttpReference{
							Url: "nottesty.com",
						},
					},
				},
			},
			want: &klcv1alpha3.RuntimeSpec{
				HttpReference: klcv1alpha3.HttpReference{
					Url: "testy.com",
				},
			},
		},
		{
			name: "default function empty & python exists ",
			def: &klcv1alpha3.KeptnTaskDefinition{
				Spec: klcv1alpha3.KeptnTaskDefinitionSpec{
					Function: &klcv1alpha3.RuntimeSpec{
						HttpReference: klcv1alpha3.HttpReference{
							Url: "",
						},
					},
					Python: &klcv1alpha3.RuntimeSpec{
						HttpReference: klcv1alpha3.HttpReference{
							Url: "testy.com",
						},
					},
				},
			},
			want: &klcv1alpha3.RuntimeSpec{
				HttpReference: klcv1alpha3.HttpReference{
					Url: "testy.com",
				},
			},
		},
		{
			name: "only container spec exists ",
			def: &klcv1alpha3.KeptnTaskDefinition{
				Spec: klcv1alpha3.KeptnTaskDefinitionSpec{
					Container: &klcv1alpha3.ContainerSpec{
						Container: &corev1.Container{
							Name: "myc",
						},
					},
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetRuntimeSpec(tt.def); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRuntimeSpec() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetRuntimeMountPath(t *testing.T) {

	tests := []struct {
		name string
		def  *klcv1alpha3.KeptnTaskDefinition
		want string
	}{
		{
			name: "default function",
			def: &klcv1alpha3.KeptnTaskDefinition{
				Spec: klcv1alpha3.KeptnTaskDefinitionSpec{
					Function: &klcv1alpha3.RuntimeSpec{
						CmdParameters: "hi",
					},
				},
			},
			want: FunctionScriptMountPath,
		},
		{
			name: "deno",
			def: &klcv1alpha3.KeptnTaskDefinition{
				Spec: klcv1alpha3.KeptnTaskDefinitionSpec{
					Deno: &klcv1alpha3.RuntimeSpec{
						CmdParameters: "hi",
					},
				},
			},
			want: FunctionScriptMountPath,
		},
		{
			name: "python",
			def: &klcv1alpha3.KeptnTaskDefinition{
				Spec: klcv1alpha3.KeptnTaskDefinitionSpec{
					Python: &klcv1alpha3.RuntimeSpec{
						CmdParameters: "hi",
					},
				},
			},
			want: PythonScriptMountPath,
		},
		{
			name: "default and python defined, default wins",
			def: &klcv1alpha3.KeptnTaskDefinition{
				Spec: klcv1alpha3.KeptnTaskDefinitionSpec{
					Function: &klcv1alpha3.RuntimeSpec{
						HttpReference: klcv1alpha3.HttpReference{
							Url: "testy.com",
						},
					},
					Python: &klcv1alpha3.RuntimeSpec{
						HttpReference: klcv1alpha3.HttpReference{
							Url: "testy.com",
						},
					},
				},
			},
			want: FunctionScriptMountPath,
		},
		{
			name: "deno and python defined, deno wins",
			def: &klcv1alpha3.KeptnTaskDefinition{
				Spec: klcv1alpha3.KeptnTaskDefinitionSpec{
					Deno: &klcv1alpha3.RuntimeSpec{
						HttpReference: klcv1alpha3.HttpReference{
							Url: "testy.com",
						},
					},
					Python: &klcv1alpha3.RuntimeSpec{
						HttpReference: klcv1alpha3.HttpReference{
							Url: "testy.com",
						},
					},
				},
			},
			want: FunctionScriptMountPath,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetRuntimeMountPath(tt.def); got != tt.want {
				t.Errorf("GetRuntimeMountPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsRuntimeEmpty(t *testing.T) {
	tests := []struct {
		name string
		spec *klcv1alpha3.RuntimeSpec
		want bool
	}{
		{
			name: "empty",
			spec: nil,
			want: true,
		},
		{
			name: "not empty",
			spec: &klcv1alpha3.RuntimeSpec{
				HttpReference: klcv1alpha3.HttpReference{Url: "hello.com"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsRuntimeEmpty(tt.spec); got != tt.want {
				t.Errorf("IsRuntimeEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsContianerEmpty(t *testing.T) {
	tests := []struct {
		name string
		spec *klcv1alpha3.ContainerSpec
		want bool
	}{
		{
			name: "empty",
			spec: nil,
			want: true,
		},
		{
			name: "not empty",
			spec: &klcv1alpha3.ContainerSpec{
				Container: &corev1.Container{
					Name: "name",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsContainerEmpty(tt.spec); got != tt.want {
				t.Errorf("IsRuntimeEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsInline(t *testing.T) {

	var tests = []struct {
		name string
		spec *klcv1alpha3.RuntimeSpec
		want bool
	}{
		{
			name: "empty inline",
			spec: &klcv1alpha3.RuntimeSpec{
				Inline: klcv1alpha3.Inline{},
			},
			want: false,
		},
		{
			name: "code in inline",
			spec: &klcv1alpha3.RuntimeSpec{
				Inline: klcv1alpha3.Inline{
					Code: "testcode",
				},
			},
			want: true,
		},
		{
			name: "nil inline",
			spec: nil,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInline(tt.spec); got != tt.want {
				t.Errorf("IsInline() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsVolumeMountPresent(t *testing.T) {

	var tests = []struct {
		name string
		spec *klcv1alpha3.ContainerSpec
		want bool
	}{
		{
			name: "with mount",
			spec: &klcv1alpha3.ContainerSpec{
				Container: &corev1.Container{
					VolumeMounts: []corev1.VolumeMount{
						{
							Name: "myvolume",
						},
					},
				},
			},
			want: true,
		},
		{
			name: "no spec",
			spec: nil,
			want: false,
		},
		{
			name: "no container",
			spec: &klcv1alpha3.ContainerSpec{
				Container: nil,
			},
			want: false,
		},
		{
			name: "no mount",
			spec: &klcv1alpha3.ContainerSpec{
				Container: &corev1.Container{},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsVolumeMountPresent(tt.spec); got != tt.want {
				t.Errorf("IsVolumeMountPresent() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestSpecExists(t *testing.T) {

	tests := []struct {
		name       string
		definition *klcv1alpha3.KeptnTaskDefinition
		want       bool
	}{
		{
			name: "container spec",
			definition: &klcv1alpha3.KeptnTaskDefinition{
				Spec: klcv1alpha3.KeptnTaskDefinitionSpec{
					Container: &klcv1alpha3.ContainerSpec{
						Container: &corev1.Container{
							Name: "mytestcontainer",
						},
					},
				},
			},
			want: true,
		},
		{
			name: "runtime spec",
			definition: &klcv1alpha3.KeptnTaskDefinition{
				Spec: klcv1alpha3.KeptnTaskDefinitionSpec{
					Python: &klcv1alpha3.RuntimeSpec{
						CmdParameters: "ciaoPy",
					},
				},
			},
			want: true,
		},
		{
			name:       "no spec",
			definition: nil,
			want:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SpecExists(tt.definition); got != tt.want {
				t.Errorf("SpecExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCmName(t *testing.T) {

	tests := []struct {
		name         string
		functionName string
		spec         *klcv1alpha3.RuntimeSpec
		want         string
	}{
		{
			name:         "inline func",
			functionName: "funcName",
			spec: &klcv1alpha3.RuntimeSpec{
				Inline: klcv1alpha3.Inline{
					Code: "code",
				},
			},
			want: "keptnfn-funcName",
		},
		{
			name:         "inline func long name",
			functionName: "funcNamelooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooong",
			spec: &klcv1alpha3.RuntimeSpec{
				Inline: klcv1alpha3.Inline{
					Code: "code",
				},
			},
			want: "keptnfn-funcNameloooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooo",
		},
		{
			name:         "non inline func",
			functionName: "funcName",
			spec: &klcv1alpha3.RuntimeSpec{
				ConfigMapReference: klcv1alpha3.ConfigMapReference{
					Name: "configMapName",
				},
			},
			want: "configMapName",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetCmName(tt.functionName, tt.spec)
			require.Equal(t, tt.want, got)
		})
	}
}
