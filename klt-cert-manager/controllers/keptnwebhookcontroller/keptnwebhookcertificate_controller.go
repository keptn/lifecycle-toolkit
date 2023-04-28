package keptnwebhookcontroller

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"reflect"

	"github.com/go-logr/logr"
	"github.com/keptn/lifecycle-toolkit/klt-cert-manager/eventfilter"
	"github.com/pkg/errors"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// KeptnWebhookCertificateReconciler reconciles a KeptnWebhookCertificate object
type KeptnWebhookCertificateReconciler struct {
	Client        client.Client
	Scheme        *runtime.Scheme
	CancelMgrFunc context.CancelFunc
	Log           logr.Logger
	Namespace     string
	MatchLabels   labels.Set
}

//clusterrole
// +kubebuilder:rbac:groups=admissionregistration.k8s.io,resources=validatingwebhookconfigurations,verbs=get;list;watch;update;patch;
// +kubebuilder:rbac:groups=admissionregistration.k8s.io,resources=mutatingwebhookconfigurations,verbs=get;list;watch;update;patch;
// +kubebuilder:rbac:groups="apiextensions.k8s.io",resources=customresourcedefinitions,verbs=get;list;watch;update;patch;
// +kubebuilder:rbac:groups="apps",resources=deployments,verbs=get;list;watch;

//role
// +kubebuilder:rbac:groups="",namespace=keptn-lifecycle-toolkit-system,resources=secrets,verbs=get;update;patch;delete,resourceNames=klt-certs
// +kubebuilder:rbac:groups="",namespace=keptn-lifecycle-toolkit-system,resources=secrets,verbs=create;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *KeptnWebhookCertificateReconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	r.Log.Info("reconciling webhook certificates",
		"namespace", request.Namespace, "name", request.Name)

	mutatingWebhookConfigurations, err := r.getMutatingWebhookConfigurations(ctx)
	if err != nil {
		r.Log.Error(err, "could not find mutating webhook configuration")
	}

	validatingWebhookConfigurations, err := r.getValidatingWebhookConfigurations(ctx)
	if err != nil {
		r.Log.Error(err, "could not find validating webhook configuration")
	}

	crds, err := r.getCRDConfigurations(ctx)
	if err != nil {
		r.Log.Error(err, "could not find CRDs")
	}

	certSecret := newCertificateSecret(r.Client)

	if err := r.setCertificates(ctx, certSecret); err != nil {
		return reconcile.Result{}, errors.WithStack(err)
	}

	mutatingWebhookConfigs := getClientConfigsFromMutatingWebhook(mutatingWebhookConfigurations)

	validatingWebhookConfigs := getClientConfigsFromValidatingWebhook(validatingWebhookConfigurations)

	areMutatingWebhookConfigsValid := certSecret.areWebhookConfigsValid(mutatingWebhookConfigs)
	areValidatingWebhookConfigsValid := certSecret.areWebhookConfigsValid(validatingWebhookConfigs)
	areCRDConversionsConfigValid := certSecret.areCRDConversionsValid(crds)
	isCertSecretRecent := certSecret.isRecent()

	if isCertSecretRecent && areMutatingWebhookConfigsValid && areValidatingWebhookConfigsValid && areCRDConversionsConfigValid {
		r.Log.Info("secret for certificates up to date, skipping update")
		r.cancelMgr()
		return reconcile.Result{RequeueAfter: successDuration}, nil
	}

	if err = r.updateConfigurations(ctx, certSecret, crds, mutatingWebhookConfigs, mutatingWebhookConfigurations, validatingWebhookConfigs, validatingWebhookConfigurations); err != nil {
		return reconcile.Result{}, errors.WithStack(err)
	}

	r.cancelMgr()
	return reconcile.Result{RequeueAfter: successDuration}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KeptnWebhookCertificateReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.Deployment{}).
		WithEventFilter(eventfilter.ForLabelsAndNamespace(labels.SelectorFromSet(r.MatchLabels), r.Namespace)).
		Owns(&corev1.Secret{}).
		Complete(r)
}

func (r *KeptnWebhookCertificateReconciler) setCertificates(ctx context.Context, certSecret *certificateSecret) error {
	err := certSecret.setSecretFromReader(ctx, r.Namespace, r.Log)
	if err != nil {
		r.Log.Error(err, "could not get secret")
		return err
	}

	err = certSecret.setCertificates(r.Namespace)
	if err != nil {
		r.Log.Error(err, "could not validate certificate")
		return err
	}

	return nil
}

func (r *KeptnWebhookCertificateReconciler) updateConfigurations(ctx context.Context, certSecret *certificateSecret, crds *apiv1.CustomResourceDefinitionList, mutatingWebhookConfigs []*admissionregistrationv1.WebhookClientConfig, mutatingWebhookConfigurationList *admissionregistrationv1.MutatingWebhookConfigurationList, validatingWebhookConfigs []*admissionregistrationv1.WebhookClientConfig, validatingWebhookConfigurationList *admissionregistrationv1.ValidatingWebhookConfigurationList) error {
	if err := certSecret.createOrUpdateIfNecessary(ctx); err != nil {
		return err
	}

	bundle, err := certSecret.loadCombinedBundle()
	if err != nil {
		return err
	}

	for i := range mutatingWebhookConfigurationList.Items {
		r.Log.Info("injecting certificate into mutating webhook config", "mwc", mutatingWebhookConfigurationList.Items[i].Name)
		if err := r.updateClientConfigurations(ctx, bundle, mutatingWebhookConfigs, &mutatingWebhookConfigurationList.Items[i]); err != nil {
			return err
		}
	}

	for i := range validatingWebhookConfigurationList.Items {
		r.Log.Info("injecting certificate into validating webhook config", "vwc", mutatingWebhookConfigurationList.Items[i].Name)
		if err := r.updateClientConfigurations(ctx, bundle, validatingWebhookConfigs, &validatingWebhookConfigurationList.Items[i]); err != nil {
			return err
		}
	}

	if err = r.updateCRDsConfiguration(ctx, crds, bundle); err != nil {
		return err
	}
	return nil
}

func (r *KeptnWebhookCertificateReconciler) cancelMgr() {
	if r.CancelMgrFunc != nil {
		r.Log.Info("stopping manager after certificates creation")
		r.CancelMgrFunc()
	}
}

func (r *KeptnWebhookCertificateReconciler) getMutatingWebhookConfigurations(ctx context.Context) (
	*admissionregistrationv1.MutatingWebhookConfigurationList, error) {
	var mutatingWebhooks admissionregistrationv1.MutatingWebhookConfigurationList

	if err := r.Client.List(ctx, &mutatingWebhooks, client.MatchingLabels(r.MatchLabels)); err != nil {
		return nil, err
	}
	return &mutatingWebhooks, nil
}

func (r *KeptnWebhookCertificateReconciler) getValidatingWebhookConfigurations(ctx context.Context) (
	*admissionregistrationv1.ValidatingWebhookConfigurationList, error) {
	var validatingWebhooks admissionregistrationv1.ValidatingWebhookConfigurationList

	if err := r.Client.List(ctx, &validatingWebhooks, client.MatchingLabels(r.MatchLabels)); err != nil {
		return nil, err
	}
	return &validatingWebhooks, nil
}

func (r *KeptnWebhookCertificateReconciler) updateClientConfigurations(ctx context.Context, bundle []byte,
	webhookClientConfigs []*admissionregistrationv1.WebhookClientConfig, webhookConfig client.Object) error {
	if webhookConfig == nil || reflect.ValueOf(webhookConfig).IsNil() {
		return nil
	}

	for i := range webhookClientConfigs {
		webhookClientConfigs[i].CABundle = bundle
	}

	if err := r.Client.Update(ctx, webhookConfig); err != nil {
		return err
	}
	return nil
}

func (r *KeptnWebhookCertificateReconciler) getCRDConfigurations(ctx context.Context) (
	*apiv1.CustomResourceDefinitionList, error) {
	var crds apiv1.CustomResourceDefinitionList
	opt := client.MatchingLabels(r.MatchLabels)
	if err := r.Client.List(ctx, &crds, opt); err != nil {
		return nil, err
	}

	return &crds, nil
}

func (r *KeptnWebhookCertificateReconciler) updateCRDsConfiguration(ctx context.Context, crds *apiv1.CustomResourceDefinitionList, bundle []byte) error {
	fail := false
	for _, crd := range crds.Items {
		if err := r.updateCRDConfiguration(ctx, crd.Name, bundle); err != nil {
			fail = true
		}

	}
	if fail {
		return fmt.Errorf(couldNotUpdateCRDErr)
	}
	return nil
}

func (r *KeptnWebhookCertificateReconciler) updateCRDConfiguration(ctx context.Context, crdName string, bundle []byte) error {
	var crd apiv1.CustomResourceDefinition
	if err := r.Client.Get(ctx, types.NamespacedName{Name: crdName}, &crd); err != nil {
		return err
	}

	if !hasConversionWebhook(crd) {
		r.Log.Info(fmt.Sprintf("no conversion webhook config for %s, no cert will be provided", crdName))
		return nil
	}

	// update crd
	crd.Spec.Conversion.Webhook.ClientConfig.CABundle = bundle
	if err := r.Client.Update(ctx, &crd); err != nil {
		return err
	}
	return nil
}

func hasConversionWebhook(crd apiv1.CustomResourceDefinition) bool {
	return crd.Spec.Conversion != nil && crd.Spec.Conversion.Webhook != nil && crd.Spec.Conversion.Webhook.ClientConfig != nil
}
