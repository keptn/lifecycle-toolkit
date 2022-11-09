package api_test

import (
	"testing"
	"time"

	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1"
	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1/common"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestKeptnAppVersion(t *testing.T) {
	app := &v1alpha1.KeptnAppVersion{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "app",
			Namespace: "namespace",
		},
		Status: v1alpha1.KeptnAppVersionStatus{
			PreDeploymentStatus:            common.StateFailed,
			PreDeploymentEvaluationStatus:  common.StateFailed,
			PostDeploymentStatus:           common.StateFailed,
			PostDeploymentEvaluationStatus: common.StateFailed,
			WorkloadOverallStatus:          common.StateFailed,
			Status:                         common.StateFailed,
			PreDeploymentTaskStatus: []v1alpha1.TaskStatus{
				{
					TaskDefinitionName: "defname",
					Status:             common.StateFailed,
					TaskName:           "taskname",
				},
			},
			PostDeploymentTaskStatus: []v1alpha1.TaskStatus{
				{
					TaskDefinitionName: "defname2",
					Status:             common.StateFailed,
					TaskName:           "taskname2",
				},
			},
			PreDeploymentEvaluationTaskStatus: []v1alpha1.EvaluationStatus{
				{
					EvaluationDefinitionName: "defname3",
					Status:                   common.StateFailed,
					EvaluationName:           "taskname3",
				},
			},
			PostDeploymentEvaluationTaskStatus: []v1alpha1.EvaluationStatus{
				{
					EvaluationDefinitionName: "defname4",
					Status:                   common.StateFailed,
					EvaluationName:           "taskname4",
				},
			},
			CurrentPhase: common.PhaseAppDeployment.ShortName,
		},
		Spec: v1alpha1.KeptnAppVersionSpec{
			KeptnAppSpec: v1alpha1.KeptnAppSpec{
				PreDeploymentTasks:        []string{"task1", "task2"},
				PostDeploymentTasks:       []string{"task3", "task4"},
				PreDeploymentEvaluations:  []string{"task5", "task6"},
				PostDeploymentEvaluations: []string{"task7", "task8"},
				Version:                   "version",
			},
			PreviousVersion: "prev",
			AppName:         "appname",
			TraceId:         map[string]string{"traceparent": "trace1"},
		},
	}

	require.True(t, app.IsPreDeploymentCompleted())
	require.False(t, app.IsPreDeploymentSucceeded())
	require.True(t, app.IsPreDeploymentFailed())

	require.True(t, app.IsPreDeploymentEvaluationCompleted())
	require.False(t, app.IsPreDeploymentEvaluationSucceeded())
	require.True(t, app.IsPreDeploymentEvaluationFailed())

	require.True(t, app.IsPostDeploymentCompleted())
	require.False(t, app.IsPostDeploymentSucceeded())
	require.True(t, app.IsPostDeploymentFailed())

	require.True(t, app.IsPostDeploymentEvaluationCompleted())
	require.False(t, app.IsPostDeploymentEvaluationSucceeded())
	require.True(t, app.IsPostDeploymentEvaluationFailed())

	require.True(t, app.AreWorkloadsCompleted())
	require.False(t, app.AreWorkloadsSucceeded())
	require.True(t, app.AreWorkloadsFailed())

	require.False(t, app.IsEndTimeSet())
	require.False(t, app.IsStartTimeSet())

	app.SetStartTime()
	app.SetEndTime()

	require.True(t, app.IsEndTimeSet())
	require.True(t, app.IsStartTimeSet())

	require.Equal(t, []attribute.KeyValue{
		common.AppName.String("appname"),
		common.AppVersion.String("version"),
		common.AppNamespace.String("namespace"),
	}, app.GetActiveMetricsAttributes())

	require.Equal(t, []attribute.KeyValue{
		common.AppName.String("appname"),
		common.AppVersion.String("version"),
		common.AppNamespace.String("namespace"),
		common.AppStatus.String(string(common.StateFailed)),
	}, app.GetMetricsAttributes())

	require.Equal(t, []attribute.KeyValue{
		common.AppName.String("appname"),
		common.AppVersion.String("version"),
		common.AppPreviousVersion.String("prev"),
	}, app.GetDurationMetricsAttributes())

	require.Equal(t, common.StateFailed, app.GetState())

	require.Equal(t, []string{"task1", "task2"}, app.GetPreDeploymentTasks())
	require.Equal(t, []string{"task3", "task4"}, app.GetPostDeploymentTasks())
	require.Equal(t, []string{"task5", "task6"}, app.GetPreDeploymentEvaluations())
	require.Equal(t, []string{"task7", "task8"}, app.GetPostDeploymentEvaluations())

	require.Equal(t, []v1alpha1.TaskStatus{
		{
			TaskDefinitionName: "defname",
			Status:             common.StateFailed,
			TaskName:           "taskname",
		},
	}, app.GetPreDeploymentTaskStatus())

	require.Equal(t, []v1alpha1.TaskStatus{
		{
			TaskDefinitionName: "defname2",
			Status:             common.StateFailed,
			TaskName:           "taskname2",
		},
	}, app.GetPostDeploymentTaskStatus())

	require.Equal(t, []v1alpha1.EvaluationStatus{
		{
			EvaluationDefinitionName: "defname3",
			Status:                   common.StateFailed,
			EvaluationName:           "taskname3",
		},
	}, app.GetPreDeploymentEvaluationTaskStatus())

	require.Equal(t, []v1alpha1.EvaluationStatus{
		{
			EvaluationDefinitionName: "defname4",
			Status:                   common.StateFailed,
			EvaluationName:           "taskname4",
		},
	}, app.GetPostDeploymentEvaluationTaskStatus())

	require.Equal(t, "appname", app.GetAppName())
	require.Equal(t, "prev", app.GetPreviousVersion())
	require.Equal(t, "appname", app.GetParentName())
	require.Equal(t, "namespace", app.GetNamespace())

	app.SetState(common.StatePending)
	require.Equal(t, common.StatePending, app.GetState())

	require.True(t, !app.GetStartTime().IsZero())
	require.True(t, !app.GetEndTime().IsZero())

	app.SetCurrentPhase(common.PhaseAppDeployment.LongName)
	require.Equal(t, common.PhaseAppDeployment.LongName, app.GetCurrentPhase())

	app.Status.EndTime = v1.Time{Time: time.Time{}}
	app.Complete()
	require.True(t, !app.GetEndTime().IsZero())

	require.Equal(t, "version", app.GetVersion())

	require.Equal(t, "trace1.appname.version.phase", app.GetSpanKey("phase"))

	task := app.GenerateTask(map[string]string{}, "taskdef", common.PostDeploymentCheckType)
	require.Equal(t, v1alpha1.KeptnTaskSpec{
		AppVersion:       app.GetVersion(),
		AppName:          app.GetParentName(),
		TaskDefinition:   "taskdef",
		Parameters:       v1alpha1.TaskParameters{},
		SecureParameters: v1alpha1.SecureParameters{},
		Type:             common.PostDeploymentCheckType,
	}, task.Spec)

	evaluation := app.GenerateEvaluation(map[string]string{}, "taskdef", common.PostDeploymentCheckType)
	require.Equal(t, v1alpha1.KeptnEvaluationSpec{
		AppVersion:           app.GetVersion(),
		AppName:              app.GetParentName(),
		EvaluationDefinition: "taskdef",
		Type:                 common.PostDeploymentCheckType,
		RetryInterval: metav1.Duration{
			Duration: 5 * time.Second,
		},
	}, evaluation.Spec)

	require.Equal(t, "phase", app.GetSpanName("phase"))

	require.Equal(t, []attribute.KeyValue{
		common.AppName.String("appname"),
		common.AppVersion.String("version"),
		common.AppNamespace.String("namespace"),
	}, app.GetSpanAttributes())
}
