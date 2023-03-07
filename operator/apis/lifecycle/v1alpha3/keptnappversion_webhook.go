package v1alpha3

import (
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *KeptnAppVersion) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}
