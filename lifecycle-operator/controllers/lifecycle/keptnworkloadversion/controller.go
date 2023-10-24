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

package keptnworkloadversion

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	klcv1alpha4 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha4"
	controllercommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/evaluation"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/telemetry"
	controllererrors "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const traceComponentName = "keptn/lifecycle-operator/workloadversion"

// KeptnWorkloadVersionReconciler reconciles a KeptnWorkloadVersion object
type KeptnWorkloadVersionReconciler struct {
	client.Client
	Scheme                 *runtime.Scheme
	EventSender            controllercommon.IEvent
	Log                    logr.Logger
	Meters                 apicommon.KeptnMeters
	SpanHandler            *telemetry.SpanHandler
	TracerFactory          telemetry.TracerFactory
	SchedulingGatesHandler controllercommon.ISchedulingGatesHandler
	EvaluationHandler      evaluation.EvaluationHandlerInterface
}

// +kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnworkloadversions,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnworkloadversions/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnworkloadversions/finalizers,verbs=update
// +kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptntasks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptntasks/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptntasks/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=events,verbs=create;watch;patch
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;update
// +kubebuilder:rbac:groups=apps,resources=replicasets;deployments;statefulsets;daemonsets,verbs=get;list;watch
// +kubebuilder:rbac:groups=argoproj.io,resources=rollouts,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func NewReconciler(client client.Client, scheme *runtime.Scheme, iEvent controllercommon.IEvent, log logr.Logger, keptnMeters apicommon.KeptnMeters, spanHandler *telemetry.SpanHandler, tracerFactory telemetry.TracerFactory, iSchedulingGatesHandler controllercommon.ISchedulingGatesHandler, tracer trace.Tracer) *KeptnWorkloadVersionReconciler {
	return &KeptnWorkloadVersionReconciler{
		Client:                 client,
		Scheme:                 scheme,
		EventSender:            iEvent,
		Log:                    log,
		Meters:                 keptnMeters,
		SpanHandler:            spanHandler,
		TracerFactory:          tracerFactory,
		SchedulingGatesHandler: iSchedulingGatesHandler,
		EvaluationHandler:      evaluation.NewEvaluationHandler(client, iEvent, log, tracer, scheme, spanHandler),
	}
}

//nolint:gocyclo,gocognit
func (r *KeptnWorkloadVersionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	requestInfo := controllercommon.GetRequestInfo(req)
	r.Log.Info("Searching for KeptnWorkloadVersion", "requestInfo", requestInfo)

	// retrieve workload version
	workloadVersion := &klcv1alpha4.KeptnWorkloadVersion{}
	err := r.Get(ctx, req.NamespacedName, workloadVersion)
	if errors.IsNotFound(err) {
		return reconcile.Result{}, nil
	}

	if err != nil {
		r.Log.Error(err, "Could not retrieve KeptnWorkloadVersion", "requestInfo", requestInfo)
		return reconcile.Result{}, fmt.Errorf(controllererrors.ErrCannotRetrieveWorkloadVersionMsg, err)
	}

	completionFunc := r.getCompletionFunc(ctx, workloadVersion)
	defer completionFunc(workloadVersion)

	if requeue, err := r.checkPreEvaluationStatusOfApp(ctx, workloadVersion); requeue {
		return ctrl.Result{Requeue: true, RequeueAfter: 10 * time.Second}, controllererrors.IgnoreReferencedResourceNotFound(err)
	}

	appTraceContextCarrier := propagation.MapCarrier(workloadVersion.Spec.TraceId)
	ctxAppTrace := otel.GetTextMapPropagator().Extract(context.TODO(), appTraceContextCarrier)

	// Wait for pre-deployment checks of Workload
	phase := apicommon.PhaseWorkloadPreDeployment
	phaseHandler := controllercommon.PhaseHandler{
		Client:      r.Client,
		EventSender: r.EventSender,
		Log:         r.Log,
		SpanHandler: r.SpanHandler,
	}

	// this will be the parent span for all phases of the WorkloadVersion
	ctxWorkloadTrace, spanWorkloadTrace, err := r.SpanHandler.GetSpan(ctxAppTrace, r.getTracer(), workloadVersion, "")
	if err != nil {
		r.Log.Error(err, "could not get span")
	}

	if workloadVersion.Status.CurrentPhase == "" {
		spanWorkloadTrace.AddEvent("WorkloadVersion Pre-Deployment Tasks started", trace.WithTimestamp(time.Now()))
	}

	if !workloadVersion.IsPreDeploymentSucceeded() {
		reconcilePre := func(phaseCtx context.Context) (apicommon.KeptnState, error) {
			return r.reconcilePrePostDeployment(ctx, phaseCtx, workloadVersion, apicommon.PreDeploymentCheckType)
		}
		result, err := phaseHandler.HandlePhase(ctx, ctxWorkloadTrace, r.getTracer(), workloadVersion, phase, reconcilePre)
		if !result.Continue {
			return result.Result, err
		}
	}

	// Wait for pre-evaluation checks of Workload
	phase = apicommon.PhaseWorkloadPreEvaluation
	if !workloadVersion.IsPreDeploymentEvaluationSucceeded() {
		reconcilePreEval := func(phaseCtx context.Context) (apicommon.KeptnState, error) {
			return r.reconcilePrePostEvaluation(ctx, phaseCtx, workloadVersion, apicommon.PreDeploymentEvaluationCheckType)
		}
		result, err := phaseHandler.HandlePhase(ctx, ctxWorkloadTrace, r.getTracer(), workloadVersion, phase, reconcilePreEval)
		if !result.Continue {
			return result.Result, err
		}
	}

	if r.SchedulingGatesHandler.Enabled() {
		// pre-evaluation checks done at this moment, we can remove the gate
		if err := r.SchedulingGatesHandler.RemoveGates(ctx, workloadVersion); err != nil {
			r.Log.Error(err, "could not remove SchedulingGates")
			return ctrl.Result{Requeue: true, RequeueAfter: 10 * time.Second}, err
		}
	}

	// Wait for deployment of Workload
	phase = apicommon.PhaseWorkloadDeployment
	if !workloadVersion.IsDeploymentSucceeded() {
		reconcileWorkloadVersion := func(phaseCtx context.Context) (apicommon.KeptnState, error) {
			return r.reconcileDeployment(ctx, workloadVersion)
		}
		result, err := phaseHandler.HandlePhase(ctx, ctxWorkloadTrace, r.getTracer(), workloadVersion, phase, reconcileWorkloadVersion)
		if !result.Continue {
			return result.Result, err
		}
	}

	// Wait for post-deployment checks of Workload
	phase = apicommon.PhaseWorkloadPostDeployment
	if !workloadVersion.IsPostDeploymentSucceeded() {
		reconcilePostDeployment := func(phaseCtx context.Context) (apicommon.KeptnState, error) {
			return r.reconcilePrePostDeployment(ctx, phaseCtx, workloadVersion, apicommon.PostDeploymentCheckType)
		}
		result, err := phaseHandler.HandlePhase(ctx, ctxWorkloadTrace, r.getTracer(), workloadVersion, phase, reconcilePostDeployment)
		if !result.Continue {
			return result.Result, err
		}
	}

	// Wait for post-evaluation checks of Workload
	phase = apicommon.PhaseWorkloadPostEvaluation
	if !workloadVersion.IsPostDeploymentEvaluationSucceeded() {
		reconcilePostEval := func(phaseCtx context.Context) (apicommon.KeptnState, error) {
			return r.reconcilePrePostEvaluation(ctx, phaseCtx, workloadVersion, apicommon.PostDeploymentEvaluationCheckType)
		}
		result, err := phaseHandler.HandlePhase(ctx, ctxWorkloadTrace, r.getTracer(), workloadVersion, phase, reconcilePostEval)
		if !result.Continue {
			return result.Result, err
		}
	}

	// WorkloadVersion is completed at this place
	return r.finishKeptnWorkloadVersionReconcile(ctx, workloadVersion, spanWorkloadTrace)
}

func (r *KeptnWorkloadVersionReconciler) finishKeptnWorkloadVersionReconcile(ctx context.Context, workloadVersion *klcv1alpha4.KeptnWorkloadVersion, spanWorkloadTrace trace.Span) (ctrl.Result, error) {
	if !workloadVersion.IsEndTimeSet() {
		workloadVersion.Status.CurrentPhase = apicommon.PhaseCompleted.ShortName
		workloadVersion.Status.Status = apicommon.StateSucceeded
		workloadVersion.SetEndTime()
	}

	err := r.Client.Status().Update(ctx, workloadVersion)
	if err != nil {
		return ctrl.Result{Requeue: true}, err
	}

	r.EventSender.Emit(apicommon.PhaseWorkloadCompleted, "Normal", workloadVersion, apicommon.PhaseStateFinished, "has finished", workloadVersion.GetVersion())

	attrs := workloadVersion.GetMetricsAttributes()

	// metrics: add deployment duration
	duration := workloadVersion.Status.EndTime.Time.Sub(workloadVersion.Status.StartTime.Time)
	r.Meters.DeploymentDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(attrs...))

	spanWorkloadTrace.AddEvent(workloadVersion.Name + " has finished")
	spanWorkloadTrace.SetStatus(codes.Ok, "Finished")
	spanWorkloadTrace.End()
	if err := r.SpanHandler.UnbindSpan(workloadVersion, ""); err != nil {
		r.Log.Error(err, controllererrors.ErrCouldNotUnbindSpan, workloadVersion.Name)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KeptnWorkloadVersionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &klcv1alpha4.KeptnWorkloadVersion{}, "spec.app", func(rawObj client.Object) []string {
		workloadVersion := rawObj.(*klcv1alpha4.KeptnWorkloadVersion)
		return []string{workloadVersion.Spec.AppName}
	}); err != nil {
		return err
	}
	return ctrl.NewControllerManagedBy(mgr).
		// predicate disabling the auto reconciliation after updating the object status
		For(&klcv1alpha4.KeptnWorkloadVersion{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Complete(r)
}

func (r *KeptnWorkloadVersionReconciler) sendUnfinishedPreEvaluationEvents(appPreEvalStatus apicommon.KeptnState, phase apicommon.KeptnPhaseType, workloadVersion *klcv1alpha4.KeptnWorkloadVersion) {
	if appPreEvalStatus.IsFailed() {
		r.EventSender.Emit(phase, "Warning", workloadVersion, apicommon.PhaseStateFailed, "has failed since app has failed", workloadVersion.GetVersion())
	}
}

func (r *KeptnWorkloadVersionReconciler) getCompletionFunc(ctx context.Context, workloadVersion *klcv1alpha4.KeptnWorkloadVersion) func(workloadVersion *klcv1alpha4.KeptnWorkloadVersion) {
	workloadVersion.SetStartTime()

	endFunc := func(workloadVersion *klcv1alpha4.KeptnWorkloadVersion) {
		if workloadVersion.IsEndTimeSet() {
			r.Log.Info("Increasing deployment count")
			attrs := workloadVersion.GetMetricsAttributes()
			r.Meters.DeploymentCount.Add(ctx, 1, metric.WithAttributes(attrs...))
		}
	}

	return endFunc
}

func (r *KeptnWorkloadVersionReconciler) checkPreEvaluationStatusOfApp(ctx context.Context, workloadVersion *klcv1alpha4.KeptnWorkloadVersion) (bool, error) {
	// Wait for pre-evaluation checks of App
	// Only check if we have not begun with the first phase of the workload version, to avoid retrieving the KeptnAppVersion
	// in each reconciliation loop
	if workloadVersion.GetCurrentPhase() != "" {
		return false, nil
	}
	phase := apicommon.PhaseAppPreEvaluation
	found, appVersion, err := r.getAppVersionForWorkloadVersion(ctx, workloadVersion)
	if err != nil {
		r.EventSender.Emit(phase, "Warning", workloadVersion, "GetAppVersionFailed", "has failed since app could not be retrieved", workloadVersion.GetVersion())
		return true, fmt.Errorf(controllererrors.ErrCannotFetchAppVersionForWorkloadVersionMsg + err.Error())
	} else if !found {
		r.EventSender.Emit(phase, "Warning", workloadVersion, "AppVersionNotFound", "has failed since app could not be found", workloadVersion.GetVersion())
		return true, controllererrors.ErrNoMatchingAppVersionFound
	}

	appPreEvalStatus := appVersion.Status.PreDeploymentEvaluationStatus
	if !appPreEvalStatus.IsSucceeded() {
		r.sendUnfinishedPreEvaluationEvents(appPreEvalStatus, phase, workloadVersion)
		return true, nil
	}

	// set the App trace id if not already set
	if len(workloadVersion.Spec.TraceId) < 1 {
		appDeploymentTraceID := appVersion.Status.PhaseTraceIDs[apicommon.PhaseAppDeployment.ShortName]
		if appDeploymentTraceID != nil {
			workloadVersion.Spec.TraceId = appDeploymentTraceID
		} else {
			workloadVersion.Spec.TraceId = appVersion.Spec.TraceId
		}
		if err := r.Update(ctx, workloadVersion); err != nil {
			return true, err
		}
	}
	return false, nil
}

func (r *KeptnWorkloadVersionReconciler) getAppVersionForWorkloadVersion(ctx context.Context, wli *klcv1alpha4.KeptnWorkloadVersion) (bool, klcv1alpha3.KeptnAppVersion, error) {
	apps := &klcv1alpha3.KeptnAppVersionList{}

	// TODO add label selector for looking up by name?
	if err := r.Client.List(ctx, apps, client.InNamespace(wli.Namespace)); err != nil {
		return false, klcv1alpha3.KeptnAppVersion{}, err
	}

	// due to effectivity reasons deprecated KeptnAppVersions are removed from the list, as there is
	// no point in iterating through them in the next steps
	apps.RemoveDeprecated()

	workloadFound, latestVersion, err := getLatestAppVersion(apps, wli)
	if err != nil {
		r.Log.Error(err, "could not look up KeptnAppVersion for WorkloadVersion")
		return false, latestVersion, err
	}

	// If the latest version is empty or the workload is not found, return false and empty result
	if latestVersion.Spec.Version == "" || !workloadFound {
		return false, klcv1alpha3.KeptnAppVersion{}, nil
	}
	return true, latestVersion, nil
}

func (r *KeptnWorkloadVersionReconciler) getTracer() telemetry.ITracer {
	return r.TracerFactory.GetTracer(traceComponentName)
}

func getLatestAppVersion(apps *klcv1alpha3.KeptnAppVersionList, wli *klcv1alpha4.KeptnWorkloadVersion) (bool, klcv1alpha3.KeptnAppVersion, error) {
	latestVersion := klcv1alpha3.KeptnAppVersion{}

	workloadFound := false
	for _, app := range apps.Items {
		if app.Spec.AppName == wli.Spec.AppName {
			for _, appWorkload := range app.Spec.Workloads {
				if workloadMatchesApp(appWorkload, wli, app) {
					workloadFound = true

					if isNewer(app, latestVersion) {
						latestVersion = app
					}
				}
			}
		}
	}
	return workloadFound, latestVersion, nil
}

func isNewer(app klcv1alpha3.KeptnAppVersion, latestVersion klcv1alpha3.KeptnAppVersion) bool {
	return app.ObjectMeta.CreationTimestamp.Time.After(latestVersion.ObjectMeta.CreationTimestamp.Time) || latestVersion.CreationTimestamp.Time.IsZero()
}

func workloadMatchesApp(appWorkload klcv1alpha3.KeptnWorkloadRef, wli *klcv1alpha4.KeptnWorkloadVersion, app klcv1alpha3.KeptnAppVersion) bool {
	return appWorkload.Version == wli.Spec.Version && app.GetWorkloadNameOfApp(appWorkload.Name) == wli.Spec.WorkloadName
}
