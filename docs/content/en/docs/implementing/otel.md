---
title: OpenTelemetry observability
description: How to standardize access to OpenTelemetry observability data
weight: 140
---

To access OpenTelemetry metrics with the Keptn Lifecycle Toolkit,
you must:

- Install an OpenTelemetry collector on your cluster.
  See
  [OpenTelemetry Collector](https://opentelemetry.io/docs/collector/)
  for more information.
- Apply
  [basic annotations](../implementing/integrate/#basic-annotations)
  for your `Deployment` resource
  to integrate the Lifecycle Toolkit into your Kubernetes cluster.

KLT begins to collect OpenTelemetry metrics
as soon as the `Deployment` resource
has the basic annotations to integrate KLT in the cluster.

To expose OpenTelemetry metrics,
define a [KeptnConfig](../yaml-crd-ref/config.md) resource
that has the `spec.OTelCollectorUrl` field populated
with the URL of the OpenTelemetry collector.

Keptn metrics can be exposed as OpenTelemetry (OTel) metrics
via port `9999` of the KLT metrics-operator.

To access the metrics, use the following command:

```shell
kubectl port-forward deployment/metrics-operator 9999 -n keptn-lifecycle-toolkit-system
```

You can access the metrics from your browser at: `http://localhost:9999`

For an introduction to using OpenTelemetry with Keptn metrics, see the
[Standardize observability](../getting-started/observability)
getting started guide.
