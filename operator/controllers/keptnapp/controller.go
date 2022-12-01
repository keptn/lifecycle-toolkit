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
	"strconv"

	"github.com/go-logr/logr"
	klcv1alpha1 "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1"
	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1/common"
	controllererrors "github.com/keptn/lifecycle-toolkit/operator/controllers/errors"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/interfaces"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

// KeptnAppReconciler reconciles a KeptnApp object

type KeptnAppReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
	Log      logr.Logger
	Tracer   interfaces.ITracer
}

//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnapps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnapps/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnapps/finalizers,verbs=update
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnappversion,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnappversion/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnappversion/finalizers,verbs=update

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

	app := &klcv1alpha1.KeptnApp{}
	err := r.Get(ctx, req.NamespacedName, app)
	if errors.IsNotFound(err) {
		return reconcile.Result{}, nil
	}
	if err != nil {
		return reconcile.Result{}, fmt.Errorf(controllererrors.ErrCannotFetchAppMsg, err)
	}

	traceContextCarrier := propagation.MapCarrier(app.Annotations)
	ctx = otel.GetTextMapPropagator().Extract(ctx, traceContextCarrier)

	ctx, span := r.Tracer.Start(ctx, "reconcile_app", trace.WithSpanKind(trace.SpanKindConsumer))
	defer span.End()

	app.SetSpanAttributes(span)

	r.Log.Info("Reconciling Keptn App", "app", app.Name)

	appVersion := &klcv1alpha1.KeptnAppVersion{}

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
			r.Recorder.Event(app, "Warning", "AppVersionNotCreated", fmt.Sprintf("Could not create KeptnAppVersion / Namespace: %s, Name: %s ", appVersion.Namespace, appVersion.Name))
			return ctrl.Result{}, err
		}
		r.Recorder.Event(app, "Normal", "AppVersionCreated", fmt.Sprintf("Created KeptnAppVersion / Namespace: %s, Name: %s ", appVersion.Namespace, appVersion.Name))

		app.Status.CurrentVersion = app.Spec.Version
		if err := r.Client.Status().Update(ctx, app); err != nil {
			r.Log.Error(err, "could not update Current Version of App")
			return ctrl.Result{}, err
		}
		if app.Generation != 1 {
			if err := r.cancelDeprecatedAppVersions(ctx, *app); err != nil {
				r.Log.Error(err, "could not cancel deprecated appVersions for appVersion %s", appVersion.Name)
				return ctrl.Result{Requeue: true}, nil
			}
			//only cancel this in future
			if err := r.deleteWorkloadInstancesOfApp(ctx, *app); err != nil {
				r.Log.Error(err, "could not delete WIs for appVersion %s", appVersion.Name)
				return ctrl.Result{Requeue: true}, nil
			}
			if err := r.bumpRevisionOfWorkloadForApp(ctx, *app); err != nil {
				r.Log.Error(err, "could not bump workload revision for appVersion %s", appVersion.Name)
				return ctrl.Result{Requeue: true}, nil
			}

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
		For(&klcv1alpha1.KeptnApp{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Complete(r)
}

func (r *KeptnAppReconciler) createAppVersion(ctx context.Context, app *klcv1alpha1.KeptnApp) (*klcv1alpha1.KeptnAppVersion, error) {
	ctx, span := r.Tracer.Start(ctx, "create_app_version", trace.WithSpanKind(trace.SpanKindProducer))
	defer span.End()

	ctxAppTrace, spanAppTrace := r.Tracer.Start(ctx, app.GetAppVersionName(), trace.WithNewRoot(), trace.WithSpanKind(trace.SpanKindServer))
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
	appVersion.Spec.TraceId = appTraceContextCarrier
	err := controllerutil.SetControllerReference(app, &appVersion, r.Scheme)
	if err != nil {
		r.Log.Error(err, "could not set controller reference for AppVersion: "+appVersion.Name)
	}

	return &appVersion, err
}

func (r *KeptnAppReconciler) cancelDeprecatedAppVersions(ctx context.Context, app klcv1alpha1.KeptnApp) error {
	var resultErr error
	resultErr = nil
	for i := 1; i < int(app.Generation); i++ {
		deprecatedAppVersion := &klcv1alpha1.KeptnAppVersion{}
		err := r.Get(ctx, types.NamespacedName{Namespace: app.Namespace, Name: app.Name + "-" + app.Spec.Version + "-" + strconv.Itoa(i)}, deprecatedAppVersion)
		if errors.IsNotFound(err) {
			continue
		}

		if err != nil {
			r.Log.Error(err, "AppVersion not found")
			resultErr = err
			continue
		}

		deprecatedAppVersion.DeprecateRemainingPhases(common.PhaseDeprecated)
		if err := r.Client.Status().Update(ctx, deprecatedAppVersion); err != nil {
			r.Log.Error(err, "could not update appVersion %s status", deprecatedAppVersion.Name)
			resultErr = err
			continue
		}
	}
	return resultErr
}

func (r *KeptnAppReconciler) deleteWorkloadInstancesOfApp(ctx context.Context, app klcv1alpha1.KeptnApp) error {
	var resultErr error
	resultErr = nil

	//delete only if workloads of deprecated appVersion are failed or are stuck -> its current phase is AppDeploy
	previousAppVersion := r.getPreviousAppVersion(ctx, app)
	if previousAppVersion.Status.CurrentPhase != common.PhaseAppDeployment.ShortName {
		return resultErr
	}

	for _, w := range app.Spec.Workloads {
		bak := v1.DeletePropagationBackground
		wi := &klcv1alpha1.KeptnWorkloadInstance{
			ObjectMeta: v1.ObjectMeta{
				Name:      app.Name + "-" + w.Name + "-" + w.Version,
				Namespace: app.Namespace,
			},
		}
		if err := r.Client.Delete(ctx, wi, &client.DeleteOptions{PropagationPolicy: &bak}); err != nil {
			r.Log.Error(err, "could not delete WI %s", wi.Name)
			resultErr = err
			continue
		}
	}

	return resultErr
}

func (r *KeptnAppReconciler) bumpRevisionOfWorkloadForApp(ctx context.Context, app klcv1alpha1.KeptnApp) error {
	var resultErr error
	resultErr = nil

	for _, w := range app.Spec.Workloads {
		workload := &klcv1alpha1.KeptnWorkload{}
		if err := r.Client.Get(ctx, types.NamespacedName{Namespace: app.Namespace, Name: app.Name + "-" + w.Name}, workload); err != nil {
			r.Log.Error(err, "could not get workload %s", workload.Name)
			resultErr = err
			continue
		}

		workload.Spec.AppGeneration = app.Generation

		if err := r.Client.Update(ctx, workload); err != nil {
			r.Log.Error(err, "could not update workload %s", workload.Name)
			resultErr = err
			continue
		}

	}

	return resultErr
}

func (r *KeptnAppReconciler) getPreviousAppVersion(ctx context.Context, app klcv1alpha1.KeptnApp) *klcv1alpha1.KeptnAppVersion {
	appVersion := &klcv1alpha1.KeptnAppVersion{}
	for i := app.Generation - 1; i > 0; i-- {
		err := r.Client.Get(ctx, types.NamespacedName{Namespace: app.Namespace, Name: app.Name + "-" + app.Spec.Version + "-" + strconv.FormatInt(i, 10)}, appVersion)
		if err == nil && appVersion != nil {
			return appVersion
		}
	}
	return &klcv1alpha1.KeptnAppVersion{}
}
