# Installation Instructions

## Install version 0.6.0 and above

In version 0.6.0 and later, you can install the Lifecycle Toolkit using helm charts:

```shell
helm repo add keptn-lifecycle-toolkit https://charts.lifecycle.keptn.sh
helm repo update
helm upgrade --install keptn-lifecycle-toolkit keptn/lifecycle-toolkit -n keptn-lifecycle-toolkit-system --create-namespace --wait
```

To install a specific version use `--version <version>` as part of `helm upgrade --install` command.

To list available versions:

```shell
helm repo update
helm search repo keptn-lifecycle-toolkit
```

The `helm upgrade --install` command offers a flag called `--set`, which can be used to specify configuration options using the format key1=value1,key2=value2,....
The full list of available flags can be found in the [helm-charts](https://github.com/keptn/lifecycle-toolkit/blob/main/helm/chart/README.md).

For installing the Lifecycle Toolkit via manifests use:

<!---x-release-please-start-version-->

```shell
kubectl apply -f https://github.com/keptn/lifecycle-toolkit/releases/download/v0.6.0/manifest.yaml
kubectl wait --for=condition=Available deployment/lifecycle-operator -n keptn-lifecycle-toolkit-system --timeout=120s
```

<!---x-release-please-end-->

The Lifecycle Toolkit and its dependencies are now installed and ready to use.

**Note:** Installation of the Lifecycle Toolkit version 0.5.0 and lower is not supported via helm charts.

## Install version 0.5.0 and earlier

You must first install *cert-manager* with the following commands:

```shell
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.11.0/cert-manager.yaml
kubectl wait --for=condition=Available deployment/cert-manager-webhook -n cert-manager --timeout=60s
```

After that, you can install the Lifecycle Toolkit `<oldversion>` with:

```shell
kubectl apply -f https://github.com/keptn/lifecycle-toolkit/releases/download/<oldversion>/manifest.yaml
kubectl wait --for=condition=Available deployment/lifecycle-operator -n keptn-lifecycle-toolkit-system --timeout=120s
```
