package elasticsearch

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-logr/logr"
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestNewElasticProvider(t *testing.T) {
	log := logr.Discard()
	scheme := runtime.NewScheme()
	k8sClient := fake.NewClientBuilder().WithScheme(scheme).Build()
	elasticURL := "http://localhost:9200"

	provider, err := NewElasticProvider(log, k8sClient, elasticURL)

	assert.NoError(t, err)
	assert.NotNil(t, provider)
	assert.Equal(t, log, provider.Log)
	assert.Equal(t, k8sClient, provider.K8sClient)
	assert.NotNil(t, provider.Elastic)
}

func TestKeptnElasticProvider_FetchAnalysisValue(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := `{
			"hits": {
				"total": {
					"value": 42
				}
			}
		}`
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	cfg := elasticsearch.Config{
		Addresses: []string{server.URL},
	}
	es, err := elasticsearch.NewClient(cfg)
	require.NoError(t, err)

	scheme := runtime.NewScheme()
	k8sClient := fake.NewClientBuilder().WithScheme(scheme).Build()

	provider := &KeptnElasticProvider{
		Log:       logr.Discard(),
		K8sClient: k8sClient,
		Elastic:   es,
	}

	ctx := context.Background()
	query := "test_query"
	analysis := metricsapi.Analysis{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-analysis",
		},
		Spec: metricsapi.AnalysisSpec{
			Timeframe: metricsapi.Timeframe{
				From: metav1.NewTime(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
				To:   metav1.NewTime(time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)),
			},
		},
	}
	metricProvider := &metricsapi.KeptnMetricsProvider{}

	result, err := provider.FetchAnalysisValue(ctx, query, analysis, metricProvider)

	assert.NoError(t, err)
	assert.Equal(t, "42", result)
}

func TestKeptnElasticProvider_EvaluateQuery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := `{
			"hits": {
				"total": {
					"value": 100
				}
			}
		}`
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	cfg := elasticsearch.Config{
		Addresses: []string{server.URL},
	}
	es, err := elasticsearch.NewClient(cfg)
	require.NoError(t, err)

	scheme := runtime.NewScheme()
	k8sClient := fake.NewClientBuilder().WithScheme(scheme).Build()

	provider := &KeptnElasticProvider{
		Log:       logr.Discard(),
		K8sClient: k8sClient,
		Elastic:   es,
	}

	ctx := context.Background()
	metric := metricsapi.KeptnMetric{
		Spec: metricsapi.KeptnMetricSpec{
			Query: "test_query",
		},
	}
	metricProvider := metricsapi.KeptnMetricsProvider{}

	result, _, err := provider.EvaluateQuery(ctx, metric, metricProvider)

	assert.NoError(t, err)
	assert.Equal(t, "100", result)
}

func TestKeptnElasticProvider_runElasticQuery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := `{
			"hits": {
				"total": {
					"value": 200
				}
			}
		}`
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	cfg := elasticsearch.Config{
		Addresses: []string{server.URL},
	}
	es, err := elasticsearch.NewClient(cfg)
	require.NoError(t, err)

	scheme := runtime.NewScheme()
	k8sClient := fake.NewClientBuilder().WithScheme(scheme).Build()

	provider := &KeptnElasticProvider{
		Log:       logr.Discard(),
		K8sClient: k8sClient,
		Elastic:   es,
	}

	ctx := context.Background()
	query := "test_query"
	from := time.Now().Add(-1 * time.Hour)
	to := time.Now()

	result, err := provider.runElasticQuery(ctx, query, from, to)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	hits, ok := result["hits"].(map[string]interface{})
	assert.True(t, ok)

	total, ok := hits["total"].(map[string]interface{})
	assert.True(t, ok)

	value, ok := total["value"].(float64)
	assert.True(t, ok)
	assert.Equal(t, float64(200), value)
}

func TestKeptnElasticProvider_extractMetric(t *testing.T) {
	provider := &KeptnElasticProvider{
		Log: logr.Discard(),
	}

	testCases := []struct {
		name          string
		input         map[string]interface{}
		expectedValue string
		expectedError string
	}{
		{
			name: "Valid input",
			input: map[string]interface{}{
				"hits": map[string]interface{}{
					"total": map[string]interface{}{
						"value": 42,
					},
				},
			},
			expectedValue: "42",
			expectedError: "",
		},
		{
			name:          "Missing hits",
			input:         map[string]interface{}{},
			expectedValue: "",
			expectedError: "invalid result format: missing 'hits' field",
		},
		{
			name: "Missing total",
			input: map[string]interface{}{
				"hits": map[string]interface{}{},
			},
			expectedValue: "",
			expectedError: "invalid result format: missing 'total' field in 'hits'",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			value, err := provider.extractMetric(tc.input)

			if tc.expectedError != "" {
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedValue, value)
			}
		})
	}
}
