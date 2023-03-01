package provider

import (
	"fmt"
	"sync"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/metrics/pkg/apis/custom_metrics"
	"sigs.k8s.io/custom-metrics-apiserver/pkg/provider"
)

type metricKey string

type CustomMetricValue struct {
	Value  custom_metrics.MetricValue
	Labels map[string]string
}

type CustomMetricsCache struct {
	mtx     sync.RWMutex
	metrics map[metricKey]CustomMetricValue
}

// Update adds a new metricValue for the given metricName to the cache. If an item has already been present for the provided
// metricName, the previous value will be replaced.
func (cm *CustomMetricsCache) Update(metricName string, metricValue CustomMetricValue) {
	cm.mtx.Lock()
	defer cm.mtx.Unlock()
	if cm.metrics == nil {
		cm.metrics = map[metricKey]CustomMetricValue{}
	}
	metricNamespace := metricValue.Value.DescribedObject.Namespace

	metricKey := getMetricKey(metricName, metricNamespace)
	cm.metrics[metricKey] = metricValue
}

// Delete will delete the value for the given metricName
func (cm *CustomMetricsCache) Delete(metricName types.NamespacedName) {
	cm.mtx.Lock()
	defer cm.mtx.Unlock()

	delete(cm.metrics, getMetricKey(metricName.Name, metricName.Namespace))
}

// List returns a slice of provider.CustomMetricInfo objects containing all the available metrics
// that are currently present in the cache
func (cm *CustomMetricsCache) List() []provider.CustomMetricInfo {
	cm.mtx.RLock()
	defer cm.mtx.RUnlock()
	res := make([]provider.CustomMetricInfo, len(cm.metrics))

	i := 0
	for metricInfo := range cm.metrics {
		res[i] = generateCustomMetricInfo(cm.metrics[metricInfo].Value.Metric.Name)
		i++
	}
	return res
}

// ListByLabelSelector returns a slice of provider.CustomMetricInfo objects containing all the available metrics
// that are currently present in the cache and match with the provided labels
func (cm *CustomMetricsCache) ListByLabelSelector(selector labels.Selector) []provider.CustomMetricInfo {
	cm.mtx.RLock()
	defer cm.mtx.RUnlock()
	res := []provider.CustomMetricInfo{}
	for _, metricValue := range cm.metrics {
		if selector.Matches(labels.Set(metricValue.Labels)) {
			res = append(res, generateCustomMetricInfo(metricValue.Value.Metric.Name))
		}
	}
	return res
}

// Get returns the metric value for the given metric name
func (cm *CustomMetricsCache) Get(metricName types.NamespacedName) (*CustomMetricValue, error) {
	cm.mtx.RLock()
	defer cm.mtx.RUnlock()
	metric, ok := cm.metrics[getMetricKey(metricName.Name, metricName.Namespace)]
	if !ok {
		return nil, ErrMetricNotFound
	}
	return &metric, nil
}

// GetValuesByLabel returns a slice of CustomMetricValue objects containing the values of all
// available metrics that match with the given label
func (cm *CustomMetricsCache) GetValuesByLabel(selector labels.Selector) []CustomMetricValue {
	cm.mtx.RLock()
	defer cm.mtx.RUnlock()

	res := []CustomMetricValue{}
	for _, value := range cm.metrics {
		if selector.Matches(labels.Set(value.Labels)) {
			res = append(res, value)
		}
	}
	return res
}

func generateCustomMetricInfo(name string) provider.CustomMetricInfo {
	return provider.CustomMetricInfo{
		GroupResource: schema.GroupResource{
			Group:    metricsGroup,
			Resource: metricsResource,
		},
		Metric:     name,
		Namespaced: true,
	}
}

func getMetricKey(metricName, metricNamespace string) metricKey {
	if metricNamespace == "" {
		metricNamespace = "default"
	}
	mk := fmt.Sprintf("%s-%s", metricNamespace, metricName)
	return metricKey(mk)
}
