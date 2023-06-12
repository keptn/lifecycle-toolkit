package keptntask

import (
	"testing"

	"github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
)

func Test_getJobRunnerBuilder(t *testing.T) {
	runtimeBuilderOptions := BuilderOptions{
		funcSpec: &v1alpha3.RuntimeSpec{
			Inline: v1alpha3.Inline{
				Code: "some code",
			},
		},
	}
	containerBuilderOptions := BuilderOptions{
		containerSpec: &v1alpha3.ContainerSpec{
			Container: &v1.Container{
				Image: "image",
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
			options: runtimeBuilderOptions,
			want:    NewRuntimeBuilder(runtimeBuilderOptions),
		},
		{
			name:    "container builder",
			options: containerBuilderOptions,
			want:    NewContainerBuilder(containerBuilderOptions),
		},
		{
			name:    "invalid builder",
			options: BuilderOptions{},
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, NewJobRunnerBuilder(tt.options))
		})
	}
}
