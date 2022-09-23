package keptntaskdefinition

import (
	"context"
	"fmt"
	klcv1alpha1 "github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *KeptnTaskDefinitionReconciler) ReconcileFunction(ctx context.Context, req ctrl.Request, definition *klcv1alpha1.KeptnTaskDefinition) error {
	if definition.Spec.Function.InlineReference != (klcv1alpha1.InlineReference{}) {
		err := r.ReconcileFunctionInline(ctx, req, definition)
		if err != nil {
			return err
		}
	}
	if definition.Spec.Function.ConfigMapReference != (klcv1alpha1.ConfigMapReference{}) {
		err := r.ReconcileFunctionConfigMap(ctx, req, definition)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *KeptnTaskDefinitionReconciler) ReconcileFunctionInline(ctx context.Context, req ctrl.Request, definition *klcv1alpha1.KeptnTaskDefinition) error {
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

	definition.Status.Function.ConfigMap = functionCm.Name
	err = r.Client.Status().Update(ctx, definition)
	if err != nil {
		r.Log.Error(err, "could not update configmap status reference for: "+definition.Name)
		return err
	}
	r.Log.Info("updated configmap status reference for: " + definition.Name)
	return nil
}

func (r *KeptnTaskDefinitionReconciler) ReconcileFunctionConfigMap(ctx context.Context, req ctrl.Request, definition *klcv1alpha1.KeptnTaskDefinition) error {
	if definition.Spec.Function.ConfigMapReference.Name != definition.Status.Function.ConfigMap {
		definition.Status.Function.ConfigMap = definition.Spec.Function.ConfigMapReference.Name
		err := r.Client.Status().Update(ctx, definition)
		if err != nil {
			r.Log.Error(err, "could not update configmap status reference for: "+definition.Name)
			return err
		}
		r.Log.Info("updated configmap status reference for: " + definition.Name)
	}
	return nil
}

func (r *KeptnTaskDefinitionReconciler) getFunctionConfigMap(ctx context.Context, functionName string, namespace string) (*corev1.ConfigMap, error) {
	cm := &corev1.ConfigMap{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: functionName, Namespace: namespace}, cm)
	if err != nil {
		return cm, err
	}
	return cm, nil
}
