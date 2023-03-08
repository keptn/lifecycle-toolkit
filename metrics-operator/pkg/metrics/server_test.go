package server

import (
	"bytes"
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha2"
	"github.com/open-feature/go-sdk/pkg/openfeature"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestMain(m *testing.M) {
	cancel := setup()
	code := m.Run()
	cancel()
	os.Exit(code)
}

var k8sClient client.WithWatch

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

	err := k8sClient.Create(context.TODO(), &metric)
	require.Nil(t, err)

	var resp *http.Response

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

}

func TestMetricServer_noMetric(t *testing.T) {

	var resp *http.Response
	var err error
	require.Eventually(t, func() bool {
		cli := &http.Client{}
		req, err2 := http.NewRequestWithContext(context.TODO(), http.MethodGet, "http://localhost:9999/api/v1/metrics/default/sample", nil)
		require.Nil(t, err2)
		resp, err = cli.Do(req)
		return err == nil
	}, 10*time.Second, time.Second)

	defer resp.Body.Close()
	stat := resp.StatusCode
	require.Equal(t, 404, stat)

}

func TestMetricServer_disabledServer(t *testing.T) {

	var err error

	require.Eventually(t, func() bool {
		cli := &http.Client{}
		req, err2 := http.NewRequestWithContext(context.TODO(), http.MethodGet, "http://localhost:9999/metrics", nil)
		require.Nil(t, err2)
		_, err = cli.Do(req)
		return err != nil
	}, 30*time.Second, 3*time.Second)

	require.Contains(t, err.Error(), "connection refused")

}

func setup() context.CancelFunc {
	err2 := metricsapi.AddToScheme(scheme.Scheme)
	if err2 != nil {
		panic("BAD SCHEME!")
	}
	k8sClient = fake.NewClientBuilder().WithScheme(scheme.Scheme).Build()
	ctx, cancel := context.WithCancel(context.Background())

	StartServerManager(ctx, k8sClient, openfeature.NewClient("klt-test"), false, 3*time.Second)
	return cancel
}
