---
comments: true
---

# Kubernetes cluster

Keptn is meant to be installed
into an existing Kubernetes cluster
that runs your deployment software.
See [Requirements](index.md#supported-kubernetes-versions) for information about supported releases
and advice about resources required.

## Create local Kubernetes cluster

You can also create a local cluster using packages such as
[KinD](https://kind.sigs.k8s.io/),
[k3d](https://k3d.io/),
[k3s](https://k3s.io/),
and [Minikube](https://minikube.sigs.k8s.io/docs/)
to set up a local, lightweight Kubernetes cluster
where you can install Keptn
for personal study, demonstrations, and testing.
For more information, see the Kubernetes
[Install Tools](https://kubernetes.io/docs/tasks/tools/)
documentation.

The [Keptn: Installation and KeptnTask Creation in Minutes](https://www.youtube.com/watch?v=Hh01bBwZ_qM)
video  demonstrates how to create a KinD cluster.
on which you can install Keptn.
The basic steps are:

1. Download, install, and run [Docker](https://docs.docker.com/get-docker/)
2. Download [KinD](https://kind.sigs.k8s.io/)
3. Create the local KinD cluster with the following command:

    ```shell
    kind create cluster
    ```

   See the
   [KinD Quick Start Guide](https://kind.sigs.k8s.io/docs/user/quick-start/)
   for more information

4. When the cluster has been created,
   run the following to verify that the cluster is working
   and that it is running a supported version of Kubernetes
   with the following command:

    ```shell
    kubectl version --short
    ```

## Prepare your cluster for Keptn

Keptn installs into an existing Kubernetes cluster.
When setting up a local Kubernetes cluster to study or demonstrate Keptn,
you need to provide these components.

Your cluster should include the following:

* A supported version of Kubernetes.
  See [Supported Kubernetes versions](index.md#supported-kubernetes-versions)
  for details.

* The
  [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl)
  CLI that is used to interact with Kubernetes clusters.

* The
  [Helm](https://helm.sh/docs/intro/install/)
  CLI that is used to install and configure Keptn.

* Deployment tools of your choice,
  such as
  [Argo CD](https://argo-cd.readthedocs.io/en/stable/) or
  [Flux](https://fluxcd.io/).
  Alternatively, Keptn also works with just `kubectl apply` for deployment.

* At least one observability data provider such as
  [Prometheus](https://prometheus.io/),
  [Thanos](https://thanos.io/),
  [Dynatrace](https://www.dynatrace.com/),
  or [Datadog](https://www.datadoghq.com/);
  you can use multiple instances of different data providers.
  These provide:

    * Metrics used for [Keptn Metrics](../guides/evaluatemetrics.md)
    * Metrics used for [OpenTelemetry](../guides/otel.md) observability
    * SLIs for pre- and post-deployment [evaluations](../guides/evaluations.md)
    * SLIs used for [analyses](../guides/slo.md)

* If you want to use the standardized observability feature,
  you must have an OpenTelemetry collector
  as well as a Prometheus operator installed on your cluster.
  For more information, see
  [Requirements for OpenTelemetry](../guides/otel.md/#requirements-for-opentelemetry).

* If you want a dashboard for reviewing metrics and traces,
  install the dashboard tools of your choice;
  we primarily use Grafana.
  For more information, see
  [Requirements for Open Telemetry](../guides/otel.md/#requirements-for-opentelemetry).

* Keptn includes a lightweight `cert-manager` that, by default,
  is installed as part of the Keptn software.
  If you are using another certificate manager in the cluster,
  you can configure Keptn to instead use your cert-manager.
  See [Use Keptn with cert-manager.io](./configuration/cert-manager.md)
  for detailed instructions.
