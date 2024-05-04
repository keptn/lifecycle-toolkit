---
comments: true
---

# Upgrade

If you installed the previous version of Keptn using `helm`,
you can upgrade to the latest version
by running the same command sequence used to install Keptn:

```shell
helm repo add keptn https://charts.lifecycle.keptn.sh
helm repo update
helm upgrade --install keptn keptn/keptn \
   -n keptn-system --create-namespace --wait
```

Use the `--set` flag or download and edit the `values.yaml` file
to modify the configuration as discussed on the
[Install Keptn](./index.md) page.

> **Warning**
If you installed your Keptn instance from the Manifest,
additional steps are required to use the Helm Chart to upgrade.
Contact us on Slack for assistance.

## Upgrade to v1 version

If you have previously used Keptn Lifecycle Operator with API
resources of version `v1alpha3` and `v1alpha4`, you need to
edit manually created or edited `KeptnApp` resources.
For further information please refer to the
[migration section](../migrate/keptnapp/index.md).
