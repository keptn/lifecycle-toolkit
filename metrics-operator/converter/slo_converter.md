# SLO Converter

## Description

SLO converter is a tool to convert `slo.yaml` files used in [KeptnV1](https://v1.keptn.sh/) into the new
`AnalysisDefinition` resources used in the new kubernetes-native [Keptn](https://lifecycle.keptn.sh/).
The converter is part of `metrics-operator` image.

## Usage

The converter will convert a single `slo.yaml` file into single `AnalysisDefintion` resource.

To run the converter, execute the following command:

```shell
METRICS_OPERATOR_IMAGE=<METRICS_OPERATOR_IMAGE>
PATH_TO_SLO=<PATH_TO_SLO>
ANALYSIS_VALUE_TEMPLATE_NAMESPACE=<ANALYSIS_VALUE_TEMPLATE_NAMESPACE>
DEFINITION_NAME=<DEFINITION_NAME>

docker-run $METRICS_OPERATOR_IMAGE manager --convert-slo=$PATH_TO_SLO --slo-namespace=$ANALYSIS_VALUE_TEMPLATE_NAMESPACE --definition=$DEFINITION_NAME
```

Please be aware, you need to substitute the placeholders with the following information:

* **PATH_TO_SLO** - path to your `slo.yaml` file
* **ANALYSIS_VALUE_TEMPLATE_NAMESPACE** - namespace of `AnalysisValueTemplate` which will be referenced in objectives
* **DEFINITION_NAME** - name of created `AnalysisDefinition`

> **Note**

All the SLOs present in `slo.yaml` file will reference `AnalysisValueTemplate` resources from the namespace defined
by `ANALYSIS_VALUE_TEMPLATE_NAMESPACE` argument.

## Conversion details

We have multiple use-cases which are and which are not supported.
There is a need to convert the use-cases that make
logical sense and are common, but in some cases, where it is problematic and these cases will not be supported.

> **Note** Please be aware, that comparison criteria containing `%` symbol ware not supported and will be ignored.

### Unsupported use-cases

Criteria with 3 and more inputs won't be supported, only the first 2 non-percentage inputs (those not containing `%`
as we do not support comparison rules) will be taken and converted.
In the example below, only `<600` and `>400` rules will be converted.
Rule `>800` will be ignored.

```yaml
objectives:
- sli: response_time_p95
  displayName: "Response Time P95"
  pass:
  - criteria
    - "<600"
    - ">400"
    - ">800"
  weight: 2
  key_sli: true
```

Secondly, conversion of multiple criteria elements won't be supported at all.
Rules containing multiple criteria elements combined with logical OR (see documentation
[here](https://github.com/keptn/spec/blob/master/service_level_objective.md#objectives))
can not be converted and will be turned into informative items.
Informative in this context means that there are no objectives that have to be met and
therefore the item will not affect the overall score.
However, their values are still retrieved to provide additional insights for the execution of an Analysis.
An example of an unsupported SLOs for conversion is the following:

```yaml
objectives:
- sli: response_time_p95
  displayName: "Response Time P95"
  pass:
  - criteria:
    - "<=+10%"
  - criteria
    - "<600"
  weight: 2
  key_sli: true
```

Thirdly, a case where `pass` criteria are set, `warn` criteria are set, but pass criteria interval and warn
criteria interval do not intercept

```yaml
objectives:
- sli: response_time_p95
  displayName: "Response Time P95"
  pass:
  - criteria:
    - ">200"
    - "<400"
  warn:
  - criteria:
    - ">600"
    - "<800" 
```

### Supported use-cases

The basic objective with a single rule for pass or warning criteria

```yaml
objectives:
- sli: response_time_p95
  displayName: "Response Time P95"
  pass:
  - criteria:
    - ">400"
  warn:
  - criteria:
    - ">200"
```

will be converted to

```yaml
spec:
  objectives:
  - analysisValueTemplateRef:
      name: response_time_p95
      namespace: default
    target:
      failure:
        lessThanOrEqual:
          fixedValue: "200"
      warning:
        lessThanOrEqual:
          fixedValue: "400"
```

The buckets for rules with single criteria element but with one or more criteria combined with logical AND
operator (see documentation [here](https://github.com/keptn/spec/blob/master/service_level_objective.md#objectives)):

1. `pass` criteria set, `warn` criteria not set

```yaml
objectives:
- sli: response_time_p95
  displayName: "Response Time P95"
  pass:
  - criteria:
    - ">400"
    - "<600"
```

will be converted to

```yaml
spec:
  objectives:
  - analysisValueTemplateRef:
      name: response_time_p95
      namespace: default
    target:
      failure:
        notInRange:
          highBound: "600"
          lowBound: "400"
```

1. `pass` criteria set, `warn` criteria set, warn criteria interval is superset of pass criteria interval

```yaml
objectives:
- sli: response_time_p95
  displayName: "Response Time P95"
  pass:
  - criteria:
    - ">400"
    - "<600"
  warn:
  - criteria:
    - ">200"
    - "<800" 
```

will be converted to

```yaml
spec:
  objectives:
  - analysisValueTemplateRef:
      name: response_time_p95
      namespace: default
    target:
      failure:
        notInRange:
          highBound: "200"
          lowBound: "800"
      warning:
        notInRange:
          highBound: "400"
          lowBound: "600"
```

1. `pass` criteria set, `warn` criteria set, pass criteria interval is superset of warn criteria interval

```yaml
objectives:
- sli: response_time_p95
  displayName: "Response Time P95"
  pass:
  - criteria:
    - ">200"
    - "<800"
  warn:
  - criteria:
    - ">400"
    - "<600" 
```

will be converted to

```yaml
spec:
  objectives:
  - analysisValueTemplateRef:
      name: response_time_p95
      namespace: default
    target:
      failure:
        notInRange:
          highBound: "200"
          lowBound: "800"
      warning:
        inRange:
          highBound: "400"
          lowBound: "600"
```

## Example

The following content of a full example of `slo.yaml` file

```yaml
spec_version: "0.1.1"
comparison:
  aggregate_function: "avg"
  compare_with: "single_result"
  include_result_with_score: "pass"
  number_of_comparison_results: 1
objectives:
  - sli: "response_time_p90"
    key_sli: false
    pass:
    - criteria:
        - ">600"
        - "<800"
    warning:
    - criteria:
        - "<=1000"
        - ">500"
    weight: 2
  - sli: "response_time_p80"
    key_sli: false
    pass:
      - criteria:
          - ">600"
          - "<800"
    warning:
      - criteria:
          - "<=1000"
    weight: 2
  - sli: "response_time_p70"
    key_sli: false
    warning:
      - criteria:
          - ">600"
          - "<800"
    pass:
      - criteria:
          - "<=1000"
    weight: 2
  - sli: "response_time_p95"
    key_sli: false
    pass:
      - criteria:
          - "<=+75%"
          - "<800"
    warning:
      - criteria:
          - "<=1000"
          - "<=+100%"
    weight: 1
  - sli: "cpu"
    pass:
      - criteria:
          - "<=+100%"
          - ">=80"
      - criteria:
          - "<=+100%"
          - ">=80"
  - sli: "throughput"
    pass:
      - criteria:
          - "<=+100%"
          - ">=-80%"
  - sli: "error_rate"
total_score:
  pass: "100%"
  warning: "65%"
```

with the following command

```shell
docker-run $METRICS_OPERATOR_IMAGE manager --convert-slo=./slo.yaml --slo-namespace=default --definition=defName
```

will be converted to:

```yaml
apiVersion: metrics.keptn.sh/v1alpha3
kind: AnalysisDefinition
metadata:
  creationTimestamp: null
  name: defName
spec:
  objectives:
  - analysisValueTemplateRef:
      name: response_time_p90
      namespace: default
    target:
      failure:
        notInRange:
          highBound: 1k
          lowBound: "500"
      warning:
        notInRange:
          highBound: "800"
          lowBound: "600"
    weight: 2
  - analysisValueTemplateRef:
      name: response_time_p80
      namespace: default
    target:
      failure:
        greaterThan:
          fixedValue: 1k
      warning:
        notInRange:
          highBound: "800"
          lowBound: "600"
    weight: 2
  - analysisValueTemplateRef:
      name: response_time_p70
      namespace: default
    target:
      failure:
        greaterThan:
          fixedValue: 1k
      warning:
        inRange:
          highBound: "800"
          lowBound: "600"
    weight: 2
  - analysisValueTemplateRef:
      name: response_time_p95
      namespace: default
    target:
      failure:
        greaterThan:
          fixedValue: 1k
      warning:
        greaterThanOrEqual:
          fixedValue: "800"
    weight: 1
  - analysisValueTemplateRef:
      name: cpu
      namespace: default
    target: {}
  - analysisValueTemplateRef:
      name: throughput
      namespace: default
    target: {}
  - analysisValueTemplateRef:
      name: error_rate
      namespace: default
    target: {}
  totalScore:
    passPercentage: 100
    warningPercentage: 65
```
