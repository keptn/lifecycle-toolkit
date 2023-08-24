package v1alpha2

import (
	"testing"
	"time"

	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha2/common"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestKeptnWorkloadInstance(t *testing.T) {
	workload := &KeptnWorkloadInstance{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "workload",
			Namespace: "namespace",
		},
		Status: KeptnWorkloadInstanceStatus{
			PreDeploymentStatus:            common.StateFailed,
			PreDeploymentEvaluationStatus:  common.StateFailed,
			PostDeploymentStatus:           common.StateFailed,
			PostDeploymentEvaluationStatus: common.StateFailed,
			DeploymentStatus:               common.StateFailed,
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
		Spec: KeptnWorkloadInstanceSpec{
			KeptnWorkloadSpec: KeptnWorkloadSpec{
				PreDeploymentTasks:        []string{"task1", "task2"},
				PostDeploymentTasks:       []string{"task3", "task4"},
				PreDeploymentEvaluations:  []string{"task5", "task6"},
				PostDeploymentEvaluations: []string{"task7", "task8"},
				Version:                   "version",
				AppName:                   "appname",
			},
			PreviousVersion: "prev",
			WorkloadName:    "workloadname",
			TraceId:         map[string]string{"traceparent": "trace1"},
		},
	}

	require.True(t, workload.IsPreDeploymentCompleted())
	require.False(t, workload.IsPreDeploymentSucceeded())
	require.True(t, workload.IsPreDeploymentFailed())

	require.True(t, workload.IsPreDeploymentEvaluationCompleted())
	require.False(t, workload.IsPreDeploymentEvaluationSucceeded())
	require.True(t, workload.IsPreDeploymentEvaluationFailed())

	require.True(t, workload.IsPostDeploymentCompleted())
	require.False(t, workload.IsPostDeploymentSucceeded())
	require.True(t, workload.IsPostDeploymentFailed())

	require.True(t, workload.IsPostDeploymentEvaluationCompleted())
	require.False(t, workload.IsPostDeploymentEvaluationSucceeded())
	require.True(t, workload.IsPostDeploymentEvaluationFailed())

	require.True(t, workload.IsDeploymentCompleted())
	require.False(t, workload.IsDeploymentSucceeded())
	require.True(t, workload.IsDeploymentFailed())

	require.False(t, workload.IsEndTimeSet())
	require.False(t, workload.IsStartTimeSet())

	workload.SetStartTime()
	workload.SetEndTime()

	require.True(t, workload.IsEndTimeSet())
	require.True(t, workload.IsStartTimeSet())

	require.Equal(t, []attribute.KeyValue{
		common.AppName.String("appname"),
		common.WorkloadName.String("workloadname"),
		common.WorkloadVersion.String("version"),
		common.WorkloadNamespace.String("namespace"),
	}, workload.GetActiveMetricsAttributes())

	require.Equal(t, []attribute.KeyValue{
		common.AppName.String("appname"),
		common.WorkloadName.String("workloadname"),
		common.WorkloadVersion.String("version"),
		common.WorkloadNamespace.String("namespace"),
		common.WorkloadStatus.String(string(common.StateFailed)),
	}, workload.GetMetricsAttributes())

	require.Equal(t, []attribute.KeyValue{
		common.AppName.String("appname"),
		common.WorkloadName.String("workloadname"),
		common.WorkloadVersion.String("version"),
		common.WorkloadPreviousVersion.String("prev"),
	}, workload.GetDurationMetricsAttributes())

	require.Equal(t, common.StateFailed, workload.GetState())

	require.Equal(t, []string{"task1", "task2"}, workload.GetPreDeploymentTasks())
	require.Equal(t, []string{"task3", "task4"}, workload.GetPostDeploymentTasks())
	require.Equal(t, []string{"task5", "task6"}, workload.GetPreDeploymentEvaluations())
	require.Equal(t, []string{"task7", "task8"}, workload.GetPostDeploymentEvaluations())

	require.Equal(t, []ItemStatus{
		{
			DefinitionName: "defname",
			Status:         common.StateFailed,
			Name:           "taskname",
		},
	}, workload.GetPreDeploymentTaskStatus())

	require.Equal(t, []ItemStatus{
		{
			DefinitionName: "defname2",
			Status:         common.StateFailed,
			Name:           "taskname2",
		},
	}, workload.GetPostDeploymentTaskStatus())

	require.Equal(t, []ItemStatus{
		{
			DefinitionName: "defname3",
			Status:         common.StateFailed,
			Name:           "taskname3",
		},
	}, workload.GetPreDeploymentEvaluationTaskStatus())

	require.Equal(t, []ItemStatus{
		{
			DefinitionName: "defname4",
			Status:         common.StateFailed,
			Name:           "taskname4",
		},
	}, workload.GetPostDeploymentEvaluationTaskStatus())

	require.Equal(t, "appname", workload.GetAppName())
	require.Equal(t, "prev", workload.GetPreviousVersion())
	require.Equal(t, "workloadname", workload.GetParentName())
	require.Equal(t, "namespace", workload.GetNamespace())

	workload.SetState(common.StatePending)
	require.Equal(t, common.StatePending, workload.GetState())

	require.True(t, !workload.GetStartTime().IsZero())
	require.True(t, !workload.GetEndTime().IsZero())

	workload.SetCurrentPhase(common.PhaseAppDeployment.LongName)
	require.Equal(t, common.PhaseAppDeployment.LongName, workload.GetCurrentPhase())

	workload.Status.EndTime = v1.Time{Time: time.Time{}}
	workload.Complete()
	require.True(t, !workload.GetEndTime().IsZero())

	require.Equal(t, "version", workload.GetVersion())

	require.Equal(t, "trace1.workloadname.version.phase", workload.GetSpanKey("phase"))

	task := workload.GenerateTask("taskdef", common.PostDeploymentCheckType)
	require.Equal(t, KeptnTaskSpec{
		AppName:          workload.GetAppName(),
		WorkloadVersion:  workload.GetVersion(),
		Workload:         workload.GetParentName(),
		TaskDefinition:   "taskdef",
		Parameters:       TaskParameters{},
		SecureParameters: SecureParameters{},
		Type:             common.PostDeploymentCheckType,
	}, task.Spec)

	evaluation := workload.GenerateEvaluation("taskdef", common.PostDeploymentCheckType)
	require.Equal(t, KeptnEvaluationSpec{
		AppName:              workload.GetAppName(),
		WorkloadVersion:      workload.GetVersion(),
		Workload:             workload.GetParentName(),
		EvaluationDefinition: "taskdef",
		Type:                 common.PostDeploymentCheckType,
		RetryInterval: metav1.Duration{
			Duration: 5 * time.Second,
		},
	}, evaluation.Spec)

	require.Equal(t, "workload", workload.GetSpanName(""))

	require.Equal(t, "workloadname/phase", workload.GetSpanName("phase"))

	require.Equal(t, []attribute.KeyValue{
		common.AppName.String("appname"),
		common.WorkloadName.String("workloadname"),
		common.WorkloadVersion.String("version"),
		common.WorkloadNamespace.String("namespace"),
	}, workload.GetSpanAttributes())

	require.Equal(t, map[string]string{
		"appName":              "appname",
		"workloadName":         "workloadname",
		"workloadVersion":      "version",
		"workloadInstanceName": "workload",
	}, workload.GetEventAnnotations())
}

//nolint:dupl
func TestKeptnWorkloadInstance_DeprecateRemainingPhases(t *testing.T) {
	workloadInstance := KeptnWorkloadInstance{
		Status: KeptnWorkloadInstanceStatus{
			PreDeploymentStatus:            common.StatePending,
			PreDeploymentEvaluationStatus:  common.StatePending,
			PostDeploymentStatus:           common.StatePending,
			PostDeploymentEvaluationStatus: common.StatePending,
			DeploymentStatus:               common.StatePending,
			Status:                         common.StatePending,
		},
	}

	tests := []struct {
		workloadInstance KeptnWorkloadInstance
		phase            common.KeptnPhaseType
		want             KeptnWorkloadInstance
	}{
		{
			workloadInstance: workloadInstance,
			phase:            common.PhaseWorkloadPostEvaluation,
			want: KeptnWorkloadInstance{
				Status: KeptnWorkloadInstanceStatus{
					PreDeploymentStatus:            common.StatePending,
					PreDeploymentEvaluationStatus:  common.StatePending,
					PostDeploymentStatus:           common.StatePending,
					PostDeploymentEvaluationStatus: common.StatePending,
					DeploymentStatus:               common.StatePending,
					Status:                         common.StatePending,
				},
			},
		},
		{
			workloadInstance: workloadInstance,
			phase:            common.PhaseWorkloadPostDeployment,
			want: KeptnWorkloadInstance{
				Status: KeptnWorkloadInstanceStatus{
					PreDeploymentStatus:            common.StatePending,
					PreDeploymentEvaluationStatus:  common.StatePending,
					PostDeploymentStatus:           common.StatePending,
					PostDeploymentEvaluationStatus: common.StateDeprecated,
					DeploymentStatus:               common.StatePending,
					Status:                         common.StateFailed,
				},
			},
		},
		{
			workloadInstance: workloadInstance,
			phase:            common.PhaseWorkloadDeployment,
			want: KeptnWorkloadInstance{
				Status: KeptnWorkloadInstanceStatus{
					PreDeploymentStatus:            common.StatePending,
					PreDeploymentEvaluationStatus:  common.StatePending,
					PostDeploymentStatus:           common.StateDeprecated,
					PostDeploymentEvaluationStatus: common.StateDeprecated,
					DeploymentStatus:               common.StatePending,
					Status:                         common.StateFailed,
				},
			},
		},
		{
			workloadInstance: workloadInstance,
			phase:            common.PhaseWorkloadPreEvaluation,
			want: KeptnWorkloadInstance{
				Status: KeptnWorkloadInstanceStatus{
					PreDeploymentStatus:            common.StatePending,
					PreDeploymentEvaluationStatus:  common.StatePending,
					PostDeploymentStatus:           common.StateDeprecated,
					PostDeploymentEvaluationStatus: common.StateDeprecated,
					DeploymentStatus:               common.StateDeprecated,
					Status:                         common.StateFailed,
				},
			},
		},
		{
			workloadInstance: workloadInstance,
			phase:            common.PhaseWorkloadPreDeployment,
			want: KeptnWorkloadInstance{
				Status: KeptnWorkloadInstanceStatus{
					PreDeploymentStatus:            common.StatePending,
					PreDeploymentEvaluationStatus:  common.StateDeprecated,
					PostDeploymentStatus:           common.StateDeprecated,
					PostDeploymentEvaluationStatus: common.StateDeprecated,
					DeploymentStatus:               common.StateDeprecated,
					Status:                         common.StateFailed,
				},
			},
		},
		{
			workloadInstance: workloadInstance,
			phase:            common.PhaseDeprecated,
			want: KeptnWorkloadInstance{
				Status: KeptnWorkloadInstanceStatus{
					PreDeploymentStatus:            common.StateDeprecated,
					PreDeploymentEvaluationStatus:  common.StateDeprecated,
					PostDeploymentStatus:           common.StateDeprecated,
					PostDeploymentEvaluationStatus: common.StateDeprecated,
					DeploymentStatus:               common.StateDeprecated,
					Status:                         common.StateDeprecated,
				},
			},
		},
		{
			workloadInstance: workloadInstance,
			phase:            common.PhaseAppPreDeployment,
			want: KeptnWorkloadInstance{
				Status: KeptnWorkloadInstanceStatus{
					PreDeploymentStatus:            common.StatePending,
					PreDeploymentEvaluationStatus:  common.StatePending,
					PostDeploymentStatus:           common.StatePending,
					PostDeploymentEvaluationStatus: common.StatePending,
					DeploymentStatus:               common.StatePending,
					Status:                         common.StateFailed,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			tt.workloadInstance.DeprecateRemainingPhases(tt.phase)
			require.Equal(t, tt.want, tt.workloadInstance)
		})
	}
}

func TestKeptnWorkloadInstance_SetPhaseTraceID(t *testing.T) {
	app := KeptnWorkloadInstance{
		Status: KeptnWorkloadInstanceStatus{},
	}

	app.SetPhaseTraceID(common.PhaseAppDeployment.ShortName, propagation.MapCarrier{
		"name3": "trace3",
	})

	require.Equal(t, KeptnWorkloadInstance{
		Status: KeptnWorkloadInstanceStatus{
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

	require.Equal(t, KeptnWorkloadInstance{
		Status: KeptnWorkloadInstanceStatus{
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

func TestKeptnWorkloadInstanceList(t *testing.T) {
	list := KeptnWorkloadInstanceList{
		Items: []KeptnWorkloadInstance{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "obj1",
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "obj2",
				},
			},
		},
	}

	got := list.GetItems()
	require.Len(t, got, 2)
	require.Equal(t, "obj1", got[0].GetName())
	require.Equal(t, "obj2", got[1].GetName())
}
