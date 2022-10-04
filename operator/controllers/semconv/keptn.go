package semconv

import "go.opentelemetry.io/otel/attribute"

const (
	ApplicationName attribute.Key = attribute.Key("keptn.deployment.app_name")
	Workload        attribute.Key = attribute.Key("keptn.deployment.workload")
	Version         attribute.Key = attribute.Key("keptn.deployment.version")
)
