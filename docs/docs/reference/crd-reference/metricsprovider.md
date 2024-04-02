---
comments: true
---

# KeptnMetricsProvider

A `KeptnMetricsProvider` resource defines an instance of a data provider
(such as Prometheus, Thanos, Dynatrace, or Datadog)
that is used by one or more [KeptnMetric](metric.md) resources.

One Keptn application can perform
traces, evaluations, and analyses based on metrics
from more than one data provider
and, beginning with the API version v1alpha3,
can use more than one instance of each data provider.
To implement this, create a `KeptnMetricsProvider` resource
for each instance of each data provider being used,
then reference the appropriate provider
for each metric definition by its name.

## Yaml Synopsis

```yaml
{% include "../../assets/crd/examples/yaml-synopsis.yaml"  %}
```

## Fields

* **apiVersion** -- API version being used.

* **kind** -- Resource type.
  Must be set to KeptnMetricsProvider

* **metadata**
    * **name** -- Unique name of this provider,
      used to reference the provider for the
      [KeptnEvaluationDefinition](evaluationdefinition.md)
      and [KeptnMetric](metric.md) resources.
      Names must comply with the
      [Kubernetes Object Names and IDs](https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-subdomain-names)
      specification.

        For example, you might define `dev-prometheus`
        for the Prometheus instance that monitors the development deployment,
        and `qa-prometheus` for the Prometheus instance
        that monitors the Prometheus instance that monitors the QA deployment,
        and `prod-dynatrace` for the Dynatrace instance
        that monitors the production deployment.

    * **namespace** -- Namespace where this provider is used.

* **spec**

    * **type** (required) -- The type of data provider for this instance
    * **targetServer** (required) -- URL of the data provider, enclosed in double quotes
    * **secretKeyRef**
        * **name:** -- Name of the Secret used by the provider
        * **key:** -- Key of the Secret from which to select
        * **optional** -- Specify whether the Secret or its key must be defined

## Usage

A KeptnMetricsProvider resource must be located
in the same namespace as the
[KeptnMetric](metric.md)
resources that are associated with it.
KeptnMetric resources are used to generate metrics for the cluster
and are used as the SLI (Service Level Indicator) for the
[KeptnEvaluationDefinition](evaluationdefinition.md)
and
[AnalysisValueTemplate](analysisvaluetemplate.md)
resources.

`KeptnEvaluationDefinition` resources can reference metrics
from any namespace.
This means that you can create `KeptnMetricsProvider`
and `KeptnMetric` resources
in a centralized namespace (e.g. in `keptn-system`)
and access those metrics in evaluations
on all namespaces in the cluster.

Metrics providers may require authentication data.
These can be stored in a Secret and referenced in the `KeptnMetricsProvider`.
Those Secret resources need to be configured differently for each metrics provider.
The Secret holding the data must have a specific structure
and contain the required data fields.
For detailed information please look at the [Examples section](#examples).

## Examples

<!-- markdownlint-disable MD046 -->

=== "Prometheus and Thanos"

    An example of Prometheus as a metrics provider with a Secret holding
    the authentication data looks like the following:

    ```yaml
    {% include "./assets/keptnmetricsprovider-prometheus.yaml" %}
    ```
    > **Note**
    Setting the `.spec.secretKeyRef.key` field in `KeptnMetricsProvider` is not necessary,
    as `user` and `password` key names are required to be present in the linked Secret.
    Setting this field won't have any effect.

=== "Datadog"

    An example of Datadog as a metrics provider with a Secret holding
    the authentication data looks like the following:
 
    ```yaml
    {% include "./assets/keptnmetricsprovider-datadog.yaml" %}
    ```
    > **Note**
    Setting the `.spec.secretKeyRef.key` field in `KeptnMetricsProvider` is not necessary,
    as `DD_CLIENT_API_KEY` and `DD_CLIENT_API_KEY` key names must be
    present in the linked Secret.
    Setting this field has no effect.

=== "Dynatrace and DQL"

    An example of Dynatrace as a metrics provider with a Secret holding
    the authentication data looks like the following:

    ```yaml
    {% include "./assets/keptnmetricsprovider-dynatrace.yaml" %}
    ```

    > **Note**
    When using Dynatrace as metrics provider you can
    define the key name of your DT token stored in a secret,
    which is not possible for Datadog, Prometheus or Thanos.
    For this example `myCustomTokenKey` was used.

<!-- markdownlint-enable MD046 -->

## Files

API Reference:

* [KeptnEvaluationDefinition](../api-reference/lifecycle/index.md)

## Differences between versions

For the `v1alpha2` API version,
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
{% include "../../assets/crd/examples/old-version-synopsis.yaml"  %}
```

Also note that, for the v1alpha1 and v1alpha2 API versions,
`KeptnMetricsProvider` only specifies the provider
for the `KeptnMetric` resource.
Beginning with `v1alpha3` API version,
`KeptnMetricsProvider` is also used to specify the provider
for the KeptnEvaluationDefinition resource.

## See also

* [KeptnEvaluationDefinition](evaluationdefinition.md)
* [AnalysisValueTemplate](analysisvaluetemplate.md)
* [KeptnMetric](metric.md)
