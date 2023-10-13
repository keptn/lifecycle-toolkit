---
title: v1alpha4
description: Reference information for lifecycle.keptn.sh/v1alpha4
---
<!-- markdownlint-disable -->

## Packages
- [lifecycle.keptn.sh/v1alpha4](#lifecyclekeptnshv1alpha4)


## lifecycle.keptn.sh/v1alpha4

Package v1alpha4 contains API Schema definitions for the lifecycle v1alpha4 API group

### Resource Types
- [KeptnWorkloadVersion](#keptnworkloadversion)
- [KeptnWorkloadVersionList](#keptnworkloadversionlist)



#### KeptnWorkloadVersion



KeptnWorkloadVersion is the Schema for the keptnworkloadversions API

_Appears in:_
- [KeptnWorkloadVersionList](#keptnworkloadversionlist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `lifecycle.keptn.sh/v1alpha4`
| `kind` _string_ | `KeptnWorkloadVersion`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[KeptnWorkloadVersionSpec](#keptnworkloadversionspec)_ | Spec describes the desired state of the KeptnWorkloadVersion. |
| `status` _[KeptnWorkloadVersionStatus](#keptnworkloadversionstatus)_ | Status describes the current state of the KeptnWorkloadVersion. |


#### KeptnWorkloadVersionList



KeptnWorkloadVersionList contains a list of KeptnWorkloadVersion



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `lifecycle.keptn.sh/v1alpha4`
| `kind` _string_ | `KeptnWorkloadVersionList`
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[KeptnWorkloadVersion](#keptnworkloadversion) array_ |  |


#### KeptnWorkloadVersionSpec



KeptnWorkloadVersionSpec defines the desired state of KeptnWorkloadVersion

_Appears in:_
- [KeptnWorkloadVersion](#keptnworkloadversion)

| Field | Description |
| --- | --- |
| `app` _string_ | AppName is the name of the KeptnApp containing the KeptnWorkload. |
| `version` _string_ | Version defines the version of the KeptnWorkload. |
| `preDeploymentTasks` _string array_ | PreDeploymentTasks is a list of all tasks to be performed during the pre-deployment phase of the KeptnWorkload. The items of this list refer to the names of KeptnTaskDefinitions located in the same namespace as the KeptnApp, or in the Keptn namespace. |
| `postDeploymentTasks` _string array_ | PostDeploymentTasks is a list of all tasks to be performed during the post-deployment phase of the KeptnWorkload. The items of this list refer to the names of KeptnTaskDefinitions located in the same namespace as the KeptnWorkload, or in the Keptn namespace. |
| `preDeploymentEvaluations` _string array_ | PreDeploymentEvaluations is a list of all evaluations to be performed during the pre-deployment phase of the KeptnWorkload. The items of this list refer to the names of KeptnEvaluationDefinitions located in the same namespace as the KeptnWorkload, or in the Keptn namespace. |
| `postDeploymentEvaluations` _string array_ | PostDeploymentEvaluations is a list of all evaluations to be performed during the post-deployment phase of the KeptnWorkload. The items of this list refer to the names of KeptnEvaluationDefinitions located in the same namespace as the KeptnWorkload, or in the Keptn namespace. |
| `resourceReference` _[ResourceReference](../v1alpha3/#resourcereference)_ | ResourceReference is a reference to the Kubernetes resource (Deployment, DaemonSet, StatefulSet or ReplicaSet) the KeptnWorkload is representing. |
| `workloadName` _string_ | WorkloadName is the name of the KeptnWorkload. |
| `previousVersion` _string_ | PreviousVersion is the version of the KeptnWorkload that has been deployed prior to this version. |
| `traceId` _object (keys:string, values:string)_ | TraceId contains the OpenTelemetry trace ID. |


#### KeptnWorkloadVersionStatus



KeptnWorkloadVersionStatus defines the observed state of KeptnWorkloadVersion

_Appears in:_
- [KeptnWorkloadVersion](#keptnworkloadversion)

| Field | Description |
| --- | --- |
| `preDeploymentStatus` _KeptnState_ | PreDeploymentStatus indicates the current status of the KeptnWorkloadVersion's PreDeployment phase. |
| `deploymentStatus` _KeptnState_ | DeploymentStatus indicates the current status of the KeptnWorkloadVersion's Deployment phase. |
| `preDeploymentEvaluationStatus` _KeptnState_ | PreDeploymentEvaluationStatus indicates the current status of the KeptnWorkloadVersion's PreDeploymentEvaluation phase. |
| `postDeploymentEvaluationStatus` _KeptnState_ | PostDeploymentEvaluationStatus indicates the current status of the KeptnWorkloadVersion's PostDeploymentEvaluation phase. |
| `postDeploymentStatus` _KeptnState_ | PostDeploymentStatus indicates the current status of the KeptnWorkloadVersion's PostDeployment phase. |
| `preDeploymentTaskStatus` _[ItemStatus](../v1alpha3/#itemstatus) array_ | PreDeploymentTaskStatus indicates the current state of each preDeploymentTask of the KeptnWorkloadVersion. |
| `postDeploymentTaskStatus` _[ItemStatus](../v1alpha3/#itemstatus) array_ | PostDeploymentTaskStatus indicates the current state of each postDeploymentTask of the KeptnWorkloadVersion. |
| `preDeploymentEvaluationTaskStatus` _[ItemStatus](../v1alpha3/#itemstatus) array_ | PreDeploymentEvaluationTaskStatus indicates the current state of each preDeploymentEvaluation of the KeptnWorkloadVersion. |
| `postDeploymentEvaluationTaskStatus` _[ItemStatus](../v1alpha3/#itemstatus) array_ | PostDeploymentEvaluationTaskStatus indicates the current state of each postDeploymentEvaluation of the KeptnWorkloadVersion. |
| `startTime` _[Time](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#time-v1-meta)_ | StartTime represents the time at which the deployment of the KeptnWorkloadVersion started. |
| `endTime` _[Time](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#time-v1-meta)_ | EndTime represents the time at which the deployment of the KeptnWorkloadVersion finished. |
| `currentPhase` _string_ | CurrentPhase indicates the current phase of the KeptnWorkloadVersion. This can be: - PreDeploymentTasks - PreDeploymentEvaluations - Deployment - PostDeploymentTasks - PostDeploymentEvaluations |
| `phaseTraceIDs` _object (keys:string, values:object)_ | PhaseTraceIDs contains the trace IDs of the OpenTelemetry spans of each phase of the KeptnWorkloadVersion |
| `status` _KeptnState_ | Status represents the overall status of the KeptnWorkloadVersion. |


