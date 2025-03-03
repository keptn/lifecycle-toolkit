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
	defaultTimeRange        = 30 * time.Minute
)

type KeptnElasticProvider struct {
	Log       logr.Logger
	K8sClient client.Client
	Elastic   *elastic.Client
}

type ElasticsearchResponse struct {
	Hits struct {
		Total struct {
			Value int `json:"value"`
		} `json:"total"`
	} `json:"hits"`
}

func NewElasticProvider(log logr.Logger, k8sClient client.Client, elasticURL string) (*KeptnElasticProvider, error) {
	log.Info("Initializing Elasticsearch client with TLS disabled...")

	// Custom Transport to Skip TLS Verification
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	// Elasticsearch Client Config
	cfg := elastic.Config{
		Addresses: []string{"https://quickstart-es-http:9200"},
		Username:  "elastic",
		Password:  "073wS1l5ct65LdktJ03M9oqa",
		Transport: transport,
	}

	// Create Elasticsearch Client
	es, err := elastic.NewClient(cfg)
	if err != nil {
		log.Error(err, "Failed to create Elasticsearch client")
		return nil, fmt.Errorf("failed to create Elasticsearch client: %w", err)
	}

	log.Info("Successfully initialized Elasticsearch client with TLS disabled")
	return &KeptnElasticProvider{
		Log:       log,
		K8sClient: k8sClient,
		Elastic:   es,
	}, nil
}

func (r *KeptnElasticProvider) FetchAnalysisValue(ctx context.Context, query string, analysis metricsapi.Analysis, provider *metricsapi.KeptnMetricsProvider) (string, error) {
	r.Log.Info("Fetching analysis value from Elasticsearch", "query", query)
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	result, err := r.runElasticQuery(ctx, query, analysis.GetFrom(), analysis.GetTo())
	if err != nil {
		r.Log.Error(err, "Failed to fetch analysis value")
		return "", err
	}

	r.Log.Info("Elasticsearch query result", "result", result)
	return r.extractMetric(result)
}

func (r *KeptnElasticProvider) EvaluateQuery(ctx context.Context, metric metricsapi.KeptnMetric, provider metricsapi.KeptnMetricsProvider) (string, []byte, error) {
	r.Log.Info("Evaluating query for KeptnMetric", "metric", metric.Name)
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	timeRange := getTimeRangeFromSpec(metric.Spec.Range)

	result, err := r.runElasticQuery(ctx, metric.Spec.Query, time.Now().Add(-timeRange), time.Now())
	if err != nil {
		r.Log.Error(err, "Failed to evaluate query")
		return "", nil, err
	}

	metricValue, err := r.extractMetric(result)
	if err != nil {
		return "", nil, err
	}

	r.Log.Info("Successfully evaluated metric", "value", metricValue)
	return metricValue, []byte{}, nil
}

func (r *KeptnElasticProvider) EvaluateQueryForStep(ctx context.Context, metric metricsapi.KeptnMetric, provider metricsapi.KeptnMetricsProvider) ([]string, []byte, error) {
	r.Log.Info("EvaluateQueryForStep called but not implemented")
	return nil, nil, nil
}

func getTimeRangeFromSpec(rangeSpec *metricsapi.RangeSpec) time.Duration {
	if rangeSpec == nil || rangeSpec.Interval == "" {
		return defaultTimeRange
	}
	duration, err := time.ParseDuration(rangeSpec.Interval)
	if err != nil {
		return defaultTimeRange
	}
	return duration
}

func (r *KeptnElasticProvider) runElasticQuery(ctx context.Context, query string, from, to time.Time) (map[string]interface{}, error) {
	r.Log.Info("Running Elasticsearch query", "query", query, "from", from, "to", to)

	queryBody := fmt.Sprintf(`
	{
		"query": {
			"bool": {
				"must": [
					%s,
					{
						"range": {
							"@timestamp": {
								"gte": "%s",
								"lte": "%s"
							}
						}
					}
				]
			}
		}
	}`, query, from.Format(time.RFC3339), to.Format(time.RFC3339))

	res, err := r.Elastic.Search(
		r.Elastic.Search.WithContext(ctx),
		r.Elastic.Search.WithBody(strings.NewReader(queryBody)),
	)
	if err != nil {
		r.Log.Error(err, "Failed to execute Elasticsearch query")
		return nil, fmt.Errorf("failed to execute Elasticsearch query: %w", err)
	}
	defer res.Body.Close()

	if warnings, ok := res.Header["Warning"]; ok {
		r.Log.Info(fmt.Sprintf(warningLogStringElastic, "Elasticsearch", warnings))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		r.Log.Error(err, "Failed to parse Elasticsearch response")
		return nil, fmt.Errorf("failed to parse Elasticsearch response: %w", err)
	}

	r.Log.Info("Successfully executed query", "result", result)
	return result, nil
}

func (r *KeptnElasticProvider) extractMetric(result map[string]interface{}) (string, error) {
	r.Log.Info("Extracting metric from Elasticsearch response")
	var response ElasticsearchResponse
	jsonData, err := json.Marshal(result)
	if err != nil {
		r.Log.Error(err, "Failed to marshal Elasticsearch result")
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	if err := json.Unmarshal(jsonData, &response); err != nil {
		r.Log.Error(err, "Failed to unmarshal result into struct")
		return "", fmt.Errorf("failed to unmarshal result into struct: %w", err)
	}

	value := fmt.Sprintf("%d", response.Hits.Total.Value)
	r.Log.Info("Extracted metric successfully", "value", value)
	return value, nil
}
