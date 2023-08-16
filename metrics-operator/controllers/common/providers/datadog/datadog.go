package datadog

import (
	"context"
	"encoding/json"
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

var errCreateReq = "Error while creating request"
var errNoValues = "No values in query result"
var errNoMetricPoints = "No metric points in query result"
var errParseRes = "Error while parsing response"
var errCloseBody = "Could not close request body"
var errDatadogProvider = "Error from DataDog provider"

type KeptnDataDogProvider struct {
	Log        logr.Logger
	HttpClient http.Client
	K8sClient  client.Client
}

// EvaluateQuery fetches the SLI values from datadog provider
func (d *KeptnDataDogProvider) EvaluateQuery(ctx context.Context, metric metricsapi.KeptnMetric, provider metricsapi.KeptnMetricsProvider) (string, []byte, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	fromTime, toTime, err := getTimeRange(metric)
	if err != nil {
		return "", nil, err
	}
	qURL := provider.Spec.TargetServer + "/api/v1/query?from=" + strconv.Itoa(int(fromTime)) + "&to=" + strconv.Itoa(int(toTime)) + "&query=" + url.QueryEscape(metric.Spec.Query)
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
		d.Log.Error(err, errParseRes)
		return "", b, err
	}

	if result.Error != nil {
		err = fmt.Errorf("%s", *result.Error)
		d.Log.Error(err, errDatadogProvider)
		return "", b, err
	}

	if len(result.Series) == 0 {
		d.Log.Info(errNoValues)
		return "", nil, fmt.Errorf(errNoValues)
	}

	points := (result.Series)[0].Pointlist
	if len(points) == 0 {
		d.Log.Info(errNoMetricPoints)
		return "", b, fmt.Errorf(errNoMetricPoints)
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
		d.Log.Error(err, errParseRes)
		return nil, b, err
	}

	if result.Error != nil {
		err = fmt.Errorf("%s", *result.Error)
		d.Log.Error(err, errDatadogProvider)
		return nil, b, err
	}

	if len(result.Series) == 0 {
		d.Log.Info(errNoValues)
		return nil, nil, fmt.Errorf(errNoValues)
	}

	points := (result.Series)[0].Pointlist
	if len(points) == 0 {
		d.Log.Info(errNoMetricPoints)
		return nil, b, fmt.Errorf(errNoMetricPoints)
	}

	r := d.getResultSlice(points)
	return r, b, nil
}

func (d *KeptnDataDogProvider) executeQuery(ctx context.Context, qURL string, apiKeyVal string, appKeyVal string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", qURL, nil)
	if err != nil {
		d.Log.Error(err, errCreateReq)
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Dd-Api-Key", apiKeyVal)
	req.Header.Set("Dd-Application-Key", appKeyVal)

	res, err := d.HttpClient.Do(req)
	if err != nil {
		d.Log.Error(err, errCreateReq)
		return nil, err
	}
	defer func() {
		err := res.Body.Close()
		if err != nil {
			d.Log.Error(err, errCloseBody)
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
		if point[1] != nil {
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
