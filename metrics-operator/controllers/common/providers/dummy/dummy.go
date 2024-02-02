package dummy

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/go-logr/logr"
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1beta1"
)

type KeptnDummyProvider struct {
	Log        logr.Logger
	HttpClient http.Client
}

func (d *KeptnDummyProvider) FetchAnalysisValue(ctx context.Context, query string, analysis metricsapi.Analysis, provider *metricsapi.KeptnMetricsProvider) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	res, _, err := d.query(ctx, query, *provider, analysis.GetFrom().Unix(), analysis.GetTo().Unix())
	return res, err
}

// EvaluateQuery evaluates the query against the random number API endpoint.
func (d *KeptnDummyProvider) EvaluateQuery(ctx context.Context, metric metricsapi.KeptnMetric, provider metricsapi.KeptnMetricsProvider) (string, []byte, error) {
	// create a context for cancelling the request if it takes too long.
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	fromTime, toTime, err := getTimeRange(metric)
	if err != nil {
		return "", nil, err
	}
	return d.query(ctx, metric.Spec.Query, provider, fromTime, toTime)
}

func (d *KeptnDummyProvider) query(ctx context.Context, query string, provider metricsapi.KeptnMetricsProvider, fromTime int64, toTime int64) (string, []byte, error) {
	// create a new request with context
	//baseURL := "http://www.randomnumberapi.com/api/v1.0/"
	qURL := provider.Spec.TargetServer
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, qURL+query+"?min=42&max=43", nil)
	if err != nil {
		d.Log.Error(err, "Error in creating the request")
		return "", nil, err
	}

	// make an http call using the provided client.
	response, err := d.HttpClient.Do(request)
	if err != nil {
		d.Log.Error(err, "Error in making the request")
		return "", nil, err
	}
	defer response.Body.Close()

	// parse the response data
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		d.Log.Error(err, "Error in reading the response")
	}
	responseStr := d.bytesToString(responseData)

	// Return the metric
	return responseStr, responseData, nil

}

func (d *KeptnDummyProvider) EvaluateQueryForStep(ctx context.Context, metric metricsapi.KeptnMetric, provider metricsapi.KeptnMetricsProvider) ([]string, []byte, error) {
	// create a context for cancelling the request if it takes too long.
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	var result []string
	fromTime, toTime, err := getTimeRange(metric)
	if err != nil {
		return result, nil, err
	}
	resStr, resData, err := d.query(ctx, metric.Spec.Query, provider, fromTime, toTime)

	// Append strings to the slice
	result = append(result, resStr)
	return result, resData, err
}

func (d *KeptnDummyProvider) bytesToString(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	return strconv.Itoa(int(data[0]))
}

func getTimeRange(metric metricsapi.KeptnMetric) (int64, int64, error) {
	var intervalInMin string
	if metric.Spec.Range != nil {
		intervalInMin = metric.Spec.Range.Interval
	} else {
		intervalInMin = "5m"
	}
	intervalDuration, err := time.ParseDuration(intervalInMin)
	if err != nil {
		return 0, 0, err
	}
	return time.Now().Add(-intervalDuration).Unix(), time.Now().Unix(), nil
}
