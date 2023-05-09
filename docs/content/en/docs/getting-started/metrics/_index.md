---
title: Getting started with Keptn metrics
description: Learn how Keptn metrics enhances your deployment
weight: 25
---

The Keptn metrics component of the Keptn Lifecycle Toolkit
allow you to define any type of metric
from multiple instances of any type of data source in your Kubernetes cluster.
You may have tools like Argo, Flux, KEDA, HPA, or Keptn
that need observability data to make automated decisions.
Whether a rollout is good, whether to scale up or down.
Your observability data may come
from multiple observability solutions --
Datadog, Dynatrace, Lightstep, Honeycomb, Splunk, or data directly from your cloud provider such as AWS, Google, and Azure.

## Using this exercise

This exercise is based on the
[simplenode-dev](https://github.com/keptn-sandbox/klt-on-k3s-with-argocd)
example.
You can clone that repo to access it locally
or just look at it for examples
as you implement the functionality "from scratch"
on your local Kubernetes deployment cluster.

The steps to implement pre- and post-deployment orchestration are:

1. [Bring or create a Kubernetes cluster](#bring-or-create-a-kubernetes-deployment-cluster)
1. [Install the Keptn Lifecycle Toolkit on your cluster](#install-klt-on-your-cluster)
1. [Enable KLT for your cluster](#enable-klt-for-your-cluster)
1. [Integrate KLT with your cluster](#integrate-klt-with-your-cluster)
1. Configure metrics to use
   * [Define metrics providers](#define-metrics-providers)
   * [Define KeptnMetric information](#define-keptnmetric-information)
   * [View available metrics](#view-available-metrics)

See the
[Introducing Keptn Lifecycle Toolkit](https://youtu.be/449HAFYkUlY)
video for a demonstration of this exercise.

## Bring or create a Kubernetes deployment cluster

You can run this exercise on an existing Kubernetes cluster
or you can create a new cluster.
For personal study and demonstrations,
this exercise runs well on a local Kubernetes cluster.
See [Bring or Install a Kubernetes Cluster](../../install/k8s.md).

## Install KLT on your cluster

Install the Keptn Lifecycle Toolkit on your cluster
by executing the following command sequence:

```shell
helm repo add klt https://charts.lifecycle.keptn.sh
helm repo update
helm upgrade --install keptn klt/klt \
   -n keptn-lifecycle-toolkit-system --create-namespace --wait
```

If you only want to use Keptn's metrics features,
you can install just the `metrics-operator`
bu modifying Helm values.
See
[Install KLT](../../install/install.md)
for more information about installing the Lifecycle Toolkit.

To verify that the `metrics-operator` is installed in your cluster,
run the following command:

```shell
kubectl get pods -n keptn-lifecycle-toolkit-system
```

The output shows all components that are running on your system.

## Enable KLT for your cluster

To enable KLT for your cluster, annotate the Kubernetes
[Namespace](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/)
resource.
In this example, this is defined in the
[simplenode-dev-ns.yaml](https://github.com/keptn-sandbox/klt-on-k3s-with-argocd/blob/main/simplenode-dev/simplenode-dev-ns.yaml)
file, which looks like this:

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: simplenode-dev
  annotations:
    keptn.sh/lifecycle-toolkit: "enabled"
```

You see the annotation line that enables `lifecycle-toolkit`.
This line tells the webhook to handle the namespace

## Integrate KLT with your cluster

To integrate KLT with your cluster, annotate the Kubernetes
[Deployment](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/)
resource.
In this example, this is defined in the
[simplenode-dev-deployment.yaml](https://github.com/keptn-sandbox/klt-on-k3s-with-argocd/blob/main/simplenode-dev/simplenode-dev-deployment.yaml)
file, which includes the following lines:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: simplenode
  namespace: simplenode-dev
...
template:
    metadata:
      labels:
        app: simplenode
        app.kubernetes.io/name: simplenodeservice
      annotations:
        # keptn.sh/app: simpleapp
        keptn.sh/workload: simplenode
        keptn.sh/version: 1.0.2
        keptn.sh/pre-deployment-evaluations: evaluate-dependencies
        keptn.sh/pre-deployment-tasks: notify
        keptn.sh/post-deployment-evaluations: evaluate-deployment
        keptn.sh/post-deployment-tasks: notify
...
```

For more information about using annotations and labels
to integrate KLT into your deployment cluster, see
[Integrate KLT with your applications](../../implementing/integrate.md).

The [Kubernetes metric server](https://github.com/kubernetes-sigs/metrics-server)
requires that you maintain point-to-point integrations
from Argo Rollouts, Flux, KEDA, and HPA.
Each has plugins but it is difficult to maintain them,
especially if you are using multiple tools
and multible observability platforms.
The Keptn Metrics Server unifies and standardizes access to this data.

The steps to implement Keptn metrics are:

## Install and configure KLT

Use the Helm Chart to install the Keptn Metrics Server
as part of the Lifecycle Toolkit
or completely stand-alone.
 See
[Install KLT using the Helm Chart](../../install/install.md/#use-helm-chart).
-- End of Probably goes --

## Define metrics to use

You need to define the external observability platforms
from which you want to pull data
and then the specific data you want to pull.
This data is pulled and fetched continuously
at an interval you specify for each specific bit of data.
Data is available through the resource and through the data provider itself,
as well as the Kubernetes CLI.

### Define metrics providers

Populate a
[KeptnMetricsProvider](../../yaml-crd-ref/metricsprovider.md)
resource for each external observability platform you want to use.

For our example, we define two observability platforms:

* `dev-prometheus`
* `dev-dynatrace`

You can specify a virtually unlimited number of providers,
including multiple instances of each observability platform.
Each one must be assigned a unique name,
identified by the type of platform it is
and the URL.

> Note: The video and example application use an older syntax
  of the `KeptnMetricsProvider` and `KeptnMetric` resources.
  The syntax shown in this document is correct for v0.7.1 and later.

Definition of
[dev-prometheus](https://github.com/keptn-sandbox/klt-on-k3s-with-argocd/blob/main/simplenode-dev/keptn-prometheus-provider.yaml)
data source:

```yaml
kind: KeptnMetricsProvider
metadata:
  name: dev-prometheus
  namespace: simplenode-dev
spec:
  type: prometheus
  targetserver: "http://prometheus-k8s-monitoring-svc.cluster.local:9090"
```

Definition of the
[dev-dynatrace](https://github.com/keptn-sandbox/klt-on-k3s-with-argocd/blob/main/simplenode-dev/dynatrace-provider.yaml.tmp)
data source.
Note that the `dev-dynatrace` server is protected by a secret key
so that information is included in the provider definition:

```yaml
kind: KeptnMetricsProvider
metadata:
  name: dev-dynatrace
  namespace: simplenode-dev
spec:
  type: dynatrace
  targetServer: "https://hci34192.live.dynatrace.com
  secretKeyRef
    name: dynatrace
    key: DT_TOKEN
...
```

### Define KeptnMetric information

The [KeptnMetric](../../yaml-crd-ref/metric.md) resource
defines the information you want to gather,
specified as a query for the particular observability platform
you are using.
You can define any type of metric from any data source.

In our example, we define two bits of information to retrieve:

* Number of CPUs, derived from the `dev-prometheus` data platform
* `availability` SLO, derived from the `dev-dynatrace` data platform

Each of these are configured to fetch data every 10 seconds
but you could configure a different `fetchIntervalSeconds` value
for each metric.

The
[keptn-metric.yaml](https://github.com/keptn-sandbox/klt-on-k3s-with-argocd/blob/main/simplenode-dev/keptn-metric.yaml)
file for our example looks like:

```yaml
apiVersion: metrics.keptn.sh/v1alpha2
kind: Keptnmetric
metadata:
  name: available-cpus
  namespace: simplenode-dev
spec:
  provider:
    name: dev-prometheus
  query: "sum(kube_node_status_cvapacity{resources`cpu`})"
  fetchIntervalSeconds" 10
...
apiVersion: metrics.keptn.sh/v1alpha2
kind: Keptnmetric
metadata:
  name: availability-slo
  namespace: simplenode-dev
spec:
  provider:
    name: dev-dynatrace
  query: "func:slo.availability_simplenodeservice"
  fetchIntervalSeconds" 10
...
```

Note the following:

* You populate one YAML file
that includes all the metrics for your cluster.
* Each metric is assigned a unique `name`.
* The value of the `spec.provider.name` field
  must correspond to the name assigned in a
  the `metadata.name` field of a `KeptnMetricsProvider` resource.
* Information is fetched in on a continuous basis
at a rate specified by the value of the `spec.fetchIntervalSeconds` field.

### View available metrics

Use the following command to view
the metrics that are configured in your cluster.
This example displays the two metrics we configured above:

TODO: Is the syntax of the output changed
to include the `name` of the provider (`dev-prometheus` or `dev-dynatrace`?

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
See Using the HorizontalPodAutoscaler](../../implmenting/evaluatemetrics)
for detailed information.

TODO: Link to HPA subsection after that content is merged

## Learn more

To learn more about the Keptn Metrics Server, see:

* Architecture:
  [Keptn Metrics Operator](../../concepts/architecture/components/metrics-operator/)
* More information about implementing Keptn Metrics:
  [Keptn Metrics](../../implementing/evaluatemetrics.md/)
