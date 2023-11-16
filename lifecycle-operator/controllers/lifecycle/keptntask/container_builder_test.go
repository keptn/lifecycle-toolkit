package keptntask

import (
	"context"
	"testing"

	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestContainerBuilder_CreateContainerWithVolumes(t *testing.T) {
	tests := []struct {
		name          string
		builder       ContainerBuilder
		wantContainer *v1.Container
		wantError     bool
	}{
		{
			name: "defined, no task spec",
			builder: ContainerBuilder{
				containerSpec: v1alpha3.ContainerSpec{
					Container: &v1.Container{
						Image: "image",
					},
				},
			},
			wantContainer: &v1.Container{
				Image: "image",
				Env: []v1.EnvVar{
					{
						Name:  KeptnContextEnvVarName,
						Value: `{"workloadName":"","appName":"","appVersion":"","workloadVersion":"","taskType":"","objectType":""}`,
					},
				},
			},
		},
		{
			name: "defined, adding context",
			builder: ContainerBuilder{
				containerSpec: v1alpha3.ContainerSpec{
					Container: &v1.Container{
						Image: "image",
					},
				},
				taskSpec: v1alpha3.KeptnTaskSpec{
					Context: v1alpha3.TaskContext{
						WorkloadName: "my-workload",
					},
				},
			},
			wantContainer: &v1.Container{
				Image: "image",
				Env: []v1.EnvVar{
					{
						Name:  KeptnContextEnvVarName,
						Value: `{"workloadName":"my-workload","appName":"","appVersion":"","workloadVersion":"","taskType":"","objectType":""}`,
					},
				},
			},
		},
		{
			name: "defined, replacing context",
			builder: ContainerBuilder{
				containerSpec: v1alpha3.ContainerSpec{
					Container: &v1.Container{
						Image: "image",
						Env: []v1.EnvVar{
							{
								Name:  KeptnContextEnvVarName,
								Value: `foo`,
							},
						},
					},
				},
				taskSpec: v1alpha3.KeptnTaskSpec{
					Context: v1alpha3.TaskContext{
						WorkloadName: "my-workload",
					},
				},
			},
			wantContainer: &v1.Container{
				Image: "image",
				Env: []v1.EnvVar{
					{
						Name:  KeptnContextEnvVarName,
						Value: `{"workloadName":"my-workload","appName":"","appVersion":"","workloadVersion":"","taskType":"","objectType":""}`,
					},
				},
			},
		},
		{
			name: "nil",
			builder: ContainerBuilder{
				containerSpec: v1alpha3.ContainerSpec{
					Container: nil,
				},
			},
			wantContainer: nil,
			wantError:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container, err := tt.builder.CreateContainer(context.TODO())
			require.Equal(t, tt.wantContainer, container)
			if tt.wantError {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
			}
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
				containerSpec: v1alpha3.ContainerSpec{
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
				containerSpec: v1alpha3.ContainerSpec{
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
				containerSpec: v1alpha3.ContainerSpec{
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
		spec *v1alpha3.ContainerSpec
		want *v1.Volume
	}{
		{
			name: "defined",
			spec: &v1alpha3.ContainerSpec{
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
			spec: &v1alpha3.ContainerSpec{},
			want: nil,
		},
	}
	for _, tt := range tests {
		builder := ContainerBuilder{
			containerSpec: *tt.spec,
		}
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, builder.generateVolume())
		})
	}
}

func Test_GetVolumeSource(t *testing.T) {
	tests := []struct {
		name string
		spec *v1alpha3.ContainerSpec
		want *v1.EmptyDirVolumeSource
	}{
		{
			name: "not set limits",
			spec: &v1alpha3.ContainerSpec{
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
			spec: &v1alpha3.ContainerSpec{
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
			containerSpec: *tt.spec,
		}
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, builder.getVolumeSource())
		})
	}
}
