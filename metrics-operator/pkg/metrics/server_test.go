package server

import (
	"bytes"
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha2"
	"github.com/open-feature/go-sdk/pkg/openfeature"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestMain(m *testing.M) {
	err := metricsapi.AddToScheme(scheme.Scheme)
	if err != nil {
		panic("BAD SCHEME!")
	}
	code := m.Run()
	os.Exit(code)
}

func TestMetricServer_disabledServer(t *testing.T) {

	tInstance := &serverManager{
		ticker:        clock.New().Ticker(3 * time.Second),
		ofClient:      openfeature.NewClient("klt-test"),
		exposeMetrics: false,
		k8sClient:     fake.NewClientBuilder().WithScheme(scheme.Scheme).Build(),
	}
	tInstance.start(context.Background())

	var err error
	require.Eventually(t, func() bool {
		cli := &http.Client{}
		req, _ := http.NewRequestWithContext(context.TODO(), http.MethodGet, "http://localhost:9999/metrics", nil)
		_, err = cli.Do(req)
		return err != nil
	}, 30*time.Second, 3*time.Second)

	require.Contains(t, err.Error(), "connection refused")

}

func TestMetricServer_happyPath(t *testing.T) {

	var metric = metricsapi.KeptnMetric{
		ObjectMeta: v1.ObjectMeta{
			Name:      "sample-metric",
			Namespace: "keptn-lifecycle-toolkit-system",
		},
		Spec: metricsapi.KeptnMetricSpec{
			Provider: metricsapi.ProviderRef{
				Name: "dynatrace",
			},
			Query:                "query",
			FetchIntervalSeconds: 5,
		},
		Status: metricsapi.KeptnMetricStatus{
			Value:    "12",
			RawValue: nil,
			LastUpdated: v1.Time{
				Time: time.Now(),
			},
		},
	}
	k8sClient := fake.NewClientBuilder().WithScheme(scheme.Scheme).WithObjects(&metric).Build()

	tInstance := &serverManager{
		ticker:        clock.New().Ticker(3 * time.Second),
		ofClient:      openfeature.NewClient("klt-test"),
		exposeMetrics: true,
		k8sClient:     k8sClient,
	}
	tInstance.start(context.Background())

	require.Eventually(t, func() bool {
		return tInstance.server != nil
	}, 30*time.Second, time.Second)

	var resp *http.Response
	var err error

	require.Eventually(t, func() bool {
		cli := &http.Client{}
		req, err2 := http.NewRequestWithContext(context.TODO(), http.MethodGet, "http://localhost:9999/metrics", nil)
		require.Nil(t, err2)
		resp, err = cli.Do(req)
		return err == nil
	}, 10*time.Second, time.Second)

	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	require.Nil(t, err)
	newStr := buf.String()

	require.Contains(t, newStr, "# TYPE sample_metric gauge")

	require.Eventually(t, func() bool {
		cli := &http.Client{}
		req, _ := http.NewRequestWithContext(context.TODO(), http.MethodGet, "http://localhost:9999/api/v1/metrics/keptn-lifecycle-toolkit-system/sample-metric", nil)
		resp, err = cli.Do(req)
		return err == nil
	}, 10*time.Second, time.Second)

	defer resp.Body.Close()

	buf = new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	require.Nil(t, err)
	newStr = buf.String()

	require.Contains(t, newStr, "\"metric\":\"sample-metric\",\"namespace\":\"keptn-lifecycle-toolkit-system\",\"value\":\"12\"")
}

func TestMetricServer_noMetric(t *testing.T) {

	k8sClient := fake.NewClientBuilder().WithScheme(scheme.Scheme).Build()

	tInstance := &serverManager{
		ticker:        clock.New().Ticker(3 * time.Second),
		ofClient:      openfeature.NewClient("klt-test"),
		exposeMetrics: true,
		k8sClient:     k8sClient,
	}
	tInstance.start(context.Background())

	require.Eventually(t, func() bool {
		return tInstance.server != nil
	}, 30*time.Second, time.Second)

	var resp *http.Response
	var err error
	require.Eventually(t, func() bool {
		cli := &http.Client{}
		req, _ := http.NewRequestWithContext(context.TODO(), http.MethodGet, "http://localhost:9999/api/v1/metrics/default/sample", nil)
		resp, err = cli.Do(req)
		return err == nil
	}, 10*time.Second, time.Second)

	defer resp.Body.Close()
	stat := resp.StatusCode
	require.Equal(t, 404, stat)

}
