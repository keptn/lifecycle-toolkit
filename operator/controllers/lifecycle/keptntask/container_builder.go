package keptntask

import (
	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/common"
	"golang.org/x/net/context"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// ContainerBuilder implements container builder interface for python
type ContainerBuilder struct {
	spec *klcv1alpha3.ContainerSpec
}

func NewContainerBuilder(spec *klcv1alpha3.ContainerSpec) *ContainerBuilder {
	return &ContainerBuilder{
		spec: spec,
	}
}

func (c *ContainerBuilder) CreateContainer(ctx context.Context) (*corev1.Container, error) {
	return c.spec.Container, nil
}

func (c *ContainerBuilder) CreateVolume(ctx context.Context) (*corev1.Volume, error) {
	return c.generateVolume(), nil
}

func (c *ContainerBuilder) getVolumeSource() *corev1.EmptyDirVolumeSource {
	quantity, ok := c.spec.Resources.Limits["memory"]
	if ok {
		return &corev1.EmptyDirVolumeSource{
			SizeLimit: &quantity,
			Medium:    corev1.StorageMedium("Memory"),
		}
	}

	return &corev1.EmptyDirVolumeSource{
		// Default 50% of the memory of the node, max 1Gi
		SizeLimit: resource.NewQuantity(1, resource.Format("Gi")),
		Medium:    corev1.StorageMedium("Memory"),
	}
}

func (c *ContainerBuilder) generateVolume() *corev1.Volume {
	if !common.IsVolumeMountPresent(c.spec) {
		return nil
	}
	return &corev1.Volume{
		Name: c.spec.VolumeMounts[0].Name,
		VolumeSource: corev1.VolumeSource{
			EmptyDir: c.getVolumeSource(),
		},
	}
}
