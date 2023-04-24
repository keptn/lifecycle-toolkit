package keptnwebhookcontroller

import (
	"context"
)

type IResourceRetriever interface {
	GetMutatingWebhooks(ctx context.Context)
	GetValidatingWebhooks(ctx context.Context)
	GetCRDs(ctx context.Context)
}
