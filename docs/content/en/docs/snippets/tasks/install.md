# Installation Instructions

## Install version 0.7.0 and above

In version 0.7.0 and later, you can install the Lifecycle Toolkit using either helm charts or manifests.

For installing the Lifecycle Toolkit via Helm chart:

```shell
helm repo add klt https://charts.lifecycle.keptn.sh
helm repo update
helm upgrade --install keptn klt/klt -n keptn-lifecycle-toolkit-system --create-namespace --wait
```

To install a specific version, use the `--version <version>` flag as part of the
`helm upgrade --install` command.

To list available versions:

```shell
helm repo update
helm search repo keptn-lifecycle-toolkit
```

The `helm upgrade --install` command offers a flag called `--set`, which can be used to specify
configuration options using the format key1=value1,key2=value2,....

Or you could download the chart value file and modify it using

```shell
helm get values RELEASE_NAME [flags] > values.yaml
```

and install adding `--values=values.yaml` to your `helm upgrade` command (official documentation
available [here](https://helm.sh/docs/helm/helm_get_values/)).

The full list of available flags can be found in the [helm-charts](https://github.com/keptn/lifecycle-toolkit/blob/main/helm/chart/README.md).

> **Note**
Installation of the Lifecycle Toolkit version 0.6.0 and lower is not supported via helm charts.

<details>
<summary>Install Keptn using Manifests</summary>

All versions of the Lifecycle Toolkit can be installed using manifests,
with a command like the following:

<!---x-release-please-start-version-->

```shell
kubectl apply -f https://github.com/keptn/lifecycle-toolkit/releases/download/v0.8.0/manifest.yaml
kubectl wait --for=condition=Available deployment/lifecycle-operator -n keptn-lifecycle-toolkit-system --timeout=120s
```

<!---x-release-please-end-->

The Lifecycle Toolkit and its dependencies are now installed and ready to use.

</details>
