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
  To fully implement the `KeptnSLIProvider` interface, it's necessary to implement the
  `EvaluateQuery`, `EvaluateQueryForStep` and `FetchAnalysisValue` functions.
  You can follow other existing implementations,
 such as [prometheus.go](https://github.com/keptn/lifecycle-toolkit/blob/main/metrics-operator/controllers/common/providers/prometheus/prometheus.go),
 as an example.
  Below is an example of a dummy provider implementation.

    ```go
    {% include "./dummy-code-example.go" %}
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

5. **Add Test Cases:**

- Write a unit test to validate your implementation at the function level.
  Unit tests ensure that individual
 functions behave as expected and meet their functional requirements.
  Below is a unit test example for our dummy provider

  ```go
  {% include "./dummy-test-example.go" %}
  ```

- Include a KUTTL test to validate the behavior of Kubernetes resources managed by your code.
  KUTTL tests simulate real-world scenarios and interactions within a Kubernetes cluster, ensuring
  the correctness of your Kubernetes configurations and deployments.