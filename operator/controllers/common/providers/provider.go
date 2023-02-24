package providers

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/common/providers/keptnmetric"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// KeptnSLIProvider is the interface that describes the operations that an SLI provider must implement
type KeptnSLIProvider interface {
	EvaluateQuery(ctx context.Context, objective klcv1alpha3.Objective, namespace string) (string, []byte, error)
}

// NewProvider is a factory method that chooses the right implementation of KeptnSLIProvider
func NewProvider(provider string, log logr.Logger, k8sClient client.Client) (KeptnSLIProvider, error) {
	switch strings.ToLower(provider) {
	case KeptnMetricProviderName:
		return &keptnmetric.KeptnMetricProvider{
			Log:       log,
			K8sClient: k8sClient,
		}, nil
	default:
		return nil, fmt.Errorf("provider %s not supported", provider)
	}
}
