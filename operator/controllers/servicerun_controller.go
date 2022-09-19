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

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	types "k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/google/uuid"
	"github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1"
)

// ServiceRunRunReconciler reconciles a ServiceRunRun object
type ServiceRunReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=serviceruns,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=serviceruns/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=serviceruns/finalizers,verbs=update
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=service,verbs=get
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=events,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=events/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=events/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ServiceRunRun object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *ServiceRunReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("Searching for serviceRun")

	serviceRun := &v1alpha1.ServiceRun{}
	err := r.Get(ctx, req.NamespacedName, serviceRun)
	if errors.IsNotFound(err) {
		return reconcile.Result{}, nil
	}

	if err != nil {
		return reconcile.Result{}, fmt.Errorf("could not fetch ServiceRun: %+v", err)
	}

	logger.Info("Searching for service")

	service := &v1alpha1.Service{}
	err = r.Get(ctx, types.NamespacedName{Name: serviceRun.Spec.ServiceName, Namespace: serviceRun.Namespace}, service)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("could not fetch Service: %+v", err)
	}

	if serviceRun.IsCompleted() {
		return reconcile.Result{}, nil
	}

	logger.Info("Reconciling ServiceRun", "serviceRun", serviceRun.Name)

	if serviceRun.IsDeploymentCheckNotCreated() {
		logger.Info("Deployment checks do not exist, creating")

		preDeploymentCheckName, err := r.startPreDeploymentChecks(ctx, service)
		if err != nil {
			logger.Error(err, "Could not start pre-deployment checks")
			return reconcile.Result{}, err
		}

		serviceRun.Status.PreDeploymentCheckName = preDeploymentCheckName
		serviceRun.Status.Phase = v1alpha1.ServiceRunRunning

		k8sEvent := r.generateK8sEvent(service, serviceRun, "started")
		if err := r.Create(ctx, k8sEvent); err != nil {
			logger.Error(err, "Could not send started pre-deployment checks event")
			return reconcile.Result{}, err
		}

		if err := r.Status().Update(ctx, serviceRun); err != nil {
			logger.Error(err, "Could not update ServiceRun")
			return reconcile.Result{}, err
		}
		return ctrl.Result{Requeue: true, RequeueAfter: 5 * time.Second}, nil
	}

	preDeploymentChecksEvent, err := r.getPreDeploymentChecksEvent(ctx, serviceRun)
	if err != nil {
		logger.Error(err, "Could not retrieve pre-deployment checks Event")
		return reconcile.Result{}, err
	}

	logger.Info("Checking status")

	if preDeploymentChecksEvent.IsCompleted() {
		if preDeploymentChecksEvent.Status.Phase == v1alpha1.EventFailed {
			serviceRun.Status.Phase = v1alpha1.ServiceRunFailed
		} else {
			serviceRun.Status.Phase = v1alpha1.ServiceRunSucceeded
		}

		if err := r.Delete(ctx, preDeploymentChecksEvent); err != nil {
			logger.Error(err, "Could not delete Event")
			return reconcile.Result{}, err
		}

		if err := r.Status().Update(ctx, serviceRun); err != nil {
			logger.Error(err, "Could not update ServiceRun")
			return reconcile.Result{}, err
		}

		k8sEvent := r.generateK8sEvent(service, serviceRun, "finished")
		if err := r.Create(ctx, k8sEvent); err != nil {
			logger.Error(err, "Could not send finished pre-deployment checks event")
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}

	return ctrl.Result{Requeue: true, RequeueAfter: 5 * time.Second}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ServiceRunReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.ServiceRun{}).
		Complete(r)
}

func (r *ServiceRunReconciler) generateSuffix() string {
	uid := uuid.New().String()
	return uid[:10]
}

func (r *ServiceRunReconciler) startPreDeploymentChecks(ctx context.Context, service *v1alpha1.Service) (string, error) {
	event := &v1alpha1.Event{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{
				"keptn.sh/application": service.Spec.ApplicationName,
				"keptn.sh/service":     service.Name,
			},
			Name:      service.Name + "-" + r.generateSuffix(),
			Namespace: service.Namespace,
		},
		Spec: v1alpha1.EventSpec{
			Service:     service.Name,
			Application: service.Spec.ApplicationName,
			JobSpec:     service.Spec.PreDeplymentCheck.JobSpec,
		},
	}
	for i := 0; i < 5; i++ {
		if err := r.Create(ctx, event); err != nil {
			if errors.IsAlreadyExists(err) {
				event.Name = service.Name + "-" + r.generateSuffix()
				continue
			}
			return "", err
		}
		break
	}
	return event.Name, nil
}

func (r *ServiceRunReconciler) generateK8sEvent(serviceRun *v1alpha1.Service, serviceRunRun *v1alpha1.ServiceRun, eventType string) *corev1.Event {
	return &corev1.Event{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName:    serviceRun.Name + "-" + eventType + "-",
			Namespace:       serviceRun.Namespace,
			ResourceVersion: "v1alpha1",
			Labels: map[string]string{
				"keptn.sh/application": serviceRun.Spec.ApplicationName,
				"keptn.sh/serviceRun":  serviceRun.Name,
			},
		},
		InvolvedObject: corev1.ObjectReference{
			Kind:      serviceRun.Kind,
			Namespace: serviceRun.Namespace,
			Name:      serviceRun.Name,
		},
		Reason:  string(serviceRunRun.Status.Phase),
		Message: "pre-deployment checks are " + eventType,
		Source: corev1.EventSource{
			Component: serviceRun.Kind,
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
		ReportingController: "serviceRun-controller",
		ReportingInstance:   "serviceRun-controller",
	}
}

func (r *ServiceRunReconciler) getPreDeploymentChecksEvent(ctx context.Context, serviceRun *v1alpha1.ServiceRun) (*v1alpha1.Event, error) {
	event := &v1alpha1.Event{}
	err := r.Get(ctx, types.NamespacedName{Name: serviceRun.Status.PreDeploymentCheckName, Namespace: serviceRun.Namespace}, event)
	if errors.IsNotFound(err) {
		return nil, err
	}

	return event, nil
}
