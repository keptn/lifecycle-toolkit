---
comments: true
---

# KeptnEvaluationDefinition

A `KeptnEvaluationDefinition` assigns target values
to [KeptnMetric](metric.md) queries.
These are used as part of evaluation tasks
that Keptn runs
as part of pre- and post-analysis phases of a `KeptnApp` or workload.

## Yaml Synopsis

```yaml
apiVersion: lifecycle.keptn.sh/v1
kind: KeptnEvaluationDefinition
metadata:
  name: <evaluation-name>
spec:
  retries: <number-of-retries>
  retryInterval: <duration>
  objectives:
    - evaluationTarget: "<value>"
      keptnMetricRef:
        name: available-cpus
        namespace: some-namespace
```

## Fields

* **apiVersion** -- API version being used.
  Must be `v1beta1` or later for this syntax.
* **kind** -- Resource type.
  Must be set to `KeptnEvaluationDefinition`

* **metadata**
    * **name** -- Unique name of this evaluation
      such as `pre-deploy-eval` or `post-deploy-resource-eval`.
      Names must comply with the
      [Kubernetes Object Names and IDs](https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-subdomain-names)
      specification.

* **spec**

    * **objectives** (required) -- define the evaluations to be performed.
      Each objective is expressed as a `keptnMetricRef`
      and an `evaluationTarget` value.

        * **keptnMetricRef** (required) -- A reference to the [KeptnMetric](metric.md) object

            * **name** (required) -- Name of the referenced [KeptnMetric](metric.md) object

            * **namespace** -- Name of the referenced [KeptnMetric](metric.md) object

        * **evaluationTarget** (required) -- Desired value of the query,
          expressed as an arithmetic formula, usually less than (`<`) or greater than (`>`)
          This is used to define success or failure criteria for the referenced `KeptnMetric` in order to pass or fail
          the pre- and post-evaluation stages

    * **retries** -- specifies the number of times
      an `Keptnevaluation` defined by the `KeptnEvaluationDefinition`
      should be restarted if an attempt is unsuccessful.
      The default value is `10`.
    * **retryInterval** -- specifies the time
      to wait between the retries.
      The value supplied should specify the unit of measurement;
      for example, `5s` indicates 5 seconds and `5m` indicates 5 minutes.
      The default value is `5s`.

## Usage

A `KeptnEvaluationDefinition` references one or more [KeptnMetric](metric.md) resources.
When multiple `KeptnMetric`s are used, Keptn considers the evaluation successful
if **all** metrics meet their `evaluationTarget`.

The `KeptnMetric` resource and associated [KeptnMetricsProvider](metricsprovider.md)
resource must be located in the same namespace but the `KeptnEvaluationDefinition` resources
can reference metrics from any namespace in the cluster.
This means that you can create `KeptnMetricsProvider` and `KeptnMetric` resources
in a centralized namespace (e.g. in `keptn-system`) and use those metrics in evaluations
on all namespaces in the cluster.

## Example

```yaml
apiVersion: lifecycle.keptn.sh/v1
kind: KeptnEvaluationDefinition
metadata:
  name: my-prometheus-evaluation
  namespace: example
spec:
  retries: 5
  retryInterval: 10s
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

In the `v1alpha1` and `v1alpha2` API versions, `KeptnEvaluationDefinition` referenced the `KeptnEvaluationProvider` CR
to identify the data source associated with this definition and itself contained the queries
that are now taken from the specified [KeptnMetric](metric.md) CRD.
The synopsis was:

```yaml
apiVersion: lifecycle.keptn.sh/v1
kind: KeptnEvaluationDefinition
metadata:
  name: <evaluation-name>
spec:
  source: prometheus | dynatrace | datadog | thanos
  objectives:
    - name: query-1
      query: "xxxx"
      evaluationTarget: <20
    - name: query-2
      query: "yyyy"
      evaluationTarget: >4
```

Beginning with `v1beta1` API version, `KeptnEvaluationDefinition` references a `keptnMetricRef`
that points to a [KeptnMetric](metric.md) CR, that defines the data source, the query and the namespace to use.
The `KeptnEvaluationDefinition` merely specifies the evaluation target.

> **Note**
Versions `v1beta1` and `v1` are fully compatible.

## See also

* [KeptnMetricsProvider](metricsprovider.md)
* [KeptnMetric](metric.md)
