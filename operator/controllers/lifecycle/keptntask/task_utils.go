package keptntask

import (
	"context"

	klcv1alpha2 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2"
	"k8s.io/apimachinery/pkg/types"
)

func (r *KeptnTaskReconciler) getTaskDefinition(ctx context.Context, definitionName string, namespace string) (*klcv1alpha2.KeptnTaskDefinition, error) {
	definition := &klcv1alpha2.KeptnTaskDefinition{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: definitionName, Namespace: namespace}, definition)
	if err != nil {
		return definition, err
	}
	return definition, nil
}
