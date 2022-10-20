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
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
	Log      logr.Logger
	Meters   common.KeptnMeters
	Tracer   trace.Tracer
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

	//Wait for pre-evaluation checks of App
	phase := common.PhaseAppPreEvaluation

	appVersion, err := r.getAppVersion(ctx, types.NamespacedName{Namespace: req.Namespace, Name: workloadInstance.Spec.AppName})
	if errors.IsNotFound(err) {
		r.recordEvent(phase, "Warning", workloadInstance, "AppVersionNotFound", "has failed since app could not be found")
		r.Log.Error(err, "Related App Version not found")
		return reconcile.Result{Requeue: true, RequeueAfter: 20 * time.Second}, nil
	} else if err != nil {
		r.recordEvent(phase, "Warning", workloadInstance, "GetAppVersionFailed", "has failed since app could not be retrieved")
		r.Log.Error(err, "Could not retrieve App Version")
		return reconcile.Result{}, fmt.Errorf("could not fetch AppVersion: %+v", err)
	}

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
	if !workloadInstance.IsPreDeploymentSucceeded() {
		reconcilePre := func() (common.KeptnState, error) {
			return r.reconcilePrePostDeployment(ctx, workloadInstance, common.PreDeploymentCheckType)
		}
		return r.handlePhase(ctx, workloadInstance, phase, span, workloadInstance.IsPreDeploymentFailed, reconcilePre)
	}

	//Wait for pre-evaluation checks of Workload
	phase = common.PhaseAppPreEvaluation
	if !workloadInstance.IsPreDeploymentEvaluationSucceeded() {
		reconcilePreEval := func() (common.KeptnState, error) {
			return r.reconcilePrePostEvaluation(ctx, workloadInstance, common.PreDeploymentEvaluationCheckType)
		}
		return r.handlePhase(ctx, workloadInstance, phase, span, workloadInstance.IsPreDeploymentEvaluationFailed, reconcilePreEval)
	}

	//Wait for deployment of Workload
	phase = common.PhaseWorkloadDeployment
	if !workloadInstance.IsDeploymentSucceeded() {
		reconcileWorkloadInstance := func() (common.KeptnState, error) {
			return r.reconcileDeployment(ctx, workloadInstance)
		}
		return r.handlePhase(ctx, workloadInstance, phase, span, workloadInstance.IsDeploymentFailed, reconcileWorkloadInstance)
	}

	//Wait for post-deployment checks of Workload
	phase = common.PhaseWorkloadPostDeployment
	if !workloadInstance.IsPostDeploymentSucceeded() {
		reconcilePostDeployment := func() (common.KeptnState, error) {
			return r.reconcilePrePostDeployment(ctx, workloadInstance, common.PostDeploymentCheckType)
		}
		return r.handlePhase(ctx, workloadInstance, phase, span, workloadInstance.IsPostDeploymentFailed, reconcilePostDeployment)
	}

	//Wait for post-evaluation checks of Workload
	phase = common.PhaseAppPostEvaluation
	if !workloadInstance.IsPostDeploymentEvaluationSucceeded() {
		reconcilePostEval := func() (common.KeptnState, error) {
			return r.reconcilePrePostEvaluation(ctx, workloadInstance, common.PostDeploymentEvaluationCheckType)
		}
		return r.handlePhase(ctx, workloadInstance, phase, span, workloadInstance.IsPostDeploymentEvaluationFailed, reconcilePostEval)
	}

	// WorkloadInstance is completed at this place
	if !workloadInstance.IsEndTimeSet() {
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

func (r *KeptnWorkloadInstanceReconciler) handlePhase(ctx context.Context, workloadInstance *klcv1alpha1.KeptnWorkloadInstance, phase common.KeptnPhaseType, span trace.Span, phaseFailed func() bool, reconcilePhase func() (common.KeptnState, error)) (ctrl.Result, error) {
	r.Log.Info(phase.LongName + " not finished")
	oldPhase := workloadInstance.Status.CurrentPhase
	workloadInstance.Status.CurrentPhase = phase.ShortName
	if phaseFailed() { //TODO eventually we should decide whether a task returns FAILED, currently we never have this status set
		r.recordEvent(phase, "Warning", workloadInstance, "Failed", "has failed")
		return ctrl.Result{Requeue: true, RequeueAfter: 60 * time.Second}, nil
	}
	state, err := reconcilePhase()
	if err != nil {
		r.recordEvent(phase, "Warning", workloadInstance, "ReconcileErrored", "could not get reconciled")
		span.SetStatus(codes.Error, err.Error())
		return ctrl.Result{Requeue: true}, err
	}
	if state.IsSucceeded() {
		r.recordEvent(phase, "Normal", workloadInstance, "Succeeded", "has succeeded")
	} else {
		r.recordEvent(phase, "Warning", workloadInstance, "NotFinished", "has not finished")
	}
	if oldPhase != workloadInstance.Status.CurrentPhase {
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
