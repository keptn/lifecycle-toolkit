package context

import (
	"context"
)

const keptnAppContextKey = "keptnAppContextMeta"

func ContextWithAppMetadata(ctx context.Context, appContextMeta map[string]string) context.Context {
	return context.WithValue(ctx, keptnAppContextKey, appContextMeta)
}

func GetAppMetadataFromContext(ctx context.Context) (map[string]string, bool) {
	value := ctx.Value(keptnAppContextKey)

	appContextMeta, ok := value.(map[string]string)
	if ok {
		return appContextMeta, true
	}
	return map[string]string{}, false
}
