package keptntask

import (
	"context"
	"testing"

	lifecycle "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1beta1"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestContainerBuilder_CreateContainerWithVolumes(t *testing.T) {
	tests := []struct {
		name          string
		builder       ContainerBuilder
		wantContainer *v1.Container
	}{
		{
			name: "defined",
			builder: ContainerBuilder{
				spec: &lifecycle.ContainerSpec{
					Container: &v1.Container{
						Image: "image",
					},
				},
			},
			wantContainer: &v1.Container{
				Image: "image",
			},
		},
		{
			name: "nil",
			builder: ContainerBuilder{
				spec: &lifecycle.ContainerSpec{
					Container: nil,
				},
			},
			wantContainer: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container, _ := tt.builder.CreateContainer(context.TODO())
			require.Equal(t, tt.wantContainer, container)
		})
	}
}

func TestContainerBuilder_CreateVolume(t *testing.T) {
	tests := []struct {
		name       string
		builder    ContainerBuilder
		wantVolume *v1.Volume
	}{
		{
			name: "defined without volume",
			builder: ContainerBuilder{
				spec: &lifecycle.ContainerSpec{
					Container: &v1.Container{
						Image: "image",
					},
				},
			},
			wantVolume: nil,
		},
		{
			name: "defined with volume",
			builder: ContainerBuilder{
				spec: &lifecycle.ContainerSpec{
					Container: &v1.Container{
						Image: "image",
						VolumeMounts: []v1.VolumeMount{
							{
								Name:      "test-volume",
								MountPath: "path",
							},
						},
					},
				},
			},
			wantVolume: &v1.Volume{
				Name: "test-volume",
				VolumeSource: v1.VolumeSource{
					EmptyDir: &v1.EmptyDirVolumeSource{
						SizeLimit: resource.NewQuantity(1, resource.Format("Gi")),
						Medium:    v1.StorageMedium("Memory"),
					},
				},
			},
		},
		{
			name: "defined with volume and limits",
			builder: ContainerBuilder{
				spec: &lifecycle.ContainerSpec{
					Container: &v1.Container{
						Image: "image",
						Resources: v1.ResourceRequirements{
							Limits: v1.ResourceList{
								"memory": *resource.NewQuantity(100, resource.Format("Mi")),
							},
						},
						VolumeMounts: []v1.VolumeMount{
							{
								Name:      "test-volume",
								MountPath: "path",
							},
						},
					},
				},
			},
			wantVolume: &v1.Volume{
				Name: "test-volume",
				VolumeSource: v1.VolumeSource{
					EmptyDir: &v1.EmptyDirVolumeSource{
						SizeLimit: resource.NewQuantity(100, resource.Format("Mi")),
						Medium:    v1.StorageMedium("Memory"),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			volume, _ := tt.builder.CreateVolume(context.TODO())
			require.Equal(t, tt.wantVolume, volume)
		})
	}
}

func Test_GenerateVolumes(t *testing.T) {
	tests := []struct {
		name string
		spec *lifecycle.ContainerSpec
		want *v1.Volume
	}{
		{
			name: "defined",
			spec: &lifecycle.ContainerSpec{
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
			want: &v1.Volume{
				Name: "name",
				VolumeSource: v1.VolumeSource{
					EmptyDir: &v1.EmptyDirVolumeSource{
						SizeLimit: resource.NewQuantity(1, resource.Format("Gi")),
						Medium:    v1.StorageMedium("Memory"),
					},
				},
			},
		},
		{
			name: "empty",
			spec: &lifecycle.ContainerSpec{},
			want: nil,
		},
	}
	for _, tt := range tests {
		builder := ContainerBuilder{
			spec: tt.spec,
		}
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, builder.generateVolume())
		})
	}
}

func Test_GetVolumeSource(t *testing.T) {
	tests := []struct {
		name string
		spec *lifecycle.ContainerSpec
		want *v1.EmptyDirVolumeSource
	}{
		{
			name: "not set limits",
			spec: &lifecycle.ContainerSpec{
				Container: &v1.Container{
					Image: "image",
					Resources: v1.ResourceRequirements{
						Limits: v1.ResourceList{},
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
			spec: &lifecycle.ContainerSpec{
				Container: &v1.Container{
					Image: "image",
					Resources: v1.ResourceRequirements{
						Limits: v1.ResourceList{
							"memory": *resource.NewQuantity(100, resource.Format("Mi")),
						},
					},
				},
			},
			want: &v1.EmptyDirVolumeSource{
				SizeLimit: resource.NewQuantity(100, resource.Format("Mi")),
				Medium:    v1.StorageMedium("Memory"),
			},
		},
	}
	for _, tt := range tests {
		builder := ContainerBuilder{
			spec: tt.spec,
		}
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, builder.getVolumeSource())
		})
	}
}
