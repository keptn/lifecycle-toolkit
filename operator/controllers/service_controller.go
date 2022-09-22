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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1"
)

// ServiceReconciler reconciles a Service object
type ServiceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=services/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=services/finalizers,verbs=update
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=serviceruns,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=serviceruns/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=serviceruns/finalizers,verbs=update

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
	_ = log.FromContext(ctx)

	// logger.Info("Searching for service")

	// service := &v1alpha1.Service{}
	// err := r.Get(ctx, req.NamespacedName, service)
	// if errors.IsNotFound(err) {
	// 	return reconcile.Result{}, nil
	// }

	// if err != nil {
	// 	return reconcile.Result{}, fmt.Errorf("could not fetch Service: %+v", err)
	// }

	// logger.Info("Reconciling Service", "service", service.Name)

	// serviceRun := &v1alpha1.ServiceRun{}
	// err = r.Get(ctx, types.NamespacedName{Namespace: service.Namespace, Name: service.GetServiceRunName()}, serviceRun)
	// if errors.IsNotFound(err) {
	// 	logger.Info("Creating serviceRun from service", "service", service.Name)
	// 	serviceRun, err := r.createServiceRun(ctx, service)
	// 	if err != nil {
	// 		logger.Error(err, "Could not create ServiceRun")
	// 		return reconcile.Result{}, err
	// 	}

	// 	k8sEvent := r.generateK8sEvent(service, serviceRun)
	// 	if err := r.Create(ctx, k8sEvent); err != nil {
	// 		logger.Error(err, "Could not send serviceRun created K8s event")
	// 		return reconcile.Result{}, err
	// 	}

	// 	if err := r.Status().Update(ctx, service); err != nil {
	// 		logger.Error(err, "Could not update Service")
	// 		return reconcile.Result{}, err
	// 	}
	// 	return ctrl.Result{}, nil
	// }

	// if err != nil {
	// 	return reconcile.Result{}, fmt.Errorf("could not fetch ServiceRun: %+v", err)
	// }

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Service{}).
		Complete(r)
}

// func (r *ServiceReconciler) createServiceRun(ctx context.Context, service *v1alpha1.Service) (*v1alpha1.ServiceRun, error) {
// 	serviceRun := &v1alpha1.ServiceRun{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Annotations: map[string]string{
// 				"keptn.sh/application": service.Spec.ApplicationName,
// 				"keptn.sh/service":     service.Name,
// 			},
// 			Name:      service.GetServiceRunName(),
// 			Namespace: service.Namespace,
// 			OwnerReferences: []metav1.OwnerReference{
// 				{
// 					APIVersion: service.APIVersion,
// 					Kind:       service.Kind,
// 					Name:       service.Name,
// 					UID:        service.UID,
// 				},
// 			},
// 		},
// 	}
// 	return serviceRun, r.Create(ctx, serviceRun)
// }

// func (r *ServiceReconciler) generateK8sEvent(service *v1alpha1.Service, serviceRun *v1alpha1.ServiceRun) *corev1.Event {
// 	return &corev1.Event{
// 		ObjectMeta: metav1.ObjectMeta{
// 			GenerateName:    serviceRun.Name + "-created-",
// 			Namespace:       serviceRun.Namespace,
// 			ResourceVersion: "v1alpha1",
// 			Labels: map[string]string{
// 				"keptn.sh/application": service.Spec.ApplicationName,
// 				"keptn.sh/service":     serviceRun.Name,
// 			},
// 		},
// 		InvolvedObject: corev1.ObjectReference{
// 			Kind:      serviceRun.Kind,
// 			Namespace: serviceRun.Namespace,
// 			Name:      serviceRun.Name,
// 		},
// 		Reason:  "created",
// 		Message: "serviceRun " + serviceRun.Name + " was created",
// 		Source: corev1.EventSource{
// 			Component: serviceRun.Kind,
// 		},
// 		Type: "Normal",
// 		EventTime: metav1.MicroTime{
// 			Time: time.Now().UTC(),
// 		},
// 		FirstTimestamp: metav1.Time{
// 			Time: time.Now().UTC(),
// 		},
// 		LastTimestamp: metav1.Time{
// 			Time: time.Now().UTC(),
// 		},
// 		Action:              "created",
// 		ReportingController: "serviceRun-controller",
// 		ReportingInstance:   "serviceRun-controller",
// 	}
// }
