package keptnwebhookcontroller

import (
	"context"
	"fmt"
	"reflect"

	"github.com/go-logr/logr"
	"github.com/keptn/lifecycle-toolkit/klt-cert-manager/eventfilter"
	"github.com/pkg/errors"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apiv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// KeptnWebhookCertificateReconciler reconciles a KeptnWebhookCertificate object
type KeptnWebhookCertificateReconciler struct {
	ctx           context.Context
	Client        client.Client
	Scheme        *runtime.Scheme
	CancelMgrFunc context.CancelFunc
	Log           logr.Logger
}

//clusterrole
// +kubebuilder:rbac:groups=admissionregistration.k8s.io,resources=validatingwebhookconfigurations,verbs=get;list;watch;update;patch;
// +kubebuilder:rbac:groups=admissionregistration.k8s.io,resources=mutatingwebhookconfigurations,verbs=get;list;watch;update;patch;
// +kubebuilder:rbac:groups="apiextensions.k8s.io",resources=customresourcedefinitions,verbs=get;list;watch;update;patch;
// +kubebuilder:rbac:groups="apps",resources=deployments,verbs=get;list;watch;

//role
// +kubebuilder:rbac:groups="",namespace=keptn-lifecycle-toolkit-system,resources=secrets,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *KeptnWebhookCertificateReconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	r.Log.Info("reconciling webhook certificates",
		"namespace", request.Namespace, "name", request.Name)

	r.ctx = ctx

	mutatingWebhookConfiguration, err := r.getMutatingWebhookConfiguration()
	if err != nil {
		r.Log.Error(err, "could not find mutating webhook configuration")
	}

	validatingWebhookConfiguration, err := r.getValidatingWebhookConfiguration()
	if err != nil {
		r.Log.Error(err, "could not find validating webhook configuration")
	}

	crds := &apiv1.CustomResourceDefinitionList{}
	crds, err = r.getCRDConfigurations()
	if err != nil {
		r.Log.Error(err, "could not find CRDs")
	}

	certSecret := newCertificateSecret(r.Client)

	err = certSecret.setSecretFromReader(r.ctx, namespace, r.Log)
	if err != nil {
		r.Log.Error(err, "could not get secret")
		return reconcile.Result{}, errors.WithStack(err)
	}

	err = certSecret.setCertificates(namespace)
	if err != nil {
		r.Log.Error(err, "could not validate certificate")
		return reconcile.Result{}, errors.WithStack(err)
	}

	mutatingWebhookConfigs := getClientConfigsFromMutatingWebhook(mutatingWebhookConfiguration)

	validatingWebhookConfigs := getClientConfigsFromValidatingWebhook(validatingWebhookConfiguration)

	areMutatingWebhookConfigsValid := certSecret.areWebhookConfigsValid(mutatingWebhookConfigs)
	areValidatingWebhookConfigsValid := certSecret.areWebhookConfigsValid(validatingWebhookConfigs)
	areCRDConversionsConfigValid := certSecret.areCRDConversionsValid(crds)
	isCertSecretRecent := certSecret.isRecent()

	if isCertSecretRecent && areMutatingWebhookConfigsValid && areValidatingWebhookConfigsValid && areCRDConversionsConfigValid {
		r.Log.Info("secret for certificates up to date, skipping update")
		r.cancelMgr()
		return reconcile.Result{RequeueAfter: SuccessDuration}, nil
	}

	if err = certSecret.createOrUpdateIfNecessary(r.ctx); err != nil {
		return reconcile.Result{}, errors.WithStack(err)
	}

	bundle, err := certSecret.loadCombinedBundle()
	if err != nil {
		return reconcile.Result{}, errors.WithStack(err)
	}

	if err := r.updateClientConfigurations(bundle, mutatingWebhookConfigs, mutatingWebhookConfiguration); err != nil {
		return reconcile.Result{}, errors.WithStack(err)
	}

	if err := r.updateClientConfigurations(bundle, validatingWebhookConfigs, validatingWebhookConfiguration); err != nil {
		return reconcile.Result{}, errors.WithStack(err)
	}

	if err = r.updateCRDsConfiguration(crds, bundle); err != nil {
		return reconcile.Result{}, errors.WithStack(err)
	}

	r.cancelMgr()
	return reconcile.Result{RequeueAfter: SuccessDuration}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KeptnWebhookCertificateReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.Deployment{}).
		WithEventFilter(eventfilter.ForObjectNameAndNamespace(DeploymentName, namespace)).
		Owns(&corev1.Secret{}).
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
	if err := r.Client.Get(r.ctx, client.ObjectKey{
		Name: MutatingWebhookconfig,
	}, &mutatingWebhook); err != nil {
		return nil, err
	}

	if len(mutatingWebhook.Webhooks) <= 0 {
		return nil, errors.New("mutating webhook configuration has no registered webhooks")
	}
	return &mutatingWebhook, nil
}

func (r *KeptnWebhookCertificateReconciler) getValidatingWebhookConfiguration() (
	*admissionregistrationv1.ValidatingWebhookConfiguration, error) {
	var validatingWebhook admissionregistrationv1.ValidatingWebhookConfiguration
	if err := r.Client.Get(r.ctx, client.ObjectKey{
		Name: ValidatingWebhookconfig,
	}, &validatingWebhook); err != nil {
		return nil, err
	}

	if len(validatingWebhook.Webhooks) <= 0 {
		return nil, errors.New("validating webhook configuration has no registered webhooks")
	}
	return &validatingWebhook, nil
}

func (r *KeptnWebhookCertificateReconciler) updateClientConfigurations(bundle []byte,
	webhookClientConfigs []*admissionregistrationv1.WebhookClientConfig, webhookConfig client.Object) error {
	if webhookConfig == nil || reflect.ValueOf(webhookConfig).IsNil() {
		return nil
	}

	for i := range webhookClientConfigs {
		webhookClientConfigs[i].CABundle = bundle
	}

	if err := r.Client.Update(r.ctx, webhookConfig); err != nil {
		return err
	}
	return nil
}

func (r *KeptnWebhookCertificateReconciler) getCRDConfigurations() (
	*apiv1.CustomResourceDefinitionList, error) {
	var crds apiv1.CustomResourceDefinitionList
	opt := client.MatchingLabels{
		"crdGroup": crdGroup,
	}
	if err := r.Client.List(r.ctx, &crds, opt); err != nil {
		return nil, err
	}

	return &crds, nil
}

func (r *KeptnWebhookCertificateReconciler) updateCRDsConfiguration(crds *apiv1.CustomResourceDefinitionList, bundle []byte) error {
	fail := false
	for _, crd := range crds.Items {
		if err := r.updateCRDConfiguration(crd.Name, bundle); err != nil {
			fail = true
		}

	}
	if fail {
		return fmt.Errorf(couldNotUpdateCRDErr)
	}
	return nil
}

func (r *KeptnWebhookCertificateReconciler) updateCRDConfiguration(crdName string, bundle []byte) error {
	var crd apiv1.CustomResourceDefinition
	if err := r.Client.Get(r.ctx, types.NamespacedName{Name: crdName}, &crd); err != nil {
		return err
	}

	if !hasConversionWebhook(crd) {
		r.Log.Info(fmt.Sprintf("no conversion webhook config for %s, no cert will be provided", crdName))
		return nil
	}

	// update crd
	crd.Spec.Conversion.Webhook.ClientConfig.CABundle = bundle
	if err := r.Client.Update(r.ctx, &crd); err != nil {
		return err
	}
	return nil
}

func hasConversionWebhook(crd apiv1.CustomResourceDefinition) bool {
	return crd.Spec.Conversion != nil && crd.Spec.Conversion.Webhook != nil && crd.Spec.Conversion.Webhook.ClientConfig != nil
}
