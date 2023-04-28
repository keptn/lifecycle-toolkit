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
	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3/common"
	controllercommon "github.com/keptn/lifecycle-toolkit/operator/controllers/common"
	controllererrors "github.com/keptn/lifecycle-toolkit/operator/controllers/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const traceComponentName = "keptn/operator/app"

// KeptnAppReconciler reconciles a KeptnApp object
type KeptnAppReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	Recorder      record.EventRecorder
	Log           logr.Logger
	TracerFactory controllercommon.TracerFactory
}

// +kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnapps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnapps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnapps/finalizers,verbs=update
// +kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnappversion,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnappversion/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnappversion/finalizers,verbs=update

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
	r.Log.Info("Searching for App")

	app := &klcv1alpha3.KeptnApp{}
	err := r.Get(ctx, req.NamespacedName, app)
	if errors.IsNotFound(err) {
		return reconcile.Result{}, nil
	}
	if err != nil {
		return reconcile.Result{}, fmt.Errorf(controllererrors.ErrCannotFetchAppMsg, err)
	}

	traceContextCarrier := propagation.MapCarrier(app.Annotations)
	ctx = otel.GetTextMapPropagator().Extract(ctx, traceContextCarrier)

	ctx, span := r.getTracer().Start(ctx, "reconcile_app", trace.WithSpanKind(trace.SpanKindConsumer))
	defer span.End()

	app.SetSpanAttributes(span)

	r.Log.Info("Reconciling Keptn App", "app", app.Name)

	appVersion := &klcv1alpha3.KeptnAppVersion{}

	// Try to find the AppVersion
	err = r.Get(ctx, types.NamespacedName{Namespace: app.Namespace, Name: app.GetAppVersionName()}, appVersion)
	// If the app instance does not exist, create it
	if errors.IsNotFound(err) {
		appVersion, err := r.createAppVersion(ctx, app)
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			return reconcile.Result{}, err
		}
		err = r.Client.Create(ctx, appVersion)
		if err != nil {
			r.Log.Error(err, "could not create AppVersion")
			span.SetStatus(codes.Error, err.Error())
			controllercommon.RecordEvent(r.Recorder, common.PhaseCreateAppVersion, "Warning", appVersion, "AppVersionNotCreated", "Could not create KeptnAppVersion", appVersion.Spec.Version)
			return ctrl.Result{}, err
		}
		controllercommon.RecordEvent(r.Recorder, common.PhaseCreateAppVersion, "Normal", appVersion, "AppVersionCreated", "created KeptnAppVersion", appVersion.Spec.Version)

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
		span.SetStatus(codes.Error, err.Error())
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KeptnAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&klcv1alpha3.KeptnApp{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Complete(r)
}

func (r *KeptnAppReconciler) createAppVersion(ctx context.Context, app *klcv1alpha3.KeptnApp) (*klcv1alpha3.KeptnAppVersion, error) {
	ctx, span := r.getTracer().Start(ctx, "create_app_version", trace.WithSpanKind(trace.SpanKindProducer))
	defer span.End()

	ctxAppTrace, spanAppTrace := r.getTracer().Start(ctx, app.GetAppVersionName(), trace.WithNewRoot(), trace.WithSpanKind(trace.SpanKindServer))
	defer spanAppTrace.End()

	app.SetSpanAttributes(span)
	app.SetSpanAttributes(spanAppTrace)

	// create TraceContext
	// follow up with a Keptn propagator that JSON-encoded the OTel map into our own key
	traceContextCarrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, traceContextCarrier)
	appTraceContextCarrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctxAppTrace, appTraceContextCarrier)

	previousVersion := ""
	if app.Spec.Version != app.Status.CurrentVersion {
		previousVersion = app.Status.CurrentVersion
	}

	appVersion := app.GenerateAppVersion(previousVersion, traceContextCarrier)
	if len(appVersion.ObjectMeta.Name) > 253 {
		appVersion.ObjectMeta.Name = common.TruncateString(appVersion.ObjectMeta.Name, 253)
	}
	appVersion.Spec.TraceId = appTraceContextCarrier
	err := controllerutil.SetControllerReference(app, &appVersion, r.Scheme)
	if err != nil {
		r.Log.Error(err, "could not set controller reference for AppVersion: "+appVersion.Name)
	}

	return &appVersion, err
}

func (r *KeptnAppReconciler) handleGenerationBump(ctx context.Context, app *klcv1alpha3.KeptnApp) error {
	if app.Generation != 1 {
		if err := r.deprecateAppVersions(ctx, app); err != nil {
			r.Log.Error(err, "could not deprecate appVersions for appVersion %s", app.GetAppVersionName())
			controllercommon.RecordEvent(r.Recorder, common.PhaseCreateAppVersion, "Warning", app, "AppVersionNotDeprecated", fmt.Sprintf("could not deprecate KeptnAppVersions for KeptnAppVersion: %s", app.GetAppVersionName()), app.Spec.Version)
			return err
		}
		controllercommon.RecordEvent(r.Recorder, common.PhaseCreateAppVersion, "Normal", app, "AppVersionDeprecated", fmt.Sprintf("deprecated KeptnAppVersions for KeptnAppVersion: %s", app.GetAppVersionName()), app.Spec.Version)
	}
	return nil
}

func (r *KeptnAppReconciler) deprecateAppVersions(ctx context.Context, app *klcv1alpha3.KeptnApp) error {
	var lastResultErr error
	lastResultErr = nil
	for i := app.Generation - 1; i > 0; i-- {
		deprecatedAppVersion := &klcv1alpha3.KeptnAppVersion{}
		if err := r.Get(ctx, types.NamespacedName{Namespace: app.Namespace, Name: app.Name + "-" + app.Spec.Version + "-" + common.Hash(i)}, deprecatedAppVersion); err != nil {
			if !errors.IsNotFound(err) {
				r.Log.Error(err, fmt.Sprintf("Could not get KeptnAppVersion: %s", app.Name+"-"+app.Spec.Version+"-"+common.Hash(i)))
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

func (r *KeptnAppReconciler) getTracer() controllercommon.ITracer {
	return r.TracerFactory.GetTracer(traceComponentName)
}
