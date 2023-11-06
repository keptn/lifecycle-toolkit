---
title: Install and enable Keptn
description: Install Keptn
weight: 35
hidechildren: false # this flag hides all sub-pages in the sidebar-multicard.html
---

Keptn must be installed and integrated
into each Kubernetes cluster you want to monitor.
Keptn v0.9.0 and later is installed using
an umbrella [Helm Chart](#basic-installation).
This means that the Helm Chart that installs all of Keptn
actually groups subcharts for individual services
and you can install one of these services
without installing all of Keptn.

> **Note** Earlier releases could also be installed using the manifest.
> See
[Upgrade to Helm from a manifest installation](upgrade.md/#upgrade-to-helm-from-a-manifest-installation)
> if you need to upgrade from a manifest installation.

This page covers the following:

* [Basic installation](#basic-installation) discusses the command sequence
  used to install Keptn.

* To modify the Keptn configuration,
  you must modify the appropriate `values.yaml` files:

  * [Modify Helm configuration options](#modify-helm-configuration-options)
    summarizes the mechanics of modifying and applying a `values.yaml` file.
  * [Control what components are installed](#control-what-components-are-installed)
    explains how to modify the umbrella `values.yaml` file
    to control which components are included in your Keptn configuration.
    By default, all components are installed.

    This section then gives explicit notes about the components that must be configured
    for common use-cases.
  * [Customizing the configuration of components](#customizing-the-configuration-of-components)
    explains how to modify the configuration of individual components.

* [Running Keptn with vCluster](#running-keptn-with-vcluster)
  discusses how to configure vCluster in order to run Keptn on it.

After you install Keptn, you are ready to
[Integrate Keptn with your applications](../implementing/integrate.md).

## Basic installation

Keptn is installed onto an existing Kubernetes cluster
using an umbrella Helm Chart and standard
[helm](https://helm.sh/docs/intro/cheatsheet/)
commands.

> **Note** Keptn works on virtually any type of Kubernetes cluster
  with some standard tools installed.
  See
  [Kubernetes cluster](k8s.md)
  for details about preparing your Kubernetes cluster for Keptn.
>

The command sequence to fetch and install the latest release of Keptn is:

```shell
helm repo add keptn https://charts.lifecycle.keptn.sh
helm repo update
helm upgrade --install keptn keptn/keptn \
   -n keptn-lifecycle-toolkit-system --create-namespace --wait
```

Note that the `helm repo update` command is used for fresh installs
as well as for upgrades.

Some helpful hints:

* Use the `--version <version>` flag on the
  `helm upgrade --install` command to install a Keptn version
  other than the latest.
  You must specify the Chart Version of `keptn/keptn` here,
  not the actual Keptn release number.

  To get the appropriate chart version for the Keptn version you want,
  use the following command sequence:

    ```shell
    helm repo add keptn https://charts.lifecycle.keptn.sh
    helm repo update
    helm search repo keptn
    ```

  The output should look like this:
  ```
  NAME                            CHART VERSION   APP VERSION     DESCRIPTION                                       
  keptn/keptn                     0.3.0           v0.9.0          A Helm chart for Keptn Lifecycle Toolkit, a set...
  keptn/keptn-cert-manager        0.2.0           v1.2.0          A Helm chart for Keptn Certificate Manager, a s...
  keptn/keptn-lifecycle-operator  0.1.0           v0.8.3          A Helm chart for Keptn Lifecycle Operator, a su...
  keptn/keptn-metrics-operator    0.1.0           v0.8.3          A Helm chart for Keptn Metrics Operator, a subp...
  keptn/common                    0.1.0                           A Helm chart containing common templates for Keptn
  keptn/klt                       0.2.6           v0.8.2          A Helm chart for Keptn Lifecycle Toolkit, a set...
  ```

  You see that the "CHART VERSION" for `keptn/keptn`
  for v0.9.0 is 0.3.0
  so the following command sequence explicitly installs Keptn v0.9.0.

  ```shell
  helm repo add keptn https://charts.lifecycle.keptn.sh
  helm repo update
  helm upgrade --install keptn keptn/keptn \
   --version version 0.3.0 \
   -n keptn-lifecycle-toolkit-system --create-namespace --wait
  ```

* If you modify the `values.yaml` file for the umbrella chart
  and/or the component sub-charts,
  apply the modified configuration with the `-- values values.yaml` flag.
  For example, if you created a `my.values.yaml`
  to exclude the `lifecycle-operator` from your cluster
  and created a `metrics.values.yaml` to modify the configuration
  of the `metrics-operator`, the command sequence is:

  ```shell
  helm repo add keptn https://charts.lifecycle.keptn.sh
  helm repo update
  helm upgrade --install keptn keptn/keptn \
   --values my.values.yaml --values metrics.values.yaml \
   -n keptn-lifecycle-toolkit-system --create-namespace --wait
  ```

  You see that you can specify multiple `--values` options in one command.
  The precedence is from right to left.

* To view which Keptn components are installed in your cluster
  and verify that they are the correct ones,
  run the following command:

  ```shell
  kubectl get pods -n keptn-lifecycle-toolkit-system
  ```

## Modify Helm configuration options

To modify the Keptn configuration,
you must modify the appropriate `values.yaml` file.
Values can be modified before the installation
or can be modified and reapplied after the initial installation.
This is useful if you do not want to install all components,
if you need to modify the configuration in some way,
such as to turn verbose logging on or off.

To modify configuration options,
download a copy of the appropriate `values.yaml` file
(see the table in the
[Control what components are installed](#control-what-components-are-installed)
section below).
Edit the local copy of that `value.yaml` file to make the modifications,
and use the modified file to install Keptn:

1. Download the `values.yaml` file:

   ```shell
   helm get values RELEASE_NAME [flags] > values.yaml
   ```

1. Edit your local copy to modify the appropriate parameters.
   The `README` file for each component
   (which is linked to the text in the "Components" column in the table)
   describes each parameter in the `values.yaml` file.

1. Upgrade your Keptn installation to use the new configuration
   by adding the following string to your `helm upgrade` command:

   ```shell
   --values=values.yaml
   ```

You can also use the `--set` flag
to specify a value change for the `helm upgrade --install` command.
Configuration options are specified using the format:

```shell
--set key1=value1,key2=value2,....
```

## Control what components are installed

Keptn presents a toolkit composed of multiple operators,
each of which enables a specific use-case.
Keptn v.0.9 and later releases use an umbrella chart
with separate charts for each operator.

The installation command is always the same;
only the content of the `values.yaml` file(s) changes:

* Use the Keptn umbrella `values.yaml` file to specify which components are installed
  and to set some global configuration parameters for Keptn
* Each component has its own `values.yaml` file
  where you set configuration parameters that are specific to that component.

The following table summarizes the Keptn umbrella chart scheme;
click the name in the "Components" column
to view the README file that documents the parameters used
in the `values.yaml` file for that component:

| Component | Used for | Configuration file |
| --------- | -------- | --------------------|
| [Keptn umbrella](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn/README.md) | Installs subcharts, global configuration | [keptn/values.yaml](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn/values.yaml) |
| [lifecycle-operator](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn-lifecycle-operator/README.md) | [Observability](../implementing/otel.md), Lifecycle management [tasks](../implementing/tasks.md) and [evaluations](../implementing/evaluatemetrics.md) | [keptn-lifecycle-operator/values.yaml](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn/values.yaml) |
| [metrics-operator](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn-metrics-operator/README.md) | [Keptn metrics](../implementing/evaluatemetrics.md), [Analysis](../implementing/slo.md) | [keptn-metrics-operator/values.yaml](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn-metrics-operator/values.yaml) |
| [cert-manager](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn-cert-manager/README.md)  | Configures TLS certificates | [keptn-cert-manager/values.yaml](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn-cert-manager/values.yaml) |

By default, all components are included when you install Keptn.
To specify the components that are included,
you need to modify the
[keptn/values.yaml](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn/values.yaml)
file.
The following sections summarize the configurations needed for different use cases.

Note that the umbrella scheme is quite flexible.
You can install all of Keptn on your cluster,
then modify the configuration to exclude some components
and update your configuration.
Conversely, you can exclude some components when you install Keptn
then later add them in.

### Enable Keptn Lifecycle Operator (Observability and/or Release Lifecycle Management)

If you only want to run the Keptn Observability
and/or Release Lifecycle Management use-cases in your cluster,
you do not need to install the Keptn Metrics Operator.
To disable it, set the `metricsOperator.enabled` value
to `false` as in the following:

{{< embed path="/docs/content/en/docs/install/assets/values-only-lifecycle.yaml" >}}

Note that, if you want to run pre- and/or post-deployment
[evaluations](../implementing/evaluations.md)
as part of the Release Lifecycle Management use-case,
you need to run the Keptn Metrics Operator.

You must also enable Keptn for each
[namespace](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/)
on which you want to run either
the Observability or Release Lifecycle Management use-case.
This is because Keptn communicates with the Kubernetes scheduling mechanism
for tasks such as enforcing checks natively,
stopping a deployment from proceeding when criteria are not met,
doing post-deployment evaluations
and tracing all activities of all deployment
[workloads](https://kubernetes.io/docs/concepts/workloads/)
on the cluster.

To enable Keptn in your namespace,
annotate the appropriate `Namespace` resource(s).
For example, for the `simplenode-dev` namespace,
the annotation looks like this:

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: simplenode-dev
  annotations:
    keptn.sh/lifecycle-toolkit: "enabled"
```

You see the annotation line `keptn.sh/lifecycle-toolkit: "enabled"`.
This annotation tells the webhook to handle the namespace.

After enabling Keptn for your namespace(s),
you are ready to
[Integrate Keptn with your applications](../implementing/integrate.md).

For more information about implementing Observability, see the
[Observability User Guide](../implementing/otel.md).

For more information about implementing Keptn Release Management, see the
[Deployment tasks](../implementing/tasks.md)
and
[Evaluations](../implementing/evaluations.md)
User Guides.

### Enable Keptn Metrics Operator (Metrics)

If you are only interested in Metrics,
you do not need the Keptn Lifecycle Operator.
Disable it using the following values.yaml:

{{< embed path="/docs/content/en/docs/install/assets/values-only-metrics.yaml" >}}

For more information about implementing Metrics, see the
[Metrics User Guide](../implementing/evaluatemetrics.md).

### Enable Keptn Analysis (SLOs/SLIs)

To enable Keptn Analysis in your cluster,
you again do not need the Keptn Lifcycle Operator.
Disable it using the following values.yaml:

{{< embed path="/docs/content/en/docs/install/assets/values-only-metrics.yaml" >}}

> **Note** The Analysis use-case is currently behind a feature flag.
  To enable it, add the following to your `helm upgrade` command line:

  ```shell
  --set metricsOperator.env.enableKeptnAnalysis=true
  ```
>

For more information about implementing Keptn Analysis, see the
[Analysis User Guide](../implementing/slo.md).

### Disable Keptn Certificate Manager (Certificates)

If you wish to use your own custom certificate manager,
you can disable the Keptn `cert-manager` by using the
`--set " certificateManager.enabled=false"` argument
to the `helm upgrade` command line
or you can modify the `keptn/values.yaml` file:

{{< embed path="/docs/content/en/docs/install/assets/values-remove-certmanager.yaml" >}}

For more information about using `cert-manager` with Keptn, see
[Use Keptn with cert-manager.io](../operate/cert-manager.md).

## Customizing the configuration of components

To access and modify the configuration of a subcomponent,
modify the component's `values.yaml` file.
You can use the subchart name
as written in the `chart.yaml` file.
The `README` file for each component documents the parameters
that are supported for that component.

Here is an example `values.yaml` that alters values for the metrics operator:

{{< embed path="/docs/content/en/docs/install/assets/values-advance-changes.yaml" >}}

## Running Keptn with vCluster

Keptn running on Kubernetes versions 1.26 and older
uses a custom
[scheduler](../architecture/components/scheduler.md),
so, out of the box, it does not work with
[Virtual Kubernetes Clusters](https://www.vcluster.com/)
("vClusters").
This is also an issue
if the `schedulingGatesEnabled` Helm chart value is set to `false`
for Kubernetes version 1.27 and later.
See
[Keptn integration with Scheduling](../architecture/components/scheduler.md)
for details.

To solve this problem:

1. Follow the instructions in
   [Separate vCluster Scheduler](https://www.vcluster.com/docs/architecture/scheduling#separate-vcluster-scheduler)
   to modify the vCluster `values.yaml` file
   to use a virtual scheduler.

1. Create or upgrade the vCluster,
   following the instructions in that same document.

1. Follow the instructions in the section below
   to install Keptn in that vCluster.
