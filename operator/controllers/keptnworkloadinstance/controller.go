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
	"github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1/semconv"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"math/rand"
	"time"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1"
	klcv1alpha1 "github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1"
	"github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1/common"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type StatusSummary struct {
	failed    int
	succeeded int
	running   int
	pending   int
}

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
	r.Log = log.FromContext(ctx)
	r.Log.Info("Searching for Keptn Workload Instance")

	workloadInstance := &klcv1alpha1.KeptnWorkloadInstance{}
	err := r.Get(ctx, req.NamespacedName, workloadInstance)
	if errors.IsNotFound(err) {
		return reconcile.Result{}, nil
	}

	if err != nil {
		r.Log.Error(err, "Workload Instance not found")
		return reconcile.Result{}, fmt.Errorf("could not fetch KeptnWorkloadInstance: %+v", err)
	}

	traceContextCarrier := propagation.MapCarrier(workloadInstance.Annotations)
	ctx = otel.GetTextMapPropagator().Extract(ctx, traceContextCarrier)

	ctx, span := r.Tracer.Start(ctx, "reconcile_workload_instance", trace.WithSpanKind(trace.SpanKindConsumer))
	defer span.End()

	semconv.AddAttributeFromWorkloadInstance(span, *workloadInstance)

	if !workloadInstance.IsStartTimeSet() {
		// metrics: increment active deployment counter
		r.Meters.DeploymentActive.Add(ctx, 1, workloadInstance.GetActiveMetricsAttributes()...)
		workloadInstance.SetStartTime()
	}

	if !workloadInstance.IsPreDeploymentCompleted() {
		r.Log.Info("Pre deployment checks not finished")
		err := r.reconcilePreDeployment(ctx, req, workloadInstance)
		if err != nil {
			r.Log.Error(err, "Error reconciling pre-deployment checks")
			span.SetStatus(codes.Error, err.Error())
			return ctrl.Result{Requeue: true}, err
		}
		return ctrl.Result{Requeue: true, RequeueAfter: 5 * time.Second}, nil
	}

	if !workloadInstance.IsDeploymentCompleted() {
		r.Log.Info("Deployment not finished")
		err := r.reconcileDeployment(ctx, workloadInstance)
		if err != nil {
			r.Log.Error(err, "Error reconciling deployment")
			span.SetStatus(codes.Error, err.Error())
			return ctrl.Result{Requeue: true}, err
		}
		return ctrl.Result{Requeue: true, RequeueAfter: 10 * time.Second}, nil
	}

	if !workloadInstance.IsPostDeploymentCompleted() {
		r.Log.Info("Post deployment checks not finished")
		err = r.reconcilePostDeployment(ctx, req, workloadInstance)
		if err != nil {
			r.Log.Error(err, "Error reconciling post-deployment checks")
			span.SetStatus(codes.Error, err.Error())
			return ctrl.Result{Requeue: true}, err
		}
		return ctrl.Result{Requeue: true, RequeueAfter: 5 * time.Second}, nil
	}

	// WorkloadInstance is completed at this place

	if !workloadInstance.IsEndTimeSet() {
		// metrics: decrement active deployment counter
		r.Meters.DeploymentActive.Add(ctx, -1, workloadInstance.GetActiveMetricsAttributes()...)
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

	return ctrl.Result{}, nil
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

func (r *KeptnWorkloadInstanceReconciler) getTaskStatus(taskName string, instanceStatus []klcv1alpha1.WorkloadTaskStatus) klcv1alpha1.WorkloadTaskStatus {
	for _, status := range instanceStatus {
		if status.TaskDefinitionName == taskName {
			return status
		}
	}
	return klcv1alpha1.WorkloadTaskStatus{
		TaskDefinitionName: taskName,
		Status:             common.StatePending,
		TaskName:           "",
	}
}

func (r *KeptnWorkloadInstanceReconciler) getKeptnTask(ctx context.Context, taskName string, namespace string) (*klcv1alpha1.KeptnTask, error) {
	task := &klcv1alpha1.KeptnTask{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: taskName, Namespace: namespace}, task)
	if err != nil {
		return task, err
	}
	return task, nil
}

func updateStatusSummary(status common.KeptnState, summary StatusSummary) StatusSummary {
	switch status {
	case common.StateFailed:
		summary.failed++
	case common.StateSucceeded:
		summary.succeeded++
	case common.StateRunning:
		summary.running++
	case common.StatePending, "":
		summary.pending++
	}
	return summary
}

func getOverallState(summary StatusSummary) common.KeptnState {
	if summary.failed > 0 {
		return common.StateFailed
	}
	if summary.pending > 0 {
		return common.StatePending
	}
	if summary.running > 0 {
		return common.StateRunning
	}
	return common.StateSucceeded
}

func generateTaskName(checkType common.CheckType, taskName string) string {
	randomId := rand.Intn(99999-10000) + 10000
	return fmt.Sprintf("%s-%s-%d", checkType, common.TruncateString(taskName, 32), randomId)
}

func (r *KeptnWorkloadInstanceReconciler) createKeptnTask(ctx context.Context, namespace string, workloadInstance *klcv1alpha1.KeptnWorkloadInstance, taskDefinition string, checkType common.CheckType) (string, error) {
	ctx, span := r.Tracer.Start(ctx, fmt.Sprintf("create_%s_task", checkType), trace.WithSpanKind(trace.SpanKindProducer))
	defer span.End()

	semconv.AddAttributeFromWorkloadInstance(span, *workloadInstance)

	// create TraceContext
	// follow up with a Keptn propagator that JSON-encoded the OTel map into our own key
	traceContextCarrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, traceContextCarrier)
	newTask := &klcv1alpha1.KeptnTask{
		ObjectMeta: metav1.ObjectMeta{
			Name:        generateTaskName(checkType, taskDefinition),
			Namespace:   namespace,
			Annotations: traceContextCarrier,
		},
		Spec: klcv1alpha1.KeptnTaskSpec{
			Workload:         workloadInstance.Spec.WorkloadName,
			WorkloadVersion:  workloadInstance.Spec.Version,
			AppName:          workloadInstance.Spec.AppName,
			TaskDefinition:   taskDefinition,
			Parameters:       klcv1alpha1.TaskParameters{},
			SecureParameters: klcv1alpha1.SecureParameters{},
			Type:             checkType,
		},
	}
	err := controllerutil.SetControllerReference(workloadInstance, newTask, r.Scheme)
	if err != nil {
		r.Log.Error(err, "could not set controller reference:")
	}
	err = r.Client.Create(ctx, newTask)
	if err != nil {
		r.Log.Error(err, "could not create KeptnTask")
		r.Recorder.Event(workloadInstance, "Warning", "KeptnTaskNotCreated", fmt.Sprintf("Could not create KeptnTask / Namespace: %s, Name: %s ", newTask.Namespace, newTask.Name))
		return "", err
	}
	r.Recorder.Event(workloadInstance, "Normal", "KeptnTaskCreated", fmt.Sprintf("Created KeptnTask / Namespace: %s, Name: %s ", newTask.Namespace, newTask.Name))

	return newTask.Name, nil
}

func (r *KeptnWorkloadInstanceReconciler) reconcileChecks(ctx context.Context, checkType common.CheckType, workloadInstance *klcv1alpha1.KeptnWorkloadInstance) ([]v1alpha1.WorkloadTaskStatus, StatusSummary, error) {
	tasks := workloadInstance.Spec.PreDeploymentTasks
	statuses := workloadInstance.Status.PreDeploymentTaskStatus
	if checkType == common.PostDeploymentCheckType {
		tasks = workloadInstance.Spec.PostDeploymentTasks
		statuses = workloadInstance.Status.PostDeploymentTaskStatus
	}
	var summary StatusSummary
	// Check current state of the PrePostDeploymentTasks
	var newStatus []klcv1alpha1.WorkloadTaskStatus
	for _, taskDefinitionName := range tasks {
		taskStatus := r.getTaskStatus(taskDefinitionName, statuses)
		task := &klcv1alpha1.KeptnTask{}
		taskExists := false

		// Check if task has already succeeded or failed
		if taskStatus.Status == common.StateSucceeded || taskStatus.Status == common.StateFailed {
			newStatus = append(newStatus, taskStatus)
			continue
		}

		// Check if Task is already created
		if taskStatus.TaskName != "" {
			err := r.Client.Get(ctx, types.NamespacedName{Name: taskStatus.TaskName, Namespace: workloadInstance.Namespace}, task)
			if err != nil && errors.IsNotFound(err) {
				taskStatus.TaskName = ""
			} else if err != nil {
				return nil, summary, err
			}
			taskExists = true
		}

		// Create new Task if it does not exist
		if !taskExists {
			taskName, err := r.createKeptnTask(ctx, workloadInstance.Namespace, workloadInstance, taskDefinitionName, checkType)
			if err != nil {
				return nil, summary, err
			}
			taskStatus.TaskName = taskName
			taskStatus.SetStartTime()
		} else {
			// Update state of Task if it is already created
			taskStatus.Status = task.Status.Status
			if taskStatus.Status.IsCompleted() {
				taskStatus.SetEndTime()
			}
		}
		// Update state of the Check
		newStatus = append(newStatus, taskStatus)
	}

	for _, ns := range newStatus {
		summary = updateStatusSummary(ns.Status, summary)
	}
	return newStatus, summary, nil
}
