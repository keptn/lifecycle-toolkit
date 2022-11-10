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
	"fmt"
	"time"

	"github.com/go-logr/logr"
	klcv1alpha1 "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1"
	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1/common"
	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1/semconv"
	controllercommon "github.com/keptn/lifecycle-toolkit/operator/controllers/common"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
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

// KeptnTaskReconciler reconciles a KeptnTask object
type KeptnTaskReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
	Log      logr.Logger
	Meters   common.KeptnMeters
	Tracer   trace.Tracer
}

//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptntasks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptntasks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptntasks/finalizers,verbs=update
//+kubebuilder:rbac:groups=batch,resources=jobs,verbs=create;get;update;list;watch
//+kubebuilder:rbac:groups=batch,resources=jobs/status,verbs=get;list

func (r *KeptnTaskReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	r.Log.Info("Reconciling KeptnTask")
	task := &klcv1alpha1.KeptnTask{}

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

	ctx, span := r.Tracer.Start(ctx, "reconcile_task", trace.WithSpanKind(trace.SpanKindConsumer))
	defer span.End()

	semconv.AddAttributeFromTask(span, *task)

	task.SetStartTime()

	if task.Status.Status.IsPending() {
		task.Status.Status = common.StateProgressing
	}

	defer func(task *klcv1alpha1.KeptnTask) {
		err := r.Client.Status().Update(ctx, task)
		if err != nil {
			r.Log.Error(err, "could not update status")
		}

	}(task)

	jobExists, err := r.JobExists(ctx, *task, req.Namespace)
	if err != nil {
		r.Log.Error(err, "Could not check if job is running")
		span.SetStatus(codes.Error, err.Error())
		return ctrl.Result{Requeue: true, RequeueAfter: 30 * time.Second}, nil
	}

	if !jobExists {
		err = r.createJob(ctx, req, task)
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			r.Log.Error(err, "could not create Job")
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
	r.Meters.TaskCount.Add(ctx, 1, attrs...)

	// metrics: add task duration
	duration := task.Status.EndTime.Time.Sub(task.Status.StartTime.Time)
	r.Meters.TaskDuration.Record(ctx, duration.Seconds(), attrs...)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KeptnTaskReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		// predicate disabling the auto reconciliation after updating the object status
		For(&klcv1alpha1.KeptnTask{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Owns(&batchv1.Job{}).
		Complete(r)
}

func (r *KeptnTaskReconciler) JobExists(ctx context.Context, task klcv1alpha1.KeptnTask, namespace string) (bool, error) {
	jobList := &batchv1.JobList{}

	jobLabels := client.MatchingLabels{}
	for k, v := range createKeptnLabels(task) {
		jobLabels[k] = v
	}

	if len(jobLabels) == 0 {
		return false, fmt.Errorf(controllercommon.ErrNoLabelsFoundTask, task.Name)
	}

	if err := r.Client.List(ctx, jobList, client.InNamespace(namespace), jobLabels); err != nil {
		return false, err
	}

	if len(jobList.Items) > 0 {
		return true, nil
	}

	return false, nil
}

func (r *KeptnTaskReconciler) GetActiveTasks(ctx context.Context) ([]common.GaugeValue, error) {
	tasks := &klcv1alpha1.KeptnTaskList{}
	err := r.List(ctx, tasks)
	if err != nil {
		return nil, fmt.Errorf(controllercommon.ErrCannotRetrieveWorkloadInstancesMsg, err)
	}

	res := []common.GaugeValue{}

	for _, task := range tasks.Items {
		gaugeValue := int64(0)
		if !task.IsEndTimeSet() {
			gaugeValue = int64(1)
		}
		res = append(res, common.GaugeValue{
			Value:      gaugeValue,
			Attributes: task.GetActiveMetricsAttributes(),
		})
	}

	return res, nil
}
