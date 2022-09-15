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

package controllers

import (
	"context"
	"fmt"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	types "k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/google/uuid"
	"github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EventReconciler reconciles a Event object
type EventReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=events,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=events/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=events/finalizers,verbs=update
//+kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=batch,resources=jobs/status,verbs=get;create;delete
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Event object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *EventReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("reconciling event")

	event := &v1alpha1.Event{}
	err := r.Get(ctx, req.NamespacedName, event)
	if errors.IsNotFound(err) {
		logger.Error(nil, "Could not find Event")
		return reconcile.Result{}, nil
	}

	if err != nil {
		return reconcile.Result{}, fmt.Errorf("could not fetch Event: %+v", err)
	}

	logger.Info("Reconciling Event", event.Spec.Application, event.Spec.Service)

	if event.IsCompleted() {
		return reconcile.Result{}, nil
	}

	job := &batchv1.Job{}
	if event.IsJobNotCreated() {
		logger.Info("job does not exists, creating")
		job, err = r.createJob(ctx, event)
		if err != nil {
			logger.Error(err, "Could not create Job")
			return reconcile.Result{}, err
		}
		event.Status.JobName = job.Name
		event.Status.Phase = v1alpha1.EventRunning

		k8sEvent := r.generateK8sEvent(event, "started")
		if err := r.Create(ctx, k8sEvent); err != nil {
			logger.Error(err, "Could not send started Event event")
			return reconcile.Result{}, err
		}

		if err := r.Status().Update(ctx, event); err != nil {
			logger.Error(err, "Could not update Event")
			return reconcile.Result{}, err
		}
		return ctrl.Result{Requeue: true, RequeueAfter: 5 * time.Second}, nil
	}

	err = r.Get(ctx, types.NamespacedName{Namespace: req.Namespace, Name: event.Status.JobName}, job)
	if err != nil {
		logger.Error(err, "Could not get Job")
		return reconcile.Result{}, fmt.Errorf("could not fetch Job: %+v", err)
	}

	logger.Info("checking status")

	if job.Status.Active == 0 {
		if job.Status.Failed == 0 {
			event.Status.Phase = v1alpha1.EventSucceeded
		} else {
			event.Status.Phase = v1alpha1.EventFailed
		}
		if err := r.Delete(ctx, job); err != nil {
			logger.Error(err, "Could not delete Job")
			return reconcile.Result{}, err
		}
		if err := r.Status().Update(ctx, event); err != nil {
			logger.Error(err, "Could not update Event")
			return reconcile.Result{}, err
		}

		k8sEvent := r.generateK8sEvent(event, "finished")
		if err := r.Create(ctx, k8sEvent); err != nil {
			logger.Error(err, "Could not send finished Event event")
			return reconcile.Result{}, err
		}

		return ctrl.Result{}, nil
	}

	return ctrl.Result{Requeue: true, RequeueAfter: 5 * time.Second}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *EventReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Event{}).
		Complete(r)
}

func (r *EventReconciler) createJob(ctx context.Context, event *v1alpha1.Event) (*batchv1.Job, error) {
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{
				"keptn.sh/application": event.Spec.Application,
				"keptn.sh/service":     event.Spec.Service,
			},
			Name:      event.Name + "-" + r.generateSuffix(),
			Namespace: event.Namespace,
		},
		Spec: event.Spec.JobSpec,
	}
	for i := 0; i < 5; i++ {
		if err := r.Create(ctx, job); err != nil {
			if errors.IsAlreadyExists(err) {
				job.Name = event.Name + "-" + r.generateSuffix()
				continue
			}
			return nil, err
		}
		break
	}
	return job, nil
}

func (r *EventReconciler) generateSuffix() string {
	uid := uuid.New().String()
	return uid[:10]
}

func (r *EventReconciler) generateK8sEvent(event *v1alpha1.Event, eventType string) *corev1.Event {
	return &corev1.Event{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName:    event.Name + "-" + eventType + "-",
			Namespace:       event.Namespace,
			ResourceVersion: "v1alpha1",
			Labels: map[string]string{
				"keptn.sh/application": event.Spec.Application,
				"keptn.sh/service":     event.Spec.Service,
				"keptn.sh/event":       event.Name,
			},
		},
		InvolvedObject: corev1.ObjectReference{
			Kind:      event.Kind,
			Namespace: event.Namespace,
			Name:      event.Name,
		},
		Reason:  string(event.Status.Phase),
		Message: "job is " + eventType,
		Source: corev1.EventSource{
			Component: event.Kind,
		},
		Type: "Normal",
		EventTime: metav1.MicroTime{
			Time: time.Now().UTC(),
		},
		FirstTimestamp: metav1.Time{
			Time: time.Now().UTC(),
		},
		LastTimestamp: metav1.Time{
			Time: time.Now().UTC(),
		},
		Action:              eventType,
		ReportingController: "event-controller",
		ReportingInstance:   "event-controller",
	}
}
