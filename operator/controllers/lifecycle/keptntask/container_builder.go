package keptntask

import (
	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	"golang.org/x/net/context"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// JSBuilder implements container builder interface for javascript deno
type ContainerBuilder struct {
	taskDef *klcv1alpha3.KeptnTaskDefinition
}

func newContainerBuilder(taskDef *klcv1alpha3.KeptnTaskDefinition) *ContainerBuilder {
	return &ContainerBuilder{
		taskDef: taskDef,
	}
}

func (c *ContainerBuilder) CreateContainerWithVolumes(ctx context.Context) (*corev1.Container, []corev1.Volume, error) {
	jobVolumes := []corev1.Volume{
		{
			Name: "default-volume",
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{
					// Default 50% of the memory of the node, max 1Gi
					SizeLimit: resource.NewQuantity(1, resource.Format("Gi")),
					Medium:    corev1.StorageMedium("Memory"),
				},
			},
		},
	}
	return c.taskDef.Spec.Container.Container, jobVolumes, nil
}
