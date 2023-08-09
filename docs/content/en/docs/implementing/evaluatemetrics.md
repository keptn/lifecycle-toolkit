---
title: Keptn Metrics
description: Implement Keptn metrics
weight: 85
---

The Keptn Metrics Operator provides a single entry point
to all metrics in the cluster
and allows you to define metrics based on multiple data platforms
and multiple instances of any data platform.
Metrics are fetched independently
and can be used for an evaluation at workload- and application-level, or for scaling your workloads.

This data can be displayed on Grafana
or another standard dashboard application that you configure
or can be retrieved using standard Kubernetes commands.

For an introduction to Keptn metrics, see
[Getting started with Keptn metrics](../getting-started/metrics).

## Keptn metric basics

Keptn metrics are implemented with two resources:

* [KeptnMetric](../yaml-crd-ref/metric.md) --
  define the metric to report
* [KeptnMetricsProvider](../yaml-crd-ref/metricsprovider.md) --
  define the configuration for a data provider

### Define KeptnMetricsProvider resources

You must define a
[KeptnMetricsProvider](../yaml-crd-ref/metricsprovider.md) resource
for each instance of each data provider you are using.

Note the following:

* Each `KeptnMetricsProvider` resource is bound to a specific namespace.
* Each `KeptnMetric` resource must be located in the same namespace
  as the associated `KeptnMetricsProvider` resource.
* `KeptnEvaluationDefinition` resources can reference metrics
  from any namespace in the cluster.
* To define metrics that can be used in evaluations
  on all namespaces in the cluster,
  create `KeptnMetricsProvider` and `KeptnMetric` resources
  in a centralized namespace
  such as `keptn-lifecycle-toolkit-system`.

To configure a data provider into your KLT cluster:

1. Create a secret if your data source uses one.
   See
   [Create secret text](../implementing/tasks/#create-secret-text).
1. Install and configure each instance of each data source
   into your KLT cluster,
   following the instructions provided by the data source provider.
   See
   [Prepare your cluster for KLT](../install/k8s.md/#prepare-your-cluster-for-klt)
for links.
   KLT supports using multiple instances of multiple data sources.
1. Define a
   [KeptnMetricsProvider](../yaml-crd-ref/metricsprovider.md)
   resource for each data source.

For example, the `KeptnMetricProvider` resource
for a Prometheus data source that does not use a secret
could look like:

```yaml
apiVersion: metrics.keptn.sh/v1alpha2
kind: KeptnMetricsProvider
metadata:
  name: prometheus
  namespace: simplenode-dev
spec:
  targetServer: "http://prometheus-k8s.monitoring.svc.cluster.local:9090"
```

The `KeptnMetricProvider` resource for a Dynatrace data source
that uses a secret could look like:

```yaml
apiVersion: metrics.keptn.sh/v1alpha3
kind: KeptnMetricsProvider
metadata:
  name: dynatrace
  namespace: podtato-kubectl
spec:
  targetServer: "<dynatrace-tenant-url>"
  secretKeyRef:
    name: dt-api-token
    key: DT_TOKEN
```

## Accessing Metrics via the Kubernetes Custom Metrics API

`KeptnMetrics` can also be retrieved via the Kubernetes Custom Metrics API.

### Retrieve KeptnMetric values with kubectl

Use the `kubectl get --raw` command
to retrieve the values of a `KeptnMetric`, as in the following example:

```shell
$ kubectl get --raw "/apis/custom.metrics.k8s.io/v1beta2/namespaces/podtato-kubectl/keptnmetrics.metrics.sh/keptnmetric-sample/keptnmetric-sample" | jq .

{
  "kind": "MetricValueList",
  "apiVersion": "custom.metrics.k8s.io/v1beta2",
  "metadata": {},
  "items": [
    {
      "describedObject": {
        "kind": "KeptnMetric",
        "namespace": "podtato-kubectl",
        "name": "keptnmetric-sample",
        "apiVersion": "metrics.keptn.sh/v1alpha1"
      },
      "metric": {
        "name": "keptnmetric-sample",
        "selector": {
          "matchLabels": {
            "app": "frontend"
          }
        }
      },
      "timestamp": "2023-01-25T09:26:15Z",
      "value": "10"
    }
  ]
}
```

### Filter on matching labels

You can filter based on matching labels.
For example, to retrieve all metrics
that are labelled with `app=frontend`,
use the following command:

```shell
$ kubectl get --raw "/apis/custom.metrics.k8s.io/v1beta2/namespaces/podtato-kubectl/keptnmetrics.metrics.sh/*/*?labelSelector=app%3Dfrontend" | jq .

{
  "kind": "MetricValueList",
  "apiVersion": "custom.metrics.k8s.io/v1beta2",
  "metadata": {},
  "items": [
    {
      "describedObject": {
        "kind": "KeptnMetric",
        "namespace": "keptn-lifecycle-toolkit-system",
        "name": "keptnmetric-sample",
        "apiVersion": "metrics.keptn.sh/v1alpha3"
      },
      "metric": {
        "name": "keptnmetric-sample",
        "selector": {
          "matchLabels": {
            "app": "frontend"
          }
        }
      },
      "timestamp": "2023-01-25T09:26:15Z",
      "value": "10"
    }
  ]
}
```

## Using the HorizontalPodAutoscaler

Use the Kubernetes Custom Metrics API
to refer to `KeptnMetric` via the
[Kubernetes HorizontalPodAutoscaler](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/)
(HPA),
as in the following example:

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: podtato-head-entry
  namespace: podtato-kubectl
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: podtato-head-entry
  minReplicas: 1
  maxReplicas: 10
  metrics:
    - type: Object
      object:
        metric:
          name: keptnmetric-sample
        describedObject:
          apiVersion: metrics.keptn.sh/v1alpha1
          kind: KeptnMetric
          name: keptnmetric-sample
        target:
          type: Value
          value: "10"
```

See the
[Scaling Kubernetes Workloads based on Dynatrace Metrics](https://www.linkedin.com/pulse/scaling-kubernetes-workloads-based-dynatrace-metrics-keptnproject/)
blog post
for a detailed discussion of doing this with Dynatrace metrics.
The same approach could be used to implement HPA with other data providers.
