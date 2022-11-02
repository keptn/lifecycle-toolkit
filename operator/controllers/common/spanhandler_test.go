package common

import (
	"context"
	"github.com/keptn/lifecycle-controller/operator/api/v1alpha1"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"testing"
)

func TestSpanHandler_GetAndUnbindSpan(t *testing.T) {
	r := SpanHandler{}

	tracer := otel.Tracer("keptn/operator/workloadinstance")

	wi := &v1alpha1.KeptnWorkloadInstance{}
	ctx, span, err := r.GetSpan(context.TODO(), tracer, wi, "pre")

	require.Nil(t, err)
	require.NotNil(t, t, span)
	require.NotNil(t, ctx)

	require.Len(t, r.bindCRDSpan, 1)

	err = r.UnbindSpan(wi, "pre")

	require.Nil(t, err)

	require.Empty(t, r.bindCRDSpan)
}
