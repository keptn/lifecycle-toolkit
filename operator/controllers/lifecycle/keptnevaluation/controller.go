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
	klcv1alpha2 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2/common"
	controllererrors "github.com/keptn/lifecycle-toolkit/operator/controllers/errors"
	providers "github.com/keptn/lifecycle-toolkit/operator/controllers/lifecycle/keptnevaluation/providers"
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
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// KeptnEvaluationReconciler reconciles a KeptnEvaluation object
type KeptnEvaluationReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
	Log      logr.Logger
	Meters   apicommon.KeptnMeters
	Tracer   trace.Tracer
}

//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnevaluations,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnevaluations/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnevaluations/finalizers,verbs=update
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnevaluationproviders,verbs=get;list;watch
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnevaluationdefinitions,verbs=get;list;watch
//+kubebuilder:rbac:groups=core,resources=secrets,verbs=get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the KeptnEvaluation object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *KeptnEvaluationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	r.Log.Info("Reconciling KeptnEvaluation")
	evaluation := &klcv1alpha2.KeptnEvaluation{}

	if err := r.Client.Get(ctx, req.NamespacedName, evaluation); err != nil {
		if errors.IsNotFound(err) {
			// taking down all associated K8s resources is handled by K8s
			r.Log.Info("KeptnEvaluation resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		r.Log.Error(err, "Failed to get the KeptnEvaluation")
		return ctrl.Result{}, nil
	}

	traceContextCarrier := propagation.MapCarrier(evaluation.Annotations)
	ctx = otel.GetTextMapPropagator().Extract(ctx, traceContextCarrier)

	ctx, span := r.Tracer.Start(ctx, "reconcile_evaluation", trace.WithSpanKind(trace.SpanKindConsumer))
	defer span.End()

	evaluation.SetSpanAttributes(span)

	evaluation.SetStartTime()

	if evaluation.Status.RetryCount >= evaluation.Spec.Retries {
		r.recordEvent("Warning", evaluation, "ReconcileTimeOut", "retryCount exceeded")
		err := controllererrors.ErrRetryCountExceeded
		span.SetStatus(codes.Error, err.Error())
		evaluation.Status.OverallStatus = apicommon.StateFailed
		err2 := r.updateFinishedEvaluationMetrics(ctx, evaluation, span)
		if err2 != nil {
			r.Log.Error(err2, "failed to update finished evaluation metrics")
		}
		return ctrl.Result{}, nil
	}

	if !evaluation.Status.OverallStatus.IsSucceeded() {
		namespacedDefinition := types.NamespacedName{
			Namespace: req.NamespacedName.Namespace,
			Name:      evaluation.Spec.EvaluationDefinition,
		}
		evaluationDefinition, evaluationProvider, err := r.fetchDefinitionAndProvider(ctx, namespacedDefinition)
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
		// load the provider
		provider, err2 := providers.NewProvider(evaluationDefinition.Spec.Source, r.Log, r.Client)
		if err2 != nil {
			r.recordEvent("Error", evaluation, "ProviderNotFound", "evaluation provider was not found")
			r.Log.Error(err2, "Failed to get the correct Metric Provider")
			span.SetStatus(codes.Error, err2.Error())
			return ctrl.Result{Requeue: false}, err2
		}

		evaluation = r.performEvaluation(ctx, evaluation, evaluationDefinition, provider, evaluationProvider)

	}

	if !evaluation.Status.OverallStatus.IsSucceeded() {
		// Evaluation is uncompleted, update status anyway this avoids updating twice in case of completion
		err := r.Client.Status().Update(ctx, evaluation)
		if err != nil {
			r.recordEvent("Warning", evaluation, "ReconcileErrored", "could not update status")
			span.SetStatus(codes.Error, err.Error())
			return ctrl.Result{Requeue: true}, err
		}

		r.recordEvent("Normal", evaluation, "NotFinished", "has not finished")

		return ctrl.Result{Requeue: true, RequeueAfter: evaluation.Spec.RetryInterval.Duration}, nil

	}

	r.Log.Info("Finished Reconciling KeptnEvaluation")

	err := r.updateFinishedEvaluationMetrics(ctx, evaluation, span)

	return ctrl.Result{}, err

}

func (r *KeptnEvaluationReconciler) performEvaluation(ctx context.Context, evaluation *klcv1alpha2.KeptnEvaluation, evaluationDefinition *klcv1alpha2.KeptnEvaluationDefinition, provider providers.KeptnSLIProvider, evaluationProvider *klcv1alpha2.KeptnEvaluationProvider) *klcv1alpha2.KeptnEvaluation {
	statusSummary := apicommon.StatusSummary{Total: len(evaluationDefinition.Spec.Objectives)}
	newStatus := make(map[string]klcv1alpha2.EvaluationStatusItem)

	if evaluation.Status.EvaluationStatus == nil {
		evaluation.Status.EvaluationStatus = make(map[string]klcv1alpha2.EvaluationStatusItem)
	}

	for _, query := range evaluationDefinition.Spec.Objectives {
		newStatus, statusSummary = r.evaluateObjective(ctx, evaluation, statusSummary, newStatus, query, provider, evaluationProvider)
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

func (r *KeptnEvaluationReconciler) evaluateObjective(ctx context.Context, evaluation *klcv1alpha2.KeptnEvaluation, statusSummary apicommon.StatusSummary, newStatus map[string]klcv1alpha2.EvaluationStatusItem, query klcv1alpha2.Objective, provider providers.KeptnSLIProvider, evaluationProvider *klcv1alpha2.KeptnEvaluationProvider) (map[string]klcv1alpha2.EvaluationStatusItem, apicommon.StatusSummary) {
	if _, ok := evaluation.Status.EvaluationStatus[query.Name]; !ok {
		evaluation.AddEvaluationStatus(query)
	}
	if evaluation.Status.EvaluationStatus[query.Name].Status.IsSucceeded() {
		statusSummary = apicommon.UpdateStatusSummary(apicommon.StateSucceeded, statusSummary)
		newStatus[query.Name] = evaluation.Status.EvaluationStatus[query.Name]
		return newStatus, statusSummary
	}
	// resolving the SLI value
	value, err := provider.EvaluateQuery(ctx, query, *evaluationProvider)
	statusItem := &klcv1alpha2.EvaluationStatusItem{
		Value:  value,
		Status: apicommon.StateFailed,
	}
	if err != nil {
		statusItem.Message = err.Error()
		statusItem.Status = apicommon.StateFailed
	}
	// Evaluating SLO
	check, err := checkValue(query, statusItem)
	if err != nil {
		statusItem.Message = err.Error()
		r.Log.Error(err, "Could not check query result")
	}
	if check {
		statusItem.Status = apicommon.StateSucceeded
	}
	statusSummary = apicommon.UpdateStatusSummary(statusItem.Status, statusSummary)
	newStatus[query.Name] = *statusItem

	return newStatus, statusSummary
}

func (r *KeptnEvaluationReconciler) updateFinishedEvaluationMetrics(ctx context.Context, evaluation *klcv1alpha2.KeptnEvaluation, span trace.Span) error {
	evaluation.SetEndTime()

	err := r.Client.Status().Update(ctx, evaluation)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		r.recordEvent("Warning", evaluation, "ReconcileErrored", "could not update status")
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
		For(&klcv1alpha2.KeptnEvaluation{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Complete(r)
}

func (r *KeptnEvaluationReconciler) fetchDefinitionAndProvider(ctx context.Context, namespacedDefinition types.NamespacedName) (*klcv1alpha2.KeptnEvaluationDefinition, *klcv1alpha2.KeptnEvaluationProvider, error) {
	evaluationDefinition := &klcv1alpha2.KeptnEvaluationDefinition{}
	if err := r.Client.Get(ctx, namespacedDefinition, evaluationDefinition); err != nil {
		return nil, nil, err
	}
	namespacedProvider := types.NamespacedName{
		Namespace: namespacedDefinition.Namespace,
		Name:      evaluationDefinition.Spec.Source,
	}
	evaluationProvider := &klcv1alpha2.KeptnEvaluationProvider{}
	if err := r.Client.Get(ctx, namespacedProvider, evaluationProvider); err != nil {
		return nil, nil, err
	}
	return evaluationDefinition, evaluationProvider, nil
}

func (r *KeptnEvaluationReconciler) recordEvent(eventType string, evaluation *klcv1alpha2.KeptnEvaluation, shortReason string, longReason string) {
	r.Recorder.Event(evaluation, eventType, shortReason, fmt.Sprintf("%s / Namespace: %s, Name: %s, WorkloadVersion: %s ", longReason, evaluation.Namespace, evaluation.Name, evaluation.Spec.WorkloadVersion))
}
