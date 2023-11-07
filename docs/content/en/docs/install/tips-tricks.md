---
title: Installation Tips and Tricks
description: Supplemental installation information
weight: 38
hidechildren: false # this flag hides all sub-pages in the sidebar-multicard.html
--- 

The
[Install Keptn](install.md)
page documents how to install Keptn.
This page provides some background and more examples
that supplement that information.

## Umbrella chart

Keptn v0.9.0 and later is installed using a Helm umbrella chart.
This means that the Helm Chart that installs all of Keptn
actually groups subcharts for individual services
and you can install individual services
without installing all of Keptn.

Keptn is installed using the
[keptn/values.yaml](https://github.com/keptn/lifecycle-toolkit-charts/blob/main/charts/keptn/values.yaml)
umbrella chart.
Each subchart has its own `values.yaml` file
with parameters that are documented in the README file for each,
but configuration changes for the subcharts
are added to the umbrella chart.
See
[Customizing the configuration of components](install.md/#customizing-the-configuration-of-components)
for an example.

## Installing older versions of Keptn

Installation of Keptn v0.9.0 has two significant differences
compared to installation of earlier releases:

* Keptn v0.9.0 and later releases use
  the umbrella charts whereas earlier versions did not
* Keptn v0.9.0 and later releases use
  `keptn` as the value to the `NAME` option to the
  [helm repo add](https://helm.sh/docs/helm/helm_repo_add/)
  command whereas earlier releases used `klt`.

To install a version prior to v0.9.0,
use the install command sequence that is documented for that release.
Use the same command sequence documented for v0.9.0
to determine the CHART version for the release you want to install:
  
```shell
helm repo add keptn https://charts.lifecycle.keptn.sh
helm repo update
helm search repo keptn
```

Then specify that CHART version to the `--version` flag
in the `helm update` command documented for the release you are installing.

## Example configurations by use-case

[Control what components are installed](install/#customizing-the-configuration-of-components)
discusses how to configure Keptn to include only the components you want.
The following sections summarize and give examples
of the configurations needed for different use cases.

### Enable Keptn Lifecycle Operator (Observability and/or Release Lifecycle Management)

If you only want to run the Keptn Observability
and/or Release Lifecycle Management use-cases in your cluster,
you do not need to install the Keptn Metrics Operator.
To disable it, set the `metricsOperator.enabled` value
to `false` as in the following:

{{< embed path="/docs/content/en/docs/install/assets/values-only-lifecycle.yaml" >}}

Note that, if you want to run pre- and/or post-deployment
[evaluations](../implementing/evaluations.md)
as part of the Release Lifecycle Management use-case,
you need to run the Keptn Metrics Operator.

You must also enable Keptn for each
[namespace](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/)
on which you want to run either
the Observability or Release Lifecycle Management use-case.
This is because Keptn communicates with the Kubernetes scheduling mechanism
for tasks such as enforcing checks natively,
stopping a deployment from proceeding when criteria are not met,
doing post-deployment evaluations
and tracing all activities of all deployment
[workloads](https://kubernetes.io/docs/concepts/workloads/)
on the cluster.

To enable Keptn in your namespace,
annotate the appropriate `Namespace` resource(s).
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
This annotation tells the webhook to handle the namespace.

After enabling Keptn for your namespace(s),
you are ready to
[Integrate Keptn with your applications](../implementing/integrate.md).

For more information about implementing Observability, see the
[Observability User Guide](../implementing/otel.md).

For more information about implementing Keptn Release Management, see the
[Deployment tasks](../implementing/tasks.md)
and
[Evaluations](../implementing/evaluations.md)
User Guides.

### Enable Keptn Metrics Operator (Metrics)

If you are only interested in Metrics,
you do not need the Keptn Lifecycle Operator.
Disable it using the following values.yaml:

{{< embed path="/docs/content/en/docs/install/assets/values-only-metrics.yaml" >}}

For more information about implementing Metrics, see the
[Metrics User Guide](../implementing/evaluatemetrics.md).

### Enable Keptn Analysis (SLOs/SLIs)

To enable Keptn Analysis in your cluster,
you again do not need the Keptn Lifcycle Operator.
Disable it using the following values.yaml:

{{< embed path="/docs/content/en/docs/install/assets/values-only-metrics.yaml" >}}

> **Note** A preliminary release of the Keptn Analysis feature
  is included in Keptn v0.8.3 and v0.9.0 but is hidden behind a feature flag.
  See the
  [Analysis](../yaml-crd-ref/analysis.md/#differences-between-versions)
  reference page for how to activate the preview of this feature.
>

For more information about implementing Keptn Analysis, see the
[Analysis User Guide](../implementing/slo.md).

### Disable Keptn Certificate Manager (Certificates)

If you wish to use your own custom certificate manager,
you can disable the Keptn `cert-manager` by using the
`--set "certificateManager.enabled=false"` argument
to the `helm upgrade` command line
or you can modify the `keptn/values.yaml` file:

{{< embed path="/docs/content/en/docs/install/assets/values-remove-certmanager.yaml" >}}

For more information about using `cert-manager` with Keptn, see
[Use Keptn with cert-manager.io](../operate/cert-manager.md).
