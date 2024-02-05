# Adding a New Metrics Provider for Dummy Endpoint

To create a provider for the dummy endpoint, follow these steps:

1. Fork the [Keptn repo](https://github.com/keptn/lifecycle-toolkit)

2. **Define the Provider Type:** In the `metrics-operator/controllers/common/providers/common.go` file,
 define the constant `KeptnDummyProviderType` with the value `"dummy"`.

    ```go
    const KeptnDummyProviderType = "dummy"
    ```

3. **Implement the Provider:** Create your own new folder inside the
[metrics-operator/controllers/common/providers](https://github.com/keptn/lifecycle-toolkit/tree/main/metrics-operator/controllers/common/providers)
 matching the new service name: dummy and a new Go package for the dummy provider.
  This package should contain
 a struct that implements the `KeptnSLIProvider` interface.
  You can follow other existing implementations,
 such as [prometheus.go](https://github.com/keptn/lifecycle-toolkit/blob/main/metrics-operator/controllers/common/providers/prometheus/prometheus.go),
 as an example.
  Ensure that you implement the `EvaluateQuery` function to fetch the metrics accurately.
  In the implementation, make a request to the dummy endpoint and return the response.

```go
// Inside the dummy package

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
```

4. **Instantiate the Provider:** In the `providers.NewProvider` function
 in the `metrics-operator/controllers/common/providers/provider.go` file,
 add a case for the `KeptnDummyProviderType`.
  Instantiate the dummy provider struct and return it.

    ```go
    // Inside the providers package

    // NewProvider function
    func NewProvider(providerType string, log logr.Logger, k8sClient client.Client) (KeptnSLIProvider, error) {
        switch strings.ToLower(providerType) {
        case KeptnDummyProviderType:
            return &dummy.KeptnDummyProvider{
                Log:        log,
                HttpClient: http.Client{},
            }, nil
        // Other cases...
        }
    }
    ```

5. **Add Test Cases:** Write test cases to validate your implementation and ensure it works correctly.
 This step is crucial for maintaining code quality and reliability.

6. **Test:** Thoroughly test your implementation to verify that it functions as expected.
 Make sure to cover various scenarios and edge cases to ensure robustness.
