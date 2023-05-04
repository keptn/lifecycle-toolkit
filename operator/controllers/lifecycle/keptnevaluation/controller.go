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

package keptnevaluation

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3/common"
	controllercommon "github.com/keptn/lifecycle-toolkit/operator/controllers/common"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/common/providers/keptnmetric"
	controllererrors "github.com/keptn/lifecycle-toolkit/operator/controllers/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

const traceComponentName = "keptn/operator/evaluation"

// KeptnEvaluationReconciler reconciles a KeptnEvaluation object
type KeptnEvaluationReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	Recorder      record.EventRecorder
	Log           logr.Logger
	Meters        apicommon.KeptnMeters
	TracerFactory controllercommon.TracerFactory
	Namespace     string
}

// clusterrole
// +kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnevaluations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnevaluations/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnevaluations/finalizers,verbs=update
// +kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnevaluationproviders,verbs=get;list;watch
// +kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnevaluationdefinitions,verbs=get;list;watch
// +kubebuilder:rbac:groups=metrics.keptn.sh,resources=keptnmetrics,verbs=get;list;watch

// role
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *KeptnEvaluationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	r.Log.Info("Reconciling KeptnEvaluation")
	evaluation := &klcv1alpha3.KeptnEvaluation{}

	if err := r.Client.Get(ctx, req.NamespacedName, evaluation); err != nil {
		if errors.IsNotFound(err) {
			// taking down all associated K8s resources is handled by K8s
			r.Log.Info("KeptnEvaluation resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		r.Log.Error(err, "Failed to get the KeptnEvaluation")
		return ctrl.Result{}, nil
	}

	ctx, span := r.setupEvaluationSpans(ctx, evaluation)
	defer span.End()

	if evaluation.Status.RetryCount >= evaluation.Spec.Retries {
		r.handleEvaluationExceededRetries(ctx, evaluation, span)
		return ctrl.Result{}, nil
	}

	if !evaluation.Status.OverallStatus.IsSucceeded() {
		evaluationDefinition, err := controllercommon.GetEvaluationDefinition(r.Client, r.Log, ctx, evaluation.Spec.EvaluationDefinition, req.NamespacedName.Namespace)
		if err != nil {
			if errors.IsNotFound(err) {
				r.Log.Info(err.Error() + ", ignoring error since object must be deleted")
				span.SetStatus(codes.Error, err.Error())
				return ctrl.Result{Requeue: true, RequeueAfter: 10 * time.Second}, nil
			}
			r.Log.Error(err, "Failed to retrieve a resource")
			span.SetStatus(codes.Error, err.Error())
			return ctrl.Result{}, nil
		}

		evaluation = r.performEvaluation(ctx, evaluation, evaluationDefinition)

	}

	if !evaluation.Status.OverallStatus.IsSucceeded() {
		if err := r.handleEvaluationIncomplete(ctx, evaluation, span); err != nil {
			return ctrl.Result{Requeue: true}, err
		}
		return ctrl.Result{Requeue: true, RequeueAfter: evaluation.Spec.RetryInterval.Duration}, nil
	}

	r.Log.Info("Finished Reconciling KeptnEvaluation")

	err := r.updateFinishedEvaluationMetrics(ctx, evaluation, span)

	return ctrl.Result{}, err

}

func (r *KeptnEvaluationReconciler) setupEvaluationSpans(ctx context.Context, evaluation *klcv1alpha3.KeptnEvaluation) (context.Context, trace.Span) {
	traceContextCarrier := propagation.MapCarrier(evaluation.Annotations)
	ctx = otel.GetTextMapPropagator().Extract(ctx, traceContextCarrier)
	ctx, span := r.getTracer().Start(ctx, "reconcile_evaluation", trace.WithSpanKind(trace.SpanKindConsumer))
	evaluation.SetSpanAttributes(span)
	evaluation.SetStartTime()

	return ctx, span
}

func (r *KeptnEvaluationReconciler) handleEvaluationIncomplete(ctx context.Context, evaluation *klcv1alpha3.KeptnEvaluation, span trace.Span) error {
	// Evaluation is uncompleted, update status anyway this avoids updating twice in case of completion
	err := r.Client.Status().Update(ctx, evaluation)
	if err != nil {
		controllercommon.RecordEvent(r.Recorder, apicommon.PhaseReconcileEvaluation, "Warning", evaluation, "ReconcileErrored", "could not update status", "")
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	controllercommon.RecordEvent(r.Recorder, apicommon.PhaseReconcileEvaluation, "Normal", evaluation, "NotFinished", "has not finished", "")

	return nil

}

func (r *KeptnEvaluationReconciler) handleEvaluationExceededRetries(ctx context.Context, evaluation *klcv1alpha3.KeptnEvaluation, span trace.Span) {
	controllercommon.RecordEvent(r.Recorder, apicommon.PhaseReconcileEvaluation, "Warning", evaluation, "ReconcileTimeOut", "retryCount exceeded", "")
	err := controllererrors.ErrRetryCountExceeded
	span.SetStatus(codes.Error, err.Error())
	evaluation.Status.OverallStatus = apicommon.StateFailed
	err2 := r.updateFinishedEvaluationMetrics(ctx, evaluation, span)
	if err2 != nil {
		r.Log.Error(err2, "failed to update finished evaluation metrics")
	}
}

func (r *KeptnEvaluationReconciler) performEvaluation(ctx context.Context, evaluation *klcv1alpha3.KeptnEvaluation, evaluationDefinition *klcv1alpha3.KeptnEvaluationDefinition) *klcv1alpha3.KeptnEvaluation {
	statusSummary := apicommon.StatusSummary{Total: len(evaluationDefinition.Spec.Objectives)}
	newStatus := make(map[string]klcv1alpha3.EvaluationStatusItem)

	if evaluation.Status.EvaluationStatus == nil {
		evaluation.Status.EvaluationStatus = make(map[string]klcv1alpha3.EvaluationStatusItem)
	}

	provider := &keptnmetric.KeptnMetricProvider{
		Log:       r.Log,
		K8sClient: r.Client,
	}

	for _, query := range evaluationDefinition.Spec.Objectives {
		newStatus, statusSummary = r.evaluateObjective(ctx, evaluation, statusSummary, newStatus, query, provider)
	}

	evaluation.Status.RetryCount++
	evaluation.Status.EvaluationStatus = newStatus
	if apicommon.GetOverallState(statusSummary) == apicommon.StateSucceeded {
		evaluation.Status.OverallStatus = apicommon.StateSucceeded
	} else {
		evaluation.Status.OverallStatus = apicommon.StateProgressing
	}

	return evaluation
}

func (r *KeptnEvaluationReconciler) evaluateObjective(ctx context.Context, evaluation *klcv1alpha3.KeptnEvaluation, statusSummary apicommon.StatusSummary, newStatus map[string]klcv1alpha3.EvaluationStatusItem, objective klcv1alpha3.Objective, provider *keptnmetric.KeptnMetricProvider) (map[string]klcv1alpha3.EvaluationStatusItem, apicommon.StatusSummary) {
	if _, ok := evaluation.Status.EvaluationStatus[objective.KeptnMetricRef.Name]; !ok {
		evaluation.AddEvaluationStatus(objective)
	}
	if evaluation.Status.EvaluationStatus[objective.KeptnMetricRef.Name].Status.IsSucceeded() {
		statusSummary = apicommon.UpdateStatusSummary(apicommon.StateSucceeded, statusSummary)
		newStatus[objective.KeptnMetricRef.Name] = evaluation.Status.EvaluationStatus[objective.KeptnMetricRef.Name]
		return newStatus, statusSummary
	}
	// resolving the SLI value
	value, _, err := provider.FetchData(ctx, objective, evaluation.Namespace)
	statusItem := &klcv1alpha3.EvaluationStatusItem{
		Value:  value,
		Status: apicommon.StateFailed,
	}
	if err != nil {
		statusItem.Message = err.Error()
	}
	// Evaluating SLO
	check, err := checkValue(objective, statusItem)
	if err != nil {
		statusItem.Message = err.Error()
		r.Log.Error(err, "Could not check objective result")
	} else {
		// if there is no error, we set the message depending on if the value passed the objective, or not
		if check {
			statusItem.Status = apicommon.StateSucceeded
			statusItem.Message = fmt.Sprintf("value '%s' met objective '%s'", value, objective.EvaluationTarget)
		} else {
			statusItem.Message = fmt.Sprintf("value '%s' did not meet objective '%s'", value, objective.EvaluationTarget)
		}
	}
	statusSummary = apicommon.UpdateStatusSummary(statusItem.Status, statusSummary)
	newStatus[objective.KeptnMetricRef.Name] = *statusItem

	return newStatus, statusSummary
}

func (r *KeptnEvaluationReconciler) updateFinishedEvaluationMetrics(ctx context.Context, evaluation *klcv1alpha3.KeptnEvaluation, span trace.Span) error {
	evaluation.SetEndTime()

	err := r.Client.Status().Update(ctx, evaluation)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		controllercommon.RecordEvent(r.Recorder, apicommon.PhaseReconcileEvaluation, "Warning", evaluation, "ReconcileErrored", "could not update status", "")
		return err
	}

	attrs := evaluation.GetMetricsAttributes()

	r.Log.Info("Increasing evaluation count")

	// metrics: increment evaluation counter
	r.Meters.EvaluationCount.Add(ctx, 1, attrs...)

	// metrics: add evaluation duration
	duration := evaluation.Status.EndTime.Time.Sub(evaluation.Status.StartTime.Time)
	r.Meters.EvaluationDuration.Record(ctx, duration.Seconds(), attrs...)
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KeptnEvaluationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&klcv1alpha3.KeptnEvaluation{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Complete(r)
}

func (r *KeptnEvaluationReconciler) getTracer() controllercommon.ITracer {
	return r.TracerFactory.GetTracer(traceComponentName)
}
