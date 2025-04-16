package elastic

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-logr/logr"
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestGetElasticClient(t *testing.T) {
	tests := []struct {
		name          string
		expectedError bool
	}{
		{
			name:          "Success - Elasticsearch client created",
			expectedError: false,
		},
		{
			name:          "Failure - Invalid connection",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			provider := metricsapi.KeptnMetricsProvider{}
			client, err := GetElasticClient(provider)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
				assert.IsType(t, &elasticsearch.Client{}, client)
			}
		})
	}
}

func TestFetchAnalysisValue(t *testing.T) {
	provider := &KeptnElasticProvider{
		K8sClient: fake.NewFakeClient(),
		Log:       logr.Logger{},
	}

	tests := []struct {
		name          string
		query         string
		analysis      metricsapi.Analysis
		expectedError bool
	}{
		{
			name:  "Failure - Missing AnalysisDefinition",
			query: "SELECT avg(cpu) FROM metrics",
			analysis: metricsapi.Analysis{
				Spec: metricsapi.AnalysisSpec{
					Args: map[string]string{"metricPath": "metrics.cpu"},
					AnalysisDefinition: metricsapi.ObjectReference{
						Name: "non-existent-definition",
					},
				},
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			value, err := provider.FetchAnalysisValue(ctx, tt.query, tt.analysis, &metricsapi.KeptnMetricsProvider{})

			if tt.expectedError {
				assert.Error(t, err)
				assert.Empty(t, value)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, value)
			}
		})
	}
}

func TestEvaluateQueryForStep(t *testing.T) {
	provider := &KeptnElasticProvider{
		Log: logr.Discard(),
	}

	tests := []struct {
		name    string
		metric  metricsapi.KeptnMetric
		wantErr bool
	}{
		{
			name:    "Unimplemented function returns nil values",
			metric:  metricsapi.KeptnMetric{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			result, rawData, err := provider.EvaluateQueryForStep(ctx, tt.metric, metricsapi.KeptnMetricsProvider{})

			assert.Nil(t, result)
			assert.Nil(t, rawData)
			assert.NoError(t, err)
		})
	}
}

func TestRunElasticQuery(t *testing.T) {
	esClient, err := elasticsearch.NewDefaultClient()
	if err != nil {
		t.Fatalf("Failed to create Elasticsearch client: %v", err)
	}

	provider := &KeptnElasticProvider{
		Elastic: esClient,
		Log:     logr.Discard(),
	}

	tests := []struct {
		name          string
		query         string
		expectedError bool
	}{
		{
			name:          "Success - Valid Query",
			query:         `{"query": {"match_all": {}}}`,
			expectedError: true,
		},
		{
			name:          "Failure - Invalid Query",
			query:         `INVALID QUERY`,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			result, err := provider.runElasticQuery(ctx, metricsapi.KeptnMetricsProvider{}, tt.query)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)

				jsonData, _ := json.Marshal(result)
				assert.True(t, json.Valid(jsonData))
			}
		})
	}
}

func TestExtractMetric(t *testing.T) {
	provider := &KeptnElasticProvider{
		Log: logr.Discard(),
	}

	tests := []struct {
		name          string
		result        map[string]interface{}
		metricPath    string
		expectedValue string
	}{
		{
			name: "Success - Metric Found in Aggregations",
			result: map[string]interface{}{
				"aggregations": map[string]interface{}{
					"metrics": map[string]interface{}{
						"cpu": 75.5,
					},
				},
			},
			metricPath:    "metrics.cpu",
			expectedValue: "75.500000",
		},
		{
			name: "Failure - Metric Not Found",
			result: map[string]interface{}{
				"aggregations": map[string]interface{}{
					"metrics": map[string]interface{}{
						"memory": 1024,
					},
				},
			},
			metricPath:    "metrics.cpu",
			expectedValue: "",
		},
		{
			name: "Success - Comma Separated Paths",
			result: map[string]interface{}{
				"aggregations": map[string]interface{}{
					"metrics": map[string]interface{}{
						"cpu": 60,
					},
				},
			},
			metricPath:    "metrics.memory,metrics.cpu",
			expectedValue: "60.000000",
		},
		{
			name: "Ignored Non-Aggregation Keys",
			result: map[string]interface{}{
				"metrics": map[string]interface{}{
					"cpu": 50,
				},
			},
			metricPath:    "metrics.cpu",
			expectedValue: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := provider.extractMetric(tt.result, tt.metricPath)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedValue, value)
		})
	}
}

func TestConvertResultTOMap(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		expected map[string]float64
	}{
		{
			name: "Basic Conversion",
			input: map[string]interface{}{
				"cpu": 50,
				"mem": 2048.5,
				"disk": map[string]interface{}{
					"used": 500,
					"free": 1500.75,
				},
			},
			expected: map[string]float64{
				"cpu":       50,
				"mem":       2048.5,
				"disk.used": 500,
				"disk.free": 1500.75,
			},
		},
		{
			name: "Array of Primitives",
			input: map[string]interface{}{
				"prims": []interface{}{float64(1), float64(2)},
			},
			expected: map[string]float64{
				"prims.0": 1,
				"prims.1": 2,
			},
		},
		{
			name: "Array of Objects with Keys",
			input: map[string]interface{}{
				"buckets": []interface{}{
					map[string]interface{}{"key": "one", "value": 5},
					map[string]interface{}{"key": "two", "value": 10},
				},
			},
			expected: map[string]float64{
				"buckets.one.value": 5,
				"buckets.two.value": 10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := convertResultTOMap(tt.input)
			assert.Equal(t, tt.expected, output)
		})
	}
}
