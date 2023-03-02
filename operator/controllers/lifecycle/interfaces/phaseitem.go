package interfaces

import (
	"time"

	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3/common"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// PhaseItem represents an object which has reconciled phases
//
//go:generate moq -pkg fake --skip-ensure -out ./fake/phaseitem_mock.go . PhaseItem
type PhaseItem interface {
	GetState() apicommon.KeptnState
	SetState(apicommon.KeptnState)
	GetCurrentPhase() string
	SetCurrentPhase(string)
	Complete()
	IsEndTimeSet() bool
	GetEndTime() time.Time
	GetStartTime() time.Time
	GetVersion() string
	GetPreviousVersion() string
	GetParentName() string
	GetNamespace() string
	GetAppName() string
	GetPreDeploymentTasks() []string
	GetPostDeploymentTasks() []string
	GetPreDeploymentTaskStatus() []klcv1alpha3.ItemStatus
	GetPostDeploymentTaskStatus() []klcv1alpha3.ItemStatus
	GetPreDeploymentEvaluations() []string
	GetPostDeploymentEvaluations() []string
	GetPreDeploymentEvaluationTaskStatus() []klcv1alpha3.ItemStatus
	GetPostDeploymentEvaluationTaskStatus() []klcv1alpha3.ItemStatus
	GenerateTask(taskDefinition string, checkType apicommon.CheckType) klcv1alpha3.KeptnTask
	GenerateEvaluation(evaluationDefinition string, checkType apicommon.CheckType) klcv1alpha3.KeptnEvaluation
	GetSpanAttributes() []attribute.KeyValue
	SetSpanAttributes(span trace.Span)
	DeprecateRemainingPhases(phase apicommon.KeptnPhaseType)
}

type PhaseItemWrapper struct {
	Obj PhaseItem
}

func NewPhaseItemWrapperFromClientObject(object client.Object) (*PhaseItemWrapper, error) {
	pi, ok := object.(PhaseItem)
	if !ok {
		return nil, errors.ErrCannotWrapToPhaseItem
	}
	return &PhaseItemWrapper{Obj: pi}, nil
}

func (pw PhaseItemWrapper) GetState() apicommon.KeptnState {
	return pw.Obj.GetState()
}

func (pw *PhaseItemWrapper) SetState(state apicommon.KeptnState) {
	pw.Obj.SetState(state)
}

func (pw PhaseItemWrapper) GetCurrentPhase() string {
	return pw.Obj.GetCurrentPhase()
}

func (pw *PhaseItemWrapper) SetCurrentPhase(phase string) {
	pw.Obj.SetCurrentPhase(phase)
}

func (pw PhaseItemWrapper) GetEndTime() time.Time {
	return pw.Obj.GetEndTime()
}

func (pw PhaseItemWrapper) GetStartTime() time.Time {
	return pw.Obj.GetStartTime()
}

func (pw PhaseItemWrapper) IsEndTimeSet() bool {
	return pw.Obj.IsEndTimeSet()
}

func (pw *PhaseItemWrapper) Complete() {
	pw.Obj.Complete()
}

func (pw PhaseItemWrapper) GetVersion() string {
	return pw.Obj.GetVersion()
}

func (pw PhaseItemWrapper) GetPreviousVersion() string {
	return pw.Obj.GetPreviousVersion()
}

func (pw PhaseItemWrapper) GetParentName() string {
	return pw.Obj.GetParentName()
}

func (pw PhaseItemWrapper) GetNamespace() string {
	return pw.Obj.GetNamespace()
}

func (pw PhaseItemWrapper) GetAppName() string {
	return pw.Obj.GetAppName()
}

func (pw PhaseItemWrapper) GetPreDeploymentTasks() []string {
	return pw.Obj.GetPreDeploymentTasks()
}

func (pw PhaseItemWrapper) GetPostDeploymentTasks() []string {
	return pw.Obj.GetPostDeploymentTasks()
}

func (pw PhaseItemWrapper) GetPreDeploymentTaskStatus() []klcv1alpha3.ItemStatus {
	return pw.Obj.GetPreDeploymentTaskStatus()
}

func (pw PhaseItemWrapper) GetPostDeploymentTaskStatus() []klcv1alpha3.ItemStatus {
	return pw.Obj.GetPostDeploymentTaskStatus()
}

func (pw PhaseItemWrapper) GetPreDeploymentEvaluations() []string {
	return pw.Obj.GetPreDeploymentEvaluations()
}

func (pw PhaseItemWrapper) GetPostDeploymentEvaluations() []string {
	return pw.Obj.GetPostDeploymentEvaluations()
}

func (pw PhaseItemWrapper) GetPreDeploymentEvaluationTaskStatus() []klcv1alpha3.ItemStatus {
	return pw.Obj.GetPreDeploymentEvaluationTaskStatus()
}

func (pw PhaseItemWrapper) GetPostDeploymentEvaluationTaskStatus() []klcv1alpha3.ItemStatus {
	return pw.Obj.GetPostDeploymentEvaluationTaskStatus()
}

func (pw PhaseItemWrapper) GenerateTask(taskDefinition string, checkType apicommon.CheckType) klcv1alpha3.KeptnTask {
	return pw.Obj.GenerateTask(taskDefinition, checkType)
}

func (pw PhaseItemWrapper) GenerateEvaluation(evaluationDefinition string, checkType apicommon.CheckType) klcv1alpha3.KeptnEvaluation {
	return pw.Obj.GenerateEvaluation(evaluationDefinition, checkType)
}

func (pw PhaseItemWrapper) SetSpanAttributes(span trace.Span) {
	pw.Obj.SetSpanAttributes(span)
}

func (pw PhaseItemWrapper) GetSpanAttributes() []attribute.KeyValue {
	return pw.Obj.GetSpanAttributes()
}

func (pw PhaseItemWrapper) DeprecateRemainingPhases(phase apicommon.KeptnPhaseType) {
	pw.Obj.DeprecateRemainingPhases(phase)
}
