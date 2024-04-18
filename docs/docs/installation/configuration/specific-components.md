---
comments: true
---

# Control what components are installed

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

## Disable Keptn Certificate Manager (Certificates)

If you wish to use your custom certificate manager,
you can disable Keptn `cert-manager` by setting the
`certificateManager.enabled` Helm value to `false`:

```yaml
{% include "../assets/values-remove-certmanager.yaml" %}
```

For more information on using `cert-manager` with Keptn, see
[Use Keptn with cert-manager.io](../../components/certificate-operator.md).

For the full list of Helm values, see the
[keptn-cert-manager Helm chart README](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn-cert-manager/README.md).

## Uninstalling Keptn

To uninstall Keptn from your cluster, please follow the steps
on the [Uninstall page](../uninstall.md).
