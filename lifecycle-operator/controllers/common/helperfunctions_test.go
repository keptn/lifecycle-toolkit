package common

import (
	"context"
	"testing"

	klcv1beta1 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1beta1"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1beta1/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/config"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/testcommon"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func Test_GetItemStatus(t *testing.T) {
	tests := []struct {
		name     string
		inStatus []klcv1beta1.ItemStatus
		want     klcv1beta1.ItemStatus
	}{
		{
			name: "non-existing",
			inStatus: []klcv1beta1.ItemStatus{
				{
					DefinitionName: "def-name",
					Name:           "name",
					Status:         apicommon.StatePending,
				},
			},
			want: klcv1beta1.ItemStatus{
				DefinitionName: "non-existing",
				Status:         apicommon.StatePending,
				Name:           "",
			},
		},
		{
			name: "def-name",
			inStatus: []klcv1beta1.ItemStatus{
				{
					DefinitionName: "def-name",
					Name:           "name",
					Status:         apicommon.StateProgressing,
				},
			},
			want: klcv1beta1.ItemStatus{
				DefinitionName: "def-name",
				Name:           "name",
				Status:         apicommon.StateProgressing,
			},
		},
		{
			name:     "empty",
			inStatus: []klcv1beta1.ItemStatus{},
			want: klcv1beta1.ItemStatus{
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

func Test_GetOldStatus(t *testing.T) {
	tests := []struct {
		statuses       []klcv1beta1.ItemStatus
		definitionName string
		want           apicommon.KeptnState
	}{
		{
			statuses:       []klcv1beta1.ItemStatus{},
			definitionName: "",
			want:           "",
		},
		{
			statuses:       []klcv1beta1.ItemStatus{},
			definitionName: "defName",
			want:           "",
		},
		{
			statuses: []klcv1beta1.ItemStatus{
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
			statuses: []klcv1beta1.ItemStatus{
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
			require.Equal(t, GetOldStatus(tt.definitionName, tt.statuses), tt.want)
		})
	}
}

//nolint:dupl
func Test_GetTaskDefinition(t *testing.T) {
	tests := []struct {
		name             string
		taskDef          *klcv1beta1.KeptnTaskDefinition
		taskDefName      string
		taskDefNamespace string
		out              *klcv1beta1.KeptnTaskDefinition
		wantError        bool
	}{
		{
			name: "taskDef not found",
			taskDef: &klcv1beta1.KeptnTaskDefinition{
				ObjectMeta: v1.ObjectMeta{
					Name:      "taskDef",
					Namespace: "some-other-namespace",
				},
			},
			taskDefName:      "taskDef",
			taskDefNamespace: "some-namespace",
			out:              nil,
			wantError:        true,
		},
		{
			name: "taskDef found",
			taskDef: &klcv1beta1.KeptnTaskDefinition{
				ObjectMeta: v1.ObjectMeta{
					Name:      "taskDef",
					Namespace: "some-namespace",
				},
			},
			taskDefName:      "taskDef",
			taskDefNamespace: "some-namespace",
			out: &klcv1beta1.KeptnTaskDefinition{
				ObjectMeta: v1.ObjectMeta{
					Name:      "taskDef",
					Namespace: "some-namespace",
				},
			},
			wantError: false,
		},
		{
			name: "taskDef found in default Keptn namespace",
			taskDef: &klcv1beta1.KeptnTaskDefinition{
				ObjectMeta: v1.ObjectMeta{
					Name:      "taskDef",
					Namespace: testcommon.KeptnNamespace,
				},
			},
			taskDefName:      "taskDef",
			taskDefNamespace: "some-namespace",
			out: &klcv1beta1.KeptnTaskDefinition{
				ObjectMeta: v1.ObjectMeta{
					Name:      "taskDef",
					Namespace: testcommon.KeptnNamespace,
				},
			},
			wantError: false,
		},
	}

	err := klcv1beta1.AddToScheme(scheme.Scheme)
	require.Nil(t, err)

	config.Instance().SetDefaultNamespace(testcommon.KeptnNamespace)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := fake.NewClientBuilder().WithObjects(tt.taskDef).Build()
			d, err := GetTaskDefinition(client, ctrl.Log.WithName("testytest"), context.TODO(), tt.taskDefName, tt.taskDefNamespace)
			if tt.out != nil && d != nil {
				require.Equal(t, tt.out.Name, d.Name)
				require.Equal(t, tt.out.Namespace, d.Namespace)
			} else if tt.out != d {
				t.Errorf("want: %v, got: %v", tt.out, d)
			}
			if tt.wantError != (err != nil) {
				t.Errorf("want error: %t, got: %v", tt.wantError, err)
			}

		})
	}
}

//nolint:dupl
func Test_GetEvaluationDefinition(t *testing.T) {
	tests := []struct {
		name             string
		evalDef          *klcv1beta1.KeptnEvaluationDefinition
		evalDefName      string
		evalDefNamespace string
		out              *klcv1beta1.KeptnEvaluationDefinition
		wantError        bool
	}{
		{
			name: "evalDef not found",
			evalDef: &klcv1beta1.KeptnEvaluationDefinition{
				ObjectMeta: v1.ObjectMeta{
					Name:      "evalDef",
					Namespace: "some-other-namespace",
				},
			},
			evalDefName:      "evalDef",
			evalDefNamespace: "some-namespace",
			out:              nil,
			wantError:        true,
		},
		{
			name: "evalDef found",
			evalDef: &klcv1beta1.KeptnEvaluationDefinition{
				ObjectMeta: v1.ObjectMeta{
					Name:      "evalDef",
					Namespace: "some-namespace",
				},
			},
			evalDefName:      "evalDef",
			evalDefNamespace: "some-namespace",
			out: &klcv1beta1.KeptnEvaluationDefinition{
				ObjectMeta: v1.ObjectMeta{
					Name:      "evalDef",
					Namespace: "some-namespace",
				},
			},
			wantError: false,
		},
		{
			name: "evalDef found in default Keptn namespace",
			evalDef: &klcv1beta1.KeptnEvaluationDefinition{
				ObjectMeta: v1.ObjectMeta{
					Name:      "evalDef",
					Namespace: testcommon.KeptnNamespace,
				},
			},
			evalDefName:      "evalDef",
			evalDefNamespace: "some-namespace",
			out: &klcv1beta1.KeptnEvaluationDefinition{
				ObjectMeta: v1.ObjectMeta{
					Name:      "evalDef",
					Namespace: testcommon.KeptnNamespace,
				},
			},
			wantError: false,
		},
	}

	err := klcv1beta1.AddToScheme(scheme.Scheme)
	require.Nil(t, err)
	config.Instance().SetDefaultNamespace(testcommon.KeptnNamespace)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := fake.NewClientBuilder().WithObjects(tt.evalDef).Build()
			d, err := GetEvaluationDefinition(client, ctrl.Log.WithName("testytest"), context.TODO(), tt.evalDefName, tt.evalDefNamespace)
			if tt.out != nil && d != nil {
				require.Equal(t, tt.out.Name, d.Name)
				require.Equal(t, tt.out.Namespace, d.Namespace)
			} else if tt.out != d {
				t.Errorf("want: %v, got: %v", tt.out, d)
			}
			if tt.wantError != (err != nil) {
				t.Errorf("want error: %t, got: %v", tt.wantError, err)
			}

		})
	}
}

func TestGetRequestInfo(t *testing.T) {
	req := ctrl.Request{
		NamespacedName: types.NamespacedName{
			Name:      "example",
			Namespace: "test-namespace",
		}}

	info := GetRequestInfo(req)
	expected := map[string]string{
		"name":      "example",
		"namespace": "test-namespace",
	}
	require.Equal(t, expected, info)
}
