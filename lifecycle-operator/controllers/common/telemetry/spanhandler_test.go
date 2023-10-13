package telemetry

import (
	"context"
	"testing"

	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha4"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestSpanHandler_GetAndUnbindSpan_WorkloadVersion(t *testing.T) {
	wi := &v1alpha4.KeptnWorkloadVersion{}
	wi.Spec.TraceId = make(map[string]string, 1)
	wi.Spec.TraceId["test"] = "test"
	wi.Spec.AppName = "test"
	wi.Spec.WorkloadName = "test"
	wi.Spec.Version = "test"
	doAssert(t, wi)
}

func TestSpanHandler_GetAndUnbindSpan_AppVersion(t *testing.T) {
	av := &v1alpha3.KeptnAppVersion{}
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

func TestSpanHandler_GetSpan(t *testing.T) {
	wi := &v1alpha4.KeptnWorkloadVersion{}
	wi.Spec.TraceId = make(map[string]string, 1)
	wi.Spec.TraceId["traceparent"] = "test-parent"
	wi.Spec.AppName = "test"
	wi.Spec.WorkloadName = "test"
	wi.Spec.Version = "test"

	r := SpanHandler{}
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

	wi2 := &v1alpha4.KeptnWorkloadVersion{}
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
