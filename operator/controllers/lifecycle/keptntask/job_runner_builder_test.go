package keptntask

import (
	"testing"

	"github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
)

var jsTaskDefSpec = v1alpha3.KeptnTaskDefinitionSpec{
	Function: &v1alpha3.FunctionSpec{
		Inline: v1alpha3.Inline{
			Code: "some code",
		},
	},
}

var containerTaskDefSpec = v1alpha3.KeptnTaskDefinitionSpec{
	Container: &v1alpha3.ContainerSpec{
		Container: &v1.Container{
			Image: "image",
		},
	},
}

func Test_getJobRunnerBuilder(t *testing.T) {
	jsBuilderOptions := BuilderOptions{
		taskDef: &v1alpha3.KeptnTaskDefinition{
			Spec: jsTaskDefSpec,
		},
	}
	containerBuilderOptions := BuilderOptions{
		taskDef: &v1alpha3.KeptnTaskDefinition{
			Spec: containerTaskDefSpec,
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

func Test_specExists(t *testing.T) {
	tests := []struct {
		name    string
		taskDef *v1alpha3.KeptnTaskDefinition
		want    bool
	}{
		{
			name: "js builder",
			taskDef: &v1alpha3.KeptnTaskDefinition{
				Spec: jsTaskDefSpec,
			},
			want: true,
		},
		{
			name: "container builder",
			taskDef: &v1alpha3.KeptnTaskDefinition{
				Spec: containerTaskDefSpec,
			},
			want: true,
		},
		{
			name: "empty builder",
			taskDef: &v1alpha3.KeptnTaskDefinition{
				Spec: v1alpha3.KeptnTaskDefinitionSpec{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, specExists(tt.taskDef))
		})
	}
}

func Test_isJSSpecDefined(t *testing.T) {
	tests := []struct {
		name        string
		taskDefSpec v1alpha3.KeptnTaskDefinitionSpec
		want        bool
	}{
		{
			name:        "defined",
			taskDefSpec: jsTaskDefSpec,
			want:        true,
		},
		{
			name:        "empty",
			taskDefSpec: v1alpha3.KeptnTaskDefinitionSpec{},
			want:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, isJSSpecDefined(tt.taskDefSpec))
		})
	}
}

func Test_isContainerSpecDefined(t *testing.T) {
	tests := []struct {
		name        string
		taskDefSpec v1alpha3.KeptnTaskDefinitionSpec
		want        bool
	}{
		{
			name:        "defined",
			taskDefSpec: containerTaskDefSpec,
			want:        true,
		},
		{
			name:        "empty",
			taskDefSpec: v1alpha3.KeptnTaskDefinitionSpec{},
			want:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, isContainerSpecDefined(tt.taskDefSpec))
		})
	}
}
