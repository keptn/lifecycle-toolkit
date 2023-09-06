package common

import (
	"context"
	"testing"

	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func Test_GetItemStatus(t *testing.T) {
	tests := []struct {
		name     string
		inStatus []klcv1alpha3.ItemStatus
		want     klcv1alpha3.ItemStatus
	}{
		{
			name: "non-existing",
			inStatus: []klcv1alpha3.ItemStatus{
				{
					DefinitionName: "def-name",
					Name:           "name",
					Status:         apicommon.StatePending,
				},
			},
			want: klcv1alpha3.ItemStatus{
				DefinitionName: "non-existing",
				Status:         apicommon.StatePending,
				Name:           "",
			},
		},
		{
			name: "def-name",
			inStatus: []klcv1alpha3.ItemStatus{
				{
					DefinitionName: "def-name",
					Name:           "name",
					Status:         apicommon.StateProgressing,
				},
			},
			want: klcv1alpha3.ItemStatus{
				DefinitionName: "def-name",
				Name:           "name",
				Status:         apicommon.StateProgressing,
			},
		},
		{
			name:     "empty",
			inStatus: []klcv1alpha3.ItemStatus{},
			want: klcv1alpha3.ItemStatus{
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
		statuses       []klcv1alpha3.ItemStatus
		definitionName string
		want           apicommon.KeptnState
	}{
		{
			statuses:       []klcv1alpha3.ItemStatus{},
			definitionName: "",
			want:           "",
		},
		{
			statuses:       []klcv1alpha3.ItemStatus{},
			definitionName: "defName",
			want:           "",
		},
		{
			statuses: []klcv1alpha3.ItemStatus{
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
			statuses: []klcv1alpha3.ItemStatus{
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

func Test_setEventMessage(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    string
	}{
		{
			name:    "version empty",
			version: "",
			want:    "App Deployment: longReason / Namespace: namespace, Name: app",
		},
		{
			name:    "version set",
			version: "1.0.0",
			want:    "App Deployment: longReason / Namespace: namespace, Name: app, Version: 1.0.0",
		},
	}

	appVersion := &klcv1alpha3.KeptnAppVersion{
		ObjectMeta: v1.ObjectMeta{
			Name:      "app",
			Namespace: "namespace",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, setEventMessage(apicommon.PhaseAppDeployment, appVersion, "longReason", tt.version), tt.want)
		})
	}
}

func Test_setAnnotations(t *testing.T) {
	tests := []struct {
		name   string
		object client.Object
		want   map[string]string
	}{
		{
			name:   "nil object",
			object: nil,
			want:   nil,
		},
		{
			name:   "empty object",
			object: &klcv1alpha3.KeptnEvaluationDefinition{},
			want:   nil,
		},
		{
			name: "unknown object",
			object: &klcv1alpha3.KeptnEvaluationDefinition{
				ObjectMeta: v1.ObjectMeta{
					Name:      "def",
					Namespace: "namespace",
				},
			},
			want: map[string]string{
				"namespace":   "namespace",
				"name":        "def",
				"phase":       "AppDeploy",
				"traceparent": "",
			},
		},
		{
			name: "object with traceparent",
			object: &klcv1alpha3.KeptnEvaluationDefinition{
				ObjectMeta: v1.ObjectMeta{
					Name:      "def",
					Namespace: "namespace",
					Annotations: map[string]string{
						"traceparent": "23232333",
					},
				},
			},
			want: map[string]string{
				"namespace":   "namespace",
				"name":        "def",
				"phase":       "AppDeploy",
				"traceparent": "23232333",
			},
		},
		{
			name: "KeptnApp",
			object: &klcv1alpha3.KeptnApp{
				ObjectMeta: v1.ObjectMeta{
					Name:       "app",
					Namespace:  "namespace",
					Generation: 1,
				},
				Spec: klcv1alpha3.KeptnAppSpec{
					Version: "1.0.0",
				},
			},
			want: map[string]string{
				"namespace":   "namespace",
				"name":        "app",
				"phase":       "AppDeploy",
				"appName":     "app",
				"appVersion":  "1.0.0",
				"appRevision": "6b86b273",
				"traceparent": "",
			},
		},
		{
			name: "KeptnAppVersion",
			object: &klcv1alpha3.KeptnAppVersion{
				ObjectMeta: v1.ObjectMeta{
					Name:      "appVersion",
					Namespace: "namespace",
				},
				Spec: klcv1alpha3.KeptnAppVersionSpec{
					AppName: "app",
					KeptnAppSpec: klcv1alpha3.KeptnAppSpec{
						Version: "1.0.0",
					},
				},
			},
			want: map[string]string{
				"namespace":      "namespace",
				"name":           "appVersion",
				"phase":          "AppDeploy",
				"appName":        "app",
				"appVersion":     "1.0.0",
				"appVersionName": "appVersion",
				"traceparent":    "",
			},
		},
		{
			name: "KeptnWorkload",
			object: &klcv1alpha3.KeptnWorkload{
				ObjectMeta: v1.ObjectMeta{
					Name:      "workload",
					Namespace: "namespace",
				},
				Spec: klcv1alpha3.KeptnWorkloadSpec{
					AppName: "app",
					Version: "1.0.0",
				},
			},
			want: map[string]string{
				"namespace":       "namespace",
				"name":            "workload",
				"phase":           "AppDeploy",
				"appName":         "app",
				"workloadVersion": "1.0.0",
				"workloadName":    "workload",
				"traceparent":     "",
			},
		},
		{
			name: "KeptnWorkloadInstance",
			object: &klcv1alpha3.KeptnWorkloadInstance{
				ObjectMeta: v1.ObjectMeta{
					Name:      "workloadInstance",
					Namespace: "namespace",
				},
				Spec: klcv1alpha3.KeptnWorkloadInstanceSpec{
					KeptnWorkloadSpec: klcv1alpha3.KeptnWorkloadSpec{
						AppName: "app",
						Version: "1.0.0",
					},
					WorkloadName: "workload",
				},
			},
			want: map[string]string{
				"namespace":            "namespace",
				"name":                 "workloadInstance",
				"phase":                "AppDeploy",
				"appName":              "app",
				"workloadVersion":      "1.0.0",
				"workloadName":         "workload",
				"workloadInstanceName": "workloadInstance",
				"traceparent":          "",
			},
		},
		{
			name: "KeptnTask",
			object: &klcv1alpha3.KeptnTask{
				ObjectMeta: v1.ObjectMeta{
					Name:      "task",
					Namespace: "namespace",
				},
				Spec: klcv1alpha3.KeptnTaskSpec{
					TaskDefinition: "def",
					Context: klcv1alpha3.TaskContext{
						WorkloadName:    "workload",
						AppName:         "app",
						AppVersion:      "1.0.0",
						WorkloadVersion: "2.0.0",
					},
				},
			},
			want: map[string]string{
				"namespace":          "namespace",
				"name":               "task",
				"phase":              "AppDeploy",
				"appName":            "app",
				"appVersion":         "1.0.0",
				"workloadName":       "workload",
				"workloadVersion":    "2.0.0",
				"taskDefinitionName": "def",
				"taskName":           "task",
				"traceparent":        "",
			},
		},
		{
			name: "KeptnEvaluation",
			object: &klcv1alpha3.KeptnEvaluation{
				ObjectMeta: v1.ObjectMeta{
					Name:      "eval",
					Namespace: "namespace",
				},
				Spec: klcv1alpha3.KeptnEvaluationSpec{
					AppName:              "app",
					AppVersion:           "1.0.0",
					Workload:             "workload",
					WorkloadVersion:      "2.0.0",
					EvaluationDefinition: "def",
				},
			},
			want: map[string]string{
				"namespace":                "namespace",
				"name":                     "eval",
				"phase":                    "AppDeploy",
				"appName":                  "app",
				"appVersion":               "1.0.0",
				"workloadName":             "workload",
				"workloadVersion":          "2.0.0",
				"evaluationDefinitionName": "def",
				"evaluationName":           "eval",
				"traceparent":              "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, setAnnotations(tt.object, apicommon.PhaseAppDeployment), tt.want)
		})
	}
}

//nolint:dupl
func Test_GetTaskDefinition(t *testing.T) {
	tests := []struct {
		name             string
		taskDef          *klcv1alpha3.KeptnTaskDefinition
		taskDefName      string
		taskDefNamespace string
		out              *klcv1alpha3.KeptnTaskDefinition
		wantError        bool
	}{
		{
			name: "taskDef not found",
			taskDef: &klcv1alpha3.KeptnTaskDefinition{
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
			taskDef: &klcv1alpha3.KeptnTaskDefinition{
				ObjectMeta: v1.ObjectMeta{
					Name:      "taskDef",
					Namespace: "some-namespace",
				},
			},
			taskDefName:      "taskDef",
			taskDefNamespace: "some-namespace",
			out: &klcv1alpha3.KeptnTaskDefinition{
				ObjectMeta: v1.ObjectMeta{
					Name:      "taskDef",
					Namespace: "some-namespace",
				},
			},
			wantError: false,
		},
		{
			name: "taskDef found in default KLT namespace",
			taskDef: &klcv1alpha3.KeptnTaskDefinition{
				ObjectMeta: v1.ObjectMeta{
					Name:      "taskDef",
					Namespace: KLTNamespace,
				},
			},
			taskDefName:      "taskDef",
			taskDefNamespace: "some-namespace",
			out: &klcv1alpha3.KeptnTaskDefinition{
				ObjectMeta: v1.ObjectMeta{
					Name:      "taskDef",
					Namespace: KLTNamespace,
				},
			},
			wantError: false,
		},
	}

	err := klcv1alpha3.AddToScheme(scheme.Scheme)
	require.Nil(t, err)

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
		evalDef          *klcv1alpha3.KeptnEvaluationDefinition
		evalDefName      string
		evalDefNamespace string
		out              *klcv1alpha3.KeptnEvaluationDefinition
		wantError        bool
	}{
		{
			name: "evalDef not found",
			evalDef: &klcv1alpha3.KeptnEvaluationDefinition{
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
			evalDef: &klcv1alpha3.KeptnEvaluationDefinition{
				ObjectMeta: v1.ObjectMeta{
					Name:      "evalDef",
					Namespace: "some-namespace",
				},
			},
			evalDefName:      "evalDef",
			evalDefNamespace: "some-namespace",
			out: &klcv1alpha3.KeptnEvaluationDefinition{
				ObjectMeta: v1.ObjectMeta{
					Name:      "evalDef",
					Namespace: "some-namespace",
				},
			},
			wantError: false,
		},
		{
			name: "evalDef found in default KLT namespace",
			evalDef: &klcv1alpha3.KeptnEvaluationDefinition{
				ObjectMeta: v1.ObjectMeta{
					Name:      "evalDef",
					Namespace: KLTNamespace,
				},
			},
			evalDefName:      "evalDef",
			evalDefNamespace: "some-namespace",
			out: &klcv1alpha3.KeptnEvaluationDefinition{
				ObjectMeta: v1.ObjectMeta{
					Name:      "evalDef",
					Namespace: KLTNamespace,
				},
			},
			wantError: false,
		},
	}

	err := klcv1alpha3.AddToScheme(scheme.Scheme)
	require.Nil(t, err)

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
