---
title: Install and enable Keptn
description: Install Keptn
weight: 35
hidechildren: false # this flag hides all sub-pages in the sidebar-multicard.html
---

Keptn must be installed, enabled, and integrated
into each cluster you want to monitor.
This is because Keptn communicates with the Kubernetes scheduler
for tasks such as enforcing checks natively,
stopping a deployment from proceeding when criteria are not met,
doing post-deployment evaluations
and tracing all activities of all deployment workloads on the cluster.

Two methods are supported for installing Keptn:

* Releases v0.7.0 and later can be installed using
  the [Helm Chart](#use-helm-chart).
  This is the preferred strategy because it allows you to customize your cluster.

* Releases v0.8.2 and earlier can be installed using
  the [manifests](#use-manifests).
  This is the less-preferred way because it does not support customization.

After Keptn is installed, you must
[Enable Keptn for your cluster](#enable-keptn-for-your-cluster)
in order to run some Keptn functionality.

You are then ready to
[Integrate Keptn with your applications](../implementing/integrate).

## Use Helm Chart

Version v0.7.0 and later of Keptn
should be installed using Helm Charts.
The command sequence to fetch and install the latest release is:

```shell
helm repo add klt https://charts.lifecycle.keptn.sh
helm repo update
helm upgrade --install keptn klt/klt \
   -n keptn-lifecycle-toolkit-system --create-namespace --wait
```

Note that the `helm repo update` command is used for fresh installs
as well as for upgrades.

Some helpful hints:

* Use the `--version <version>` flag on the
  `helm upgrade --install` command line to specify a different Keptn version.

* Use the following command sequence to see a list of available versions:

  ```shell
  helm repo update
  helm search repo klt
  ```

* To verify that the Keptn components are installed in your cluster,
  run the following command:

  ```shell
  kubectl get pods -n keptn-lifecycle-toolkit-system
  ```

  The output shows all components that are running on your system.

### Modify Helm configuration options

Helm chart values can be modified before the installation.
This is useful if you want to install only the `metrics-operator`
rather than the full Toolkit
or if you need to change the size of the installation.

To modify configuration options, download a copy of the
[helm/chart/values.yaml](https://github.com/keptn/lifecycle-toolkit/blob/main/helm/chart/values.yaml)
file, modify some values, and use the modified file to install Keptn:

1. Download the `values.yaml` file:

   ```shell
   helm get values RELEASE_NAME [flags] > values.yaml
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
--set key1=value1,key2=value2,....
```

For more information,see

* The [Helm Get Values](https://helm.sh/docs/helm/helm_get_values/)) document

* The [helm-charts](https://github.com/keptn/lifecycle-toolkit/blob/main/helm/chart/README.md) page
  contains the full list of available values.

## Use manifests

Versions v0.8.2 and earlier of Keptn can be installed using manifests,
although we recommend that you use Helm Charts
because they allow you to easily customize your configuration.

Versions 0.6.0 and earlier can only be installed using manifests.

> **Note** When installing Version 0.6.0,
you must first install the `cert-manager` with the following command sequence:

```shell
kubectl apply \
   -f https://github.com/cert-manager/cert-manager/releases/download/v1.11.0/cert-manager.yaml
kubectl wait \
   --for=condition=Available deployment/cert-manager-webhook -n cert-manager --timeout=60s
```

Use a command sequence like the following
to install Keptn from the manifest,
specifying the version you want to install.

```shell
kubectl apply \
   -f https://github.com/keptn/lifecycle-toolkit/releases/download/v0.6.0/manifest.yaml
kubectl wait --for=condition=Available deployment/lifecycle-operator \
   -n keptn-lifecycle-toolkit-system --timeout=120s
```

Keptn and its dependencies are now installed and ready to use.

## Enable Keptn for your cluster

To enable the Keptn in your cluster,
annotate the Kubernetes
[Namespace](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/)
resource for each namespace in the cluster.
For an example of this, see
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

You see the annotation line `keptn.sh/lifecycle-toolkit: "enabled"`.
This annotation tells the webhook to handle the namespace.

After enabling Keptn for your namespace(s),
you are ready to
[Integrate Keptn with your applications](../implementing/integrate).
