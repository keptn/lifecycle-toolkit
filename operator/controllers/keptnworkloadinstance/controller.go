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

package keptnworkloadinstance

import (
	"context"
	"fmt"
	"golang.org/x/mod/semver"
	"time"

	"github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1/semconv"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	klcv1alpha1 "github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1"
	"github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1/common"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// KeptnWorkloadInstanceReconciler reconciles a KeptnWorkloadInstance object
type KeptnWorkloadInstanceReconciler struct {
	client.Client
	Scheme    *runtime.Scheme
	Recorder  record.EventRecorder
	Log       logr.Logger
	Meters    common.KeptnMeters
	Tracer    trace.Tracer
	AppTracer trace.Tracer
}

//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnworkloadinstances,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnworkloadinstances/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnworkloadinstances/finalizers,verbs=update
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptntasks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptntasks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptntasks/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;watch;patch
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch
//+kubebuilder:rbac:groups=apps,resources=replicasets;deployments;statefulsets,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the KeptnWorkloadInstance object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile

var bindCRDSpan = make(map[string]trace.Span, 100)

func (r *KeptnWorkloadInstanceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	r.Log.Info("Searching for Keptn Workload Instance")

	//retrieve workload instance
	workloadInstance := &klcv1alpha1.KeptnWorkloadInstance{}
	err := r.Get(ctx, req.NamespacedName, workloadInstance)
	if errors.IsNotFound(err) {
		return reconcile.Result{}, nil
	}

	if err != nil {
		r.Log.Error(err, "Workload Instance not found")
		return reconcile.Result{}, fmt.Errorf("could not fetch KeptnWorkloadInstance: %+v", err)
	}

	found, appVersion, err := r.getAppVersionForWorkloadInstance(ctx, workloadInstance)
	if err != nil {
		return reconcile.Result{Requeue: true, RequeueAfter: 10 * time.Second}, fmt.Errorf("could not fetch AppVersion for KeptnWorkloadInstance: %+v", err)
	} else if !found {
		return reconcile.Result{Requeue: true, RequeueAfter: 10 * time.Second}, fmt.Errorf("could not find AppVersion for KeptnWorkloadInstance")
	}

	//setup otel
	traceContextCarrier := propagation.MapCarrier(workloadInstance.Annotations)
	appTraceContextCarrier := propagation.MapCarrier(appVersion.Spec.TraceId)

	ctx = otel.GetTextMapPropagator().Extract(ctx, traceContextCarrier)
	ctxAppTrace := otel.GetTextMapPropagator().Extract(context.TODO(), appTraceContextCarrier)

	ctx, span := r.Tracer.Start(ctx, "reconcile_workload_instance", trace.WithSpanKind(trace.SpanKindConsumer))
	defer span.End()

	var spanAppTrace trace.Span

	semconv.AddAttributeFromWorkloadInstance(span, *workloadInstance)

	if !workloadInstance.IsStartTimeSet() {
		// metrics: increment active deployment counter
		r.Meters.DeploymentActive.Add(ctx, 1, workloadInstance.GetActiveMetricsAttributes()...)
		workloadInstance.SetStartTime()
	}

	//Wait for pre-evaluation checks of App
	phase := common.PhaseAppPreEvaluation

	appPreEvalStatus := appVersion.Status.PreDeploymentEvaluationStatus
	if !appPreEvalStatus.IsSucceeded() {
		if appPreEvalStatus.IsFailed() {
			r.recordEvent(phase, "Warning", workloadInstance, "Failed", "has failed since app has failed")
			return ctrl.Result{Requeue: true, RequeueAfter: 60 * time.Second}, nil
		}
		r.recordEvent(phase, "Normal", workloadInstance, "NotFinished", "Pre evaluations tasks for app not finished")
		return ctrl.Result{Requeue: true, RequeueAfter: 20 * time.Second}, nil
	}

	//Wait for pre-deployment checks of Workload
	phase = common.PhaseWorkloadPreDeployment

	if appVersion.Status.CurrentPhase == "" {
		ctx, spanAppTrace = r.getSpan(ctxAppTrace, fmt.Sprintf("%s/%s", workloadInstance.Spec.WorkloadName, phase.ShortName))
		semconv.AddAttributeFromAppVersion(spanAppTrace, appVersion)
		spanAppTrace.AddEvent("WorkloadInstance Pre-Deployment Tasks started", trace.WithTimestamp(time.Now()))
		r.recordEvent(phase, "Normal", workloadInstance, "Started", "have started")
	}

	if !workloadInstance.IsPreDeploymentSucceeded() {
		reconcilePre := func() (common.KeptnState, error) {
			return r.reconcilePrePostDeployment(ctx, workloadInstance, common.PreDeploymentCheckType)
		}
		return r.handlePhase(ctx, ctxAppTrace, workloadInstance, phase, span, workloadInstance.IsPreDeploymentFailed, reconcilePre)
	}

	//Wait for pre-evaluation checks of Workload
	phase = common.PhaseWorkloadPreEvaluation
	if !workloadInstance.IsPreDeploymentEvaluationSucceeded() {
		reconcilePreEval := func() (common.KeptnState, error) {
			return r.reconcilePrePostEvaluation(ctx, workloadInstance, common.PreDeploymentEvaluationCheckType)
		}
		return r.handlePhase(ctx, ctxAppTrace, workloadInstance, phase, span, workloadInstance.IsPreDeploymentEvaluationFailed, reconcilePreEval)
	}

	//Wait for deployment of Workload
	phase = common.PhaseWorkloadDeployment
	if !workloadInstance.IsDeploymentSucceeded() {
		reconcileWorkloadInstance := func() (common.KeptnState, error) {
			return r.reconcileDeployment(ctx, workloadInstance)
		}
		return r.handlePhase(ctx, ctxAppTrace, workloadInstance, phase, span, workloadInstance.IsDeploymentFailed, reconcileWorkloadInstance)
	}

	//Wait for post-deployment checks of Workload
	phase = common.PhaseWorkloadPostDeployment
	if !workloadInstance.IsPostDeploymentSucceeded() {
		reconcilePostDeployment := func() (common.KeptnState, error) {
			return r.reconcilePrePostDeployment(ctx, workloadInstance, common.PostDeploymentCheckType)
		}
		return r.handlePhase(ctx, ctxAppTrace, workloadInstance, phase, span, workloadInstance.IsPostDeploymentFailed, reconcilePostDeployment)
	}

	//Wait for post-evaluation checks of Workload
	phase = common.PhaseAppPostEvaluation
	if !workloadInstance.IsPostDeploymentEvaluationSucceeded() {
		reconcilePostEval := func() (common.KeptnState, error) {
			return r.reconcilePrePostEvaluation(ctx, workloadInstance, common.PostDeploymentEvaluationCheckType)
		}
		return r.handlePhase(ctx, ctxAppTrace, workloadInstance, phase, span, workloadInstance.IsPostDeploymentEvaluationFailed, reconcilePostEval)
	}

	// WorkloadInstance is completed at this place
	if !workloadInstance.IsEndTimeSet() {
		// metrics: decrement active deployment counter
		r.Meters.DeploymentActive.Add(ctx, -1, workloadInstance.GetActiveMetricsAttributes()...)
		workloadInstance.Status.CurrentPhase = common.PhaseCompleted.ShortName
		workloadInstance.SetEndTime()
	}

	err = r.Client.Status().Update(ctx, workloadInstance)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return ctrl.Result{Requeue: true}, err
	}

	attrs := workloadInstance.GetMetricsAttributes()

	r.Log.Info("Increasing deployment count")
	// metrics: increment deployment counter
	r.Meters.DeploymentCount.Add(ctx, 1, attrs...)

	// metrics: add deployment duration
	duration := workloadInstance.Status.EndTime.Time.Sub(workloadInstance.Status.StartTime.Time)
	r.Meters.DeploymentDuration.Record(ctx, duration.Seconds(), attrs...)

	r.recordEvent(phase, "Normal", workloadInstance, "Finished", "is finished")

	return ctrl.Result{}, nil
}

func (r *KeptnWorkloadInstanceReconciler) handlePhase(ctx context.Context, ctxAppTrace context.Context, workloadInstance *klcv1alpha1.KeptnWorkloadInstance, phase common.KeptnPhaseType, span trace.Span, phaseFailed func() bool, reconcilePhase func() (common.KeptnState, error)) (ctrl.Result, error) {
	r.Log.Info(phase.LongName + " not finished")
	ctx, spanAppTrace := r.getSpan(ctxAppTrace, fmt.Sprintf("%s/%s", workloadInstance.Spec.WorkloadName, phase.ShortName))
	oldPhase := workloadInstance.Status.CurrentPhase
	workloadInstance.Status.CurrentPhase = phase.ShortName
	if phaseFailed() { //TODO eventually we should decide whether a task returns FAILED, currently we never have this status set
		r.recordEvent(phase, "Warning", workloadInstance, "Failed", "has failed")
		return ctrl.Result{Requeue: true, RequeueAfter: 60 * time.Second}, nil
	}
	state, err := reconcilePhase()
	if err != nil {
		spanAppTrace.AddEvent(phase.LongName + " could not get reconciled")
		r.recordEvent(phase, "Warning", workloadInstance, "ReconcileErrored", "could not get reconciled")
		span.SetStatus(codes.Error, err.Error())
		return ctrl.Result{Requeue: true}, err
	}
	if state.IsSucceeded() {
		spanAppTrace.AddEvent(phase.LongName + " has succeeded")
		spanAppTrace.SetStatus(codes.Ok, "Succeeded")
		spanAppTrace.End()
		r.recordEvent(phase, "Normal", workloadInstance, "Succeeded", "has succeeded")
	} else {
		spanAppTrace.AddEvent(phase.LongName + " not finished")
		r.recordEvent(phase, "Warning", workloadInstance, "NotFinished", "has not finished")
	}
	if oldPhase != workloadInstance.Status.CurrentPhase {
		ctx, spanAppTrace = r.getSpan(ctxAppTrace, fmt.Sprintf("%s/%s", workloadInstance.Spec.WorkloadName, workloadInstance.Status.CurrentPhase))
		semconv.AddAttributeFromWorkloadInstance(spanAppTrace, *workloadInstance)

		if err := r.Status().Update(ctx, workloadInstance); err != nil {
			r.Log.Error(err, "could not update status")
		}
	}
	return ctrl.Result{Requeue: true, RequeueAfter: 5 * time.Second}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KeptnWorkloadInstanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		// predicate disabling the auto reconciliation after updating the object status
		For(&klcv1alpha1.KeptnWorkloadInstance{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Complete(r)
}

func (r *KeptnWorkloadInstanceReconciler) generateSuffix() string {
	uid := uuid.New().String()
	return uid[:10]
}

func (r *KeptnWorkloadInstanceReconciler) recordEvent(phase common.KeptnPhaseType, eventType string, workloadInstance *klcv1alpha1.KeptnWorkloadInstance, shortReason string, longReason string) {
	r.Recorder.Event(workloadInstance, eventType, fmt.Sprintf("%s%s", phase.ShortName, shortReason), fmt.Sprintf("%s %s / Namespace: %s, Name: %s, Version: %s ", phase.LongName, longReason, workloadInstance.Namespace, workloadInstance.Name, workloadInstance.Spec.Version))
}

func GetAppVersionName(namespace string, appName string, version string) types.NamespacedName {
	return types.NamespacedName{Namespace: namespace, Name: appName + "-" + version}
}

func (r *KeptnWorkloadInstanceReconciler) getAppVersion(ctx context.Context, appName types.NamespacedName) (*klcv1alpha1.KeptnAppVersion, error) {
	app := &klcv1alpha1.KeptnApp{}
	err := r.Get(ctx, appName, app)
	if err != nil {
		return nil, err
	}

	appVersion := &klcv1alpha1.KeptnAppVersion{}
	err = r.Get(ctx, GetAppVersionName(appName.Namespace, appName.Name, app.Spec.Version), appVersion)
	return appVersion, err
}

func (r *KeptnWorkloadInstanceReconciler) getAppVersionForWorkloadInstance(ctx context.Context, wli *klcv1alpha1.KeptnWorkloadInstance) (bool, klcv1alpha1.KeptnAppVersion, error) {
	apps := &klcv1alpha1.KeptnAppVersionList{}
	if err := r.Client.List(ctx, apps, client.InNamespace(wli.Namespace)); err != nil {
		return false, klcv1alpha1.KeptnAppVersion{}, err
	}
	latestVersion := klcv1alpha1.KeptnAppVersion{}
	for _, app := range apps.Items {
		r.Log.Info("DEBUG: " + app.Name + "==" + wli.Spec.AppName)
		if app.Spec.AppName == wli.Spec.AppName {
			r.Log.Info("DEBUG: FOUND - " + wli.Spec.AppName)

			for _, appWorkload := range app.Spec.Workloads {
				r.Log.Info("DEBUG: " + appWorkload.Version + "==" + wli.Spec.Version + " && " + fmt.Sprintf("%s-%s", app.Spec.AppName, appWorkload.Name) + "==" + wli.Spec.WorkloadName)
				if appWorkload.Version == wli.Spec.Version && fmt.Sprintf("%s-%s", app.Spec.AppName, appWorkload.Name) == wli.Spec.WorkloadName {
					r.Log.Info("DEBUG: FOUND - " + fmt.Sprintf("%s-%s", app.Spec.AppName, appWorkload.Name))
					r.Log.Info("DEBUG: Version - " + app.Spec.Version)
					if latestVersion.Spec.Version == "" {
						r.Log.Info("DEBUG: Set latest version - " + app.Spec.Version)
						latestVersion = app
					} else {
						if semver.Compare(latestVersion.Spec.Version, app.Spec.Version) < 0 {
							r.Log.Info("DEBUG: Set latest version - " + app.Spec.Version)
							latestVersion = app
						}
					}
				}
			}
		}
	}

	r.Log.Info("DEBUG: Selected Version - " + latestVersion.Spec.Version)
	if latestVersion.Spec.Version == "" {
		return false, klcv1alpha1.KeptnAppVersion{}, nil
	}
	return true, latestVersion, nil
}

func (r *KeptnWorkloadInstanceReconciler) getSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	if bindCRDSpan == nil {
		bindCRDSpan = make(map[string]trace.Span)
	}
	if span, ok := bindCRDSpan[name]; ok {
		return ctx, span
	}
	ctx, span := r.AppTracer.Start(ctx, name, trace.WithSpanKind(trace.SpanKindConsumer))
	bindCRDSpan[name] = span
	return ctx, span
}
