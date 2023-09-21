package provider

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic/fake"
	provider2 "sigs.k8s.io/custom-metrics-apiserver/pkg/provider"
)

const KeptnNamespace = "my-namespace"

func TestProvider(t *testing.T) {
	metricObj1 := getSampleKeptnMetric("my-metric", map[string]interface{}{})

	km := &unstructured.Unstructured{}
	km.SetUnstructuredContent(metricObj1)

	scheme := runtime.NewScheme()
	fakeClient := fake.NewSimpleDynamicClient(scheme, km)

	provider := NewProvider(context.TODO(), fakeClient, KeptnNamespace)

	require.NotNil(t, provider)

	// make sure that the informer picks up what is already there
	cmis := []provider2.CustomMetricInfo{}
	require.Eventually(t, func() bool {
		cmis = provider.ListAllMetrics()

		return len(cmis) != 0
	}, 10*time.Second, 100*time.Millisecond)

	// update the metric
	metricObj1["status"] = map[string]interface{}{
		"value": "10.0",
	}
	km.SetUnstructuredContent(metricObj1)
	_, err := fakeClient.Resource(keptnMetricGroupVersionResource).Namespace(KeptnNamespace).UpdateStatus(context.TODO(), km, metav1.UpdateOptions{})

	require.Nil(t, err)

	// eventually the updated value should be reflected
	require.Eventually(t, func() bool {
		metricValue, err := provider.GetMetricByName(context.TODO(), types.NamespacedName{
			Namespace: KeptnNamespace,
			Name:      "my-metric",
		}, provider2.CustomMetricInfo{}, nil)

		if err != nil || metricValue.Value.Value() != 10 {
			return false
		}
		return true
	}, 10*time.Second, 100*time.Millisecond)

	// look for an unknown metric
	metricValue, err := provider.GetMetricByName(context.TODO(), types.NamespacedName{
		Namespace: KeptnNamespace,
		Name:      "my-unknown-metric",
	}, provider2.CustomMetricInfo{}, labels.Set{}.AsSelector())

	require.NotNil(t, err)
	require.Nil(t, metricValue)

	// look for metrics based on a label selector
	metrics, err := provider.GetMetricBySelector(
		context.TODO(),
		KeptnNamespace,
		labels.Set(map[string]string{"app": "frontend"}).AsSelector(),
		provider2.CustomMetricInfo{},
		nil,
	)

	require.Nil(t, err)
	require.Empty(t, metrics.Items)

	// now, create a matching metric
	metricObj2 := getSampleKeptnMetric("my-metric-2", map[string]interface{}{"app": "frontend"})
	km2 := &unstructured.Unstructured{}
	km2.SetUnstructuredContent(metricObj2)

	_, err = fakeClient.Resource(keptnMetricGroupVersionResource).Namespace(KeptnNamespace).Create(context.TODO(), km2, metav1.CreateOptions{})
	require.Nil(t, err)

	// wait for the new metric to be registered
	require.Eventually(t, func() bool {
		cmis = provider.ListAllMetrics()

		return len(cmis) == 2
	}, 10*time.Second, 100*time.Millisecond)

	// retrieve based on the selector again
	metrics, err = provider.GetMetricBySelector(
		context.TODO(),
		KeptnNamespace,
		labels.Set(map[string]string{"app": "frontend"}).AsSelector(),
		provider2.CustomMetricInfo{},
		nil,
	)

	require.Nil(t, err)
	require.Len(t, metrics.Items, 1)

	// delete the metrics again
	err = fakeClient.Resource(keptnMetricGroupVersionResource).Namespace(KeptnNamespace).Delete(context.TODO(), "my-metric", metav1.DeleteOptions{})
	require.Nil(t, err)
	err = fakeClient.Resource(keptnMetricGroupVersionResource).Namespace(KeptnNamespace).Delete(context.TODO(), "my-metric-2", metav1.DeleteOptions{})
	require.Nil(t, err)

	// wait for the length of the returned list to be 0
	require.Eventually(t, func() bool {
		cmis = provider.ListAllMetrics()

		return len(cmis) == 0
	}, 10*time.Second, 100*time.Millisecond)
}

func getSampleKeptnMetric(metricName string, labels map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"apiVersion": "metrics.keptn.sh/v1alpha2",
		"kind":       "KeptnMetric",
		"metadata": map[string]interface{}{
			"name":      metricName,
			"namespace": KeptnNamespace,
			"labels":    labels,
		},
		"spec": map[string]interface{}{
			"fetchIntervalSeconds": "2",
			"provider": map[string]interface{}{
				"name": "my-provider",
			},
			"query": "my-query",
		},
	}
}
