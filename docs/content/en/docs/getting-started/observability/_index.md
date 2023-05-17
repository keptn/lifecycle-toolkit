---
title: Standardize access to observability data
description: Learn how the Keptn Lifecycle Toolkit provides observability for Kubernetes deployments
weight: 45
---

The Keptn Lifecycle Toolkit (KLT) makes any Kubernetes deployment observable.
You can readily see why a deployment takes so long or why it fails,
even when using multiple deployment tools.
Keptn provides a `KeptnApplication` resource
which is an abstraction that connects multiple
Workloads belonging together.
In other words, KLT, creates a distributed end-to-end trace
of everything Kubernetes does in the context of a Deployment.

The observability data is an amalgamation of the following:

- DORA metrics are collected out of the box
  when the Lifecycle Toolkit is enabled
- OpenTelemetry runs traces that show everything that happens in the Kubernetes pod scheduler
  and can display this information with dashboard tools
  such as Grafana.
- Specific metrics that you can define to monitor
  information from all the data providers configured in your cluster.

The Keptn Lifecycle Toolkit can provide this information
for all deployments running in your cluster,
even if they are using different deployment tools.
And it can capture metrics from multiple data sources
using multiple data platforms.
With KLT implmented on your cluster,
you can easily monitor what is happening during a deployment into your Kuberenetes cluster,
and quickly get data to help you understand issues such as
why a deployment took so long or why it failed.

## Using this exercise

This exercise shows how to standardize access
to the observability data for your cluster.
It is based on the
[simplenode-dev](https://github.com/keptn-sandbox/klt-on-k3s-with-argocd)
example.
You can clone that repo to access it locally
or just look at it for examples
as you implement the functionality "from scratch"
on your local Kubernetes deployment cluster.
The
[README](https://github.com/keptn-sandbox/klt-on-k3s-with-argocd/blob/main/setup/observability/README.md)
file for that repo contains useful information.

Two videos are available
to walk you through this exercise if you prefer:

* [Introducing Keptn Lifecycle Toolkit](https://youtu.be/449HAFYkUlY)
* [Use SLOs and get DORA the Native K8s way!](https://www.youtube.com/watch?v=zeEC0475SOU)

In the
[Getting started with Keptn metrics](../metrics)
exercise, you learn how to define and use Keptn metrics.
You may want to complete that exercise before doing this exercise
although that is not required.

This exercise shows how to standardize access
to the observability data for your cluster.
The steps are:

1. Install and configure the Lifecycle Toolkit on your cluster
1. Configure OpenTelemetry for the Lifecycle Toolkit
1. Create a Keptn application that aggregates your deployment workloads
1. DORA metrics
1. Using OpenTelemetry data
1. Keptn metrics

TODO: Make the list items links to the subsections.

- Increment the version number for either your Workload
  or your application to start aggragating data.
- View the aggregated metrics and traces as text.
  You can also use Grafana or the dashboard of your choice
  to view this information.
- If you like, define `KeptnMetrics` for additional data you want to monitor.

## Install and configure KLT

The steps to install and configure KLT for this exercise are:

* [Install KLT](#install-klt)
* [Enable KLT on your cluster](#enable-klt-on-your-cluster)
* [Annotate Deployment resource](#enable-klt-on-your-cluster)

### Install KLT

To install the Keptn Lifecycle Toolkit on your cluster:

   - Be sure that your cluster includes the components discussed in
     [Prepare your cluster for KLT](../../install/k8s.md/#prepare-your-cluster-for-klt)
   - [Install and enable the Keptn Lifecycle Toolkit](../../install/install.md/#use-helm-chart)
     on your cluster using the Helm chart

If you installed KLT on your cluster for the
[Getting started with Keptn metrics](../metrics)
exercise, you do not need to re-install it for this exercise.
However, if you only installed the `metrics-operator` for that exercise,
you now need to install the full KLT.

### Enable KLT on your cluster

Follow the instructions in
[Enable KLT for your cluster](../../install/install.md/#enable-klt-for-your-cluster)
to enable KLT on your cluster.

The
[simplenode-dev-ns.yaml](https://github.com/keptn-sandbox/klt-on-k3s-with-argocd/blob/main/simplenode-dev/simplenode-dev-ns.yaml)
file defines the `Namespace` resource for our example.

### Verify that your cluster is ready

Run the following command to ensure that your Kuberetes cluster
is ready to complete this exercise by running:

```shell
kubectl get pods -n keptn-lifecycle-toolkit-system
```

You should see pods for the following components:
- certificate-operator (or another cert manager)
- lifecycle-operator
- scheduler
- otel-collector
- metrics-operator

### Annotate Deployment resource

Follow the instructions in
[Annotate workload](../../implementing/integrate/#basic-annotations)
to apply basic annotations to your `Deployment` resource.

The
[simplenode-dev-deployment.yaml](https://github.com/keptn-sandbox/klt-on-k3s-with-argocd/blob/main/simplenode-dev/simplenode-dev-deployment.yaml/)
file defines the `Deployment` resource for our example.

The example file also includes annotations for
pre- and post-deployment activities.
We will discuss those in a separate exercise.

## Define the Keptn application

A Keptn application defines the workloads
to be included in your Keptn Application.
We will use the application discovery feature
to automatically generate a Keptn Application
that includes all workloads on the cluster,
regardless of the tools being used.

A Keptn application aggregates multiple workloads
that belong to a logical app into a single
You can create the YAML file to define the
[KeptnApp](../../yaml-crd-ref/app.md)
resource manually
but the easier approach is to let KLT create this definition for you.
This requires that you annotate all your workloads
(`Deployments`, `Stateful Sets`, `DaemonSets`, and `ReplicaSets`
as described in
[Use Keptn automatic app discovery](../../implementing/integrate.md/#use-keptn-automatic-app-discovery).

## DORA metrics

DORA metrics are an industry-standard set of measurements;
see the following for a description:

* [What are DORA Metrics and Why Do They Matter?](https://codeclimate.com/blog/dora-metrics)
* [Are you an Elite DevOps Performer?  Find out with the Four Keys Project](https://cloud.google.com/blog/products/devops-sre/using-the-four-keys-to-measure-your-devops-performance)

DORA metrics provide information such as:

* How many deployments happened in the last six hours?
* Time between deployments
* Deployment time between versions
* Average time between versions.

The Keptn Lifecycle Toolkit starts collecting these metrics
as soon as you annotate the `Deployment` resource.
Metrics are collected only for the `Deployment` resources
that are annotated.

To view DORA metrics, run the following command:

```shell
kubectl port-forward -n keptn-lifecycle-toolkit-system \
   svc/lifecycle-operator-metrics-service 2222
```

Then view the metrics at:

```shell
http://localhost:2222/metrics
```

DORA metrics are also displayed on Grafana
or whatever dashboard application you choose.
See the videos that accompany this exercise
for examples.

## Using OpenTelemetry

The Keptn Lifecycle Toolkit extends the Kubernetes
primitives to create OpenTelemetry data
that connects all your deployment and observability tools
without worrying about where it is stored and where it is managed.
OpenTelemetry traces collect data as Kubernetes is deploying the changes,
which allows you to trace everything done in the context of that deployment.

* You must have an OpenTelemetry collector installed on your cluster.
  See
  [OpenTelemetry Collector](https://opentelemetry.io/docs/collector/)
  for more information.
* Follow the instructions in
  [OpenTelemetry observability](../../implementing/otel.md)
  to configure where your OpenTelemetry data is sent.
  * Define a [KeptnConfig](../../yaml-crd-ref/config.md) resource
  that defines the URL and port of the OpenTelemetry collector.
  For our example, this is in the
  [keptnconfig.yaml](https://github.com/keptn-sandbox/klt-on-k3s-with-argocd/blob/main/setup/keptn/keptnconfig.yaml)
  file.

## Keptn metrics

You can supplement the DORA Metrics and OpenTelemetry information
with information you explicitly define using Keptn metrics.
The
[Getting started with Keptn metrics](../metrics)
exercise discusses how to define Keptn metrics.

## View the results

To start feeding observability data for your deployments
onto a dashboard of your choice,
modify either your `Deployment` or `KeptnApp` resource yaml file
to increment the version number
and commit that change to your repository.

TODO: How to view DORA and OpenTelemetry data without a dashboard.
      Are the DORA metrics viewed using the same commands shown for
      Keptn metrics.

TODO: talk about the Grafana display.

The videos that go with this exercise show how the
DORA, OpenTelemetry, and Keptn metrics information
appears on a Grafana dashboard with
[Jaeger](https://grafana.com/docs/grafana-cloud/data-configuration/metrics/prometheus-config-examples/the-jaeger-authors-jaeger/).

If you also have Jaeger extension for Grafana installed on your cluster,
you can view full end-to-end trace for everything
that happens in your deployment.
For more information, see
