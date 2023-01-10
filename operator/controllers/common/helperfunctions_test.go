package common

import (
	"testing"

	klcv1alpha2 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2/common"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/types"
)

func Test_GetItemStatus(t *testing.T) {
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
			require.Equal(t, GetItemStatus(tt.name, tt.inStatus), tt.want)
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
