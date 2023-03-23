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

package options

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	optionsv1alpha1 "github.com/keptn/lifecycle-toolkit/operator/apis/options/v1alpha1"
	controllercommon "github.com/keptn/lifecycle-toolkit/operator/controllers/common"
	controllererrors "github.com/keptn/lifecycle-toolkit/operator/controllers/errors"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// KeptnConfigReconciler reconciles a KeptnConfig object
type KeptnConfigReconciler struct {
	client.Client
	Scheme          *runtime.Scheme
	Log             logr.Logger
	LastAppliedSpec *optionsv1alpha1.KeptnConfigSpec
}

// +kubebuilder:rbac:groups=options.keptn.sh,resources=keptnconfigs,verbs=get;list;watch;create
// +kubebuilder:rbac:groups=options.keptn.sh,resources=keptnconfigs/status,verbs=get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *KeptnConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	r.Log.Info("Searching for KeptnConfig")

	config := &optionsv1alpha1.KeptnConfig{}
	err := r.Get(ctx, req.NamespacedName, config)
	if errors.IsNotFound(err) {
		return reconcile.Result{}, nil
	}
	if err != nil {
		return reconcile.Result{}, fmt.Errorf(controllererrors.ErrCannotRetrieveConfigMsg, err)
	}

	if r.LastAppliedSpec == nil {
		r.Log.Info("initializing KeptnConfig since no config was there before")
		r.initConfig()
	}

	result, err := r.reconcileOtelCollectorUrl(config)
	if err != nil {
		return result, err
	}

	r.LastAppliedSpec = &config.Spec
	return ctrl.Result{}, nil
}

func (r *KeptnConfigReconciler) reconcileOtelCollectorUrl(config *optionsv1alpha1.KeptnConfig) (ctrl.Result, error) {
	r.Log.Info(fmt.Sprintf("reconciling Keptn Config: %s", config.Name))
	otelConfig := controllercommon.GetOtelInstance()

	// collector URL changed, so we need to re-initialize the exporter
	if r.LastAppliedSpec.OTelCollectorUrl != config.Spec.OTelCollectorUrl {
		if err := otelConfig.InitOtelCollector(config.Spec.OTelCollectorUrl); err != nil {
			r.Log.Error(err, "unable to initialize OTel tracer options")
			return ctrl.Result{Requeue: true, RequeueAfter: 10 * time.Second}, err
		}
	}
	return ctrl.Result{}, nil
}

func (r *KeptnConfigReconciler) initConfig() {
	r.LastAppliedSpec = &optionsv1alpha1.KeptnConfigSpec{
		OTelCollectorUrl: "",
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *KeptnConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&optionsv1alpha1.KeptnConfig{}).
		Complete(r)
}
