package v1

import (
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *KeptnApp) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}
