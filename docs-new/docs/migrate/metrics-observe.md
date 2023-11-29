# Migrate Quality Gates

The SLIs and SLOs used for Keptn v1 quality gates can be ported to
appropriate Keptn facilities:

* [Keptn Metrics](../guides/evaluatemetrics.md/)
  allow you to define and view metrics
  from multiple data sources in your Kubernetes cluster.
* Use
  [Keptn Evaluations](../guides/evaluations.md)
  to do a simple evaluation of the metrics data you capture.
  To implement this, transfer the information from the Keptn v1
  [sli.yaml](https://keptn.sh/docs/1.0.x/reference/files/sli/)
  and
  [slo.yaml](https://keptn.sh/docs/1.0.x/reference/files/slo/)
  files into
  [KeptnEvaluationDefinition](../reference/crd-reference/evaluationdefinition.md)
  resources.

* Keptn v1 calculations that use weighting and scoring
  can be converted to use the
  [Keptn Analysis](../guides/slo.md)
  feature.
  Tools are provided to help with this conversion;
  see below.

By default, Keptn includes additional observability features
that are not included in Keptn v1:

* [DORA metrics](../guides/dora.md)
* [OpenTelemetry observability](../guides/otel.md)

## Paradigm changes

The Keptn paradigm for evaluations and analyses
differs from that of Keptn v1 quality gates:

* Data providers are installed and configured into your Kubernetes cluster
  using Helm charts and standard practices.
* Keptn supports multiple instances of multiple data providers.
* You must populate a
  [KeptnMetricsProvider](../reference/crd-reference/metricsprovider.md) resource
  for each instance of each data provider.
  This resource specifies the URL and namespace for the data provider
  and gives it a unique `name` that can be referenced
  for Keptn Metrics (which are also used for Evaluations) and Analysis.
* Queries and objectives that are specified in your Keptn v1
  [sli.yaml](https://keptn.sh/docs/1.0.x/reference/files/sli/)
  and
  [slo.yaml](https://keptn.sh/docs/1.0.x/reference/files/slo/)
  files are transferred/converted into Keptn resources.

## Transfer Keptn v1 SLIs/SLOs to evaluation resources

Simple comparisons of data can be implemented as
[Keptn Evaluations](../guides/evaluations.md).
To implement this:

* Transfer the information from the Keptn v1
  [sli.yaml](https://keptn.sh/docs/1.0.x/reference/files/sli/)
  files into
  [KeptnMetric](../reference/crd-reference/metric.md) resources
* Transfer the information from the Keptn v1
  [slo.yaml](https://keptn.sh/docs/1.0.x/reference/files/slo/)
  files into
  [KeptnEvaluationDefinition](../reference/crd-reference/evaluationdefinition.md)
  resources.

## Convert Keptn v1 SLIs/SLOs to Analysis resources

The Keptn Analysis feature provides capabilities
similar to those of the Keptn v1
[Quality Gates](https://keptn.sh/docs/1.0.x/define/quality-gates/)
feature
but it uses Kubernetes resources to define the analysis to be done
rather than the configuration files used for Keptn v1.
Tools are provided to convert Keptn v1 SLI and SLO definitions
to Keptn Analysis resources.

These instructions assume that the same SLO file
has been used for all services in the project,
so only one `AnalysisDefinition` resource
(named `my-project-ad`) is created.
If your Keptn v1 project has multiple SLOs,
you need to create a separate `AnalysisDefinition`
with a unique name for each one.

The process is:

1. Convert the SLIs to `AnalysisValueTemplates` resources

   The following command sequence converts a Keptn v1
   [sli.yaml](https://keptn.sh/docs/1.0.x/reference/files/sli/)
   file to a Keptn
   [AnalysisValueTemplate](../reference/crd-reference/analysisvaluetemplate.md)
   resource:

   <!
