---
title: Analysis
description: Define specific configurations and the Analysis to report
weight: 4
hide: true
---

An `Analysis` resource customizes the templates
that are defined in an
[AnalysisDefinition](analysisdefinition) resource
by identifying the time for which the analysis should be done
and the appropriate values to use for variables
that are used in the `AnalysisDefinition` query.

## Synopsis

```yaml
apiVersion: metrics.keptn.sh/v1alpha3
kind: Analysis
metadata:
  labels:
    app.kubernetes.io/name: <name-of-generated-`KeptnWorkload` resource>
    app.kubernetes.io/instance: analysis-sample
    app.kubernetes.io/part-of: <name-of-associated-KeptnApp-resource>
    app.kuberentes.io/managed-by: kustomize | ??
    app.kubernetes.io/created-by: metrics-operator
  name: analysis-sample
spec:
  timeframe: from: <start-time> to: <end-time> | `recent`
  args:
    <variable1>: <value1>
    <variable2>: <value2>
    ...
  analysisDefinition:
    name: <name of associated `analysisDefinition` resource
    namespace: <namespace of associated `analysisDefinition` resource
```

## Fields

* **apiVersion** -- API version being used
* **kind** -- Resource type.
   Must be set to `Analysis`
* **metadata**
  * **labels** -- The Analysis feature uses the
    `name` and `part-of` labels that are discussed in
    [Basic annotations](../implementing/integrate/#basic-annotations)
    plus the following:
    * **app.kubernetes.io/instance** analysis-sample
    * **app.kuberentes.io/managed-by** kustomize
    * **app.kubernetes.io/created-by** metrics-operator

    TODO: In your example, `name` is set to analysis.
    Does this mean that your workload is called `analysis`
    or is this label here meant to refer to the analysis resource or feature??
    Maybe better to use `my-analysis` for the example?

    TODO: Need descriptions of `instance`, `managed-by` and `created-by`
    as well as valid values.
    Especially alternatives to `kustomize` for `managed-by`

    TODO: Are there Keptn equivalents or annotations (Keptn or k8s)?
    If not, can you mix Keptn and k8s labels/annotations?

  * **name** -- Unique name of this analysis.
    Names must comply with the
    [Kubernetes Object Names and IDs](https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-subdomain-names)
    specification.
* **spec**
  * **timeframe** -- Specifies the range  for the corresponding query
    in the AnalysisValueTemplate.
    This can be populated as one of the following:

    * A combination of ‘from’ and ’to’
      to specify the start and stop times for the analysis
    * Set the ‘recent’ property.

    If neither is set, the Analysis can not be added to the cluster.
  * **args** -- Map of key/value pairs that can be used
    to substitute variables in the `AnalysisValueTemplate` query.
  * **analysisDefinition** -- Identify the `AnalysisDefinition` resource
    that stores the `AnalysisValuesTemplate` associated with this `Analysis`
    * **name: Name of the `AnalysisDefinition` resource
    * **namespace: Namespace of the `AnalysisDefinition` resource.

## Usage

## Examples

{{< embed path="/metrics-operator/config/samples/metrics_v1alpha3_analysis.yaml" >}}

This `Analysis` resource:

* Defines the `timeframe` for which the analysis is done
  as between 5 am and 10 am on the 5th of May 2023
* Adds a few specific key-value pairs that will be substituted in the query.
  For instance, the query could contain the `{{.nodename}}` variable.
  The value of the `args.nodename` field (`test`)
  will be substituted for this string.

For a full example of how to implement the Keptn Analysis feature, see the
[Analysis](../implementing/slo)
guide page.

## Files

API reference: [Analysis](../../crd-ref/metrics/v1alpha3/#analysis)

## Differences between versions

The Analysis feature was first introduced in Keptn v.0.9.0.

## See also

* [AnalysisDefinition](analysisdefinition.md)
* [AnalysisValueTemplate](analysisvaluetemplate.md)
* [Analysis](../implementing/slo) guide
