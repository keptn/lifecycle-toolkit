# v1alpha1

Reference information for metrics.keptn.sh/v1alpha1

<!-- markdownlint-disable -->

## Packages
- [metrics.keptn.sh/v1alpha1](#metricskeptnshv1alpha1)


## metrics.keptn.sh/v1alpha1

Package v1alpha1 contains API Schema definitions for the metrics v1alpha1 API group

### Resource Types
- [KeptnMetric](#keptnmetric)
- [KeptnMetricList](#keptnmetriclist)



#### KeptnMetric



KeptnMetric is the Schema for the keptnmetrics API



_Appears in:_
- [KeptnMetricList](#keptnmetriclist)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `apiVersion` _string_ | `metrics.keptn.sh/v1alpha1` | | | |
| `kind` _string_ | `KeptnMetric` | | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation about [`metadata`](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/#attaching-metadata-to-objects). || ✓ |  |
| `spec` _[KeptnMetricSpec](#keptnmetricspec)_ |  || ✓ |  |
| `status` _[KeptnMetricStatus](#keptnmetricstatus)_ |  || ✓ |  |


#### KeptnMetricList



KeptnMetricList contains a list of KeptnMetric





| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `apiVersion` _string_ | `metrics.keptn.sh/v1alpha1` | | | |
| `kind` _string_ | `KeptnMetricList` | | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#listmeta-v1-meta)_ |  || ✓ |  |
| `items` _[KeptnMetric](#keptnmetric) array_ |  || x |  |


#### KeptnMetricSpec



KeptnMetricSpec defines the desired state of KeptnMetric



_Appears in:_
- [KeptnMetric](#keptnmetric)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `provider` _[ProviderRef](#providerref)_ | Provider represents the provider object || x |  |
| `query` _string_ | Query represents the query to be run || x |  |
| `fetchIntervalSeconds` _integer_ | FetchIntervalSeconds represents the update frequency in seconds that is used to update the metric || x |  |


#### KeptnMetricStatus



KeptnMetricStatus defines the observed state of KeptnMetric



_Appears in:_
- [KeptnMetric](#keptnmetric)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `value` _string_ | Value represents the resulting value || x |  |
| `rawValue` _integer array_ | RawValue represents the resulting value in raw format || x |  |
| `lastUpdated` _[Time](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#time-v1-meta)_ | LastUpdated represents the time when the status data was last updated || x |  |


#### ProviderRef



ProviderRef represents the provider object



_Appears in:_
- [KeptnMetricSpec](#keptnmetricspec)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `name` _string_ | Name of the provider || x |  |


