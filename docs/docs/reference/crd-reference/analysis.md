---
comments: true
---

# Analysis

An `Analysis` is a snapshot of your current application status.
Based on your defined SLIs, it can validate that your SLOs are satisfied,
using the data coming from your defined set of `KeptnMetricsProvider` resources.

The `Analysis` resource is an instance of an
[AnalysisDefinition](analysisdefinition.md) resource
which defines specific data like
the time for which the analysis should be done
and the appropriate values to use for variables
that are used in the `AnalysisDefinition` query.

## Synopsis

```yaml
apiVersion: metrics.keptn.sh/v1
kind: Analysis
metadata:
  name: <name-of-analysis>
spec:
  timeframe: from: <start-time> to: <end-time> | `recent <timespan>`
  args:
    <variable1>: <value1>
    <variable2>: <value2>
    ...
  analysisDefinition:
    name: <name of associated `analysisDefinition` resource
    namespace: <namespace of associated `analysisDefinition` resource
status:
  pass: true | false
  warning: true | false
  raw: <JSON object>
  state: Completed | Progressing
```

## Fields

- **apiVersion** -- API version being used
- **kind** -- Resource type.
   Must be set to `Analysis`
- **metadata**
    - **name** -- Unique name of this analysis.
       Names must comply with the
       [Kubernetes Object Names and IDs](https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-subdomain-names)
       specification.
- **spec**
    - **timeframe** (required) -- Specifies the range  for the corresponding query
       in the AnalysisValueTemplate.
       This can be populated as one of the following:

        - A combination of ‘from’ and ‘to’
            to specify the start and stop times for the analysis.
            These fields follow the
            [RFC 3339](https://www.ietf.org/rfc/rfc3339.txt)
            timestamp format.
        - Set the ‘recent’ property to a time span.
            This causes the Analysis to use data going back that amount of time.
            For example, if `recent: 10m` is set,
            the Analysis studies data from the last ten minutes.

        If neither is set, the Analysis can not be added to the cluster.

    - **args** -- Map of key/value pairs that can be used
       to substitute variables in the `AnalysisValueTemplate` query.
    - **analysisDefinition** (required) -- Identify the `AnalysisDefinition` resource
       that stores the `AnalysisValuesTemplate` associated with this `Analysis`
        - **name** -- Name of the `AnalysisDefinition` resource
        - **namespace** (optional) --
            Namespace of the `AnalysisDefinition` resource.
            The `AnalysisDefinition` resource can be located in any namespace.
            If the namespace is not specified,
            the analysis controller looks for the `AnalysisDefinition` resource
            in the same namespace as the `Analysis` resource.
- **status** -- results of this Analysis run,
   added to the resource by Keptn,
   based on criteria defined in the `AnalysisDefinition` resource.

    - **warning** -- Whether the analysis returned a warning.
    - **raw** --  String-encoded JSON object that reports the results
        of evaluating one or more objectives or metrics.
        See
        [Interpreting Analysis results](#interpreting-analysis-results)
        for details.
    - **state** -- Set to `Completed` or `Progressing` as appropriate.

## Interpreting Analysis results

The `status.raw` field is a string encoded JSON object object that represents the
results of evaluating one or more performance objectives or metrics.
It shows whether these objectives have passed or failed, their actual values, and the associated scores.
In this example, the objectives include response time and error rate analysis,
each with its own criteria for passing or failing.
The overall evaluation has passed, and no warnings have been issued.

> Note: Please check the inline annotations to get more information about the particular lines you are interested in.

```json
{
    "objectiveResults"/*(1)!*/: [
        {
            "result"/*(2)!*/: {
                "failResult"/*(3)!*/: {
                    "operator": {
                        "greaterThan": {
                            "fixedValue": "500m"
                        }
                    },
                    "fulfilled": false
                },
                "warnResult"/*(4)!*/: {
                    "operator": {
                        "greaterThan": {
                            "fixedValue": "300m"
                        }
                    },
                    "fulfilled": false
                },
                "warning"/*(5)!*/: false,
                "pass"/*(6)!*/: true
            },
            "objective"/*(7)!*/: {
                "analysisValueTemplateRef"/*(8)!*/: {
                    "name": "response-time-p95"
                },
                "target"/*(9)!*/: {
                    "failure": {
                        "greaterThan": {
                            "fixedValue": "500m"
                        }
                    },
                    "warning" : {
                        "greaterThan": {
                            "fixedValue": "300m"
                        }
                    }
                },
                "weight"/*(10)!*/: 1
            },
            "value"/*(11)!*/: 0.00475,
            "score"/*(12)!*/: 1
        },
        {
            "result"/*(13)!*/: {
                "failResult": {
                    "operator": {
                        "greaterThan": {
                            "fixedValue": "0"
                        }
                    },
                    "fulfilled": false
                },
                "warnResult": {
                    "operator": {

                    },
                    "fulfilled": false
                },
                "warning": false,
                "pass": true
            },
            "objective"/*(14)!*/: {
                "analysisValueTemplateRef"/*(15)!*/: {
                    "name": "error-rate"
                },
                "target"/*(16)!*/: {
                    "failure": {
                        "greaterThan": {
                            "fixedValue": "0"
                        }
                    }
                },
                "weight"/*(17)!*/: 1,
                "keyObjective"/*(18)!*/: true
            },
            "value"/*(19)!*/: 0,
            "score"/*(20)!*/: 1
        }
    ],
    "totalScore"/*(21)!*/: 2,
    "maximumScore"/*(22)!*/: 2,
    "pass"/*(23)!*/: true,
    "warning"/*(24)!*/: false
}
```

1. **`objectiveResults`**: This is an array containing one or more objects,
    each representing the results of a specific objective or performance metric.
2. **`result`** -- This object contains information about whether the objective has passed or failed.
    It has two sub-objects **`failResult`** & **`warnResult`**
3. **`failResult`** -- Indicates whether the objective has failed.
    In this case, it checks if a value is greater than 500 milliseconds and it has not been fulfilled (`fulfilled: false`).
4. **`warnResult`** -- Indicates whether the objective has issued a warning.
    It checks if a value is greater than 300 milliseconds
5. **`warning`** (false in this case).
6. **`pass`** -- Indicates whether the objective has passed (true in this case).
7. **`objective`** -- Describes the objective being evaluated.
    It includes: **`analysisValueTemplateRef`** , **`target`** & **`weight`**
8. **`analysisValueTemplateRef`** -- Refers to the template used for analysis (`response-time-p95`).
9. **`target`** -- Sets the target value for failure (failure occurs if the value is greater than 0).
    In this case, failure occurs
    if the value is greater than 500 milliseconds and warning occurs if it's greater than 300 milliseconds.
10. **`weight`** -- Specifies the weight assigned to this objective (weight: 1).
11. **`value`** -- Indicates the actual value measured for this objective (value: 0.00475).
12. **`score`** -- Indicates the score assigned to this objective (score: 1).
13. **`result`** -- Similar to the first objective,
    it checks whether a value is greater than 0 and has not been fulfilled (`fulfilled: false`).
    There are no warning conditions in this case.
14. **`objective`** -- Describes the objective related to error rate analysis.
15. **`analysisValueTemplateRef`** -- Refers to the template used for analysis (`error-rate`).
16. **`target`** -- Sets the target value for failure (failure occurs if the value is greater than 0).
17. **`weight`** -- Specifies the weight assigned to this objective (weight: 1).
18. **`keyObjective`** -- Indicates that this is a key objective (true).
19. **`value`** -- Indicates the actual value measured for this objective (value: 0).
20. **`score`** -- Indicates the score assigned to this objective (score: 1).
21. **`totalScore`** -- Represents the total score achieved based on the objectives evaluated (totalScore: 2).
22. **`maximumScore`** -- Indicates the maximum possible score (maximumScore: 2).
23. **`pass`** -- Indicates whether the overall evaluation has passed (true in this case).
24. **`warning`** -- Indicates whether any warnings have been issued during the evaluation (false in this case).

## Usage

An `Analysis` resource specifies a single Analysis run.
It specifies the `AnalysisValueTemplate` resource
that defines the calculations to use,
the timeframe for which to report information,
and values to use for variables for this run.

The result of this analysis stays in the cluster
until the `Analysis` is deleted.
That also means that, if another analysis should be performed,
the new analysis must be given a new, unique name within the namespace.

To perform an Analysis (or "trigger an evaluation" in Keptn v1 jargon),
apply the `analysis-instance.yaml` file:

```shell
kubectl apply -f analysis-instance.yaml -n keptn-lifecycle-poc
```

Retrieve the current status of the Analysis with the following command:

```shell
kubectl get analysis - n keptn-lifecycle-poc
```

This yields an output that looks like the following:

```shell
NAME                ANALYSISDEFINITION      WARNING   PASS
analysis-sample-1   my-project-ad                     true
```

This shows that the analysis passed successfully.

To get the detailed result of the evaluation,
use the `-oyaml` argument to inspect the full state of the analysis:

This displays the `Analysis` resource
with the definition of the analysis
as well as the `status` (results) of the analysis; for example:

```shell
kubectl get analysis - n keptn-lifecycle-poc -oyaml
```

## Examples

```yaml
{% include "../../assets/crd/analysis.yaml" %}
```

This `Analysis` resource:

- Defines the `timeframe` for which the analysis is done
  as between 5 am and 10 am on the 5th of May 2023
- Adds a few specific key-value pairs that will be substituted in the query.
  For instance, the query could contain the `{{.nodename}}` variable.
  The value of the `args.nodename` field (`test`)
  will be substituted for this string.

For a full example of how to implement the Keptn Analysis feature, see the
[Analysis](../../guides/slo.md)
guide page.

## Files

API reference: [Analysis](../api-reference/metrics/v1/index.md#analysis)

## Differences between versions

The Keptn Analysis feature is an official part of Keptn v0.10.0 and later.
Keptn v0.8.3 included a preliminary release of this feature
but it was hidden behind a feature flag.
The behavior of this feature is unchanged since v0.8.3.

## See also

- [AnalysisDefinition](analysisdefinition.md)
- [AnalysisValueTemplate](analysisvaluetemplate.md)
- [Analysis](../../guides/slo.md) guide
