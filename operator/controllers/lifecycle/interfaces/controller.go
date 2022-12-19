package interfaces

import ctrl "sigs.k8s.io/controller-runtime"

type Controller interface {
	SetupWithManager(mgr ctrl.Manager) error
}
