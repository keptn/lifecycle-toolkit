package provider

import (
	"context"
	"sync"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	"k8s.io/metrics/pkg/apis/custom_metrics"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/custom-metrics-apiserver/pkg/provider"
)

const kltNamespace = "keptn-lifecycle-toolkit-system"

var keptnMetricGroupVersionResource = schema.GroupVersionResource{Group: "metrics.keptn.sh", Version: "v1alpha1", Resource: "keptnmetrics"}

var providerInstance *keptnMetricsProvider

var providerOnce sync.Once

type keptnMetricsProvider struct {
	client    dynamic.Interface
	k8sClient ctrlclient.Client
	scheme    *runtime.Scheme
	logger    logr.Logger

	// cache is being populated via the updates received by the provider's dynamic informer
	// this way, we avoid sending a request to the Kubernetes API each time a custom metric value should be retrieved
	cache CustomMetricsCache
}

func NewProvider(ctx context.Context, client dynamic.Interface) provider.CustomMetricsProvider {
	providerOnce.Do(func() {
		scheme := runtime.NewScheme()

		cl, err := ctrlclient.New(config.GetConfigOrDie(), ctrlclient.Options{Scheme: scheme})

		if err != nil {
			klog.Fatalf("failed to create client: %v", err)
		}

		providerInstance = &keptnMetricsProvider{
			client:    client,
			k8sClient: cl,
			scheme:    scheme,
			cache: CustomMetricsCache{
				metrics: map[string]CustomMetricValue{},
			},
			logger: ctrl.Log.WithName("provider"),
		}

		if err := providerInstance.watchMetrics(ctx); err != nil {
			klog.Fatalf("failed to start informer: %v", err)
		}
	})

	return providerInstance
}

func (p *keptnMetricsProvider) ListAllMetrics() []provider.CustomMetricInfo {
	return p.cache.List()
}

// GetMetricByName retrieves a metric based on its name.
// Used for requests such as e.g. /apis/custom.metrics.k8s.io/v1beta2/namespaces/keptn-lifecycle-toolkit/keptnmetrics.metrics.sh/keptnmetric-sample/keptnmetric-sample
func (p *keptnMetricsProvider) GetMetricByName(ctx context.Context, name types.NamespacedName, info provider.CustomMetricInfo, metricSelector labels.Selector) (*custom_metrics.MetricValue, error) {
	klog.InfoS("GetMetricByName()", "name", name, "metricSelector", metricSelector, "context", ctx)
	val, err := p.cache.Get(name.Name)
	if err != nil {
		return nil, err
	}
	return &val.Value, nil
}

// GetMetricBySelector retrieves a list of metrics based on the given selectors.
// Used for requests such as e.g. /apis/custom.metrics.k8s.io/v1beta2/namespaces/keptn-lifecycle-toolkit/keptnmetrics.metrics.sh/*/*?labelSelector=<key>%3D<value>
func (p *keptnMetricsProvider) GetMetricBySelector(ctx context.Context, _ string, selector labels.Selector, _ provider.CustomMetricInfo, metricSelector labels.Selector) (*custom_metrics.MetricValueList, error) {
	klog.InfoS("GetMetricBySelector()", "selector", selector, "metricSelector", metricSelector, "context", ctx)

	metricValues := p.cache.GetValuesByLabel(selector)

	res := []custom_metrics.MetricValue{}
	for _, metricValue := range metricValues {
		res = append(res, metricValue.Value)
	}

	return &custom_metrics.MetricValueList{
		Items: res,
	}, nil
}

func (p *keptnMetricsProvider) watchMetrics(ctx context.Context) error {
	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(p.client, 0, kltNamespace, nil)

	informer := factory.ForResource(keptnMetricGroupVersionResource).Informer()

	handlers := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			klog.InfoS("AddFunc", "obj", obj)
			p.updateMetric(obj)
		},
		UpdateFunc: func(oldObj, obj interface{}) {
			klog.InfoS("UpdateFunc", "obj", obj)
			p.updateMetric(obj)
		},
		DeleteFunc: func(obj interface{}) {
			klog.InfoS("DeleteFunc", "obj", obj)
			unstructuredKeptnMetric := obj.(*unstructured.Unstructured)
			p.cache.Delete(unstructuredKeptnMetric.GetName())
		},
	}
	if _, err := informer.AddEventHandler(handlers); err != nil {
		return err
	}
	go func() {
		informer.Run(ctx.Done())
	}()
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
		klog.InfoS("No value available, defaulting to 0.0", "name", unstructuredKeptnMetric.GetName())
	}

	metricValue, err := resource.ParseQuantity(value)
	if err != nil {
		klog.ErrorS(err, "Could not parse metric", "name", unstructuredKeptnMetric.GetName())
		return
	}
	metricObj := CustomMetricValue{
		Value: custom_metrics.MetricValue{
			Metric: custom_metrics.MetricIdentifier{
				Name:     unstructuredKeptnMetric.GetName(),
				Selector: &metav1.LabelSelector{MatchLabels: unstructuredKeptnMetric.GetLabels()},
			},
			Timestamp: metav1.Time{Time: time.Now().UTC()},
			Value:     metricValue,
			DescribedObject: custom_metrics.ObjectReference{
				APIVersion: keptnMetricGroupVersionResource.Group + "/" + keptnMetricGroupVersionResource.Version,
				Kind:       "KeptnMetric",
				Name:       unstructuredKeptnMetric.GetName(),
				Namespace:  unstructuredKeptnMetric.GetNamespace(),
			},
		},
		Labels: unstructuredKeptnMetric.GetLabels(),
	}

	p.cache.Update(unstructuredKeptnMetric.GetName(), metricObj)
}
