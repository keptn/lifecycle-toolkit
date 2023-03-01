package server

import (
	"bytes"
	"context"
	"net/http"
	"testing"
	"time"

	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha2"
	"github.com/open-feature/go-sdk/pkg/openfeature"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestMetricServer_happyPath(t *testing.T) {
	metric := metricsapi.KeptnMetric{
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
	}

	err := metricsapi.AddToScheme(scheme.Scheme)
	require.Nil(t, err)
	k8sClient := fake.NewClientBuilder().WithScheme(scheme.Scheme).WithObjects(&metric).Build()

	ctx, cancel := context.WithCancel(context.Background())

	StartServerManager(ctx, k8sClient, openfeature.NewClient("klt-test"), true, 3*time.Second)

	require.Eventually(t, func() bool {
		return instance.server != nil
	}, 10*time.Second, time.Second)

	var resp *http.Response

	require.Eventually(t, func() bool {
		cli := &http.Client{}
		req, err2 := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:9999/metrics", nil)
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

	cancel()

	require.Eventually(t, func() bool {
		return instance.server == nil
	}, 10*time.Second, time.Second)
}

func TestMetricServer_disabledServer(t *testing.T) {
	err2 := metricsapi.AddToScheme(scheme.Scheme)
	require.Nil(t, err2)
	k8sClient := fake.NewClientBuilder().WithScheme(scheme.Scheme).Build()

	ctx, cancel := context.WithCancel(context.Background())

	StartServerManager(ctx, k8sClient, openfeature.NewClient("klt-test"), false, 3*time.Second)

	require.Eventually(t, func() bool {
		return instance.server == nil
	}, 30*time.Second, 3*time.Second)

	var err error

	require.Eventually(t, func() bool {
		cli := &http.Client{}
		req, err2 := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:9999/metrics", nil)
		require.Nil(t, err2)
		_, err = cli.Do(req)
		return err != nil
	}, 30*time.Second, 3*time.Second)

	require.Contains(t, err.Error(), "connection refused")

	cancel()

	require.Eventually(t, func() bool {
		return instance.server == nil
	}, 30*time.Second, 3*time.Second)
}
