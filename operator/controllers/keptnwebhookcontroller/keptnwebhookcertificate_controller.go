/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package keptnwebhookcontroller

import (
	"context"
	"github.com/go-logr/logr"
	lifecyclev1alpha1 "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"time"

	"github.com/keptn/lifecycle-toolkit/operator/webhooks"
	"github.com/pkg/errors"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	apiv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	SuccessDuration              = 3 * time.Hour
	secretPostfix                = "-certs"
	errorCertificatesSecretEmpty = "certificates secret is empty"
	crdName                      = "keptnwebhookcertificates.lifecycle.keptn.sh"
)

// KeptnWebhookCertificateReconciler reconciles a KeptnWebhookCertificate object
type KeptnWebhookCertificateReconciler struct {
	ctx           context.Context
	Client        client.Client
	Scheme        *runtime.Scheme
	ApiReader     client.Reader
	namespace     string
	CancelMgrFunc context.CancelFunc
	Log           logr.Logger
}

//+kubebuilder:rbac:groups=admissionregistration.k8s.io,resources=mutatingwebhookconfigurations,resourceNames=klc-mutating-webhook-configuration,verbs=get;list;watch
//+kubebuilder:rbac:groups="",namespace=keptn-lifecycle-toolkit-system,resources=secrets,verbs=get;list;watch;create;update;patch;delete

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
	r.namespace = request.Namespace
	r.ctx = ctx

	mutatingWebhookConfiguration, err := r.getMutatingWebhookConfiguration()
	if err != nil {
		// Generation must not be skipped because webhook startup routine listens for the secret
		// See cmd/operator/manager.go and cmd/operator/watcher.go
		r.Log.Info("could not find mutating webhook configuration, this is normal when deployed using OLM")
	}

	certSecret := newCertificateSecret()

	err = certSecret.setSecretFromReader(r.ctx, r.ApiReader, r.namespace, r.Log)
	if err != nil {
		return reconcile.Result{}, errors.WithStack(err)
	}

	err = certSecret.validateCertificates(r.namespace)
	if err != nil {
		return reconcile.Result{}, errors.WithStack(err)
	}

	mutatingWebhookConfigs := getClientConfigsFromMutatingWebhook(mutatingWebhookConfiguration)

	areMutatingWebhookConfigsValid := certSecret.areWebhookConfigsValid(mutatingWebhookConfigs)

	if certSecret.isRecent() &&
		areMutatingWebhookConfigsValid {
		r.Log.Info("secret for certificates up to date, skipping update")
		r.cancelMgr()
		return reconcile.Result{RequeueAfter: SuccessDuration}, nil
	}

	if err = certSecret.createOrUpdateIfNecessary(r.ctx, r.Client); err != nil {
		return reconcile.Result{}, errors.WithStack(err)
	}

	bundle, err := certSecret.loadCombinedBundle()
	if err != nil {
		return reconcile.Result{}, errors.WithStack(err)
	}

	err = r.updateClientConfigurations(bundle, mutatingWebhookConfigs, mutatingWebhookConfiguration)
	if err != nil {
		return reconcile.Result{}, errors.WithStack(err)
	}

	r.cancelMgr()
	return reconcile.Result{RequeueAfter: SuccessDuration}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KeptnWebhookCertificateReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&lifecyclev1alpha1.KeptnWebhookCertificate{}).
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
		Name: webhooks.DeploymentName,
	}, &mutatingWebhook)
	if err != nil {
		return nil, err
	}

	if len(mutatingWebhook.Webhooks) <= 0 {
		return nil, errors.New("mutating webhook configuration has no registered webhooks")
	}
	return &mutatingWebhook, nil
}

func (r *KeptnWebhookCertificateReconciler) getValidatingWebhookConfiguration() (
	*admissionregistrationv1.ValidatingWebhookConfiguration, error) {
	var mutatingWebhook admissionregistrationv1.ValidatingWebhookConfiguration
	err := r.ApiReader.Get(r.ctx, client.ObjectKey{
		Name: webhooks.DeploymentName,
	}, &mutatingWebhook)
	if err != nil {
		return nil, err
	}

	if len(mutatingWebhook.Webhooks) <= 0 {
		return nil, errors.New("validating webhook configuration has no registered webhooks")
	}
	return &mutatingWebhook, nil
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

func (r *KeptnWebhookCertificateReconciler) updateCRDConfiguration(crdName string, bundle []byte) error {
	var crd apiv1.CustomResourceDefinition
	if err := r.ApiReader.Get(r.ctx, types.NamespacedName{Name: crdName}, &crd); err != nil {
		return err
	}

	if !hasConversionWebhook(crd) {
		r.Log.Info("no conversion webhook config, no cert will be provided")
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
