---
title: OpenTelemetry observability
description: How to standardize access to OpenTelemetry observability data
weight: 40
---


Keptn makes any Kubernetes deployment observable.
In other words, it creates a distributed, end-to-end trace
of what Kubernetes does in the context of a Deployment.
To do this,
Keptn introduces the concept of an `application`,
which is an abstraction that connects multiple
[Workloads](https://kubernetes.io/docs/concepts/workloads/) that logically belong together,
even if they use different deployment strategies.

This means that:

- You can readily see why a deployment takes so long
  or why it fails, even when using multiple deployment strategies.
- Keptn can capture DORA metrics and expose them as OpenTelemetry metrics

The observability data is an amalgamation of the following:

- DORA metrics are collected out of the box
  when Keptn is enabled
- OpenTelemetry runs traces that show
  everything that happens in the Kubernetes cluster
- Custom Keptn metrics that you can use to monitor
  information from all the data providers configured in your cluster

All this information can be displayed with dashboard tools
such as Grafana.

For an introduction to using OpenTelemetry with Keptn metrics, see the
[Standardize observability](../intro/usecase-observability.md)
getting started guide.

## DORA metrics

DORA metrics are an industry-standard set of measurements;
see the following for a description:

- [What are DORA Metrics and Why Do They Matter?](https://codeclimate.com/blog/dora-metrics)
- [Are you an Elite DevOps Performer?
   Find out with the Four Keys Project](https://cloud.google.com/blog/products/devops-sre/using-the-four-keys-to-measure-your-devops-performance)

DORA metrics provide information such as:

- How many deployments happened in the last six hours?
- Time between deployments
- Deployment time between versions
- Average time between versions.

Keptn starts collecting these metrics
as soon as you apply
[basic annotations](./integrate.md#basic-annotations)
to the [workload](https://kubernetes.io/docs/concepts/workloads/).
Metrics are collected only for the resources
that are annotated.

To view DORA metrics, run the following two commands:

- Retrieve the service name with:

```shell
kubectl -n keptn-lifecycle-toolkit-system get service -l control-plane=lifecycle-operator
```

- Then port-forward to the name of your service:

```shell
kubectl -n keptn-lifecycle-toolkit-system port-forward service/<YOURNAME> 2222
```

Then view the metrics at:

```shell
http://localhost:2222/metrics
```

DORA metrics can be displayed on Grafana
or whatever dashboard application you choose.

## OpenTelemetry

### Requirements for OpenTelemetry

To access OpenTelemetry metrics with Keptn,
you must have the following on your cluster:

- An OpenTelemetry collector.
  See
  [OpenTelemetry Collector](https://opentelemetry.io/docs/collector/)
  for more information.
- A Prometheus Operator.
  See [Prometheus Operator Setup](https://github.com/prometheus-operator/kube-prometheus/blob/main/docs/customizing.md).

  - The Prometheus Operator must have the required permissions
    to watch resources of your Keptn namespace (default is `keptn-lifecycle-toolkit-system`) (see
    [Setup for Monitoring other Namespaces](https://prometheus-operator.dev/docs/kube/monitoring-other-namespaces/)).

  - To install Prometheus into the `monitoring` namespace
    using the example configuration included with Keptn,
    use the following command sequence.
    Use similar commands if you define a different configuration:

    > **Note**
    You must clone the `lifecycle-toolkit` repository
    and `cd` into the correct directory
    (`examples/support/observability`) before running the following commands.

    ```shell
    kubectl create namespace monitoring
    kubectl apply --server-side -f config/prometheus/setup/
    kubectl apply -f config/prometheus/
    ```

- If you want a dashboard for reviewing metrics and traces:

  - Install
    [Grafana](https://grafana.com/grafana/)
    or the visualization tool of your choice, following the instructions in
    [Grafana Setup](https://grafana.com/docs/grafana/latest/setup-grafana/).
  - Install
    [Jaeger](https://www.jaegertracing.io/)
    or a similar tool for traces following the instructions in
    [Jaeger Setup](https://www.jaegertracing.io/docs/1.50/getting-started/).

  - Follow the instructions in the Grafana
    [README](https://github.com/keptn/lifecycle-toolkit/blob/main/dashboards/grafana/README.md)
    file to configure the Grafana dashboard(s) for Keptn..

    Metrics can also be retrieved without a dashboard.
    See
    [Accessing Metrics via the Kubernetes Custom Metrics API](evaluatemetrics.md/#accessing-metrics-via-the-kubernetes-custom-metrics-api)

### Integrate OpenTelemetry into Keptn

To integrate OpenTelemetry into Keptn:

- Apply
  [basic annotations](./integrate.md#basic-annotations)
  for your `Deployment` resource
  to integrate Keptn into your Kubernetes cluster.
- To expose OpenTelemetry metrics,
  define a [KeptnConfig](../yaml-crd-ref/config.md) resource
  that has the `spec.OTelCollectorUrl` field populated
  with the URL of the OpenTelemetry collector.

The
[otel-collector.yaml](https://github.com/keptn/lifecycle-toolkit/blob/main/examples/support/observability/config/otel-collector.yaml)
is the OpenTelemetry manifest file for the PodtatoHead example,
located in the `config` directory.
To deploy and configure the OpenTelemetry collector
using this manifest, the command is:

```shell
kubectl apply -f config/otel-collector.yaml \
    -n keptn-lifecycle-toolkit-system
```

Use the following command to confirm that the pod
for the `otel-collector` deployment is up and running:

```shell
$ kubectl get pods -lapp=opentelemetry \
    -n keptn-lifecycle-toolkit-system

NAME                              READY   STATUS    RESTARTS      AGE
otel-collector-6fc4cc84d6-7hnvp   1/1     Running   0             92m
```

If you want to extend the OTel Collector configuration
to send your telemetry data to other Observability platform,
you can edit the Collector `ConfigMap` with the following command:

```shell
kubectl edit configmap otel-collector-conf \
    -n keptn-lifecycle-toolkit-system
```

When the `otel-collector` pod is up and running,
restart the `keptn-scheduler` (if installed) and `lifecycle-operator`
so they can pick up the new configuration:

```shell
kubectl rollout restart deployment \
    -n keptn-lifecycle-toolkit-system keptn-scheduler lifecycle-operator
```

Keptn begins to collect OpenTelemetry metrics
as soon as the `Deployment` resource
has the basic annotations to integrate Keptn in the cluster.

## Access Keptn metrics as OpenTelemetry metrics

Keptn metrics can be exposed as OpenTelemetry (OTel) metrics
via port `9999` of the Keptn metrics-operator.

To access the metrics, use the following command:

```shell
kubectl port-forward deployment/metrics-operator 9999 -n keptn-lifecycle-toolkit-system
```

You can access the metrics from your browser at: `http://localhost:9999`
