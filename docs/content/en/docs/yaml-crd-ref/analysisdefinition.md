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
  name: <name-of-this-resource>
  namespace: <namespace-where-this-resource-resides>
spec:
  objectives:
    - analysisValueTemplateRef:
        name: <name-of-referenced-analysisValueTemplateRef-resource>
        namespace: <namespace-of-the-template-ref>
      target:
        failure | warning:
          <operator>:
            <operatorValue>: <quantity> |
            <RangeValue>:
                lowbound: <quantity>
                highBound: <quantity>
      weight: <integer>
      keyObjective: <boolean>
  totalScore:
    passPercentage: 90
    warningPercentage: 75
```

## Fields

* **apiVersion** -- API version being used
* **kind** -- Resource type.
   Must be set to `AnalysisDefinition`.
* **metadata**
  * **name** -- Unique name of this analysis definition.
    Names must comply with the
    [Kubernetes Object Names and IDs](https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-subdomain-names)
    specification.
  * **namespace** -- Namespace where this resource is located.
    `Analysis` resources must specify this namespace
    when referencing this definition,
    unless it resides in the same namespace as the `Analysis` resource.
* **spec**
  * **objectives**
    This is a list of objectives whose results are combined
    to determine whether the analysis fails, passes, or passes with a warning.
    * **analysisValueTemplateRef** (required) --
      This string marks the beginning of each objective
      * **name** (required) -- The `metadata.name` value of the
      [AnalysisDefinition](analysisdefinition.md)
      resource used for this objective.
      That resource defines the data provider and the query to use.
      * **namespace** --
        Namespace of the `analysisValueTemplateRef` resource.
        If the namespace is not specified,
        the analysis controller looks for the `AnalysisValueTemplateRef` resource
        in the same namespace as the `Analysis` resource.
      * **target** -- defines failure or, optionally, warning criteria.
        Values not specified for failure or warning result in a pass.
        Keptn writes the results of the analysis to the `status` section
        of the
        [Analysis](analysis.md)
        resource after the analysis runs.
        * **failure** -- criteria for failure, specified as
          `operator: <quantity>`.
          This can be specified either as an absolute value
          or as a range of values.

          Valid operators for absolute values are:
          * `lessThan` -- `<` operator
          * `lessThanOrEqual` -- `<=` operator
          * `greaterThan` -- `>` operator
          * `greaterThanOrEqual` -- `>=` operator
          * `equalTo` -- `==` operator

          Valid operators for specifying ranges are:
          * `inRange` -- value is inclusively in the defined range
          * `notInRange` --  value is exclusively out of the defined range

            Each of these operators require two arguments:

            * `lowBound` -- minimum value of the range included or excluded
            * `highBound` -- maximum value of the range included or excluded
        <!-- markdownlint-disable -->
        * **warning** -- criteria for a warning,
          specified in the same way as the `failure` field.
      * **weight**  -- used to emphasize the importance
        of one `objective` over others
      * **keyObjective** -- If set to `true`,
        the entire analysis fails if this objective fails
  * **totalScore** (required) --
    * **passPercentage** -- threshold to reach for the full analysis
      (all objectives) to pass
    <!-- markdownlint-disable -->
    * **warningPercentage** -- threshold to reach
      for the full analysis (all objectives) to pass with  `warning` status

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

## Example

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
          <operator>:
            fixedValue: integer> |
            inRange: | notInRange:
              lowBound: <integer-quantity>
              highBound: <integer-quantity>
        warning:
          <operator>:
            fixedValue: integer> |
            inRange: | notInRange:
              lowBound: <integer-quantity>
              highBound: <integer-quantity>
      weight: <integer>
      keyObjective: <boolean>
  totalScore:
    passPercentage: <integer-percentage>
    warningPercentage: <integer-percentage>
```

For an example of how to implement the Keptn Analysis feature, see the
[Analysis](../implementing/slo.md)
guide page.

## Files

API reference:
[AnalysisDefinition](../crd-ref/metrics/v1alpha3/#analysisdefinition)

## Differences between versions

A preliminary release of the Keptn Analysis feature
is included in Keptn v.0.8.3 but is hidden behind a feature flag.
To preview these features, set the environment `ENABLE_ANALYSIS` to `true`
in the `metrics-operator` deployment.

## See also

* [Analysis](analysis.md)
* [AnalysisValueTemplate](analysisvaluetemplate.md)
* [Analysis](../implementing/slo.md) guide
