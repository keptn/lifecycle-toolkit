package dynatrace

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
}

type DynatraceResult struct {
	MetricId string          `json:"metricId"`
	Data     []DynatraceData `json:"data"`
}

type DynatraceData struct {
	Timestamps []int64    `json:"timestamps"`
	Values     []*float64 `json:"values"`
}

// EvaluateQuery fetches the SLI values from dynatrace provider
func (d *KeptnDynatraceProvider) EvaluateQuery(ctx context.Context, metric metricsapi.KeptnMetric, provider metricsapi.KeptnMetricsProvider) (string, []byte, error) {
	baseURL := d.normalizeAPIURL(provider.Spec.TargetServer)
	qURL := baseURL + "v2/metrics/query?metricSelector=" + metric.Spec.Query

	d.Log.Info("Running query: " + qURL)
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", qURL, nil)
	if err != nil {
		d.Log.Error(err, "Error while creating request")
		return "", nil, err
	}

	token, err := getDTSecret(ctx, provider, d.K8sClient)
	if err != nil {
		return "", nil, err
	}

	req.Header.Set("Authorization", "Api-Token "+token)
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

	// we ignore the error here because we fail later while unmarshalling
	b, _ := io.ReadAll(res.Body)
	result := DynatraceResponse{}
	err = json.Unmarshal(b, &result)
	if err != nil {
		d.Log.Error(err, "Error while parsing response")
		return "", nil, err
	}

	r := fmt.Sprintf("%f", d.getSingleValue(result))
	return r, b, nil
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

func (d *KeptnDynatraceProvider) getSingleValue(result DynatraceResponse) float64 {
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
