package v1alpha3

import (
	"testing"

	"github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3/common"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestKeptnEvaluation(t *testing.T) {
	evaluation := &KeptnEvaluation{
		ObjectMeta: metav1.ObjectMeta{
			Name: "evaluation",
		},
		Spec: KeptnEvaluationSpec{
			AppName:              "app",
			AppVersion:           "appversion",
			Type:                 common.PostDeploymentCheckType,
			EvaluationDefinition: "def",
		},
		Status: KeptnEvaluationStatus{
			OverallStatus: common.StateFailed,
		},
	}

	evaluation.SetPhaseTraceID("", nil)
	require.Equal(t, KeptnEvaluation{
		ObjectMeta: metav1.ObjectMeta{
			Name: "evaluation",
		},
		Spec: KeptnEvaluationSpec{
			AppName:              "app",
			AppVersion:           "appversion",
			Type:                 common.PostDeploymentCheckType,
			EvaluationDefinition: "def",
		},
		Status: KeptnEvaluationStatus{
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

	evaluation.AddEvaluationStatus(Objective{KeptnMetricRef: KeptnMetricReference{Name: "objName"}})
	require.Equal(t, EvaluationStatusItem{
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

	require.Equal(t, map[string]string{
		"appName":                  "app",
		"appVersion":               "appversion",
		"workloadName":             "",
		"workloadVersion":          "",
		"evaluationName":           "evaluation",
		"evaluationDefinitionName": "def",
	}, evaluation.GetEventAnnotations())
}

func TestKeptnEvaluationList(t *testing.T) {
	list := KeptnEvaluationList{
		Items: []KeptnEvaluation{
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
