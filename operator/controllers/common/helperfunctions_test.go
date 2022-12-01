package common

import (
	"testing"

	klcv1alpha2 "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha2"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha2/common"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/types"
)

func Test_GetTaskStatus(t *testing.T) {
	tests := []struct {
		name     string
		inStatus []klcv1alpha2.TaskStatus
		want     klcv1alpha2.TaskStatus
	}{
		{
			name: "non-existing",
			inStatus: []klcv1alpha2.TaskStatus{
				{
					TaskDefinitionName: "def-name",
					TaskName:           "name",
					Status:             apicommon.StatePending,
				},
			},
			want: klcv1alpha2.TaskStatus{
				TaskDefinitionName: "non-existing",
				Status:             apicommon.StatePending,
				TaskName:           "",
			},
		},
		{
			name: "def-name",
			inStatus: []klcv1alpha2.TaskStatus{
				{
					TaskDefinitionName: "def-name",
					TaskName:           "name",
					Status:             apicommon.StateProgressing,
				},
			},
			want: klcv1alpha2.TaskStatus{
				TaskDefinitionName: "def-name",
				TaskName:           "name",
				Status:             apicommon.StateProgressing,
			},
		},
		{
			name:     "empty",
			inStatus: []klcv1alpha2.TaskStatus{},
			want: klcv1alpha2.TaskStatus{
				TaskDefinitionName: "empty",
				Status:             apicommon.StatePending,
				TaskName:           "",
			},
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			require.Equal(t, GetTaskStatus(tt.name, tt.inStatus), tt.want)
		})
	}
}

func Test_GetEvaluationStatus(t *testing.T) {
	tests := []struct {
		name     string
		inStatus []klcv1alpha2.EvaluationStatus
		want     klcv1alpha2.EvaluationStatus
	}{
		{
			name: "non-existing",
			inStatus: []klcv1alpha2.EvaluationStatus{
				{
					EvaluationDefinitionName: "def-name",
					EvaluationName:           "name",
					Status:                   apicommon.StatePending,
				},
			},
			want: klcv1alpha2.EvaluationStatus{
				EvaluationDefinitionName: "non-existing",
				Status:                   apicommon.StatePending,
				EvaluationName:           "",
			},
		},
		{
			name: "def-name",
			inStatus: []klcv1alpha2.EvaluationStatus{
				{
					EvaluationDefinitionName: "def-name",
					EvaluationName:           "name",
					Status:                   apicommon.StateProgressing,
				},
			},
			want: klcv1alpha2.EvaluationStatus{
				EvaluationDefinitionName: "def-name",
				EvaluationName:           "name",
				Status:                   apicommon.StateProgressing,
			},
		},
		{
			name:     "empty",
			inStatus: []klcv1alpha2.EvaluationStatus{},
			want: klcv1alpha2.EvaluationStatus{
				EvaluationDefinitionName: "empty",
				Status:                   apicommon.StatePending,
				EvaluationName:           "",
			},
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			require.Equal(t, GetEvaluationStatus(tt.name, tt.inStatus), tt.want)
		})
	}
}

func Test_GetAppVersionName(t *testing.T) {
	tests := []struct {
		namespace string
		appName   string
		version   string
		want      types.NamespacedName
	}{
		{
			namespace: "namespace",
			appName:   "name",
			version:   "version",
			want: types.NamespacedName{
				Namespace: "namespace",
				Name:      "name-version",
			},
		},
		{
			namespace: "",
			appName:   "name",
			version:   "version",
			want: types.NamespacedName{
				Namespace: "",
				Name:      "name-version",
			},
		},
		{
			namespace: "namespace",
			appName:   "",
			version:   "version",
			want: types.NamespacedName{
				Namespace: "namespace",
				Name:      "-version",
			},
		},
		{
			namespace: "namespace",
			appName:   "name",
			version:   "",
			want: types.NamespacedName{
				Namespace: "namespace",
				Name:      "name-",
			},
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			require.Equal(t, GetAppVersionName(tt.namespace, tt.appName, tt.version), tt.want)
		})
	}
}
