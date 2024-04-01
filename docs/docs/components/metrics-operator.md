---
comments: true
---

# Keptn Metrics Operator

The Keptn Metrics Operator collects, processes,
and analyzes metrics data from a variety of sources.
Once collected, this data can be used
to generate a variety of reports and dashboards
that provide insights into the health and performance
of the application and infrastructure.

While Kubernetes has ways to extend its metrics APIs, there are limitations,
especially that they allow you to use only a single observability platform
such as Prometheus, Thanos, Dynatrace or Datadog.
The Keptn Metrics Operator solves this problem
by providing a single entry point for
all your metrics data, regardless of its source.
This means that data multiple observability platforms are available via
a single access point.

Keptn metrics are integrated with the Kubernetes
[Custom Metrics API](https://github.com/kubernetes/metrics#custom-metrics-api),
so they are compatible with the Kubernetes
[HorizontalPodAutoscaler](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/)
(HPA), which enables the horizontal scaling of workloads
based on metrics collected from observability platforms.
See
[Using the HorizontalPodAutoscaler](../use-cases/hpa.md)
for instructions.

The Metrics Operator consists of the following components:

* Metrics Controller
* Analysis Controller
* Metrics Adapter

```mermaid
graph TD;
K((CRs)) -- apply --> L[Kubernetes API]
X[Metrics Adapter] <--> L
Y[Metrics Controller] <--> L
Z[Analysis Controller] <--> L

P3((<svg height="80" width="80"><image href="https://upload.wikimedia.org/wikipedia/commons/thumb/3/38/Prometheus_software_logo.svg/2066px-Prometheus_software_logo.svg.png" height="80" width="80" /></svg>))
P1[<svg height="80" width="100"><image href="https://imgix.datadoghq.com/img/about/presskit/usage/logousage_white.png" height="100" width="100" /></svg>]  <--> Y
P3 <--> Y
P3 <--> Z
P2[<svg height="70" width="100"><image href="https://seeklogo.com/images/D/dynatrace-logo-0B89594073-seeklogo.com.png" height="70" width="100" /></svg>] <--> Z

style L fill:#006bb8,stroke:#fff,stroke-width:px,color:#fff
style Y fill:#d8e6f4,stroke:#fff,stroke-width:px,color:#006bb8
style Z fill:#d8e6f4,stroke:#fff,stroke-width:px,color:#006bb8
style X fill:#d8e6f4,stroke:#fff,stroke-width:px,color:#006bb8
style K fill:#fff,stroke:#123,stroke-width:px,color:#006bb8
style P1 fill:#fff,stroke:#fff,stroke-width:px,color:#fff
style P2 fill:#fff,stroke:#fff,stroke-width:px,color:#fff
style P3 fill:#fff,stroke:#fff,stroke-width:px,color:#fff
```

The **Metrics adapter** exposes custom metrics from an application
to external monitoring and alerting tools.
The adapter exposes custom metrics on a specific endpoint
where external monitoring and alerting tools can scrape them.
It is an important component of the metrics operator
as it allows for the collection and exposure of custom metrics,
which can be used to gain insight into the behavior and performance
of applications running on a Kubernetes cluster.

The **Metrics controller** fetches metrics from an SLI provider.
The controller reconciles a [`KeptnMetric`](../reference/crd-reference/metric.md)
resource and updates its status with the metric value
provided by the selected metric provider.
Each `KeptnMetric` is identified by `name`
and is associated with an instance of an observability platform
that is defined in a
[KeptnMetricsProvider](../reference/crd-reference/metricsprovider.md)
resource.

The steps in which the controller fetches metrics are given below:

1. When a [`KeptnMetric`](../reference/crd-reference/metric.md)
   resource is found or modified,
   the controller checks whether the metric has been updated
   within the interval that is defined in the `spec.fetchintervalseconds` field.
   * If not, it skips the reconciliation process
     and queues the request for later.

2. The controller attempts to fetch the provider defined in the
   `spec.provider.name` field.
   * If this is not possible, the controller reconciles
     and queues the request for later.

3. If the provider is found,
   the controller loads the provider and evaluates the query
   defined in the `spec.query` field.
   * If the evaluation is successful,
     it stores the fetched value
     in the `status` field of the `KeptnMetric` object.
   * If the evaluation fails,
     the error and reason is written to the
     [KeptnMetricStatus](../reference/api-reference/metrics/v1/index.md#keptnmetricstatus)
     resource.
     The error is described in both human-readable language
     and as raw data to help identify the source of the problem
     (such as a forbidden code).
