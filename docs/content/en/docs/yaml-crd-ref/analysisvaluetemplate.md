---
title: AnalysisValueTemplate
description: Define the data source and query for each SLI
weight: 8
hide: true
---

An `AnalysisValueTemplate` resource
identifies the data source and query that define an SLI.
One Analysis can use data from multiple instances
of multiple types of data provider.

## Synopsis

```yaml
apiVersion: metrics.keptn.sh/v1alpha3
kind: AnalysisValueTemplate
metadata:
  labels:
    app.kubernetes.io/name: analysisvaluetemplate
    app.kubernetes.io/instance: analysisvaluetemplate-sample
    app.kubernetes.io/part-of: metrics-operator
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/created-by: metrics-operator
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
  * **name** -- Unique name of this template.
    Names must comply with the
    [Kubernetes Object Names and IDs](https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-subdomain-names)
    specification.
  * **namespace** -- Namespace where this templante lives
* **spec**
  * **provider**
    * **name** prometheus -- The `spec.name` value of the
      [KeptnMetricsProvider](metricsprovider.md) resource to use.

    TODO: Can this template use multiple data providers?
  * **query** -- query to be made.
    This is done in the data provider's query language.
    It can include variables in the form `{{.nodename}}'}`;
    The value to substitute for that variable for this Analysis
    is defined in the `spec.args` section of the `AnalysisTemplate` resource.

## Usage

You must define a
[KeptnMetricsProvider](metricsprovider.md]
for each instance of each data provider you are using.
The `AnalysisValueTemplate` refers to that provider and queries it.

## Examples

```yaml
apiVersion: metrics.keptn.sh/v1alpha3
kind: AnalysisValueTemplate
metadata:
  labels:
    app.kubernetes.io/name: analysisvaluetemplate
    app.kubernetes.io/instance: analysisvaluetemplate-sample
    app.kubernetes.io/part-of: metrics-operator
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/created-by: metrics-operator
  name: response-time-p95
  namespace: keptn-lifecycle-toolkit-system
spec:
  provider:
    name: prometheus
  query: <query>
```

For a full example of how the `AnalysisValueTemplate` is used
to implement the Keptn Analysis feature, see the
[Analysis](../implementing/slo)
guide page.

## Files

[AnalysisValueTemplate](../../crd-ref/metrics/v1alpha3/#analysisvaluetemplate)
API reference

## Differences between versions

The Analysis feature was first introduced in Keptn v.0.9.0.

## See also

* [Analysis](analysis.md)
* [AnalysisDefinition](analysisdefinition.md)
* [Analysis](../implementing/slo) guide
