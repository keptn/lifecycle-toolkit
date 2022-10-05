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
	klcv1alpha1 "github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1"
	"github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1/common"
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

	if !task.IsStartTimeSet() {
		// metrics: increment active task counter
		r.Meters.TaskActive.Add(ctx, 1, task.GetActiveMetricsAttributes()...)
		task.SetStartTime()
	}

	err := r.Client.Status().Update(ctx, task)
	if err != nil {
		return ctrl.Result{Requeue: true}, err
	}

	jobExists, err := r.JobExists(ctx, *task, req.Namespace)
	if err != nil {
		r.Log.Error(err, "Could not check if job is running")
		return ctrl.Result{Requeue: true, RequeueAfter: 30 * time.Second}, nil
	}

	if !jobExists {
		err = r.createJob(ctx, req, task)
		if err != nil {
			return ctrl.Result{Requeue: true}, err
		}
		return ctrl.Result{Requeue: true, RequeueAfter: 10 * time.Second}, nil
	}

	if !task.Status.Status.IsCompleted() {
		err := r.updateJob(ctx, req, task)
		if err != nil {
			return ctrl.Result{Requeue: true, RequeueAfter: 10 * time.Second}, err
		}
		return ctrl.Result{Requeue: true, RequeueAfter: 10 * time.Second}, nil
	}

	r.Log.Info("Finished Reconciling KeptnTask")

	// WorkloadInstance is completed at this place

	if !task.IsEndTimeSet() {
		// metrics: decrement active task counter
		r.Meters.TaskActive.Add(ctx, -1, task.GetActiveMetricsAttributes()...)
		task.SetEndTime()
	}

	err = r.Client.Status().Update(ctx, task)
	if err != nil {
		return ctrl.Result{Requeue: true}, err
	}

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
		return false, fmt.Errorf("no labels found for task: %s", task.Name)
	}

	if err := r.Client.List(ctx, jobList, client.InNamespace(namespace), jobLabels); err != nil {
		return false, err
	}

	if len(jobList.Items) > 0 {
		return true, nil
	}

	return false, nil
}
