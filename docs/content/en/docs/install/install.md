---
title: Install Keptn
description: Install and enable Keptn
weight: 35
hidechildren: false # this flag hides all sub-pages in the sidebar-multicard.html
---

Keptn must be installed onto each Kubernetes cluster you want to monitor.
Additionally, Keptn needs to be enabled on your namespaces.
This gives you flexibility in how and where you want to use Keptn.

Keptn v0.9.0 and later is installed using [Helm](https://helm.sh/).

> **Note** Earlier releases could also be installed using the manifest.
> See
[Upgrade to Helm from a manifest installation](upgrade.md/#upgrade-to-helm-from-a-manifest-installation)
> if you need to upgrade from a manifest installation.

After you install Keptn, you are ready to
[Integrate Keptn with your applications](../implementing/integrate.md).

## Basic installation

Keptn is installed onto an existing Kubernetes cluster
using a Helm chart.
To modify the Keptn configuration,
you must modify the `values.yaml` file of the chart.

> **Note** Keptn works on virtually any type of Kubernetes cluster.
  See
  [Requirements](reqs.md)
  for specific requirements for your Kubernetes cluster.
>

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
  use the following command sequence:

  ```shell
  helm repo add keptn https://charts.lifecycle.keptn.sh
  helm repo update
  helm search repo keptn --versions
  ```
  
  You see that the "CHART VERSION" for `keptn/keptn` v0.9.0 is 0.3.0
  so the following command sequence explicitly installs Keptn v0.9.0:

  ```shell
  helm repo add keptn https://charts.lifecycle.keptn.sh
  helm repo update
  helm upgrade --install keptn keptn/keptn \
   --version version 0.3.0 \
   -n keptn-lifecycle-toolkit-system --create-namespace --wait
  ```

* To view which Keptn components are installed in your cluster
  and verify that they are the correct ones,
  run the following command:

  ```shell
  kubectl get pods -n keptn-system
  ```

  The output shows all Keptn components that are running on your cluster.

## Customizing the configuration of components

To modify configuration of the subcomponents of Keptn,
use the corresponding Helm values in your main `values.yaml`.
Use the subcomponent's parent value as the root for your configuration.

Here is an example `values.yaml` altering global and metrics operator values:

{{< docsembed path="content/en/docs/install/assets/values-advance-changes.yaml" >}}

Note the additional parameters that are specified
in the `metricsOperator` section.
These are documented in the README file for that operator,
which is linked from the `metrics-operator` item under "Component"
in the table in
[Control what components are installed](install/#control-what-components-are-installed).

### Modify Helm configuration options

Helm values can be modified before the installation.

To modify configuration options, download a copy of the
default `values.yaml` file,
modify some values, and use the modified file to install Keptn:

1. Download the `values.yaml` file:

   ```shell
   helm show values keptn/keptn > values.yaml
   ```

1. Edit your local copy to modify some values

1. To apply the modified configuration.
   install Keptn by adding the following string
   to your `helm upgrade` command line:

   ```shell
   --values=values.yaml
   ```

   For example, if you create a `my.values.yaml`
   and modified some configuration values,
   the command sequence to apply your configuration is:

   ```shell
   helm repo add keptn https://charts.lifecycle.keptn.sh
   helm repo update
   helm upgrade --install keptn keptn/keptn \
    --values my.values.yaml \
    -n keptn-lifecycle-toolkit-system --create-namespace --wait
   ```

   You can also use the `--set` flag
   to specify a value change for the `helm upgrade --install` command.
   Configuration options are specified using the format:

   ```shell
   --set key1=value1 \
   --set key2=value2 ...
   ```

## Control what components are installed

Keptn consists of multiple components,
each of which enables a specific use-case.
Each of the components is packaged into a subchart of Keptn
which has parameters that can be configured for that component
in the umbrella `values.yaml` file.

The following table summarizes the Keptn umbrella chart scheme;
click the name in the "Component" column
to view the README file that documents the parameters
that can be set for that component:

| Component                                                                                                                  | Used for                                                                                                                  | Configuration file |
|----------------------------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------------------------| --------------------|
| [Keptn](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn/README.md)                           | Installs subcharts, global configuration                                                                                  | [keptn/values.yaml](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn/values.yaml) |
| [lifecycle-operator](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn-lifecycle-operator/README.md) | [Observability](../implementing/otel.md), [Release Lifecycle Management](../intro/_index.md#release-lifecycle-management) | [keptn-lifecycle-operator/values.yaml](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn-lifecycle-operator/values.yaml) |
| [metrics-operator](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn-metrics-operator/README.md)    | [Keptn metrics](../implementing/evaluatemetrics.md), [Analysis](../implementing/slo.md)                                   | [keptn-metrics-operator/values.yaml](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn-metrics-operator/values.yaml) |
| [cert-manager](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn-cert-manager/README.md)            | TLS Certificate management for all Keptn components                                                                       | [keptn-cert-manager/values.yaml](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn-cert-manager/values.yaml) |

By default, all components are included when you install Keptn.
To specify the components that are included,
you need to modify the Keptn `values.yaml` file.

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

{{< docsembed path="content/en/docs/install/assets/values-remove-certmanager.yaml" >}}

For more information on using `cert-manager` with Keptn, see
[Use Keptn with cert-manager.io](../operate/cert-manager.md).

For the full list of Helm values, see the
[keptn-cert-manager Helm chart README](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn-cert-manager/README.md).
