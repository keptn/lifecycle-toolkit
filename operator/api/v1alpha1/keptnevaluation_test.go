package v1alpha1_test

import (
	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1"
	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1/common"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestKeptnEvaluation(t *testing.T) {
	evaluation := &v1alpha1.KeptnEvaluation{
		ObjectMeta: metav1.ObjectMeta{
			Name: "evaluation",
		},
		Spec: v1alpha1.KeptnEvaluationSpec{
			AppName:    "app",
			AppVersion: "appversion",
			Type:       common.PostDeploymentCheckType,
		},
		Status: v1alpha1.KeptnEvaluationStatus{
			OverallStatus: common.StateFailed,
		},
	}

	evaluation.SetPhaseTraceID("", nil)
	require.Equal(t, v1alpha1.KeptnEvaluation{
		ObjectMeta: metav1.ObjectMeta{
			Name: "evaluation",
		},
		Spec: v1alpha1.KeptnEvaluationSpec{
			AppName:    "app",
			AppVersion: "appversion",
			Type:       common.PostDeploymentCheckType,
		},
		Status: v1alpha1.KeptnEvaluationStatus{
			OverallStatus: common.StateFailed,
		},
	}, *evaluation)

	require.Equal(t, "evaluation", evaluation.GetSpanKey(""))
	require.Equal(t, "evaluation", evaluation.GetSpanName(""))

	require.False(t, evaluation.IsEndTimeSet())
	require.False(t, evaluation.IsStartTimeSet())

	evaluation.SetStartTime()
	evaluation.SetEndTime()

	require.True(t, evaluation.IsEndTimeSet())
	require.True(t, evaluation.IsStartTimeSet())

	require.Equal(t, []attribute.KeyValue{
		common.AppName.String("app"),
		common.AppVersion.String("appversion"),
		common.WorkloadName.String(""),
		common.WorkloadVersion.String(""),
		common.EvaluationName.String("evaluation"),
		common.EvaluationType.String(string(common.PostDeploymentCheckType)),
	}, evaluation.GetActiveMetricsAttributes())

	require.Equal(t, []attribute.KeyValue{
		common.AppName.String("app"),
		common.AppVersion.String("appversion"),
		common.WorkloadName.String(""),
		common.WorkloadVersion.String(""),
		common.EvaluationName.String("evaluation"),
		common.EvaluationType.String(string(common.PostDeploymentCheckType)),
		common.EvaluationStatus.String(string(common.StateFailed)),
	}, evaluation.GetMetricsAttributes())

	evaluation.AddEvaluationStatus(v1alpha1.Objective{Name: "objName"})
	require.Equal(t, v1alpha1.EvaluationStatusItem{
		Status: common.StatePending,
	}, evaluation.Status.EvaluationStatus["objName"])

	require.Equal(t, []attribute.KeyValue{
		common.AppName.String("app"),
		common.AppVersion.String("appversion"),
		common.WorkloadName.String(""),
		common.WorkloadVersion.String(""),
		common.EvaluationName.String("evaluation"),
		common.EvaluationType.String(string(common.PostDeploymentCheckType)),
	}, evaluation.GetSpanAttributes())

}

func TestKeptnEvaluationList(t *testing.T) {
	list := v1alpha1.KeptnEvaluationList{
		Items: []v1alpha1.KeptnEvaluation{
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
