---
title: AnalysisDefinition
description: Define SLOs for an Analysis
weight: 6
hide: true
---

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

* **apiVersion** metrics.keptn.sh/v1alpha3
* **kind** AnalysisDefinition
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

[AnalysisDefinition](../../crd-ref/metrics/v1alpha3/#analysisdefinition)
API reference

## Differences between versions

## See also

* [Analysis](analysis.md)
* [AnalysisValueTemplate](analysisvaluetemplate.md)
* [Analysis](../implementing/slo) guide
