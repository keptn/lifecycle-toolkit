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
- [KeptnMetric](#keptnmetric)
- [KeptnMetricList](#keptnmetriclist)
- [KeptnMetricsProvider](#keptnmetricsprovider)
- [KeptnMetricsProviderList](#keptnmetricsproviderlist)



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


