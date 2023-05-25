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

// EvaluateQuery fetches the SLI values from dynatrace provider
func (d *keptnDynatraceDQLProvider) EvaluateQuery(ctx context.Context, metric metricsapi.KeptnMetric, provider metricsapi.KeptnMetricsProvider) (string, []byte, error) {
	if err := d.ensureDTClientIsSetUp(ctx, provider); err != nil {
		return "", nil, err
	}
	// submit DQL
	dqlHandler, err := d.postDQL(ctx, metric.Spec.Query)
	if err != nil {
		d.log.Error(err, "Error while posting the DQL query", "query", metric.Spec.Query)
		return "", nil, err
	}
	// attend result
	results, err := d.getDQL(ctx, *dqlHandler)
	if err != nil {
		d.log.Error(err, "Error while waiting for DQL query", "query", dqlHandler)
		return "", nil, err
	}

	// parse result
	if len(results.Records) > 1 {
		d.log.Info("More than a single result, the first one will be used")
	}
	if len(results.Records) == 0 {
		return "", nil, ErrInvalidResult
	}
	r := fmt.Sprintf("%f", results.Records[0].Value.Avg)
	b, err := json.Marshal(results)
	if err != nil {
		d.log.Error(err, "Error marshaling DQL results")
	}
	return r, b, nil
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

func (d *keptnDynatraceDQLProvider) postDQL(ctx context.Context, query string) (*DynatraceDQLHandler, error) {
	d.log.V(10).Info("posting DQL")
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	values := url.Values{}
	values.Add("query", query)

	path := fmt.Sprintf("/platform/storage/query/v0.7/query:execute?%s", values.Encode())

	b, err := d.dtClient.Do(ctx, path, http.MethodPost, []byte(`{}`))
	if err != nil {
		return nil, err
	}

	dqlHandler := &DynatraceDQLHandler{}
	err = json.Unmarshal(b, &dqlHandler)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response %s: %w", string(b), err)
	}
	return dqlHandler, nil
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

	path := fmt.Sprintf("/platform/storage/query/v0.7/query:poll?%s", values.Encode())

	b, err := d.dtClient.Do(ctx, path, http.MethodGet, nil)
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
		err = fmt.Errorf(ErrAPI, result.Error.Message)
		d.log.Error(err, "Error from Dynatrace DQL provider")
		return nil, err
	}
	return result, nil
}
