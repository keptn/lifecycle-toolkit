---
title: v1beta1
description: Reference information for metrics.keptn.sh/v1beta1
---
<!-- markdownlint-disable -->

## Packages
- [metrics.keptn.sh/v1beta1](#metricskeptnshv1beta1)


## metrics.keptn.sh/v1beta1

Package v1beta1 contains API Schema definitions for the metrics v1beta1 API group

### Resource Types
- [Analysis](#analysis)
- [AnalysisDefinition](#analysisdefinition)
- [AnalysisDefinitionList](#analysisdefinitionlist)
- [AnalysisList](#analysislist)
- [AnalysisValueTemplate](#analysisvaluetemplate)
- [AnalysisValueTemplateList](#analysisvaluetemplatelist)
- [KeptnMetric](#keptnmetric)
- [KeptnMetricList](#keptnmetriclist)
- [KeptnMetricsProvider](#keptnmetricsprovider)
- [KeptnMetricsProviderList](#keptnmetricsproviderlist)



#### Analysis



Analysis is the Schema for the analyses API

_Appears in:_
- [AnalysisList](#analysislist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `metrics.keptn.sh/v1beta1`
| `kind` _string_ | `Analysis`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[AnalysisSpec](#analysisspec)_ |  |
| `status` _[AnalysisStatus](#analysisstatus)_ |  |


#### AnalysisDefinition



AnalysisDefinition is the Schema for the analysisdefinitions APIs

_Appears in:_
- [AnalysisDefinitionList](#analysisdefinitionlist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `metrics.keptn.sh/v1beta1`
| `kind` _string_ | `AnalysisDefinition`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[AnalysisDefinitionSpec](#analysisdefinitionspec)_ |  |
| `status` _string_ | unused field |


#### AnalysisDefinitionList



AnalysisDefinitionList contains a list of AnalysisDefinition resources



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `metrics.keptn.sh/v1beta1`
| `kind` _string_ | `AnalysisDefinitionList`
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[AnalysisDefinition](#analysisdefinition) array_ |  |


#### AnalysisDefinitionSpec



AnalysisDefinitionSpec defines the desired state of AnalysisDefinition

_Appears in:_
- [AnalysisDefinition](#analysisdefinition)

| Field | Description |
| --- | --- |
| `objectives` _[Objective](#objective) array_ | Objectives defines a list of objectives to evaluate for an analysis |
| `totalScore` _[TotalScore](#totalscore)_ | TotalScore defines the required score for an analysis to be successful |


#### AnalysisList



AnalysisList contains a list of Analysis resources



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `metrics.keptn.sh/v1beta1`
| `kind` _string_ | `AnalysisList`
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[Analysis](#analysis) array_ |  |


#### AnalysisSpec



AnalysisSpec defines the desired state of Analysis

_Appears in:_
- [Analysis](#analysis)

| Field | Description |
| --- | --- |
| `timeframe` _[Timeframe](#timeframe)_ | Timeframe specifies the range for the corresponding query in the AnalysisValueTemplate. Please note that either a combination of 'from' and 'to' or the 'recent' property may be set. If neither is set, the Analysis can not be added to the cluster. |
| `args` _object (keys:string, values:string)_ | Args corresponds to a map of key/value pairs that can be used to substitute placeholders in the AnalysisValueTemplate query. i.e. for args foo:bar the query could be "query:percentile(95)?scope=tag(my_foo_label:{{.foo}})". |
| `analysisDefinition` _[ObjectReference](#objectreference)_ | AnalysisDefinition refers to the AnalysisDefinition, a CRD that stores the AnalysisValuesTemplates |


#### AnalysisState

_Underlying type:_ `string`

AnalysisState represents the state of the analysis

_Appears in:_
- [AnalysisStatus](#analysisstatus)



#### AnalysisStatus



AnalysisStatus stores the status of the overall analysis returns also pass or warnings

_Appears in:_
- [Analysis](#analysis)

| Field | Description |
| --- | --- |
| `timeframe` _[Timeframe](#timeframe)_ | Timeframe describes the time frame which is evaluated by the Analysis |
| `raw` _string_ | Raw contains the raw result of the SLO computation |
| `pass` _boolean_ | Pass returns whether the SLO is satisfied |
| `warning` _boolean_ | Warning returns whether the analysis returned a warning |
| `state` _[AnalysisState](#analysisstate)_ | State describes the current state of the Analysis (Pending/Progressing/Completed) |
| `storedValues` _object (keys:string, values:[ProviderResult](#providerresult))_ | StoredValues contains all analysis values that have already been retrieved successfully |


#### AnalysisValueTemplate



AnalysisValueTemplate is the Schema for the analysisvaluetemplates API

_Appears in:_
- [AnalysisValueTemplateList](#analysisvaluetemplatelist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `metrics.keptn.sh/v1beta1`
| `kind` _string_ | `AnalysisValueTemplate`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[AnalysisValueTemplateSpec](#analysisvaluetemplatespec)_ | Spec contains the specification for the AnalysisValueTemplate |
| `status` _string_ | unused field |


#### AnalysisValueTemplateList



AnalysisValueTemplateList contains a list of AnalysisValueTemplate resources



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `metrics.keptn.sh/v1beta1`
| `kind` _string_ | `AnalysisValueTemplateList`
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[AnalysisValueTemplate](#analysisvaluetemplate) array_ |  |


#### AnalysisValueTemplateSpec



AnalysisValueTemplateSpec defines the desired state of AnalysisValueTemplate

_Appears in:_
- [AnalysisValueTemplate](#analysisvaluetemplate)

| Field | Description |
| --- | --- |
| `provider` _[ObjectReference](#objectreference)_ | Provider refers to the KeptnMetricsProvider which should be used to retrieve the data |
| `query` _string_ | Query represents the query to be run. It can include placeholders that are defined using the go template syntax. More info on go templating - https://pkg.go.dev/text/template |


#### IntervalResult





_Appears in:_
- [KeptnMetricStatus](#keptnmetricstatus)

| Field | Description |
| --- | --- |
| `value` _string_ | Value represents the resulting value |
| `range` _[RangeSpec](#rangespec)_ | Range represents the time range for which this data was queried |
| `lastUpdated` _[Time](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#time-v1-meta)_ | LastUpdated represents the time when the status data was last updated |
| `errMsg` _string_ | ErrMsg represents the error details when the query could not be evaluated |


#### KeptnMetric



KeptnMetric is the Schema for the keptnmetrics API

_Appears in:_
- [KeptnMetricList](#keptnmetriclist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `metrics.keptn.sh/v1beta1`
| `kind` _string_ | `KeptnMetric`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[KeptnMetricSpec](#keptnmetricspec)_ |  |
| `status` _[KeptnMetricStatus](#keptnmetricstatus)_ |  |


#### KeptnMetricList



KeptnMetricList contains a list of KeptnMetric resources



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `metrics.keptn.sh/v1beta1`
| `kind` _string_ | `KeptnMetricList`
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[KeptnMetric](#keptnmetric) array_ |  |


#### KeptnMetricSpec



KeptnMetricSpec defines the desired state of KeptnMetric

_Appears in:_
- [KeptnMetric](#keptnmetric)

| Field | Description |
| --- | --- |
| `provider` _[ProviderRef](#providerref)_ | Provider represents the provider object |
| `query` _string_ | Query represents the query to be run |
| `fetchIntervalSeconds` _integer_ | FetchIntervalSeconds represents the update frequency in seconds that is used to update the metric |
| `range` _[RangeSpec](#rangespec)_ | Range represents the time range for which data is to be queried |


#### KeptnMetricStatus



KeptnMetricStatus defines the observed state of KeptnMetric

_Appears in:_
- [KeptnMetric](#keptnmetric)

| Field | Description |
| --- | --- |
| `value` _string_ | Value represents the resulting value |
| `rawValue` _integer array_ | RawValue represents the resulting value in raw format |
| `lastUpdated` _[Time](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#time-v1-meta)_ | LastUpdated represents the time when the status data was last updated |
| `errMsg` _string_ | ErrMsg represents the error details when the query could not be evaluated |
| `intervalResults` _[IntervalResult](#intervalresult) array_ | IntervalResults contain a slice of all the interval results |


#### KeptnMetricsProvider



KeptnMetricsProvider is the Schema for the keptnmetricsproviders API

_Appears in:_
- [KeptnMetricsProviderList](#keptnmetricsproviderlist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `metrics.keptn.sh/v1beta1`
| `kind` _string_ | `KeptnMetricsProvider`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[KeptnMetricsProviderSpec](#keptnmetricsproviderspec)_ |  |
| `status` _string_ | unused field |


#### KeptnMetricsProviderList



KeptnMetricsProviderList contains a list of KeptnMetricsProvider resources



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `metrics.keptn.sh/v1beta1`
| `kind` _string_ | `KeptnMetricsProviderList`
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[KeptnMetricsProvider](#keptnmetricsprovider) array_ |  |


#### KeptnMetricsProviderSpec



KeptnMetricsProviderSpec defines the desired state of KeptnMetricsProvider

_Appears in:_
- [KeptnMetricsProvider](#keptnmetricsprovider)

| Field | Description |
| --- | --- |
| `type` _string_ | Type represents the provider type. This can be one of prometheus, dynatrace, datadog, dql. |
| `targetServer` _string_ | TargetServer defines URL (including port and protocol) at which the metrics provider is reachable. |
| `secretKeyRef` _[SecretKeySelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#secretkeyselector-v1-core)_ | SecretKeyRef defines an optional secret for access credentials to the metrics provider. |


#### ObjectReference





_Appears in:_
- [AnalysisSpec](#analysisspec)
- [AnalysisValueTemplateSpec](#analysisvaluetemplatespec)
- [Objective](#objective)
- [ProviderResult](#providerresult)

| Field | Description |
| --- | --- |
| `name` _string_ | Name defines the name of the referenced object |
| `namespace` _string_ | Namespace defines the namespace of the referenced object |


#### Objective



Objective defines an objective for analysis

_Appears in:_
- [AnalysisDefinitionSpec](#analysisdefinitionspec)

| Field | Description |
| --- | --- |
| `analysisValueTemplateRef` _[ObjectReference](#objectreference)_ | AnalysisValueTemplateRef refers to the appropriate AnalysisValueTemplate |
| `target` _[Target](#target)_ | Target defines failure or warning criteria |
| `weight` _integer_ | Weight can be used to emphasize the importance of one Objective over the others |
| `keyObjective` _boolean_ | KeyObjective defines whether the whole analysis fails when this objective's target is not met |


#### Operator



Operator specifies the supported operators for value comparisons

_Appears in:_
- [Target](#target)

| Field | Description |
| --- | --- |
| `lessThanOrEqual` _[OperatorValue](#operatorvalue)_ | LessThanOrEqual represents '<=' operator |
| `lessThan` _[OperatorValue](#operatorvalue)_ | LessThan represents '<' operator |
| `greaterThan` _[OperatorValue](#operatorvalue)_ | GreaterThan represents '>' operator |
| `greaterThanOrEqual` _[OperatorValue](#operatorvalue)_ | GreaterThanOrEqual represents '>=' operator |
| `equalTo` _[OperatorValue](#operatorvalue)_ | EqualTo represents '==' operator |
| `inRange` _[RangeValue](#rangevalue)_ | InRange represents operator checking the value is inclusively in the defined range, e.g. 2 <= x <= 5 |
| `notInRange` _[RangeValue](#rangevalue)_ | NotInRange represents operator checking the value is exclusively out of the defined range, e.g. x < 2 AND x > 5 |


#### OperatorValue



OperatorValue represents the value to which the result is compared

_Appears in:_
- [Operator](#operator)

| Field | Description |
| --- | --- |
| `fixedValue` _Quantity_ | FixedValue defines the value for comparison |


#### ProviderRef



ProviderRef represents the provider object

_Appears in:_
- [KeptnMetricSpec](#keptnmetricspec)

| Field | Description |
| --- | --- |
| `name` _string_ | Name of the provider |


#### ProviderResult



ProviderResult stores reference of already collected provider query associated to its objective template

_Appears in:_
- [AnalysisStatus](#analysisstatus)

| Field | Description |
| --- | --- |
| `objectiveReference` _[ObjectReference](#objectreference)_ | Objective store reference to corresponding objective template |
| `query` _string_ | Query represents the executed query |
| `value` _string_ | Value is the value the provider returned |
| `errMsg` _string_ | ErrMsg stores any possible error at retrieval time |


#### RangeSpec



RangeSpec defines the time range for which data is to be queried

_Appears in:_
- [IntervalResult](#intervalresult)
- [KeptnMetricSpec](#keptnmetricspec)

| Field | Description |
| --- | --- |
| `interval` _string_ | Interval specifies the duration of the time interval for the data query |
| `step` _string_ | Step represents the query resolution step width for the data query |
| `aggregation` _string_ | Aggregation defines the type of aggregation function to be applied on the data. Accepted values: p90, p95, p99, max, min, avg, median |
| `storedResults` _integer_ | StoredResults indicates the upper limit of how many past results should be stored in the status of a KeptnMetric |


#### RangeValue



RangeValue represents a range which the value should fit

_Appears in:_
- [Operator](#operator)

| Field | Description |
| --- | --- |
| `lowBound` _Quantity_ | LowBound defines the lower bound of the range |
| `highBound` _Quantity_ | HighBound defines the higher bound of the range |


#### Target



Target defines the failure and warning criteria

_Appears in:_
- [Objective](#objective)

| Field | Description |
| --- | --- |
| `failure` _[Operator](#operator)_ | Failure defines limits up to which an analysis fails |
| `warning` _[Operator](#operator)_ | Warning defines limits where the result does not pass or fail |


#### Timeframe





_Appears in:_
- [AnalysisSpec](#analysisspec)
- [AnalysisStatus](#analysisstatus)

| Field | Description |
| --- | --- |
| `from` _[Time](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#time-v1-meta)_ | From is the time of start for the query. This field follows RFC3339 time format |
| `to` _[Time](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#time-v1-meta)_ | To is the time of end for the query. This field follows RFC3339 time format |
| `recent` _[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#duration-v1-meta)_ | Recent describes a recent timeframe using a duration string. E.g. Setting this to '5m' provides an Analysis for the last five minutes |


#### TotalScore



TotalScore defines the required score for an analysis to be successful

_Appears in:_
- [AnalysisDefinitionSpec](#analysisdefinitionspec)

| Field | Description |
| --- | --- |
| `passPercentage` _integer_ | PassPercentage defines the threshold to reach for an analysis to pass |
| `warningPercentage` _integer_ | WarningPercentage defines the threshold to reach for an analysis to pass with a 'warning' status |


