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

The Keptn Lifecycle Toolkit can be migrated to
Keptn Metrics and Keptn Evaluations.

Note that KLT includes additional observability features
that are not included in Keptn v1 by default:

* [Dora metrics](../../implementing/dora)
* [OpenTelemetry observability](../../implementing/otel)

Keptn v1 Quality Gates can be migrated into KLT metrics
and KLT evaluations.
Notice the paradigm differences:

* Data providers are installed and configured into your Kubernetes cluster
  using Helm charts and standard practices.
* You must populate a
  [KeptnMetricsProvider](../../yaml-crd-ref/metricsprovider.md) resource
  for each instance of each data source.
  This resource specifies the URL and namespace for the data source
  and gives it a unique `name` that can be referenced by other resources.
  This is the only KLT specific configuration that is required.
* Each query that is specified in your Keptn v1
  [slo.yaml](https://keptn.sh/docs/1.0.x/reference/files/sli/) file
  is converted into its own
  [KeptnMetric](../../yaml-crd-ref/metric) resource.
  Note that KLT supports using multiple instances of multiple data providers.
* Simple Keptn v1 comparisons that are defined in
  [slo.yaml](https://keptn.sh/docs/1.0.x/reference/files/slo/)
  files can be converted to
  [KeptnEvaluationDefinition](../../yaml-crd-ref/evaluationdefinition)
  resources.
  Keptn v1 calculations that use weighting and scoring
  cannot currently be converted to `KeptnEvaluationDefinition` resources.

## Define KeptnMetricsProvider resources

You must define a
[KeptnMetricsProvider](../../yaml-crd-ref/metricsprovider.md) resource
for each instance of each data provider you are using.

Note the following:

* Each `KeptnMetricsProvider` resource is bound to a specific namespace.
* Each `KeptnMetric` resource must be located in the same namespace
  as the associated `KeptnMetricsProvider` resource.
* `KeptnEvaluationDefinition` resources can reference metrics
  from any namespace in the cluster.
* To define metrics that can be used in evaluations
  on all namespaces in the cluster,
  create `KeptnMetricsProvider` and KeptnMetric resources
  in a centralized namespace
  such as `keptn-lifecycle-toolkit-system`.

To configure a data source into your KLT cluster:

1. Create a secret if your data source uses one.  See
   [Create secret text](../../implementing/tasks/#create-secret-text).
1. Install and configure each instance of each data source
   into your KLT cluster,
   following the instructions provided by the data source provider.
   See
   [Prepare your cluster for KLT](../../install/k8s/#prepare-your-cluster-for-klt)
for links.
   KLT supports using multiple instances of multiple data sources.
1. Define a
   [KeptnMetricsProvider](../../yaml-crd-ref/metricsprovider.md)
   resource for each data source.

For example, the `KeptnMetricProvider` resource
for a Prometheus data source that does not use a secret
could look like:

```
apiVersion: metrics.keptn.sh/v1alpha2
kind: KeptnMetricsProvider
metadata:
  name: prometheus
  namespace: simplenode-dev
spec:
  targetServer: "http://prometheus-k8s.monitoring.svc.cluster.local:9090"
```

The `KeptnMetricProvider resource for a Dynatrace data source
that uses a secret could look like:

```
```yaml
apiVersion: metrics.keptn.sh/v1alpha3
kind: KeptnMetricsProvider
metadata:
  name: dynatrace
  namespace: podtato-kubectl
spec:
  targetServer: "<dynatrace-tenant-url>"
  secretKeyRef:
    name: dt-api-token
    key: DT_TOKEN
```

## Create a KeptnMetric resource for each SLI

## Create KeptnEvaluationDefinition resources for SLOs
