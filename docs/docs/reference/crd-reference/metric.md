---
comments: true
---

# KeptnMetric

A `KeptnMetric` represents a metric that is collected from a provider.
Providing the metrics as a custom resource
facilitates the reusability of this data across multiple components
and allows using multiple observability platforms
for different metrics at the same time.

`KeptnMetric` resources are also used as targets for
[EvaluationDefinition](evaluationdefinition.md) resources.

## Yaml Synopsis

```yaml
apiVersion: metrics.keptn.sh/v1
kind: KeptnMetric
metadata:
  name: <metric-name>
  namespace: <application-namespace>
spec:
  provider:
    name: "<named-provider>"
  query: "<query>"
  fetchIntervalSeconds: <#-seconds>
  range:
    interval: "<timeframe>"
    step: <query-resolution-step-width>
    aggregation: p90 | p95 | p99 | max | min | avg | median
    storedResults: <integer>
  status:
    properties:
      value: <resulting value in human-readable language>
      rawValue: <resulting value, in raw format>
      errMsg: <error details if the query could not be evaluated>
      lastUpdated: <time when the status data was last updated>
```

## Fields

- **apiVersion** -- API version being used.

- **kind** -- Resource type.
  Must be set to `KeptnMetric`.

- **metadata**
    - **name** -- Unique name of this metric.
      Names must comply with the
      [Kubernetes Object Names and IDs](https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-subdomain-names)
      specification.
    - **namespace** -- Namespace of the application using this metric.

- **spec**
    - **provider.name** (required) --
      Name of this instance of the data source
      from which the metric is collected.
      This value must match the value of the `metadata.name` field
      of the corresponding [KeptnMetricsProvider](metricsprovider.md) CRD.

        Assigning your own name to the provider
        rather than just the type of provider
        enables you to support multiple instances of a data provider.
        For example, you might have `dev-prometheus`
        as the name of the Prometheus server that monitors the dev deployment
        and `prod-prometheus` as the name of the Prometheus server
        that monitors the production deployment.

    - **query** (required) -- String in the provider-specific query language,
      used to obtain a metric.

    - **fetchIntervalSeconds** (required) --
      Number of seconds between updates of the metric.
    - **range**
        - **interval** -- Timeframe for which the metric is queried.
          Defaults to 5m.
        - **step** -- A string that represents
          the query resolution step width for the data query
        - **aggregation** -- type of aggregation function
          to be applied to the data.
          Valid values are `p90`, `p95`, `p99`,
          `max`, `min`, `avg`, `median`.
        - **storedResults** -- Maximum number of past results
          to store in the status of a `KeptnMetric` resource.
          This can be set to an integer that is less than or equal to 255.
          When set to a value greater than 1,
          the user can see a slice of this number of metrics
          in the`status.intervalResults` field.

    - **status** --
      Keptn fills in this information when the metric is evaluated.
      It always records the time the metric was last evaluated.
      If the evaluation is successful,
      this stores the result in both human-readable and raw format.
      If the evaluation is not successful,
      this stores error details that you can use to understand the problem
      such as a forbidden code.

        By default, Keptn stores the most recent metric that was run.
        If the value of the `spec.range.storedResults` field
        is set to a value greater than 1 and no larger than 255,
        Keptn stores that number of metrics.

        - **value** -- A string that represents the resulting value
          in human-readable format.
        - **rawValue** -- An array that represents the resulting value
          in raw format.
        - **lastUpdated** -- Time when the status data was last updated.
        - **errMsg** -- Error details if the query could not be evaluated.
        - **intervalResults** -- Slice of all interval results.
          Up to 255 results can be stored,
          determined by the value of the `spec.range` field.

## Usage

As soon as you define and apply your `KeptnMetricsProvider` and `KeptnMetric` resources,
Keptn begins collecting the metrics you defined.
You do not need to do anything else.

A `KeptnMetric` resource must be located
in the same namespace as the associated
[KeptnMetricsProvider](metricsprovider.md)
resource.
`KeptnMetric` resources are used to generate metrics for the cluster
and are used as the SLI (Service Level Indicator) for
[KeptnEvaluationDefinition](evaluationdefinition.md)
resources that are used for pre- and post-deployment evaluations.

`KeptnEvaluationDefinition` resources can reference metrics
from any namespace.
This means that you can create `KeptnMetricsProvider`
and `KeptnMetric` resources
in a centralized namespace (e.g. in `keptn-system`)
and access those metrics in evaluations
on all namespaces in the cluster.

## Example

This example pulls metrics from the data provider
defined as `my-provider` in the `spec.provider.name` field
of the corresponding `KeptnMetricsProvider` CR.

```yaml
apiVersion: metrics.keptn.sh/v1
kind: KeptnMetric
metadata:
  name: keptnmetric-sample
  namespace: podtato-kubectl
spec:
  provider:
    name: "my-provider"
  query: "sum(kube_pod_container_resource_limits{resource='cpu'})"
  fetchIntervalSeconds: 5
  range:
    interval: "5m"
```

## Files

API Reference:

## Differences between versions

- Versions `v1beta1` and `v1` are fully compatible.

- Beginning with the `v1beta1` API version,
  the metrics controller supports multiple metrics in its `status` field
  if the value of the `spec.range.storedResults` field is greater than 1.

- Beginning with the `v1alpha3` API version,
  Keptn allows you to define multiple instances of the same data source.
  In earlier versions, you could use multiple data sources
  but only one instance of each.
  Consequently, the `v1alpha1` and `v1alpha2` API versions
  define the `provider` field with the type of the data provider
  (`prometheus`, `dynatrace`, or `dql`)
  rather than the particular name assigned
  to the instance of the data provider
  that is assigned in the
  [KeptnMetricsProvider](metricsprovider.md) CR.

    So the `v1alpha1` and `v1alpha2` synopsis
    of the `spec` field is:

    ```yaml
    ...
    spec:
      provider:
        name: "prometheus | dynatrace | dql"
      fetchIntervalSeconds: <seconds>
      query: >-
        "<query-from-provider>"
    ```

## See also

- [KeptnEvaluationDefinition](evaluationdefinition.md)
- [KeptnMetricsProvider](metricsprovider.md)
- Implementing [Keptn Metrics](../../guides/evaluatemetrics.md)
- [Getting started with Keptn metrics](../../getting-started/metrics.md)
- Architecture of the [Keptn Metrics Operator](../../components/metrics-operator.md)
