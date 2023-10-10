---
title: AnalysisValueTemplate
description: Define the data source and query for each SLI
weight: 8
hide: true
---

An `AnalysisValueTemplate` resource
defines a Service Level Indicator (SLI),
which identifies the data to be analyzed
by data source to use and query to issue.
One Analysis can use data from multiple instances
of multiple types of data provider.

## Synopsis

```yaml
apiVersion: metrics.keptn.sh/v1alpha3
kind: AnalysisValueTemplate
metadata:
  labels:
    app.kubernetes.io/name: <name-of-application>
    app.kubernetes.io/instance: <name-of-instance>
    app.kubernetes.io/part-of: <name-of-higher-level-app>
    app.kubernetes.io/managed-by: helm | kustomize
    app.kubernetes.io/created-by: metrics-operator

    TODO: `created-by` is not in the list of k8s Recommended Labels

    TODO: We have the same five annotations in the `Analysis` CRD.
    Why?
    And do the values in `AnalysisValueTemplate` and `Analysis`
    need to be identical or could they be different?

  name: response-time-p95
  namespace: keptn-lifecycle-toolkit-system
spec:
  provider:
    name: prometheus
  query: <query>
```

## Fields

* **apiVersion** -- API version being used
* **kind** -- Resource type.
  Must be set to `AnalysisValueTemplate`
* **metadata**
  * **labels** -- The Analysis feature uses the
    `name` and `part-of` labels that are discussed in
    [Basic annotations](../implementing/integrate/#basic-annotations)
    plus the following:
    * **app.kubernetes.io/instance** analysis-sample
    * **app.kuberentes.io/managed-by** -- Tool used to manage
      the operation of the application.
      Valid values are `helm` and `kustomize`.
    * **app.kubernetes.io/created-by** metrics-operator

      TODO: Need to clarify how to use these annotations
  * **name** -- Unique name of this template.
    Names must comply with the
    [Kubernetes Object Names and IDs](https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-subdomain-names)
    specification.
  * **namespace** -- Namespace where this template lives
* **spec**
  * **provider**
    * **name** -- The `spec.name` value of the
      [KeptnMetricsProvider](metricsprovider.md) resource to use.
  * **query** -- query to be made.
    This is done in the data provider's query language.
    It can include variables in the form `{{.nodename}}'}`;
    The value to substitute for that variable for this Analysis
    is defined in the `spec.args` section of the `AnalysisTemplate` resource.

## Usage

You must define a
[KeptnMetricsProvider](metricsprovider.md)
for each instance of each data provider you are using.
The `AnalysisValueTemplate` refers to that provider and queries it.

One `Analysis` can use data from multiple instances
of multiple types of data provider;
you must define a
[KeptnMetricsProvider](../../yaml-crd-ref/metricsprovider.md)
resource for each instance of each data provider you are using.
The template refers to that provider and queries it.

## Examples

{{< embed path="/metrics-operator/config/samples/metrics_v1alpha3_analysisvaluetemplate.yaml" >}}

For a full example of how the `AnalysisValueTemplate` is used
to implement the Keptn Analysis feature, see the
[Analysis](../implementing/slo)
guide page.

## Files

[AnalysisValueTemplate](../crd-ref/metrics/v1alpha3/#analysisvaluetemplate)
API reference

## Differences between versions

A preliminary release of the Keptn Analysis feature
but is hidden behind a feature flag.
To preview these features, set the environment `ENABLE_ANALYSIS` to `true`
in the `metrics-operator` deployment.

## See also

* [Analysis](analysis.md)
* [AnalysisDefinition](analysisdefinition.md)
* [Analysis](../implementing/slo) guide
