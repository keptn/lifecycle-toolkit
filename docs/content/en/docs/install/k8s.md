---
title: Kubernetes cluster
description: Bring or install a Kubernetes cluster 
icon: concepts
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
  This is used for the metrics used for the observability features.

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
  [Requirements for Open Telemetry](../implementing/otel/#requirements-for-opentelemetry).

Also note that the Keptn Lifecycle Toolkit includes
a light-weight cert-manager that, by default, is installed
as part of the KLT software.
If you are using another cert-manager in the cluster,
you can configure KLT to instead use your cert-manager.
See [Use your own cert-manager](cert-manager.md)
for detailed instructions.
