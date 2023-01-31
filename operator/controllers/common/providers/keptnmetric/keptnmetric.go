package keptnmetric

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	klcv1alpha2 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2"
	metricsv1alpha1 "github.com/keptn/lifecycle-toolkit/operator/apis/metrics/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type KeptnMetricProvider struct {
	Log       logr.Logger
	k8sClient client.Client
}

// EvaluateQuery fetches the SLI values from KeptnMetric resource
func (p *KeptnMetricProvider) EvaluateQuery(ctx context.Context, objective klcv1alpha2.Objective, provider klcv1alpha2.KeptnEvaluationProvider) (string, []byte, error) {
	metric := &metricsv1alpha1.KeptnMetric{}
	if err := p.k8sClient.Get(ctx, types.NamespacedName{Name: objective.Name, Namespace: provider.Namespace}, metric); err != nil {
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
