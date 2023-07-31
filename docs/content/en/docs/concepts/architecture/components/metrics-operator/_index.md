---
title: Keptn Metrics Operator
linktitle: Metrics Operator
description: Basic understanding of Keptn's Metrics Operator
weight: 80
cascade:
---


The Keptn Metrics Operator collects, processes,
and analyzes metrics data from a variety of sources.
Once collected, this data can be used
to generate a variety of reports and dashboards
that provide insights into the health and performance
of the application and infrastructure.

While Kubernetes has ways to extend its metrics APIs, they have limitations,
especially that they only allow you to use a single observability platform
such as Prometheus, Dynatrace or Datadog.
The Keptn Metrics Operator solves this problem
by providing a single entry point for
all your metrics data, regardless of its source,
so you can use multiple instances of multiple observability platforms.

Keptn metrics are integrated with the Kubernetes
[Custom Metrics API](https://github.com/kubernetes/metrics#custom-metrics-api)
so are compatible with the Kubernetes
[HorizontalPodAutoscaler](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/)
(HPA), which enables the horizontal scaling of workloads
based on metrics collected from multiple observability platforms.
See
[Using the HorizontalPodAutoscaler](../../../../implementing/evaluatemetrics.md/#using-the-horizontalpodautoscaler)
for instructions.

The Metrics Operator consists of the following components:

* Metrics Controller
* Metrics Adapter

```mermaid
graph TD;
    Metrics-Operator-->Metrics-Adapter;
    Metrics-Operator-->Metrics-Controller
style Metrics-Operator fill:#006bb8,stroke:#fff,stroke-width:px,color:#fff
style Metrics-Adapter fill:#d8e6f4,stroke:#fff,stroke-width:px,color:#006bb8
style Metrics-Controller fill:#d8e6f4,stroke:#fff,stroke-width:px,color:#006bb8
```

The **Metrics adapter** exposes custom metrics from an application
to external monitoring and alerting tools.
The adapter exposes custom metrics on a specific endpoint
where external monitoring and alerting tools can scrape them.
It is an important component of the metrics operator
as it allows for the collection and exposure of custom metrics,
which can be used to gain insight into the behavior and performance
of applications running on a Kubenetes cluster.

The **Metrics controller** fetches metrics from an SLI provider.
The controller reconciles a [`KeptnMetric`](../../../../yaml-crd-ref/metric.md)
resource and updates its status with the metric value
provided by the selected metric provider.
Each `KeptnMetric` is identified by `name`
and is associated with an instance of an observability platform
that is defined in a
[KeptnMetricsProvider](../../../../yaml-crd-ref/metricsprovider.md)
resource.

The steps in which the controller fetches metrics are given below:

1. Fetch the `KeptnMetric` object to be reconciled.
   * If the object is not found,
     it returns and lets Kubernetes handle deleting all associated resources.
1. If the object is found, the code checks
   whether the metric has been updated within the interval
   that is defined in the `spec.fetchintervalseconds` field.
   * If not, it skips reconciling and requeues the request for later.
1. If the metric should be reconciled,
   it fetches the provider defined in the `spec.provider.name` field.
   * If the provider is not found,
     it returns and requeues the request for later.
1. If the provider is found,
   it loads the provider and evaluates the query
   defined in the `spec.query` field.
1. If the evaluation is successful,
   it stores the fetched value
   in the metric status field of the `KeptnMetric` object.
1. If the evaluation fails,
   the error and reason is also provided in the metric status.
   The result is also stored as a raw value
   so that, if we get something unexpected
   (such as a forbidden code),
   the user can also view it there.

   The metric status field includes the following information:

   ```yaml
   properties:
     errMsg:
       description: provides error details if the query could not be evaluated
       type: string
     lastUpdated:
       description: time when the status data was last updated
       format: date-time
       type: string
     rawValue:
       description: shows the resulting value, in raw format
       format: byte
       type: string
   ```
