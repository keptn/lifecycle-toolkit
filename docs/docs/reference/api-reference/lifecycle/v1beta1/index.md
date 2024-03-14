# v1beta1

Reference information for lifecycle.keptn.sh/v1beta1

<!-- markdownlint-disable -->

## Packages
- [lifecycle.keptn.sh/v1beta1](#lifecyclekeptnshv1beta1)


## lifecycle.keptn.sh/v1beta1

Package v1beta1 contains API Schema definitions for the lifecycle v1beta1 API group

### Resource Types
- [KeptnApp](#keptnapp)
- [KeptnAppContext](#keptnappcontext)
- [KeptnAppContextList](#keptnappcontextlist)
- [KeptnAppCreationRequest](#keptnappcreationrequest)
- [KeptnAppCreationRequestList](#keptnappcreationrequestlist)
- [KeptnAppList](#keptnapplist)
- [KeptnAppVersion](#keptnappversion)
- [KeptnAppVersionList](#keptnappversionlist)
- [KeptnEvaluation](#keptnevaluation)
- [KeptnEvaluationDefinition](#keptnevaluationdefinition)
- [KeptnEvaluationDefinitionList](#keptnevaluationdefinitionlist)
- [KeptnEvaluationList](#keptnevaluationlist)
- [KeptnTask](#keptntask)
- [KeptnTaskDefinition](#keptntaskdefinition)
- [KeptnTaskDefinitionList](#keptntaskdefinitionlist)
- [KeptnTaskList](#keptntasklist)
- [KeptnWorkload](#keptnworkload)
- [KeptnWorkloadList](#keptnworkloadlist)
- [KeptnWorkloadVersion](#keptnworkloadversion)
- [KeptnWorkloadVersionList](#keptnworkloadversionlist)





#### AutomountServiceAccountTokenSpec







_Appears in:_
- [KeptnTaskDefinitionSpec](#keptntaskdefinitionspec)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `type` _boolean_ |  || x |  |


#### CheckType

_Underlying type:_ _string_





_Appears in:_
- [KeptnEvaluationSpec](#keptnevaluationspec)
- [KeptnTaskSpec](#keptntaskspec)



#### ConfigMapReference







_Appears in:_
- [RuntimeSpec](#runtimespec)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `name` _string_ | Name is the name of the referenced ConfigMap. || ✓ |  |


#### ContainerSpec







_Appears in:_
- [KeptnTaskDefinitionSpec](#keptntaskdefinitionspec)



#### DeploymentTaskSpec







_Appears in:_
- [KeptnAppContextSpec](#keptnappcontextspec)
- [KeptnAppVersionSpec](#keptnappversionspec)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `preDeploymentTasks` _string array_ | PreDeploymentTasks is a list of all tasks to be performed during the pre-deployment phase of the KeptnApp.<br />The items of this list refer to the names of KeptnTaskDefinitions<br />located in the same namespace as the KeptnApp, or in the Keptn namespace. || x |  |
| `postDeploymentTasks` _string array_ | PostDeploymentTasks is a list of all tasks to be performed during the post-deployment phase of the KeptnApp.<br />The items of this list refer to the names of KeptnTaskDefinitions<br />located in the same namespace as the KeptnApp, or in the Keptn namespace. || x |  |
| `preDeploymentEvaluations` _string array_ | PreDeploymentEvaluations is a list of all evaluations to be performed<br />during the pre-deployment phase of the KeptnApp.<br />The items of this list refer to the names of KeptnEvaluationDefinitions<br />located in the same namespace as the KeptnApp, or in the Keptn namespace. || x |  |
| `postDeploymentEvaluations` _string array_ | PostDeploymentEvaluations is a list of all evaluations to be performed<br />during the post-deployment phase of the KeptnApp.<br />The items of this list refer to the names of KeptnEvaluationDefinitions<br />located in the same namespace as the KeptnApp, or in the Keptn namespace. || x |  |
| `promotionTasks` _string array_ | PromotionTasks is a list of all tasks to be performed during the promotion phase of the KeptnApp.<br />The items of this list refer to the names of KeptnTaskDefinitions<br />located in the same namespace as the KeptnApp, or in the Keptn namespace. || x |  |


#### EvaluationStatusItem







_Appears in:_
- [KeptnEvaluationStatus](#keptnevaluationstatus)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `value` _string_ | Value represents the value of the KeptnMetric being evaluated. || x |  |
| `status` _string_ | Status indicates the status of the objective being evaluated. || x |  |
| `message` _string_ | Message contains additional information about the evaluation of an objective.<br />This can include explanations about why an evaluation has failed (e.g. due to a missed objective),<br />or if there was any error during the evaluation of the objective. || ✓ |  |


#### FailureConditions



FailureConditions represent the failure conditions (number of retries and retry interval)
for the evaluation to be considered as failed



_Appears in:_
- [KeptnEvaluationDefinitionSpec](#keptnevaluationdefinitionspec)
- [KeptnEvaluationSpec](#keptnevaluationspec)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `retries` _integer_ | Retries indicates how many times the KeptnEvaluation can be attempted in the case of an error or<br />missed evaluation objective, before considering the KeptnEvaluation to be failed. |10| ✓ |  |
| `retryInterval` _[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#duration-v1-meta)_ | RetryInterval specifies the interval at which the KeptnEvaluation is retried in the case of an error<br />or a missed objective. |5s| ✓ | Pattern: `^0|([0-9]+(\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$` <br />Type: string <br /> |


#### FunctionReference







_Appears in:_
- [RuntimeSpec](#runtimespec)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `name` _string_ | Name is the name of the referenced KeptnTaskDefinition. || ✓ |  |


#### FunctionStatus







_Appears in:_
- [KeptnTaskDefinitionStatus](#keptntaskdefinitionstatus)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `configMap` _string_ | ConfigMap indicates the ConfigMap in which the function code is stored. || ✓ |  |


#### HttpReference







_Appears in:_
- [RuntimeSpec](#runtimespec)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `url` _string_ | Url is the URL containing the code of the function. || ✓ |  |


#### Inline







_Appears in:_
- [RuntimeSpec](#runtimespec)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `code` _string_ | Code contains the code of the function. || ✓ |  |


#### ItemStatus







_Appears in:_
- [KeptnAppVersionStatus](#keptnappversionstatus)
- [KeptnWorkloadVersionStatus](#keptnworkloadversionstatus)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `definitionName` _string_ | DefinitionName is the name of the EvaluationDefinition/TaskDefinition || ✓ |  |
| `status` _string_ |  |Pending| ✓ |  |
| `name` _string_ | Name is the name of the Evaluation/Task || ✓ |  |
| `startTime` _[Time](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#time-v1-meta)_ | StartTime represents the time at which the Item (Evaluation/Task) started. || ✓ |  |
| `endTime` _[Time](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#time-v1-meta)_ | EndTime represents the time at which the Item (Evaluation/Task) started. || ✓ |  |


#### KeptnApp



KeptnApp is the Schema for the keptnapps API



_Appears in:_
- [KeptnAppList](#keptnapplist)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `apiVersion` _string_ | `lifecycle.keptn.sh/v1beta1` | | | |
| `kind` _string_ | `KeptnApp` | | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation about [`metadata`](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/#attaching-metadata-to-objects). || ✓ |  |
| `spec` _[KeptnAppSpec](#keptnappspec)_ | Spec describes the desired state of the KeptnApp. || ✓ |  |
| `status` _[KeptnAppStatus](#keptnappstatus)_ | Status describes the current state of the KeptnApp. || ✓ |  |


#### KeptnAppContext



KeptnAppContext is the Schema for the keptnappcontexts API



_Appears in:_
- [KeptnAppContextList](#keptnappcontextlist)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `apiVersion` _string_ | `lifecycle.keptn.sh/v1beta1` | | | |
| `kind` _string_ | `KeptnAppContext` | | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation about [`metadata`](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/#attaching-metadata-to-objects). || x |  |
| `spec` _[KeptnAppContextSpec](#keptnappcontextspec)_ |  || x |  |
| `status` _[KeptnAppContextStatus](#keptnappcontextstatus)_ |  || x |  |


#### KeptnAppContextList



KeptnAppContextList contains a list of KeptnAppContext





| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `apiVersion` _string_ | `lifecycle.keptn.sh/v1beta1` | | | |
| `kind` _string_ | `KeptnAppContextList` | | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#listmeta-v1-meta)_ |  || x |  |
| `items` _[KeptnAppContext](#keptnappcontext) array_ |  || x |  |


#### KeptnAppContextSpec



KeptnAppContextSpec defines the desired state of KeptnAppContext



_Appears in:_
- [KeptnAppContext](#keptnappcontext)
- [KeptnAppVersionSpec](#keptnappversionspec)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `preDeploymentTasks` _string array_ | PreDeploymentTasks is a list of all tasks to be performed during the pre-deployment phase of the KeptnApp.<br />The items of this list refer to the names of KeptnTaskDefinitions<br />located in the same namespace as the KeptnApp, or in the Keptn namespace. || x |  |
| `postDeploymentTasks` _string array_ | PostDeploymentTasks is a list of all tasks to be performed during the post-deployment phase of the KeptnApp.<br />The items of this list refer to the names of KeptnTaskDefinitions<br />located in the same namespace as the KeptnApp, or in the Keptn namespace. || x |  |
| `preDeploymentEvaluations` _string array_ | PreDeploymentEvaluations is a list of all evaluations to be performed<br />during the pre-deployment phase of the KeptnApp.<br />The items of this list refer to the names of KeptnEvaluationDefinitions<br />located in the same namespace as the KeptnApp, or in the Keptn namespace. || x |  |
| `postDeploymentEvaluations` _string array_ | PostDeploymentEvaluations is a list of all evaluations to be performed<br />during the post-deployment phase of the KeptnApp.<br />The items of this list refer to the names of KeptnEvaluationDefinitions<br />located in the same namespace as the KeptnApp, or in the Keptn namespace. || x |  |
| `promotionTasks` _string array_ | PromotionTasks is a list of all tasks to be performed during the promotion phase of the KeptnApp.<br />The items of this list refer to the names of KeptnTaskDefinitions<br />located in the same namespace as the KeptnApp, or in the Keptn namespace. || x |  |
| `metadata` _object (keys:string, values:string)_ | Metadata contains additional key-value pairs for contextual information. || ✓ |  |
| `spanLinks` _string array_ | SpanLinks are links to OpenTelemetry span IDs for tracking. These links establish relationships between spans across different services, enabling distributed tracing.<br />For more information on OpenTelemetry span links, refer to the documentation: https://opentelemetry.io/docs/concepts/signals/traces/#span-links || ✓ |  |


#### KeptnAppContextStatus



KeptnAppContextStatus defines the observed state of KeptnAppContext



_Appears in:_
- [KeptnAppContext](#keptnappcontext)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `status` _string_ | unused field || ✓ |  |


#### KeptnAppCreationRequest



KeptnAppCreationRequest is the Schema for the keptnappcreationrequests API



_Appears in:_
- [KeptnAppCreationRequestList](#keptnappcreationrequestlist)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `apiVersion` _string_ | `lifecycle.keptn.sh/v1beta1` | | | |
| `kind` _string_ | `KeptnAppCreationRequest` | | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation about [`metadata`](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/#attaching-metadata-to-objects). || ✓ |  |
| `spec` _[KeptnAppCreationRequestSpec](#keptnappcreationrequestspec)_ | Spec describes the desired state of the KeptnAppCreationRequest. || ✓ |  |
| `status` _string_ | Status describes the current state of the KeptnAppCreationRequest. || ✓ |  |


#### KeptnAppCreationRequestList



KeptnAppCreationRequestList contains a list of KeptnAppCreationRequest





| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `apiVersion` _string_ | `lifecycle.keptn.sh/v1beta1` | | | |
| `kind` _string_ | `KeptnAppCreationRequestList` | | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#listmeta-v1-meta)_ |  || ✓ |  |
| `items` _[KeptnAppCreationRequest](#keptnappcreationrequest) array_ |  || x |  |


#### KeptnAppCreationRequestSpec



KeptnAppCreationRequestSpec defines the desired state of KeptnAppCreationRequest



_Appears in:_
- [KeptnAppCreationRequest](#keptnappcreationrequest)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `appName` _string_ | AppName is the name of the KeptnApp the KeptnAppCreationRequest should create if no user-defined object with that name is found. || x |  |


#### KeptnAppList



KeptnAppList contains a list of KeptnApp





| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `apiVersion` _string_ | `lifecycle.keptn.sh/v1beta1` | | | |
| `kind` _string_ | `KeptnAppList` | | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#listmeta-v1-meta)_ |  || ✓ |  |
| `items` _[KeptnApp](#keptnapp) array_ |  || x |  |


#### KeptnAppSpec



KeptnAppSpec defines the desired state of KeptnApp



_Appears in:_
- [KeptnApp](#keptnapp)
- [KeptnAppVersionSpec](#keptnappversionspec)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `version` _string_ | Version defines the version of the application. For automatically created KeptnApps,<br />the version is a function of all KeptnWorkloads that are part of the KeptnApp. || x |  |
| `revision` _integer_ | Revision can be modified to trigger another deployment of a KeptnApp of the same version.<br />This can be used for restarting a KeptnApp which failed to deploy,<br />e.g. due to a failed preDeploymentEvaluation/preDeploymentTask. |1| ✓ |  |
| `workloads` _[KeptnWorkloadRef](#keptnworkloadref) array_ | Workloads is a list of all KeptnWorkloads that are part of the KeptnApp. || ✓ |  |


#### KeptnAppStatus



KeptnAppStatus defines the observed state of KeptnApp



_Appears in:_
- [KeptnApp](#keptnapp)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `currentVersion` _string_ | CurrentVersion indicates the version that is currently deployed or being reconciled. || ✓ |  |


#### KeptnAppVersion



KeptnAppVersion is the Schema for the keptnappversions API



_Appears in:_
- [KeptnAppVersionList](#keptnappversionlist)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `apiVersion` _string_ | `lifecycle.keptn.sh/v1beta1` | | | |
| `kind` _string_ | `KeptnAppVersion` | | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation about [`metadata`](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/#attaching-metadata-to-objects). || ✓ |  |
| `spec` _[KeptnAppVersionSpec](#keptnappversionspec)_ | Spec describes the desired state of the KeptnAppVersion. || ✓ |  |
| `status` _[KeptnAppVersionStatus](#keptnappversionstatus)_ | Status describes the current state of the KeptnAppVersion. || ✓ |  |


#### KeptnAppVersionList



KeptnAppVersionList contains a list of KeptnAppVersion





| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `apiVersion` _string_ | `lifecycle.keptn.sh/v1beta1` | | | |
| `kind` _string_ | `KeptnAppVersionList` | | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#listmeta-v1-meta)_ |  || ✓ |  |
| `items` _[KeptnAppVersion](#keptnappversion) array_ |  || x |  |


#### KeptnAppVersionSpec



KeptnAppVersionSpec defines the desired state of KeptnAppVersion



_Appears in:_
- [KeptnAppVersion](#keptnappversion)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `preDeploymentTasks` _string array_ | PreDeploymentTasks is a list of all tasks to be performed during the pre-deployment phase of the KeptnApp.<br />The items of this list refer to the names of KeptnTaskDefinitions<br />located in the same namespace as the KeptnApp, or in the Keptn namespace. || x |  |
| `postDeploymentTasks` _string array_ | PostDeploymentTasks is a list of all tasks to be performed during the post-deployment phase of the KeptnApp.<br />The items of this list refer to the names of KeptnTaskDefinitions<br />located in the same namespace as the KeptnApp, or in the Keptn namespace. || x |  |
| `preDeploymentEvaluations` _string array_ | PreDeploymentEvaluations is a list of all evaluations to be performed<br />during the pre-deployment phase of the KeptnApp.<br />The items of this list refer to the names of KeptnEvaluationDefinitions<br />located in the same namespace as the KeptnApp, or in the Keptn namespace. || x |  |
| `postDeploymentEvaluations` _string array_ | PostDeploymentEvaluations is a list of all evaluations to be performed<br />during the post-deployment phase of the KeptnApp.<br />The items of this list refer to the names of KeptnEvaluationDefinitions<br />located in the same namespace as the KeptnApp, or in the Keptn namespace. || x |  |
| `promotionTasks` _string array_ | PromotionTasks is a list of all tasks to be performed during the promotion phase of the KeptnApp.<br />The items of this list refer to the names of KeptnTaskDefinitions<br />located in the same namespace as the KeptnApp, or in the Keptn namespace. || x |  |
| `metadata` _object (keys:string, values:string)_ | Metadata contains additional key-value pairs for contextual information. || ✓ |  |
| `spanLinks` _string array_ | SpanLinks are links to OpenTelemetry span IDs for tracking. These links establish relationships between spans across different services, enabling distributed tracing.<br />For more information on OpenTelemetry span links, refer to the documentation: https://opentelemetry.io/docs/concepts/signals/traces/#span-links || ✓ |  |
| `version` _string_ | Version defines the version of the application. For automatically created KeptnApps,<br />the version is a function of all KeptnWorkloads that are part of the KeptnApp. || x |  |
| `revision` _integer_ | Revision can be modified to trigger another deployment of a KeptnApp of the same version.<br />This can be used for restarting a KeptnApp which failed to deploy,<br />e.g. due to a failed preDeploymentEvaluation/preDeploymentTask. |1| ✓ |  |
| `workloads` _[KeptnWorkloadRef](#keptnworkloadref) array_ | Workloads is a list of all KeptnWorkloads that are part of the KeptnApp. || ✓ |  |
| `appName` _string_ | AppName is the name of the KeptnApp. || x |  |
| `previousVersion` _string_ | PreviousVersion is the version of the KeptnApp that has been deployed prior to this version. || ✓ |  |
| `traceId` _object (keys:string, values:string)_ | TraceId contains the OpenTelemetry trace ID. || ✓ |  |


#### KeptnAppVersionStatus



KeptnAppVersionStatus defines the observed state of KeptnAppVersion



_Appears in:_
- [KeptnAppVersion](#keptnappversion)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `preDeploymentStatus` _string_ | PreDeploymentStatus indicates the current status of the KeptnAppVersion's PreDeployment phase. |Pending| ✓ |  |
| `postDeploymentStatus` _string_ | PostDeploymentStatus indicates the current status of the KeptnAppVersion's PostDeployment phase. |Pending| ✓ |  |
| `promotionStatus` _string_ | PromotionStatus indicates the current status of the KeptnAppVersion's Promotion phase. |Pending| ✓ |  |
| `preDeploymentEvaluationStatus` _string_ | PreDeploymentEvaluationStatus indicates the current status of the KeptnAppVersion's PreDeploymentEvaluation phase. |Pending| ✓ |  |
| `postDeploymentEvaluationStatus` _string_ | PostDeploymentEvaluationStatus indicates the current status of the KeptnAppVersion's PostDeploymentEvaluation phase. |Pending| ✓ |  |
| `workloadOverallStatus` _string_ | WorkloadOverallStatus indicates the current status of the KeptnAppVersion's Workload deployment phase. |Pending| ✓ |  |
| `workloadStatus` _[WorkloadStatus](#workloadstatus) array_ | WorkloadStatus contains the current status of each KeptnWorkload that is part of the KeptnAppVersion. || ✓ |  |
| `currentPhase` _string_ | CurrentPhase indicates the current phase of the KeptnAppVersion. || ✓ |  |
| `preDeploymentTaskStatus` _[ItemStatus](#itemstatus) array_ | PreDeploymentTaskStatus indicates the current state of each preDeploymentTask of the KeptnAppVersion. || ✓ |  |
| `postDeploymentTaskStatus` _[ItemStatus](#itemstatus) array_ | PostDeploymentTaskStatus indicates the current state of each postDeploymentTask of the KeptnAppVersion. || ✓ |  |
| `promotionTaskStatus` _[ItemStatus](#itemstatus) array_ | PromotionTaskStatus indicates the current state of each promotionTask of the KeptnAppVersion. || ✓ |  |
| `preDeploymentEvaluationTaskStatus` _[ItemStatus](#itemstatus) array_ | PreDeploymentEvaluationTaskStatus indicates the current state of each preDeploymentEvaluation of the KeptnAppVersion. || ✓ |  |
| `postDeploymentEvaluationTaskStatus` _[ItemStatus](#itemstatus) array_ | PostDeploymentEvaluationTaskStatus indicates the current state of each postDeploymentEvaluation of the KeptnAppVersion. || ✓ |  |
| `phaseTraceIDs` _[MapCarrier](https://pkg.go.dev/go.opentelemetry.io/otel/propagation#MapCarrier)_ | PhaseTraceIDs contains the trace IDs of the OpenTelemetry spans of each phase of the KeptnAppVersion. || ✓ |  |
| `status` _string_ | Status represents the overall status of the KeptnAppVersion. |Pending| ✓ |  |
| `startTime` _[Time](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#time-v1-meta)_ | StartTime represents the time at which the deployment of the KeptnAppVersion started. || ✓ |  |
| `endTime` _[Time](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#time-v1-meta)_ | EndTime represents the time at which the deployment of the KeptnAppVersion finished. || ✓ |  |


#### KeptnEvaluation



KeptnEvaluation is the Schema for the keptnevaluations API



_Appears in:_
- [KeptnEvaluationList](#keptnevaluationlist)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `apiVersion` _string_ | `lifecycle.keptn.sh/v1beta1` | | | |
| `kind` _string_ | `KeptnEvaluation` | | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation about [`metadata`](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/#attaching-metadata-to-objects). || ✓ |  |
| `spec` _[KeptnEvaluationSpec](#keptnevaluationspec)_ | Spec describes the desired state of the KeptnEvaluation. || ✓ |  |
| `status` _[KeptnEvaluationStatus](#keptnevaluationstatus)_ | Status describes the current state of the KeptnEvaluation. || ✓ |  |


#### KeptnEvaluationDefinition



KeptnEvaluationDefinition is the Schema for the keptnevaluationdefinitions API



_Appears in:_
- [KeptnEvaluationDefinitionList](#keptnevaluationdefinitionlist)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `apiVersion` _string_ | `lifecycle.keptn.sh/v1beta1` | | | |
| `kind` _string_ | `KeptnEvaluationDefinition` | | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation about [`metadata`](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/#attaching-metadata-to-objects). || ✓ |  |
| `spec` _[KeptnEvaluationDefinitionSpec](#keptnevaluationdefinitionspec)_ | Spec describes the desired state of the KeptnEvaluationDefinition. || ✓ |  |
| `status` _string_ | unused field || ✓ |  |


#### KeptnEvaluationDefinitionList



KeptnEvaluationDefinitionList contains a list of KeptnEvaluationDefinition





| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `apiVersion` _string_ | `lifecycle.keptn.sh/v1beta1` | | | |
| `kind` _string_ | `KeptnEvaluationDefinitionList` | | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#listmeta-v1-meta)_ |  || ✓ |  |
| `items` _[KeptnEvaluationDefinition](#keptnevaluationdefinition) array_ |  || x |  |


#### KeptnEvaluationDefinitionSpec



KeptnEvaluationDefinitionSpec defines the desired state of KeptnEvaluationDefinition



_Appears in:_
- [KeptnEvaluationDefinition](#keptnevaluationdefinition)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `objectives` _[Objective](#objective) array_ | Objectives is a list of objectives that have to be met for a KeptnEvaluation referencing this<br />KeptnEvaluationDefinition to be successful. || x |  |
| `retries` _integer_ | Retries indicates how many times the KeptnEvaluation can be attempted in the case of an error or<br />missed evaluation objective, before considering the KeptnEvaluation to be failed. |10| ✓ |  |
| `retryInterval` _[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#duration-v1-meta)_ | RetryInterval specifies the interval at which the KeptnEvaluation is retried in the case of an error<br />or a missed objective. |5s| ✓ | Pattern: `^0|([0-9]+(\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$` <br />Type: string <br /> |


#### KeptnEvaluationList



KeptnEvaluationList contains a list of KeptnEvaluation





| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `apiVersion` _string_ | `lifecycle.keptn.sh/v1beta1` | | | |
| `kind` _string_ | `KeptnEvaluationList` | | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#listmeta-v1-meta)_ |  || ✓ |  |
| `items` _[KeptnEvaluation](#keptnevaluation) array_ |  || x |  |


#### KeptnEvaluationSpec



KeptnEvaluationSpec defines the desired state of KeptnEvaluation



_Appears in:_
- [KeptnEvaluation](#keptnevaluation)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `workload` _string_ | Workload defines the KeptnWorkload for which the KeptnEvaluation is done. || ✓ |  |
| `workloadVersion` _string_ | WorkloadVersion defines the version of the KeptnWorkload for which the KeptnEvaluation is done. || x |  |
| `appName` _string_ | AppName defines the KeptnApp for which the KeptnEvaluation is done. || ✓ |  |
| `appVersion` _string_ | AppVersion defines the version of the KeptnApp for which the KeptnEvaluation is done. || ✓ |  |
| `evaluationDefinition` _string_ | EvaluationDefinition refers to the name of the KeptnEvaluationDefinition<br />which includes the objectives for the KeptnEvaluation.<br />The KeptnEvaluationDefinition can be<br />located in the same namespace as the KeptnEvaluation, or in the Keptn namespace. || x |  |
| `checkType` _string_ | Type indicates whether the KeptnEvaluation is part of the pre- or postDeployment phase. || ✓ |  |
| `retries` _integer_ | Retries indicates how many times the KeptnEvaluation can be attempted in the case of an error or<br />missed evaluation objective, before considering the KeptnEvaluation to be failed. |10| ✓ |  |
| `retryInterval` _[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#duration-v1-meta)_ | RetryInterval specifies the interval at which the KeptnEvaluation is retried in the case of an error<br />or a missed objective. |5s| ✓ | Pattern: `^0|([0-9]+(\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$` <br />Type: string <br /> |


#### KeptnEvaluationStatus



KeptnEvaluationStatus defines the observed state of KeptnEvaluation



_Appears in:_
- [KeptnEvaluation](#keptnevaluation)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `retryCount` _integer_ | RetryCount indicates how many times the KeptnEvaluation has been attempted already. |0| x |  |
| `evaluationStatus` _object (keys:string, values:[EvaluationStatusItem](#evaluationstatusitem))_ | EvaluationStatus describes the status of each objective of the KeptnEvaluationDefinition<br />referenced by the KeptnEvaluation. || x |  |
| `overallStatus` _string_ | OverallStatus describes the overall status of the KeptnEvaluation. The Overall status is derived<br />from the status of the individual objectives of the KeptnEvaluationDefinition<br />referenced by the KeptnEvaluation. |Pending| x |  |
| `startTime` _[Time](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#time-v1-meta)_ | StartTime represents the time at which the KeptnEvaluation started. || ✓ |  |
| `endTime` _[Time](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#time-v1-meta)_ | EndTime represents the time at which the KeptnEvaluation finished. || ✓ |  |




#### KeptnMetricReference







_Appears in:_
- [Objective](#objective)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `name` _string_ | Name is the name of the referenced KeptnMetric. || x |  |
| `namespace` _string_ | Namespace is the namespace where the referenced KeptnMetric is located. || ✓ |  |




#### KeptnState

_Underlying type:_ _string_

KeptnState  is a string containing current Phase state  (Progressing/Succeeded/Failed/Unknown/Pending/Deprecated/Warning)



_Appears in:_
- [EvaluationStatusItem](#evaluationstatusitem)
- [ItemStatus](#itemstatus)
- [KeptnAppVersionStatus](#keptnappversionstatus)
- [KeptnEvaluationStatus](#keptnevaluationstatus)
- [KeptnTaskStatus](#keptntaskstatus)
- [KeptnWorkloadVersionStatus](#keptnworkloadversionstatus)
- [WorkloadStatus](#workloadstatus)



#### KeptnTask



KeptnTask is the Schema for the keptntasks API



_Appears in:_
- [KeptnTaskList](#keptntasklist)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `apiVersion` _string_ | `lifecycle.keptn.sh/v1beta1` | | | |
| `kind` _string_ | `KeptnTask` | | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation about [`metadata`](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/#attaching-metadata-to-objects). || ✓ |  |
| `spec` _[KeptnTaskSpec](#keptntaskspec)_ | Spec describes the desired state of the KeptnTask. || ✓ |  |
| `status` _[KeptnTaskStatus](#keptntaskstatus)_ | Status describes the current state of the KeptnTask. || ✓ |  |


#### KeptnTaskDefinition



KeptnTaskDefinition is the Schema for the keptntaskdefinitions API



_Appears in:_
- [KeptnTaskDefinitionList](#keptntaskdefinitionlist)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `apiVersion` _string_ | `lifecycle.keptn.sh/v1beta1` | | | |
| `kind` _string_ | `KeptnTaskDefinition` | | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation about [`metadata`](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/#attaching-metadata-to-objects). || ✓ |  |
| `spec` _[KeptnTaskDefinitionSpec](#keptntaskdefinitionspec)_ | Spec describes the desired state of the KeptnTaskDefinition. || ✓ |  |
| `status` _[KeptnTaskDefinitionStatus](#keptntaskdefinitionstatus)_ | Status describes the current state of the KeptnTaskDefinition. || ✓ |  |


#### KeptnTaskDefinitionList



KeptnTaskDefinitionList contains a list of KeptnTaskDefinition





| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `apiVersion` _string_ | `lifecycle.keptn.sh/v1beta1` | | | |
| `kind` _string_ | `KeptnTaskDefinitionList` | | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#listmeta-v1-meta)_ |  || ✓ |  |
| `items` _[KeptnTaskDefinition](#keptntaskdefinition) array_ |  || x |  |


#### KeptnTaskDefinitionSpec



KeptnTaskDefinitionSpec defines the desired state of KeptnTaskDefinition



_Appears in:_
- [KeptnTaskDefinition](#keptntaskdefinition)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `function` _[RuntimeSpec](#runtimespec)_ | Deprecated<br />Function contains the definition for the function that is to be executed in KeptnTasks. || ✓ |  |
| `python` _[RuntimeSpec](#runtimespec)_ | Python contains the definition for the python function that is to be executed in KeptnTasks. || ✓ |  |
| `deno` _[RuntimeSpec](#runtimespec)_ | Deno contains the definition for the Deno function that is to be executed in KeptnTasks. || ✓ |  |
| `container` _[ContainerSpec](#containerspec)_ | Container contains the definition for the container that is to be used in Job. || ✓ |  |
| `retries` _integer_ | Retries specifies how many times a job executing the KeptnTaskDefinition should be restarted in the case<br />of an unsuccessful attempt. |10| ✓ |  |
| `timeout` _[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#duration-v1-meta)_ | Timeout specifies the maximum time to wait for the task to be completed successfully.<br />If the task does not complete successfully within this time frame, it will be<br />considered to be failed. |5m| ✓ | Pattern: `^0|([0-9]+(\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$` <br />Type: string <br /> |
| `serviceAccount` _[ServiceAccountSpec](#serviceaccountspec)_ | ServiceAccount specifies the service account to be used in jobs to authenticate with the Kubernetes API and access cluster resources. || ✓ |  |
| `automountServiceAccountToken` _[AutomountServiceAccountTokenSpec](#automountserviceaccounttokenspec)_ | AutomountServiceAccountToken allows to enable K8s to assign cluster API credentials to a pod, if set to false<br />the pod will decline the service account || ✓ |  |
| `ttlSecondsAfterFinished` _integer_ | TTLSecondsAfterFinished controller makes a job eligible to be cleaned up after it is finished.<br />The timer starts when the status shows up to be Complete or Failed. |300| ✓ |  |
| `imagePullSecrets` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#localobjectreference-v1-core) array_ | ImagePullSecrets is an optional field to specify the names of secrets to use for pulling container images || ✓ |  |


#### KeptnTaskDefinitionStatus



KeptnTaskDefinitionStatus defines the observed state of KeptnTaskDefinition



_Appears in:_
- [KeptnTaskDefinition](#keptntaskdefinition)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `function` _[FunctionStatus](#functionstatus)_ | Function contains status information of the function definition for the task. || ✓ |  |


#### KeptnTaskList



KeptnTaskList contains a list of KeptnTask





| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `apiVersion` _string_ | `lifecycle.keptn.sh/v1beta1` | | | |
| `kind` _string_ | `KeptnTaskList` | | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#listmeta-v1-meta)_ |  || ✓ |  |
| `items` _[KeptnTask](#keptntask) array_ |  || x |  |


#### KeptnTaskSpec



KeptnTaskSpec defines the desired state of KeptnTask



_Appears in:_
- [KeptnTask](#keptntask)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `taskDefinition` _string_ | TaskDefinition refers to the name of the KeptnTaskDefinition<br />which includes the specification for the task to be performed.<br />The KeptnTaskDefinition can be<br />located in the same namespace as the KeptnTask, or in the Keptn namespace. || x |  |
| `context` _[TaskContext](#taskcontext)_ | Context contains contextual information about the task execution. || ✓ |  |
| `parameters` _[TaskParameters](#taskparameters)_ | Parameters contains parameters that will be passed to the job that executes the task. || ✓ |  |
| `secureParameters` _[SecureParameters](#secureparameters)_ | SecureParameters contains secure parameters that will be passed to the job that executes the task.<br />These will be stored and accessed as secrets in the cluster. || ✓ |  |
| `checkType` _string_ | Type indicates whether the KeptnTask is part of the pre- or postDeployment phase. || ✓ |  |
| `retries` _integer_ | Retries indicates how many times the KeptnTask can be attempted in the case of an error<br />before considering the KeptnTask to be failed. |10| ✓ |  |
| `timeout` _[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#duration-v1-meta)_ | Timeout specifies the maximum time to wait for the task to be completed successfully.<br />If the task does not complete successfully within this time frame, it will be<br />considered to be failed. |5m| ✓ | Pattern: `^0|([0-9]+(\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$` <br />Type: string <br /> |


#### KeptnTaskStatus



KeptnTaskStatus defines the observed state of KeptnTask



_Appears in:_
- [KeptnTask](#keptntask)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `jobName` _string_ | JobName is the name of the Job executing the Task. || ✓ |  |
| `status` _string_ | Status represents the overall state of the KeptnTask. |Pending| ✓ |  |
| `message` _string_ | Message contains information about unexpected errors encountered during the execution of the KeptnTask. || ✓ |  |
| `startTime` _[Time](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#time-v1-meta)_ | StartTime represents the time at which the KeptnTask started. || ✓ |  |
| `endTime` _[Time](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#time-v1-meta)_ | EndTime represents the time at which the KeptnTask finished. || ✓ |  |
| `reason` _string_ | Reason contains more information about the reason for the last transition of the Job executing the KeptnTask. || ✓ |  |


#### KeptnWorkload



KeptnWorkload is the Schema for the keptnworkloads API



_Appears in:_
- [KeptnWorkloadList](#keptnworkloadlist)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `apiVersion` _string_ | `lifecycle.keptn.sh/v1beta1` | | | |
| `kind` _string_ | `KeptnWorkload` | | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation about [`metadata`](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/#attaching-metadata-to-objects). || ✓ |  |
| `spec` _[KeptnWorkloadSpec](#keptnworkloadspec)_ | Spec describes the desired state of the KeptnWorkload. || ✓ |  |
| `status` _[KeptnWorkloadStatus](#keptnworkloadstatus)_ | Status describes the current state of the KeptnWorkload. || ✓ |  |


#### KeptnWorkloadList



KeptnWorkloadList contains a list of KeptnWorkload





| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `apiVersion` _string_ | `lifecycle.keptn.sh/v1beta1` | | | |
| `kind` _string_ | `KeptnWorkloadList` | | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#listmeta-v1-meta)_ |  || ✓ |  |
| `items` _[KeptnWorkload](#keptnworkload) array_ |  || x |  |


#### KeptnWorkloadRef



KeptnWorkloadRef refers to a KeptnWorkload that is part of a KeptnApp



_Appears in:_
- [KeptnAppSpec](#keptnappspec)
- [KeptnAppVersionSpec](#keptnappversionspec)
- [WorkloadStatus](#workloadstatus)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `name` _string_ | Name is the name of the KeptnWorkload. || x |  |
| `version` _string_ | Version is the version of the KeptnWorkload. || x |  |


#### KeptnWorkloadSpec



KeptnWorkloadSpec defines the desired state of KeptnWorkload



_Appears in:_
- [KeptnWorkload](#keptnworkload)
- [KeptnWorkloadVersionSpec](#keptnworkloadversionspec)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `app` _string_ | AppName is the name of the KeptnApp containing the KeptnWorkload. || x |  |
| `version` _string_ | Version defines the version of the KeptnWorkload. || x |  |
| `preDeploymentTasks` _string array_ | PreDeploymentTasks is a list of all tasks to be performed during the pre-deployment phase of the KeptnWorkload.<br />The items of this list refer to the names of KeptnTaskDefinitions<br />located in the same namespace as the KeptnApp, or in the Keptn namespace. || ✓ |  |
| `postDeploymentTasks` _string array_ | PostDeploymentTasks is a list of all tasks to be performed during the post-deployment phase of the KeptnWorkload.<br />The items of this list refer to the names of KeptnTaskDefinitions<br />located in the same namespace as the KeptnWorkload, or in the Keptn namespace. || ✓ |  |
| `preDeploymentEvaluations` _string array_ | PreDeploymentEvaluations is a list of all evaluations to be performed<br />during the pre-deployment phase of the KeptnWorkload.<br />The items of this list refer to the names of KeptnEvaluationDefinitions<br />located in the same namespace as the KeptnWorkload, or in the Keptn namespace. || ✓ |  |
| `postDeploymentEvaluations` _string array_ | PostDeploymentEvaluations is a list of all evaluations to be performed<br />during the post-deployment phase of the KeptnWorkload.<br />The items of this list refer to the names of KeptnEvaluationDefinitions<br />located in the same namespace as the KeptnWorkload, or in the Keptn namespace. || ✓ |  |
| `resourceReference` _[ResourceReference](#resourcereference)_ | ResourceReference is a reference to the Kubernetes resource<br />(Deployment, DaemonSet, StatefulSet or ReplicaSet) the KeptnWorkload is representing. || x |  |
| `metadata` _object (keys:string, values:string)_ | Metadata contains additional key-value pairs for contextual information. || ✓ |  |


#### KeptnWorkloadStatus



KeptnWorkloadStatus defines the observed state of KeptnWorkload



_Appears in:_
- [KeptnWorkload](#keptnworkload)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `currentVersion` _string_ | CurrentVersion indicates the version that is currently deployed or being reconciled. || ✓ |  |


#### KeptnWorkloadVersion



KeptnWorkloadVersion is the Schema for the keptnworkloadversions API



_Appears in:_
- [KeptnWorkloadVersionList](#keptnworkloadversionlist)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `apiVersion` _string_ | `lifecycle.keptn.sh/v1beta1` | | | |
| `kind` _string_ | `KeptnWorkloadVersion` | | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation about [`metadata`](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/#attaching-metadata-to-objects). || ✓ |  |
| `spec` _[KeptnWorkloadVersionSpec](#keptnworkloadversionspec)_ | Spec describes the desired state of the KeptnWorkloadVersion. || ✓ |  |
| `status` _[KeptnWorkloadVersionStatus](#keptnworkloadversionstatus)_ | Status describes the current state of the KeptnWorkloadVersion. || ✓ |  |


#### KeptnWorkloadVersionList



KeptnWorkloadVersionList contains a list of KeptnWorkloadVersion





| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `apiVersion` _string_ | `lifecycle.keptn.sh/v1beta1` | | | |
| `kind` _string_ | `KeptnWorkloadVersionList` | | | |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#listmeta-v1-meta)_ |  || ✓ |  |
| `items` _[KeptnWorkloadVersion](#keptnworkloadversion) array_ |  || x |  |


#### KeptnWorkloadVersionSpec



KeptnWorkloadVersionSpec defines the desired state of KeptnWorkloadVersion



_Appears in:_
- [KeptnWorkloadVersion](#keptnworkloadversion)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `app` _string_ | AppName is the name of the KeptnApp containing the KeptnWorkload. || x |  |
| `version` _string_ | Version defines the version of the KeptnWorkload. || x |  |
| `preDeploymentTasks` _string array_ | PreDeploymentTasks is a list of all tasks to be performed during the pre-deployment phase of the KeptnWorkload.<br />The items of this list refer to the names of KeptnTaskDefinitions<br />located in the same namespace as the KeptnApp, or in the Keptn namespace. || ✓ |  |
| `postDeploymentTasks` _string array_ | PostDeploymentTasks is a list of all tasks to be performed during the post-deployment phase of the KeptnWorkload.<br />The items of this list refer to the names of KeptnTaskDefinitions<br />located in the same namespace as the KeptnWorkload, or in the Keptn namespace. || ✓ |  |
| `preDeploymentEvaluations` _string array_ | PreDeploymentEvaluations is a list of all evaluations to be performed<br />during the pre-deployment phase of the KeptnWorkload.<br />The items of this list refer to the names of KeptnEvaluationDefinitions<br />located in the same namespace as the KeptnWorkload, or in the Keptn namespace. || ✓ |  |
| `postDeploymentEvaluations` _string array_ | PostDeploymentEvaluations is a list of all evaluations to be performed<br />during the post-deployment phase of the KeptnWorkload.<br />The items of this list refer to the names of KeptnEvaluationDefinitions<br />located in the same namespace as the KeptnWorkload, or in the Keptn namespace. || ✓ |  |
| `resourceReference` _[ResourceReference](#resourcereference)_ | ResourceReference is a reference to the Kubernetes resource<br />(Deployment, DaemonSet, StatefulSet or ReplicaSet) the KeptnWorkload is representing. || x |  |
| `metadata` _object (keys:string, values:string)_ | Metadata contains additional key-value pairs for contextual information. || ✓ |  |
| `workloadName` _string_ | WorkloadName is the name of the KeptnWorkload. || x |  |
| `previousVersion` _string_ | PreviousVersion is the version of the KeptnWorkload that has been deployed prior to this version. || ✓ |  |
| `traceId` _object (keys:string, values:string)_ | TraceId contains the OpenTelemetry trace ID. || ✓ |  |


#### KeptnWorkloadVersionStatus



KeptnWorkloadVersionStatus defines the observed state of KeptnWorkloadVersion



_Appears in:_
- [KeptnWorkloadVersion](#keptnworkloadversion)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `preDeploymentStatus` _string_ | PreDeploymentStatus indicates the current status of the KeptnWorkloadVersion's PreDeployment phase. |Pending| ✓ |  |
| `deploymentStatus` _string_ | DeploymentStatus indicates the current status of the KeptnWorkloadVersion's Deployment phase. |Pending| ✓ |  |
| `preDeploymentEvaluationStatus` _string_ | PreDeploymentEvaluationStatus indicates the current status of the KeptnWorkloadVersion's PreDeploymentEvaluation phase. |Pending| ✓ |  |
| `postDeploymentEvaluationStatus` _string_ | PostDeploymentEvaluationStatus indicates the current status of the KeptnWorkloadVersion's PostDeploymentEvaluation phase. |Pending| ✓ |  |
| `postDeploymentStatus` _string_ | PostDeploymentStatus indicates the current status of the KeptnWorkloadVersion's PostDeployment phase. |Pending| ✓ |  |
| `preDeploymentTaskStatus` _[ItemStatus](#itemstatus) array_ | PreDeploymentTaskStatus indicates the current state of each preDeploymentTask of the KeptnWorkloadVersion. || ✓ |  |
| `postDeploymentTaskStatus` _[ItemStatus](#itemstatus) array_ | PostDeploymentTaskStatus indicates the current state of each postDeploymentTask of the KeptnWorkloadVersion. || ✓ |  |
| `preDeploymentEvaluationTaskStatus` _[ItemStatus](#itemstatus) array_ | PreDeploymentEvaluationTaskStatus indicates the current state of each preDeploymentEvaluation of the KeptnWorkloadVersion. || ✓ |  |
| `postDeploymentEvaluationTaskStatus` _[ItemStatus](#itemstatus) array_ | PostDeploymentEvaluationTaskStatus indicates the current state of each postDeploymentEvaluation of the KeptnWorkloadVersion. || ✓ |  |
| `startTime` _[Time](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#time-v1-meta)_ | StartTime represents the time at which the deployment of the KeptnWorkloadVersion started. || ✓ |  |
| `endTime` _[Time](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#time-v1-meta)_ | EndTime represents the time at which the deployment of the KeptnWorkloadVersion finished. || ✓ |  |
| `currentPhase` _string_ | CurrentPhase indicates the current phase of the KeptnWorkloadVersion. This can be:<br />- PreDeploymentTasks<br />- PreDeploymentEvaluations<br />- Deployment<br />- PostDeploymentTasks<br />- PostDeploymentEvaluations || ✓ |  |
| `phaseTraceIDs` _[MapCarrier](https://pkg.go.dev/go.opentelemetry.io/otel/propagation#MapCarrier)_ | PhaseTraceIDs contains the trace IDs of the OpenTelemetry spans of each phase of the KeptnWorkloadVersion || ✓ |  |
| `status` _string_ | Status represents the overall status of the KeptnWorkloadVersion. |Pending| ✓ |  |
| `appContextMetadata` _object (keys:string, values:string)_ | AppContextMetadata contains metadata from the related KeptnAppVersion. || ✓ |  |
| `deploymentStartTime` _[Time](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#time-v1-meta)_ | DeploymentStartTime represents the start time of the deployment phase || ✓ |  |


#### Objective







_Appears in:_
- [KeptnEvaluationDefinitionSpec](#keptnevaluationdefinitionspec)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `keptnMetricRef` _[KeptnMetricReference](#keptnmetricreference)_ | KeptnMetricRef references the KeptnMetric that should be evaluated. || x |  |
| `evaluationTarget` _string_ | EvaluationTarget specifies the target value for the references KeptnMetric.<br />Needs to start with either '<' or '>', followed by the target value (e.g. '<10'). || x |  |


#### PhaseTraceID

_Underlying type:_ _[MapCarrier](https://pkg.go.dev/go.opentelemetry.io/otel/propagation#MapCarrier)_





_Appears in:_
- [KeptnAppVersionStatus](#keptnappversionstatus)
- [KeptnWorkloadVersionStatus](#keptnworkloadversionstatus)



#### ResourceReference



ResourceReference represents the parent resource of Workload



_Appears in:_
- [KeptnWorkloadSpec](#keptnworkloadspec)
- [KeptnWorkloadVersionSpec](#keptnworkloadversionspec)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `uid` _string_ |  || x |  |
| `kind` _string_ |  || x |  |
| `name` _string_ |  || x |  |


#### RuntimeSpec







_Appears in:_
- [KeptnTaskDefinitionSpec](#keptntaskdefinitionspec)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `functionRef` _[FunctionReference](#functionreference)_ | FunctionReference allows to reference another KeptnTaskDefinition which contains the source code of the<br />function to be executes for KeptnTasks based on this KeptnTaskDefinition. This can be useful when you have<br />multiple KeptnTaskDefinitions that should execute the same logic, but each with different parameters. || ✓ |  |
| `inline` _[Inline](#inline)_ | Inline allows to specify the code that should be executed directly in the KeptnTaskDefinition, as a multi-line<br />string. || ✓ |  |
| `httpRef` _[HttpReference](#httpreference)_ | HttpReference allows to point to an HTTP URL containing the code of the function. || ✓ |  |
| `configMapRef` _[ConfigMapReference](#configmapreference)_ | ConfigMapReference allows to reference a ConfigMap containing the code of the function.<br />When referencing a ConfigMap, the code of the function must be available as a value of the 'code' key<br />of the referenced ConfigMap. || ✓ |  |
| `parameters` _[TaskParameters](#taskparameters)_ | Parameters contains parameters that will be passed to the job that executes the task as env variables. || ✓ |  |
| `secureParameters` _[SecureParameters](#secureparameters)_ | SecureParameters contains secure parameters that will be passed to the job that executes the task.<br />These will be stored and accessed as secrets in the cluster. || ✓ |  |
| `cmdParameters` _string_ | CmdParameters contains parameters that will be passed to the command || ✓ |  |


#### SecureParameters







_Appears in:_
- [KeptnTaskSpec](#keptntaskspec)
- [RuntimeSpec](#runtimespec)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `secret` _string_ | Secret contains the parameters that will be made available to the job<br />executing the KeptnTask via the 'SECRET_DATA' environment variable.<br />The 'SECRET_DATA'  environment variable's content will the same as value of the 'SECRET_DATA'<br />key of the referenced secret. || ✓ |  |


#### ServiceAccountSpec







_Appears in:_
- [KeptnTaskDefinitionSpec](#keptntaskdefinitionspec)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `name` _string_ |  || x |  |




#### TaskContext







_Appears in:_
- [KeptnTaskSpec](#keptntaskspec)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `workloadName` _string_ | WorkloadName the name of the KeptnWorkload the KeptnTask is being executed for. || ✓ |  |
| `appName` _string_ | AppName the name of the KeptnApp the KeptnTask is being executed for. || ✓ |  |
| `appVersion` _string_ | AppVersion the version of the KeptnApp the KeptnTask is being executed for. || ✓ |  |
| `workloadVersion` _string_ | WorkloadVersion the version of the KeptnWorkload the KeptnTask is being executed for. || ✓ |  |
| `taskType` _string_ | TaskType indicates whether the KeptnTask is part of the pre- or postDeployment phase. || ✓ |  |
| `objectType` _string_ | ObjectType indicates whether the KeptnTask is being executed for a KeptnApp or KeptnWorkload. || ✓ |  |
| `metadata` _object (keys:string, values:string)_ | Metadata contains additional key-value pairs for contextual information. || ✓ |  |


#### TaskParameters







_Appears in:_
- [KeptnTaskSpec](#keptntaskspec)
- [RuntimeSpec](#runtimespec)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `map` _object (keys:string, values:string)_ | Inline contains the parameters that will be made available to the job<br />executing the KeptnTask via the 'DATA' environment variable.<br />The 'DATA'  environment variable's content will be a json<br />encoded string containing all properties of the map provided. || ✓ |  |


#### WorkloadStatus







_Appears in:_
- [KeptnAppVersionStatus](#keptnappversionstatus)

| Field | Description | Default | Optional |Validation |
| --- | --- | --- | --- | --- |
| `workload` _[KeptnWorkloadRef](#keptnworkloadref)_ | Workload refers to a KeptnWorkload that is part of the KeptnAppVersion. || ✓ |  |
| `status` _string_ | Status indicates the current status of the KeptnWorkload. |Pending| ✓ |  |


