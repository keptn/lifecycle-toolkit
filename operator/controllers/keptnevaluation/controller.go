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
	"math"
	"net/http"
	"strconv"
	"time"

	promapi "github.com/prometheus/client_golang/api"
	prometheus "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
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

	"github.com/go-logr/logr"
	klcv1alpha1 "github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1"
	"github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1/common"
	"github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1/semconv"
)

// KeptnEvaluationReconciler reconciles a KeptnEvaluation object
type KeptnEvaluationReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
	Log      logr.Logger
	Meters   common.KeptnMeters
	Tracer   trace.Tracer
}

//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnevaluations,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnevaluations/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnevaluations/finalizers,verbs=update
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnevaluationproviders,verbs=get;list;watch
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnevaluationdefinitions,verbs=get;list;watch

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
	evaluation := &klcv1alpha1.KeptnEvaluation{}

	if err := r.Client.Get(ctx, req.NamespacedName, evaluation); err != nil {
		if errors.IsNotFound(err) {
			// taking down all associated K8s resources is handled by K8s
			r.Log.Info("KeptnEvaluation resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		r.Log.Error(err, "Failed to get the KeptnEvaluation")
		return ctrl.Result{Requeue: true, RequeueAfter: 30 * time.Second}, nil
	}

	traceContextCarrier := propagation.MapCarrier(evaluation.Annotations)
	ctx = otel.GetTextMapPropagator().Extract(ctx, traceContextCarrier)

	ctx, span := r.Tracer.Start(ctx, "reconcile_evaluation", trace.WithSpanKind(trace.SpanKindConsumer))
	defer span.End()

	semconv.AddAttributeFromEvaluation(span, *evaluation)

	if !evaluation.IsStartTimeSet() {
		// metrics: increment active evaluation counter
		r.Meters.AnalysisActive.Add(ctx, 1, evaluation.GetActiveMetricsAttributes()...)
		evaluation.SetStartTime()
	}

	if evaluation.Status.RetryCount >= evaluation.Spec.Retries {
		r.recordEvent("Warning", evaluation, "ReconcileTimeOut", "retryCount exceeded")
		err := fmt.Errorf("retryCount for evaluation exceeded")
		span.SetStatus(codes.Error, err.Error())
		evaluation.Status.OverallStatus = common.StateFailed
		r.updateFinishedEvaluationMetrics(ctx, evaluation, span)
		return ctrl.Result{}, err
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
			} else {
				r.Log.Error(err, "Failed to retrieve a resource")
			}
			return ctrl.Result{Requeue: true, RequeueAfter: 30 * time.Second}, nil
		}

		if evaluationDefinition == nil || evaluationProvider == nil {
			return ctrl.Result{}, nil
		}

		statusSummary := common.StatusSummary{}
		statusSummary.Total = len(evaluationDefinition.Spec.Objectives)
		newStatus := make(map[string]klcv1alpha1.EvaluationStatusItem)

		if evaluation.Status.EvaluationStatus == nil {
			evaluation.Status.EvaluationStatus = make(map[string]klcv1alpha1.EvaluationStatusItem)
		}

		for _, query := range evaluationDefinition.Spec.Objectives {
			if _, ok := evaluation.Status.EvaluationStatus[query.Name]; !ok {
				evaluation.AddEvaluationStatus(query)
			}
			if evaluation.Status.EvaluationStatus[query.Name].Status.IsSucceeded() {
				statusSummary = common.UpdateStatusSummary(common.StateSucceeded, statusSummary)
				newStatus[query.Name] = evaluation.Status.EvaluationStatus[query.Name]
				continue
			}
			statusItem := r.queryEvaluation(query, *evaluationProvider)
			statusSummary = common.UpdateStatusSummary(statusItem.Status, statusSummary)
			newStatus[query.Name] = *statusItem
		}

		evaluation.Status.RetryCount++
		evaluation.Status.EvaluationStatus = newStatus
		if common.GetOverallState(statusSummary) == common.StateSucceeded {
			evaluation.Status.OverallStatus = common.StateSucceeded
		} else {
			evaluation.Status.OverallStatus = common.StatePending
		}

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

func (r *KeptnEvaluationReconciler) updateFinishedEvaluationMetrics(ctx context.Context, evaluation *klcv1alpha1.KeptnEvaluation, span trace.Span) error {
	r.recordEvent("Normal", evaluation, string(evaluation.Status.OverallStatus), "the evaluation has "+string(evaluation.Status.OverallStatus))

	if !evaluation.IsEndTimeSet() {
		// metrics: decrement active evaluation counter
		r.Meters.AnalysisActive.Add(ctx, -1, evaluation.GetActiveMetricsAttributes()...)
		evaluation.SetEndTime()
	}

	err := r.Client.Status().Update(ctx, evaluation)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		r.recordEvent("Warning", evaluation, "ReconcileErrored", "could not update status")
		return err
	}

	attrs := evaluation.GetMetricsAttributes()

	r.Log.Info("Increasing evaluation count")

	// metrics: increment evaluation counter
	r.Meters.AnalysisCount.Add(ctx, 1, attrs...)

	// metrics: add evaluation duration
	duration := evaluation.Status.EndTime.Time.Sub(evaluation.Status.StartTime.Time)
	r.Meters.AnalysisDuration.Record(ctx, duration.Seconds(), attrs...)
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KeptnEvaluationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&klcv1alpha1.KeptnEvaluation{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Complete(r)
}

func (r *KeptnEvaluationReconciler) fetchDefinitionAndProvider(ctx context.Context, namespacedDefinition types.NamespacedName) (*klcv1alpha1.KeptnEvaluationDefinition, *klcv1alpha1.KeptnEvaluationProvider, error) {
	evaluationDefinition := &klcv1alpha1.KeptnEvaluationDefinition{}

	if err := r.Client.Get(ctx, namespacedDefinition, evaluationDefinition); err != nil {
		return nil, nil, err
	}

	namespacedProvider := types.NamespacedName{
		Namespace: namespacedDefinition.Namespace,
		Name:      evaluationDefinition.Spec.Source,
	}

	evaluationProvider := &klcv1alpha1.KeptnEvaluationProvider{}

	if err := r.Client.Get(ctx, namespacedProvider, evaluationProvider); err != nil {
		return nil, nil, err
	}

	return evaluationDefinition, evaluationProvider, nil
}

func (r *KeptnEvaluationReconciler) queryEvaluation(objective klcv1alpha1.Objective, provider klcv1alpha1.KeptnEvaluationProvider) *klcv1alpha1.EvaluationStatusItem {
	query := &klcv1alpha1.EvaluationStatusItem{
		Value:  "",
		Status: common.StateFailed, //setting status per default to failed
	}

	queryTime := time.Now().UTC()
	r.Log.Info("Running query: /api/v1/query?query=" + objective.Query + "&time=" + queryTime.String())

	client, err := promapi.NewClient(promapi.Config{Address: provider.Spec.TargetServer, Client: &http.Client{}})
	api := prometheus.NewAPI(client)
	result, w, err := api.Query(
		context.Background(),
		objective.Query,
		queryTime,
		[]prometheus.Option{}...,
	)

	if err != nil {
		query.Error = err.Error()
		return query
	}

	if len(w) != 0 {
		query.Error = w[0]
		r.Log.Info("Prometheus API returned warnings: " + w[0])
	}

	// check if we can cast the result to a vector, it might be another data struct which we can't process
	resultVector, ok := result.(model.Vector)
	if !ok {
		query.Error = "could not cast result"
		return query
	}

	// We are only allowed to return one value, if not the query may be malformed
	// we are using two different errors to give the user more information about the result
	if len(resultVector) == 0 {
		r.Log.Info("No values in query result")
		query.Error = "No values in query result"
		return query
	} else if len(resultVector) > 1 {
		r.Log.Info("Too many values in the query result")
		query.Error = "Too many values in the query result"
		return query
	}

	// parse the first entry as float and return the value if it's a valid float value
	query.Value = resultVector[0].Value.String()

	check, err := r.checkValue(objective, query)

	if err != nil {
		query.Error = err.Error()
		r.Log.Error(err, "Could not check query result")
	}
	if check {
		query.Status = common.StateSucceeded
	}
	return query
}

func (r *KeptnEvaluationReconciler) checkValue(objective klcv1alpha1.Objective, query *klcv1alpha1.EvaluationStatusItem) (bool, error) {

	if len(query.Value) == 0 || len(objective.EvaluationTarget) == 0 {
		return false, fmt.Errorf("no values")
	}

	eval := objective.EvaluationTarget[1:]
	sign := objective.EvaluationTarget[:1]

	resultValue, err := strconv.ParseFloat(query.Value, 64)
	if err != nil || math.IsNaN(resultValue) {
		return false, err
	}

	compareValue, err := strconv.ParseFloat(eval, 64)
	if err != nil || math.IsNaN(compareValue) {
		return false, err
	}

	// choose comparator
	switch sign {
	case ">":
		return resultValue > compareValue, nil
	case "<":
		return resultValue < compareValue, nil
	default:
		return false, fmt.Errorf("invalid operator")
	}
}

func (r *KeptnEvaluationReconciler) recordEvent(eventType string, evaluation *klcv1alpha1.KeptnEvaluation, shortReason string, longReason string) {
	r.Recorder.Event(evaluation, eventType, shortReason, fmt.Sprintf("%s / Namespace: %s, Name: %s, WorkloadVersion: %s ", longReason, evaluation.Namespace, evaluation.Name, evaluation.Spec.WorkloadVersion))
}
