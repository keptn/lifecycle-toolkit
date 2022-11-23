package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	ApplicationName attribute.Key = attribute.Key("keptn.deployment.app_name")
	Workload        attribute.Key = attribute.Key("keptn.deployment.workload")
	Version         attribute.Key = attribute.Key("keptn.deployment.version")
	Namespace       attribute.Key = attribute.Key("keptn.deployment.namespace")
	Status          attribute.Key = attribute.Key("keptn.deployment.status")
)

func CreateSpan(ctx context.Context, crd *unstructured.Unstructured, t trace.Tracer, ns string) (context.Context, trace.Span) {
	// search for annotations
	annotations, found, _ := unstructured.NestedMap(crd.UnstructuredContent(), "metadata", "annotations")
	if found {
		ctx = otel.GetTextMapPropagator().Extract(ctx, KeptnCarrier(annotations))
	}
	ctx, span := t.Start(ctx, "schedule")

	appName, found, _ := unstructured.NestedString(crd.UnstructuredContent(), "spec", "app")
	if found {
		span.SetAttributes(ApplicationName.String(appName))
	}
	workload, found, _ := unstructured.NestedString(crd.UnstructuredContent(), "spec", "workloadName")
	if found {
		span.SetAttributes(Workload.String(workload))
	}
	version, found, _ := unstructured.NestedString(crd.UnstructuredContent(), "spec", "version")
	if found {
		span.SetAttributes(Version.String(version))
	}
	span.SetAttributes(Namespace.String(ns))
	return ctx, span
}
