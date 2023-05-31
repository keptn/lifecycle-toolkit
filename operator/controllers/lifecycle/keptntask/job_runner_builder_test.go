package keptntask

import (
	"testing"

	"github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
)

func Test_getJobRunnerBuilder(t *testing.T) {
	jsBuilderOptions := BuilderOptions{
		taskDef: &v1alpha3.KeptnTaskDefinition{
			Spec: v1alpha3.KeptnTaskDefinitionSpec{
				Function: &v1alpha3.FunctionSpec{
					Inline: v1alpha3.Inline{
						Code: "some code",
					},
				},
			},
		},
	}
	containerBuilderOptions := BuilderOptions{
		taskDef: &v1alpha3.KeptnTaskDefinition{
			Spec: v1alpha3.KeptnTaskDefinitionSpec{
				Container: &v1alpha3.ContainerSpec{
					Container: &v1.Container{
						Image: "image",
					},
				},
			},
		},
	}
	tests := []struct {
		name    string
		options BuilderOptions
		want    JobRunnerBuilder
	}{
		{
			name:    "js builder",
			options: jsBuilderOptions,
			want:    newJSBuilder(jsBuilderOptions),
		},
		{
			name:    "container builder",
			options: containerBuilderOptions,
			want:    newContainerBuilder(containerBuilderOptions.taskDef),
		},
		{
			name: "invalid builder",
			options: BuilderOptions{
				taskDef: &v1alpha3.KeptnTaskDefinition{},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, getJobRunnerBuilder(tt.options))
		})
	}
}
