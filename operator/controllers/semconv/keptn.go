package semconv

import (
	klcv1alpha1 "github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1"
	"github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1/common"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	ApplicationName attribute.Key = attribute.Key("keptn.deployment.app_name")
	Workload        attribute.Key = attribute.Key("keptn.deployment.workload")
	Version         attribute.Key = attribute.Key("keptn.deployment.version")
)

func AddAttributeFromWorkload(s trace.Span, w klcv1alpha1.KeptnWorkload) {
	s.SetAttributes(ApplicationName.String(w.Spec.AppName))
	s.SetAttributes(Workload.String(w.Name))
	s.SetAttributes(Version.String(w.Spec.Version))
}

func AddAttributeFromAnnotations(s trace.Span, annotations map[string]string) {
	s.SetAttributes(ApplicationName.String(annotations[common.AppAnnotation]))
	s.SetAttributes(Workload.String(annotations[common.WorkloadAnnotation]))
	s.SetAttributes(Version.String(annotations[common.VersionAnnotation]))
}
