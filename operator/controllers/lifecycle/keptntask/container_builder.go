package keptntask

import (
	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	"golang.org/x/net/context"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// ContainerBuilder implements container builder interface for javascript deno
type ContainerBuilder struct {
	taskDef *klcv1alpha3.KeptnTaskDefinition
}

func newContainerBuilder(taskDef *klcv1alpha3.KeptnTaskDefinition) *ContainerBuilder {
	return &ContainerBuilder{
		taskDef: taskDef,
	}
}

func (c *ContainerBuilder) CreateContainerWithVolumes(ctx context.Context) (*corev1.Container, []corev1.Volume, error) {
	return c.taskDef.Spec.Container.Container, c.generateVolumes(), nil
}

func (c *ContainerBuilder) getVolumeSource() *corev1.EmptyDirVolumeSource {
	if c.taskDef.IsContainerSpecDefined() {
		quantity, ok := c.taskDef.Spec.Container.Resources.Limits["memory"]
		if ok {
			return &corev1.EmptyDirVolumeSource{
				SizeLimit: &quantity,
				Medium:    corev1.StorageMedium("Memory"),
			}
		}
	}

	return &corev1.EmptyDirVolumeSource{
		// Default 50% of the memory of the node, max 1Gi
		SizeLimit: resource.NewQuantity(1, resource.Format("Gi")),
		Medium:    corev1.StorageMedium("Memory"),
	}
}

func (c *ContainerBuilder) generateVolumes() []corev1.Volume {
	if !c.taskDef.IsVolumeMountPresent() {
		return []corev1.Volume{}
	}
	return []corev1.Volume{
		{
			Name: c.taskDef.Spec.Container.VolumeMounts[0].Name,
			VolumeSource: corev1.VolumeSource{
				EmptyDir: c.getVolumeSource(),
			},
		},
	}
}
