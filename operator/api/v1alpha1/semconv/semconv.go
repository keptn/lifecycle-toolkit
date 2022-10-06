package semconv

import (
	"github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1"
	"github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1/common"
	"go.opentelemetry.io/otel/trace"
)

func AddAttributeFromWorkload(s trace.Span, w v1alpha1.KeptnWorkload) {
	s.SetAttributes(common.ApplicationName.String(w.Spec.AppName))
	s.SetAttributes(common.Workload.String(w.Name))
	s.SetAttributes(common.Version.String(w.Spec.Version))
}

func AddAttributeFromAnnotations(s trace.Span, annotations map[string]string) {
	s.SetAttributes(common.ApplicationName.String(annotations[common.AppAnnotation]))
	s.SetAttributes(common.Workload.String(annotations[common.WorkloadAnnotation]))
	s.SetAttributes(common.Version.String(annotations[common.VersionAnnotation]))
}
