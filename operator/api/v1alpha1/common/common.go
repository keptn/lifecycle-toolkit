package common

import "fmt"

const WorkloadAnnotation = "keptn.sh/workload"
const VersionAnnotation = "keptn.sh/version"
const AppAnnotation = "keptn.sh/app"
const PreDeploymentTaskAnnotation = "keptn.sh/pre-deployment-tasks"
const PostDeploymentTaskAnnotation = "keptn.sh/post-deployment-tasks"
const PreDeploymentAnalysisAnnotation = "keptn.sh/pre-deployment-analysis"
const PostDeploymentAnalysisAnnotation = "keptn.sh/post-deployment-analysis"

type KeptnState string

const (
	StateRunning   KeptnState = "Running"
	StateSucceeded KeptnState = "Succeeded"
	StateFailed    KeptnState = "Failed"
	StateUnknown   KeptnState = "Unknown"
	StatePending   KeptnState = "Pending"
)

var ErrTooLongAnnotationsErr = fmt.Errorf("Too long annotations, maximum length for app and workload is 25 characters, for version 12 characters")
