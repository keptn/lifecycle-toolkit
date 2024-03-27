---
comments: true
---

# Add a metrics provider

The
[KeptnMetric](../../guides/evaluatemetrics.md)
feature works with almost any data platform
but Keptn requires that a metrics provider be defined
for any data platform it uses as a data source.
This guide gives instructions for creating a new metrics provider.
For these instructions,
we define a sample provider called `placeholder`.
The steps to create your own metrics provider are:

1. **Fork and clone** the
[Keptn repository](https://github.com/keptn/lifecycle-toolkit).
 For more information, see
 [Fork and clone the repository](https://keptn.sh/stable/docs/contribute/general/git/fork-clone/).

2. **Define the Provider Type:** In the `metrics-operator/controllers/common/providers/common.go` file,
 define the constant `KeptnPlaceholderProviderType`.
  In our example we use `"placeholder"`.

    ```go
    const KeptnPlaceholderProviderType = "placeholder"
    ```

3. **Implement the Provider:** Create a new folder inside the
   [metrics-operator/controllers/common/providers](<https://github.com/keptn/lifecycle-toolkit/tree/main/metrics-operator/controllers/common/providers>).
   Use the provider name as the name of the folder.
   This name defines the string used to identify this provider
   in the `spec.type` field of the
   [KeptnMetricsProvider](../../reference/crd-reference/metricsprovider.md)
   resource.
   In this example, the folder is named `placeholder`.
   Create a new Go package for the placeholder provider in that folder.
   This package should contain a `struct` that implements the `KeptnSLIProvider` interface.
   To fully implement the `KeptnSLIProvider` interface, it's necessary to implement the following functions.
  
      * `EvaluateQuery`(Fetches metric values from the provider)
        * This function fetches metric values based on the provided
          metric query from the provider.
          It evaluates the query and returns the metric values
          along with any additional data if required.
        * It takes as input a [KeptnMetric](../../reference/crd-reference/metric.md)
          and [KeptnMetricsProvider](../../reference/crd-reference/metricsprovider.md)
      * `EvaluateQueryForStep`(Fetches metric values with step interval from the provider)
        * This function fetches metric values with a specified step interval from the placeholder provider.
          It takes into account the metric query and the step interval provided, executes the query,
          and returns the metric values along with any additional data if required.
        * It takes as input a [KeptnMetric](../../reference/crd-reference/metric.md)
          and [KeptnMetricsProvider](../../reference/crd-reference/metricsprovider.md)
      * `FetchAnalysisValue`(Fetches analysis values from the provider) functions.
        * This function fetches analysis values based on the provided query and time range from the
          provider.
          It evaluates the query within the specified time range and returns the analysis
          values along with any additional data if required.
        * It takes as input an [Analysis](../../reference/crd-reference/analysis.md),
          resource that contains a `query` and a
          [KeptnMetricsProvider](../../reference/crd-reference/metricsprovider.md) resource.

     You can follow other existing implementations,
     such as [prometheus.go](https://github.com/keptn/lifecycle-toolkit/blob/main/metrics-operator/controllers/common/providers/prometheus/prometheus.go),
     as an example.

     Each of the three functions expects a string containing a float value in it.
     But for example purposes we returned some of the data accessible in the function.

     Below is an example of a placeholder provider implementation.

      ```go
      {% include "./assets/example-code/placeholder-code-example.go" %}
      ```

      > **Note** Refer to the documentation of the
      > [KeptnMetric](../../reference/crd-reference/metric.md)
      > and
      > [Analysis](../../reference/crd-reference/analysis.md)
      > resources
      > to understand what data should be retrieved from the methods inputs to compute accurate results.

4. **Instantiate the Provider** in the `providers.NewProvider` function
   in the `metrics-operator/controllers/common/providers/provider.go` file.
   add a case for the `KeptnPlaceholderProviderType`.
   Instantiate the placeholder provider struct and return it.

    ```go
    {% include "./assets/example-code/new-provider-function.go" %}
    ```

5. **Update the validation webhook and crd config:** To update the validation webhook and crd config of the metrics operator.
   Add the provider name next to last providers on this
   [line](https://github.com/keptn/lifecycle-toolkit/blob/main/metrics-operator/api/v1/keptnmetricsprovider_types.go#L29)
   to look like this

    `// +kubebuilder:validation:Pattern:=prometheus|thanos|dynatrace|datadog|dql|placeholder`.

     In the metric-operator directory run `make generate manifests` to update the metrics-operator crd config
     Then modify the helm chart and the helm chart crd validation to match the update in the metrics-operator crd config
  
6. **Add Test Cases:**
     * Write a unit test to validate your implementation at the function level.
       Unit tests ensure that individual
       functions behave as expected and meet their functional requirements.

     * Include a Chainsaw test to validate the behavior of Kubernetes resources managed by your code.
       Chainsaw tests simulate real-world scenarios and interactions within a Kubernetes cluster, ensuring
       the correctness of your Kubernetes configurations and deployments.

        Below are the steps for adding an integration test.

        * In the directory `test/chainsaw/testmetrics`, create a folder `metrics-provider-placeholder` in our case.
        * Within the `keptn-metrics-placeholder` folder, create YAML file `00-install.yaml`.
          * `00-install.yaml` contains a sample configuration that installs a valid `KeptnMetricsProvider`
              in our case `placeholder` and it also defines a sample `KeptnMetric` configuration
              representing a valid use case, while.
        * Create a file named `chainsaw-test.yaml` and define the steps for the integration test in chainsaw-test.yaml.

        > For more information checkout [an already existing integration test](https://github.com/keptn/lifecycle-toolkit/tree/main/test/chainsaw/testmetrics/metrics-provider)
