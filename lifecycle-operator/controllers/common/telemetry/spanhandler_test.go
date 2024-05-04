package telemetry

import (
	"context"
	"testing"

	apilifecycle "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1/common"
	keptncontext "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/context"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestSpanHandler_GetAndUnbindSpan_WorkloadVersion(t *testing.T) {
	wi := &apilifecycle.KeptnWorkloadVersion{}
	wi.Spec.TraceId = make(map[string]string, 1)
	wi.Spec.TraceId["test"] = "test"
	wi.Spec.AppName = "test"
	wi.Spec.WorkloadName = "test"
	wi.Spec.Version = "test"
	doAssert(t, wi)
}

func TestSpanHandler_GetAndUnbindSpan_AppVersion(t *testing.T) {
	av := &apilifecycle.KeptnAppVersion{}
	av.Spec.TraceId = make(map[string]string, 1)
	av.Spec.TraceId["test"] = "test"
	av.Spec.AppName = "test"
	av.Spec.Version = "test"
	doAssert(t, av)
}

func doAssert(t *testing.T, obj client.Object) {
	r := Handler{}
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

func TestSpanHandler_GetSpan(t *testing.T) {
	wi := &apilifecycle.KeptnWorkloadVersion{}
	wi.Spec.TraceId = make(map[string]string, 1)
	wi.Spec.TraceId["traceparent"] = "test-parent"
	wi.Spec.AppName = "test"
	wi.Spec.WorkloadName = "test"
	wi.Spec.Version = "test"

	r := Handler{}
	phase := apicommon.PhaseAppDeployment.ShortName
	tracer := otel.Tracer("keptn/test")

	ctx, span, err := r.GetSpan(context.TODO(), tracer, wi, phase)

	require.Nil(t, err)
	require.NotNil(t, span)
	require.NotNil(t, ctx)

	ctx2, span2, err2 := r.GetSpan(context.TODO(), tracer, wi, phase)

	require.Nil(t, err2)
	require.Equal(t, ctx, ctx2)
	require.Equal(t, span, span2)

	wi2 := &apilifecycle.KeptnWorkloadVersion{}
	wi2.Spec.TraceId = make(map[string]string, 1)
	wi2.Spec.TraceId["traceparent"] = "test-parent2"
	wi2.Spec.AppName = "test2"
	wi2.Spec.WorkloadName = "test2"
	wi2.Spec.Version = "test2"

	tracer2 := otel.Tracer("keptn/test2")
	phase2 := apicommon.PhaseWorkloadPreDeployment.LongName

	ctx3, span3, err3 := r.GetSpan(context.TODO(), tracer2, wi2, phase2)

	require.Nil(t, err3)
	require.NotNil(t, span3)
	require.NotNil(t, ctx3)
	require.NotEqual(t, ctx, ctx3)
	require.NotEqual(t, span, span3)

	ctx4, span4, err4 := r.GetSpan(context.TODO(), tracer2, wi2, phase2)

	require.Nil(t, err4)
	require.Equal(t, ctx3, ctx4)
	require.Equal(t, span3, span4)

	ctx5, span5, err5 := r.GetSpan(context.TODO(), tracer, wi, phase)

	require.Nil(t, err5)
	require.Equal(t, ctx, ctx5)
	require.Equal(t, span, span5)

}

func TestSpanHandler_GetSpanWithAttributes(t *testing.T) {
	wi := &apilifecycle.KeptnWorkloadVersion{}
	wi.Spec.TraceId = make(map[string]string, 1)
	wi.Spec.TraceId["traceparent"] = "test-parent"
	wi.Spec.AppName = "test"
	wi.Spec.WorkloadName = "test"
	wi.Spec.Version = "test"

	r := Handler{}
	phase := apicommon.PhaseAppDeployment.ShortName

	tp := trace.NewTracerProvider()

	tracer := tp.Tracer("keptn")

	ctx := context.TODO()

	ctx = keptncontext.WithAppMetadata(ctx, map[string]string{"foo": "bar"})
	ctx, span, err := r.GetSpan(ctx, tracer, wi, phase)

	require.Nil(t, err)
	require.NotNil(t, span)
	require.NotNil(t, ctx)

	attributes := span.(trace.ReadOnlySpan).Attributes()
	require.NotNil(t, attributes)

	// the total number of attributes should be 5 (i.e. the workload specific ones + the additional one)
	require.Len(t, attributes, 5)
	require.Equal(t, "bar", attributes[4].Value.AsString())
}
