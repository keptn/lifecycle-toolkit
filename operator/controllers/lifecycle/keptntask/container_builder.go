package keptntask

import (
	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	"golang.org/x/net/context"
	corev1 "k8s.io/api/core/v1"
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
	if !c.taskDef.IsVolumeMountPresent() {
		return c.taskDef.Spec.Container.Container, []corev1.Volume{}, nil
	}
	return c.taskDef.Spec.Container.Container, c.taskDef.GenerateVolumes(), nil
}
