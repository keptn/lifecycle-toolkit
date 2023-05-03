package keptnmetric

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha2"
	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/common"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type KeptnMetricProvider struct {
	Log       logr.Logger
	K8sClient client.Client
}

// FetchData fetches the SLI values from KeptnMetric resource
func (p *KeptnMetricProvider) FetchData(ctx context.Context, objective klcv1alpha3.Objective, namespace string) (string, []byte, error) {
	metric, err := p.GetKeptnMetric(ctx, objective, namespace)
	if err != nil {
		return "", nil, err
	}

	if !metric.IsStatusSet() {
		err := fmt.Errorf("empty value for: %s", metric.Name)
		p.Log.Error(err, "KeptnMetric has no value")
		return "", nil, err
	}

	return metric.Status.Value, metric.Status.RawValue, nil
}

func (p *KeptnMetricProvider) GetKeptnMetric(ctx context.Context, objective klcv1alpha3.Objective, namespace string) (*metricsapi.KeptnMetric, error) {
	metric := &metricsapi.KeptnMetric{}

	if objective.KeptnMetricRef.Namespace != "" {
		if err := p.K8sClient.Get(ctx, types.NamespacedName{Name: objective.KeptnMetricRef.Name, Namespace: objective.KeptnMetricRef.Namespace}, metric); err != nil {
			p.Log.Error(err, "Failed to get KeptnMetric from objective namespace")
			return nil, err
		}
	} else {
		if err := p.K8sClient.Get(ctx, types.NamespacedName{Name: objective.KeptnMetricRef.Name, Namespace: namespace}, metric); err != nil {
			p.Log.Error(err, "Failed to get KeptnMetric from KeptnEvaluation resource namespace")
			if err := p.K8sClient.Get(ctx, types.NamespacedName{Name: objective.KeptnMetricRef.Name, Namespace: common.KLTNamespace}, metric); err != nil {
				p.Log.Error(err, "Failed to get KeptnMetric from "+common.KLTNamespace+" namespace")
				return nil, err
			}
		}
	}

	return metric, nil
}
