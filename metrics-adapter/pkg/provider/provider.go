package provider

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
	ctrl "sigs.k8s.io/controller-runtime"
	"sync"
	"time"

	//metricsv1alpha1 "github.com/keptn/lifecycle-toolkit/operator/api/metrics/v1alpha1"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/metrics/pkg/apis/custom_metrics"
	"os"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/custom-metrics-apiserver/pkg/provider"
	"sigs.k8s.io/custom-metrics-apiserver/pkg/provider/helpers"
)

const keptnNamespace = "keptn-lifecycle-toolkit-system"

var providerInstance *keptnMetricsProvider

var providerOnce sync.Once

type keptnMetricsProvider struct {
	client    dynamic.Interface
	mapper    apimeta.RESTMapper
	k8sClient ctrlclient.Client
	scheme    *runtime.Scheme
	logger    logr.Logger

	metrics CustomMetrics
}

func NewProvider(ctx context.Context, client dynamic.Interface, mapper apimeta.RESTMapper) (provider.CustomMetricsProvider, error) {
	var err error
	providerOnce.Do(func() {
		scheme := runtime.NewScheme()

		cl, err := ctrlclient.New(config.GetConfigOrDie(), ctrlclient.Options{Scheme: scheme})

		if err != nil {
			fmt.Println("failed to create client")
			os.Exit(1)
		}

		providerInstance = &keptnMetricsProvider{
			client:    client,
			mapper:    mapper,
			k8sClient: cl,
			scheme:    scheme,
			logger:    ctrl.Log.WithName("provider"),
		}

		err = providerInstance.watchMetrics(ctx)
	})

	return providerInstance, err
}

func (p *keptnMetricsProvider) ListAllMetrics() []provider.CustomMetricInfo {
	return p.metrics.List()
}

// metricFor constructs a result for a single metric value. TODO remove
func (p *keptnMetricsProvider) metricFor(name types.NamespacedName, info provider.CustomMetricInfo) (*custom_metrics.MetricValue, error) {
	// construct a reference referring to the described object
	fmt.Println("metricFor " + info.String())
	objRef, err := helpers.ReferenceFor(p.mapper, name, info)
	if err != nil {
		return nil, err
	}

	// TODO get metric value from p.metrics and attach DescribedObject: objRef,
	//metric := &metricsv1alpha1.Metric{}
	//err = p.k8sClient.Get(context.Background(), name, metric)
	//if err != nil {
	//return nil, err
	//}

	metricValue, err := resource.ParseQuantity("1.0") // metric.Status.Value
	if err != nil {
		return nil, err
	}

	return &custom_metrics.MetricValue{
		DescribedObject: objRef,
		Metric: custom_metrics.MetricIdentifier{
			Name: info.Metric,
		},
		// you'll want to use the actual timestamp in a real adapter
		//Timestamp: metric.Status.LastUpdated
		Value: metricValue,
	}, nil
}

func (p *keptnMetricsProvider) GetMetricByName(ctx context.Context, name types.NamespacedName, info provider.CustomMetricInfo, metricSelector labels.Selector) (*custom_metrics.MetricValue, error) {
	val, err := p.metrics.Get(generateCustomMetricInfo(name.Name))
	if err != nil {
		return nil, err
	}
	return &val.Value, nil
}

func (p *keptnMetricsProvider) GetMetricBySelector(_ context.Context, _ string, selector labels.Selector, _ provider.CustomMetricInfo, metricSelector labels.Selector) (*custom_metrics.MetricValueList, error) {

	p.logger.Info("GetMetricBySelector()", "selector", selector, "metricSelector", metricSelector)

	metricValues := p.metrics.GetValuesByLabel(selector)

	res := make([]custom_metrics.MetricValue, len(metricValues))
	for i, metricValue := range metricValues {
		res[i] = metricValue.Value
	}

	return &custom_metrics.MetricValueList{
		Items: res,
	}, nil

	/*
		names, err := helpers.ListObjectNames(p.mapper, p.client, namespace, selector, info)
		if err != nil {
			return nil, err
		}

		res := make([]custom_metrics.MetricValue, len(names))
		for i, name := range names {
			metricValue, err := p.metrics.GetByLabel(generateCustomMetricInfo(name), metricSelector)
			if err != nil {
				p.logger.Error(err, "Could not get MetricValue", "metric", name)
				continue
			}
			res[i] = metricValue.Value
		}

		return &custom_metrics.MetricValueList{
			Items: res,
		}, nil
	*/
}

func (p *keptnMetricsProvider) watchMetrics(ctx context.Context) error {
	// Define the resource we want to watch (e.g. pods)
	metricsResource := schema.GroupVersionResource{Group: "metrics.keptn.sh", Version: "v1alpha1", Resource: "keptnmetrics"}

	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(p.client, 0, keptnNamespace, func(options *metav1.ListOptions) {})

	informer := factory.ForResource(metricsResource).Informer()

	handlers := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			p.updateMetric(obj)
		},
		UpdateFunc: func(oldObj, obj interface{}) {
			p.updateMetric(obj)
		},
		DeleteFunc: func(obj interface{}) {
			unstructuredKeptnMetric := obj.(*unstructured.Unstructured)
			p.metrics.Delete(generateCustomMetricInfo(unstructuredKeptnMetric.GetName()))
		},
	}
	if _, err := informer.AddEventHandler(handlers); err != nil {
		return err
	}
	informer.Run(ctx.Done())
	return nil
}

func (p *keptnMetricsProvider) updateMetric(obj interface{}) {
	unstructuredKeptnMetric := obj.(*unstructured.Unstructured)
	value, found, err := unstructured.NestedString(unstructuredKeptnMetric.UnstructuredContent(), "status", "value")
	if err != nil {
		p.logger.Error(err, "Could not parse metric", "name", unstructuredKeptnMetric.GetName())
		return
	}
	if !found {
		// set the value to 0.0 by default, and add the metric to the list of available metrics
		value = "0.0"
		p.logger.Info("No value available, defaulting to 0.0", "name", unstructuredKeptnMetric.GetName())
	}

	metricValue, err := resource.ParseQuantity(value)
	if err != nil {
		p.logger.Error(err, "Could not parse metric", "name", unstructuredKeptnMetric.GetName())
		return
	}

	p.metrics.Update(generateCustomMetricInfo(unstructuredKeptnMetric.GetName()),
		CustomMetricValue{
			Value: custom_metrics.MetricValue{
				Value:     metricValue,
				Timestamp: metav1.Time{Time: time.Now().UTC()},
			},
			Labels: unstructuredKeptnMetric.GetLabels(),
		},
	)
}

func generateCustomMetricInfo(name string) provider.CustomMetricInfo {
	return provider.CustomMetricInfo{
		GroupResource: schema.GroupResource{
			Group:    "metrics.keptn.sh",
			Resource: "metrics",
		},
		Metric:     name,
		Namespaced: true,
	}
}
