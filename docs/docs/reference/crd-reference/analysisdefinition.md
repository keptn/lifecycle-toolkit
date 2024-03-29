---
comments: true
---

# AnalysisDefinition

An `AnalysisDefinition` resource defines the
list of Service Level Objectives (SLOs) for an `Analysis`.

## Synopsis

```yaml
apiVersion: metrics.keptn.sh/v1
kind: AnalysisDefinition
metadata:
  name: <name-of-this-resource>
  namespace: <namespace-where-this-resource-resides>
spec:
  objectives:
    - analysisValueTemplateRef:
        name: <name-of-referenced-analysisValueTemplate-resource>
        namespace: <namespace-of-referenced-analysisValueTemplate-resource>
      target:
        failure:
          <operator>:
            fixedValue: <integer> | <quantity>
          inRange: | notInRange:
            lowBound: <integer> | <quantity>
            highBound: <integer> | <quantity>
        warning:
          <operator>:
            fixedValue: <integer> | <quantity>
          inRange: | notInRange:
            lowBound: <integer> | <quantity>
            highBound: <integer> | <quantity>
      weight: <integer>
      keyObjective: <boolean>
  totalScore:
    passPercentage: <min-percentage-to-pass>
    warningPercentage: <min-percentage-for-warning>
```

## Fields

- **apiVersion** -- API version being used
- **kind** -- Resource type.
  Must be set to `AnalysisDefinition`.
    - **metadata**
        - **name** -- Unique name of this analysis definition.
          Names must comply with the
          [Kubernetes Object Names and IDs](https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-subdomain-names)
          specification.
        - **namespace** -- Namespace where this resource is located.
          `Analysis` resources must specify this namespace
          when referencing this definition,
          unless it resides in the same namespace as the `Analysis` resource.
- **spec**
    - **objectives**
      A list of objectives whose results are combined
      to determine whether the analysis fails, passes, or passes with a warning.
        - **analysisValueTemplateRef** (required) --
          This string marks the beginning of each objective
            - **name** (required) -- The `metadata.name` value of the
              [AnalysisValueTemplateRef](analysisvaluetemplate.md)
              resource that defines the SLI used for this objective.
              That resource defines the data provider and the query to use.
            - **namespace** --
              Namespace of the `analysisValueTemplateRef` resource.
              If the namespace is not specified,
              the analysis controller looks for the `AnalysisValueTemplateRef` resource
              in the same namespace as the `Analysis` resource.

        - **target** -- defines failure or, optionally, warning criteria.
          Values not specified for failure or warning result in a pass.
          Keptn writes the results of the analysis to the `status` section
          of the
          [Analysis](analysis.md)
          resource after the analysis runs.

            To use a value that includes a fraction, use a Kubernetes
            [quantity](https://kubernetes.io/docs/reference/kubernetes-api/common-definitions/quantity/)
            value rather than a `float`.
            For example, use the `3m` quantity rather than the `0.003` float;
            a `float` value here causes `Invalid value` errors.
            A whole number (integer) is also a legal `quantity` value.

            - **failure** -- criteria for failure, specified as
              `operator: <quantity>`.
              This can be specified either as an absolute `quantity` value
              or as a range of values.

                Valid operators for absolute values are:

                - `lessThan` -- `<` operator
                - `lessThanOrEqual` -- `<=` operator
                - `greaterThan` -- `>` operator
                - `greaterThanOrEqual` -- `>=` operator
                - `equalTo` -- `==` operator

                Valid operators for specifying ranges are:

                - `inRange` -- value is inclusively in the defined range
                - `notInRange` -- value is exclusively out of the defined range

                Each of these operators require two arguments:

                - `lowBound` -- minimum `quantity` value of the range included or excluded
                - `highBound` -- maximum `quantity` value of the range included or excluded

            - **warning** -- criteria for a warning, specified in the same way as the `failure` field.

        - **weight** -- used to emphasize the importance of one `objective` over others
        - **keyObjective** -- If set to `true`, the entire analysis fails if this particular objective fails,
          no matter what the actual `score` of the analysis is

    - **totalScore** (required) --
        - **passPercentage** -- threshold to reach for the full analysis (all objectives) to pass
        - **warning** Percentage
          for the full analysis (all objectives) to pass with  `warning` status

## Usage

An `AnalysisDefinition` resource contains a list of objectives to satisfy.
Each of these objectives must specify:

- The `AnalysisValueTemplate` resource that contains the SLIs,
  defining the data provider from which to gather the data
  and how to compute the Analysis
- Failure or warning target criteria
- Whether the objective is a key objective
  meaning that its failure fails the Analysis
- Weight of the objective on the overall Analysis

## Example

```yaml
{% include "https://raw.githubusercontent.com/keptn/lifecycle-toolkit/main/metrics-operator/config/samples/metrics_v1_analysisdefinition.yaml" %}
```

For a full example of how to implement the Keptn Analysis feature, see the
[Analysis](../../guides/slo.md)
guide page.

## Files

API reference:
[AnalysisDefinition](../api-reference/metrics/v1/index.md#analysisdefinition)

## Differences between versions

The Keptn Analysis feature is an official part of Keptn v0.10.0 and later.
Keptn v0.8.3 included a preliminary release of this feature
but it was hidden behind a feature flag.
The behavior of this feature is unchanged since v0.8.3.

## See also

- [Analysis](analysis.md)
- [AnalysisValueTemplate](analysisvaluetemplate.md)
- [Analysis](../../guides/slo.md) guide
