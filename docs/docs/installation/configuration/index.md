---
comments: true
---

# Configuration

The Keptn configuration is controlled by a set of Helm value files,
summarized in the following table.
The Keptn Helm chart is an umbrella chart
that contains subcharts for all components of Keptn.
Each component has its own Helm values file
(documented in its own README file),
that defines configurations specific to that component.

All configuration changes for all components
can be made in one `values.yaml` file.
This is discussed more in
[Customizing the configuration of components](#customizing-the-configuration-of-components)
below.

The following table summarizes the Keptn `values.yaml` files.

* The "Component" column leads you to the
  README files for each component where
  all Helm values are documented
* The "Configuration file" column leads you to
  the Helm values files for each component

| Component                                                                                                                   | Used for                                                                                                                         | Configuration file                                                                                                                              |
|-----------------------------------------------------------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------|
| [Keptn](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn/README.md)                                 | Installs subcharts, global configuration                                                                                         | [keptn/values.yaml](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn/values.yaml)                                       |
| [lifecycle-operator](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn-lifecycle-operator/README.md) | [Observability](../../guides/otel.md), [Release Lifecycle Management](../../core-concepts/index.md#release-lifecycle-management) | [keptn-lifecycle-operator/values.yaml](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn-lifecycle-operator/values.yaml) |
| [metrics-operator](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn-metrics-operator/README.md)     | [Keptn metrics](../../guides/evaluatemetrics.md), [Analysis](../../guides/slo.md)                                                | [keptn-metrics-operator/values.yaml](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn-metrics-operator/values.yaml)     |
| [cert-manager](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn-cert-manager/README.md)             | [TLS Certificate management for all Keptn components](../../components/certificate-operator.md)                                  | [keptn-cert-manager/values.yaml](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn-cert-manager/values.yaml)             |

## Customizing the configuration of components

To modify configuration of the subcomponents of Keptn,
use the corresponding Helm values in your `values.yaml` file.
Use the subcomponent's parent value as the root for your configuration.

Here is an example `values.yaml` altering global and metrics operator values:

```yaml
{% include "./assets/values-advance-changes.yaml" %}
```

Note the additional values that are specified
in the `metricsOperator` section.
These are documented in the README file for that operator,
which is linked from the `metrics-operator` item under "Component"
in the table above.
To implement this:

* Go into the `values.yaml` file linked under "Configuration file"
* Copy the lines for the values you want to modify
* Paste those lines into your `values.yaml` file
  and modify their values in that file.

### Modify Helm values

To modify Helm values:

1. Download a copy of the Helm values file:

    ```shell
    helm show values keptn/keptn > values.yaml
    ```

1. Edit your local copy to modify some values

1. Add the following string
   to your `helm upgrade` command to install Keptn
   with your configuration changes:

    ```shell
    --values=values.yaml
    ```

   For example, if you create a `my.values.yaml`
   and modify some configuration values,
   use the following command to apply your configuration:

    ```shell
    helm upgrade --install keptn keptn/keptn \
    --values my.values.yaml \
    -n keptn-system --create-namespace --wait
    ```

   You can also use the `--set` flag
   to specify a value change for the `helm upgrade --install` command.
   Helm values are specified using the format:

    ```shell
    --set key1=value1 \
    --set key2=value2 ...
    ```
