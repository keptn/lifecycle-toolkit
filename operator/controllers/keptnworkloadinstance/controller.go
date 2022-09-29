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
	"time"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	klcv1alpha1 "github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1"
	"github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1/common"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// KeptnWorkloadInstanceReconciler reconciles a KeptnWorkloadInstance object
type KeptnWorkloadInstanceReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
	Log      logr.Logger
}

//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnworkloadinstances,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnworkloadinstances/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnworkloadinstances/finalizers,verbs=update
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptntasks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptntasks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptntasks/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;watch
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch
//+kubebuilder:rbac:groups=apps,resources=replicasets,verbs=get;list;watch

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

	r.Log.Info("Workload Instance found", "instance", workloadInstance)

	if workloadInstance.IsPostDeploymentCompleted() {
		return reconcile.Result{}, nil
	}

	r.Log.Info("Post deployment checks not finished")

	isDeployed, err := r.IsWorkloadResourceDeployed(ctx, workloadInstance)
	if err != nil {
		return reconcile.Result{}, err
	}
	if isDeployed {
		err = r.reconcilePostDeployment(ctx, req, workloadInstance)
		if err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true, RequeueAfter: 5 * time.Second}, nil
	}

	r.Log.Info("deployment not finished")

	if workloadInstance.IsPreDeploymentCompleted() {
		r.Log.Info("Post deployment checks not finished")
		return ctrl.Result{Requeue: true}, nil
	}

	r.Log.Info("pre-deployment checks not finished")

	err = r.reconcilePreDeployment(ctx, req, workloadInstance)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{Requeue: true, RequeueAfter: 5 * time.Second}, nil

}

// SetupWithManager sets up the controller with the Manager.
func (r *KeptnWorkloadInstanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&klcv1alpha1.KeptnWorkloadInstance{}).
		Complete(r)
}

func (r *KeptnWorkloadInstanceReconciler) generateSuffix() string {
	uid := uuid.New().String()
	return uid[:10]
}

func (r *KeptnWorkloadInstanceReconciler) IsWorkloadResourceDeployed(ctx context.Context, workloadInstance *klcv1alpha1.KeptnWorkloadInstance) (bool, error) {
	if workloadInstance.Status.DeploymentStatus == common.StateSucceeded {
		return true, nil
	}
	if workloadInstance.Spec.ResourceReference.Kind == "Pod" {
		if r.IsPodRunning(ctx, workloadInstance.Spec.ResourceReference, workloadInstance.Namespace) {
			workloadInstance.Status.DeploymentStatus = common.StateSucceeded
			return true, r.Client.Status().Update(ctx, workloadInstance)
		}
		return false, nil
	}
	if r.IsReplicaSetRunning(ctx, workloadInstance.Spec.ResourceReference, workloadInstance.Namespace) {
		workloadInstance.Status.DeploymentStatus = common.StateSucceeded
		return true, r.Client.Status().Update(ctx, workloadInstance)
	}
	return false, nil
}

func (r *KeptnWorkloadInstanceReconciler) IsPodRunning(ctx context.Context, resource klcv1alpha1.ResourceReference, namespace string) bool {
	podList := &corev1.PodList{}
	r.Client.List(ctx, podList, client.InNamespace(namespace))
	for _, p := range podList.Items {
		if p.UID == resource.UID {
			if p.Status.Phase == corev1.PodRunning {
				return true
			}
			return false
		}
	}
	return false
}

func (r *KeptnWorkloadInstanceReconciler) IsReplicaSetRunning(ctx context.Context, resource klcv1alpha1.ResourceReference, namespace string) bool {

	replica := &appsv1.ReplicaSetList{}
	r.Client.List(ctx, replica, client.InNamespace(namespace))

	for _, re := range replica.Items {
		if re.UID == resource.UID {
			if re.Status.ReadyReplicas == *re.Spec.Replicas {
				return true
			}
			return false
		}
	}
	return false

}
