package provider

import (
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/metrics/pkg/apis/custom_metrics"
)

func TestCustomMetrics_Delete(t *testing.T) {
	cm := CustomMetricsCache{
		metrics: map[metricKey]CustomMetricValue{
			"my-namespace-my-metric": {
				Value: custom_metrics.MetricValue{
					Metric: custom_metrics.MetricIdentifier{
						Name: "my-metric",
					},
				},
			},
		},
	}

	cm.Delete(types.NamespacedName{
		Namespace: "my-namespace",
		Name:      "my-metric",
	})

	require.Empty(t, cm.metrics)
}

func TestCustomMetrics_DeleteWrongKey(t *testing.T) {
	cm := CustomMetricsCache{
		metrics: map[metricKey]CustomMetricValue{
			"my-metric": {
				Value: custom_metrics.MetricValue{
					Metric: custom_metrics.MetricIdentifier{
						Name: "my-metric",
					},
				},
			},
		},
	}

	cm.Delete(types.NamespacedName{
		Namespace: "",
		Name:      "something-else",
	})

	require.Len(t, cm.metrics, 1)
}

func TestCustomMetrics_DeleteFromEmptyMetrics(t *testing.T) {
	cm := CustomMetricsCache{}

	cm.Delete(types.NamespacedName{
		Namespace: "",
		Name:      "my-metric",
	})

	require.Empty(t, cm.metrics)
}

func TestCustomMetrics_Get(t *testing.T) {
	cm := CustomMetricsCache{
		metrics: map[metricKey]CustomMetricValue{
			"default-my-metric": {
				Value: custom_metrics.MetricValue{
					Metric: custom_metrics.MetricIdentifier{
						Name: "my-metric",
					},
				},
			},
		},
	}

	val, err := cm.Get(types.NamespacedName{
		Namespace: "",
		Name:      "my-metric",
	})

	require.NoError(t, err)
	require.NotNil(t, val)

	val, err = cm.Get(types.NamespacedName{
		Namespace: "",
		Name:      "nothere",
	})
	require.ErrorIs(t, err, ErrMetricNotFound)
	require.Nil(t, val)
}

func TestCustomMetrics_GetValuesByLabel(t *testing.T) {
	cm := CustomMetricsCache{
		metrics: map[metricKey]CustomMetricValue{
			"default-my-metric": {
				Value: custom_metrics.MetricValue{
					Metric: custom_metrics.MetricIdentifier{
						Name: "my-metric",
					},
				},
			},
			"my-labeled-metric": {
				Value: custom_metrics.MetricValue{
					Metric: custom_metrics.MetricIdentifier{
						Name: "my-labeled-metric",
					},
				},
				Labels: map[string]string{
					"app": "frontend",
				},
			},
		},
	}

	val, err := cm.Get(types.NamespacedName{
		Namespace: "",
		Name:      "my-metric",
	})

	require.NoError(t, err)
	require.NotNil(t, val)

	values := cm.GetValuesByLabel(labels.Set(map[string]string{"app": "frontend"}).AsSelector())

	require.Len(t, values, 1)
	require.Equal(t, "my-labeled-metric", values[0].Value.Metric.Name)

	values = cm.GetValuesByLabel(labels.Set(map[string]string{}).AsSelector())

	require.Len(t, values, 2)
}

func TestCustomMetrics_List(t *testing.T) {
	cm := CustomMetricsCache{
		metrics: map[metricKey]CustomMetricValue{
			"my-metric": {
				Value: custom_metrics.MetricValue{
					Metric: custom_metrics.MetricIdentifier{
						Name: "my-metric",
					},
				},
			},
			"my-labeled-metric": {
				Value: custom_metrics.MetricValue{
					Metric: custom_metrics.MetricIdentifier{
						Name: "my-labeled-metric",
					},
				},
				Labels: map[string]string{
					"app": "frontend",
				},
			},
		},
	}

	metricInfos := cm.List()

	require.Len(t, metricInfos, 2)

	for _, mi := range metricInfos {
		require.NotEmpty(t, mi.Metric)
		require.True(t, mi.Namespaced)
		require.Equal(t, schema.GroupResource{
			Group:    metricsGroup,
			Resource: metricsResource,
		}, mi.GroupResource)
	}
}

func TestCustomMetrics_ListByLabelSelector(t *testing.T) {
	cm := CustomMetricsCache{
		metrics: map[metricKey]CustomMetricValue{
			"my-metric": {
				Value: custom_metrics.MetricValue{
					Metric: custom_metrics.MetricIdentifier{
						Name: "my-metric",
					},
				},
			},
			"my-labeled-metric": {
				Value: custom_metrics.MetricValue{
					Metric: custom_metrics.MetricIdentifier{
						Name: "my-labeled-metric",
					},
				},
				Labels: map[string]string{
					"app": "frontend",
				},
			},
		},
	}

	metricInfos := cm.ListByLabelSelector(labels.Set(map[string]string{"app": "frontend"}).AsSelector())

	require.Len(t, metricInfos, 1)
	require.Equal(t, "my-labeled-metric", metricInfos[0].Metric)

	metricInfos = cm.ListByLabelSelector(labels.Set(map[string]string{}).AsSelector())

	require.Len(t, metricInfos, 2)
}

func TestCustomMetrics_Update(t *testing.T) {
	cm := CustomMetricsCache{}

	cm.Update("my-metric", CustomMetricValue{
		Value: custom_metrics.MetricValue{
			Metric: custom_metrics.MetricIdentifier{
				Name: "my-metric",
			},
		},
	})

	get, err := cm.Get(types.NamespacedName{
		Namespace: "",
		Name:      "my-metric",
	})

	require.Nil(t, err)
	require.Equal(t, "my-metric", get.Value.Metric.Name)

	q := resource.NewQuantity(10, resource.DecimalExponent)
	cm.Update("my-metric", CustomMetricValue{
		Value: custom_metrics.MetricValue{
			Metric: custom_metrics.MetricIdentifier{
				Name: "my-metric",
			},
			Value: *q,
		},
	})

	get, err = cm.Get(types.NamespacedName{
		Namespace: "",
		Name:      "my-metric",
	})

	require.Nil(t, err)
	require.Equal(t, "my-metric", get.Value.Metric.Name)
	require.Equal(t, *q, get.Value.Value)
}
