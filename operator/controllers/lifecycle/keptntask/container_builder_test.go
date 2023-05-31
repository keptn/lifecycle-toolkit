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
	wantVolumes := []v1.Volume{
		{
			Name: "default-volume",
			VolumeSource: v1.VolumeSource{
				EmptyDir: &v1.EmptyDirVolumeSource{
					SizeLimit: resource.NewQuantity(1, resource.Format("Gi")),
					Medium:    v1.StorageMedium("Memory"),
				},
			},
		},
	}
	tests := []struct {
		name          string
		builder       ContainerBuilder
		wantContainer *v1.Container
	}{
		{
			name: "defined",
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
		},
		{
			name: "empty",
			builder: ContainerBuilder{
				taskDef: &v1alpha3.KeptnTaskDefinition{
					Spec: v1alpha3.KeptnTaskDefinitionSpec{
						Container: &v1alpha3.ContainerSpec{},
					},
				},
			},
			wantContainer: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container, volumes, _ := tt.builder.CreateContainerWithVolumes(context.TODO())
			require.Equal(t, tt.wantContainer, container)
			require.Equal(t, wantVolumes, volumes)
		})
	}
}
