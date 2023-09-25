---
title: Analyses
description: Understand Analyses in Keptn and how to use them
weight: 150
---

The Keptn Metrics Operator implements an SLO/SLI feature set inspired by Keptn v1 under the name of Analysis.
With an Analysis Definition you can specify multiple Service Level Objectives that will be evaluated in your Analysis.  
At the end of the Analysis the status will return whether your objective failed or passed.

This data can be seen as Prometheus metrics and can be displayed on Grafana.

## Keptn Analysis basics

Keptn Analysis are implemented with three resources:

* [Analysis](https://lifecycle.keptn.sh/docs/crd-ref/metrics/v1alpha3/#analysis) --
  define the specific configurations and the analysis to report
* [AnalysisDefinition](https://lifecycle.keptn.sh/docs/crd-ref/metrics/v1alpha3/#analysisdefinition) --
  define the list of SLOs for an analysis
* [AnalysisValueTemplate](https://lifecycle.keptn.sh/docs/crd-ref/metrics/v1alpha3/#analysisvaluetemplate) --
  define the SLI: the KeptnMetricsProvider and the Query to perform for each SLO

### Define Analysis, Analysis Definition and AnalysisValueTemplate

An Analysis customizes the templates defined inside an AnalysisDefinition by adding configurations such as:
* a Timeframe that specifies the range for the corresponding query in the AnalysisValueTemplate
* a map of key/value pairs that can be used to substitute placeholders in the AnalysisValueTemplate

An AnalysisDefinition contains a list of objectives to satisfy. 
Each of these objectives:
* specifies failure or warning target criteria, 
* specifies whether the objective is a key objective (its failure would fail the analysis)
* indicates the weight of the objective on the overall analysis
* refers to an AnalysisValueTemplate that contains the SLIs, so from what provider to gather the data and how to compute the analysis

In each AnalysisValueTemplate we store the query for the analysis of the SLI. You must define a
[KeptnMetricsProvider](../yaml-crd-ref/metricsprovider.md) resource
for each instance of each data provider you are using.
The template will refer to that provider and query it.

Let's consider the following Analysis: 

{{< embed path="/lifecycle-operator/config/samples/metrics_v1alpha3_analysis.yaml" >}}

This cr setups the timeframe we are interested in as between 5am and 10am of the 5th of May 2023 , 
and adds a few specific key-value pair that will be substituted in the query 
for instance the query could contain a {{.nodename }} and this value will be substituted by test

The definition of this Analysis is referenced by its name and namespace and can be seen here:

{{< embed path="/lifecycle-operator/config/samples/metrics_v1alpha3_analysisdefinition.yaml" >}}

This simple definition contains a single objective, response-time-p95. For this objective there are both a
failure and warning criteria: 

* objective will fail if the percentile 95 is less than 600 
* there will be a warning in case the value is in between 300 and 500

The total score shows that this analysis should overall score  90% of all objectives to pass or 75 to get a warning.
Since the objective is one only, this means that we either will pass with 100% (response time is less than 600) or fail with 0%(slower response time)

The objective points to the corresponding AnalysisValueTemplate:

{{< embed path="/lifecycle-operator/config/samples/metrics_v1alpha3_analysisdvaluetemplate.yaml" >}}

This template tell us that we will query a provider called prometheus using this query:
```shell
 sum(kube_pod_container_resource_limits{node='{{.nodename}}'}) - sum(kube_node_status_capacity{node='{{.nodename}}'})
```

at runtime the metrics operator will try to substitute everything in '{{. }}' format with a key-value pair in the Analysis resource,
so in this case the query would become:

```shell
 sum(kube_pod_container_resource_limits{node='test'}) - sum(kube_node_status_capacity{node='test'})
```

For a working example you can check [here](https://github.com/keptn/lifecycle-toolkit/tree/main/test/integration/analysis-controller-multiple-providers) 

## Accessing Analysis

### Retrieve KeptnMetric values with kubectl
Use the `kubectl get` command to retrieve all the `Analysis` in your cluster: 

```shell
$  kubectl get analyses.metrics.keptn.sh -A

```
This will return something like 

```shell
NAMESPACE   NAME              ANALYSISDEFINITION    STATE   WARNING   PASS
default     analysis-sample   ed-my-proj-dev-svc1
```

You can than describe the analysis like so:

```shell
kubectl describe analyses.metrics.keptn.sh analysis-sample -n=default
```