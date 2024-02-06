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
  To fully implement the `KeptnSLIProvider` interface, it's necessary to implement the following functions.
  `EvaluateQuery`(Fetches metric values from the dummy provider)
   - This function fetches metric values based on the provided
     metric query from the dummy provider.
     It evaluates the query and returns the metric values
     along with any additional data if required.
  `EvaluateQueryForStep`(Fetches metric values with step interval from the dummy provider)
   - This function fetches metric values with a specified step interval from the dummy provider.
      It takes into account the metric query and the step interval provided, executes the query,
      and returns the metric values along with any additional data if required.
  `FetchAnalysisValue`(Fetches analysis values from the dummy provider) functions.
   - This function fetches analysis values based on the provided query and time range from the
     dummy provider.
     It evaluates the query within the specified time range and returns the analysis
     values along with any additional data if required.
  You can follow other existing implementations,
 such as [prometheus.go](https://github.com/keptn/lifecycle-toolkit/blob/main/metrics-operator/controllers/common/providers/prometheus/prometheus.go),
 as an example.
  Below is an example of a dummy provider implementation.

    ```go
    {% include "./dummy-code-example.go" %}
    ```

   **NB:** Ensure to refer to the documentation of
    [metrics](https://github.com/keptn/lifecycle-toolkit/blob/main/docs/docs/reference/crd-reference/metric.md)
    and [analysis](https://github.com/keptn/lifecycle-toolkit/blob/main/docs/docs/reference/crd-reference/analysis.md)
    to understand what data should be retrieved from the methods inputs to compute accurate results.

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

5. **Update validation webhook and crd config:** To update the validation webhook and crd config of the metrics operator.

   Add the provider name next to last provider in this
   [line](https://github.com/keptn/lifecycle-toolkit/blob/main/metrics-operator/api/v1beta1/keptnmetricsprovider_types.go#L29)
   to look like this `// +kubebuilder:validation:Pattern:=prometheus|dynatrace|datadog|dql|dummy`.

   In the metric-operator directory run `make manifests` to update the metrics-operator crd config
  
6. **Add Test Cases:**

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

    ```yaml
    {% include "./00-install.yaml" %}
    ```

    Below is the assert check

    ```yaml
    {% include "./01-assert.yaml" %}
    ```
