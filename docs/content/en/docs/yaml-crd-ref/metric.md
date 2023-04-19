---
title: KeptnMetric
description: Define all workloads and checks associated with an application
weight: 50
---

A `KeptnMetric` is represents a metric that is collected
from a provider.
Providing the metrics as CRD into a Kubernetes cluster
facilitates the reusability of this data across multiple components
and llows using multiple observability platforms for different metrics.

The metric will be collected from the provider specified in the
specs.provider.name field.
The query is a string in the provider-specific query language, used to obtain a metric.

## Yaml Synopsis

```yaml
apiVersion: metrics.keptn.sh/v?alpha?
kind: KeptnMetric
metadata:
  name: <metric-name>
  namespace: <application-namespace>
spec:
  provider:
    name: "prometheus | dynatrace | dql"
  query: "<query-from-provider>"
  fetchIntervalSeconds: <seconds>
```

## Fields

* **apiVersion** -- API version being used.
`
* **kind** -- Resource type.  Must be set to `KeptnTaskDefinition`

* **metadata**
  * **name** -- Unique name of this metric.
    * Must be an alphanumeric string and, by convention, is all lowercase.
    * Can include the special characters `_`, `-`, (others?)
    * Should not include spaces.
  * **namespace** -- namespace of the application using this metric

## Usage


## Examples


```yaml
apiVersion: metrics.keptn.sh/v1alpha1
kind: KeptnMetric
metadata:
  name: keptnmetric-sample
  namespace: podtato-kubectl
spec:
  provider:
    name: "prometheus"
  query: "sum(kube_pod_container_resource_limits{resource='cpu'})"
  fetchIntervalSeconds: 5
```

### Example 1: inline script

### More examples

See the [operator/config/samples](https://github.com/keptn/lifecycle-toolkit/tree/main/operator/config/samples)
directory for more example `KeptnTaskDefinition` YAML files.
Separate examples are provided for each API version.
For example, the `lifecycle_v1alpha3_keptntaskdefinition` file
contains examples for the `v1alpha3` version of the lifecycle API group.

## Files

API Reference:

* [KeptnTaskDefinition](../../crd-ref/lifecycle/v1alpha3/#keptntaskdefinition)

## Differences between versions

The `KeptnTaskDefinition` is the same for
all `v1alpha?` library versions.

## See also

* Link to reference pages for any related CRDs
