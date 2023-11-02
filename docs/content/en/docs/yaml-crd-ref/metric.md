---
title: KeptnMetric
description: Define all workloads and checks associated with an application
weight: 50
---

A `KeptnMetric` represents a metric that is collected from a provider.
Providing the metrics as a CR in a Kubernetes cluster
facilitates the reusability of this data across multiple components
and allows using multiple observability platforms
for different metrics at the same time.

`KeptnMetric` CRs are also used as targets for
[EvaluationDefinition](evaluationdefinition.md) CRs.

## Yaml Synopsis

```yaml
apiVersion: metrics.keptn.sh/v1alpha3
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
  status:
    properties:
      value: <resulting value in human-readable language>
      rawValue: <resulting value, in raw format>
      errMsg: <error details if the query could not be evaluated>
      lastUpdated: <time when the status data was last updated>
```

## Fields

* **apiVersion** -- API version being used.

* **kind** -- Resource type.
   Must be set to `KeptnMetric`.

* **metadata**
  * **name** -- Unique name of this metric.
    Names must comply with the
    [Kubernetes Object Names and IDs](https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-subdomain-names)
    specification.
  * **namespace** -- Namespace of the application using this metric.

* **spec**
  * **provider.name** (required) --
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
  * **query** (required) -- String in the provider-specific query language,
    used to obtain a metric.

  * **fetchIntervalSeconds** (required) -- Number of seconds between updates of the metric.
  * **range**
    * **interval** -- Timeframe for which the metric is queried.
    Defaults to 5m.

* **status**
  * Keptn fills in this information when the metric is evaluated.
    It always records the time the metric was last evaluated.
    If the evaluation is successful,
    this stores the result in both human-readable and raw format.
    If the evaluation is not successful,
    this stores error details that you can use to understand the problem
    such as a forbidden code.

## Usage

As soon as you define and apply
your `KeptnMetricsProvider` and `KeptnMetric` resources,
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
in a centralized namespace (e.g. in `keptn-lifecycle-toolkit-system`)
and access those metrics in evaluations
on all namespaces in the cluster.

## Example

This example pulls metrics from the data provider
defined as `my-provider` in the `spec.provider.name` field
of the corresponding `KeptnMetricsProvider` CR.

```yaml
apiVersion: metrics.keptn.sh/v1alpha3
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

Beginning with the `v1alpha3` API version,
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

* [KeptnEvaluationDefinition](evaluationdefinition.md)
* [KeptnMetricsProvider](metricsprovider.md)
* Implementing [Keptn Metrics](../implementing/evaluatemetrics.md)
* [Getting started with Keptn metrics](../getting-started/metrics.md)
* Architecture of the [Keptn Metrics Operator](../architecture/components/metrics-operator.md)
