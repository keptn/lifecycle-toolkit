---
comments: true
---

# AnalysisValueTemplate

An `AnalysisValueTemplate` resource
defines a Service Level Indicator (SLI),
which identifies the data to be analyzed
by a data source to use and the query to issue.
One Analysis can use data from multiple AnalysisValueTemplates.

## Synopsis

```yaml
apiVersion: metrics.keptn.sh/v1
kind: AnalysisValueTemplate
metadata:
  name: response-time-p95
  namespace: <namespace-where-this-resource-resides>
spec:
  provider:
    name: prometheus | thanos | dynatrace | dql | datadog
  query: <query>
```

## Fields

- **apiVersion** -- API version being used
- **kind** -- Resource type.
  Must be set to `AnalysisValueTemplate`
- **metadata**

    - **name** -- Unique name of this template.
       Names must comply with the
       [Kubernetes Object Names and IDs](https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-subdomain-names)
       specification.
    - **namespace** (optional) -- Namespace where this template lives.
       `Analysis` resources must specify this namespace
       when referencing this definition,
       unless it resides in the same namespace as the `Analysis` resource.

- **spec**
    - **provider** (required) -- the KeptnMetricProvider
        - **name** -- The `spec.name` value of the
            [KeptnMetricsProvider](metricsprovider.md) resource to use.
            Note that each `AnalysisValueTemplate` resource
            can use only one data source.
            However, an `Analysis` resource
            can use multiple `AnalysisValueTemplate` resources,
            each of which uses a different data source.
    - **query** (required) -- query to be made.
       This is done in the data provider's query language.
       It can include variables that use the go templating syntax
       to insert a placeholder in the query.
       For example, the query might include `{{.nodename}}'}`;
       the value to substitute for that variable for this Analysis
       is defined in the `spec.args` section of the `AnalysisTemplate` resource,
       which might be set to `nodename: test`.

## Usage

You must define a
[KeptnMetricsProvider](metricsprovider.md)
for each instance of each data provider you are using.
The `AnalysisValueTemplate` refers to that provider and queries it.

One `Analysis` can use data from multiple instances
of multiple types of data provider;
you must define a
[KeptnMetricsProvider](metricsprovider.md)
resource for each instance of each data provider you are using.
The template refers to that provider and queries it.

## Example

```yaml
{% include "../../assets/crd/analysis-template.yaml" %}
```

For a full example of how the `AnalysisValueTemplate` is used
to implement the Keptn Analysis feature, see the
[Analysis](../../guides/slo.md)
guide page.

## Files

API reference:
[AnalysisValueTemplate](../api-reference/metrics/v1/index.md#analysisvaluetemplate)

## Differences between versions

The Keptn Analysis feature is an official part of Keptn v0.10.0 and later.
Keptn v0.8.3 included a preliminary release of this feature
but it was hidden behind a feature flag.
The behavior of this feature is unchanged since v0.8.3.

## See also

- [Analysis](analysis.md)
- [AnalysisDefinition](analysisdefinition.md)
- [Analysis](../../guides/slo.md) guide
