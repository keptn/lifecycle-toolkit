---
title: Migrate Quality Gates to KLT metrics and evaluations
description: How to translate Quality Gates into Keptn Metrics and evaluations
weight: 40
hidechildren: false # this flag hides all sub-pages in the sidebar-multicard.html
---

> **Note**
This section is under development.
Information that is published here has been reviewed for technical accuracy
but the format and content is still evolving.
We hope you will contribute your experiences
and questions that you have.

The SLIs and SLOs used for Keptn v1 quality gates can be ported to
KLT KeptnMetrics and KeptnEvaluationDefintions.

By default, KLT includes additional observability features
that are not included in Keptn v1:

* [DORA metrics](../../implementing/dora)
* [OpenTelemetry observability](../../implementing/otel.md)

Keptn v1 Quality Gates can be migrated into KLT metrics
and KLT evaluations.

> **Note**
The full SLO capabilities
provided by Keptn v1 such as weighting and scoring
are currently under development for KLT.
You can follow and participate in the design and implementation process at
[GitHub Epic 1646](https://github.com/keptn/lifecycle-toolkit/issues/1646).

Notice the paradigm differences when implementing KLT evaluations:

* Data providers are installed and configured into your Kubernetes cluster
  using Helm charts and standard practices.
* You must populate a
  [KeptnMetricsProvider](../../yaml-crd-ref/metricsprovider.md) resource
  for each instance of each data provider.
  This resource specifies the URL and namespace for the data provider
  and gives it a unique `name` that can be referenced by other resources.
  This is the only KLT specific configuration that is required.
* Each query that is specified in your Keptn v1
  [slo.yaml](https://keptn.sh/docs/1.0.x/reference/files/sli/) file
  should be converted into its own
  [KeptnMetric](../../yaml-crd-ref/metric.md) resource
  if you are using it for an evaluation.
  Note that KLT supports using multiple instances of multiple data providers.
* Simple Keptn v1 comparisons that are defined in
  [slo.yaml](https://keptn.sh/docs/1.0.x/reference/files/slo/)
  files can be converted to
  [KeptnEvaluationDefinition](../../yaml-crd-ref/evaluationdefinition.md)
  resources.
* Keptn v1 calculations that use weighting and scoring
  cannot currently be converted to `KeptnEvaluationDefinition` resources
  but are under development.

For more information about working with Keptn metrics, see the
[Keptn Metrics](../../implementing/evaluatemetrics.md)
page.
