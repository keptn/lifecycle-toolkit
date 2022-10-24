package keptntask

import (
	"context"
	klcv1alpha1 "github.com/keptn/lifecycle-controller/operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
)

func (r *KeptnTaskReconciler) getTaskDefinition(ctx context.Context, definitionName string, namespace string) (*klcv1alpha1.KeptnTaskDefinition, error) {
	definition := &klcv1alpha1.KeptnTaskDefinition{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: definitionName, Namespace: namespace}, definition)
	if err != nil {
		return definition, err
	}
	return definition, nil
}
