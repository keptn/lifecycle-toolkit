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
	"reflect"
	"time"

	"golang.org/x/mod/semver"

	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1/semconv"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"github.com/go-logr/logr"
	version "github.com/hashicorp/go-version"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	klcv1alpha1 "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1"
	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1/common"
	controllercommon "github.com/keptn/lifecycle-toolkit/operator/controllers/common"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// KeptnWorkloadInstanceReconciler reconciles a KeptnWorkloadInstance object
type KeptnWorkloadInstanceReconciler struct {
	client.Client
	Scheme      *runtime.Scheme
	Recorder    record.EventRecorder
	Log         logr.Logger
	Meters      common.KeptnMeters
	Tracer      trace.Tracer
	SpanHandler controllercommon.SpanHandler
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

	//setup otel
	traceContextCarrier := propagation.MapCarrier(workloadInstance.Annotations)
	ctx = otel.GetTextMapPropagator().Extract(ctx, traceContextCarrier)

	ctx, span := r.Tracer.Start(ctx, "reconcile_workload_instance", trace.WithSpanKind(trace.SpanKindConsumer))
	defer span.End()

	semconv.AddAttributeFromWorkloadInstance(span, *workloadInstance)

	workloadInstance.SetStartTime()

	defer func(span trace.Span, workloadInstance *klcv1alpha1.KeptnWorkloadInstance) {
		if workloadInstance.IsEndTimeSet() {
			r.Log.Info("Increasing deployment count")
			attrs := workloadInstance.GetMetricsAttributes()
			r.Meters.AppCount.Add(ctx, 1, attrs...)
		}
		span.End()
	}(span, workloadInstance)

	//Wait for pre-evaluation checks of App
	phase := common.PhaseAppPreEvaluation

	found, appVersion, err := r.getAppVersionForWorkloadInstance(ctx, workloadInstance)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		controllercommon.RecordEvent(r.Recorder, phase, "Warning", workloadInstance, "GetAppVersionFailed", "has failed since app could not be retrieved", workloadInstance.GetVersion())
		return reconcile.Result{Requeue: true, RequeueAfter: 10 * time.Second}, fmt.Errorf("could not fetch AppVersion for KeptnWorkloadInstance: %+v", err)
	} else if !found {
		span.SetStatus(codes.Error, "app could not be found")
		controllercommon.RecordEvent(r.Recorder, phase, "Warning", workloadInstance, "AppVersionNotFound", "has failed since app could not be found", workloadInstance.GetVersion())
		return reconcile.Result{Requeue: true, RequeueAfter: 10 * time.Second}, fmt.Errorf("could not find AppVersion for KeptnWorkloadInstance")
	}

	appTraceContextCarrier := propagation.MapCarrier(appVersion.Spec.TraceId)
	ctxAppTrace := otel.GetTextMapPropagator().Extract(context.TODO(), appTraceContextCarrier)

	appPreEvalStatus := appVersion.Status.PreDeploymentEvaluationStatus
	if !appPreEvalStatus.IsSucceeded() {
		if appPreEvalStatus.IsFailed() {
			controllercommon.RecordEvent(r.Recorder, phase, "Warning", workloadInstance, "Failed", "has failed since app has failed", workloadInstance.GetVersion())
			return ctrl.Result{Requeue: true, RequeueAfter: 20 * time.Second}, nil
		}
		controllercommon.RecordEvent(r.Recorder, phase, "Normal", workloadInstance, "NotFinished", "Pre evaluations tasks for app not finished", workloadInstance.GetVersion())
		return ctrl.Result{Requeue: true, RequeueAfter: 20 * time.Second}, nil
	}

	controllercommon.RecordEvent(r.Recorder, phase, "Normal", workloadInstance, "FinishedSuccess", "Pre evaluations tasks for app have finished successfully", workloadInstance.GetVersion())

	//Wait for pre-deployment checks of Workload
	phase = common.PhaseWorkloadPreDeployment
	phaseHandler := controllercommon.PhaseHandler{
		Client:      r.Client,
		Recorder:    r.Recorder,
		Log:         r.Log,
		SpanHandler: r.SpanHandler,
	}

	// set the App trace id if not already set
	if len(workloadInstance.Spec.TraceId) < 1 {
		workloadInstance.Spec.TraceId = appVersion.Spec.TraceId
		if err := r.Update(ctx, workloadInstance); err != nil {
			return ctrl.Result{}, err
		}
	}

	if workloadInstance.Status.CurrentPhase == "" {
		if err := r.SpanHandler.UnbindSpan(workloadInstance, phase.ShortName); err != nil {
			r.Log.Error(err, "cannot unbind span")
		}
		var spanAppTrace trace.Span
		ctxAppTrace, spanAppTrace, err = r.SpanHandler.GetSpan(ctxAppTrace, r.Tracer, workloadInstance, phase.ShortName)
		if err != nil {
			r.Log.Error(err, "could not get span")
		}
		semconv.AddAttributeFromWorkloadInstance(spanAppTrace, *workloadInstance)
		spanAppTrace.AddEvent("WorkloadInstance Pre-Deployment Tasks started", trace.WithTimestamp(time.Now()))
		controllercommon.RecordEvent(r.Recorder, phase, "Normal", workloadInstance, "Started", "have started", workloadInstance.GetVersion())
	}

	if !workloadInstance.IsPreDeploymentSucceeded() {
		reconcilePre := func() (common.KeptnState, error) {
			return r.reconcilePrePostDeployment(ctx, workloadInstance, common.PreDeploymentCheckType)
		}
		result, err := phaseHandler.HandlePhase(ctx, ctxAppTrace, r.Tracer, workloadInstance, phase, span, reconcilePre)
		if !result.Continue {
			return result.Result, err
		}
	}

	//Wait for pre-evaluation checks of Workload
	phase = common.PhaseAppPreEvaluation
	if !workloadInstance.IsPreDeploymentEvaluationSucceeded() {
		reconcilePreEval := func() (common.KeptnState, error) {
			return r.reconcilePrePostEvaluation(ctx, workloadInstance, common.PreDeploymentEvaluationCheckType)
		}
		result, err := phaseHandler.HandlePhase(ctx, ctxAppTrace, r.Tracer, workloadInstance, phase, span, reconcilePreEval)
		if !result.Continue {
			return result.Result, err
		}
	}

	//Wait for deployment of Workload
	phase = common.PhaseWorkloadDeployment
	if !workloadInstance.IsDeploymentSucceeded() {
		reconcileWorkloadInstance := func() (common.KeptnState, error) {
			return r.reconcileDeployment(ctx, workloadInstance)
		}
		result, err := phaseHandler.HandlePhase(ctx, ctxAppTrace, r.Tracer, workloadInstance, phase, span, reconcileWorkloadInstance)
		if !result.Continue {
			return result.Result, err
		}
	}

	//Wait for post-deployment checks of Workload
	phase = common.PhaseWorkloadPostDeployment
	if !workloadInstance.IsPostDeploymentSucceeded() {
		reconcilePostDeployment := func() (common.KeptnState, error) {
			return r.reconcilePrePostDeployment(ctx, workloadInstance, common.PostDeploymentCheckType)
		}
		result, err := phaseHandler.HandlePhase(ctx, ctxAppTrace, r.Tracer, workloadInstance, phase, span, reconcilePostDeployment)
		if !result.Continue {
			return result.Result, err
		}
	}

	//Wait for post-evaluation checks of Workload
	phase = common.PhaseAppPostEvaluation
	if !workloadInstance.IsPostDeploymentEvaluationSucceeded() {
		reconcilePostEval := func() (common.KeptnState, error) {
			return r.reconcilePrePostEvaluation(ctx, workloadInstance, common.PostDeploymentEvaluationCheckType)
		}
		result, err := phaseHandler.HandlePhase(ctx, ctxAppTrace, r.Tracer, workloadInstance, phase, span, reconcilePostEval)
		if !result.Continue {
			return result.Result, err
		}
	}

	// WorkloadInstance is completed at this place
	if !workloadInstance.IsEndTimeSet() {
		workloadInstance.Status.CurrentPhase = common.PhaseCompleted.ShortName
		workloadInstance.Status.Status = common.StateSucceeded
		workloadInstance.SetEndTime()
	}

	err = r.Client.Status().Update(ctx, workloadInstance)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return ctrl.Result{Requeue: true}, err
	}

	attrs := workloadInstance.GetMetricsAttributes()

	// metrics: add deployment duration
	duration := workloadInstance.Status.EndTime.Time.Sub(workloadInstance.Status.StartTime.Time)
	r.Meters.DeploymentDuration.Record(ctx, duration.Seconds(), attrs...)

	controllercommon.RecordEvent(r.Recorder, phase, "Normal", workloadInstance, "Finished", "is finished", workloadInstance.GetVersion())

	return ctrl.Result{}, nil
}

func (r *KeptnWorkloadInstanceReconciler) GetActiveDeployments(ctx context.Context) ([]common.GaugeValue, error) {
	workloadInstances := &klcv1alpha1.KeptnWorkloadInstanceList{}
	err := r.List(ctx, workloadInstances)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve workload instances: %w", err)
	}

	res := []common.GaugeValue{}

	for _, workloadInstance := range workloadInstances.Items {
		gaugeValue := int64(0)
		if !workloadInstance.IsEndTimeSet() {
			gaugeValue = int64(1)
		}
		res = append(res, common.GaugeValue{
			Value:      gaugeValue,
			Attributes: workloadInstance.GetActiveMetricsAttributes(),
		})
	}

	return res, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KeptnWorkloadInstanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		// predicate disabling the auto reconciliation after updating the object status
		For(&klcv1alpha1.KeptnWorkloadInstance{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Complete(r)
}

func (r *KeptnWorkloadInstanceReconciler) getAppVersion(ctx context.Context, appName types.NamespacedName) (*klcv1alpha1.KeptnAppVersion, error) {
	app := &klcv1alpha1.KeptnApp{}
	err := r.Get(ctx, appName, app)
	if err != nil {
		return nil, err
	}

	appVersion := &klcv1alpha1.KeptnAppVersion{}
	err = r.Get(ctx, controllercommon.GetAppVersionName(appName.Namespace, appName.Name, app.Spec.Version), appVersion)
	return appVersion, err
}

func (r *KeptnWorkloadInstanceReconciler) getAppVersionForWorkloadInstance(ctx context.Context, wli *klcv1alpha1.KeptnWorkloadInstance) (bool, klcv1alpha1.KeptnAppVersion, error) {
	apps := &klcv1alpha1.KeptnAppVersionList{}

	if err := r.Client.List(ctx, apps, client.InNamespace(wli.Namespace)); err != nil {
		return false, klcv1alpha1.KeptnAppVersion{}, err
	}
	latestVersion := klcv1alpha1.KeptnAppVersion{}
	for _, app := range apps.Items {
		if app.Spec.AppName == wli.Spec.AppName {

			for _, appWorkload := range app.Spec.Workloads {
				if !reflect.DeepEqual(latestVersion, app) {
					latestVersion = app
				} else if appWorkload.Version == wli.Spec.Version && fmt.Sprintf("%s-%s", app.Spec.AppName, appWorkload.Name) == wli.Spec.WorkloadName {
					oldVersion, err := version.NewVersion(app.Spec.Version)
					if err != nil {
						r.Log.Error(err, "could not parse version")
					}
					newVersion, err := version.NewVersion(latestVersion.Spec.Version)
					if err != nil {
						r.Log.Error(err, "could not parse version")
					}
					if oldVersion.LessThan(newVersion) {
						latestVersion = app
					}
				}
			}
		}
	}

	if latestVersion.Spec.Version == "" {
		return false, klcv1alpha1.KeptnAppVersion{}, nil
	}
	return true, latestVersion, nil
}

func (r *KeptnWorkloadInstanceReconciler) GetDeploymentInterval(ctx context.Context) ([]common.GaugeFloatValue, error) {
	workloadInstances := &klcv1alpha1.KeptnWorkloadInstanceList{}
	err := r.List(ctx, workloadInstances)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve workload instances: %w", err)
	}

	res := []common.GaugeFloatValue{}
	for _, workloadInstance := range workloadInstances.Items {
		if workloadInstance.Spec.PreviousVersion != "" {
			previousWorkloadInstance := &klcv1alpha1.KeptnWorkloadInstance{}
			err := r.Get(ctx, types.NamespacedName{Name: fmt.Sprintf("%s-%s", workloadInstance.Spec.WorkloadName, workloadInstance.Spec.PreviousVersion), Namespace: workloadInstance.Namespace}, previousWorkloadInstance)
			if err != nil {
				r.Log.Error(err, "Previous WorkloadInstance not found")
			} else if workloadInstance.IsEndTimeSet() {
				previousInterval := workloadInstance.Status.StartTime.Time.Sub(previousWorkloadInstance.Status.EndTime.Time)
				res = append(res, common.GaugeFloatValue{
					Value:      previousInterval.Seconds(),
					Attributes: workloadInstance.GetIntervalMetricsAttributes(),
				})
			}
		}
	}
	return res, nil
}

func (r *KeptnWorkloadInstanceReconciler) GetDeploymentDuration(ctx context.Context) ([]common.GaugeFloatValue, error) {
	workloadInstances := &klcv1alpha1.KeptnWorkloadInstanceList{}
	err := r.List(ctx, workloadInstances)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve workload instances: %w", err)
	}

	res := []common.GaugeFloatValue{}

	for _, workloadInstance := range workloadInstances.Items {
		if workloadInstance.IsEndTimeSet() {
			duration := workloadInstance.Status.EndTime.Time.Sub(workloadInstance.Status.StartTime.Time)
			res = append(res, common.GaugeFloatValue{
				Value:      duration.Seconds(),
				Attributes: workloadInstance.GetIntervalMetricsAttributes(),
			})
		}
	}

	return res, nil
}
