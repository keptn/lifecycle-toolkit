package manager

import (
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type Provider interface {
	SetupWebhookServer(mgr manager.Manager)
}
