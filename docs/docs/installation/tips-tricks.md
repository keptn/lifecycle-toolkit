---
comments: true
---

# Installation Tips and Tricks

The
[Install Keptn](index.md)
page documents how to install Keptn.
This page provides some background and more examples
that supplement that information.

## Umbrella chart

Keptn v0.9.0 and later is installed using a Helm umbrella chart.
This means that the Helm Chart that installs all of Keptn
actually groups subcharts for individual components
and you can install individual components
without installing all of Keptn.

Keptn is installed using the
[Keptn](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn)
umbrella Helm chart.
Each subchart has its own README file describing possible configuration options,
but configuration changes for the subcharts are added to a single `values.yaml` file.
See
[Customizing the configuration of components](index.md#customizing-the-configuration-of-components)
for an example.

## Installing older versions of Keptn

Installation of Keptn v0.9.0 has two significant differences
compared to the installation of earlier releases:

* Keptn v0.9.0 and later releases use
  the umbrella charts whereas earlier versions did not
* Keptn v0.9.0 and later releases use
  the `keptn` Helm chart, whereas earlier
  releases used the `klt` chart.

To install a version prior to v0.9.0,
use the install command sequence that is documented for that release.
To install the latest version, use the installation commands on the
[Install Keptn](index.md#basic-installation)
[Install Keptn](index.md#basic-installation)
page.

To install an older release,
specify the chart version with the `--version` flag
in the `helm upgrade --install` command for the release you are installing.

## Example configurations by use-case

[Control what components are installed](index.md#customizing-the-configuration-of-components)
discusses how to configure Keptn to include only the components you want.
The following sections summarize and give examples
of the configurations needed for different use cases.

### Enable Keptn Lifecycle Operator (Observability and/or Release Lifecycle)

If you only want to run the Keptn Observability
and/or Release Lifecycle use-cases in your cluster,
you do not need to install the Keptn Metrics Operator.
To disable it, set the `metricsOperator.enabled` value
to `false` as in the following:

```yaml
{% include "./assets/values-only-lifecycle.yaml" %}
```

Note that, if you want to run pre- and/or post-deployment
[evaluations](../guides/evaluations.md)
as part of the Release Lifecycle use-case,
you need to have the Keptn Metrics Operator installed.

You must also enable Keptn for each
[namespace](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/)
on which you want to run either
the Observability or Release Lifecycle use-case.

To enable Keptn, annotate the appropriate `Namespace` resource(s).
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

After enabling Keptn for your namespace(s),
you are ready to
[Integrate Keptn with your applications](../guides/integrate.md).

For more information about implementing Observability, see the
[Observability User Guide](../guides/otel.md).

For more information about implementing Keptn Release Lifecycle, see the
[Deployment tasks](../guides/tasks.md)
and
[Evaluations](../guides/evaluations.md)
User Guides.

### Enable Keptn Metrics Operator (Metrics)

If you are only interested in Metrics,
you do not need the Keptn Lifecycle Operator.
Disable it using the following values.yaml:

```yaml
{% include "./assets/values-only-metrics.yaml" %}
```

For more information about implementing Metrics, see the
[Metrics User Guide](../guides/evaluatemetrics.md).

### Enable Keptn Analysis (SLOs/SLIs)

To enable Keptn Analysis in your cluster,
you again do not need the Keptn Lifcycle Operator.
Disable it using the following values.yaml:

```yaml
{% include "./assets/values-only-metrics.yaml" %}
```

> **Note** A preliminary release of the Keptn Analysis feature
  is included in Keptn v0.8.3 and v0.9.0 but is hidden behind a feature flag.
  See the
  [Analysis](../reference/crd-reference/analysis.md/#differences-between-versions)
  reference page for how to activate the preview of this feature.
>

For more information about implementing Keptn Analysis, see the
[Analysis User Guide](../guides/slo.md).

### Disable Keptn Certificate Manager (Certificates)

If you wish to use your own custom certificate manager,
you can disable the Keptn `cert-manager` by using the
`--set certificateManager.enabled=false` argument
to the `helm upgrade` command line
or you can modify the `values.yaml` file:

```yaml
{% include "./assets/values-remove-certmanager.yaml" %}
```

For more information about using `cert-manager` with Keptn, see
[Use Keptn with cert-manager.io](./configuration/cert-manager.md).
