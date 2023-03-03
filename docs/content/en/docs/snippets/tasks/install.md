# Installation Instructions

## Install version 0.6.0 and above

In version 0.6.0 and later, you can install the Lifecycle Toolkit using the current release manifest:

<!---x-release-please-start-version-->

```shell
kubectl apply -f https://github.com/keptn/lifecycle-toolkit/releases/download/v0.5.0/manifest.yaml
kubectl wait --for=condition=Available deployment/lifecycle-operator -n keptn-lifecycle-toolkit-system --timeout=120s
```

<!---x-release-please-end-->

The Lifecycle Toolkit and its dependencies are now installed and ready to use.

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
