package provider

import (
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/metrics/pkg/apis/custom_metrics"
	"sigs.k8s.io/custom-metrics-apiserver/pkg/provider"
	"sync"
)

var ErrMetricNotFound = errors.New("no metric value found")

type CustomMetricValue struct {
	Value  custom_metrics.MetricValue
	Labels map[string]string
}

type CustomMetrics struct {
	mtx     sync.RWMutex
	metrics map[provider.CustomMetricInfo]CustomMetricValue
}

func (cm *CustomMetrics) Update(metric provider.CustomMetricInfo, metricValue CustomMetricValue) {
	cm.mtx.Lock()
	defer cm.mtx.Unlock()

	cm.metrics[metric] = metricValue
}

func (cm *CustomMetrics) Delete(metric provider.CustomMetricInfo) {
	cm.mtx.Lock()
	defer cm.mtx.Unlock()

	delete(cm.metrics, metric)
}

func (cm *CustomMetrics) List() []provider.CustomMetricInfo {
	cm.mtx.RLock()
	defer cm.mtx.RUnlock()
	res := make([]provider.CustomMetricInfo, len(cm.metrics))
	for metricInfo := range cm.metrics {
		res = append(res, metricInfo)
	}
	return res
}

func (cm *CustomMetrics) ListByLabelSelector(selector labels.Selector) []provider.CustomMetricInfo {
	cm.mtx.RLock()
	defer cm.mtx.RUnlock()
	res := []provider.CustomMetricInfo{}
	for metricInfo, metricValue := range cm.metrics {
		if selector.Matches(labels.Set(metricValue.Labels)) {
			res = append(res, metricInfo)
		}
	}
	return res
}

func (cm *CustomMetrics) Get(metricsInfo provider.CustomMetricInfo) (*CustomMetricValue, error) {
	cm.mtx.RLock()
	defer cm.mtx.RUnlock()
	metric, ok := cm.metrics[metricsInfo]
	if !ok {
		return nil, ErrMetricNotFound
	}
	return &metric, nil
}

func (cm *CustomMetrics) GetByLabel(metricsInfo provider.CustomMetricInfo, selector labels.Selector) (*CustomMetricValue, error) {
	cm.mtx.RLock()
	defer cm.mtx.RUnlock()
	metric, ok := cm.metrics[metricsInfo]
	if !ok {
		return nil, ErrMetricNotFound
	} else if !selector.Matches(labels.Set(metric.Labels)) {
		return nil, ErrMetricNotFound
	}
	return &metric, nil
}
