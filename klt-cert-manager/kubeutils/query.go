package kubeutils

import (
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type kubeQuery struct {
	kubeClient client.Client
	kubeReader client.Reader
	log        logr.Logger
}

func newKubeQuery(kubeClient client.Client, kubeReader client.Reader, log logr.Logger) kubeQuery {
	return kubeQuery{
		kubeClient: kubeClient,
		kubeReader: kubeReader,
		log:        log,
	}
}
