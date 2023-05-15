package providers

import (
	"context"

	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/spanhandler_mock.go . ISpanHandler
type IProvider interface {
	CreateJob(ctx context.Context, req ctrl.Request, task *klcv1alpha3.KeptnTask, definition *klcv1alpha3.KeptnTaskDefinition) (string, error)
}

func GetTaskDefinition(ctx context.Context, c client.Client, definitionName string, namespace string) (*klcv1alpha3.KeptnTaskDefinition, error) {
	definition := &klcv1alpha3.KeptnTaskDefinition{}
	err := c.Get(ctx, types.NamespacedName{Name: definitionName, Namespace: namespace}, definition)
	if err != nil {
		return definition, err
	}
	return definition, nil
}
