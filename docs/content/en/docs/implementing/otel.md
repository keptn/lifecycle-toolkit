---
title: OpenTelemetry observability`
description: How to standardize access to OpenTelemetry observability data
weight: 140
---
## Using OpenTelemetry with Keptn metrics

Keptn metrics can be exposed as OpenTelemetry (OTel) metrics
via port `9999` of the KLT metrics-operator.

To expose OTel metrics:

* Be sure that the `EXPOSE_KEPTN_METRICS` environment variable
  in the `metrics-operator` manifest is set to `true`,
  which is the default value.
* Define a [KeptnConfig](../yaml-crd-ref/config.md) CRD
  that has the `spec.OTelCollectorUrl` field populated
  with the URL of the OpenTelemetry collector.

To access the metrics, use the following command:

```shell
kubectl port-forward deployment/metrics-operator 9999 -n keptn-lifecycle-toolkit-system
```

You can access the metrics from your browser at: `http://localhost:9999`

For an introduction to using OpenTelemetry with Keptn metrics, see the
[Standardize access to observability data](../getting-started/observability]
getting started guide.
