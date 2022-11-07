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

package keptnappversion

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/types"

	controllercommon "github.com/keptn/lifecycle-toolkit/operator/controllers/common"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/go-logr/logr"
	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1/common"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	klcv1alpha1 "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// KeptnAppVersionReconciler reconciles a KeptnAppVersion object
type KeptnAppVersionReconciler struct {
	Scheme *runtime.Scheme
	client.Client
	Log         logr.Logger
	Recorder    record.EventRecorder
	Tracer      trace.Tracer
	Meters      common.KeptnMeters
	SpanHandler *controllercommon.SpanHandler
}

//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnappversions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnappversions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnappversions/finalizers,verbs=update
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnworkloadinstances/status,verbs=get;update;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the KeptnAppVersion object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *KeptnAppVersionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	r.Log.Info("Searching for Keptn App Version")

	appVersion := &klcv1alpha1.KeptnAppVersion{}
	err := r.Get(ctx, req.NamespacedName, appVersion)
	if errors.IsNotFound(err) {
		return reconcile.Result{}, nil
	}

	if err != nil {
		r.Log.Error(err, "App Version not found")
		return reconcile.Result{}, fmt.Errorf("could not fetch KeptnappVersion: %+v", err)
	}

	appVersion.SetStartTime()

	traceContextCarrier := propagation.MapCarrier(appVersion.Annotations)
	ctx = otel.GetTextMapPropagator().Extract(ctx, traceContextCarrier)

	appTraceContextCarrier := propagation.MapCarrier(appVersion.Spec.TraceId)
	ctxAppTrace := otel.GetTextMapPropagator().Extract(context.TODO(), appTraceContextCarrier)

	ctx, span := r.Tracer.Start(ctx, "reconcile_app_version", trace.WithSpanKind(trace.SpanKindConsumer))

	defer func(span trace.Span, appVersion *klcv1alpha1.KeptnAppVersion) {
		if appVersion.IsEndTimeSet() {
			r.Log.Info("Increasing app count")
			attrs := appVersion.GetMetricsAttributes()
			r.Meters.AppCount.Add(ctx, 1, attrs...)
		}
		span.End()
	}(span, appVersion)

	appVersion.SetSpanAttributes(span)

	phase := common.PhaseAppPreDeployment

	phaseHandler := controllercommon.PhaseHandler{
		Client:      r.Client,
		Recorder:    r.Recorder,
		Log:         r.Log,
		SpanHandler: r.SpanHandler,
	}

	if appVersion.Status.CurrentPhase == "" {
		if err := r.SpanHandler.UnbindSpan(appVersion, phase.ShortName); err != nil {
			r.Log.Error(err, "cannot unbind span")
		}
		var spanAppTrace trace.Span
		ctxAppTrace, spanAppTrace, err = r.SpanHandler.GetSpan(ctxAppTrace, r.Tracer, appVersion, phase.ShortName)
		if err != nil {
			r.Log.Error(err, "could not get span")
		}

		appVersion.SetSpanAttributes(spanAppTrace)
		spanAppTrace.AddEvent("App Version Pre-Deployment Tasks started", trace.WithTimestamp(time.Now()))
		controllercommon.RecordEvent(r.Recorder, phase, "Normal", appVersion, "Started", "have started", appVersion.GetVersion())
	}

	if !appVersion.IsPreDeploymentSucceeded() {
		reconcilePreDep := func() (common.KeptnState, error) {
			return r.reconcilePrePostDeployment(ctx, appVersion, common.PreDeploymentCheckType)
		}
		result, err := phaseHandler.HandlePhase(ctx, ctxAppTrace, r.Tracer, appVersion, phase, span, reconcilePreDep)
		if !result.Continue {
			return result.Result, err
		}
	}

	phase = common.PhaseAppPreEvaluation
	if !appVersion.IsPreDeploymentEvaluationSucceeded() {
		reconcilePreEval := func() (common.KeptnState, error) {
			return r.reconcilePrePostEvaluation(ctx, appVersion, common.PreDeploymentEvaluationCheckType)
		}
		result, err := phaseHandler.HandlePhase(ctx, ctxAppTrace, r.Tracer, appVersion, phase, span, reconcilePreEval)
		if !result.Continue {
			return result.Result, err
		}
	}

	phase = common.PhaseAppDeployment
	if !appVersion.AreWorkloadsSucceeded() {
		reconcileAppDep := func() (common.KeptnState, error) {
			return r.reconcileWorkloads(ctx, appVersion)
		}
		result, err := phaseHandler.HandlePhase(ctx, ctxAppTrace, r.Tracer, appVersion, phase, span, reconcileAppDep)
		if !result.Continue {
			return result.Result, err
		}
	}

	phase = common.PhaseAppPostDeployment
	if !appVersion.IsPostDeploymentSucceeded() {
		reconcilePostDep := func() (common.KeptnState, error) {
			return r.reconcilePrePostDeployment(ctx, appVersion, common.PostDeploymentCheckType)
		}
		result, err := phaseHandler.HandlePhase(ctx, ctxAppTrace, r.Tracer, appVersion, phase, span, reconcilePostDep)
		if !result.Continue {
			return result.Result, err
		}
	}

	phase = common.PhaseAppPostEvaluation
	if !appVersion.IsPostDeploymentEvaluationCompleted() {
		reconcilePostEval := func() (common.KeptnState, error) {
			return r.reconcilePrePostEvaluation(ctx, appVersion, common.PostDeploymentEvaluationCheckType)
		}
		result, err := phaseHandler.HandlePhase(ctx, ctxAppTrace, r.Tracer, appVersion, phase, span, reconcilePostEval)
		if !result.Continue {
			return result.Result, err
		}
	}

	controllercommon.RecordEvent(r.Recorder, phase, "Normal", appVersion, "Finished", "is finished", appVersion.GetVersion())
	err = r.Client.Status().Update(ctx, appVersion)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return ctrl.Result{Requeue: true}, err
	}

	// AppVersion is completed at this place

	if !appVersion.IsEndTimeSet() {
		appVersion.Status.CurrentPhase = common.PhaseCompleted.ShortName
		appVersion.SetEndTime()
	}

	err = r.Client.Status().Update(ctx, appVersion)
	if err != nil {
		return ctrl.Result{Requeue: true}, err
	}

	attrs := appVersion.GetMetricsAttributes()

	// metrics: add app duration
	duration := appVersion.Status.EndTime.Time.Sub(appVersion.Status.StartTime.Time)
	r.Meters.AppDuration.Record(ctx, duration.Seconds(), attrs...)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KeptnAppVersionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&klcv1alpha1.KeptnAppVersion{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Complete(r)
}
