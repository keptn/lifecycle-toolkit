---
title: Analysis
description: How to implement Keptn Analyses
weight: 60
---

The Keptn Metrics Operator Analysis feature
allows you to validate a deployment or release
using data from the observability data provider(s)
that are configured for Keptn Metrics.
You define the quality criteria for the analysis with SLIs and SLOs:

* A Service Level Input (SLI) identifies the data to be analysed
  as a query to a data provider
* A Service Level Objective (SLO) defines the quality criteria
  you define for each SLI.

You can specify multiple Service Level Objectives (SLOs)
that are evaluated in your Analysis
and you can weight the different analyses appropriately.
At the end of the analysis,
the status returns whether your objective failed, passed,
or passed with a warning.
This is similar to the functionality provided by the Keptn v1
[Quality Gates](https://keptn.sh/docs/1.0.x/define/quality-gates/)
feature.

Converters are provided to
to migrate most Keptn v1
[SLIs](https://keptn.sh/docs/1.0.x/reference/files/sli/)
and
[SLOs](https://keptn.sh/docs/1.0.x/reference/files/slo/)
to Keptn Analysis SLIs and SLOs.
For more information,see:

* [SLO converter](https://github.com/keptn/lifecycle-toolkit/blob/main/metrics-operator/converter/slo_converter.md#slo-converter)
* [SLI converter](https://github.com/keptn/lifecycle-toolkit/blob/main/metrics-operator/converter/sli_converter.md#sli-converter)
* [Migrate Quality Gates](../migrate/metrics-observe.md)

The Analysis result is exposed as an OpenTelemetry metric
and can be displayed on dashboard tools, such as Grafana.

> **Note** A preliminary release of the Keptn Analysis feature
  is included in Keptn v0.8.3 and v0.9.0 but is hidden behind a feature flag.
  See the
  [Analysis](../yaml-crd-ref/analysis.md/#differences-between-versions)
  reference page for instructions to activate the preview of this feature.

## Keptn Analysis basics

A Keptn Analysis is implemented with three resources:

* [AnalysisValueTemplate](../yaml-crd-ref/analysisvaluetemplate.md)
  defines the SLI with the `KeptnMetricsProvider` (data source)
  and the query to perform for each SLI

  Each `AnalysisValueTemplate` resource identifies the data source
  and the query for the analysis of the SLI.
  One `Analysis` can use data from multiple instances
  of multiple types of data provider;
  you must define a
  [KeptnMetricsProvider](../yaml-crd-ref/metricsprovider.md)
  resource for each instance of each data provider you are using.
  The template refers to that provider and queries it.

* [AnalysisDefinition](../yaml-crd-ref/analysisdefinition.md)
  define the list of SLOs for an `Analysis`

  An `AnalysisDefinition` resource contains a list of objectives to satisfy.
  Each of these objectives must specify:

  * Failure or warning target criteria
  * Whether the objective is a key objective
    meaning that its failure fails the Analysis
  * Weight of the objective on the overall Analysis
  * The `AnalysisValueTemplate` resource that contains the SLIs,
    defining the data provider from which to gather the data
    and how to compute the Analysis

* [Analysis](../yaml-crd-ref/analysis.md)
  define the specific configurations and the Analysis to report.

  An `Analysis` resource customizes the templates
  defined inside an `AnalysisDefinition` resource
  by adding configuration information such as:

  * Timeframe that specifies the range to use
    for the corresponding query in the `AnalysisValueTemplate`
  * Map of key/value pairs that can be used
    to substitute placeholders in the `AnalysisValueTemplate`

## Example Analysis

Consider the following `Analysis` resource:

{{< embed path="/metrics-operator/config/samples/metrics_v1alpha3_analysis.yaml" >}}

This `Analysis` resource:

* Defines the `timeframe` for which the analysis is done
  as between 5 am and 10 am on the 5th of May 2023
* Adds a few specific key-value pairs that will be substituted in the query.
  For instance, the query could contain the `{{.nodename}}` variable.
  The value of the `args.nodename` field (`test`)
  will be substituted for this string.

The `AnalysisDefinition` resource references this `Analysis` resource
by its `name` and `namespace` and can be seen here:

{{< embed path="/metrics-operator/config/samples/metrics_v1alpha3_analysisdefinition.yaml" >}}

This simple definition contains a single objective, `response-time-p95`.
For this objective, both failure and warning criteria are defined:

* The objective fails if the percentile 95 is less than 600
* A warning is issued when the value is between 300 and 500

The total score shows that this `Analysis`
should have an overall score of 90% to pass or 75% to get a warning.
Since only one objective is defined,
this means that the analysis either passes with 100%
(response time is less than 600)
or fails with 0% (slower response time).

The objective points to the corresponding `AnalysisValueTemplate` resource:
{{< embed path="/metrics-operator/config/samples/metrics_v1alpha3_analysisvaluetemplate.yaml" >}}

This template defines a query to a provider called `prometheus`:

```shell
 sum(kube_pod_container_resource_limits{node='{{.nodename}}'}) - sum(kube_node_status_capacity{node='{{.nodename}}'})
```

At runtime, the metrics operator tries to substitute
everything in`{{.variableName}}` format
with a key-value pair specified in the `Analysis` resource,
so, in this case, the query becomes:

```shell
 sum(kube_pod_container_resource_limits{node='test'}) - sum(kube_node_status_capacity{node='test'})
```

The other key-value pairs such as 'project' and 'stage' are just examples of how one could pass to the provider
information similar to Keptn v1 objectives.
For a working example you can
check [here](https://github.com/keptn/lifecycle-toolkit/tree/main/test/testanalysis/analysis-controller-multiple-providers).

## Accessing Analysis

### Retrieve KeptnMetric values with kubectl

Use the `kubectl get` command to retrieve all the `Analysis` resources
in your cluster:

```shell
kubectl get analyses.metrics.keptn.sh -A

```

This returns something like

```shell
NAMESPACE   NAME              ANALYSISDEFINITION    STATE   WARNING   PASS
default     analysis-sample   ed-my-proj-dev-svc1
```

You can then describe the `Analysis` with:

```shell
kubectl describe analyses.metrics.keptn.sh analysis-sample -n=default
```
