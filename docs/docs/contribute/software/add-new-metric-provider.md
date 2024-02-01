
Copy code
## Adding a New Metrics Provider for Dummy Endpoint

To create a provider for the dummy endpoint, follow these steps:

1. **Define the Provider Type:** In the `metrics-operator/controllers/common/providers/common.go` file, define the constant `KeptnDummyProviderType` with the value `"dummy"`.

    ```go
    const KeptnDummyProviderType = "dummy"
    ```

2. **Implement the Provider:** Create your own new folder inside [this folder](https://github.com/keptn/lifecycle-toolkit/tree/main/metrics-operator/controllers/common/providers) matching the new service name: dummy and a new Go package for the dummy provider. This package should contain a struct that implements the `KeptnSLIProvider` interface. In the implementation, make a request to the dummy endpoint and return the response.

    ```go
    // Inside the dummy package

    type KeptnDummyProvider struct {
        Log        logr.Logger
        HttpClient http.Client
    }

    func (d *KeptnDummyProvider) FetchAnalysisValue(ctx context.Context, query string, analysis metricsapi.Analysis, provider *metricsapi.KeptnMetricsProvider) (string, error) {
        ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
        defer cancel()
        res, err := d.query(ctx, query, *provider, analysis.GetFrom().Unix(), analysis.GetTo().Unix())
        return string(res), err
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
        res, err := d.query(ctx, metric.Spec.Query, provider, fromTime, toTime)
        return string(res), res, err
    }

    func (r *KeptnDummyProvider) query(ctx context.Context, query string, provider metricsapi.KeptnMetricsProvider, fromTime int64, toTime int64) ([]byte, error) {
        // create a new request with context
        baseURL := "http://www.randomnumberapi.com/api/v1.0/"

        request, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+query, nil)
        if err != nil {
            r.Log.Error(err, "Error in creating the request")
            return nil, err
        }

        // make an http call using the provided client.
        response, err := r.HttpClient.Do(request)
        if err != nil {
            r.Log.Error(err, "Error in making the request")
            return nil, err
        }
        defer response.Body.Close()

        // parse the response data
        responseData, err := io.ReadAll(response.Body)
        if err != nil {
            r.Log.Error(err, "Error in reading the response")
        }

        // return the metric
        return responseData, nil
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
        res, err := d.query(ctx, metric.Spec.Query, provider, fromTime, toTime)

        // Append strings to the slice
        result = append(result, string(res))
        return result, res, err
    }
    ```

3. **Instantiate the Provider:** In the `providers.NewProvider` function in the `metrics-operator/controllers/common/providers/provider.go` file, add a case for the `KeptnDummyProviderType`. Instantiate the dummy provider struct and return it.

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
