---
title: KeptnEvaluationDefinition
description: Define all workloads and checks associated with an application
weight: 20
---

A `KeptnEvaluationDefinition` defines evaluation tasks
that can be run by the Keptn Lifecycle Toolkit
as part of pre- and post-analysis phases of a workload or application.

## Yaml Synopsis

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha3
kind: KeptnEvaluationDefinition
metadata:
  name: pre-deployment-hello
spec:
  objectives:
    - evaluationTarget: ">1"
      keptnMetricRef:
        name: available-cpus
        namespace: some-namespace
```

## Fields

* **apiVersion** -- API version being used.
`
* **kind** -- Resource type.
   Must be set to `KeptnEvaluationDefinition`

* **metadata**
  * **name** -- Unique name of this evaluation
    such as `pre-deploy-eval` or `post-deploy-resource-eval`.
    * Must be an alphanumeric string and, by convention, is all lowercase.
    * Can include the special characters `_`, `-`, (others?)
    * Should not include spaces.

* **spec**
  * **source** -- Name of the data provider being used for this evaluation.
    The value of the `source` field must match
    the string used for the `name` field
    in the corresponding [KeptnEvaluationProvider](evaluationprovider.md) CRD.

    Each `KeptnEvaluationDefinition` CRD can use only one data provider;
    if you are using multiple data provider, you must create
    `KeptnEvaluationProvider` and `KeptnEvaluationDefinition` CRDs for each.

    Currently, you can only access one occurrance of each type of data provider
    in your KLT cluster.

  * **objectives** -- define the evaluations to be performed.
     Each objective is expressed as a `query` and an `evaluationTarget` value.

    * **query** -- Any query that is supported by the data provider.
    * **value** -- Desired value of the query,
       expressed as an arithmatic formula,
       usually less than (`<`) or greater than (`>`)

## Usage

## Examples

## Files

API Reference:

## Differences between versions

In the `v1alpha1` and `v1alpha2` API versions,
`KeptnEvaluationDefinition` references the `KeptnEvaluationProvider` CRD
to identify the data source associated with this definition
and itself contained the queries
that are now taken from the specified [KeptnMetric](metric.md) CRD.
The synopsis was:

```yaml
apiVersion: lifecycle.keptn.sh/v?alpha?
kind: KeptnEvaluationDefinition
metadata:
  name: <evaluation-name>
spec:
  source: prometheus | dynatrace | datadog
  objectives:
    - name: query-1
      query: "xxxx"
      evaluationTarget: <20
    - name: query-2
      query: "yyyy"
      evaluationTarget: >4
```

Beginning with `v1alpha3` API version,
`KeptnEvaluationDefinition` references the data source defined
in the [KeptnMetricsProvider](metricsprovider.md)
and the queries are specified in the corresponding
[KeptnMetric](metric.md) CRD
although the `evaluationTarget` is defined in this CRD.

## See also

* [KeptnMetricsProvider](metricsprovider.md)
* [KeptnMetric](metric.md)
