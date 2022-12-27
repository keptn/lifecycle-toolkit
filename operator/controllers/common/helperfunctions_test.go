package common

import (
	"testing"

	klcv1alpha2 "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha2"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha2/common"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/types"
)

//nolint:dupl
func Test_GetTaskStatus(t *testing.T) {
	tests := []struct {
		name     string
		inStatus []klcv1alpha2.ItemStatus
		want     klcv1alpha2.ItemStatus
	}{
		{
			name: "non-existing",
			inStatus: []klcv1alpha2.ItemStatus{
				{
					DefinitionName: "def-name",
					Name:           "name",
					Status:         apicommon.StatePending,
				},
			},
			want: klcv1alpha2.ItemStatus{
				DefinitionName: "non-existing",
				Status:         apicommon.StatePending,
				Name:           "",
			},
		},
		{
			name: "def-name",
			inStatus: []klcv1alpha2.ItemStatus{
				{
					DefinitionName: "def-name",
					Name:           "name",
					Status:         apicommon.StateProgressing,
				},
			},
			want: klcv1alpha2.ItemStatus{
				DefinitionName: "def-name",
				Name:           "name",
				Status:         apicommon.StateProgressing,
			},
		},
		{
			name:     "empty",
			inStatus: []klcv1alpha2.ItemStatus{},
			want: klcv1alpha2.ItemStatus{
				DefinitionName: "empty",
				Status:         apicommon.StatePending,
				Name:           "",
			},
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			require.Equal(t, GetTaskStatus(tt.name, tt.inStatus), tt.want)
		})
	}
}

//nolint:dupl
func Test_GetEvaluationStatus(t *testing.T) {
	tests := []struct {
		name     string
		inStatus []klcv1alpha2.ItemStatus
		want     klcv1alpha2.ItemStatus
	}{
		{
			name: "non-existing",
			inStatus: []klcv1alpha2.ItemStatus{
				{
					DefinitionName: "def-name",
					Name:           "name",
					Status:         apicommon.StatePending,
				},
			},
			want: klcv1alpha2.ItemStatus{
				DefinitionName: "non-existing",
				Status:         apicommon.StatePending,
				Name:           "",
			},
		},
		{
			name: "def-name",
			inStatus: []klcv1alpha2.ItemStatus{
				{
					DefinitionName: "def-name",
					Name:           "name",
					Status:         apicommon.StateProgressing,
				},
			},
			want: klcv1alpha2.ItemStatus{
				DefinitionName: "def-name",
				Name:           "name",
				Status:         apicommon.StateProgressing,
			},
		},
		{
			name:     "empty",
			inStatus: []klcv1alpha2.ItemStatus{},
			want: klcv1alpha2.ItemStatus{
				DefinitionName: "empty",
				Status:         apicommon.StatePending,
				Name:           "",
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

func Test_GetOldStatus(t *testing.T) {
	tests := []struct {
		statuses       []klcv1alpha2.ItemStatus
		definitionName string
		want           apicommon.KeptnState
	}{
		{
			statuses:       []klcv1alpha2.ItemStatus{},
			definitionName: "",
			want:           "",
		},
		{
			statuses:       []klcv1alpha2.ItemStatus{},
			definitionName: "defName",
			want:           "",
		},
		{
			statuses: []klcv1alpha2.ItemStatus{
				{
					DefinitionName: "defName",
					Status:         apicommon.StateFailed,
					Name:           "name",
				},
			},
			definitionName: "defNameNon",
			want:           "",
		},
		{
			statuses: []klcv1alpha2.ItemStatus{
				{
					DefinitionName: "defName",
					Status:         apicommon.StateFailed,
					Name:           "name",
				},
			},
			definitionName: "defName",
			want:           apicommon.StateFailed,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			require.Equal(t, GetOldStatus(tt.statuses, tt.definitionName), tt.want)
		})
	}
}
