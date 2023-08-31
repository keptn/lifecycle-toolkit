---
title: Kubernetes cluster
description: Bring or install a Kubernetes cluster 
layout: quickstart
weight: 25
hidechildren: false # this flag hides all sub-pages in the sidebar-multicard.html
---

The Keptn Lifecycle Toolkit is meant to be installed
into an existing Kubernetes cluster
that runs your deployment software.
See [Requirements](reqs.md) for information about supported releases
and advice about resources required.

## Create local Kubernetes cluster

You can also create a local cluster using packages such as
[KinD](https://kind.sigs.k8s.io/),
[k3d](https://k3d.io/),
[k3s](https://k3s.io/),
and [Minikube](https://minikube.sigs.k8s.io/docs/)
to set up a local, lightweight Kubernetes cluster
where you can install the Keptn Lifecycle Toolkit
for personal study, demonstrations, and testing.
For more information, see the Kubernetes
[Install Tools](https://kubernetes.io/docs/tasks/tools/)
documentation.

The [Keptn Lifecycle Toolkit: Installation and KeptnTask Creation in Minutes](https://www.youtube.com/watch?v=Hh01bBwZ_qM)
video  demonstrates how to create a KinD cluster.
on which you can install the Lifecycle Toolkit.
The basic steps are:

1. Download, install, and run [Docker](https://docs.docker.com/get-docker/)
1. Download [KinD](https://kind.sigs.k8s.io/)
1. Create the local KinD cluster with the following command:

   ```shell
   kind create cluster
   ```

   See the
   [KinD Quick Start Guide](https://kind.sigs.k8s.io/docs/user/quick-start/)
   for more information

1. When the cluster has been created,
   run the following to verify that the cluster is working
   and that it is running a supported version of Kubernetes
   with the following command:

   ```shell
   kubectl version --short
   ```

## Prepare your cluster for KLT

The Keptn Lifecycle Toolkit installs into an existing Kubernetes cluster.
When setting up a local Kubernetes cluster
to study or demonstrate the Lifecycle Toolkit,
you need to provide these components.

Your cluster should include the following:

* A supported version of Kubernetes.
  See [Supported Kubernetes versions](reqs.md/#supported-kubernetes-versions)
  for details.

* [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl)

* Metric provider such as
  [Prometheus](https://prometheus.io/),
  [Dynatrace](https://www.dynatrace.com/),
  or [Datadog](https://www.datadoghq.com/).
  This is used for the metrics used for the observability features
  as well as the pre- and post-deployment evaluations.

* Deployment tools of your choice,
  such as
  [Argo CD](https://argo-cd.readthedocs.io/en/stable/) or
  [Flux](https://fluxcd.io/).
  Alternatively, KLT also works with just `kubctl apply` for deployment.

* If you want to use the standardized observability feature,
  you must have an OpenTelemetry collector
  and a Prometheus operator installed on your cluster.

  If you want a dashboard for reviewing metrics and traces,
  install Grafana or the dashboard of your choice.

  For traces, install Jaeger or a similar tool.

  For more information, see
  [Requirements for Open Telemetry](../implementing/otel.md/#requirements-for-opentelemetry).

Also note that the Keptn Lifecycle Toolkit includes
a light-weight cert-manager that, by default, is installed
as part of the KLT software.
If you are using another cert-manager in the cluster,
you can configure KLT to instead use your cert-manager.
See [Use Keptn with cert-manager.io](../operate/cert-manager.md)
for detailed instructions.

## How many namespaces?
  
You have significant flexibility to decide how many namespaces to use
and how to set them up.
See the Kubernetes
[Namespace](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/)
documentation for some basic information.
You can also search and find lots of "Best Practices for Namespaces"
documents published on the web.

Some considerations for KLT:
  
* KLT primarily operates on Kubernetes
  [Workload](https://kubernetes.io/docs/concepts/workloads/)
  resources and
  [KeptnApp](../yaml-crd-ref/app.md)
   resources
  that are activated and defined by annotations to each Workload.
* [KeptnMetricsProvider](../yaml-crd-ref/metricsprovider.md)
  resources need to be located
  in the same namespace as the associated
  [KeptnMetric](../yaml-crd-ref/metric.md)
  resources.
  But
  [KeptnEvaluationDefinition](../yaml-crd-ref/evaluationdefinition.md)
  resources that are used for pre- and post-deployment
  can reference metrics from any namespace.
  So you can create `KeptnMetrics` in a centralized namespace
  (such as `keptn-lifecycle-toolkit`)
  and access those metrics in evaluations on all namespaces in the cluster.
* Each `KeptnApp` resource identifies the namespace to which it belongs.
  If you configure multiple namespaces,
  you can have `KeptnApp` resources with the same name
  in multiple namespaces without having them conflict.
* You do not need separate namespaces for separate versions of your application.
  The `KeptnApp` resource includes fields to define
  the `version` as well as a `revision`
  (used if you have to rerun a deployment
  but want to retain the version number).

So, possible namespace designs run the gamut:

* Run all your KLT work in a single namespace
* Create a separate namespace for each logical grouping of your KLT work
* Create a separate namespace for each workload
