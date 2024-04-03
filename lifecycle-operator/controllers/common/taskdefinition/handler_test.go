package taskdefinition

import (
	"reflect"
	"testing"

	apilifecycle "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
)

func TestGetRuntimeImage(t *testing.T) {

	t.Setenv(FunctionRuntimeImageKey, FunctionScriptKey)
	t.Setenv(PythonRuntimeImageKey, PythonScriptKey)
	tests := []struct {
		name string
		def  *apilifecycle.KeptnTaskDefinition
		want string
	}{
		{
			name: PythonScriptKey,
			def: &apilifecycle.KeptnTaskDefinition{
				Spec: apilifecycle.KeptnTaskDefinitionSpec{
					Python: &apilifecycle.RuntimeSpec{
						HttpReference: apilifecycle.HttpReference{
							Url: "testy.com",
						},
					},
				},
			},
			want: PythonScriptKey,
		}, {
			name: FunctionScriptKey,
			def: &apilifecycle.KeptnTaskDefinition{
				Spec: apilifecycle.KeptnTaskDefinitionSpec{
					Deno: &apilifecycle.RuntimeSpec{
						HttpReference: apilifecycle.HttpReference{
							Url: "testy.com",
						},
					},
				},
			},
			want: FunctionScriptKey,
		},
		{
			name: "deno and python defined, deno wins",
			def: &apilifecycle.KeptnTaskDefinition{
				Spec: apilifecycle.KeptnTaskDefinitionSpec{
					Deno: &apilifecycle.RuntimeSpec{
						HttpReference: apilifecycle.HttpReference{
							Url: "testy.com",
						},
					},
					Python: &apilifecycle.RuntimeSpec{
						HttpReference: apilifecycle.HttpReference{
							Url: "testy.com",
						},
					},
				},
			},
			want: FunctionScriptKey,
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
		def  *apilifecycle.KeptnTaskDefinition
		want *apilifecycle.RuntimeSpec
	}{
		{
			name: PythonScriptKey,
			def: &apilifecycle.KeptnTaskDefinition{
				Spec: apilifecycle.KeptnTaskDefinitionSpec{
					Python: &apilifecycle.RuntimeSpec{
						HttpReference: apilifecycle.HttpReference{
							Url: "testy.com",
						},
					},
				},
			},
			want: &apilifecycle.RuntimeSpec{
				HttpReference: apilifecycle.HttpReference{
					Url: "testy.com",
				},
			},
		},
		{
			name: FunctionScriptKey,
			def: &apilifecycle.KeptnTaskDefinition{
				Spec: apilifecycle.KeptnTaskDefinitionSpec{
					Deno: &apilifecycle.RuntimeSpec{
						HttpReference: apilifecycle.HttpReference{
							Url: "testy.com",
						},
					},
				},
			},
			want: &apilifecycle.RuntimeSpec{
				HttpReference: apilifecycle.HttpReference{
					Url: "testy.com",
				},
			},
		},
		{
			name: "deno & python exist",
			def: &apilifecycle.KeptnTaskDefinition{
				Spec: apilifecycle.KeptnTaskDefinitionSpec{
					Deno: &apilifecycle.RuntimeSpec{
						HttpReference: apilifecycle.HttpReference{
							Url: "testy.com",
						},
					},
					Python: &apilifecycle.RuntimeSpec{
						HttpReference: apilifecycle.HttpReference{
							Url: "nottesty.com",
						},
					},
				},
			},
			want: &apilifecycle.RuntimeSpec{
				HttpReference: apilifecycle.HttpReference{
					Url: "testy.com",
				},
			},
		},
		{
			name: "only container spec exists ",
			def: &apilifecycle.KeptnTaskDefinition{
				Spec: apilifecycle.KeptnTaskDefinitionSpec{
					Container: &apilifecycle.ContainerSpec{
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
		def  *apilifecycle.KeptnTaskDefinition
		want string
	}{
		{
			name: "deno",
			def: &apilifecycle.KeptnTaskDefinition{
				Spec: apilifecycle.KeptnTaskDefinitionSpec{
					Deno: &apilifecycle.RuntimeSpec{
						CmdParameters: "hi",
					},
				},
			},
			want: FunctionScriptMountPath,
		},
		{
			name: PythonScriptKey,
			def: &apilifecycle.KeptnTaskDefinition{
				Spec: apilifecycle.KeptnTaskDefinitionSpec{
					Python: &apilifecycle.RuntimeSpec{
						CmdParameters: "hi",
					},
				},
			},
			want: PythonScriptMountPath,
		},
		{
			name: "deno and python defined, deno wins",
			def: &apilifecycle.KeptnTaskDefinition{
				Spec: apilifecycle.KeptnTaskDefinitionSpec{
					Deno: &apilifecycle.RuntimeSpec{
						HttpReference: apilifecycle.HttpReference{
							Url: "testy.com",
						},
					},
					Python: &apilifecycle.RuntimeSpec{
						HttpReference: apilifecycle.HttpReference{
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
		spec *apilifecycle.RuntimeSpec
		want bool
	}{
		{
			name: "empty",
			spec: nil,
			want: true,
		},
		{
			name: "not empty",
			spec: &apilifecycle.RuntimeSpec{
				HttpReference: apilifecycle.HttpReference{Url: "hello.com"},
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
		spec *apilifecycle.ContainerSpec
		want bool
	}{
		{
			name: "empty",
			spec: nil,
			want: true,
		},
		{
			name: "not empty",
			spec: &apilifecycle.ContainerSpec{
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
		spec *apilifecycle.RuntimeSpec
		want bool
	}{
		{
			name: "empty inline",
			spec: &apilifecycle.RuntimeSpec{
				Inline: apilifecycle.Inline{},
			},
			want: false,
		},
		{
			name: "code in inline",
			spec: &apilifecycle.RuntimeSpec{
				Inline: apilifecycle.Inline{
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
		spec *apilifecycle.ContainerSpec
		want bool
	}{
		{
			name: "with mount",
			spec: &apilifecycle.ContainerSpec{
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
			spec: &apilifecycle.ContainerSpec{
				Container: nil,
			},
			want: false,
		},
		{
			name: "no mount",
			spec: &apilifecycle.ContainerSpec{
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
		definition *apilifecycle.KeptnTaskDefinition
		want       bool
	}{
		{
			name: "container spec",
			definition: &apilifecycle.KeptnTaskDefinition{
				Spec: apilifecycle.KeptnTaskDefinitionSpec{
					Container: &apilifecycle.ContainerSpec{
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
			definition: &apilifecycle.KeptnTaskDefinition{
				Spec: apilifecycle.KeptnTaskDefinitionSpec{
					Python: &apilifecycle.RuntimeSpec{
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
		spec         *apilifecycle.RuntimeSpec
		want         string
	}{
		{
			name:         "inline func",
			functionName: "funcName",
			spec: &apilifecycle.RuntimeSpec{
				Inline: apilifecycle.Inline{
					Code: "code",
				},
			},
			want: "keptnfn-funcName",
		},
		{
			name:         "inline func long name",
			functionName: "funcNamelooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooong",
			spec: &apilifecycle.RuntimeSpec{
				Inline: apilifecycle.Inline{
					Code: "code",
				},
			},
			want: "keptnfn-funcNameloooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooo",
		},
		{
			name:         "non inline func",
			functionName: "funcName",
			spec: &apilifecycle.RuntimeSpec{
				ConfigMapReference: apilifecycle.ConfigMapReference{
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
