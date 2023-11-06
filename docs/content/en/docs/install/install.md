---
title: Install Keptn
description: Install and enable Keptn
weight: 35
hidechildren: false # this flag hides all sub-pages in the sidebar-multicard.html
---

Keptn must be installed onto each Kubernetes cluster you want to monitor.
Additionally, Keptn needs to be enabled on certain namespaces.
This is because Keptn communicates with the Kubernetes primitives
for tasks such as enforcing checks natively,
stopping a deployment from proceeding when criteria are not met,
doing post-deployment evaluations
and tracing all activities of all deployment [workloads](https://kubernetes.io/docs/concepts/workloads/) on the cluster.

Keptn v0.9.0 and later is installed using
a [Helm Chart](#basic-installation).

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

Some helpful hints:

* Use the `--version <version>` flag on the
  `helm upgrade --install` command to specify a different Keptn Helm chart version.

* Use the following command sequence to see a list of available versions:

  ```shell
  helm repo update
  helm search repo keptn --versions
  ```

* To verify that Keptn was installed in your cluster,
  run the following command:

  ```shell
  kubectl get pods -n keptn-system
  ```

  The output shows all Keptn components that are running on your cluster.

If you want to use Keptn to observe your deployments, you must enable Keptn in your
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


## Modify Helm configuration options

Helm values can be modified before the installation.

To modify configuration options, download a copy of the
default `values.yaml` file,
modify some values, and use the modified file to install Keptn:

1. Download the `values.yaml` file:

   ```shell
   helm show values keptn/keptn > values.yaml
   ```

1. Edit your local copy to modify some values

1. Install Keptn by adding the following string to your `helm upgrade` command line:

   ```shell
   --values=values.yaml
   ```

You can also use the `--set` flag
to specify a value change for the `helm upgrade --install` command.
Configuration options are specified using the format:

```shell
--set key1=value1 \
--set key2=value2 ...
```

The [helm-charts](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn/README.md) page
contains the full list of available values.

## Control what components are installed

Keptn consists of multiple operators, each of which enables a specific use-case.
Each of the operator is packaged into a subchart of Keptn which can be configured individually.

The following table summarizes the Keptn umbrella chart scheme;

| Component                                                                                                                  | Used for | Configuration file |
|----------------------------------------------------------------------------------------------------------------------------| -------- | --------------------|
| [Keptn](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn/README.md)                           | Installs subcharts, global configuration | [keptn/values.yaml](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn/values.yaml) |
| [lifecycle-operator](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn-lifecycle-operator/README.md) | [Observability](../implementing/otel.md), Lifecycle management [tasks](../implementing/tasks.md) and [evaluations](../implementing/evaluatemetrics.md) | [keptn-lifecycle-operator/values.yaml](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn-lifecycle-operator/values.yaml) |
| [metrics-operator](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn-metrics-operator/README.md)    | [Keptn metrics](../implementing/evaluatemetrics.md), [Analysis](../implementing/slo.md) | [keptn-metrics-operator/values.yaml](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn-metrics-operator/values.yaml) |
| [cert-manager](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn-cert-manager/README.md)            | Configures TLS certificates | [keptn-cert-manager/values.yaml](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn-cert-manager/values.yaml) |

By default, all components are included when you install Keptn.
To specify the components that are included,
you need to modify the Keptn `values.yaml` file.
The following sections summarize the configurations needed for different use cases.

Note that the umbrella scheme is quite flexible.
You can install all Keptn components on your cluster,
then modify the configuration to exclude some components
and update your installation.
Conversely, you can exclude some components when you install Keptn
then later add them in.

### Enable Keptn Lifecycle Operator (Observability)

If you only want to run the Keptn Observability use-case in your cluster,
you do not need to install the Keptn Metrics Operator.
To disable it, modify the `keptn/values.yaml` file like this:

{{< docsembed path="content/en/docs/install/assets/values-only-lifecycle.yaml" >}}

After enabling Keptn for your namespace(s),
you are ready to
[Integrate Keptn with your applications](../implementing/integrate.md).

For more information about implementing observability, see the
[Observability User Guide](../implementing/otel.md).

### Enable Keptn Metrics Operator (Metrics)

If you are interested in Metrics, you do not need Keptn Lifecycle Operator.
disable it using the following values.yaml:

{{< docsembed path="content/en/docs/install/assets/values-only-metrics.yaml" >}}

For more information about implementing Metrics, see the
[Metrics User Guide](../implementing/evaluatemetrics.md).

### Enable Keptn Analysis (SLOs/SLIs)

To enable Keptn Analysis in your cluster, you again do not need the Keptn Lifcycle Operator,
disable it using the following values.yaml:

{{< docsembed path="content/en/docs/install/assets/values-only-metrics.yaml" >}}

> **Note** The Analysis use-case is currently behind a feature flag.
  To enable it, add the following to your `helm upgrade` command line:

  ```shell
  --set metricsOperator.env.enableKeptnAnalysis=true
  ```
>

For more information about implementing Keptn Analysis, see the
[Analysis User Guide](../implementing/slo.md).

### Disable Keptn Certificate Manager (Certificates)

If you wish to use your custom certificate manager,
you can disable Keptn `cert-manager` by using the
`--set " certificateManager.enabled=false"` argument
to the `helm upgrade` command line
or you can modify the `keptn/values.yaml` file:

{{< docsembed path="content/en/docs/install/assets/values-remove-certmanager.yaml" >}}

For more information on using `cert-manager` with Keptn, see
[Use Keptn with cert-manager.io](../operate/cert-manager.md).

For more advanced installations configurations,see:

* [CertificateManager-README](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn-cert-manager/README.md)

## Customizing the configuration of components

To access and modify the configuration of a subcomponent,
modify the component's `values.yaml` file.
You can use the sub-chart name
as written in the `chart.yaml` file.
The `README` file for each component documents the parameters
that are supported for that component.

Here is an example `values.yaml` altering metrics operator values:

{{< docsembed path="content/en/docs/install/assets/values-advance-changes.yaml" >}}

## Running Keptn with vCluster

Keptn running on Kubernetes versions 1.26 and older
uses a custom
[scheduler](../architecture/components/scheduler.md),
so it does not work with
[Virtual Kubernetes Clusters](https://www.vcluster.com/)
("vClusters") out of the box.
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
