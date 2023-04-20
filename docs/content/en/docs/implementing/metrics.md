---
title: Keptn Metrics
description: Implement Keptn site metrics
weight: 130
---

Kubernetes provides two metrics servers,
one that is custom, the other external.
However, the Kubernetes metrics servers
only allow you to release a single service on a API;
you cannot use multiple observability platforms
in one namespace for one API.

Keptn Metrics provides a single entry point to all metrics in the cluster
and allows you to use multiple observability platforms
and, beginning with V0.8.0,
multiple instances of any observability platform.
Keptn Metrics are also application aware
so report data for all workflows that are included in the Keptn application.

This data can be presented on Grafana
or any standard dashboard application that you configure.

[More introductory info to come]

## Keptn metrics basics

Keptn metrics are implemented with two CRDs:

* [KeptnMetric](../yaml-crd-ref/metric.md) --
  define the metric to report
* [KeptnMetricsProvider](../yaml-crd-ref/metricsprovider.md) --
  define the data provider to be used for this metric

## Using OpenTelemetry with Keptn metrics

Keptn metrics can be exposed as OpenTelemetry (OTel) metrics
via port `9999` of the KLT metrics-operator.

To expose OTel metrics,
be sure that the `EXPOSE_KEPTN_METRICS` environment variable
in the `metrics-operator` manifest is set to `true`,
which is the default value.

To access the metrics, use the following command:

```shell
kubectl port-forward deployment/metrics-operator 9999 -n keptn-lifecycle-toolkit-system
```

You can access the metrics from your browser at:

```http://localhost:9999/metrics```

## Accessing Metrics via the Kubernetes Custom Metrics API

`KeptnMetrics` can also be retrieved via the Kubernetes Custom Metrics API.

### Using the HorizontalPodAutoscaler

Use the Kubernetes Custom Metrics API
to refer to Keptnmetrics via the
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

See the [Scaling Kubernetes Workloads based on Dynatrace Metrics](https://www.linkedin.com/pulse/scaling-kubernetes-workloads-based-dynatrace-metrics-keptnproject/)
blog post
for a detailed discussion of doing this with Dynatrace metrics.
A similar approach could be used to implement HPA with other data providers.

### Retrieve KeptnMetric values with kubectl raw

Use the `kubectl raw` command
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
