package provider

import (
	"context"
	"fmt"
	metricsv1alpha1 "github.com/keptn/lifecycle-toolkit/operator/api/metrics/v1alpha1"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/metrics/pkg/apis/custom_metrics"
	"os"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/custom-metrics-apiserver/pkg/provider"
	"sigs.k8s.io/custom-metrics-apiserver/pkg/provider/helpers"
)

type CustomMetricsProvider interface {
	ListAllMetrics() []provider.CustomMetricInfo

	GetMetricByName(ctx context.Context, name types.NamespacedName, info provider.CustomMetricInfo, metricSelector labels.Selector) (*custom_metrics.MetricValue, error)
	GetMetricBySelector(ctx context.Context, namespace string, selector labels.Selector, info provider.CustomMetricInfo, metricSelector labels.Selector) (*custom_metrics.MetricValueList, error)
}

type keptnMetricsProvider struct {
	client    dynamic.Interface
	mapper    apimeta.RESTMapper
	k8sClient ctrlclient.Client
	scheme    *runtime.Scheme

	// just increment values when they're requested
	values map[provider.CustomMetricInfo]int64
}

func NewProvider(client dynamic.Interface, mapper apimeta.RESTMapper) provider.CustomMetricsProvider {
	scheme := runtime.NewScheme()
	if err := metricsv1alpha1.AddToScheme(scheme); err != nil {
		fmt.Println("failed to add metrics to scheme: " + err.Error())
	}

	cl, err := ctrlclient.New(config.GetConfigOrDie(), ctrlclient.Options{Scheme: scheme})

	if err != nil {
		fmt.Println("failed to create client")
		os.Exit(1)
	}

	return &keptnMetricsProvider{
		client:    client,
		mapper:    mapper,
		values:    make(map[provider.CustomMetricInfo]int64),
		k8sClient: cl,
		scheme:    scheme,
	}
}

func (p *keptnMetricsProvider) ListAllMetrics() []provider.CustomMetricInfo {
	fmt.Println("listAllMetrics")
	list := metricsv1alpha1.MetricList{}
	var customMetrics []provider.CustomMetricInfo
	err := p.k8sClient.List(context.Background(), &list)
	if err != nil {
		fmt.Println("failed to list metrics" + err.Error())
	}

	for _, metric := range list.Items {
		customMetrics = append(customMetrics, provider.CustomMetricInfo{
			GroupResource: schema.GroupResource{
				Group:    "metrics.keptn.sh",
				Resource: "metrics",
			},
			Metric:     metric.Name,
			Namespaced: true,
		})
	}
	return customMetrics
}

// valueFor fetches a value from the fake list and increments it.
func (p *keptnMetricsProvider) valueFor(info provider.CustomMetricInfo) (int64, error) {
	// normalize the value so that you treat plural resources and singular
	// resources the same (e.g. pods vs pod)
	fmt.Println("valueFor " + info.String())
	info, _, err := info.Normalized(p.mapper)
	if err != nil {
		return 0, err
	}

	value := p.values[info]
	value += 1
	p.values[info] = value

	return value, nil
}

// metricFor constructs a result for a single metric value.
func (p *keptnMetricsProvider) metricFor(value int64, name types.NamespacedName, info provider.CustomMetricInfo) (*custom_metrics.MetricValue, error) {
	// construct a reference referring to the described object
	fmt.Println("metricFor " + info.String())
	objRef, err := helpers.ReferenceFor(p.mapper, name, info)
	if err != nil {
		return nil, err
	}

	metric := &metricsv1alpha1.Metric{}
	err = p.k8sClient.Get(context.Background(), name, metric)
	if err != nil {
		return nil, err
	}

	metricValue, err := resource.ParseQuantity(metric.Status.Value)
	if err != nil {
		return nil, err
	}

	return &custom_metrics.MetricValue{
		DescribedObject: objRef,
		Metric: custom_metrics.MetricIdentifier{
			Name: info.Metric,
		},
		// you'll want to use the actual timestamp in a real adapter
		Timestamp: metric.Status.LastUpdated,
		Value:     metricValue,
	}, nil
}

func (p *keptnMetricsProvider) GetMetricByName(ctx context.Context, name types.NamespacedName, info provider.CustomMetricInfo, metricSelector labels.Selector) (*custom_metrics.MetricValue, error) {
	value, err := p.valueFor(info)
	if err != nil {
		return nil, err
	}
	return p.metricFor(value, name, info)
}

func (p *keptnMetricsProvider) GetMetricBySelector(ctx context.Context, namespace string, selector labels.Selector, info provider.CustomMetricInfo, metricSelector labels.Selector) (*custom_metrics.MetricValueList, error) {
	totalValue, err := p.valueFor(info)
	if err != nil {
		return nil, err
	}

	names, err := helpers.ListObjectNames(p.mapper, p.client, namespace, selector, info)
	if err != nil {
		return nil, err
	}

	res := make([]custom_metrics.MetricValue, len(names))
	for i, name := range names {
		// in a real adapter, you might want to consider pre-computing the
		// object reference created in metricFor, instead of recomputing it
		// for each object.
		value, err := p.metricFor(100*totalValue/int64(len(res)), types.NamespacedName{Namespace: namespace, Name: name}, info)
		if err != nil {
			return nil, err
		}
		res[i] = *value
	}

	return &custom_metrics.MetricValueList{
		Items: res,
	}, nil
}
