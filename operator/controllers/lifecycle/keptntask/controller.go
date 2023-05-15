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
	"reflect"
	"time"

	"github.com/go-logr/logr"
	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3/common"
	controllercommon "github.com/keptn/lifecycle-toolkit/operator/controllers/common"
	controllererrors "github.com/keptn/lifecycle-toolkit/operator/controllers/errors"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/lifecycle/keptntask/providers"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
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

	defer func(task *klcv1alpha3.KeptnTask) {
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
		err = r.CreateJob(ctx, req, task)
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			r.Log.Error(err, "could not create Job")
		} else {
			task.Status.Status = apicommon.StateProgressing
		}
		return ctrl.Result{Requeue: true, RequeueAfter: 10 * time.Second}, nil
	}

	if !task.Status.Status.IsCompleted() {
		err := r.UpdateJob(ctx, req, task)
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
		For(&klcv1alpha3.KeptnTask{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Owns(&batchv1.Job{}).
		Complete(r)
}

func (r *KeptnTaskReconciler) JobExists(ctx context.Context, task klcv1alpha3.KeptnTask, namespace string) (bool, error) {
	jobList := &batchv1.JobList{}

	jobLabels := client.MatchingLabels{}
	for k, v := range task.CreateKeptnLabels() {
		jobLabels[k] = v
	}

	if len(jobLabels) == 0 {
		return false, fmt.Errorf(controllererrors.ErrNoLabelsFoundTask, task.Name)
	}

	if err := r.Client.List(ctx, jobList, client.InNamespace(namespace), jobLabels); err != nil {
		return false, err
	}

	if len(jobList.Items) > 0 {
		return true, nil
	}

	return false, nil
}

func (r *KeptnTaskReconciler) getTracer() controllercommon.ITracer {
	return r.TracerFactory.GetTracer(traceComponentName)
}

func (r *KeptnTaskReconciler) getProvider(taskDef *klcv1alpha3.KeptnTaskDefinition) providers.IProvider {

	jsProvider := &providers.JSFunctionProvider{
		Log:      r.Log,
		Client:   r.Client,
		Scheme:   r.Scheme,
		Recorder: r.Recorder,
	}
	if taskDef.Spec.ProviderType == "" {
		return jsProvider
	}

	switch taskDef.Spec.ProviderType {
	case klcv1alpha3.FUNCTION_PROVIDER:
		return jsProvider
	case klcv1alpha3.CONTAINER_RUNTIME_PROVIDER:

		return &providers.ContainerRuntimeProvider{
			Log:      r.Log,
			Client:   r.Client,
			Scheme:   r.Scheme,
			Recorder: r.Recorder,
		}
	default:
		r.Log.Error(errors.NewBadRequest("provider type does not exist"), "Failed to get task provider:")
		return nil
	}
}

func notEmpty(taskDef *klcv1alpha3.KeptnTaskDefinition) bool {
	switch taskDef.Spec.ProviderType {
	case klcv1alpha3.FUNCTION_PROVIDER:
		return !reflect.DeepEqual(taskDef.Spec.Function, klcv1alpha3.FunctionSpec{})
	case klcv1alpha3.CONTAINER_RUNTIME_PROVIDER:
		return !reflect.DeepEqual(taskDef.Spec.Container, klcv1alpha3.ContainerRuntimeSpec{})
	default:
		return false
	}
}

func (r *KeptnTaskReconciler) getJob(ctx context.Context, jobName string, namespace string) (*batchv1.Job, error) {
	job := &batchv1.Job{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: jobName, Namespace: namespace}, job)
	if err != nil {
		return job, err
	}
	return job, nil
}

func (r *KeptnTaskReconciler) UpdateJob(ctx context.Context, req ctrl.Request, task *klcv1alpha3.KeptnTask) error {
	job, err := r.getJob(ctx, task.Status.JobName, req.Namespace)
	if err != nil {
		task.Status.JobName = ""
		controllercommon.RecordEvent(r.Recorder, apicommon.PhaseReconcileTask, "Warning", task, "JobReferenceRemoved", "removed Job Reference as Job could not be found", "")
		err = r.Client.Status().Update(ctx, task)
		if err != nil {
			r.Log.Error(err, "could not remove job reference for: "+task.Name)
		}
		return err
	}
	if len(job.Status.Conditions) > 0 {
		if job.Status.Conditions[0].Type == batchv1.JobComplete {
			task.Status.Status = apicommon.StateSucceeded
		} else if job.Status.Conditions[0].Type == batchv1.JobFailed {
			task.Status.Status = apicommon.StateFailed
			task.Status.Message = job.Status.Conditions[0].Message
			task.Status.Reason = job.Status.Conditions[0].Reason
		}
	}
	return nil
}

func (r *KeptnTaskReconciler) CreateJob(ctx context.Context, req ctrl.Request, task *klcv1alpha3.KeptnTask) error {

	jobName := ""
	definition, err := providers.GetTaskDefinition(ctx, r.Client, task.Spec.TaskDefinition, req.Namespace)
	if err != nil {
		controllercommon.RecordEvent(r.Recorder, apicommon.PhaseCreateTask, "Warning", task, "TaskDefinitionNotFound", fmt.Sprintf("could not find KeptnTaskDefinition: %s ", task.Spec.TaskDefinition), "")
		return err
	}
	provider := r.getProvider(definition)
	if notEmpty(definition) {
		jobName, err = provider.CreateJob(ctx, req, task, definition)
		if err != nil {
			return err
		}
	}

	task.Status.JobName = jobName
	task.Status.Status = apicommon.StatePending
	err = r.Client.Status().Update(ctx, task)
	if err != nil {
		r.Log.Error(err, "could not update KeptnTask status reference for: "+task.Name)
	}
	r.Log.Info("updated configmap status reference for: " + definition.Name)
	return nil
}
