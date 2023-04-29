---
title: KeptnMetricsProvider
description: Define data provider used for metrics and evaluations
weight: 55
---

`KeptnMetricsProvider` defines an instance of the data provider
(such as Prometheus, Dynatrace, or Datadog)
that is used by the [KeptnMetric](metric.md)
and [KeptnEvaluationDefinition](evaluationdefinition.md) CRDs.
One Keptn application can perform evaluations and metrics
from more than one data provider
and, beginning in V0.8.0,
more than one instance of each data provider.
To implement this, create a `KeptnMetricsProvider` CRD
for each instance of each data provider being used
then reference the appropriate provider
for each evaluation or metric definition.

## Yaml Synopsis

```yaml
apiVersion: lifecycle.keptn.sh/v?alpha?
kind: KeptnMetricsProvider
metadata:
  name: <data-source-instance-name>
  namespace: <namespace>
spec:
  type: prometheus | dynatrace | dql
  targetServer: "<data-source-url>"
  secretKeyRef:
    name: <token>
    key: <TOKEN>


## Fields

* **apiVersion** -- API version being used.
`
* **kind** -- Resource type.  Must be set to `KeptnMetricsProvider`

* **metadata**
  * **name** -- Unique name of this provider,
    used to reference the provider for the
    [KeptnEvaluationDefinition](evaluationdefinition)
    and [KeptnMetric](metric.md) CRs.
    * Must be an alphanumeric string and, by convention, is all lowercase.
    * Can include the special characters `_`, `-`, (others?)
    * Should not include spaces.

  * **namespace** -- Namespace where this provider is used.

* **spec**

  * **type** -- The type of data provider for this instance
  * **targetServer** -- URL of the data provider, enclosed in double quotes
  * **secretKeyRef**
    * **name:** -- Name of the token for this data provider
    * **key:** -- Key for this data provider


## Usage


## Examples

### Example 1: Dynatrace data provider

```yaml
apiVersion: metrics.keptn.sh/v1alpha2
kind: KeptnMetricsProvider
metadata:
  name: dynatrace
  namespace: podtato-kubectl
spec:
  targetServer: "<dynatrace-tenant-url>"
  secretKeyRef:
    name: dt-api-token
    key: DT_TOKEN
```

## Files

API Reference:

* [KeptnTaskDefinition](../../crd-ref/lifecycle/v1alpha3/#keptntaskdefinition)

## Differences between versions

For the `v1alpha1` and `v1alpha2` API versions,
Keptn did not support
using more than one instance of a particular data provider
in the same namespace.
In other words, one namespace could support one instance each
of Prometheus, Dynatrace, and Datadog
but could not support, for example, two instances of Prometheus.

The synopsis in those older API versions
only specified the `metadata.name` field
that identified the data provider (`prometheus`, `dynatrace`, or `dql`):

```yaml
apiVersion: metrics.keptn.sh/v1alpha2
kind: KeptnMetricsProvider
metadata:
  name: prometheus | dynatrace |dql
  namespace: <namespace>
spec:
  targetServer: "<data-provider-url>"
  secretKeyRef:
    name: dt-api-token
    key: DT_TOKEN
```

Also note that, for the v1alpha1 and v1alpha2 API versions,
`KeptnMetricsProvider` only specifies the provider
for the `KeptnMetrics` CR.
Beginning with `v1alpha3` API version,
`KeptnMetricsProvider` is also used to specify the provider
for the `KeptnEvaluationDefinition` CR.

## See also

* [KeptnEvaluationDefinition](evaluationdefinition.md)
* [KeptnMetric](metric.md)
