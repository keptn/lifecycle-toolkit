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
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type KeptnDataDogProvider struct {
	Log        logr.Logger
	HttpClient http.Client
	K8sClient  client.Client
}

// EvaluateQuery fetches the SLI values from datadog provider
func (d *KeptnDataDogProvider) EvaluateQuery(ctx context.Context, metric metricsapi.KeptnMetric, provider metricsapi.KeptnMetricsProvider) (string, []byte, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	// Assumed default metric duration as 5 minutes
	// Think a better way to handle this
	intervalInMin := 5
	fromTime := time.Now().Add(time.Duration(-intervalInMin) * time.Minute).Unix()
	toTime := time.Now().Unix()
	qURL := provider.Spec.TargetServer + "/api/v1/query?from=" + strconv.Itoa(int(fromTime)) + "&to=" + strconv.Itoa(int(toTime)) + "&query=" + url.QueryEscape(metric.Spec.Query)
	req, err := http.NewRequestWithContext(ctx, "GET", qURL, nil)
	if err != nil {
		d.Log.Error(err, "Error while creating request")
		return "", nil, err
	}

	apiKeyVal, appKeyVal, err := getDDSecret(ctx, provider, d.K8sClient)
	if err != nil {
		return "", nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Dd-Api-Key", apiKeyVal)
	req.Header.Set("Dd-Application-Key", appKeyVal)

	res, err := d.HttpClient.Do(req)
	if err != nil {
		d.Log.Error(err, "Error while creating request")
		return "", nil, err
	}
	defer func() {
		err := res.Body.Close()
		if err != nil {
			d.Log.Error(err, "Could not close request body")
		}
	}()

	b, _ := io.ReadAll(res.Body)
	result := datadogV1.MetricsQueryResponse{}
	err = json.Unmarshal(b, &result)
	if err != nil {
		d.Log.Error(err, "Error while parsing response")
		return "", nil, err
	}

	if len(result.Series) == 0 {
		d.Log.Info("No values in query result")
		return "", nil, fmt.Errorf("no values in query result")
	}

	points := (result.Series)[0].Pointlist
	if len(points) == 0 {
		d.Log.Info("No metric points in query result")
		return "", nil, fmt.Errorf("no metric points in query result")
	}

	r := d.getSingleValue(points)
	value := strconv.FormatFloat(r, 'g', 5, 64)
	return value, b, nil
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
