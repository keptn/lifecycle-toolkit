package interfaces

import (
	"testing"
	"time"

	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/lifecycle/interfaces/fake"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func TestPhaseItemWrapper_GetState(t *testing.T) {
	appVersion := &v1alpha3.KeptnAppVersion{
		Status: v1alpha3.KeptnAppVersionStatus{
			Status:       apicommon.StateFailed,
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
		GetStateFunc: func() apicommon.KeptnState {
			return apicommon.StatePending
		},
		SetStateFunc: func(keptnState apicommon.KeptnState) {
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
		GetPreDeploymentTaskStatusFunc: func() []v1alpha3.ItemStatus {
			return nil
		},
		GetPostDeploymentTaskStatusFunc: func() []v1alpha3.ItemStatus {
			return nil
		},
		GetPreDeploymentEvaluationsFunc: func() []string {
			return nil
		},
		GetPostDeploymentEvaluationsFunc: func() []string {
			return nil
		},
		GetPreDeploymentEvaluationTaskStatusFunc: func() []v1alpha3.ItemStatus {
			return nil
		},
		GetPostDeploymentEvaluationTaskStatusFunc: func() []v1alpha3.ItemStatus {
			return nil
		},
		GenerateTaskFunc: func(taskDefinition v1alpha3.KeptnTaskDefinition, checkType apicommon.CheckType) v1alpha3.KeptnTask {
			return v1alpha3.KeptnTask{}
		},
		GenerateEvaluationFunc: func(evaluationDefinition v1alpha3.KeptnEvaluationDefinition, checkType apicommon.CheckType) v1alpha3.KeptnEvaluation {
			return v1alpha3.KeptnEvaluation{}
		},
		SetSpanAttributesFunc: func(span trace.Span) {
		},
		DeprecateRemainingPhasesFunc: func(phase apicommon.KeptnPhaseType) {
		},
	}

	wrapper := PhaseItemWrapper{Obj: &phaseItemMock}

	_ = wrapper.GetState()
	require.Len(t, phaseItemMock.GetStateCalls(), 1)

	wrapper.SetState(apicommon.StateFailed)
	require.Len(t, phaseItemMock.SetStateCalls(), 1)

	_ = wrapper.GetCurrentPhase()
	require.Len(t, phaseItemMock.GetCurrentPhaseCalls(), 1)

	wrapper.SetCurrentPhase("phase")
	require.Len(t, phaseItemMock.SetCurrentPhaseCalls(), 1)

	_ = wrapper.GetVersion()
	require.Len(t, phaseItemMock.GetVersionCalls(), 1)

	_ = wrapper.GetSpanAttributes()
	require.Len(t, phaseItemMock.GetSpanAttributesCalls(), 1)

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

	_ = wrapper.GenerateTask(v1alpha3.KeptnTaskDefinition{}, apicommon.PostDeploymentCheckType)
	require.Len(t, phaseItemMock.GenerateTaskCalls(), 1)

	_ = wrapper.GenerateEvaluation(v1alpha3.KeptnEvaluationDefinition{}, apicommon.PostDeploymentCheckType)
	require.Len(t, phaseItemMock.GenerateEvaluationCalls(), 1)

	wrapper.SetSpanAttributes(nil)
	require.Len(t, phaseItemMock.SetSpanAttributesCalls(), 1)

	wrapper.DeprecateRemainingPhases(apicommon.PhaseAppDeployment)
	require.Len(t, phaseItemMock.DeprecateRemainingPhasesCalls(), 1)

}
