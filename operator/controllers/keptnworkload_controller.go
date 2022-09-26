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
	"github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1/common"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	klcv1alpha1 "github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1"
)

// KeptnWorkloadReconciler reconciles a KeptnWorkload object
type KeptnWorkloadReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnworkloads,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnworkloads/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnworkloads/finalizers,verbs=update
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnworkloadinstances,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnworkloadinstances/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnworkloadinstances/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the KeptnWorkload object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *KeptnWorkloadReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("Searching for workload")

	workload := &klcv1alpha1.KeptnWorkload{}
	err := r.Get(ctx, req.NamespacedName, workload)
	if errors.IsNotFound(err) {
		return reconcile.Result{}, nil
	}

	if err != nil {
		return reconcile.Result{}, fmt.Errorf("could not fetch Workload: %+v", err)
	}

	logger.Info("Reconciling Keptn Workload", "workload", workload.Name)

	workloadInstance := &klcv1alpha1.KeptnWorkloadInstance{}
	err = r.Get(ctx, types.NamespacedName{Namespace: workload.Namespace, Name: workload.GetWorkloadInstanceName()}, workloadInstance)
	if errors.IsNotFound(err) {
		logger.Info("Creating KeptnWorkloadInstance from workload", "workload", workload.Name)
		workloadInstance, err := r.createWorkloadInstance(ctx, workload)
		if err != nil {
			logger.Error(err, "Could not create WorkloadInstance")
			return reconcile.Result{}, err
		}

		k8sEvent := r.generateK8sEvent(workload, workloadInstance)
		if err := r.Create(ctx, k8sEvent); err != nil {
			logger.Error(err, "Could not send WorkloadInstance created K8s event")
			return reconcile.Result{}, err
		}

		return ctrl.Result{}, nil
	}

	if err != nil {
		return reconcile.Result{}, fmt.Errorf("could not fetch WorkloadInstance: %+v", err)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KeptnWorkloadReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&klcv1alpha1.KeptnWorkload{}).
		Complete(r)
}

func (r *KeptnWorkloadReconciler) createWorkloadInstance(ctx context.Context, workload *klcv1alpha1.KeptnWorkload) (*klcv1alpha1.KeptnWorkloadInstance, error) {
	workloadInstance := &klcv1alpha1.KeptnWorkloadInstance{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: workload.Annotations,
			Name:        workload.GetWorkloadInstanceName(),
			Namespace:   workload.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: workload.APIVersion,
					Kind:       workload.Kind,
					Name:       workload.Name,
					UID:        workload.UID,
				},
			},
		},
		Spec: klcv1alpha1.KeptnWorkloadInstanceSpec{
			AppName:            workload.Spec.AppName,
			Version:            workload.Spec.Version,
			PreDeploymentCheck: workload.Spec.PreDeploymentTask,
			ResourceReference:  workload.Spec.ResourceReference,
		},
	}
	return workloadInstance, r.Create(ctx, workloadInstance)
}

func (r *KeptnWorkloadReconciler) generateK8sEvent(workload *klcv1alpha1.KeptnWorkload, workloadInstance *klcv1alpha1.KeptnWorkloadInstance) *corev1.Event {
	return &corev1.Event{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName:    workloadInstance.Name + "-created-",
			Namespace:       workloadInstance.Namespace,
			ResourceVersion: "v1alpha1",
			Labels: map[string]string{
				common.AppAnnotation:      workload.Spec.AppName,
				common.WorkloadAnnotation: workloadInstance.Name,
			},
		},
		InvolvedObject: corev1.ObjectReference{
			Kind:      workloadInstance.Kind,
			Namespace: workloadInstance.Namespace,
			Name:      workloadInstance.Name,
		},
		Reason:  "created",
		Message: "workloadInstance " + workloadInstance.Name + " was created",
		Source: corev1.EventSource{
			Component: workloadInstance.Kind,
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
		Action:              "created",
		ReportingController: "workload-controller",
		ReportingInstance:   "workload-controller",
	}
}
