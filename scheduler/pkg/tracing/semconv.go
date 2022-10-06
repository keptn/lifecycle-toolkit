package tracing

import (
	"context"
	"github.com/keptn-sandbox/lifecycle-controller/scheduler/pkg/klcpermit"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func CreateSpan(ctx context.Context, crd unstructured.Unstructured, w klcpermit.WorkloadManager) (context.Context, trace.Span) {
	// search for annotations
	annotations, found, _ := unstructured.NestedMap(crd.UnstructuredContent(), "metadata", "annotations")
	if found {
		ctx = otel.GetTextMapPropagator().Extract(ctx, KeptnCarrier(annotations))
	}
	ctx, span := w.Tracer.Start(ctx, "schedule")

	// TODO extract from unstructured data Keptn Semantic Conventions

	return ctx, span
}
