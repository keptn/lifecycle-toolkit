# v1alpha2

Reference information for metrics.keptn.sh/v1alpha2

<!-- markdownlint-disable -->

## Packages
- [metrics.keptn.sh/v1alpha2](#metricskeptnshv1alpha2)


## metrics.keptn.sh/v1alpha2

Package v1alpha2 contains API Schema definitions for the metrics v1alpha2 API group

### Resource Types
- [KeptnMetric](#keptnmetric)
- [KeptnMetricList](#keptnmetriclist)
- [KeptnMetricsProvider](#keptnmetricsprovider)
- [KeptnMetricsProviderList](#keptnmetricsproviderlist)



#### KeptnMetric



KeptnMetric is the Schema for the keptnmetrics API

_Appears in:_
- [KeptnMetricList](#keptnmetriclist)

| Field | Description | Default | Optional |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `metrics.keptn.sh/v1alpha2` | | |
| `kind` _string_ | `KeptnMetric` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. || ✓ |
| `spec` _[KeptnMetricSpec](#keptnmetricspec)_ |  || ✓ |
| `status` _[KeptnMetricStatus](#keptnmetricstatus)_ |  || ✓ |


#### KeptnMetricList



KeptnMetricList contains a list of KeptnMetric



| Field | Description | Default | Optional |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `metrics.keptn.sh/v1alpha2` | | |
| `kind` _string_ | `KeptnMetricList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. || ✓ |
| `items` _[KeptnMetric](#keptnmetric) array_ |  || x |


#### KeptnMetricSpec



KeptnMetricSpec defines the desired state of KeptnMetric

_Appears in:_
- [KeptnMetric](#keptnmetric)

| Field | Description | Default | Optional |
| --- | --- | --- | --- |
| `provider` _[ProviderRef](#providerref)_ | Provider represents the provider object || x |
| `query` _string_ | Query represents the query to be run || x |
| `fetchIntervalSeconds` _integer_ | FetchIntervalSeconds represents the update frequency in seconds that is used to update the metric || x |


#### KeptnMetricStatus



KeptnMetricStatus defines the observed state of KeptnMetric

_Appears in:_
- [KeptnMetric](#keptnmetric)

| Field | Description | Default | Optional |
| --- | --- | --- | --- |
| `value` _string_ | Value represents the resulting value || x |
| `rawValue` _integer array_ | RawValue represents the resulting value in raw format || x |
| `lastUpdated` _[Time](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#time-v1-meta)_ | LastUpdated represents the time when the status data was last updated || x |


#### KeptnMetricsProvider



KeptnMetricsProvider is the Schema for the keptnmetricsproviders API

_Appears in:_
- [KeptnMetricsProviderList](#keptnmetricsproviderlist)

| Field | Description | Default | Optional |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `metrics.keptn.sh/v1alpha2` | | |
| `kind` _string_ | `KeptnMetricsProvider` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. || ✓ |
| `spec` _[KeptnMetricsProviderSpec](#keptnmetricsproviderspec)_ |  || ✓ |
| `status` _string_ | unused field || ✓ |


#### KeptnMetricsProviderList



KeptnMetricsProviderList contains a list of KeptnMetricsProvider



| Field | Description | Default | Optional |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `metrics.keptn.sh/v1alpha2` | | |
| `kind` _string_ | `KeptnMetricsProviderList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. || ✓ |
| `items` _[KeptnMetricsProvider](#keptnmetricsprovider) array_ |  || x |


#### KeptnMetricsProviderSpec



KeptnMetricsProviderSpec defines the desired state of KeptnMetricsProvider

_Appears in:_
- [KeptnMetricsProvider](#keptnmetricsprovider)

| Field | Description | Default | Optional |
| --- | --- | --- | --- |
| `targetServer` _string_ |  || x |
| `secretKeyRef` _[SecretKeySelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#secretkeyselector-v1-core)_ |  || ✓ |


#### ProviderRef



ProviderRef represents the provider object

_Appears in:_
- [KeptnMetricSpec](#keptnmetricspec)

| Field | Description | Default | Optional |
| --- | --- | --- | --- |
| `name` _string_ | Name of the provider || x |


