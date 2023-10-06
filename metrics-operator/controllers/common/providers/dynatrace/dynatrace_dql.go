package dynatrace

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/go-logr/logr"
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	dtclient "github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/providers/dynatrace/client"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const maxRetries = 5
const retryFetchInterval = 10 * time.Second

const dqlQuerySucceeded = "SUCCEEDED"
const defaultPath = "/platform/storage/query/v1/query:"

type keptnDynatraceDQLProvider struct {
	log       logr.Logger
	k8sClient client.Client

	dtClient dtclient.DTAPIClient
	clock    clock.Clock
}

type DynatraceDQLHandler struct {
	RequestToken string `json:"requestToken"`
}

type DynatraceDQLResult struct {
	State  string    `json:"state"`
	Result DQLResult `json:"result,omitempty"`
	Error  `json:"error"`
}

type DQLResult struct {
	Records []map[string]any `json:"records"`
}

type DQLMetric struct {
	Count int64   `json:"count"`
	Sum   float64 `json:"sum"`
	Min   float64 `json:"min"`
	Avg   float64 `json:"avg"`
	Max   float64 `json:"max"`
}

type DQLRequest struct {
	Query                      string `json:"query"`
	DefaultTimeframeStart      string `json:"defaultTimeframeStart,omitempty"`
	DefaultTimeframeEnd        string `json:"defaultTimeframeEnd,omitempty"`
	Timezone                   string `json:"timezone"`
	Locale                     string `json:"locale"`
	FetchTimeoutSeconds        int    `json:"fetchTimeoutSeconds"`
	RequestTimeoutMilliseconds int    `json:"requestTimeoutMilliseconds"`
}

type timeframe struct {
	from time.Time
	to   time.Time
}

type metricRequest struct {
	query     string
	timeframe *timeframe
}

func newMetricRequestFromMetric(metric metricsapi.KeptnMetric) (*metricRequest, error) {
	res := &metricRequest{
		query: metric.Spec.Query,
	}

	if metric.Spec.Range != nil {
		intervalDuration, err := time.ParseDuration(metric.Spec.Range.Interval)
		if err != nil {
			return nil, err
		}
		res.timeframe = &timeframe{
			from: time.Now().UTC().Add(-intervalDuration),
			to:   time.Now().UTC(),
		}
	}

	return res, nil
}

func newMetricRequestFromAnalysis(query string, analysis metricsapi.Analysis) (*metricRequest, error) {
	res := &metricRequest{
		query: query,
		timeframe: &timeframe{
			from: analysis.GetFrom(),
			to:   analysis.GetTo(),
		},
	}

	return res, nil
}

type KeptnDynatraceDQLProviderOption func(provider *keptnDynatraceDQLProvider)

func WithDTAPIClient(dtApiClient dtclient.DTAPIClient) KeptnDynatraceDQLProviderOption {
	return func(provider *keptnDynatraceDQLProvider) {
		provider.dtClient = dtApiClient
	}
}

func WithLogger(logger logr.Logger) KeptnDynatraceDQLProviderOption {
	return func(provider *keptnDynatraceDQLProvider) {
		provider.log = logger
	}
}

// NewKeptnDynatraceDQLProvider creates and returns a new KeptnDynatraceDQLProvider
func NewKeptnDynatraceDQLProvider(k8sClient client.Client, opts ...KeptnDynatraceDQLProviderOption) *keptnDynatraceDQLProvider {
	provider := &keptnDynatraceDQLProvider{
		log:       logr.New(klog.NewKlogr().GetSink()),
		k8sClient: k8sClient,
		clock:     clock.New(),
	}

	for _, o := range opts {
		o(provider)
	}

	return provider
}

func (d *keptnDynatraceDQLProvider) FetchAnalysisValue(ctx context.Context, query string, analysis metricsapi.Analysis, provider *metricsapi.KeptnMetricsProvider) (string, error) {
	metricsReq, err := newMetricRequestFromAnalysis(query, analysis)
	if err != nil {
		return "", err
	}

	results, err := d.getResults(ctx, *metricsReq, *provider)
	if err != nil {
		return "", err
	}

	if len(results.Records) > 1 {
		d.log.Info("More than a single result, the first one will be used")
	}

	value := extractValueFromRecord(results.Records[0])

	return value, nil
}

func (d *keptnDynatraceDQLProvider) EvaluateQuery(ctx context.Context, metric metricsapi.KeptnMetric, provider metricsapi.KeptnMetricsProvider) (string, []byte, error) {
	metricsReq, err := newMetricRequestFromMetric(metric)
	if err != nil {
		return "", nil, err
	}
	results, err := d.getResults(ctx, *metricsReq, provider)
	if err != nil {
		return "", nil, err
	}

	if len(results.Records) == 0 {
		return "", nil, ErrInvalidResult
	}

	if len(results.Records) > 1 {
		d.log.Info("More than a single result, the first one will be used")
	}

	value := extractValueFromRecord(results.Records[0])

	b, err := json.Marshal(results)
	if err != nil {
		d.log.Error(err, "Error marshaling DQL results")
	}

	return value, b, nil
}

func (d *keptnDynatraceDQLProvider) EvaluateQueryForStep(ctx context.Context, metric metricsapi.KeptnMetric, provider metricsapi.KeptnMetricsProvider) ([]string, []byte, error) {
	metricsReq, err := newMetricRequestFromMetric(metric)
	if err != nil {
		return nil, nil, err
	}
	results, err := d.getResults(ctx, *metricsReq, provider)
	if err != nil {
		return nil, nil, err
	}

	r := extractValuesFromRecord(results.Records[0])
	b, err := json.Marshal(results)
	if err != nil {
		d.log.Error(err, "Error marshaling DQL results")
	}

	return r, b, nil
}

func (d *keptnDynatraceDQLProvider) getResults(ctx context.Context, metricsReq metricRequest, provider metricsapi.KeptnMetricsProvider) (*DQLResult, error) {
	if err := d.ensureDTClientIsSetUp(ctx, provider); err != nil {
		return nil, err
	}

	b, status, err := d.postDQL(ctx, metricsReq)
	if err != nil {
		d.log.Error(err, "Error while posting the DQL query", "query", metricsReq.query)
		return nil, err
	}

	results, err := d.parseDQLResults(b, status)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (d *keptnDynatraceDQLProvider) parseDQLResults(b []byte, status int) (*DQLResult, error) {
	results := &DQLResult{}
	if status == http.StatusOK {
		r := &DynatraceDQLResult{}
		err := json.Unmarshal(b, &r)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal response %s: %w", string(b), err)
		}
		if r.State == dqlQuerySucceeded {
			results = &r.Result
		}
	} else {
		dqlHandler := &DynatraceDQLHandler{}
		err := json.Unmarshal(b, &dqlHandler)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal response %s: %w", string(b), err)
		}
		results, err = d.getDQL(context.Background(), *dqlHandler)
		if err != nil {
			d.log.Error(err, "Error while waiting for DQL query", "query", dqlHandler)
			return nil, err
		}
	}

	if len(results.Records) == 0 {
		return nil, ErrInvalidResult
	}

	return results, nil
}

func (d *keptnDynatraceDQLProvider) ensureDTClientIsSetUp(ctx context.Context, provider metricsapi.KeptnMetricsProvider) error {
	// try to initialize the DT API Client if it has not been set in the options
	if d.dtClient == nil {
		secret, err := getDTSecret(ctx, provider, d.k8sClient)
		if err != nil {
			return err
		}
		config, err := dtclient.NewAPIConfig(
			provider.Spec.TargetServer,
			secret,
		)
		if err != nil {
			return err
		}
		d.dtClient = dtclient.NewAPIClient(*config, dtclient.WithLogger(d.log))
	}
	return nil
}

func (d *keptnDynatraceDQLProvider) postDQL(ctx context.Context, metricsReq metricRequest) ([]byte, int, error) {
	d.log.V(10).Info("posting DQL")
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	path := defaultPath + "execute"

	payload := DQLRequest{
		Query:                      metricsReq.query,
		DefaultTimeframeStart:      "",
		DefaultTimeframeEnd:        "",
		Timezone:                   "UTC",
		Locale:                     "en_US",
		FetchTimeoutSeconds:        60,
		RequestTimeoutMilliseconds: 1000,
	}

	if metricsReq.timeframe != nil {
		payload.DefaultTimeframeStart = metricsReq.timeframe.from.Format(time.RFC3339)
		payload.DefaultTimeframeEnd = metricsReq.timeframe.to.Format(time.RFC3339)
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, 0, err
	}

	b, status, err := d.dtClient.Do(ctx, path, http.MethodPost, payloadBytes)
	if err != nil {
		return nil, 0, err
	}
	return b, status, nil
}

func (d *keptnDynatraceDQLProvider) getDQL(ctx context.Context, handler DynatraceDQLHandler) (*DQLResult, error) {
	d.log.V(10).Info("posting DQL")

	for i := 0; i < maxRetries; i++ {
		r, err := d.retrieveDQLResults(ctx, handler)
		if err != nil {
			return &DQLResult{}, err
		}
		if r.State == dqlQuerySucceeded {
			return &r.Result, nil
		}
		d.log.V(10).Info("DQL not finished, got", "state", r.State)
		<-d.clock.After(retryFetchInterval)
	}
	return nil, ErrDQLQueryTimeout
}

func (d *keptnDynatraceDQLProvider) retrieveDQLResults(ctx context.Context, handler DynatraceDQLHandler) (*DynatraceDQLResult, error) {
	d.log.V(10).Info("Getting DQL")
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	values := url.Values{}
	values.Add("request-token", handler.RequestToken)

	path := defaultPath + fmt.Sprintf("poll?%s", values.Encode())

	b, _, err := d.dtClient.Do(ctx, path, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	result := &DynatraceDQLResult{}
	err = json.Unmarshal(b, &result)
	if err != nil {
		d.log.Error(err, "Error while parsing response")
		return result, err
	}

	if !reflect.DeepEqual(result.Error, Error{}) {
		err = fmt.Errorf(ErrAPIMsg, result.Error.Message)
		d.log.Error(err, "Error from Dynatrace DQL provider")
		return nil, err
	}
	return result, nil
}

// extractValueFromRecord extracts the latest value of a record.
// This is intended for timeseries queries that return a single metric
func extractValueFromRecord(record map[string]any) string {
	for _, item := range record {
		if valuesArr, ok := toFloatArray(item); ok {
			return fmt.Sprintf("%f", valuesArr[len(valuesArr)-1])
		}
	}
	return ""
}

// extractValuesFromRecord extracts all values of a record.
// This is intended for timeseries queries that return multiple values for a single metric, i.e. the individual
// data points of a time series
func extractValuesFromRecord(record map[string]any) []string {
	for _, item := range record {
		if valuesArr, ok := toFloatArray(item); ok {
			valuesStrArr := make([]string, len(valuesArr))
			for index, val := range valuesArr {
				valuesStrArr[index] = fmt.Sprintf("%f", val)
			}
			return valuesStrArr
		}
	}
	return []string{}
}

func toFloatArray(obj any) ([]float64, bool) {
	valuesArr, ok := obj.([]any)
	if !ok {
		return nil, false
	}
	res := make([]float64, len(valuesArr))
	for index, val := range valuesArr {
		if floatVal, ok := val.(float64); ok {
			res[index] = floatVal
		} else if intVal, ok := val.(int); ok {
			res[index] = float64(intVal)
		} else {
			return nil, false
		}
	}
	return res, true
}
