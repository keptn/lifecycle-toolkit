package context

import (
	"context"
	"reflect"
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

func TestGetAppMetadataFromContext(t *testing.T) {

	tests := []struct {
		name   string
		ctx    context.Context
		want   map[string]string
		exists bool
	}{
		{
			name:   "empty context",
			ctx:    context.Background(),
			want:   make(map[string]string),
			exists: false,
		},
		{
			name:   "context with metadata",
			ctx:    context.WithValue(context.TODO(), keptnAppContextKey, map[string]string{"testy": "test"}),
			want:   map[string]string{"testy": "test"},
			exists: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metadata, exist := GetAppMetadataFromContext(tt.ctx)
			if !reflect.DeepEqual(metadata, tt.want) {
				t.Errorf("GetAppMetadataFromContext() got = %v, want %v", metadata, tt.want)
			}
			if exist != tt.exists {
				t.Errorf("GetAppMetadataFromContext() got1 = %v, want %v", exist, tt.exists)
			}
		})
	}
}
