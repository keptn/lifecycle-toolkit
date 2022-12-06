package v1alpha2_test

import (
	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha2"
	"testing"

	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha2/common"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestKeptnTask(t *testing.T) {
	task := &v1alpha2.KeptnTask{
		ObjectMeta: metav1.ObjectMeta{
			Name: "task",
		},
		Spec: v1alpha2.KeptnTaskSpec{
			AppName:    "app",
			AppVersion: "appversion",
			Type:       common.PostDeploymentCheckType,
		},
		Status: v1alpha2.KeptnTaskStatus{
			Status: common.StateFailed,
		},
	}

	task.SetPhaseTraceID("", nil)
	require.Equal(t, v1alpha2.KeptnTask{
		ObjectMeta: metav1.ObjectMeta{
			Name: "task",
		},
		Spec: v1alpha2.KeptnTaskSpec{
			AppName:    "app",
			AppVersion: "appversion",
			Type:       common.PostDeploymentCheckType,
		},
		Status: v1alpha2.KeptnTaskStatus{
			Status: common.StateFailed,
		},
	}, *task)

	require.Equal(t, "task", task.GetSpanKey(""))
	require.Equal(t, "task", task.GetSpanName(""))

	require.False(t, task.IsEndTimeSet())
	require.False(t, task.IsStartTimeSet())

	task.SetStartTime()
	task.SetEndTime()

	require.True(t, task.IsEndTimeSet())
	require.True(t, task.IsStartTimeSet())

	require.Equal(t, []attribute.KeyValue{
		common.AppName.String("app"),
		common.AppVersion.String("appversion"),
		common.WorkloadName.String(""),
		common.WorkloadVersion.String(""),
		common.TaskName.String("task"),
		common.TaskType.String(string(common.PostDeploymentCheckType)),
	}, task.GetActiveMetricsAttributes())

	require.Equal(t, []attribute.KeyValue{
		common.AppName.String("app"),
		common.AppVersion.String("appversion"),
		common.WorkloadName.String(""),
		common.WorkloadVersion.String(""),
		common.TaskName.String("task"),
		common.TaskType.String(string(common.PostDeploymentCheckType)),
		common.TaskStatus.String(string(common.StateFailed)),
	}, task.GetMetricsAttributes())

	require.Equal(t, map[string]string{
		"keptn.sh/app":       "app",
		"keptn.sh/task-name": "task",
		"keptn.sh/version":   "appversion",
	}, task.CreateKeptnLabels())

	task.Spec.Workload = "workload"
	task.Spec.WorkloadVersion = "workloadversion"

	require.Equal(t, map[string]string{
		"keptn.sh/app":       "app",
		"keptn.sh/workload":  "workload",
		"keptn.sh/task-name": "task",
		"keptn.sh/version":   "workloadversion",
	}, task.CreateKeptnLabels())

	require.Equal(t, []attribute.KeyValue{
		common.AppName.String("app"),
		common.AppVersion.String("appversion"),
		common.WorkloadName.String("workload"),
		common.WorkloadVersion.String("workloadversion"),
		common.TaskName.String("task"),
		common.TaskType.String(string(common.PostDeploymentCheckType)),
	}, task.GetSpanAttributes())

}

func TestKeptnTaskList(t *testing.T) {
	list := v1alpha2.KeptnTaskList{
		Items: []v1alpha2.KeptnTask{
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
}
