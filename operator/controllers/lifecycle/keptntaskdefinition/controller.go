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

package keptntaskdefinition

import (
	"context"
	"reflect"
	"time"

	"github.com/go-logr/logr"
	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// KeptnTaskDefinitionReconciler reconciles a KeptnTaskDefinition object
type KeptnTaskDefinitionReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Log      logr.Logger
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptntaskdefinitions,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptntaskdefinitions/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptntaskdefinitions/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=create;get;update;list;watch

func (r *KeptnTaskDefinitionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	r.Log.Info("Reconciling KeptnTaskDefinition")

	definition := &klcv1alpha3.KeptnTaskDefinition{}

	if err := r.Client.Get(ctx, req.NamespacedName, definition); err != nil {
		if errors.IsNotFound(err) {
			// taking down all associated K8s resources is handled by K8s
			r.Log.Info("KeptnTaskDefinition resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		r.Log.Error(err, "Failed to get the KeptnTaskDefinition")
		return ctrl.Result{Requeue: true, RequeueAfter: 30 * time.Second}, nil
	}
	if definition.Spec.ProviderType == "" || definition.Spec.ProviderType == klcv1alpha3.FUNCTION_PROVIDER {
		if !reflect.DeepEqual(definition.Spec.Function, klcv1alpha3.FunctionSpec{}) {
			err := r.reconcileFunction(ctx, req, definition)
			if err != nil {
				return ctrl.Result{}, nil
			}
		}
	}
	r.Log.Info("Finished Reconciling KeptnTaskDefinition")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KeptnTaskDefinitionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&klcv1alpha3.KeptnTaskDefinition{}).
		Owns(&corev1.ConfigMap{}).
		Complete(r)
}
