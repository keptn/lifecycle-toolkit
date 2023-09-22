package provider

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	apierr "k8s.io/apimachinery/pkg/api/errors"
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
	"sigs.k8s.io/custom-metrics-apiserver/pkg/provider"
)

var keptnMetricGroupVersionResource = schema.GroupVersionResource{Group: "metrics.keptn.sh", Version: "v1alpha2", Resource: "keptnmetrics"}

var providerInstance *keptnMetricsProvider

var providerOnce sync.Once

type keptnMetricsProvider struct {
	client         dynamic.Interface
	scheme         *runtime.Scheme
	logger         logr.Logger
	KeptnNamespace string

	// cache is being populated via the updates received by the provider's dynamic informer
	// this way, we avoid sending a request to the Kubernetes API each time a custom metric value should be retrieved
	cache CustomMetricsCache
}

// NewProvider creates and starts a new keptnMetricsProvider. The provider will run until the given context is cancelled.
// the client passed to this function will be used to set up a dynamic informer that listens for KeptnMetric CRDs and provides metric values that reflect their states.
func NewProvider(ctx context.Context, client dynamic.Interface, namespace string) provider.CustomMetricsProvider {
	providerOnce.Do(func() {
		scheme := runtime.NewScheme()

		providerInstance = &keptnMetricsProvider{
			client: client,
			scheme: scheme,
			cache: CustomMetricsCache{
				metrics: map[metricKey]CustomMetricValue{},
			},
			logger:         ctrl.Log.WithName("provider"),
			KeptnNamespace: namespace,
		}

		if err := providerInstance.watchMetrics(ctx); err != nil {
			klog.Fatalf("failed to start informer: %v", err)
		}
	})

	return providerInstance
}

// ListAllMetrics lists all available metrics
func (p *keptnMetricsProvider) ListAllMetrics() []provider.CustomMetricInfo {
	return p.cache.List()
}

// GetMetricByName retrieves a metric based on its name.
// Used for requests such as e.g. /apis/custom.metrics.k8s.io/v1beta2/namespaces/keptn-lifecycle-toolkit/keptnmetrics.metrics.sh/keptnmetric-sample/keptnmetric-sample
func (p *keptnMetricsProvider) GetMetricByName(ctx context.Context, name types.NamespacedName, info provider.CustomMetricInfo, metricSelector labels.Selector) (*custom_metrics.MetricValue, error) {
	klog.InfoS("GetMetricByName()", "name", name, "metricSelector", metricSelector, "context", ctx)
	val, err := p.cache.Get(name)
	if err != nil {
		if errors.Is(err, ErrMetricNotFound) {
			return nil, provider.NewMetricNotFoundForSelectorError(info.GroupResource, info.Metric, name.Name, metricSelector)
		}
		return nil, &apierr.StatusError{ErrStatus: metav1.Status{
			Status:  metav1.StatusFailure,
			Code:    int32(http.StatusInternalServerError),
			Reason:  metav1.StatusReasonInternalError,
			Message: err.Error(),
		}}
	}
	return &val.Value, nil
}

// GetMetricBySelector retrieves a list of metrics based on the given selectors.
// Used for requests such as e.g. /apis/custom.metrics.k8s.io/v1beta2/namespaces/keptn-lifecycle-toolkit/keptnmetrics.metrics.sh/*/*?labelSelector=<key>%3D<value>
func (p *keptnMetricsProvider) GetMetricBySelector(ctx context.Context, _ string, selector labels.Selector, _ provider.CustomMetricInfo, metricSelector labels.Selector) (*custom_metrics.MetricValueList, error) {
	klog.InfoS("GetMetricBySelector()", "selector", selector, "metricSelector", metricSelector, "context", ctx)

	metricValues := p.cache.GetValuesByLabel(selector)

	res := make([]custom_metrics.MetricValue, len(metricValues))
	i := 0
	for _, metricValue := range metricValues {
		res[i] = metricValue.Value
		i++
	}

	return &custom_metrics.MetricValueList{
		Items: res,
	}, nil
}

func (p *keptnMetricsProvider) watchMetrics(ctx context.Context) error {
	factory := dynamicinformer.NewDynamicSharedInformerFactory(p.client, 0)

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

			p.cache.Delete(types.NamespacedName{
				Namespace: unstructuredKeptnMetric.GetNamespace(),
				Name:      unstructuredKeptnMetric.GetName(),
			})
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
		// set the value to defaultMetricsValue, and add the metric to the list of available metrics
		value = defaultMetricsValue
		klog.InfoS("No value available, using default value", "name", unstructuredKeptnMetric.GetName(), "defaultValue", defaultMetricsValue)
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
