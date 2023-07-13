package keptnmetric

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/common"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
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

	value, ok, err := unstructured.NestedString(metric.UnstructuredContent(), "status", "value")
	if !ok || err != nil || value == "" {
		err := fmt.Errorf("empty value for: %s", objective.KeptnMetricRef.Name)
		p.Log.Error(err, "KeptnMetric has no value")
		return "", nil, err
	}

	rawValue, ok, err := unstructured.NestedString(metric.UnstructuredContent(), "status", "rawValue")
	if !ok || err != nil || rawValue == "" {
		err := fmt.Errorf("empty raWvalue for: %s", objective.KeptnMetricRef.Name)
		p.Log.Error(err, "KeptnMetric has no rawValue")
		return "", nil, err
	}

	return value, []byte(rawValue), nil
}

func (p *KeptnMetricProvider) GetKeptnMetric(ctx context.Context, objective klcv1alpha3.Objective, namespace string) (*unstructured.Unstructured, error) {
	metric := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       "KeptnMetric",
			"apiVersion": "metrics.keptn.sh/v1alpha3",
		},
	}

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
