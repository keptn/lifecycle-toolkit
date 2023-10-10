---
title: AnalysisDefinition
description: Define SLOs for an Analysis
weight: 6
---

An `AnalysisDefinition` resource defines the
list of Service Level Objectives (SLOs) for an `Analysis`.

## Synopsis

```yaml
apiVersion: metrics.keptn.sh/v1alpha3
kind: AnalysisDefinition
metadata:
  name: ed-my-proj-dev-svc1
  namespace: keptn-lifecycle-toolkit-system
spec:
  objectives:
    - analysisValueTemplateRef:
        name: response-time-p95
        namespace: keptn-lifecycle-toolkit-system
      target:
        failure:
          lessThan:
            fixedValue: 600
        warning:
          inRange:
            lowBound: 300
            highBound: 500
      weight: 1
      keyObjective: false
  totalScore:
    passPercentage: 90
    warningPercentage: 75
```

## Fields

* **apiVersion** -- API version being used
* **kind** -- Resource type.
   Must be set to AnalysisDefinition.
* **metadata**
  * **name** ed-my-proj-dev-svc1
  * **namespace** keptn-lifecycle-toolkit-system
* **spec**
  * **objectives**
    * **analysisValueTemplateRef**
      * **name** response-time-p95
      * **namespace** keptn-lifecycle-toolkit-system
      * **target**
        * **failure**
          * **lessThan**
            * **fixedValue** 600
        > **Warning**
        * **inRange**
          * **lowBound** 300
          * **highBound** 500
      * **weight** 1
      * **keyObjective** false
  * **totalScore**
    * **passPercentage** 90
    > **Warning**Percentage

## Usage

An `AnalysisDefinition` resource contains a list of objectives to satisfy.
Each of these objectives must specify:

* Failure or warning target criteria
* Whether the objective is a key objective
  meaning that its failure fails the Analysis
* Weight of the objective on the overall Analysis
* The `AnalysisValueTemplate` resource that contains the SLIs,
  defining the data provider from which to gather the data
  and how to compute the Analysis

## Examples

```yaml
apiVersion: metrics.keptn.sh/v1alpha3
kind: AnalysisDefinition
metadata:
  name: ed-my-proj-dev-svc1
  namespace: keptn-lifecycle-toolkit-system
spec:
  objectives:
    - analysisValueTemplateRef:
        name: response-time-p95
        namespace: keptn-lifecycle-toolkit-system
      target:
        failure:
          lessThan:
            fixedValue: 600
        warning:
          inRange:
            lowBound: 300
            highBound: 500
      weight: 1
      keyObjective: false
  totalScore:
    passPercentage: 90
    warningPercentage: 75
```

For an example of how to implement the Keptn Analysis feature, see the
[Analysis](../implementing/slo)
guide page.

## Files

[AnalysisDefinition](../crd-ref/metrics/v1alpha3/#analysisdefinition)
API reference

## Differences between versions

A preliminary release of the Keptn Analysis feature
is included in Keptn v.0.8.3 but is hidden behind a feature flag.
To preview these features, set the environment `ENABLE_ANALYSIS` to `true`
in the `metrics-operator` deployment.

## See also

* [Analysis](analysis.md)
* [AnalysisValueTemplate](analysisvaluetemplate.md)
* [Analysis](../implementing/slo) guide
