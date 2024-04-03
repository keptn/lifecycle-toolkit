package common

import (
	"context"
	"reflect"
	"testing"

	apilifecycle "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/config"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/testcommon"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func Test_GetItemStatus(t *testing.T) {
	tests := []struct {
		name     string
		inStatus []apilifecycle.ItemStatus
		want     apilifecycle.ItemStatus
	}{
		{
			name: "non-existing",
			inStatus: []apilifecycle.ItemStatus{
				{
					DefinitionName: "def-name",
					Name:           "name",
					Status:         apicommon.StatePending,
				},
			},
			want: apilifecycle.ItemStatus{
				DefinitionName: "non-existing",
				Status:         apicommon.StatePending,
				Name:           "",
			},
		},
		{
			name: "def-name",
			inStatus: []apilifecycle.ItemStatus{
				{
					DefinitionName: "def-name",
					Name:           "name",
					Status:         apicommon.StateProgressing,
				},
			},
			want: apilifecycle.ItemStatus{
				DefinitionName: "def-name",
				Name:           "name",
				Status:         apicommon.StateProgressing,
			},
		},
		{
			name:     "empty",
			inStatus: []apilifecycle.ItemStatus{},
			want: apilifecycle.ItemStatus{
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
		statuses       []apilifecycle.ItemStatus
		definitionName string
		want           apicommon.KeptnState
	}{
		{
			statuses:       []apilifecycle.ItemStatus{},
			definitionName: "",
			want:           "",
		},
		{
			statuses:       []apilifecycle.ItemStatus{},
			definitionName: "defName",
			want:           "",
		},
		{
			statuses: []apilifecycle.ItemStatus{
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
			statuses: []apilifecycle.ItemStatus{
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
		taskDef          *apilifecycle.KeptnTaskDefinition
		taskDefName      string
		taskDefNamespace string
		out              *apilifecycle.KeptnTaskDefinition
		wantError        bool
	}{
		{
			name: "taskDef not found",
			taskDef: &apilifecycle.KeptnTaskDefinition{
				ObjectMeta: metav1.ObjectMeta{
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
			taskDef: &apilifecycle.KeptnTaskDefinition{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "taskDef",
					Namespace: "some-namespace",
				},
			},
			taskDefName:      "taskDef",
			taskDefNamespace: "some-namespace",
			out: &apilifecycle.KeptnTaskDefinition{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "taskDef",
					Namespace: "some-namespace",
				},
			},
			wantError: false,
		},
		{
			name: "taskDef found in default Keptn namespace",
			taskDef: &apilifecycle.KeptnTaskDefinition{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "taskDef",
					Namespace: testcommon.KeptnNamespace,
				},
			},
			taskDefName:      "taskDef",
			taskDefNamespace: "some-namespace",
			out: &apilifecycle.KeptnTaskDefinition{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "taskDef",
					Namespace: testcommon.KeptnNamespace,
				},
			},
			wantError: false,
		},
	}

	err := apilifecycle.AddToScheme(scheme.Scheme)
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
		evalDef          *apilifecycle.KeptnEvaluationDefinition
		evalDefName      string
		evalDefNamespace string
		out              *apilifecycle.KeptnEvaluationDefinition
		wantError        bool
	}{
		{
			name: "evalDef not found",
			evalDef: &apilifecycle.KeptnEvaluationDefinition{
				ObjectMeta: metav1.ObjectMeta{
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
			evalDef: &apilifecycle.KeptnEvaluationDefinition{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "evalDef",
					Namespace: "some-namespace",
				},
			},
			evalDefName:      "evalDef",
			evalDefNamespace: "some-namespace",
			out: &apilifecycle.KeptnEvaluationDefinition{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "evalDef",
					Namespace: "some-namespace",
				},
			},
			wantError: false,
		},
		{
			name: "evalDef found in default Keptn namespace",
			evalDef: &apilifecycle.KeptnEvaluationDefinition{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "evalDef",
					Namespace: testcommon.KeptnNamespace,
				},
			},
			evalDefName:      "evalDef",
			evalDefNamespace: "some-namespace",
			out: &apilifecycle.KeptnEvaluationDefinition{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "evalDef",
					Namespace: testcommon.KeptnNamespace,
				},
			},
			wantError: false,
		},
	}

	err := apilifecycle.AddToScheme(scheme.Scheme)
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

func Test_MergeMaps(t *testing.T) {
	tests := []struct {
		name string
		map1 map[string]string
		map2 map[string]string
		want map[string]string
	}{
		{
			name: "two empty maps",
			map1: map[string]string{},
			map2: map[string]string{},
			want: map[string]string{},
		},
		{
			name: "second map empty",
			map1: map[string]string{
				"test1": "testy1",
			},
			map2: map[string]string{},
			want: map[string]string{
				"test1": "testy1",
			},
		},
		{
			name: "first map empty",
			map1: map[string]string{},
			map2: map[string]string{
				"test1": "testy1",
			},
			want: map[string]string{
				"test1": "testy1",
			},
		},
		{
			name: "maps do not overlay",
			map1: map[string]string{
				"test2": "testy2",
			},
			map2: map[string]string{
				"test1": "testy1",
			},
			want: map[string]string{
				"test1": "testy1",
				"test2": "testy2",
			},
		},
		{
			name: "maps overlay - map2 wins",
			map1: map[string]string{
				"test2": "testy2",
				"test3": "testy4",
			},
			map2: map[string]string{
				"test1": "testy1",
				"test3": "testy3",
			},
			want: map[string]string{
				"test1": "testy1",
				"test2": "testy2",
				"test3": "testy3",
			},
		},
		{
			name: "one map is nil",
			map1: nil,
			map2: map[string]string{
				"test1": "testy1",
				"test3": "testy3",
			},
			want: map[string]string{
				"test1": "testy1",
				"test3": "testy3",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, MergeMaps(tt.map1, tt.map2), tt.want)
		})
	}
}

func Test_resourceRefUIDIndexFunc(t *testing.T) {
	type args struct {
		rawObj client.Object
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "get uid of resource reference",
			args: args{
				rawObj: &apilifecycle.KeptnWorkloadVersion{
					Spec: apilifecycle.KeptnWorkloadVersionSpec{
						KeptnWorkloadSpec: apilifecycle.KeptnWorkloadSpec{
							ResourceReference: apilifecycle.ResourceReference{
								UID: "my-uid",
							},
						},
					},
				},
			},
			want: []string{"my-uid"},
		},
		{
			name: "empty uid",
			args: args{
				rawObj: &apilifecycle.KeptnWorkloadVersion{
					Spec: apilifecycle.KeptnWorkloadVersionSpec{
						KeptnWorkloadSpec: apilifecycle.KeptnWorkloadSpec{
							ResourceReference: apilifecycle.ResourceReference{
								UID: "",
							},
						},
					},
				},
			},
			want: nil,
		},
		{
			name: "not a KeptnWorkloadVersion",
			args: args{
				rawObj: &v1.Pod{},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KeptnWorkloadVersionResourceRefUIDIndexFunc(tt.args.rawObj); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KeptnWorkloadVersionResourceRefUIDIndexFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}
