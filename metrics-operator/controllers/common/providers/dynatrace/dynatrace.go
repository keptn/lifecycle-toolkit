package dynatrace

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/go-logr/logr"
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type KeptnDynatraceProvider struct {
	Log        logr.Logger
	HttpClient http.Client
	K8sClient  client.Client
}

type DynatraceResponse struct {
	TotalCount int               `json:"totalCount"`
	Resolution string            `json:"resolution"`
	Result     []DynatraceResult `json:"result"`
	Error      `json:"error"`
}

type DynatraceResult struct {
	MetricId string          `json:"metricId"`
	Data     []DynatraceData `json:"data"`
}

type DynatraceData struct {
	Timestamps []int64    `json:"timestamps"`
	Values     []*float64 `json:"values"`
}

func (d *KeptnDynatraceProvider) FetchAnalysisValue(ctx context.Context, query string, analysis metricsapi.Analysis, provider *metricsapi.KeptnMetricsProvider) (string, error) {
	baseURL := d.normalizeAPIURL(provider.Spec.TargetServer)
	escapedQ := urlEncodeQuery(query)
	qURL := baseURL + "v2/metrics/query?metricSelector=" + escapedQ + "&from=" + analysis.GetFrom().String() + "&to=" + analysis.GetTo().String()
	res, _, err := d.runQuery(ctx, qURL, *provider)
	return res, err
}

// EvaluateQuery fetches the SLI values from dynatrace provider
func (d *KeptnDynatraceProvider) EvaluateQuery(ctx context.Context, metric metricsapi.KeptnMetric, provider metricsapi.KeptnMetricsProvider) (string, []byte, error) {
	baseURL := d.normalizeAPIURL(provider.Spec.TargetServer)

	var qURL string
	if metric.Spec.Range != nil {
		qURL = "metricSelector=" + metric.Spec.Query + "&from=now-" + metric.Spec.Range.Interval
	} else {
		qURL = "metricSelector=" + metric.Spec.Query
	}

	qURL = urlEncodeQuery(qURL)
	qURL = baseURL + "v2/metrics/query?" + qURL

	return d.runQuery(ctx, qURL, provider)
}

func (d *KeptnDynatraceProvider) runQuery(ctx context.Context, qURL string, provider metricsapi.KeptnMetricsProvider) (string, []byte, error) {
	d.Log.Info("Running query: " + qURL)
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	result, b, err := d.performRequest(ctx, provider, qURL)
	if err != nil {
		return "", b, err
	}
	r := fmt.Sprintf("%f", d.getSingleValue(result))
	return r, b, nil
}

func (d *KeptnDynatraceProvider) EvaluateQueryForStep(ctx context.Context, metric metricsapi.KeptnMetric, provider metricsapi.KeptnMetricsProvider) ([]string, []byte, error) {
	if metric.Spec.Range == nil {
		return nil, nil, fmt.Errorf("spec.range is not defined!")
	}
	baseURL := d.normalizeAPIURL(provider.Spec.TargetServer)
	qURL := "metricSelector=" + metric.Spec.Query + "&from=now-" + metric.Spec.Range.Interval + "&resolution=" + metric.Spec.Range.Step

	qURL = urlEncodeQuery(qURL)
	qURL = baseURL + "v2/metrics/query?" + qURL

	d.Log.Info("Running query: " + qURL)
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	result, b, err := d.performRequest(ctx, provider, qURL)
	if err != nil {
		return nil, b, err
	}

	r := d.getResultSlice(result)
	return r, b, nil
}

func (d *KeptnDynatraceProvider) performRequest(ctx context.Context, provider metricsapi.KeptnMetricsProvider, query string) (*DynatraceResponse, []byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", query, nil)
	if err != nil {
		d.Log.Error(err, "Error while creating request")
		return nil, nil, err
	}

	token, err := getDTSecret(ctx, provider, d.K8sClient)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Authorization", "Api-Token "+token)
	res, err := d.HttpClient.Do(req)

	if err != nil {
		d.Log.Error(err, "Error while creating request")
		return nil, nil, err
	}
	defer func() {
		err := res.Body.Close()
		if err != nil {
			d.Log.Error(err, "Could not close request body")
		}
	}()

	// we ignore the error here because we fail later while unmarshalling
	b, _ := io.ReadAll(res.Body)
	result := DynatraceResponse{}
	err = json.Unmarshal(b, &result)
	if err != nil {
		d.Log.Error(err, "Error while parsing response")
		return nil, b, err
	}
	if !reflect.DeepEqual(result.Error, Error{}) {
		err = fmt.Errorf(ErrAPIMsg, result.Error.Message)
		d.Log.Error(err, "Error from Dynatrace provider")
		return nil, b, err
	}
	return &result, b, nil
}

func (d *KeptnDynatraceProvider) normalizeAPIURL(url string) string {
	out := url
	if !strings.HasSuffix(out, "/") {
		out = out + "/"
	}
	if !strings.HasSuffix(out, "api/") {
		out = out + "api/"
	}
	return out
}

func (d *KeptnDynatraceProvider) getSingleValue(result *DynatraceResponse) float64 {
	var sum float64 = 0
	var count uint64 = 0
	for _, r := range result.Result {
		for _, points := range r.Data {
			for _, v := range points.Values {
				if v != nil {
					sum += *v
					count++
				}
			}
		}
	}
	if count < 1 {
		// cannot dive by zero
		return 0
	}
	return sum / float64(count)
}

func (d *KeptnDynatraceProvider) getResultSlice(result *DynatraceResponse) []string {
	totalValues := 0
	for _, r := range result.Result {
		for _, points := range r.Data {
			for _, v := range points.Values {
				if v != nil {
					totalValues++
				}
			}
		}
	}

	// Initialize resultSlice with the correct length
	resultSlice := make([]string, 0, totalValues) // Use a slice with capacity, but length 0
	for _, r := range result.Result {
		for _, points := range r.Data {
			for _, v := range points.Values {
				if v != nil {
					resultSlice = append(resultSlice, fmt.Sprintf("%f", *v))
				}
			}
		}
	}
	return resultSlice
}
