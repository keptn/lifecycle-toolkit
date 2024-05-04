package keptntask

import (
	"encoding/json"

	apilifecycle "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/taskdefinition"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// ContainerBuilder implements container builder interface for python
type ContainerBuilder struct {
	containerSpec apilifecycle.ContainerSpec
	taskSpec      apilifecycle.KeptnTaskSpec
}

func NewContainerBuilder(options BuilderOptions) *ContainerBuilder {
	builder := &ContainerBuilder{
		containerSpec: *options.containerSpec,
	}

	if options.task != nil {
		builder.taskSpec = options.task.Spec
	}

	return builder
}

func (c *ContainerBuilder) CreateContainer(ctx context.Context) (*corev1.Container, error) {
	if c.containerSpec.Container == nil {
		return nil, errors.New("no container definition provided")
	}
	result := c.containerSpec.Container

	taskContext := c.taskSpec.Context

	jsonContext, err := json.Marshal(taskContext)
	if err != nil {
		return nil, err
	}

	foundKeptnContextVar := false
	for i, envVar := range result.Env {
		if envVar.Name == KeptnContextEnvVar {
			foundKeptnContextVar = true
			result.Env[i].Value = string(jsonContext)
		}
	}

	if !foundKeptnContextVar {
		result.Env = append(result.Env, corev1.EnvVar{
			Name:  KeptnContextEnvVar,
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
	if !taskdefinition.IsVolumeMountPresent(&c.containerSpec) {
		return nil
	}
	return &corev1.Volume{
		Name: c.containerSpec.VolumeMounts[0].Name,
		VolumeSource: corev1.VolumeSource{
			EmptyDir: c.getVolumeSource(),
		},
	}
}
