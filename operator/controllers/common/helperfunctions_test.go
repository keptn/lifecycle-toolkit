package common

import (
	"testing"

	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3/common"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
				"appRevision": "1",
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
					AppName:         "app",
					AppVersion:      "1.0.0",
					Workload:        "workload",
					WorkloadVersion: "2.0.0",
					TaskDefinition:  "def",
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
