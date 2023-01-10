package manager

import (
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

//go:generate moq -pkg fake -skip-ensure -out ../fake/manager_mock.go . IManager:MockManager
type IManager manager.Manager

//go:generate moq -pkg fake -skip-ensure -out ../fake/webhookmanager_mock.go . Provider:MockWebhookManager
type Provider interface {
	SetupWebhookServer(mgr manager.Manager)
}
