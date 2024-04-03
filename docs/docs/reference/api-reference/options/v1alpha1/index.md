# v1alpha1

Reference information for options.keptn.sh/v1alpha1

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

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `apiVersion` _string_ | `options.keptn.sh/v1alpha1` | | | |
| `kind` _string_ | `KeptnConfig` | | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation about [`metadata`](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/#attaching-metadata-to-objects). || ✓ |  |
| `spec` _[KeptnConfigSpec](#keptnconfigspec)_ |  || ✓ |  |
| `status` _string_ | unused field || ✓ |  |


#### KeptnConfigList



KeptnConfigList contains a list of KeptnConfig





| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `apiVersion` _string_ | `options.keptn.sh/v1alpha1` | | | |
| `kind` _string_ | `KeptnConfigList` | | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#listmeta-v1-meta)_ |  || ✓ |  |
| `items` _[KeptnConfig](#keptnconfig) array_ |  || x |  |


#### KeptnConfigSpec



KeptnConfigSpec defines the desired state of KeptnConfig



_Appears in:_
- [KeptnConfig](#keptnconfig)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `OTelCollectorUrl` _string_ | OTelCollectorUrl can be used to set the Open Telemetry collector that the lifecycle operator should use || ✓ |  |
| `keptnAppCreationRequestTimeoutSeconds` _integer_ | KeptnAppCreationRequestTimeoutSeconds is used to set the interval in which automatic app discovery<br />searches for workload to put into the same auto-generated KeptnApp |30| ✓ |  |
| `cloudEventsEndpoint` _string_ | CloudEventsEndpoint can be used to set the endpoint where Cloud Events should be posted by the lifecycle operator || ✓ |  |
| `blockDeployment` _boolean_ | BlockDeployment is used to block the deployment of the application until the pre-deployment<br />tasks and evaluations succeed |true| ✓ |  |
| `observabilityTimeout` _[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#duration-v1-meta)_ | ObservabilityTimeout specifies the maximum time to observe the deployment phase of KeptnWorkload.<br />If the workload does not deploy successfully within this time frame, it will be<br />considered as failed. |5m| ✓ | Pattern: `^0|([0-9]+(\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$` <br />Type: string <br /> |


