package v1alpha3

import (
	"testing"

	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
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

func Test_GenerateVolumes(t *testing.T) {
	tests := []struct {
		name    string
		taskDef *KeptnTaskDefinition
		want    []v1.Volume
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
			want: []v1.Volume{
				{
					Name: "name",
					VolumeSource: v1.VolumeSource{
						EmptyDir: &v1.EmptyDirVolumeSource{
							SizeLimit: resource.NewQuantity(1, resource.Format("Gi")),
							Medium:    v1.StorageMedium("Memory"),
						},
					},
				},
			},
		},
		{
			name:    "empty",
			taskDef: &KeptnTaskDefinition{},
			want:    []v1.Volume{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.taskDef.GenerateVolumes())
		})
	}
}

func Test_GetVolumeSource(t *testing.T) {
	tests := []struct {
		name    string
		taskDef *KeptnTaskDefinition
		want    *v1.EmptyDirVolumeSource
	}{
		{
			name:    "not set",
			taskDef: &KeptnTaskDefinition{},
			want: &v1.EmptyDirVolumeSource{
				SizeLimit: resource.NewQuantity(1, resource.Format("Gi")),
				Medium:    v1.StorageMedium("Memory"),
			},
		},
		{
			name: "not set limits",
			taskDef: &KeptnTaskDefinition{
				Spec: KeptnTaskDefinitionSpec{
					Container: &ContainerSpec{
						Container: &v1.Container{
							Image: "image",
							Resources: v1.ResourceRequirements{
								Limits: v1.ResourceList{},
							},
						},
					},
				},
			},
			want: &v1.EmptyDirVolumeSource{
				SizeLimit: resource.NewQuantity(1, resource.Format("Gi")),
				Medium:    v1.StorageMedium("Memory"),
			},
		},
		{
			name: "set limits",
			taskDef: &KeptnTaskDefinition{
				Spec: KeptnTaskDefinitionSpec{
					Container: &ContainerSpec{
						Container: &v1.Container{
							Image: "image",
							Resources: v1.ResourceRequirements{
								Limits: v1.ResourceList{
									"memory": *resource.NewQuantity(100, resource.Format("Mi")),
								},
							},
						},
					},
				},
			},
			want: &v1.EmptyDirVolumeSource{
				SizeLimit: resource.NewQuantity(100, resource.Format("Mi")),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.taskDef.GetVolumeSource())
		})
	}
}
