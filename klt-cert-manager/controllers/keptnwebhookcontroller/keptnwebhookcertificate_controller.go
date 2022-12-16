package keptnwebhookcontroller

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	"github.com/keptn/lifecycle-toolkit/klt-cert-manager/eventfilter"
	"github.com/pkg/errors"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	SuccessDuration            = 3 * time.Hour
	Webhookconfig              = "klc-mutating-webhook-configuration"
	secretPostfix              = "-certs"
	certificatesSecretEmptyErr = "certificates secret is empty"
	namespace                  = "keptn-lifecycle-toolkit-system"
)

// KeptnWebhookCertificateReconciler reconciles a KeptnWebhookCertificate object
type KeptnWebhookCertificateReconciler struct {
	ctx           context.Context
	Client        client.Client
	Scheme        *runtime.Scheme
	ApiReader     client.Reader
	CancelMgrFunc context.CancelFunc
	Log           logr.Logger
}

// +kubebuilder:rbac:groups=admissionregistration.k8s.io,namespace=keptn-lifecycle-toolkit-system,resourceNames=klc-mutating-webhook-configuration,resources=mutatingwebhookconfigurations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",namespace=keptn-lifecycle-toolkit-system,resources=secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="apps",namespace=keptn-lifecycle-toolkit-system,resources=deployments,verbs=get;list;watch;

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the KeptnWebhookCertificate object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *KeptnWebhookCertificateReconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	r.Log.Info("reconciling webhook certificates",
		"namespace", request.Namespace, "name", request.Name)

	r.ctx = ctx

	_, err := r.getMutatingWebhookConfiguration()
	if err != nil {
		r.Log.Error(err, "could not find mutating webhook configuration")
	}

	r.cancelMgr()
	return reconcile.Result{RequeueAfter: SuccessDuration}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KeptnWebhookCertificateReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.Deployment{}).
		WithEventFilter(eventfilter.ForObjectNameAndNamespace(DeploymentName, namespace)).
		Complete(r)

}

func (r *KeptnWebhookCertificateReconciler) cancelMgr() {
	if r.CancelMgrFunc != nil {
		r.Log.Info("stopping manager after certificates creation")
		r.CancelMgrFunc()
	}
}

func (r *KeptnWebhookCertificateReconciler) getMutatingWebhookConfiguration() (
	*admissionregistrationv1.MutatingWebhookConfiguration, error) {
	var mutatingWebhook admissionregistrationv1.MutatingWebhookConfiguration
	err := r.ApiReader.Get(r.ctx, client.ObjectKey{
		Name: Webhookconfig,
	}, &mutatingWebhook)
	if err != nil {
		return nil, err
	}

	if len(mutatingWebhook.Webhooks) <= 0 {
		return nil, errors.New("mutating webhook configuration has no registered webhooks")
	}
	return &mutatingWebhook, nil
}
