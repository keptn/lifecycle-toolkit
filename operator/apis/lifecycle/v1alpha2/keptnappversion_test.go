package v1alpha2

import (
	"testing"
	"time"

	"github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2/common"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestKeptnAppVersion(t *testing.T) {
	app := &KeptnAppVersion{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "app",
			Namespace: "namespace",
		},
		Status: KeptnAppVersionStatus{
			PreDeploymentStatus:            common.StateFailed,
			PreDeploymentEvaluationStatus:  common.StateFailed,
			PostDeploymentStatus:           common.StateFailed,
			PostDeploymentEvaluationStatus: common.StateFailed,
			WorkloadOverallStatus:          common.StateFailed,
			Status:                         common.StateFailed,
			PreDeploymentTaskStatus: []ItemStatus{
				{
					DefinitionName: "defname",
					Status:         common.StateFailed,
					Name:           "taskname",
				},
			},
			PostDeploymentTaskStatus: []ItemStatus{
				{
					DefinitionName: "defname2",
					Status:         common.StateFailed,
					Name:           "taskname2",
				},
			},
			PreDeploymentEvaluationTaskStatus: []ItemStatus{
				{
					DefinitionName: "defname3",
					Status:         common.StateFailed,
					Name:           "taskname3",
				},
			},
			PostDeploymentEvaluationTaskStatus: []ItemStatus{
				{
					DefinitionName: "defname4",
					Status:         common.StateFailed,
					Name:           "taskname4",
				},
			},
			CurrentPhase: common.PhaseAppDeployment.ShortName,
		},
		Spec: KeptnAppVersionSpec{
			KeptnAppSpec: KeptnAppSpec{
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

	require.Equal(t, []ItemStatus{
		{
			DefinitionName: "defname",
			Status:         common.StateFailed,
			Name:           "taskname",
		},
	}, app.GetPreDeploymentTaskStatus())

	require.Equal(t, []ItemStatus{
		{
			DefinitionName: "defname2",
			Status:         common.StateFailed,
			Name:           "taskname2",
		},
	}, app.GetPostDeploymentTaskStatus())

	require.Equal(t, []ItemStatus{
		{
			DefinitionName: "defname3",
			Status:         common.StateFailed,
			Name:           "taskname3",
		},
	}, app.GetPreDeploymentEvaluationTaskStatus())

	require.Equal(t, []ItemStatus{
		{
			DefinitionName: "defname4",
			Status:         common.StateFailed,
			Name:           "taskname4",
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

	task := app.GenerateTask("taskdef", common.PostDeploymentCheckType)
	require.Equal(t, KeptnTaskSpec{
		AppVersion:       app.GetVersion(),
		AppName:          app.GetParentName(),
		TaskDefinition:   "taskdef",
		Parameters:       TaskParameters{},
		SecureParameters: SecureParameters{},
		Type:             common.PostDeploymentCheckType,
	}, task.Spec)

	evaluation := app.GenerateEvaluation("taskdef", common.PostDeploymentCheckType)
	require.Equal(t, KeptnEvaluationSpec{
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

	require.Equal(t, map[string]string{
		"appName":        "appname",
		"appVersion":     "version",
		"appVersionName": "app",
	}, app.GetEventAnnotations())
}

func TestKeptnAppVersion_GetWorkloadNameOfApp(t *testing.T) {
	type fields struct {
		Spec KeptnAppVersionSpec
	}
	type args struct {
		workloadName string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "",
			fields: fields{
				Spec: KeptnAppVersionSpec{AppName: "my-app"},
			},
			args: args{
				workloadName: "my-workload",
			},
			want: "my-app-my-workload",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := KeptnAppVersion{
				Spec: tt.fields.Spec,
			}
			if got := v.GetWorkloadNameOfApp(tt.args.workloadName); got != tt.want {
				t.Errorf("GetWorkloadNameOfApp() = %v, want %v", got, tt.want)
			}
		})
	}
}

//nolint:dupl
func TestKeptnAppVersion_DeprecateRemainingPhases(t *testing.T) {
	app := KeptnAppVersion{
		Status: KeptnAppVersionStatus{
			PreDeploymentStatus:            common.StatePending,
			PreDeploymentEvaluationStatus:  common.StatePending,
			PostDeploymentStatus:           common.StatePending,
			PostDeploymentEvaluationStatus: common.StatePending,
			WorkloadOverallStatus:          common.StatePending,
			Status:                         common.StatePending,
		},
	}

	tests := []struct {
		app   KeptnAppVersion
		phase common.KeptnPhaseType
		want  KeptnAppVersion
	}{
		{
			app:   app,
			phase: common.PhaseAppPostEvaluation,
			want: KeptnAppVersion{
				Status: KeptnAppVersionStatus{
					PreDeploymentStatus:            common.StatePending,
					PreDeploymentEvaluationStatus:  common.StatePending,
					PostDeploymentStatus:           common.StatePending,
					PostDeploymentEvaluationStatus: common.StatePending,
					WorkloadOverallStatus:          common.StatePending,
					Status:                         common.StatePending,
				},
			},
		},
		{
			app:   app,
			phase: common.PhaseAppPostDeployment,
			want: KeptnAppVersion{
				Status: KeptnAppVersionStatus{
					PreDeploymentStatus:            common.StatePending,
					PreDeploymentEvaluationStatus:  common.StatePending,
					PostDeploymentStatus:           common.StatePending,
					PostDeploymentEvaluationStatus: common.StateDeprecated,
					WorkloadOverallStatus:          common.StatePending,
					Status:                         common.StateFailed,
				},
			},
		},
		{
			app:   app,
			phase: common.PhaseAppDeployment,
			want: KeptnAppVersion{
				Status: KeptnAppVersionStatus{
					PreDeploymentStatus:            common.StatePending,
					PreDeploymentEvaluationStatus:  common.StatePending,
					PostDeploymentStatus:           common.StateDeprecated,
					PostDeploymentEvaluationStatus: common.StateDeprecated,
					WorkloadOverallStatus:          common.StatePending,
					Status:                         common.StateFailed,
				},
			},
		},
		{
			app:   app,
			phase: common.PhaseAppPreEvaluation,
			want: KeptnAppVersion{
				Status: KeptnAppVersionStatus{
					PreDeploymentStatus:            common.StatePending,
					PreDeploymentEvaluationStatus:  common.StatePending,
					PostDeploymentStatus:           common.StateDeprecated,
					PostDeploymentEvaluationStatus: common.StateDeprecated,
					WorkloadOverallStatus:          common.StateDeprecated,
					Status:                         common.StateFailed,
				},
			},
		},
		{
			app:   app,
			phase: common.PhaseAppPreDeployment,
			want: KeptnAppVersion{
				Status: KeptnAppVersionStatus{
					PreDeploymentStatus:            common.StatePending,
					PreDeploymentEvaluationStatus:  common.StateDeprecated,
					PostDeploymentStatus:           common.StateDeprecated,
					PostDeploymentEvaluationStatus: common.StateDeprecated,
					WorkloadOverallStatus:          common.StateDeprecated,
					Status:                         common.StateFailed,
				},
			},
		},
		{
			app:   app,
			phase: common.PhaseDeprecated,
			want: KeptnAppVersion{
				Status: KeptnAppVersionStatus{
					PreDeploymentStatus:            common.StateDeprecated,
					PreDeploymentEvaluationStatus:  common.StateDeprecated,
					PostDeploymentStatus:           common.StateDeprecated,
					PostDeploymentEvaluationStatus: common.StateDeprecated,
					WorkloadOverallStatus:          common.StateDeprecated,
					Status:                         common.StateDeprecated,
				},
			},
		},
		{
			app:   app,
			phase: common.PhaseWorkloadPreDeployment,
			want: KeptnAppVersion{
				Status: KeptnAppVersionStatus{
					PreDeploymentStatus:            common.StatePending,
					PreDeploymentEvaluationStatus:  common.StatePending,
					PostDeploymentStatus:           common.StatePending,
					PostDeploymentEvaluationStatus: common.StatePending,
					WorkloadOverallStatus:          common.StatePending,
					Status:                         common.StateFailed,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			tt.app.DeprecateRemainingPhases(tt.phase)
			require.Equal(t, tt.want, tt.app)
		})
	}
}

func TestKeptnAppVersion_SetPhaseTraceID(t *testing.T) {
	app := KeptnAppVersion{
		Status: KeptnAppVersionStatus{},
	}

	app.SetPhaseTraceID(common.PhaseAppDeployment.ShortName, propagation.MapCarrier{
		"name3": "trace3",
	})

	require.Equal(t, KeptnAppVersion{
		Status: KeptnAppVersionStatus{
			PhaseTraceIDs: common.PhaseTraceID{
				common.PhaseAppDeployment.ShortName: propagation.MapCarrier{
					"name3": "trace3",
				},
			},
		},
	}, app)

	app.SetPhaseTraceID(common.PhaseWorkloadDeployment.LongName, propagation.MapCarrier{
		"name2": "trace2",
	})

	require.Equal(t, KeptnAppVersion{
		Status: KeptnAppVersionStatus{
			PhaseTraceIDs: common.PhaseTraceID{
				common.PhaseAppDeployment.ShortName: propagation.MapCarrier{
					"name3": "trace3",
				},
				common.PhaseWorkloadDeployment.ShortName: propagation.MapCarrier{
					"name2": "trace2",
				},
			},
		},
	}, app)
}

func TestKeptnAppVersionList(t *testing.T) {
	list := KeptnAppVersionList{
		Items: []KeptnAppVersion{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "obj1",
				},
				Status: KeptnAppVersionStatus{
					Status: common.StateSucceeded,
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "obj2",
				},
				Status: KeptnAppVersionStatus{
					Status: common.StateDeprecated,
				},
			},
		},
	}

	// fetch the list items
	got := list.GetItems()
	require.Len(t, got, 2)

	require.Equal(t, "obj1", list.Items[0].GetName())
	require.Equal(t, "obj2", list.Items[1].GetName())

	// remove deprecated items from the list
	list.RemoveDeprecated()

	// check that deprecated items are not present in the list anymore
	got = list.GetItems()
	require.Len(t, got, 1)
}
