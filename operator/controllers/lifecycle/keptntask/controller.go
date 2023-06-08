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

package keptntask

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3/common"
	controllercommon "github.com/keptn/lifecycle-toolkit/operator/controllers/common"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

const traceComponentName = "keptn/operator/task"

// KeptnTaskReconciler reconciles a KeptnTask object
type KeptnTaskReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	Recorder      record.EventRecorder
	Log           logr.Logger
	Meters        apicommon.KeptnMeters
	TracerFactory controllercommon.TracerFactory
}

// +kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptntasks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptntasks/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptntasks/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch,resources=jobs,verbs=create;get;update;list;watch
// +kubebuilder:rbac:groups=batch,resources=jobs/status,verbs=get;list

func (r *KeptnTaskReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	r.Log.Info("Reconciling KeptnTask")
	task := &klcv1alpha3.KeptnTask{}

	if err := r.Client.Get(ctx, req.NamespacedName, task); err != nil {
		if errors.IsNotFound(err) {
			// taking down all associated K8s resources is handled by K8s
			r.Log.Info("KeptnTask resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		r.Log.Error(err, "Failed to get the KeptnTask")
		return ctrl.Result{Requeue: true, RequeueAfter: 30 * time.Second}, nil
	}

	traceContextCarrier := propagation.MapCarrier(task.Annotations)
	ctx = otel.GetTextMapPropagator().Extract(ctx, traceContextCarrier)

	ctx, span := r.getTracer().Start(ctx, "reconcile_task", trace.WithSpanKind(trace.SpanKindConsumer))
	defer span.End()

	task.SetSpanAttributes(span)

	task.SetStartTime()

	defer func() {
		err := r.Client.Status().Update(ctx, task)
		if err != nil {
			r.Log.Error(err, "could not update status")
		}
	}()

	job, err := r.getJob(ctx, task.Status.JobName, req.Namespace)
	if err != nil && !errors.IsNotFound(err) {
		r.Log.Error(err, "Could not check if job is running")
		span.SetStatus(codes.Error, err.Error())
		return ctrl.Result{Requeue: true, RequeueAfter: 30 * time.Second}, nil
	}

	if job == nil {
		err = r.createJob(ctx, req, task)
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			r.Log.Error(err, "could not create Job")
		} else {
			task.Status.Status = apicommon.StateProgressing
		}
		return ctrl.Result{Requeue: true, RequeueAfter: 10 * time.Second}, nil
	}

	if !task.Status.Status.IsCompleted() {
		err := r.updateJob(ctx, req, task)
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			return ctrl.Result{Requeue: true, RequeueAfter: 10 * time.Second}, err
		}
		return ctrl.Result{Requeue: true, RequeueAfter: 10 * time.Second}, nil
	}

	r.Log.Info("Finished Reconciling KeptnTask")

	// Task is completed at this place
	task.SetEndTime()

	attrs := task.GetMetricsAttributes()

	r.Log.Info("Increasing task count")

	// metrics: increment task counter
	r.Meters.TaskCount.Add(ctx, 1, metric.WithAttributes(attrs...))

	// metrics: add task duration
	duration := task.Status.EndTime.Time.Sub(task.Status.StartTime.Time)
	r.Meters.TaskDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(attrs...))

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KeptnTaskReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		// predicate disabling the auto reconciliation after updating the object status
		For(&klcv1alpha3.KeptnTask{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Owns(&batchv1.Job{}).
		Complete(r)
}

func (r *KeptnTaskReconciler) getTracer() controllercommon.ITracer {
	return r.TracerFactory.GetTracer(traceComponentName)
}
