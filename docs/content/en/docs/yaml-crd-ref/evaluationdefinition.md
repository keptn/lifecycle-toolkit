---
title: KeptnEvaluationDefinition
description: Define an evaluation query

weight: 20
---

A `KeptnEvaluationDefinition` assigns target values
to [KeptnMetric](metric.md) queries.
These are used as part of evaluation tasks
that can be run by the Keptn Lifecycle Toolkit
as part of pre- and post-analysis phases of a workload or application.

## Yaml Synopsis

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha3
kind: KeptnEvaluationDefinition
metadata:
  name: <evaluation-name>
spec:
  objectives:
    - evaluationTarget: "<value>"
      keptnMetricRef:
        name: available-cpus
        namespace: some-namespace
```

## Fields

* **apiVersion** -- API version being used.
  Must be `v1alpha3` or later for this syntax.
* **kind** -- Resource type.
   Must be set to `KeptnEvaluationDefinition`

* **metadata**
  * **name** -- Unique name of this evaluation
    such as `pre-deploy-eval` or `post-deploy-resource-eval`.
    Names must comply with the
    [Kubernetes Object Names and IDs](https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-subdomain-names)
    specification.

* **spec**

  * **objectives** -- define the evaluations to be performed.
     Each objective is expressed as a `keptnMetricRef`
     and an `evaluationTarget` value.

    * **KeptnMetricRef** -- A reference to the
      [KeptnMetric](metric.md) object that contains the value,
      identified by `name` and `namespace`
    * **evaluationTarget** -- Desired value of the query,
       expressed as an arithmetic formula,
       usually less than (`<`) or greater than (`>`)
       This is used to define success or failure criteria
       for the referenced `KeptnMetric` in order to pass or fail
       the pre- and post-evaluation stages

## Usage

A `KeptnEvaluationDefinition` references one or more
[KeptnMetric](metric.md) resources.
When multiple `KeptnMetric`s are used,
the Keptn Lifecycle Toolkit considers the evaluation successful
if **all** metrics meet their `evaluationTarget`.

The `KeptnMetric` resource and associated
[KeptnMetricsProvider](metricsprovider.md)
resource must be located in the same namespace
but the `KeptnEvaluationDefinition` resources
can reference metrics from any namespace in the cluster.
This means that you can create `KeptnMetricsProvider`
and `KeptnMetric` resources
in a centralized namespace (e.g. in `keptn-lifecycle-toolkit-system`)
and use those metrics in evaluations
on all namespaces in the cluster.

## Example

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha3
kind: KeptnEvaluationDefinition
metadata:
  name: my-prometheus-evaluation
  namespace: example
spec:
  source: prometheus
  objectives:
    - keptnMetricRef:
        name: available-cpus
        namespace: example
      evaluationTarget: ">1"
    - keptnMetricRef:
        name: cpus-throttling
        namespace: example
      evaluationTarget: "<0.01"
```

## Files

API Reference:

## Differences between versions

In the `v1alpha1` and `v1alpha2` API versions,
`KeptnEvaluationDefinition` referenced the `KeptnEvaluationProvider` CR
to identify the data source associated with this definition
and itself contained the queries
that are now taken from the specified [KeptnMetric](metric.md) CRD.
The synopsis was:

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha2
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
`KeptnEvaluationDefinition` references a `keptnMetricRef`
that points to a [KeptnMetric](metric.md) CR,
that defines the data source, the query and the namespace to use.
The `KeptnEvaluationDefinition` merely specifies the evaluation target.

## See also

* [KeptnMetricsProvider](metricsprovider.md)
* [KeptnMetric](metric.md)
