package common

import (
	"fmt"
	"math/rand"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
)

const WorkloadAnnotation = "keptn.sh/workload"
const VersionAnnotation = "keptn.sh/version"
const AppAnnotation = "keptn.sh/app"
const PreDeploymentTaskAnnotation = "keptn.sh/pre-deployment-tasks"
const PostDeploymentTaskAnnotation = "keptn.sh/post-deployment-tasks"
const K8sRecommendedWorkloadAnnotations = "app.kubernetes.io/name"
const K8sRecommendedVersionAnnotations = "app.kubernetes.io/version"
const K8sRecommendedAppAnnotations = "app.kubernetes.io/part-of"
const PreDeploymentEvaluationAnnotation = "keptn.sh/pre-deployment-evaluations"
const PostDeploymentEvaluationAnnotation = "keptn.sh/post-deployment-evaluations"
const TaskNameAnnotation = "keptn.sh/task-name"
const NamespaceEnabledAnnotation = "keptn.sh/lifecycle-toolkit"
const CreateAppTaskSpanName = "create_%s_app_task"
const CreateWorkloadTaskSpanName = "create_%s_deployment_task"
const CreateAppEvalSpanName = "create_%s_app_evaluation"
const CreateWorkloadEvalSpanName = "create_%s_deployment_evaluation"

const MaxAppNameLength = 25
const MaxWorkloadNameLength = 25
const MaxTaskNameLength = 25
const MaxVersionLength = 12

type KeptnState string

const (
	StateProgressing KeptnState = "Progressing"
	StateSucceeded   KeptnState = "Succeeded"
	StateFailed      KeptnState = "Failed"
	StateUnknown     KeptnState = "Unknown"
	StatePending     KeptnState = "Pending"
	StateDeprecated  KeptnState = "Deprecated"
	// StateCancelled represents state that was cancelled due to a previous step having failed.
	// Deprecated: Use StateDeprecated instead. Should only be used in checks for backwards compatibility reasons
	StateCancelled KeptnState = "Cancelled"
)

func (k KeptnState) IsCompleted() bool {
	return k == StateSucceeded || k == StateFailed || k == StateDeprecated || k == StateCancelled
}

func (k KeptnState) IsSucceeded() bool {
	return k == StateSucceeded
}

func (k KeptnState) IsFailed() bool {
	return k == StateFailed
}

func (k KeptnState) IsDeprecated() bool {
	return k == StateDeprecated || k == StateCancelled
}

func (k KeptnState) IsPending() bool {
	return k == StatePending
}

type StatusSummary struct {
	Total       int
	Progressing int
	Failed      int
	Succeeded   int
	Pending     int
	Unknown     int
	Deprecated  int
}

func UpdateStatusSummary(status KeptnState, summary StatusSummary) StatusSummary {
	switch status {
	case StateFailed:
		summary.Failed++
	case StateDeprecated:
		summary.Deprecated++
	case StateSucceeded:
		summary.Succeeded++
	case StateProgressing:
		summary.Progressing++
	case StatePending, "":
		summary.Pending++
	case StateUnknown:
		summary.Unknown++
	}
	return summary
}

func (s StatusSummary) GetTotalCount() int {
	return s.Failed + s.Succeeded + s.Progressing + s.Pending + s.Unknown + s.Deprecated
}

func GetOverallState(s StatusSummary) KeptnState {
	if s.Failed > 0 || s.Deprecated > 0 {
		return StateFailed
	}
	if s.Progressing > 0 {
		return StateProgressing
	}
	if s.Pending > 0 {
		return StatePending
	}
	if s.Unknown > 0 || s.GetTotalCount() != s.Total {
		return StateUnknown
	}
	return StateSucceeded
}

func TruncateString(s string, max int) string {
	if len(s) > max {
		return s[:max]
	}
	return s
}

type CheckType string

const PreDeploymentCheckType CheckType = "pre"
const PostDeploymentCheckType CheckType = "post"
const PreDeploymentEvaluationCheckType CheckType = "pre-eval"
const PostDeploymentEvaluationCheckType CheckType = "post-eval"

type KeptnMeters struct {
	TaskCount          syncint64.Counter
	TaskDuration       syncfloat64.Histogram
	DeploymentCount    syncint64.Counter
	DeploymentDuration syncfloat64.Histogram
	AppCount           syncint64.Counter
	AppDuration        syncfloat64.Histogram
	EvaluationCount    syncint64.Counter
	EvaluationDuration syncfloat64.Histogram
}

const (
	AppName                 attribute.Key = attribute.Key("keptn.deployment.app.name")
	AppVersion              attribute.Key = attribute.Key("keptn.deployment.app.version")
	AppNamespace            attribute.Key = attribute.Key("keptn.deployment.app.namespace")
	AppStatus               attribute.Key = attribute.Key("keptn.deployment.app.status")
	AppPreviousVersion      attribute.Key = attribute.Key("keptn.deployment.app.previousversion")
	WorkloadName            attribute.Key = attribute.Key("keptn.deployment.workload.name")
	WorkloadVersion         attribute.Key = attribute.Key("keptn.deployment.workload.version")
	WorkloadPreviousVersion attribute.Key = attribute.Key("keptn.deployment.workload.previousversion")
	WorkloadNamespace       attribute.Key = attribute.Key("keptn.deployment.workload.namespace")
	WorkloadStatus          attribute.Key = attribute.Key("keptn.deployment.workload.status")
	TaskStatus              attribute.Key = attribute.Key("keptn.deployment.task.status")
	TaskName                attribute.Key = attribute.Key("keptn.deployment.task.name")
	TaskType                attribute.Key = attribute.Key("keptn.deployment.task.type")
	EvaluationStatus        attribute.Key = attribute.Key("keptn.deployment.evaluation.status")
	EvaluationName          attribute.Key = attribute.Key("keptn.deployment.evaluation.name")
	EvaluationType          attribute.Key = attribute.Key("keptn.deployment.evaluation.type")
)

func GenerateTaskName(checkType CheckType, taskName string) string {
	randomId := rand.Intn(99_999-10_000) + 10000
	return fmt.Sprintf("%s-%s-%d", checkType, TruncateString(taskName, 32), randomId)
}

func GenerateEvaluationName(checkType CheckType, evalName string) string {
	randomId := rand.Intn(99_999-10_000) + 10000
	return fmt.Sprintf("%s-%s-%d", checkType, TruncateString(evalName, 27), randomId)
}
