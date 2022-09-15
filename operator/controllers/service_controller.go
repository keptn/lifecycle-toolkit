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

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/google/uuid"
	"github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

// ServiceReconciler reconciles a Service object
type ServiceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=services/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=services/finalizers,verbs=update
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=events,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=events/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=events/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Service object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *ServiceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("Searching for service")

	service := &v1alpha1.Service{}
	err := r.Get(ctx, req.NamespacedName, service)
	if errors.IsNotFound(err) {
		logger.Error(err, "Could not find Service")
		return reconcile.Result{}, nil
	}

	if err != nil {
		return reconcile.Result{}, fmt.Errorf("could not fetch Service: %+v", err)
	}

	if service.IsCompleted() {
		return reconcile.Result{}, nil
	}

	logger.Info("Reconciling Service", "service", service.Name)

	if service.IsDeploymentCheckNotCreated() {
		logger.Info("Deployment checks do not exist, creating")

		preDeploymentCheksName, err := r.startPreDeploymentChecks(ctx, service)
		if err != nil {
			logger.Error(err, "Could not start pre-deployment checks")
			return reconcile.Result{}, err
		}

		service.Status.PreDeploymentChecksName = preDeploymentCheksName
		service.Status.Phase = v1alpha1.ServiceRunning

		k8sEvent := r.generateK8sEvent(service, "started")
		if err := r.Create(ctx, k8sEvent); err != nil {
			logger.Error(err, "Could not send started pre-deployment checks event")
			return reconcile.Result{}, err
		}

		if err := r.Status().Update(ctx, service); err != nil {
			logger.Error(err, "Could not update Service")
			return reconcile.Result{}, err
		}
		return ctrl.Result{Requeue: true, RequeueAfter: 5 * time.Second}, nil
	}

	preDeploymentCheksEvent, err := r.getPreDeploymentChecksEvent(ctx, service)
	if err != nil {
		logger.Error(err, "Could not retrieve pre-deployment checks Event")
		return reconcile.Result{}, err
	}

	logger.Info("Checking status")

	if preDeploymentCheksEvent.IsCompleted() {
		if preDeploymentCheksEvent.Status.Phase == v1alpha1.EventFailed {
			service.Status.Phase = v1alpha1.ServiceFailed
		} else {
			service.Status.Phase = v1alpha1.ServiceSucceeded
		}

		if err := r.Status().Update(ctx, service); err != nil {
			logger.Error(err, "Could not update Service")
			return reconcile.Result{}, err
		}

		k8sEvent := r.generateK8sEvent(service, "finished")
		if err := r.Create(ctx, k8sEvent); err != nil {
			logger.Error(err, "Could not send finished pre-deployment checks event")
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}

	return ctrl.Result{Requeue: true, RequeueAfter: 5 * time.Second}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Service{}).
		Complete(r)
}

// func (r *ServiceReconciler) createDeployment(ctx context.Context, service *v1alpha1.Service) (*appsv1.Deployment, error) {
// 	deployment := &appsv1.Deployment{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Annotations: map[string]string{
// 				"keptn.sh/application": service.Spec.ApplicationName,
// 				"keptn.sh/service":     service.Name,
// 			},
// 			Name:      service.Name + "-" + r.generateSuffix(),
// 			Namespace: service.Namespace,
// 			Labels:    service.Spec.DeploymentSpec.Selector.MatchLabels,
// 		},
// 		Spec: service.Spec.DeploymentSpec,
// 	}
// 	for i := 0; i < 5; i++ {
// 		if err := r.Create(ctx, deployment); err != nil {
// 			if errors.IsAlreadyExists(err) {
// 				deployment.Name = service.Name + "-" + r.generateSuffix()
// 				continue
// 			}
// 			return nil, err
// 		}
// 		break
// 	}
// 	return deployment, nil
// }

func (r *ServiceReconciler) generateSuffix() string {
	uid := uuid.New().String()
	return uid[:10]
}

func (r *ServiceReconciler) startPreDeploymentChecks(ctx context.Context, service *v1alpha1.Service) (string, error) {
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
			JobSpec:     service.Spec.PreDeplymentChecks.JobSpec,
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

func (r *ServiceReconciler) generateK8sEvent(service *v1alpha1.Service, eventType string) *corev1.Event {
	return &corev1.Event{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName:    service.Name + "-" + eventType + "-",
			Namespace:       service.Namespace,
			ResourceVersion: "v1alpha1",
			Labels: map[string]string{
				"keptn.sh/application": service.Spec.ApplicationName,
				"keptn.sh/service":     service.Name,
			},
		},
		InvolvedObject: corev1.ObjectReference{
			Kind:      service.Kind,
			Namespace: service.Namespace,
			Name:      service.Name,
		},
		Reason:  string(service.Status.Phase),
		Message: "pre-deployment checks are " + eventType,
		Source: corev1.EventSource{
			Component: service.Kind,
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
		ReportingController: "service-controller",
		ReportingInstance:   "service-controller",
	}
}

func (r *ServiceReconciler) getPreDeploymentChecksEvent(ctx context.Context, service *v1alpha1.Service) (*v1alpha1.Event, error) {
	event := &v1alpha1.Event{}
	err := r.Get(ctx, types.NamespacedName{Name: service.Status.PreDeploymentChecksName, Namespace: service.Namespace}, event)
	if errors.IsNotFound(err) {
		return nil, err
	}

	return event, nil
}

// func (r *ServiceReconciler) isDeploymentRunning(d *appsv1.Deployment) bool {
// 	for _, c := range d.Status.Conditions {
// 		if c.Type == appsv1.DeploymentAvailable && c.Status == corev1.ConditionTrue {
// 			return true
// 		}
// 	}
// 	return false
// }

// func (r *ServiceReconciler) hasDeploymentFailed(d *appsv1.Deployment) bool {
// 	for _, c := range d.Status.Conditions {
// 		if c.Type == appsv1.DeploymentReplicaFailure {
// 			return false
// 		}
// 	}
// 	return false
// }
