---
title: Migrate Quality Gates
description: How to translate Quality Gates into Keptn Metrics and evaluations
weight: 40
---

The SLIs and SLOs used for Keptn v1 quality gates can be ported to
appropriate Keptn facilities:

* [Keptn Metrics](../implementing/evaluatemetrics.md/)
  allow you to define and view metrics
  from multiple data sources in your Kubernetes cluster.
* Use
  [Keptn Evaluations](../implementing/evaluations.md)
  to do a simple evaluation of the metrics data you capture.
  To implement this, transfer the information from the Keptn v1
  [sli.yaml](https://keptn.sh/docs/1.0.x/reference/files/sli/)
  and
  [slo.yaml](https://keptn.sh/docs/1.0.x/reference/files/slo/)
  files into
  [KeptnEvaluationDefinition](../yaml-crd-ref/evaluationdefinition.md)
  resources.

* Keptn v1 calculations that use weighting and scoring
  can be converted to use the
  [Keptn Analysis](../implementing/slo.md)
  feature.
  Tools are provided to help with this conversion;
  see below.

By default, Keptn includes additional observability features
that are not included in Keptn v1:

* [DORA metrics](../implementing/dora.md)
* [OpenTelemetry observability](../implementing/otel.md)

## Paradigm changes

The Keptn paradigm for evaluations and analyses
differs from that of Keptn v1 quality gates:

* Data providers are installed and configured into your Kubernetes cluster
  using Helm charts and standard practices.
* Keptn supports multiple instances of multiple data providers.
* You must populate a
  [KeptnMetricsProvider](../yaml-crd-ref/metricsprovider.md) resource
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
[Keptn Evaluations](../implementing/evaluations.md).
To implement this:

* Transfer the information from the Keptn v1
  [sli.yaml](https://keptn.sh/docs/1.0.x/reference/files/sli/)
  files into
  [KeptnMetric](../yaml-crd-ref/metric.md) resources
* Transfer the information from the Keptn v1
  [slo.yaml](https://keptn.sh/docs/1.0.x/reference/files/slo/)
  files into
  [KeptnEvaluationDefinition](../yaml-crd-ref/evaluationdefinition.md)
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
   [AnalysisValueTemplate](../yaml-crd-ref/analysisvaluetemplate.md)
   resource:

   ```shell
   METRICS_OPERATOR_IMAGE=ghcr.io/keptn/metrics-operator:v0.8.2
   PATH_TO_SLI=sli.yaml
   KEPTN_PROVIDER_NAME=my-prometheus-provider
   KEPTN_PROVIDER_NAMESPACE=keptn-lifecycle-poc

   docker run -v .:/mydata $METRICS_OPERATOR_IMAGE \
     --convert-sli=mydata/$PATH_TO_SLI \
     --keptn-provider-name=$KEPTN_PROVIDER_NAME \
     --keptn-provider-namespace=$KEPTN_PROVIDER_NAMESPACE > analysis-value-template.yaml
   ```

   This command creates an `AnalysisValueTemplate` resource
   for each SLI that is defined in the `sli.yaml` file.

   > **Note**
     The `name` of each `AnalysisValueTemplate` resource
     must adhere to the
     [Kubernetes resource naming rules](https://kubernetes.io/docs/concepts/overview/working-with-objects/names/).
     The conversion tools preserve the name of each SLI
     but modify the names to match the Kubernetes requirements.
     For example, the `Response_Time` SLI
     generates an `AnalysisValueTemplate` that is named `response-time`.
   >

   All the SLIs for a particular `sli.yaml` file
   are defined in a file called `analysis-value-template.yaml`.
   Apply this file to your cluster with a command like the following.
   Be sure to specify the namespace;
   if you omit it, the yaml file is applied to the default namespace.

   ```shell
   kubectl apply -f analysis-value-template.yaml -n keptn-lifecycle-poc
   ```

1. Convert the SLO to an `AnalysisDefinition` resource:

   The process of converting the Keptn v1
   [slo.yaml](https://keptn.sh/docs/1.0.x/reference/files/slo/)
   files to
   [AnalysisDefinition](../yaml-crd-ref/analysisdefinition.md)
   resources is similar to the process of converting the SLIs.
   Use the following command sequence:

   ```shell
   METRICS_OPERATOR_IMAGE=ghcr.io/keptn/metrics-operator:v0.8.2
   PATH_TO_SLO=slo.yaml
   ANALYSIS_VALUE_TEMPLATE_NAMESPACE=keptn-lifecycle-poc
   ANALYSIS_DEFINITION_NAME=my-project-ad

   docker run -v $(pwd):/mydata $METRICS_OPERATOR_IMAGE \
     --convert-slo=/mydata/$PATH_TO_SLO \
     --analysis-value-template-namespace=$ANALYSIS_VALUE_TEMPLATE_NAMESPACE \
     --analysis-definition-name=$ANALYSIS_DEFINITION_NAME > analysis-definition.yaml
   ```

   The result of this command yields an `AnalysisDefinition` resource
   that is defined in a file called `analysis-definition.yaml`.
   Apply this to your Keptn cluster with the following command.
   Be sure to add the namespace;
   if you omit it, the yaml file is applied to the default namespace.

   ```shell
   kubectl apply -f analysis-definition.yaml -n keptn-lifecycle-poc
   ```

1. Create a `KeptnMetricsProvider` resource

   A [KeptnMetricsProvider](../yaml-crd-ref/metricsprovider.md)
   resource configures the data provider from which the values
   for the `AnalysisValueTemplate` resource are fetched.
   This same resource is used for any metrics and evaluations you are using.
   Note that Keptn supports multiple instances of multiple data providers
   and each `AnalysisValueTemplate` resource
   can reference a different provider,

   The `KeptnMetricsProvider` resource defines:

   * The URL of the server used for your data provider
   * The namespace where the data provider is located
   * (optional) Secret information to access a protected data source

   See the `prometheus-provider` file for an example.

1. Define the Analysis resource to run

   Create a yaml file (such as `analysis-instance.yaml`)
   to populate the
   [Analysis](../yaml-crd-ref/analysis.md)
   resource that defines the specific analysis you want to run.
   Specify the following:

   * Name of the `AnalysisDefinition` resource that defines the goals
   * Timeframe for which to retrieve the metrics.
     The `from`/`to` timestamps are added automatically
     to the query defined in the `AnalysisValueTemplate` resource.
   * (optional) List of `args` that specify values for this analysis
     that replace variables used in the queries
     in the `AnalysisValueTemplate` resources.
     See
     [Passing contextual arguments to the Analysis](#passing-contextual-arguments-to-the-analysis)
     for more information.

1. Run the analysis

   To perform an Analysis (or "trigger an evaluation" in Keptn v1 jargon),
   apply the `analysis-instance.yaml` file:

   ```shell
   kubectl apply -f analysis-instance.yaml -n keptn-lifecycle-poc
   ```

   Retrieve the current status of the Analysis with the following command:

   ```shell
   kubectl get analysis -n keptn-lifecycle-poc
   ```

   This yields an output that looks like the following:

   ```shell
   NAME                ANALYSISDEFINITION      WARNING   PASS
   analysis-sample-1   my-project-ad             true
   ```

   This shows that the analysis passed successfully.

   To get the detailed result of the evaluation,
   use the `-oyaml` argument to inspect the full state of the analysis:

   This displays the `Analysis` resource
   with the definition of the analysis
   as well as the `status` (results) of the analysis; for example:

   ```shell
   kubectl get analysis -n keptn-lifecycle-poc -oyaml
   ```

   ```yaml
   apiVersion: v1
   items:
   - apiVersion: metrics.keptn.sh/v1alpha3
     kind: Analysis
     metadata:
       creationTimestamp: "2023-09-14T11:00:01Z"
       generation: 4
       name: analysis-sample-1
       namespace: keptn-lifecycle-poc
       resourceVersion: "71327"
       uid: 1c5e043d-ed5e-42f8-ba32-b7af54b55c35
     spec:
       analysisDefinition:
         name: my-project-ad
         namespace: keptn-lifecycle-poc
       args:
         ns: keptn-lifecycle-toolkit-system
         project: my-project
       timeframe:
         from: "2023-09-14T11:20:19Z"
         to: "2023-09-14T11:22:19Z"
     status:
       pass: true
       raw: '{"objectiveResults":[{"result":{"failResult":{"operator":{"greaterThan":{"fixedValue":"50"}}},"warnResult":{"operator":{"greaterThan":{"fixedValue":"50"}}},"pass":true},"value":7,"score":1}],"totalScore":1,"maximumScore":1,"pass":true}'
   kind: List
   metadata:
     resourceVersion: ""
   ```

   As can be seen in the yaml above,
   the `status.raw` property contains the detailed breakdown
   of the analysis goals, and whether or not they passed.

   The result of this analysis stays in the cluster
   until the `Analysis` is deleted.
   That also means that, if another analysis should be performed,
   the new analysis must be given a new, unique name within the namespace.

### Passing contextual arguments to the Analysis

In some cases, especially when migrating from Keptn v1,
you may want to do an analysis for different services within a project
and adapt the query for fetching the metric values,
based on which service/stage you are evaluating.
Keptn enables you to do this by using variables in the query defined in the
`AnalysisValueTemplates` resource
and arguments that are defined in the `Analysis` resource
for the specific analysis being run.

For example, you may have a query
for retrieving the response time of a service,
which is identified by its label in Prometheus.
In this case, use the go templating syntax
to insert a variable as a placeholder
(for example, in this case, `{{.service}}`)
for the service identifier in the prometheus query:

```yaml
apiVersion: metrics.keptn.sh/v1alpha3
kind: AnalysisValueTemplate
metadata:
  creationTimestamp: null
  name: response-time
spec:
  provider:
    name: my-prometheus-provider
    namespace: keptn-lifecycle-poc
  query: response_time{label="{{.service}}"}
```

Then, if an analysis for that particular service should be performed,
the name of the service can be passed to the analysis
using the `spec.args` property from the `Analysis` resource:

```yaml
apiVersion: metrics.keptn.sh/v1alpha3
kind: Analysis
metadata:
  name: analysis-sample-1
  namespace: keptn-lifecycle-poc
spec:
  timeframe:
    from: 2023-09-14T11:20:19Z
    to: 2023-09-14T11:22:19Z
  args:
    "service": "my-service"
  analysisDefinition:
    name: my-project-ad
    namespace: keptn-lifecycle-poc
```

This way, you can use the same `AnalysisDefinition`
and `AnalysisValueTemplates` for multiple services within the same project.
~
