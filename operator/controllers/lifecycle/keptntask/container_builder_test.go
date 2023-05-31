package keptntask

import (
	"context"
	"testing"

	"github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestContainerBuilder_CreateContainerWithVolumes(t *testing.T) {
	tests := []struct {
		name          string
		builder       ContainerBuilder
		wantContainer *v1.Container
		wantVolumes   []v1.Volume
	}{
		{
			name: "defined without volumes",
			builder: ContainerBuilder{
				taskDef: &v1alpha3.KeptnTaskDefinition{
					Spec: v1alpha3.KeptnTaskDefinitionSpec{
						Container: &v1alpha3.ContainerSpec{
							Container: &v1.Container{
								Image: "image",
							},
						},
					},
				},
			},
			wantContainer: &v1.Container{
				Image: "image",
			},
			wantVolumes: []v1.Volume{},
		},
		{
			name: "defined with volume",
			builder: ContainerBuilder{
				taskDef: &v1alpha3.KeptnTaskDefinition{
					Spec: v1alpha3.KeptnTaskDefinitionSpec{
						Container: &v1alpha3.ContainerSpec{
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
				},
			},
			wantContainer: &v1.Container{
				Image: "image",
				VolumeMounts: []v1.VolumeMount{
					{
						Name:      "test-volume",
						MountPath: "path",
					},
				},
			},
			wantVolumes: []v1.Volume{
				{
					Name: "test-volume",
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
			name: "defined with volume and limits",
			builder: ContainerBuilder{
				taskDef: &v1alpha3.KeptnTaskDefinition{
					Spec: v1alpha3.KeptnTaskDefinitionSpec{
						Container: &v1alpha3.ContainerSpec{
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
				},
			},
			wantContainer: &v1.Container{
				Image: "image",
				VolumeMounts: []v1.VolumeMount{
					{
						Name:      "test-volume",
						MountPath: "path",
					},
				},
				Resources: v1.ResourceRequirements{
					Limits: v1.ResourceList{
						"memory": *resource.NewQuantity(100, resource.Format("Mi")),
					},
				},
			},
			wantVolumes: []v1.Volume{
				{
					Name: "test-volume",
					VolumeSource: v1.VolumeSource{
						EmptyDir: &v1.EmptyDirVolumeSource{
							SizeLimit: resource.NewQuantity(100, resource.Format("Mi")),
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container, volumes, _ := tt.builder.CreateContainerWithVolumes(context.TODO())
			require.Equal(t, tt.wantContainer, container)
			require.Equal(t, tt.wantVolumes, volumes)
		})
	}
}
