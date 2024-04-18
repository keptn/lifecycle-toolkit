---
comments: true
---

# Installation

Keptn must be installed onto each Kubernetes cluster you want to monitor.
Additionally, Keptn needs to be enabled on your namespaces.
This gives you flexibility in how and where you want to use Keptn.

> **Note** By default, Keptn monitors all namespaces in the cluster
> except for those reserved for major components.
> See
> [Namespaces and Keptn](configuration/namespace-keptn.md)
> to learn how to limit the namespaces that Keptn monitors.

Keptn is installed using [Helm](https://helm.sh/).

After you install Keptn, you are ready to
[Integrate Keptn with your applications](../guides/integrate.md).

## Supported Kubernetes versions

Keptn requires Kubernetes v1.24.0 or later.

Run the following to ensure that both client and server versions
are running Kubernetes versions greater than or equal to v1.24.
In this example, both client and server are at v1.24.0
so Keptn will work.

```shell
kubectl version --short
```

```shell
Client Version: v1.24.0
Kustomize Version: v4.5.4
Server Version: v1.24.0
```

Keptn makes use of a custom scheduler
when running on Kubernetes v1.26 and earlier.
For Kubernetes v1.27 and later, scheduling is
implemented using
[Kubernetes scheduling gates](https://kubernetes.io/docs/concepts/scheduling-eviction/pod-scheduling-readiness/),
unless the `schedulingGatesEnabled` Helm value is set to `false`.
See
[Keptn integration with Scheduling](../components/scheduling.md)
for details.

If Keptn is installed on a [vCluster](https://www.vcluster.com/) with
Kubernetes v1.26 or earlier, some extra configuration
needs to be added for full compatibility.
See
[Running Keptn with vCluster](./configuration/vcluster.md)
for more information.

If you want to deploy Keptn via [ArgoCD](https://argoproj.github.io/cd/),
refer to the [Deploying Keptn via ArgoCD](./alternative-installs/argocd.md) section
for more information.

## Basic installation

Keptn is installed onto an existing Kubernetes cluster
using a Helm chart.
To modify the Keptn configuration,
you must modify the `values.yaml` file of the chart.

The command sequence to fetch and install the latest release of Keptn is:

```shell
helm repo add keptn https://charts.lifecycle.keptn.sh
helm repo update
helm upgrade --install keptn keptn/keptn \
   -n keptn-system --create-namespace --wait
```

If you want to use Keptn to observe your deployments
or to enhance them with lifecycle management
you must enable Keptn in your
[Namespace](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/)
via annotations.
For example, for the `testy-test` namespace:

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: testy-test
  annotations:
    keptn.sh/lifecycle-toolkit: "enabled" # this tells Keptn to watch the namespace
```

Some helpful hints:

* Use the `--version <version>` flag on the
  `helm upgrade --install` command to specify a different Keptn Helm chart version.

    To get the appropriate chart version for the Keptn version you want,
    use the following command:

    ```shell
    helm search repo keptn --versions
    ```

    You see that the "CHART VERSION" for `keptn/keptn` v0.9.0 is 0.3.0
    so use the following command to explicitly installs Keptn v0.9.0:

    ```shell
    helm upgrade --install keptn keptn/keptn --version 0.3.0 \
    -n keptn-system --create-namespace --wait
    ```

* To view which Keptn components are installed in your cluster
  and verify that they are the correct ones,
  run the following command:

    ```shell
    kubectl get pods -n keptn-system
    ```

    The output shows all Keptn components that are running on your cluster.
