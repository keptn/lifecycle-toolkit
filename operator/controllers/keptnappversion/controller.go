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

package keptnappversion

import (
	"context"
	"fmt"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"time"

	"github.com/go-logr/logr"
	"github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1/common"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	klcv1alpha1 "github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1"
)

// KeptnAppVersionReconciler reconciles a KeptnAppVersion object
type KeptnAppVersionReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Log      logr.Logger
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnappversions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnappversions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnappversions/finalizers,verbs=update
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnworkloadinstances/status,verbs=get;update;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the KeptnAppVersion object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *KeptnAppVersionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	r.Log = log.FromContext(ctx)
	r.Log.Info("Searching for Keptn App Version")

	appVersion := &klcv1alpha1.KeptnAppVersion{}
	err := r.Get(ctx, req.NamespacedName, appVersion)
	if errors.IsNotFound(err) {
		return reconcile.Result{}, nil
	}

	if err != nil {
		r.Log.Error(err, "App Version not found")
		return reconcile.Result{}, fmt.Errorf("could not fetch KeptnappVersion: %+v", err)
	}

	phase := common.PhaseAppPreDeployment
	if !appVersion.IsPreDeploymentSucceeded() {
		r.Log.Info(fmt.Sprintf("%s Tasks not finished", phase.LongName))
		if appVersion.IsPreDeploymentFailed() {
			r.recordEvent(phase, "Warning", appVersion, "Failed", "has failed")
			return ctrl.Result{Requeue: true, RequeueAfter: 60 * time.Second}, nil
		}
		r.recordEvent(phase, "Warning", appVersion, "NotFinished", "has not finished")
		state, err := r.reconcilePreDeployment(ctx, appVersion)
		if err != nil {
			r.recordEvent(phase, "Warning", appVersion, "ReconcileErrored", "could not get reconciled")
			return ctrl.Result{Requeue: true}, err
		}
		if state.IsSucceeded() {
			r.recordEvent(phase, "Normal", appVersion, "Succeeded", "has succeeded")
		}
		return ctrl.Result{Requeue: true, RequeueAfter: 5 * time.Second}, nil
	}

	phase = common.PhaseAppDeployment
	if !appVersion.AreWorkloadsSucceeded() {
		r.Log.Info("Workloads post deployments not succeeded")
		if appVersion.AreWorkloadsFailed() {
			r.recordEvent(phase, "Warning", appVersion, "Failed", "has failed")
			return ctrl.Result{Requeue: true, RequeueAfter: 60 * time.Second}, nil
		}
		r.recordEvent(phase, "Warning", appVersion, "NotFinished", "is not finished")
		state, err := r.reconcileWorkloads(ctx, appVersion)
		if err != nil {
			r.recordEvent(phase, "Warning", appVersion, "ReconcileErrored", "could not get reconciled")
			r.Log.Error(err, "Error reconciling workloads post deployments")
			return ctrl.Result{Requeue: true}, err
		}
		if state.IsSucceeded() {
			r.recordEvent(phase, "Normal", appVersion, "Succeeeded", "has succeeded")
		}
		return ctrl.Result{Requeue: true, RequeueAfter: 10 * time.Second}, nil
	}

	phase = common.PhaseAppPostDeployment
	if !appVersion.IsPostDeploymentSucceeded() {
		r.Log.Info("Post-Deployment checks not finished")
		if appVersion.IsPostDeploymentFailed() {
			r.recordEvent(phase, "Warning", appVersion, "Failed", "has failed")
			return ctrl.Result{Requeue: true, RequeueAfter: 60 * time.Second}, nil
		}
		r.recordEvent(phase, "Warning", appVersion, "NotFinished", "has not finished")
		state, err := r.reconcilePostDeployment(ctx, appVersion)
		if err != nil {
			r.recordEvent(phase, "Warning", appVersion, "ReconcileErrored", "could not get reconciled")
			r.Log.Error(err, "Error reconciling post-deployment checks")
			return ctrl.Result{Requeue: true}, err
		}
		if state.IsSucceeded() {
			r.recordEvent(phase, "Normal", appVersion, "Succeeeded", "has succeeded")
		}
		return ctrl.Result{Requeue: true, RequeueAfter: 5 * time.Second}, nil
	}

	r.recordEvent(phase, "Normal", appVersion, "Finished", "is finished")
	err = r.Client.Status().Update(ctx, appVersion)
	if err != nil {
		return ctrl.Result{Requeue: true}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KeptnAppVersionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&klcv1alpha1.KeptnAppVersion{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Complete(r)
}

func (r *KeptnAppVersionReconciler) createKeptnTask(ctx context.Context, namespace string, appVersion *klcv1alpha1.KeptnAppVersion, taskDefinition string, checkType common.CheckType) (string, error) {

	phase := common.KeptnPhaseType{
		ShortName: "KeptnTaskCreate",
		LongName:  "Keptn Task Create",
	}
	// create TraceContext
	// follow up with a Keptn propagator that JSON-encoded the OTel map into our own key
	traceContextCarrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, traceContextCarrier)
	newTask := &klcv1alpha1.KeptnTask{
		ObjectMeta: metav1.ObjectMeta{
			Name:        common.GenerateTaskName(checkType, taskDefinition),
			Namespace:   namespace,
			Annotations: traceContextCarrier,
		},
		Spec: klcv1alpha1.KeptnTaskSpec{
			Version:          appVersion.Spec.Version,
			AppName:          appVersion.Spec.AppName,
			TaskDefinition:   taskDefinition,
			Parameters:       klcv1alpha1.TaskParameters{},
			SecureParameters: klcv1alpha1.SecureParameters{},
			Type:             checkType,
		},
	}
	err := controllerutil.SetControllerReference(appVersion, newTask, r.Scheme)
	if err != nil {
		r.Log.Error(err, "could not set controller reference:")
	}
	err = r.Client.Create(ctx, newTask)
	if err != nil {
		r.Log.Error(err, "could not create KeptnTask")
		r.recordEvent(phase, "Warning", appVersion, "CreateFailed", "could not create KeptnTask")
		return "", err
	}
	r.recordEvent(phase, "Normal", appVersion, "Created", "created")

	return newTask.Name, nil
}

func (r *KeptnAppVersionReconciler) reconcileTasks(ctx context.Context, checkType common.CheckType, appVersion *klcv1alpha1.KeptnAppVersion) ([]klcv1alpha1.TaskStatus, common.StatusSummary, error) {
	phase := common.KeptnPhaseType{
		ShortName: "ReconcileTasks",
		LongName:  "Reconcile Tasks",
	}

	var tasks []string
	var statuses []klcv1alpha1.TaskStatus

	switch checkType {
	case common.PreDeploymentCheckType:
		tasks = appVersion.Spec.PreDeploymentTasks
		statuses = appVersion.Status.PreDeploymentTaskStatus
	case common.PostDeploymentCheckType:
		tasks = appVersion.Spec.PostDeploymentTasks
		statuses = appVersion.Status.PostDeploymentTaskStatus
	}

	var summary common.StatusSummary
	summary.Total = len(tasks)
	// Check current state of the PrePostDeploymentTasks
	var newStatus []klcv1alpha1.TaskStatus
	for _, taskDefinitionName := range tasks {
		taskStatus := GetTaskStatus(taskDefinitionName, statuses)
		task := &klcv1alpha1.KeptnTask{}
		taskExists := false

		// Check if task has already succeeded or failed
		if taskStatus.Status == common.StateSucceeded || taskStatus.Status == common.StateFailed {
			newStatus = append(newStatus, taskStatus)
			continue
		}

		// Check if Task is already created
		if taskStatus.TaskName != "" {
			err := r.Client.Get(ctx, types.NamespacedName{Name: taskStatus.TaskName, Namespace: appVersion.Namespace}, task)
			if err != nil && errors.IsNotFound(err) {
				taskStatus.TaskName = ""
			} else if err != nil {
				return nil, summary, err
			}
			taskExists = true
		}

		// Create new Task if it does not exist
		if !taskExists {
			taskName, err := r.createKeptnTask(ctx, appVersion.Namespace, appVersion, taskDefinitionName, checkType)
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
		summary = common.UpdateStatusSummary(ns.Status, summary)
	}
	if common.GetOverallState(summary) != common.StateSucceeded {
		r.recordEvent(phase, "Warning", appVersion, "NotFinished", "has not finished")
	}
	return newStatus, summary, nil
}

func GetTaskStatus(taskName string, instanceStatus []klcv1alpha1.TaskStatus) klcv1alpha1.TaskStatus {
	for _, status := range instanceStatus {
		if status.TaskDefinitionName == taskName {
			return status
		}
	}
	return klcv1alpha1.TaskStatus{
		TaskDefinitionName: taskName,
		Status:             common.StatePending,
		TaskName:           "",
	}
}

func (r *KeptnAppVersionReconciler) recordEvent(phase common.KeptnPhaseType, eventType string, appVersion *klcv1alpha1.KeptnAppVersion, shortReason string, longReason string) {
	r.Recorder.Event(appVersion, eventType, fmt.Sprintf("%s%s", phase.ShortName, shortReason), fmt.Sprintf("%s %s / Namespace: %s, Name: %s, Version: %s ", phase.LongName, longReason, appVersion.Namespace, appVersion.Name, appVersion.Spec.Version))
}
