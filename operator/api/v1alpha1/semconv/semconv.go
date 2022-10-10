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

func AddAttributeFromWorkloadInstance(s trace.Span, w v1alpha1.KeptnWorkloadInstance) {
	s.SetAttributes(common.ApplicationName.String(w.Spec.AppName))
	s.SetAttributes(common.Workload.String(w.Name))
	s.SetAttributes(common.Version.String(w.Spec.Version))
}

func AddAttributeFromApp(s trace.Span, a v1alpha1.KeptnApp) {
	s.SetAttributes(common.ApplicationName.String(a.Name))
	s.SetAttributes(common.Version.String(a.Spec.Version))
}

func AddAttributeFromTask(s trace.Span, t v1alpha1.KeptnTask) {
	s.SetAttributes(common.ApplicationName.String(t.Spec.AppName))
	s.SetAttributes(common.Workload.String(t.Spec.Workload))
	s.SetAttributes(common.Version.String(t.Spec.Version))
	s.SetAttributes(common.TaskName.String(t.Name))
	s.SetAttributes(common.TaskType.String(string(t.Spec.Type)))
}

func AddAttributeFromAnnotations(s trace.Span, annotations map[string]string) {
	s.SetAttributes(common.ApplicationName.String(annotations[common.AppAnnotation]))
	s.SetAttributes(common.Workload.String(annotations[common.WorkloadAnnotation]))
	s.SetAttributes(common.Version.String(annotations[common.VersionAnnotation]))
}
