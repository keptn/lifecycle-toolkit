---
comments: true
---

# Allocate namespaces for Keptn resources

Keptn primarily operates on Kubernetes
[Workload](https://kubernetes.io/docs/concepts/workloads/)
resources and
[KeptnApp](../../reference/crd-reference/app.md)
resources that are activated and defined by annotations to each workload.
You have significant flexibility to decide how many namespaces to use
and where to locate each resource.
See the Kubernetes
[Namespace](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/)
documentation for some basic information.
You can also search and find many "Best Practices for Namespaces"
documents published on the web.

Possible namespace designs for Keptn run the gamut:

* Run all your Keptn work in a single namespace
* Create a separate namespace for each logical grouping of your Keptn work
* Create a separate namespace for each [workload](https://kubernetes.io/docs/concepts/workloads/)

This page summarizes some namespace considerations
for some Keptn resources.
For information about limiting the namespaces where Keptn is allowed, see
[Namespaces and Keptn](namespace-keptn.md).

## KeptnMetric resources

[KeptnMetricsProvider](../../reference/crd-reference/metricsprovider.md)
resources need to be located
in the same namespace as the associated
[KeptnMetric](../../reference/crd-reference/metric.md)
resources.
But
[KeptnEvaluationDefinition](../../reference/crd-reference/evaluationdefinition.md)
resources that are used for pre- and post-deployment
can reference metrics from any namespace.
So you can create `KeptnMetrics` resources in a centralized namespace
(such as `keptn-system`)
and access those metrics in evaluations on all namespaces in the cluster.

## Analysis related resources

Analysis related resources
([Analysis](../../reference/crd-reference/analysis.md),
[AnalysisDefinition](../../reference/crd-reference/analysisdefinition.md),
and
[AnalysisValueTemplate](../../reference/crd-reference/analysisvaluetemplate.md))
reference each other via a `name` and, optionally, a `namespace` field.
The `Analysis` resource references the `AnalysisDefinition` resource,
which then references the `AnalysisValueTemplate` resources.

* If the `namespace` in the reference is not set explicitly,
  the `AnalysisDefinition` and `AnalysisValueTemplate` resources
  must reside in the same namespace as the `Analysis` resource.
* If the `namespace` in the reference is set for the resources,
  the `Analysis`, `AnalysisDefinition`, and `AnalysisValueTemplate` resources
  can each reside in different namespaces.

This provides configuration options such as the following:

* You can have one namespace
  with all of your `AnalysisDefinition` and `AnalysisValueTemplate` resources
  and reuse them in the different namespaces where you run analyses.

* You can have everything strictly namespaced
  and always put the `AnalysisDefinition`, `AnalysisValueTemplate`
  and the `Analysis` resources into the same namespace,
  without adding the explicit namespace selectors
  when creating references between those objects.

## KeptnApp resources

Each `KeptnApp` resource identifies the namespace to which it belongs.
If you configure multiple namespaces,
you can have `KeptnApp` resources with the same name
in multiple namespaces without having them conflict.

You do not need separate namespaces for separate versions of your application.
The `KeptnApp` resource includes fields to define
the `version` as well as a `revision`
(used if you have to rerun a deployment
but want to retain the version number).
