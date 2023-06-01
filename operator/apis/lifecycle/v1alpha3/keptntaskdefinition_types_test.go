package v1alpha3

import (
	"testing"

	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
)

var jsTaskDef = &KeptnTaskDefinition{
	Spec: KeptnTaskDefinitionSpec{
		Function: &FunctionSpec{
			Inline: Inline{
				Code: "some code",
			},
		},
	},
}

var containerTaskDef = &KeptnTaskDefinition{
	Spec: KeptnTaskDefinitionSpec{
		Container: &ContainerSpec{
			Container: &v1.Container{
				Image: "image",
			},
		},
	},
}

func Test_SpecExists(t *testing.T) {
	tests := []struct {
		name    string
		taskDef *KeptnTaskDefinition
		want    bool
	}{
		{
			name:    "js builder",
			taskDef: jsTaskDef,
			want:    true,
		},
		{
			name:    "container builder",
			taskDef: containerTaskDef,
			want:    true,
		},
		{
			name: "empty builder",
			taskDef: &KeptnTaskDefinition{
				Spec: KeptnTaskDefinitionSpec{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.taskDef.SpecExists())
		})
	}
}

func Test_IsJSSpecDefined(t *testing.T) {
	tests := []struct {
		name    string
		taskDef *KeptnTaskDefinition
		want    bool
	}{
		{
			name:    "defined",
			taskDef: jsTaskDef,
			want:    true,
		},
		{
			name:    "empty",
			taskDef: &KeptnTaskDefinition{},
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.taskDef.IsJSSpecDefined())
		})
	}
}

func Test_IsVolumeMountPresent(t *testing.T) {
	tests := []struct {
		name    string
		taskDef *KeptnTaskDefinition
		want    bool
	}{
		{
			name: "defined",
			taskDef: &KeptnTaskDefinition{
				Spec: KeptnTaskDefinitionSpec{
					Container: &ContainerSpec{
						Container: &v1.Container{
							Image: "image",
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      "name",
									MountPath: "path",
								},
							},
						},
					},
				},
			},
			want: true,
		},
		{
			name: "empty",
			taskDef: &KeptnTaskDefinition{
				Spec: KeptnTaskDefinitionSpec{
					Container: &ContainerSpec{
						Container: &v1.Container{
							Image:        "image",
							VolumeMounts: []v1.VolumeMount{},
						},
					},
				},
			},
			want: false,
		},
		{
			name: "nil",
			taskDef: &KeptnTaskDefinition{
				Spec: KeptnTaskDefinitionSpec{
					Container: &ContainerSpec{
						Container: &v1.Container{
							Image: "image",
						},
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.taskDef.IsVolumeMountPresent())
		})
	}
}

func Test_IsContainerSpecDefined(t *testing.T) {
	tests := []struct {
		name    string
		taskDef *KeptnTaskDefinition
		want    bool
	}{
		{
			name:    "defined",
			taskDef: containerTaskDef,
			want:    true,
		},
		{
			name:    "empty",
			taskDef: &KeptnTaskDefinition{},
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.taskDef.IsContainerSpecDefined())
		})
	}
}
