package context

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestContextWithAppMetadata(t *testing.T) {
	ctx := context.Background()

	metadata := map[string]string{
		"foo": "bar",
	}

	ctx = ContextWithAppMetadata(ctx, metadata)

	require.NotNil(t, ctx)

	metadata, ok := GetAppMetadataFromContext(ctx)

	require.True(t, ok)
	require.Equal(t, "bar", metadata["foo"])
}
