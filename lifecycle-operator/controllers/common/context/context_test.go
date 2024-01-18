package context

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContextWithAppMetadata(t *testing.T) {
	ctx := context.Background()

	metadata := map[string]string{
		"foo": "bar",
	}

	ctx = WithAppMetadata(ctx, metadata)

	require.NotNil(t, ctx)

	metadata, ok := GetAppMetadataFromContext(ctx)

	require.True(t, ok)
	require.Equal(t, "bar", metadata["foo"])
}
