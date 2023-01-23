package provider

import (
	"sync"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/metrics/pkg/apis/custom_metrics"
	"sigs.k8s.io/custom-metrics-apiserver/pkg/provider"
)

var ErrMetricNotFound = errors.New("no metric value found")

type CustomMetricValue struct {
	Value  custom_metrics.MetricValue
	Labels map[string]string
}

type CustomMetrics struct {
	mtx     sync.RWMutex
	metrics map[string]CustomMetricValue
}

func (cm *CustomMetrics) Update(metric string, metricValue CustomMetricValue) {
	cm.mtx.Lock()
	defer cm.mtx.Unlock()
	if cm.metrics == nil {
		cm.metrics = map[string]CustomMetricValue{}
	}

	cm.metrics[metric] = metricValue
}

func (cm *CustomMetrics) Delete(metric string) {
	cm.mtx.Lock()
	defer cm.mtx.Unlock()

	delete(cm.metrics, metric)
}

func (cm *CustomMetrics) List() []provider.CustomMetricInfo {
	cm.mtx.RLock()
	defer cm.mtx.RUnlock()
	res := []provider.CustomMetricInfo{}
	for metricInfo := range cm.metrics {
		res = append(res, generateCustomMetricInfo(metricInfo))
	}
	return res
}

func (cm *CustomMetrics) ListByLabelSelector(selector labels.Selector) []provider.CustomMetricInfo {
	cm.mtx.RLock()
	defer cm.mtx.RUnlock()
	res := []provider.CustomMetricInfo{}
	for metricInfo, metricValue := range cm.metrics {
		if selector.Matches(labels.Set(metricValue.Labels)) {
			res = append(res, generateCustomMetricInfo(metricInfo))
		}
	}
	return res
}

func (cm *CustomMetrics) Get(metricsInfo string) (*CustomMetricValue, error) {
	cm.mtx.RLock()
	defer cm.mtx.RUnlock()
	metric, ok := cm.metrics[metricsInfo]
	if !ok {
		return nil, ErrMetricNotFound
	}
	return &metric, nil
}

func (cm *CustomMetrics) GetValuesByLabel(selector labels.Selector) []CustomMetricValue {
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
			Group:    "metrics.keptn.sh",
			Resource: "keptnmetrics",
		},
		Metric:     name,
		Namespaced: true,
	}
}
