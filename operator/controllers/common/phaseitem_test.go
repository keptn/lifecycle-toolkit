package common

import (
	"testing"
	"time"

	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1"
	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1/common"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/common/fake"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func TestPhaseItemWrapper_GetState(t *testing.T) {
	appVersion := &v1alpha1.KeptnAppVersion{
		Status: v1alpha1.KeptnAppVersionStatus{
			Status:       common.StateFailed,
			CurrentPhase: "test",
		},
	}

	object, err := NewPhaseItemWrapperFromClientObject(appVersion)
	require.Nil(t, err)

	require.Equal(t, "test", object.GetCurrentPhase())

	object.Complete()

	require.NotZero(t, appVersion.Status.EndTime)
}

func TestPhaseItem(t *testing.T) {
	phaseItemMock := fake.PhaseItemMock{
		GetStateFunc: func() common.KeptnState {
			return common.StatePending
		},
		SetStateFunc: func(keptnState common.KeptnState) {
		},
		GetCurrentPhaseFunc: func() string {
			return "phase"
		},
		SetCurrentPhaseFunc: func(s string) {
		},
		GetVersionFunc: func() string {
			return "version"
		},
		GetSpanAttributesFunc: func() []attribute.KeyValue {
			return nil
		},
		GetSpanKeyFunc: func(phase string) string {
			return "span"
		},
		GetSpanNameFunc: func(phase string) string {
			return "name"
		},
		CompleteFunc: func() {
		},
		IsEndTimeSetFunc: func() bool {
			return true
		},
		GetEndTimeFunc: func() time.Time {
			return time.Now().UTC()
		},
		GetStartTimeFunc: func() time.Time {
			return time.Now().UTC()
		},
		GetPreviousVersionFunc: func() string {
			return "version"
		},
		GetParentNameFunc: func() string {
			return "parent"
		},
		GetNamespaceFunc: func() string {
			return "namespace"
		},
		GetAppNameFunc: func() string {
			return "name"
		},
		GetPreDeploymentTasksFunc: func() []string {
			return nil
		},
		GetPostDeploymentTasksFunc: func() []string {
			return nil
		},
		GetPreDeploymentTaskStatusFunc: func() []v1alpha1.TaskStatus {
			return nil
		},
		GetPostDeploymentTaskStatusFunc: func() []v1alpha1.TaskStatus {
			return nil
		},
		GetPreDeploymentEvaluationsFunc: func() []string {
			return nil
		},
		GetPostDeploymentEvaluationsFunc: func() []string {
			return nil
		},
		GetPreDeploymentEvaluationTaskStatusFunc: func() []v1alpha1.EvaluationStatus {
			return nil
		},
		GetPostDeploymentEvaluationTaskStatusFunc: func() []v1alpha1.EvaluationStatus {
			return nil
		},
		GenerateTaskFunc: func(traceContextCarrier propagation.MapCarrier, taskDefinition string, checkType common.CheckType) v1alpha1.KeptnTask {
			return v1alpha1.KeptnTask{}
		},
		GenerateEvaluationFunc: func(traceContextCarrier propagation.MapCarrier, evaluationDefinition string, checkType common.CheckType) v1alpha1.KeptnEvaluation {
			return v1alpha1.KeptnEvaluation{}
		},
		SetSpanAttributesFunc: func(span trace.Span) {
		},
		CancelRemainingPhasesFunc: func(phase common.KeptnPhaseType) {
		},
		SetPhaseTraceIDFunc: func(phase string, carrier propagation.MapCarrier) {
			return
		},
	}

	wrapper := PhaseItemWrapper{Obj: &phaseItemMock}

	_ = wrapper.GetState()
	require.Len(t, phaseItemMock.GetStateCalls(), 1)

	wrapper.SetState(common.StateFailed)
	require.Len(t, phaseItemMock.SetStateCalls(), 1)

	_ = wrapper.GetCurrentPhase()
	require.Len(t, phaseItemMock.GetCurrentPhaseCalls(), 1)

	wrapper.SetCurrentPhase("phase")
	require.Len(t, phaseItemMock.SetCurrentPhaseCalls(), 1)

	_ = wrapper.GetVersion()
	require.Len(t, phaseItemMock.GetVersionCalls(), 1)

	_ = wrapper.GetSpanAttributes()
	require.Len(t, phaseItemMock.GetSpanAttributesCalls(), 1)

	_ = wrapper.GetSpanKey("phase")
	require.Len(t, phaseItemMock.GetSpanKeyCalls(), 1)

	_ = wrapper.GetSpanName("phase")
	require.Len(t, phaseItemMock.GetSpanNameCalls(), 1)

	wrapper.Complete()
	require.Len(t, phaseItemMock.CompleteCalls(), 1)

	_ = wrapper.IsEndTimeSet()
	require.Len(t, phaseItemMock.IsEndTimeSetCalls(), 1)

	_ = wrapper.GetEndTime()
	require.Len(t, phaseItemMock.GetEndTimeCalls(), 1)

	_ = wrapper.GetStartTime()
	require.Len(t, phaseItemMock.GetStartTimeCalls(), 1)

	_ = wrapper.GetPreviousVersion()
	require.Len(t, phaseItemMock.GetPreviousVersionCalls(), 1)

	_ = wrapper.GetParentName()
	require.Len(t, phaseItemMock.GetParentNameCalls(), 1)

	_ = wrapper.GetNamespace()
	require.Len(t, phaseItemMock.GetNamespaceCalls(), 1)

	_ = wrapper.GetAppName()
	require.Len(t, phaseItemMock.GetAppNameCalls(), 1)

	_ = wrapper.GetPreDeploymentTasks()
	require.Len(t, phaseItemMock.GetPreDeploymentTasksCalls(), 1)

	_ = wrapper.GetPostDeploymentTasks()
	require.Len(t, phaseItemMock.GetPostDeploymentTasksCalls(), 1)

	_ = wrapper.GetPreDeploymentTaskStatus()
	require.Len(t, phaseItemMock.GetPreDeploymentTaskStatusCalls(), 1)

	_ = wrapper.GetPostDeploymentTaskStatus()
	require.Len(t, phaseItemMock.GetPostDeploymentTaskStatusCalls(), 1)

	_ = wrapper.GetPreDeploymentEvaluations()
	require.Len(t, phaseItemMock.GetPreDeploymentEvaluationsCalls(), 1)

	_ = wrapper.GetPostDeploymentEvaluations()
	require.Len(t, phaseItemMock.GetPostDeploymentEvaluationsCalls(), 1)

	_ = wrapper.GetPreDeploymentEvaluationTaskStatus()
	require.Len(t, phaseItemMock.GetPreDeploymentEvaluationTaskStatusCalls(), 1)

	_ = wrapper.GetPostDeploymentEvaluationTaskStatus()
	require.Len(t, phaseItemMock.GetPostDeploymentEvaluationTaskStatusCalls(), 1)

	_ = wrapper.GenerateTask(nil, "", common.PostDeploymentCheckType)
	require.Len(t, phaseItemMock.GenerateTaskCalls(), 1)

	_ = wrapper.GenerateEvaluation(nil, "", common.PostDeploymentCheckType)
	require.Len(t, phaseItemMock.GenerateEvaluationCalls(), 1)

	wrapper.SetSpanAttributes(nil)
	require.Len(t, phaseItemMock.SetSpanAttributesCalls(), 1)

	wrapper.CancelRemainingPhases(common.PhaseAppDeployment)
	require.Len(t, phaseItemMock.CancelRemainingPhasesCalls(), 1)

	wrapper.SetPhaseTraceID(common.PhaseAppDeployment.LongName, propagation.MapCarrier{})
	require.Len(t, phaseItemMock.SetPhaseTraceIDCalls(), 1)

}
