package elastic

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	elastic "github.com/elastic/go-elasticsearch/v8"
	"github.com/go-logr/logr"
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	warningLogStringElastic = "%s API returned warnings: %s"
)

type KeptnElasticProvider struct {
	Log       logr.Logger
	K8sClient client.Client
	Elastic   *elastic.Client
}

// GetElasticClient will create a new elastic client
func GetElasticClient(provider metricsapi.KeptnMetricsProvider) (*elastic.Client, error) {
	es, err := elastic.NewClient(elastic.Config{
		Addresses: []string{provider.Spec.TargetServer},
		APIKey:    provider.Spec.SecretKeyRef.Key,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: provider.Spec.InsecureSkipTlsVerify,
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Elasticsearch client: %w", err)
	}
	return es, nil
}

// FetchAnalysisValue will fetch analysis value depends on query and the metrics provided as input
func (r *KeptnElasticProvider) FetchAnalysisValue(ctx context.Context, query string, analysis metricsapi.Analysis, provider *metricsapi.KeptnMetricsProvider) (string, error) {
	// Retrieve the AnalysisDefinition referenced in Analysis
	var analysisDef metricsapi.AnalysisDefinition
	err := r.K8sClient.Get(ctx, client.ObjectKey{
		Name:      analysis.Spec.AnalysisDefinition.Name,
		Namespace: analysis.Namespace,
	}, &analysisDef)

	if err != nil {
		r.Log.Error(err, "Failed to retrieve AnalysisDefinition")
		return "", fmt.Errorf("failed to get AnalysisDefinition: %w", err)
	}

	// Extract the referenced AnalysisValueTemplate name
	if len(analysisDef.Spec.Objectives) == 0 {
		return "", fmt.Errorf("no objectives defined in AnalysisDefinition")
	}

	templateName := analysisDef.Spec.Objectives[0].AnalysisValueTemplateRef.Name
	r.Log.Info("Found referenced AnalysisValueTemplate", "templateName", templateName)
	// Retrieve the AnalysisValueTemplate using the extracted name
	var template metricsapi.AnalysisValueTemplate
	err = r.K8sClient.Get(ctx, client.ObjectKey{
		Name:      templateName,
		Namespace: analysis.Namespace,
	}, &template)

	if err != nil {
		r.Log.Error(err, "Failed to retrieve AnalysisValueTemplate")
		return "", fmt.Errorf("failed to get AnalysisValueTemplate: %w", err)
	}

	// Extract metricPath from args
	metricPathStr, exists := analysis.Spec.Args["metricPath"]
	if !exists || metricPathStr == "" {
		return "", fmt.Errorf("metric path is missing in AnalysisValueTemplate annotations")
	}

	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	result, err := r.runElasticQuery(ctx, *provider, query)
	if err != nil {
		return "", err
	}

	r.Log.Info("Elasticsearch query result", "result", result)
	return r.extractMetric(result, metricPathStr)
}

// EvaluateQuery takes query as a input but doesn't return anything
func (r *KeptnElasticProvider) EvaluateQuery(ctx context.Context, metric metricsapi.KeptnMetric, provider metricsapi.KeptnMetricsProvider) (string, []byte, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	result, err := r.runElasticQuery(ctx, provider, metric.Spec.Query)
	if err != nil {
		return "", nil, err
	}

	metricValue, err := r.extractMetric(result, "")
	if err != nil {
		return "", nil, err
	}

	return metricValue, []byte{}, nil
}

func (r *KeptnElasticProvider) EvaluateQueryForStep(ctx context.Context, metric metricsapi.KeptnMetric, provider metricsapi.KeptnMetricsProvider) ([]string, []byte, error) {
	r.Log.Info("EvaluateQueryForStep called but not implemented")
	return nil, nil, nil
}

// runElasticQuery runs query on elastic search to get output from elasticsearch
func (r *KeptnElasticProvider) runElasticQuery(ctx context.Context, provider metricsapi.KeptnMetricsProvider, query string) (map[string]interface{}, error) {
	var err error
	r.Log.Info("Running Elasticsearch query", "query", query)
	r.Elastic, err = GetElasticClient(provider)
	if err != nil {
		return nil, err
	}

	res, err := r.Elastic.Search(
		r.Elastic.Search.WithContext(ctx),
		r.Elastic.Search.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to execute Elasticsearch query: %w", err)
	}
	defer res.Body.Close()

	if warnings, ok := res.Header["Warning"]; ok {
		r.Log.Info(fmt.Sprintf(warningLogStringElastic, "Elasticsearch", warnings))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse Elasticsearch response: %w", err)
	}

	return result, nil
}

// extractMetric will parse the result and return the metrics which we input to the function
func (r *KeptnElasticProvider) extractMetric(result map[string]interface{}, metricPathStr string) (string, error) {
	convertedResult := convertResultTOMap(result)
	for k, v := range convertedResult {
		if strings.Contains(k, metricPathStr) {
			return fmt.Sprintf("%f", v), nil
		}
	}
	return "", nil
}

// convertResultTOMap recursively converts map[string]interface{} to map[string]float64
func convertResultTOMap(input map[string]interface{}) map[string]float64 {
	output := make(map[string]float64)
	for key, value := range input {
		switch v := value.(type) {
		case float64:
			output[key] = v
		case float32:
			output[key] = float64(v)
		case int:
			output[key] = float64(v)
		case int32:
			output[key] = float64(v)
		case int64:
			output[key] = float64(v)
		case map[string]interface{}:
			nestedMap := convertResultTOMap(v)
			for nestedKey, nestedValue := range nestedMap {
				output[key+"."+nestedKey] = nestedValue
			}
		default:
		}
	}
	return output
}
