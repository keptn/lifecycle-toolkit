---
title: KeptnMetric
description: Define all workloads and checks associated with an application
weight: 50
---

A `KeptnMetric` represents a metric that is collected from a provider.
Providing the metrics as a CRD into a Kubernetes cluster
facilitates the reusability of this data across multiple components
and allows using multiple observability platforms
for different metrics at the same time.

A `KeptnMetric` looks like the following:

## Yaml Synopsis

```yaml
apiVersion: metrics.keptn.sh/v?alpha?
kind: KeptnMetric
metadata:
  name: <metric-name>
  namespace: <application-namespace>
spec:
  provider:
    name: "<named-provider>"
  query: "<query>"
  fetchIntervalSeconds: <#-seconds>
```

## Fields

* **apiVersion** -- API version being used.
`
* **kind** -- Resource type.
   Must be set to `KeptnTaskDefinition`

* **metadata**
  * **name** -- Unique name of this metric.
    * Must be an alphanumeric string and, by convention, is all lowercase.
    * Can include the special characters `_`, `-`, (others?)
    * Should not include spaces.
  * **namespace** -- namespace of the application using this metric

* **spec**
  * **provider.name** --
    Name of this instance of the data source
    from which the metric is collected.
    This value must match the value of the `spec.provider.name` field
    of the corresponding [KeptnMetricsProvider](metricsprovider.md) CRD
    Assigning your own name to the provider
    rather than just the type of provider
    enables you to support multiple instances of a data provider.
    For example, you might have `dev-prometheus`
    as the name of the Prometheus server that monitors the dev deployment
    and `prod-prometheus` as the name of the Prometheus server
    that monitors the production deployment
  * **query** -- String in the provider-specific query language,
    used to obtain a metric.
  * **fetchIntervalSeconds** -- Number of seconds between ??

## Usage

## Example

This example pulls metrics from the data provider
defined as `my-provider` in the `spec.provider.name` field
of the corresponding `KeptnMetricsProvider` CRD.

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


## Files

API Reference:

## Differences between versions

Beginning with the `v1alpha3` API version,
Keptn allows you to define multiple instances of the same data source.
In earlier versions, you could use multiple data sources
but only one instance of each.
Consequently the `v1alpha1` and `v1alpha2` library versions
define the `provider` field with the name of the data provider
(`prometheus`, `dynatrace`, or `dql`)
rather than the particular name assigned
to the instance of the data provider
that is assigned in the
[KeptnMetricsProvider](metricsprovider.md) CRD.

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
* Implementing [Keptn Metrics](../implementing/metrics.md)
* Architecture of the [Keptn Metrics Operator](../concepts/architecture/components/metrics-operator/_index.md)
