package common

import (
	"fmt"

	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
)

const WorkloadAnnotation = "keptn.sh/workload"
const VersionAnnotation = "keptn.sh/version"
const AppAnnotation = "keptn.sh/app"
const PreDeploymentTaskAnnotation = "keptn.sh/pre-deployment-tasks"
const PostDeploymentTaskAnnotation = "keptn.sh/post-deployment-tasks"
const PreDeploymentAnalysisAnnotation = "keptn.sh/pre-deployment-analysis"
const PostDeploymentAnalysisAnnotation = "keptn.sh/post-deployment-analysis"
const TaskNameAnnotation = "keptn.sh/task-name"

const MaxAppNameLength = 25
const MaxWorkloadNameLength = 25
const MaxTaskNameLength = 25
const MaxVersionLength = 12

type KeptnState string

const (
	StateRunning   KeptnState = "Running"
	StateSucceeded KeptnState = "Succeeded"
	StateFailed    KeptnState = "Failed"
	StateUnknown   KeptnState = "Unknown"
	StatePending   KeptnState = "Pending"
)

var ErrTooLongAnnotations = fmt.Errorf("too long annotations, maximum length for app and workload is 25 characters, for version 12 characters")

func (k KeptnState) IsCompleted() bool {
	return k == StateSucceeded || k == StateFailed || k == StateUnknown
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

type KeptnMeters struct {
	DeploymentCount syncfloat64.Counter
}
