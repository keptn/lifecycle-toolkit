package context

import (
	"context"
)

type keptnAppContextKeyType string

var keptnAppContextKey = keptnAppContextKeyType("keptnAppContextMeta")

func WithAppMetadata(ctx context.Context, appContextMeta ...map[string]string) context.Context {
	mergedMap := map[string]string{}
	for _, meta := range appContextMeta {
		for key, value := range meta {
			mergedMap[key] = value
		}
	}
	return context.WithValue(ctx, keptnAppContextKey, mergedMap)
}

func GetAppMetadataFromContext(ctx context.Context) (map[string]string, bool) {
	value := ctx.Value(keptnAppContextKey)

	appContextMeta, ok := value.(map[string]string)
	if ok {
		return appContextMeta, true
	}
	return map[string]string{}, false
}
