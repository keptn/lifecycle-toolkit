---
title: v1alpha3
description: Reference information for metrics.keptn.sh/v1alpha3
---
<!-- markdownlint-disable -->

## Packages
- [metrics.keptn.sh/v1alpha3](#metricskeptnshv1alpha3)


## metrics.keptn.sh/v1alpha3

Package v1alpha3 contains API Schema definitions for the metrics v1alpha3 API group

### Resource Types
- [AnalysisDefinition](#analysisdefinition)
- [AnalysisDefinitionList](#analysisdefinitionlist)
- [KeptnMetric](#keptnmetric)
- [KeptnMetricList](#keptnmetriclist)
- [KeptnMetricsProvider](#keptnmetricsprovider)
- [KeptnMetricsProviderList](#keptnmetricsproviderlist)



#### AnalysisDefinition



AnalysisDefinition is the Schema for the analysisdefinitions API

_Appears in:_
- [AnalysisDefinitionList](#analysisdefinitionlist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `metrics.keptn.sh/v1alpha3`
| `kind` _string_ | `AnalysisDefinition`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[AnalysisDefinitionSpec](#analysisdefinitionspec)_ |  |
| `status` _[AnalysisDefinitionStatus](#analysisdefinitionstatus)_ |  |


#### AnalysisDefinitionList



AnalysisDefinitionList contains a list of AnalysisDefinition



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `metrics.keptn.sh/v1alpha3`
| `kind` _string_ | `AnalysisDefinitionList`
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[AnalysisDefinition](#analysisdefinition) array_ |  |


#### AnalysisDefinitionSpec



AnalysisDefinitionSpec defines the desired state of AnalysisDefinition

_Appears in:_
- [AnalysisDefinition](#analysisdefinition)

| Field | Description |
| --- | --- |
| `objective` _[Objective](#objective) array_ |  |
| `totalScore` _[Score](#score)_ |  |




#### Criteria





_Appears in:_
- [CriteriaSet](#criteriaset)



#### CriteriaSet





_Appears in:_
- [SLOTarget](#slotarget)

| Field | Description |
| --- | --- |
| `anyOf` _[Criteria](#criteria) array_ | AnyOf contains a list of targets where any of them needs to be successful for the Criteria to pass |
| `allOf` _[Criteria](#criteria) array_ | AllOf contains a list of targets where all of them need to be successful for the Criteria to pass |


#### KeptnMetric



KeptnMetric is the Schema for the keptnmetrics API

_Appears in:_
- [KeptnMetricList](#keptnmetriclist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `metrics.keptn.sh/v1alpha3`
| `kind` _string_ | `KeptnMetric`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[KeptnMetricSpec](#keptnmetricspec)_ |  |
| `status` _[KeptnMetricStatus](#keptnmetricstatus)_ |  |


#### KeptnMetricList



KeptnMetricList contains a list of KeptnMetric



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `metrics.keptn.sh/v1alpha3`
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


#### KeptnMetricsProvider



KeptnMetricsProvider is the Schema for the keptnmetricsproviders API

_Appears in:_
- [KeptnMetricsProviderList](#keptnmetricsproviderlist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `metrics.keptn.sh/v1alpha3`
| `kind` _string_ | `KeptnMetricsProvider`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[KeptnMetricsProviderSpec](#keptnmetricsproviderspec)_ |  |
| `status` _string_ | unused field |


#### KeptnMetricsProviderList



KeptnMetricsProviderList contains a list of KeptnMetricsProvider



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `metrics.keptn.sh/v1alpha3`
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
| `targetServer` _string_ | TargetServer defined the URL at which the metrics provider is reachable with included port and protocol. |
| `secretKeyRef` _[SecretKeySelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#secretkeyselector-v1-core)_ | SecretKeyRef defines an optional secret for access credentials to the metrics provider. |


#### ObjectReference





_Appears in:_
- [Objective](#objective)

| Field | Description |
| --- | --- |
| `name` _string_ |  |
| `namespace` _string_ |  |


#### Objective



Objective defines a list of objectives

_Appears in:_
- [AnalysisDefinitionSpec](#analysisdefinitionspec)

| Field | Description |
| --- | --- |
| `analysisValueTemplateRef` _[ObjectReference](#objectreference)_ | AnalysisValueTemplateRef defines a reference to the used AnalysisValueTemplate |
| `sloTargets` _[SLOTarget](#slotarget)_ | SLOTargets defines a list of SLOTargests |
| `weight` _integer_ | Weigeht defines the importance of one SLI over the others |
| `keyObjective` _boolean_ | KeyObjective defines the meaning that the analysis fails if the objective is not met |


#### ProviderRef



ProviderRef represents the provider object

_Appears in:_
- [KeptnMetricSpec](#keptnmetricspec)

| Field | Description |
| --- | --- |
| `name` _string_ | Name of the provider |


#### RangeSpec



RangeSpec defines the time range for which data is to be queried

_Appears in:_
- [KeptnMetricSpec](#keptnmetricspec)

| Field | Description |
| --- | --- |
| `interval` _string_ | Interval specifies the duration of the time interval for the data query |
| `step` _string_ | Step represents the query resolution step width for the data query |
| `aggregation` _string_ | Aggregation defines as the type of aggregation function to be applied on the data. Accepted values: p90, p95, p99, max, min, avg, median |


#### SLOTarget



SLOTarget defines the Criteria

_Appears in:_
- [Objective](#objective)

| Field | Description |
| --- | --- |
| `pass` _[CriteriaSet](#criteriaset)_ | Pass defines limit up to which an evaluation is successful |
| `warning` _[CriteriaSet](#criteriaset)_ | Warning defines the border where the result is not pass and not fail |


#### Score



Score defines the required score for an evaluation to be successful

_Appears in:_
- [AnalysisDefinitionSpec](#analysisdefinitionspec)

| Field | Description |
| --- | --- |
| `passPercentage` _integer_ | PassPercentage defines the threshold which needs to be reached for an evaluation to pass. |
| `warningPercentage` _integer_ | WarningPercentage defines the threshold which needs to be reached for an evaluation to pass with a 'warning' status. |




#### TargetValue





_Appears in:_
- [Target](#target)

| Field | Description |
| --- | --- |
| `fixedValue` _integer_ |  |


