package keptntaskdefinition

import (
	"context"
	"reflect"

	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3/common"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *KeptnTaskDefinitionReconciler) generateConfigMap(spec *klcv1alpha3.RuntimeSpec, name string, namespace string) *corev1.ConfigMap {

	functionCm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: map[string]string{
			"code": spec.Inline.Code,
		},
	}
	return functionCm
}

func (r *KeptnTaskDefinitionReconciler) reconcileConfigMap(ctx context.Context, functionCm *corev1.ConfigMap, cm *corev1.ConfigMap) {

	if (cm == nil || reflect.DeepEqual(cm, &corev1.ConfigMap{})) && functionCm != nil { //cm does not exist or new taskdef with inline func
		err := r.Client.Create(ctx, functionCm)
		if err != nil {
			r.Log.Error(err, "could not create ConfigMap")
			r.EventSender.SendK8sEvent(apicommon.PhaseReconcileTask, "Warning", functionCm, apicommon.PhaseStateFailed, "could not create configmap", "")
			return
		}
	} else if !reflect.DeepEqual(cm, functionCm) && functionCm != nil { //cm and inline func exists but differ
		err := r.Client.Update(ctx, functionCm)
		if err != nil {
			r.Log.Error(err, "could not update ConfigMap")
			r.EventSender.SendK8sEvent(apicommon.PhaseReconcileTask, "Warning", functionCm, apicommon.PhaseStateFailed, "uould not update configmap", "")
			return
		}
	}
	//nothing changed
}

func (r *KeptnTaskDefinitionReconciler) getConfigMap(ctx context.Context, cmName string, namespace string) (*corev1.ConfigMap, error) {
	cm := &corev1.ConfigMap{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: cmName, Namespace: namespace}, cm)
	if err != nil {
		r.Log.Info("could not retrieve ConfigMap '%s': %s", cmName, err.Error())
		return nil, err
	}
	return cm, nil
}

func (r *KeptnTaskDefinitionReconciler) updateTaskDefinitionStatus(functionCm *corev1.ConfigMap, definition *klcv1alpha3.KeptnTaskDefinition) {
	// config map referenced but does not exist we can use the status to signify that
	if functionCm != nil && definition.Status.Function.ConfigMap != functionCm.Name { //configmap referenced exists but old
		definition.Status.Function.ConfigMap = functionCm.Name
		//and  make sure that the definition controls the config map
		err := controllerutil.SetControllerReference(definition, functionCm, r.Scheme)
		if err != nil {
			r.Log.Error(err, "could not set controller reference for ConfigMap: "+functionCm.Name)
		}
	}
}
