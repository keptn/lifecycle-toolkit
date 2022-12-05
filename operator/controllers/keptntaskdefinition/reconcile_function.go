package keptntaskdefinition

import (
	"context"
	"fmt"
	"reflect"

	klcv1alpha2 "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha2"
	controllererrors "github.com/keptn/lifecycle-toolkit/operator/controllers/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *KeptnTaskDefinitionReconciler) reconcileFunction(ctx context.Context, req ctrl.Request, definition *klcv1alpha2.KeptnTaskDefinition) error {
	if definition.Spec.Function.Inline != (klcv1alpha2.Inline{}) {
		err := r.reconcileFunctionInline(ctx, req, definition)
		if err != nil {
			return err
		}
	}
	if definition.Spec.Function.ConfigMapReference != (klcv1alpha2.ConfigMapReference{}) {
		err := r.reconcileFunctionConfigMap(ctx, req, definition)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *KeptnTaskDefinitionReconciler) reconcileFunctionInline(ctx context.Context, req ctrl.Request, definition *klcv1alpha2.KeptnTaskDefinition) error {
	cmIsNew := false
	functionSpec := definition.Spec.Function
	functionName := "keptnfn-" + definition.Name

	cm, err := r.getFunctionConfigMap(ctx, functionName, req.Namespace)
	if err != nil {
		if errors.IsNotFound(err) {
			cmIsNew = true
		} else {
			return fmt.Errorf(controllererrors.ErrCannotGetFunctionConfigMap, err)
		}
	}

	functionCm := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      functionName,
			Namespace: definition.Namespace,
		},
		Data: map[string]string{
			"code": functionSpec.Inline.Code,
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

func (r *KeptnTaskDefinitionReconciler) reconcileFunctionConfigMap(ctx context.Context, req ctrl.Request, definition *klcv1alpha2.KeptnTaskDefinition) error {
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
