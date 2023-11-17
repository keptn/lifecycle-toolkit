package v1beta1

import (
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *KeptnEvaluationProvider) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}
