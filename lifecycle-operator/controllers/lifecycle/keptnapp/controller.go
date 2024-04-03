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

package keptnapp

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	apilifecycle "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1/common"
	operatorcommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/common"
	controllercommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/eventsender"
	controllererrors "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// KeptnAppReconciler reconciles a KeptnApp object
type KeptnAppReconciler struct {
	client.Client
	Scheme      *runtime.Scheme
	EventSender eventsender.IEvent
	Log         logr.Logger
}

// +kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnapps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnapps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnapps/finalizers,verbs=update
// +kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnappversion,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnappversion/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnappversion/finalizers,verbs=update
// +kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnappcontexts,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the KeptnApp object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *KeptnAppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	requestInfo := controllercommon.GetRequestInfo(req)
	r.Log.Info("Searching for App", "requestInfo", requestInfo)

	app := &apilifecycle.KeptnApp{}
	err := r.Get(ctx, req.NamespacedName, app)
	if errors.IsNotFound(err) {
		return reconcile.Result{}, nil
	}
	if err != nil {
		return reconcile.Result{}, fmt.Errorf(controllererrors.ErrCannotFetchAppMsg, err)
	}

	traceContextCarrier := propagation.MapCarrier(app.Annotations)
	ctx = otel.GetTextMapPropagator().Extract(ctx, traceContextCarrier)

	r.Log.Info("Reconciling Keptn App", "app", app.Name)

	appVersion := &apilifecycle.KeptnAppVersion{}

	// Try to find the AppVersion
	err = r.Get(ctx, types.NamespacedName{Namespace: app.Namespace, Name: app.GetAppVersionName()}, appVersion)
	// If the app instance does not exist, create it
	if errors.IsNotFound(err) {
		appContext := &apilifecycle.KeptnAppContext{}
		err := r.Get(ctx, types.NamespacedName{
			Namespace: app.Namespace,
			Name:      app.Name,
		}, appContext)
		if client.IgnoreNotFound(err) != nil {
			r.Log.Error(err, "Could not look up related KeptnAppContext", "requestInfo", requestInfo)
		}

		appVersion, err := r.createAppVersion(ctx, app, appContext)
		if err != nil {
			return reconcile.Result{}, err
		}
		err = r.Client.Create(ctx, appVersion)
		if err != nil {
			r.Log.Error(err, "could not create AppVersion")
			r.EventSender.Emit(common.PhaseCreateAppVersion, "Warning", appVersion, common.PhaseStateFailed, "Could not create KeptnAppVersion", appVersion.Spec.Version)
			return ctrl.Result{}, err
		}

		app.Status.CurrentVersion = app.Spec.Version
		if err := r.Client.Status().Update(ctx, app); err != nil {
			r.Log.Error(err, "could not update Current Version of App")
			return ctrl.Result{}, err
		}
		if err := r.handleGenerationBump(ctx, app); err != nil {
			return ctrl.Result{Requeue: true}, nil
		}
		return ctrl.Result{}, nil
	}
	if err != nil {
		r.Log.Error(err, "could not get AppVersion")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KeptnAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apilifecycle.KeptnApp{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Complete(r)
}

func (r *KeptnAppReconciler) createAppVersion(ctx context.Context, app *apilifecycle.KeptnApp, appContext *apilifecycle.KeptnAppContext) (*apilifecycle.KeptnAppVersion, error) {

	previousVersion := ""
	if app.Spec.Version != app.Status.CurrentVersion {
		previousVersion = app.Status.CurrentVersion
	}

	appVersion := app.GenerateAppVersion(previousVersion)

	appVersion.Spec.KeptnAppContextSpec = appContext.Spec

	err := controllerutil.SetControllerReference(app, &appVersion, r.Scheme)
	if err != nil {
		r.Log.Error(err, "could not set controller reference for AppVersion: "+appVersion.Name)
	}

	return &appVersion, err
}

func (r *KeptnAppReconciler) handleGenerationBump(ctx context.Context, app *apilifecycle.KeptnApp) error {
	if app.Generation != 1 {
		if err := r.deprecateAppVersions(ctx, app); err != nil {
			r.Log.Error(err, "could not deprecate appVersions for appVersion %s", app.GetAppVersionName())
			r.EventSender.Emit(common.PhaseDeprecateAppVersion, "Warning", app, common.PhaseStateFailed, fmt.Sprintf("could not deprecate outdated revisions of KeptnAppVersion: %s", app.GetAppVersionName()), app.Spec.Version)
			return err
		}
	}
	return nil
}

func (r *KeptnAppReconciler) deprecateAppVersions(ctx context.Context, app *apilifecycle.KeptnApp) error {
	var lastResultErr error
	lastResultErr = nil
	for i := app.Generation - 1; i > 0; i-- {
		deprecatedAppVersion := &apilifecycle.KeptnAppVersion{}
		appVersionName := operatorcommon.CreateResourceName(common.MaxK8sObjectLength, common.MinKeptnNameLen, app.Name, app.Spec.Version, common.Hash(i))
		if err := r.Get(ctx, types.NamespacedName{Namespace: app.Namespace, Name: appVersionName}, deprecatedAppVersion); err != nil {
			if !errors.IsNotFound(err) {
				r.Log.Error(err, fmt.Sprintf("Could not get KeptnAppVersion: %s", appVersionName))
				lastResultErr = err
			}
		} else if !deprecatedAppVersion.Status.Status.IsDeprecated() {
			deprecatedAppVersion.DeprecateRemainingPhases(common.PhaseDeprecated)
			if err := r.Client.Status().Update(ctx, deprecatedAppVersion); err != nil {
				r.Log.Error(err, "could not update appVersion %s status", deprecatedAppVersion.Name)
				lastResultErr = err
			}
		}
	}
	return lastResultErr
}
