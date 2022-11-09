package common

import (
	"context"
	"testing"

	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestSpanHandler_GetAndUnbindSpan_WorkloadInstance(t *testing.T) {
	wi := &v1alpha1.KeptnWorkloadInstance{}
	wi.Spec.TraceId = make(map[string]string, 1)
	wi.Spec.TraceId["test"] = "test"
	wi.Spec.AppName = "test"
	wi.Spec.WorkloadName = "test"
	wi.Spec.Version = "test"
	doAssert(t, wi)
}

func TestSpanHandler_GetAndUnbindSpan_AppVersion(t *testing.T) {
	av := &v1alpha1.KeptnAppVersion{}
	av.Spec.TraceId = make(map[string]string, 1)
	av.Spec.TraceId["test"] = "test"
	av.Spec.AppName = "test"
	av.Spec.Version = "test"
	doAssert(t, av)
}

func doAssert(t *testing.T, obj client.Object) {
	r := SpanHandler{}
	phase := "pre"
	tracer := otel.Tracer("keptn/test")
	ctx, span, err := r.GetSpan(context.TODO(), tracer, obj, phase)

	require.Nil(t, err)
	require.NotNil(t, t, span)
	require.NotNil(t, ctx)

	require.Len(t, r.bindCRDSpan, 1)
	err = r.UnbindSpan(obj, phase)

	require.Nil(t, err)

	require.Empty(t, r.bindCRDSpan)
}
