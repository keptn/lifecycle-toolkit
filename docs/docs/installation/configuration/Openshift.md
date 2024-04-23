---
comments: true
---

# Installing on Openshift

To install on OpenShift, set the value `global.openShift.enabled` in the `values.yaml` file to true.
In practice this means that `runAsUser` and `runAsGroup` are removed, since
Openshift sets those automatically.

You can set the `global.openShift.enabled` parameter when running the `helm install` command:

```shell
helm install keptn keptn/keptn -n keptn-system --create-namespace --set global.openShift.enabled=true
```

or you can define it in your `values.yaml` file:

```yaml
global:
  openShift:
    enabled: true
```
