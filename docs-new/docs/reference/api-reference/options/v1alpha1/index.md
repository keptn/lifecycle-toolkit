# v1alpha1
<!-- markdownlint-disable -->

## Packages
- [options.keptn.sh/v1alpha1](#optionskeptnshv1alpha1)


## options.keptn.sh/v1alpha1

Package v1alpha1 contains API Schema definitions for the options v1alpha1 API group

### Resource Types
- [KeptnConfig](#keptnconfig)
- [KeptnConfigList](#keptnconfiglist)



#### KeptnConfig



KeptnConfig is the Schema for the keptnconfigs API

_Appears in:_
- [KeptnConfigList](#keptnconfiglist)

| Field | Description | Default | Optional |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `options.keptn.sh/v1alpha1` | | |
| `kind` _string_ | `KeptnConfig` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. || ✓ |
| `spec` _[KeptnConfigSpec](#keptnconfigspec)_ |  || ✓ |
| `status` _string_ | unused field || ✓ |


#### KeptnConfigList



KeptnConfigList contains a list of KeptnConfig



| Field | Description | Default | Optional |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `options.keptn.sh/v1alpha1` | | |
| `kind` _string_ | `KeptnConfigList` | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. || ✓ |
| `items` _[KeptnConfig](#keptnconfig) array_ |  || x |


#### KeptnConfigSpec



KeptnConfigSpec defines the desired state of KeptnConfig

_Appears in:_
- [KeptnConfig](#keptnconfig)

| Field | Description | Default | Optional |
| --- | --- | --- | --- |
| `OTelCollectorUrl` _string_ | OTelCollectorUrl can be used to set the Open Telemetry collector that the lifecycle operator should use || ✓ |
| `keptnAppCreationRequestTimeoutSeconds` _integer_ | KeptnAppCreationRequestTimeoutSeconds is used to set the interval in which automatic app discovery searches for workload to put into the same auto-generated KeptnApp |30| ✓ |
| `cloudEventsEndpoint` _string_ | CloudEventsEndpoint can be used to set the endpoint where Cloud Events should be posted by the lifecycle operator || ✓ |


