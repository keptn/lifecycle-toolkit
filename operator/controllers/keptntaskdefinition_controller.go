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
	"github.com/go-logr/logr"
	klcv1alpha1 "github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
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

//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptntaskdefinitions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptntaskdefinitions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptntaskdefinitions/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=create;get;update;list

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the KeptnTaskDefinition object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *KeptnTaskDefinitionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	r.Log.Info("Reconciling KeptnTaskDefinition")

	definition := &klcv1alpha1.KeptnTaskDefinition{}

	if err := r.Client.Get(ctx, req.NamespacedName, definition); err != nil {
		if errors.IsNotFound(err) {
			// taking down all associated K8s resources is handled by K8s
			r.Log.Info("KeptnTaskDefinition resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		r.Log.Error(err, "Failed to get the KeptnTaskDefinition")
		return ctrl.Result{Requeue: true, RequeueAfter: 30 * time.Second}, nil
	}

	if !reflect.DeepEqual(definition.Spec.Function, klcv1alpha1.FunctionSpec{}) {
		if definition.Spec.Function.InlineReference != (klcv1alpha1.InlineReference{}) {
			err := r.ReconcileInlineConfigMap(ctx, req, definition)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
		if definition.Spec.Function.ConfigMapReference != (klcv1alpha1.ConfigMapReference{}) {
			err := r.ReconcileCmConfigMap(ctx, req, definition)
			if err != nil {
				return ctrl.Result{}, err
			}
		}

	}
	r.Log.Info("Finished Reconciling KeptnTaskDefinition")
	return ctrl.Result{}, nil
}

func (r *KeptnTaskDefinitionReconciler) getFunctionConfigMap(ctx context.Context, functionName string, namespace string) (*corev1.ConfigMap, error) {
	cm := &corev1.ConfigMap{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: functionName, Namespace: namespace}, cm)
	if err != nil {
		return cm, err
	}
	return cm, nil
}

func (r *KeptnTaskDefinitionReconciler) ReconcileCmConfigMap(ctx context.Context, req ctrl.Request, definition *klcv1alpha1.KeptnTaskDefinition) error {
	if definition.Spec.Function.ConfigMapReference.Name != definition.Status.ConfigMap {
		definition.Status.ConfigMap = definition.Spec.Function.ConfigMapReference.Name
		err := r.Client.Status().Update(ctx, definition)
		if err != nil {
			r.Log.Error(err, "could not update configmap status reference for: "+definition.Name)
			return err
		}
		r.Log.Info("updated configmap status reference for: " + definition.Name)
	}
	return nil
}

func (r *KeptnTaskDefinitionReconciler) ReconcileInlineConfigMap(ctx context.Context, req ctrl.Request, definition *klcv1alpha1.KeptnTaskDefinition) error {
	cmIsNew := false
	functionSpec := definition.Spec.Function
	functionName := "keptnfn-" + definition.Name

	cm, err := r.getFunctionConfigMap(ctx, functionName, req.Namespace)
	if err != nil {
		if errors.IsNotFound(err) {
			cmIsNew = true
		} else {
			return fmt.Errorf("could not get function configMap: %w", err)
		}
	}

	functionCm := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      functionName,
			Namespace: definition.Namespace,
		},
		Data: map[string]string{
			"code": functionSpec.InlineReference.Code,
		},
	}
	err = controllerutil.SetControllerReference(definition, &functionCm, r.Scheme)
	if err != nil {
		r.Log.Error(err, "could not set controller reference for ConfigMap: "+functionCm.Name)
	}

	if cmIsNew {
		err := r.Client.Create(ctx, &functionCm)
		if err != nil {
			r.Recorder.Event(definition, "Warning", "ConfigMapNotCreated", fmt.Sprintf("Could not create configmap / Namespace: %s, Name: %s ", functionCm.Namespace, functionCm.Name))
			return err
		}
		r.Recorder.Event(definition, "Normal", "ConfigMapCreated", fmt.Sprintf("Created configmap / Namespace: %s, Name: %s ", functionCm.Namespace, functionCm.Name))

	} else {
		if !reflect.DeepEqual(cm, functionCm) {
			err := r.Client.Update(ctx, &functionCm)
			if err != nil {
				r.Recorder.Event(definition, "Warning", "ConfigMapNotUpdated", fmt.Sprintf("Could not update configmap / Namespace: %s, Name: %s ", functionCm.Namespace, functionCm.Name))
				return err
			}
			r.Recorder.Event(definition, "Normal", "ConfigMapUpdated", fmt.Sprintf("Updated configmap / Namespace: %s, Name: %s ", functionCm.Namespace, functionCm.Name))
		}
	}

	definition.Status.ConfigMap = functionCm.Name
	err = r.Client.Status().Update(ctx, definition)
	if err != nil {
		r.Log.Error(err, "could not update configmap status reference for: "+definition.Name)
		return err
	}
	r.Log.Info("updated configmap status reference for: " + definition.Name)
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KeptnTaskDefinitionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&klcv1alpha1.KeptnTaskDefinition{}).
		Owns(&corev1.ConfigMap{}).
		Complete(r)
}
