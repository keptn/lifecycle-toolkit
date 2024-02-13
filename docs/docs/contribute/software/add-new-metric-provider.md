---
comments: true
---

# Adding a New Metrics Provider for Placeholder Endpoint

In this guide, we will create a placeholder provider.
The following steps are a starting point to create your own custom provider:

1. Fork and clone the [Keptn repo](https://github.com/keptn/lifecycle-toolkit)
  for more information [checkout this link](https://keptn.sh/stable/docs/contribute/general/git/fork-clone/)

2. **Define the Provider Type:** In the `metrics-operator/controllers/common/providers/common.go` file,
 define the constant `KeptnPlaceholderProviderType`.
  In our example we use `"placeholder"`.

    ```go
    const KeptnPlaceholderProviderType = "placeholder"
    ```

3. **Implement the Provider:** Create your own new folder inside the
[metrics-operator/controllers/common/providers]
(<https://github.com/keptn/lifecycle-toolkit/tree/main/metrics-operator/controllers/common/providers>),
matching the new provider name (`placeholder` in our example).
Create a new Go package for the placeholder provider in that folder.
This package should contain a `struct` that implements the `KeptnSLIProvider` interface.
To fully implement the `KeptnSLIProvider` interface, it's necessary to implement the following functions.

  ```markdown
  <!-- markdownlint-disable MD007 -->
  * `EvaluateQuery`(Fetches metric values from the provider)
    * This function fetches metric values based on the provided
      metric query from the provider.
      It evaluates the query and returns the metric values
      along with any additional data if required.
    * It takes as input a [KeptnMetric]
      (<https://github.com/keptn/lifecycle-toolkit/blob/main/metrics-operator/api/v1beta1/keptnmetric_types.go>)
      and [KeptnMetricsProvider](https://github.com/keptn/lifecycle-toolkit/blob/main/metrics-operator/api/v1beta1/keptnmetricsprovider_types.go)
  * `EvaluateQueryForStep`(Fetches metric values with step interval from the provider)
    * This function fetches metric values with a specified step interval from the placeholder provider.
      It takes into account the metric query and the step interval provided, executes the query,
      and returns the metric values along with any additional data if required.
    * It takes as input a [KeptnMetric]
      (<https://github.com/keptn/lifecycle-toolkit/blob/main/metrics-operator/api/v1beta1/keptnmetric_types.go>)
      and [KeptnMetricsProvider](https://github.com/keptn/lifecycle-toolkit/blob/main/metrics-operator/api/v1beta1/keptnmetricsprovider_types.go)
  * `FetchAnalysisValue`(Fetches analysis values from the provider) functions.
    * This function fetches analysis values based on the provided query and time range from the
      provider.
      It evaluates the query within the specified time range and returns the analysis
      values along with any additional data if required.
    * It takes as input a `query`, [Analysis]
      (<https://github.com/keptn/lifecycle-toolkit/blob/main/metrics-operator/api/v1beta1/analysis_types.go>) and [KeptnMetricsProvider](https://github.com/keptn/lifecycle-toolkit/blob/main/metrics-operator/api/v1beta1/keptnmetricsprovider_types.go)
    <!-- markdownlint-enable MD007 -->
    ```

  You can follow other existing implementations,
 such as [prometheus.go](https://github.com/keptn/lifecycle-toolkit/blob/main/metrics-operator/controllers/common/providers/prometheus/prometheus.go),
 as an example.
   **NB:** Each of the three functions expects a string containing a float value in it.
  But for example purposes
           we returned some of the data accessible in the function.
  Below is an example of a placeholder provider implementation.

    ```go
    {% include "./assets/example-code/placeholder-code-example.go" %}
    ```

   **NB:** Ensure to refer to the documentation of
    [metrics](https://github.com/keptn/lifecycle-toolkit/blob/main/docs/docs/reference/crd-reference/metric.md)
    and [analysis](https://github.com/keptn/lifecycle-toolkit/blob/main/docs/docs/reference/crd-reference/analysis.md)
    to understand what data should be retrieved from the methods inputs to compute accurate results.

4. **Instantiate the Provider** in the `providers.NewProvider` function
 in the `metrics-operator/controllers/common/providers/provider.go` file.
 add a case for the `KeptnPlaceholderProviderType`.
  Instantiate the placeholder provider struct and return it.

    ```go
    // Inside the providers package

    // NewProvider function
    func NewProvider(providerType string, log logr.Logger, k8sClient client.Client) (KeptnSLIProvider, error) {
        switch strings.ToLower(providerType) {
        case KeptnPlaceholderProviderType:
            return &placeholder.KeptnPlaceholderProvider{
                Log:        log,
                HttpClient: http.Client{},
            }, nil
        // Other cases...
        }
    }
    ```

5. **Update the validation webhook and crd config:** To update the validation webhook and crd config of the metrics operator.

   Add the provider name next to last providers on this
   [line](https://github.com/keptn/lifecycle-toolkit/blob/main/metrics-operator/api/v1beta1/keptnmetricsprovider_types.go#L29)
   to look like this `// +kubebuilder:validation:Pattern:=prometheus|dynatrace|datadog|dql|placeholder`.

   In the metric-operator directory run `make manifests` to update the metrics-operator crd config

   Then modify the helm chart and the helm chart crd validation to match the update in the metrics-operator crd config
  
6. **Add Test Cases:**

  * Write a unit test to validate your implementation at the function level.
    Unit tests ensure that individual
    functions behave as expected and meet their functional requirements.

  * Include a Chainsaw test to validate the behavior of Kubernetes resources managed by your code.
    Chainsaw tests simulate real-world scenarios and interactions within a Kubernetes cluster, ensuring
    the correctness of your Kubernetes configurations and deployments.
