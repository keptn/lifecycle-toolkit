# Keptn Metrics

The Keptn Metrics Operator provides a single entry point
to all metrics in the cluster
and allows you to define metrics based on multiple data platforms
and multiple instances of any data platform.
Metrics are fetched independently
and can be used for an evaluation at [workload-](https://kubernetes.io/docs/concepts/workloads/)
and application-level, or for scaling your [workloads](https://kubernetes.io/docs/concepts/workloads/).

This data can be displayed on Grafana
or another standard dashboard application that you configure
or can be retrieved using standard Kubernetes commands.

For an introduction to Keptn metrics, see
[Getting started with Keptn metrics](../getting-started/metrics.md).

## Keptn metric basics

Keptn metrics are implemented with two resources:

* [KeptnMetric](../reference/crd-reference/metric.md) --
  define the metric to report
* [KeptnMetricsProvider](../reference/crd-reference/metricsprovider.md) --
  define the configuration for a data provider

As soon as you define and apply
your `KeptnMetricsProvider` and `KeptnMetric` resources,
Keptn begins collecting the metrics you defined.
You do not need to do anything else.

### Define KeptnMetricsProvider resources

You must define a
[KeptnMetricsProvider](../reference/crd-reference/metricsprovider.md) resource
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
  such as `keptn-system`.

To configure a data provider into your Keptn cluster:

1. Create a secret if your data provider uses one.
   See
   [Create secret text](./tasks.md#create-secret-text).
1. Install and configure each instance of each data provider
   into your Keptn cluster,
   following the instructions provided by the data source provider.
   See
   [Prepare your cluster for Keptn](../installation/k8s.md/#prepare-your-cluster-for-keptn)
for links.
   Keptn supports using multiple instances of multiple data providers.
1. Define a
   [KeptnMetricsProvider](../reference/crd-reference/metricsprovider.md)
   resource for each data source.

For example, the `KeptnMetricProvider` resource
for a Prometheus data source that does not use a secret
could look like:

```yaml
apiVersion: metrics.keptn.sh/v1beta1
kind: KeptnMetricsProvider
metadata:
  name: prometheus-provider
  namespace: simplenode-dev
spec:
  type: prometheus
  targetServer: "http://prometheus-k8s.monitoring.svc.cluster.local:9090"
```

The `KeptnMetricProvider` resource for a Dynatrace provider
that uses a secret could look like:

```yaml
apiVersion: metrics.keptn.sh/v1beta1
kind: KeptnMetricsProvider
metadata:
  name: dynatrace-provider
  namespace: podtato-kubectl
spec:
  type: dynatrace
  targetServer: "<dynatrace-tenant-url>"
  secretKeyRef:
    name: dt-api-token
    key: DT_TOKEN
```

### Define KeptnMetric information

The [KeptnMetric](../reference/crd-reference/metric.md) resource
defines the information you want to gather,
specified as a query for the particular observability platform
you are using.
You can define any type of metric from any data source.

In our example, we define two bits of information to retrieve:

* Number of CPUs, fetched from the `dev-prometheus` data platform
* `availability` SLO, fetched from the `dev-dynatrace` data platform

Each of these are configured to fetch data every 10 seconds
but you could configure a different `fetchIntervalSeconds` value
for each metric.

The
[keptn-metric.yaml](https://github.com/keptn-sandbox/klt-on-k3s-with-argocd/blob/main/simplenode-dev/keptn-metric.yaml)
file for our example looks like:

```yaml
apiVersion: metrics.keptn.sh/v1beta1
kind: Keptnmetric
metadata:
  name: available-cpus
  namespace: simplenode-dev
spec:
  provider:
    name: dev-prometheus
  query: "sum(kube_node_status_capacity{resources`cpu`})"
  fetchIntervalSeconds: 10
