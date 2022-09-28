package common

import "go.opentelemetry.io/otel/metric/instrument/syncfloat64"

const WorkloadAnnotation = "keptn.sh/workload"
const VersionAnnotation = "keptn.sh/version"
const AppAnnotation = "keptn.sh/app"
const EventAnnotation = "keptn.sh/event"

type KeptnMeters struct {
	DeploymentCount syncfloat64.Counter
}
