package datadog

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	"github.com/go-logr/logr"
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var errNoValues = errors.New("no values in query result")
var errNoMetricPoints = errors.New("no metric points in query result")

type KeptnDataDogProvider struct {
	Log        logr.Logger
	HttpClient http.Client
	K8sClient  client.Client
}

func (d *KeptnDataDogProvider) RunAnalysis(ctx context.Context, query string, spec metricsapi.AnalysisSpec, provider *metricsapi.KeptnMetricsProvider) (string, []byte, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	return d.query(ctx, query, *provider, spec.From.Unix(), spec.To.Unix())
}

// EvaluateQuery fetches the SLI values from datadog provider
func (d *KeptnDataDogProvider) EvaluateQuery(ctx context.Context, metric metricsapi.KeptnMetric, provider metricsapi.KeptnMetricsProvider) (string, []byte, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	fromTime, toTime, err := getTimeRange(metric)
	if err != nil {
		return "", nil, err
	}
	return d.query(ctx, metric.Spec.Query, provider, fromTime, toTime)
}

func (d *KeptnDataDogProvider) query(ctx context.Context, query string, provider metricsapi.KeptnMetricsProvider, fromTime int64, toTime int64) (string, []byte, error) {
	qURL := provider.Spec.TargetServer + "/api/v1/query?from=" + strconv.Itoa(int(fromTime)) + "&to=" + strconv.Itoa(int(toTime)) + "&query=" + url.QueryEscape(query)
	apiKeyVal, appKeyVal, err := getDDSecret(ctx, provider, d.K8sClient)
	if err != nil {
		return "", nil, err
	}

	b, err := d.executeQuery(ctx, qURL, apiKeyVal, appKeyVal)
	if err != nil {
		return "", nil, err
	}

	result := datadogV1.MetricsQueryResponse{}
	err = json.Unmarshal(b, &result)
	if err != nil {
		return "", b, err
	}

	if result.Error != nil {
		err = fmt.Errorf("%s", *result.Error)
		return "", b, err
	}

	if len(result.Series) == 0 {
		return "", nil, errNoValues
	}

	points := (result.Series)[0].Pointlist
	if len(points) == 0 {
		return "", b, errNoMetricPoints
	}

	r := d.getSingleValue(points)
	value := strconv.FormatFloat(r, 'g', 5, 64)
	return value, b, nil
}

func (d *KeptnDataDogProvider) EvaluateQueryForStep(ctx context.Context, metric metricsapi.KeptnMetric, provider metricsapi.KeptnMetricsProvider) ([]string, []byte, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	fromTime, toTime, stepInterval, err := getTimeRangeForStep(metric)
	if err != nil {
		return nil, nil, err
	}
	qURL := provider.Spec.TargetServer + "/api/v1/query?from=" + strconv.Itoa(int(fromTime)) + "&to=" + strconv.Itoa(int(toTime)) + "&interval=" + strconv.Itoa(int(stepInterval)) + "&query=" + url.QueryEscape(metric.Spec.Query)
	apiKeyVal, appKeyVal, err := getDDSecret(ctx, provider, d.K8sClient)
	if err != nil {
		return nil, nil, err
	}

	b, err := d.executeQuery(ctx, qURL, apiKeyVal, appKeyVal)
	if err != nil {
		return nil, nil, err
	}

	result := datadogV1.MetricsQueryResponse{}
	err = json.Unmarshal(b, &result)
	if err != nil {
		return nil, b, err
	}

	if result.Error != nil {
		err = fmt.Errorf("%s", *result.Error)
		return nil, b, err
	}

	if len(result.Series) == 0 {
		return nil, nil, errNoValues
	}

	points := (result.Series)[0].Pointlist
	if len(points) == 0 {
		return nil, b, errNoMetricPoints
	}

	r := d.getResultSlice(points)
	return r, b, nil
}

func (d *KeptnDataDogProvider) executeQuery(ctx context.Context, qURL string, apiKeyVal string, appKeyVal string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", qURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Dd-Api-Key", apiKeyVal)
	req.Header.Set("Dd-Application-Key", appKeyVal)

	res, err := d.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := res.Body.Close()
		if err != nil {
			d.Log.Error(err, "could not close request body")
		}
	}()

	b, _ := io.ReadAll(res.Body)
	return b, nil
}

func (d *KeptnDataDogProvider) getSingleValue(points [][]*float64) float64 {
	var sum float64 = 0
	var count uint64 = 0
	for _, point := range points {
		if point[1] != nil {
			sum += *point[1]
			count++
		}
	}
	if count < 1 {
		// cannot dive by zero
		return 0
	}
	return sum / float64(count)
}

func (d *KeptnDataDogProvider) getResultSlice(points [][]*float64) []string {
	resultSlice := make([]string, 0, len(points))
	for _, point := range points {
		if len(point) > 1 && point[1] != nil {
			valueAsString := fmt.Sprintf("%f", *point[1])
			resultSlice = append(resultSlice, valueAsString)
		}
	}
	return resultSlice
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

func getTimeRangeForStep(metric metricsapi.KeptnMetric) (int64, int64, int64, error) {
	intervalDuration, err := time.ParseDuration(metric.Spec.Range.Interval)
	if err != nil {
		return 0, 0, 0, err
	}
	stepDuration, err := time.ParseDuration(metric.Spec.Range.Step)
	if err != nil {
		return 0, 0, 0, err
	}
	return time.Now().Add(-intervalDuration).Unix(), time.Now().Unix(), stepDuration.Milliseconds(), nil
}
