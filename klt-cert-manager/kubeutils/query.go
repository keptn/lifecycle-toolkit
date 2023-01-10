package kubeutils

import (
	"context"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type kubeQuery struct {
	kubeClient client.Client
	kubeReader client.Reader
	ctx        context.Context
	log        logr.Logger
}

func newKubeQuery(ctx context.Context, kubeClient client.Client, kubeReader client.Reader, log logr.Logger) kubeQuery {
	return kubeQuery{
		kubeClient: kubeClient,
		kubeReader: kubeReader,
		ctx:        ctx,
		log:        log,
	}
}
