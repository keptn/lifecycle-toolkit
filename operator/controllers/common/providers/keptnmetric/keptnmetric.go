package keptnmetric

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	metricsv1alpha1 "github.com/keptn/lifecycle-toolkit/metrics-operator/apis/metrics/v1alpha2"
	klcv1alpha2 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type KeptnMetricProvider struct {
	Log       logr.Logger
	K8sClient client.Client
}

// EvaluateQuery fetches the SLI values from KeptnMetric resource
func (p *KeptnMetricProvider) EvaluateQuery(ctx context.Context, objective klcv1alpha2.Objective, namespace string) (string, []byte, error) {
	metric := &metricsv1alpha1.KeptnMetric{}
	if err := p.K8sClient.Get(ctx, types.NamespacedName{Name: objective.Name, Namespace: namespace}, metric); err != nil {
		p.Log.Error(err, "Could not retrieve KeptnMetric")
		return "", nil, err
	}

	if !metric.IsStatusSet() {
		err := fmt.Errorf("empty value for: %s", metric.Name)
		p.Log.Error(err, "KeptnMetric has no value")
		return "", nil, err
	}

	return metric.Status.Value, metric.Status.RawValue, nil
}
