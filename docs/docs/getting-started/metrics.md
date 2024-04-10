---
comments: true
---

# Keptn Metrics

The Keptn metrics component
allows you to define any type of metric
from multiple instances
of any type of data source in your Kubernetes cluster.
You may have deployment tools like Argo, Flux, KEDA, HPA, or Keptn
that need observability data to make automated decisions
such as whether a rollout is good, or whether to scale up or down.

Your observability data may come
from multiple observability solutions --
Prometheus, Thanos, Cortex, Dynatrace, Datadog and others --
or may be data that comes directly
from your cloud provider such as AWS, Google, or Azure.
The Keptn Metrics Server unifies and standardizes access to all this data.
Minimal configuration is required
because Keptn hooks directly into Kubernetes primitives.

The Keptn metrics feature
integrates metrics from all these sources into a single set of metrics.
This makes it easier to use than the
[Kubernetes metric server](https://github.com/kubernetes-sigs/metrics-server),
which requires that you maintain point-to-point integrations
from each source -- Argo Rollouts, Flux, KEDA, HPA, etc.
Each has plugins but it is difficult to maintain them,
especially if you are using multiple tools,
multiple observability platforms,
and multiple instances of some tools or observability platforms.

## Using this exercise

This exercise runs on a
[kind](https://kind.sigs.k8s.io/)
cluster that you set up locally
but could be run on any Kubernetes cluster.

1. [Set up the cluster](#set-up-the-cluster)
   by creating a local Kubernetes cluster,
   installing an instance of Prometheus to use as a data source,
   and installing Keptn on the cluster.
1. [Configure Keptn for your metrics](#configure-keptn-for-your-metrics)
   by defining custom resources
   for each data provider and each piece of data you use.
1. Run the metrics

If you want to create your own cluster to run this exercise,
follow the instructions in [Installation](../installation/index.md).

## Set up the cluster

1. Create the kind cluster:

     ```shell
     kind create cluster
     ```

1. Install and configure Prometheus as your data source

     For this simple exercise, we use on instance of
     [kube-prometheus-stack](https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack)
     as our data source.
     In actual practice, you can use multiple instances
     of multiple types of data sources to gather your metrics.

     Use the following command sequence to install an instance of Prometheus:

     ```shell
     helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
     helm repo update
     helm install kube-prom-stack prometheus-community/kube-prometheus-stack

     ```

1. Install Keptn on your cluster

     Use the following command sequence to install Keptn on your cluster:

     ```shell
     helm repo add keptn https://charts.lifecycle.keptn.sh/
     helm repo update
     helm upgrade --install keptn keptn/keptn -n keptn-system --create-namespace --wait
     ```

     For more details about how to install Keptn, see the
     [Installation Guide](../installation/index.md).

### Expose Prometheus and get an existing metric

Now we need to expose Prometheus
and chose an existing metric to use for this exercise.

1. Use port forwarding to expose prometheus and find a suitable metric:

     ```shell
     kubectl get svc
     ```

     This command shows multiple services.
     Look for the Prometheus UI (the one on port `9090`)
     which should be called `prometheus-operated` (no typo).

1. Port forward to access the UI:

     ```shell
     kubectl -n default port-forward svc/prometheus-operated 9090
     ```

1. Go to `http://localhost:9090/` and, in the search box,
   type `prometheus_` and select `pr`.
   This produces a list of metrics that are defined

1. For this exercise, we are going to use a single metric.
   In real practice, you would be defining many metrics.
   Choose a metric that has a non-zero value.
   For example:

     ```shell
     prometheus_http_requests_total{code="200", container="prometheus", \
          endpoint="http-web", handler="/-/healthy", instance="10.244.0.15:9090", \
          job="kube-prom-stack-kube-prome-prometheus", namespace="default",\
          pod="prometheus-kube-prom-stack-kube-prome-prometheus-0", \
          service="kube-prom-stack-kube-prome-prometheus"}

     ```

     For a cleaner display, remove some attributes:

     ```shell
     prometheus_http_requests_total{code="200", container="prometheus", \
          endpoint="http-web", handler="/-/healthy", \
          job="kube-prom-stack-kube-prome-prometheus", namespace="default", \
          service="kube-prom-stack-kube-prome-prometheus"}
     ```

     This is your metric and value as it is stored in Prometheus.

## Configure Keptn for your metrics

Now you must tell Keptn about the metrics you are using.j
You do this by defining:

- A
  [KeptnMetricsProvider](../reference/crd-reference/metricsprovider.md)
  resource to define the external observability platform you are using
  as a data source.
  For this exercise, this is the Prometheus server
  you installed and exposed above.
- A
  [KeptnMetric](../reference/crd-reference/metric.md)
  resource to define each metric query you want to pull.

The steps are:

1. Create a
   [KeptnMetricsProvider](../../reference/crd-reference/metricsprovider.md/)
   resource for the observability platform you are using as a data source.
   To do this, create a new `.yml` file with content like the following:

      ```shell
        apiVersion: metrics.keptn.sh/v1alpha3
        kind: KeptnMetricsProvider
        metadata:
          name: local-prometheus
          namespace: default
        spec:
          type: prometheus
          targetServer: "http://prometheus-operated.default.svc.cluster.local:9090/"
      ```

     You can specify a virtually unlimited number of providers,
     including multiple instances of each observability platform.
     Each one must be assigned a unique `name`,
     identified by the `type` of platform it is
     and the URL of the target server.
     If the target server is protected by a Secret,
     provide information about the token and key.

     Keptn uses the `name` you assign to reference this data source.

     The `targetServer` field tells Keptn where to find this data source;
     in this case, it points to the UI on port 9090.

     Apply this file with the following command:

     ```shell
      kubectl apply -f YOUR-KEPTN-METRIC-PROVIDER.yml
      ```

1. Define your
   [KeptnMetric](../../reference/crd-reference/metric.md/)
   custom resource for each piece of data you want to pull
   (in this case, the Prometheus query for the metric you selected).

     This is where you use the Prometheus query from before.
     Note: The namespaces of `KeptnMetricsProvider` and `KeptnMetric` must match.

     ```yaml
     apiVersion: metrics.keptn.sh/v1alpha3
     kind: KeptnMetric
     metadata:
       name: prometheus-http-requests-total
       namespace: default
     spec:
       provider:
         name: local-prometheus
       query: 'prometheus_http_requests_total{code="200", handler="/-/healthy", job="kube-prom-stack-kube-prome-prometheus", namespace="default", service="kube-prom-stack-kube-prome-prometheus"}'
       fetchIntervalSeconds: 10
     ```

     This data is pulled and fetched continuously
     at an interval you specify for each specific bit of data.
     Data is available through the resource
     and through the data provider itself,
     as well as the Kubernetes CLI.

### Define KeptnMetric information

The [KeptnMetric](../reference/crd-reference/metric.md) resource
defines the information you want to gather,
specified as a query for the particular observability platform
you are using.
You can define any type of metric from any data source.

In our example, we define two bits of information to retrieve:

- Number of CPUs, derived from the `dev-prometheus` data platform
- `availability` SLO, derived from the `dev-dynatrace` data platform

Each of these are configured to fetch data every 10 seconds
but you could configure a different `fetchIntervalSeconds` value
for each metric.

Note the following:

- Populate one YAML file per metric
  then apply all of them.
- Each metric is assigned a unique `name`.
- The value of the `spec.provider.name` field in the `KeptnMetric` resource
  must correspond to the name assigned in
  the `metadata.name` field of a `KeptnMetricsProvider` resource.
- Information is fetched in on a continuous basis
  at a rate specified
  by the value of the `spec.fetchIntervalSeconds` field.

## View available metrics

Keptn automatically starts pulling the metrics
(every `fetchIntervalSeconds` seconds):

```shell
kubectl -n default get keptnmetrics
```

The above command should show:

```shell
NAME                             PROVIDER           QUERY INTERVAL   VALUE
prometheus-http-requests-total   local-prometheus   prometheus_http_requests_total{code="200", handler="/-/healthy", job="kube-prom-stack-kube-prome-prometheus", namespace="default", service="kube-prom-stack-kube-prome-prometheus"}              756
```

or to retrieve a specific KeptnMetric value:

```shell
kubectl -n default get keptnmetric/prometheus-http-requests-total 
```

Use the following command to view
the metrics that are configured in your cluster:

```shell
kubectl get KeptnMetrics -A
```

## Run the metrics

As soon as you define your `KeptnMetricsProvider` and `KeptnMetric` resources,
Keptn begins collecting the metrics you defined.
You do not need to do anything else.

## Observing the metrics

The metrics can be retrieved
through CRs and through the Kubernetes Metric API.

The syntax to retrieve metrics from the CR is:

```shell
kubectl get keptnmetrics.metrics.keptn.sh -n <namespace> <metric-name>
```

The syntax to retrieve metrics through the Kubernetes API  is:

```yaml
kubectl get --raw "/apis/custom.metrics.k8s.io/v1beta2/namespaces/<namespace>/keptnmetrics.metrics.sh/<metric-name>/<metric-name>"
```

You can also display the metrics graphically using a dashboard such as Grafana.

## Learn more

To learn more about the Keptn Metrics Server, see:

- Architecture:
  [Keptn Metrics Operator](../components/metrics-operator.md)
- More information about implementing Keptn Metrics:
  [Keptn Metrics](../guides/evaluatemetrics.md)
- How to integrate Keptn metrics with
  the Kubernetes HorizontalPodAutoscaler (HPA)
  to provide autoscaling for the cluster:
  [Using the HorizontalPodAutoscaler](../use-cases/hpa.md)
