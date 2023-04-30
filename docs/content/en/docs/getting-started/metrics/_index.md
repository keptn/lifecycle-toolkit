---
title: Getting started with Keptn metrics
description: Learn how Keptn metrics enhances your deployment
weight: 25
---

Intro points:

* Standardize way you can access observability data
in your Kubernetes cluster.
* May have tools like Argo, Flux, KEDA, HBA, Keptn
that need observability data to make automated decisions.
Whether a rollout is good, whether to scale up or down.
* If you have only one data source such as Prometheus
for your observability data, not the greatest use case.
But, if you have multiple observability solution --
Datadog, Dynatrace, data in AWS, Google, and Azure.
Data in Lifestep (?) and Honeycomb or Splunk.

  * In Kubernetes, you would need to maintain point-to-point integrations
from Argo Rollouts, Flux, KEDA, and HPA.
  * They all have some plugins but really hard to maintain,
especially as you move from one tool to another,
one observability platform to another,
it's really hard to maintain this.
  * The Keptn Metrics Server unifies and standardizes access to this data.

## Install KLT or just metrics server

You can install the Keptn Metrics Server completely stand-alone
or as part of the toolkit

Demo:
* Kubernetes cluster installed
* Install full KLT or modify Helm chart to only install the Keptn Metrics Server

Specify metrics I want to pull in from an external observability platfor

Two metrics identified
Specify through a CRD the type of data I want to input to my Keptn Metrics Server.
I can pull and fetch that data continuously into Prometheus
Data is available through the CRD and through Prometheus itself
as well as the Kubernetes CLI

## Set up secret for Dynatrace

Andi's demo environment: my-klt-demo-with-argo/simplenode-dev

## Define metrics to use

### Define metrics providers

**Need to redo these to use the v1alpha3 synatax**

```yaml
kind: KeptnMetricsProvider
metadata:
  name: prometheus
  namespace: simplenode-dev
spec:
  targetserver: "http://prometheus-k8s-monitoring-svc.cluster.local:9090"
```

```yaml
kind: KeptnMetricsProvider
metadata:
  name: dynatrace
  namespace: simplenode-dev
spec:
  targetServer: "https://hci34192.live.dynatrace.com
  secretKeyRef
    name: dynatrace
    key: DT_TOKEN
...
```

### Define KeptnMetric information

Define the information I want to retrieve

You can define multiple metrics from different metric providers
and multiple instances of each provider

All in one `keptn-metric.yaml` file

**Terminology question: is this one file that includes
two CRD's or one CRD that includes multiple metrics?
What if one did multiple queries for a metric provider?**

**Need to redo these to use the v1alpha3 synatax**

Check available CPUs using Prometheus

Check the availability SLO metric,
retrieved from Dynatrace:

```yaml
apiVersion: metrics.keptn.sh/v1alpha2
kind: Keptnmetric
metadata:
  name: available-cpus
  namespace: simplenode-dev
spec:
  provider:
    name: prometheus
  query: "sum(kube_node_status_cvapacity{resources`cpu`})
  fetchIntervalSeconds" 10
...
apiVersion: metrics.keptn.sh/v1alpha2
kind: Keptnmetric
metadata:
  name: availability-slo
  namespace: simplenode-dev
spec:
  provider:
    name: dynatrace
  query: "func:slo.availability_simplenodeservice"
  fetchIntervalSeconds" 10
...
```

One file with two definitions.
Information is fetched in on a continuous basis;
both fetch every 10 seconds.

Summary: you can define any type of metric
from any data source

### View available metrics

```shell
get KeptnMetrics -A
```
```shell
NAMESPACE       NAME              PROVIDER   QUERY
simplenode-dev  availability-slo  dynatrace  func:slo.availability_simplenodeservice
simplenode-dev  available-cpus    prometheus sum(kube_node_status_capacity{resource=`coy})
```

## Run the metrics

* Do I need to start and stop anything to start gathering metrics
  or could I theoretically just put these pieces into my cluster
  and would it start gathering metrics that I could then view?

## Observing the metrics

Do we want to do anything about running these metrics,
viewing the results, perhaps from CLI and from Grafana?

## Implementing autoscaling with HPA

The Kubernetes HorizontalPodAutoscaler (HPA)
uses metrics to provide autoscaling for the cluster.
HPA can retrieve KeptnMetrics and use it to implement HPA.

TODO: Link to HPA section in "Implementing"
(which includes link to Flo's blog post)

