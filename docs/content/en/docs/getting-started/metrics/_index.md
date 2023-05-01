---
title: Getting started with Keptn metrics
description: Learn how Keptn metrics enhances your deployment
weight: 25
---

The Keptn metrics component of the Keptn Lifecycle Toolkit
allow you to define any type of metric
from multiple instances of any type of data source in your Kubernetes cluster.
You may have tools like Argo, Flux, KEDA, HBA, Keptn
that need observability data to make automated decisions.
Whether a rollout is good, whether to scale up or down.
Your observability data may come
from multiple observability solutions --
Datadog, Dynatrace, data in AWS, Google, and Azure --
and include data in Lifestep (?) and Honeycomb or Splunk.

The [Kubernetes metric server](https://github.com/kubernetes-sigs/metrics-server)
requires that you maintain point-to-point integrations
from Argo Rollouts, Flux, KEDA, and HPA.
Each has plugins but it is difficult to maintain them,
especially if you are using multiple tools
and multible observability platforms.
The Keptn Metrics Server unifies and standardizes access to this data.

This guide walks you through the steps required
to implement Keptn metrics:

1. Install the Keptn Lifecycle Toolkit
   or just the Keptn Metrics Server in your cluster.
1. Create a CRD to define each observability platform
   that is implemented in this namespace.
   You can define a mix of platforms --
   Prometheus, Dynatrace, Datadog, etc. --
   and multiple instances of each.
1. Create a CRD that Defines
   the type of data to pull from each observability platform.
   This data is pulled and fetched continuously
   at an interval you specify for each data query.
1. Run the metrics
1. View metrics

Andi's demo environment: my-klt-demo-with-argo/simplenode-dev

## Install KLT or just metrics server

Use the Helm Chart to install the Keptn Metrics Server
as part of the Lifecycle Toolkit
or completely stand-alone.  See
[Install KLT using the Helm Chart](../../install/install.md/#use-helm-chart).

## Set up secret data providers?

## Define metrics to use

### Define metrics providers

Specify metrics I want to pull in from an external observability platfor

Two metrics identified
Specify through a CRD the type of data I want to input to my Keptn Metrics Server.
I can pull and fetch that data continuously into Prometheus
Data is available through the CRD and through Prometheus itself
as well as the Kubernetes CLI

TODO: Need to redo these to use the v1alpha3 synatax

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

TODO: Terminology question: is this one file that includes
two CRD's or one CRD that includes multiple metrics?
What if one did multiple queries for a metric provider?

TODO: Need to redo these to use the v1alpha3 synatax

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

TODO: Do I need to start and stop anything to start gathering metrics
or could I theoretically just put these pieces into my cluster
and would it start gathering metrics that I could then view?

## Observing the metrics

TODO: Do we want to say anything about running these metrics,
viewing the results, perhaps from CLI and from Grafana?

## Implementing autoscaling with HPA

The Kubernetes HorizontalPodAutoscaler (HPA)
uses metrics to provide autoscaling for the cluster.
HPA can retrieve KeptnMetrics and use it to implement HPA.

TODO: Link to HPA section in "Implementing"
(which includes link to Flo's blog post)

## Learn more

To learn more about the Keptn Metrics Server, see:

* Architecture:
  [Keptn Metrics Operator](../../concepts/architecture/components/metrics-operator/)
* More information about implementing Keptn Metrics:
  [Keptn Metrics](../../implementing/metrics.md/)
