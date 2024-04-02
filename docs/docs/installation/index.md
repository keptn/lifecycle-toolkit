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
refer to the [Deploying Keptn via ArgoCD](./configuration/argocd.md) section
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

## Keptn Helm configuration

The Keptn configuration is controlled by a set of Helm value files,
summarized in the following table.
The Keptn Helm chart is an umbrella chart
that contains subcharts for all components of Keptn.
Each component has its own Helm values file
(documented in its own README file),
that defines configurations specific to that component.

All configuration changes for all components
can be made in one `values.yaml` file.
This is discussed more in
[Customizing the configuration of components](#customizing-the-configuration-of-components)
below.

The following table summarizes the Keptn `values.yaml` files.

* The "Component" column leads you to the
  README files for each component where
  all Helm values are documented
* The "Configuration file" column leads you to
  the Helm values files for each component

| Component                                                                                                                  | Used for                                                                                                                    | Configuration file |
|----------------------------------------------------------------------------------------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------| --------------------|
| [Keptn](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn/README.md)                           | Installs subcharts, global configuration                                                                                    | [keptn/values.yaml](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn/values.yaml) |
| [lifecycle-operator](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn-lifecycle-operator/README.md) | [Observability](../guides/otel.md), [Release Lifecycle Management](../core-concepts/index.md#release-lifecycle-management) | [keptn-lifecycle-operator/values.yaml](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn-lifecycle-operator/values.yaml) |
| [metrics-operator](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn-metrics-operator/README.md)    | [Keptn metrics](../guides/evaluatemetrics.md), [Analysis](../guides/slo.md)                                                 | [keptn-metrics-operator/values.yaml](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn-metrics-operator/values.yaml) |
| [cert-manager](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn-cert-manager/README.md)            | [TLS Certificate management for all Keptn components](../components/certificate-operator.md)                         | [keptn-cert-manager/values.yaml](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn-cert-manager/values.yaml) |

## Customizing the configuration of components

To modify configuration of the subcomponents of Keptn,
use the corresponding Helm values in your `values.yaml` file.
Use the subcomponent's parent value as the root for your configuration.

Here is an example `values.yaml` altering global and metrics operator values:

```yaml
{% include "./assets/values-advance-changes.yaml" %}
```

Note the additional values that are specified
in the `metricsOperator` section.
These are documented in the README file for that operator,
which is linked from the `metrics-operator` item under "Component"
in the table above.
To implement this:

* Go into the `values.yaml` file linked under "Configuration file"
* Copy the lines for the values you want to modify
* Paste those lines into your `values.yaml` file
  and modify their values in that file.

### Modify Helm values

To modify Helm values:

1. Download a copy of the Helm values file:

    ```shell
    helm show values keptn/keptn > values.yaml
    ```

1. Edit your local copy to modify some values

1. Add the following string
   to your `helm upgrade` command to install Keptn
   with your configuration changes:

    ```shell
    --values=values.yaml
    ```

    For example, if you create a `my.values.yaml`
    and modify some configuration values,
    use the following command to apply your configuration:

    ```shell
    helm upgrade --install keptn keptn/keptn \
    --values my.values.yaml \
    -n keptn-system --create-namespace --wait
    ```

    You can also use the `--set` flag
    to specify a value change for the `helm upgrade --install` command.
    Helm values are specified using the format:

    ```shell
    --set key1=value1 \
    --set key2=value2 ...
    ```

## Control what components are installed

By default, all components are included when you install Keptn.
To specify the components that are included,
you need to modify the Keptn `values.yaml` file
to disable the components you do not want to install.

Note that the Keptn Helm chart is quite flexible.
You can install all Keptn components on your cluster,
then modify the configuration to exclude some components
and update your installation.
Conversely, you can exclude some components when you install Keptn
then later add them in.

### Disable Keptn Certificate Manager (Certificates)

If you wish to use your custom certificate manager,
you can disable Keptn `cert-manager` by setting the
`certificateManager.enabled` Helm value to `false`:

```yaml
{% include "./assets/values-remove-certmanager.yaml" %}
```

For more information on using `cert-manager` with Keptn, see
[Use Keptn with cert-manager.io](../components/certificate-operator.md).

For the full list of Helm values, see the
[keptn-cert-manager Helm chart README](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn-cert-manager/README.md).

## Uninstalling Keptn

To uninstall Keptn from your cluster, please follow the steps
on the [Uninstall page](./uninstall.md).
