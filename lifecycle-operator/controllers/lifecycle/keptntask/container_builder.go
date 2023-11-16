package keptntask

import (
	"encoding/json"
	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/taskdefinition"
	"golang.org/x/net/context"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// ContainerBuilder implements container builder interface for python
type ContainerBuilder struct {
	containerSpec *klcv1alpha3.ContainerSpec
	taskSpec      *klcv1alpha3.KeptnTaskSpec
}

func NewContainerBuilder(options BuilderOptions) *ContainerBuilder {
	return &ContainerBuilder{
		containerSpec: options.containerSpec,
	}
}

func (c *ContainerBuilder) CreateContainer(ctx context.Context) (*corev1.Container, error) {
	result := c.containerSpec.Container

	if c.taskSpec == nil {
		return result, nil
	}

	taskContext := c.taskSpec.Context

	jsonContext, err := json.Marshal(taskContext)
	if err != nil {
		return nil, err
	}

	foundKeptnContextVar := false
	for i, envVar := range result.Env {
		if envVar.Name == KeptnContextEnvVarName {
			foundKeptnContextVar = true
			result.Env[i].Value = string(jsonContext)
		}
	}

	if !foundKeptnContextVar {
		result.Env = append(result.Env, corev1.EnvVar{
			Name:  KeptnContextEnvVarName,
			Value: string(jsonContext),
		})
	}

	return result, nil
}

func (c *ContainerBuilder) CreateVolume(ctx context.Context) (*corev1.Volume, error) {
	return c.generateVolume(), nil
}

func (c *ContainerBuilder) getVolumeSource() *corev1.EmptyDirVolumeSource {
	quantity, ok := c.containerSpec.Resources.Limits["memory"]
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
	if !taskdefinition.IsVolumeMountPresent(c.containerSpec) {
		return nil
	}
	return &corev1.Volume{
		Name: c.containerSpec.VolumeMounts[0].Name,
		VolumeSource: corev1.VolumeSource{
			EmptyDir: c.getVolumeSource(),
		},
	}
}
