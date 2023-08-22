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
	Records []DQLRecord `json:"records"`
}

type DQLRecord struct {
	Value DQLMetric `json:"value"`
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
	DefaultTimeframeStart      string `json:"defaultTimeframeStart"`
	DefaultTimeframeEnd        string `json:"defaultTimeframeEnd"`
	Timezone                   string `json:"timezone"`
	Locale                     string `json:"locale"`
	FetchTimeoutSeconds        int    `json:"fetchTimeoutSeconds"`
	RequestTimeoutMilliseconds int    `json:"requestTimeoutMilliseconds"`
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

func (d *keptnDynatraceDQLProvider) EvaluateQuery(ctx context.Context, metric metricsapi.KeptnMetric, provider metricsapi.KeptnMetricsProvider) (string, []byte, error) {
	if err := d.ensureDTClientIsSetUp(ctx, provider); err != nil {
		return "", nil, err
	}

	b, status, err := d.postDQL(ctx, metric)
	if err != nil {
		d.log.Error(err, "Error while posting the DQL query", "query", metric.Spec.Query)
		return "", nil, err
	}

	results, err := d.parseDQLResults(b, status)
	if err != nil {
		return "", nil, err
	}

	if len(results.Records) > 1 {
		d.log.Info("More than a single result, the first one will be used")
	}
	if len(results.Records) == 0 {
		return "", nil, ErrInvalidResult
	}
	r := fmt.Sprintf("%f", results.Records[0].Value.Avg)
	b, err = json.Marshal(results)
	if err != nil {
		d.log.Error(err, "Error marshaling DQL results")
	}

	return r, b, nil
}

func (d *keptnDynatraceDQLProvider) EvaluateQueryForStep(ctx context.Context, metric metricsapi.KeptnMetric, provider metricsapi.KeptnMetricsProvider) ([]string, []byte, error) {
	if err := d.ensureDTClientIsSetUp(ctx, provider); err != nil {
		return nil, nil, err
	}

	b, status, err := d.postDQL(ctx, metric)
	if err != nil {
		d.log.Error(err, "Error while posting the DQL query", "query", metric.Spec.Query)
		return nil, nil, err
	}

	results, err := d.parseDQLResults(b, status)
	if err != nil {
		return nil, nil, err
	}

	r := d.getResultSlice(results)
	b, err = json.Marshal(results)
	if err != nil {
		d.log.Error(err, "Error marshaling DQL results")
	}

	return r, b, nil
}

func (d *keptnDynatraceDQLProvider) parseDQLResults(b []byte, status int) (*DQLResult, error) {
	results := &DQLResult{}
	if status == 200 {
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

func (d *keptnDynatraceDQLProvider) postDQL(ctx context.Context, metric metricsapi.KeptnMetric) ([]byte, int, error) {
	d.log.V(10).Info("posting DQL")
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	path := "/platform/storage/query/v0.7/query:execute"

	payload := DQLRequest{
		Query:                      metric.Spec.Query,
		DefaultTimeframeStart:      "",
		DefaultTimeframeEnd:        "",
		Timezone:                   "UTC",
		Locale:                     "en_US",
		FetchTimeoutSeconds:        60,
		RequestTimeoutMilliseconds: 1000,
	}

	if metric.Spec.Range != nil {
		intervalDuration, err := time.ParseDuration(metric.Spec.Range.Interval)
		if err != nil {
			return nil, 0, err
		}
		payload.DefaultTimeframeStart = time.Now().Add(-intervalDuration).Format(time.RFC3339)
		payload.DefaultTimeframeEnd = time.Now().Format(time.RFC3339)
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

	path := fmt.Sprintf("/platform/storage/query/v1/query:poll?%s", values.Encode())

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

func (d *keptnDynatraceDQLProvider) getResultSlice(result *DQLResult) []string {
	if len(result.Records) == 0 {
		return nil
	}
	// Initialize resultSlice with the correct length
	resultSlice := make([]string, 0, len(result.Records)) // Use a slice with capacity, but length 0
	for _, r := range result.Records {
		resultSlice = append(resultSlice, fmt.Sprintf("%f", r.Value.Max))
	}
	return resultSlice
}
