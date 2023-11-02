---
title: v1alpha3
description: Reference information for lifecycle.keptn.sh/v1alpha3
---
<!-- markdownlint-disable -->

## Packages
- [lifecycle.keptn.sh/v1alpha3](#lifecyclekeptnshv1alpha3)


## lifecycle.keptn.sh/v1alpha3

Package v1alpha3 contains API Schema definitions for the lifecycle v1alpha3 API group

### Resource Types
- [KeptnApp](#keptnapp)
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
- [KeptnWorkloadInstance](#keptnworkloadinstance)
- [KeptnWorkloadInstanceList](#keptnworkloadinstancelist)
- [KeptnWorkloadList](#keptnworkloadlist)



#### AutomountServiceAccountTokenSpec





_Appears in:_
- [KeptnTaskDefinitionSpec](#keptntaskdefinitionspec)

| Field | Description |
| --- | --- |
| `type` _boolean_ |  |


#### ConfigMapReference





_Appears in:_
- [RuntimeSpec](#runtimespec)

| Field | Description |
| --- | --- |
| `name` _string_ | Name is the name of the referenced ConfigMap. |


#### ContainerSpec





_Appears in:_
- [KeptnTaskDefinitionSpec](#keptntaskdefinitionspec)

| Field | Description |
| --- | --- |
| `name` _string_ | Name of the container specified as a DNS_LABEL. Each container in a pod must have a unique name (DNS_LABEL). Cannot be updated. |
| `image` _string_ | Container image name. More info: https://kubernetes.io/docs/concepts/containers/images This field is optional to allow higher level config management to default or override container images in workload controllers like Deployments and StatefulSets. |
| `command` _string array_ | Entrypoint array. Not executed within a shell. The container image's ENTRYPOINT is used if this is not provided. Variable references $(VAR_NAME) are expanded using the container's environment. If a variable cannot be resolved, the reference in the input string will be unchanged. Double $$ are reduced to a single $, which allows for escaping the $(VAR_NAME) syntax: i.e. "$$(VAR_NAME)" will produce the string literal "$(VAR_NAME)". Escaped references will never be expanded, regardless of whether the variable exists or not. Cannot be updated. More info: https://kubernetes.io/docs/tasks/inject-data-application/define-command-argument-container/#running-a-command-in-a-shell |
| `args` _string array_ | Arguments to the entrypoint. The container image's CMD is used if this is not provided. Variable references $(VAR_NAME) are expanded using the container's environment. If a variable cannot be resolved, the reference in the input string will be unchanged. Double $$ are reduced to a single $, which allows for escaping the $(VAR_NAME) syntax: i.e. "$$(VAR_NAME)" will produce the string literal "$(VAR_NAME)". Escaped references will never be expanded, regardless of whether the variable exists or not. Cannot be updated. More info: https://kubernetes.io/docs/tasks/inject-data-application/define-command-argument-container/#running-a-command-in-a-shell |
| `workingDir` _string_ | Container's working directory. If not specified, the container runtime's default will be used, which might be configured in the container image. Cannot be updated. |
| `ports` _[ContainerPort](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#containerport-v1-core) array_ | List of ports to expose from the container. Not specifying a port here DOES NOT prevent that port from being exposed. Any port which is listening on the default "0.0.0.0" address inside a container will be accessible from the network. Modifying this array with strategic merge patch may corrupt the data. For more information See https://github.com/kubernetes/kubernetes/issues/108255. Cannot be updated. |
| `envFrom` _[EnvFromSource](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#envfromsource-v1-core) array_ | List of sources to populate environment variables in the container. The keys defined within a source must be a C_IDENTIFIER. All invalid keys will be reported as an event when the container is starting. When a key exists in multiple sources, the value associated with the last source will take precedence. Values defined by an Env with a duplicate key will take precedence. Cannot be updated. |
| `env` _[EnvVar](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#envvar-v1-core) array_ | List of environment variables to set in the container. Cannot be updated. |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#resourcerequirements-v1-core)_ | Compute Resources required by this container. Cannot be updated. More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/ |
| `resizePolicy` _[ContainerResizePolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#containerresizepolicy-v1-core) array_ | Resources resize policy for the container. |
| `restartPolicy` _[ContainerRestartPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#containerrestartpolicy-v1-core)_ | RestartPolicy defines the restart behavior of individual containers in a pod. This field may only be set for init containers, and the only allowed value is "Always". For non-init containers or when this field is not specified, the restart behavior is defined by the Pod's restart policy and the container type. Setting the RestartPolicy as "Always" for the init container will have the following effect: this init container will be continually restarted on exit until all regular containers have terminated. Once all regular containers have completed, all init containers with restartPolicy "Always" will be shut down. This lifecycle differs from normal init containers and is often referred to as a "sidecar" container. Although this init container still starts in the init container sequence, it does not wait for the container to complete before proceeding to the next init container. Instead, the next init container starts immediately after this init container is started, or after any startupProbe has successfully completed. |
| `volumeMounts` _[VolumeMount](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#volumemount-v1-core) array_ | Pod volumes to mount into the container's filesystem. Cannot be updated. |
| `volumeDevices` _[VolumeDevice](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#volumedevice-v1-core) array_ | volumeDevices is the list of block devices to be used by the container. |
| `livenessProbe` _[Probe](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#probe-v1-core)_ | Periodic probe of container liveness. Container will be restarted if the probe fails. Cannot be updated. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes |
| `readinessProbe` _[Probe](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#probe-v1-core)_ | Periodic probe of container service readiness. Container will be removed from service endpoints if the probe fails. Cannot be updated. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes |
| `startupProbe` _[Probe](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#probe-v1-core)_ | StartupProbe indicates that the Pod has successfully initialized. If specified, no other probes are executed until this completes successfully. If this probe fails, the Pod will be restarted, just as if the livenessProbe failed. This can be used to provide different probe parameters at the beginning of a Pod's lifecycle, when it might take a long time to load data or warm a cache, than during steady-state operation. This cannot be updated. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes |
| `lifecycle` _[Lifecycle](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#lifecycle-v1-core)_ | Actions that the management system should take in response to container lifecycle events. Cannot be updated. |
| `terminationMessagePath` _string_ | Optional: Path at which the file to which the container's termination message will be written is mounted into the container's filesystem. Message written is intended to be brief final status, such as an assertion failure message. Will be truncated by the node if greater than 4096 bytes. The total message length across all containers will be limited to 12kb. Defaults to /dev/termination-log. Cannot be updated. |
| `terminationMessagePolicy` _[TerminationMessagePolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#terminationmessagepolicy-v1-core)_ | Indicate how the termination message should be populated. File will use the contents of terminationMessagePath to populate the container status message on both success and failure. FallbackToLogsOnError will use the last chunk of container log output if the termination message file is empty and the container exited with an error. The log output is limited to 2048 bytes or 80 lines, whichever is smaller. Defaults to File. Cannot be updated. |
| `imagePullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#pullpolicy-v1-core)_ | Image pull policy. One of Always, Never, IfNotPresent. Defaults to Always if :latest tag is specified, or IfNotPresent otherwise. Cannot be updated. More info: https://kubernetes.io/docs/concepts/containers/images#updating-images |
| `securityContext` _[SecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#securitycontext-v1-core)_ | SecurityContext defines the security options the container should be run with. If set, the fields of SecurityContext override the equivalent fields of PodSecurityContext. More info: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/ |
| `stdin` _boolean_ | Whether this container should allocate a buffer for stdin in the container runtime. If this is not set, reads from stdin in the container will always result in EOF. Default is false. |
| `stdinOnce` _boolean_ | Whether the container runtime should close the stdin channel after it has been opened by a single attach. When stdin is true the stdin stream will remain open across multiple attach sessions. If stdinOnce is set to true, stdin is opened on container start, is empty until the first client attaches to stdin, and then remains open and accepts data until the client disconnects, at which time stdin is closed and remains closed until the container is restarted. If this flag is false, a container processes that reads from stdin will never receive an EOF. Default is false |
| `tty` _boolean_ | Whether this container should allocate a TTY for itself, also requires 'stdin' to be true. Default is false. |


#### EvaluationStatusItem





_Appears in:_
- [KeptnEvaluationStatus](#keptnevaluationstatus)

| Field | Description |
| --- | --- |
| `value` _string_ | Value represents the value of the KeptnMetric being evaluated. |
| `status` _KeptnState_ | Status indicates the status of the objective being evaluated. |
| `message` _string_ | Message contains additional information about the evaluation of an objective. This can include explanations about why an evaluation has failed (e.g. due to a missed objective), or if there was any error during the evaluation of the objective. |


#### FunctionReference





_Appears in:_
- [RuntimeSpec](#runtimespec)

| Field | Description |
| --- | --- |
| `name` _string_ | Name is the name of the referenced KeptnTaskDefinition. |


#### FunctionStatus





_Appears in:_
- [KeptnTaskDefinitionStatus](#keptntaskdefinitionstatus)

| Field | Description |
| --- | --- |
| `configMap` _string_ | ConfigMap indicates the ConfigMap in which the function code is stored. |


#### HttpReference





_Appears in:_
- [RuntimeSpec](#runtimespec)

| Field | Description |
| --- | --- |
| `url` _string_ | Url is the URL containing the code of the function. |


#### Inline





_Appears in:_
- [RuntimeSpec](#runtimespec)

| Field | Description |
| --- | --- |
| `code` _string_ | Code contains the code of the function. |


#### ItemStatus





_Appears in:_
- [KeptnAppVersionStatus](#keptnappversionstatus)
- [KeptnWorkloadInstanceStatus](#keptnworkloadinstancestatus)

| Field | Description |
| --- | --- |
| `definitionName` _string_ | DefinitionName is the name of the EvaluationDefinition/TaskDefinition |
| `status` _KeptnState_ |  |
| `name` _string_ | Name is the name of the Evaluation/Task |
| `startTime` _[Time](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#time-v1-meta)_ | StartTime represents the time at which the Item (Evaluation/Task) started. |
| `endTime` _[Time](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#time-v1-meta)_ | EndTime represents the time at which the Item (Evaluation/Task) started. |


#### KeptnApp



KeptnApp is the Schema for the keptnapps API

_Appears in:_
- [KeptnAppList](#keptnapplist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `lifecycle.keptn.sh/v1alpha3`
| `kind` _string_ | `KeptnApp`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[KeptnAppSpec](#keptnappspec)_ | Spec describes the desired state of the KeptnApp. |
| `status` _[KeptnAppStatus](#keptnappstatus)_ | Status describes the current state of the KeptnApp. |


#### KeptnAppCreationRequest



KeptnAppCreationRequest is the Schema for the keptnappcreationrequests API

_Appears in:_
- [KeptnAppCreationRequestList](#keptnappcreationrequestlist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `lifecycle.keptn.sh/v1alpha3`
| `kind` _string_ | `KeptnAppCreationRequest`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[KeptnAppCreationRequestSpec](#keptnappcreationrequestspec)_ | Spec describes the desired state of the KeptnAppCreationRequest. |
| `status` _string_ | Status describes the current state of the KeptnAppCreationRequest. |


#### KeptnAppCreationRequestList



KeptnAppCreationRequestList contains a list of KeptnAppCreationRequest



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `lifecycle.keptn.sh/v1alpha3`
| `kind` _string_ | `KeptnAppCreationRequestList`
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[KeptnAppCreationRequest](#keptnappcreationrequest) array_ |  |


#### KeptnAppCreationRequestSpec



KeptnAppCreationRequestSpec defines the desired state of KeptnAppCreationRequest

_Appears in:_
- [KeptnAppCreationRequest](#keptnappcreationrequest)

| Field | Description |
| --- | --- |
| `appName` _string_ | AppName is the name of the KeptnApp the KeptnAppCreationRequest should create if no user-defined object with that name is found. |


#### KeptnAppList



KeptnAppList contains a list of KeptnApp



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `lifecycle.keptn.sh/v1alpha3`
| `kind` _string_ | `KeptnAppList`
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[KeptnApp](#keptnapp) array_ |  |


#### KeptnAppSpec



KeptnAppSpec defines the desired state of KeptnApp

_Appears in:_
- [KeptnApp](#keptnapp)
- [KeptnAppVersionSpec](#keptnappversionspec)

| Field | Description |
| --- | --- |
| `version` _string_ | Version defines the version of the application. For automatically created KeptnApps, the version is a function of all KeptnWorkloads that are part of the KeptnApp. |
| `revision` _integer_ | Revision can be modified to trigger another deployment of a KeptnApp of the same version. This can be used for restarting a KeptnApp which failed to deploy, e.g. due to a failed preDeploymentEvaluation/preDeploymentTask. |
| `workloads` _[KeptnWorkloadRef](#keptnworkloadref) array_ | Workloads is a list of all KeptnWorkloads that are part of the KeptnApp. |
| `preDeploymentTasks` _string array_ | PreDeploymentTasks is a list of all tasks to be performed during the pre-deployment phase of the KeptnApp. The items of this list refer to the names of KeptnTaskDefinitions located in the same namespace as the KeptnApp, or in the Keptn namespace. |
| `postDeploymentTasks` _string array_ | PostDeploymentTasks is a list of all tasks to be performed during the post-deployment phase of the KeptnApp. The items of this list refer to the names of KeptnTaskDefinitions located in the same namespace as the KeptnApp, or in the Keptn namespace. |
| `preDeploymentEvaluations` _string array_ | PreDeploymentEvaluations is a list of all evaluations to be performed during the pre-deployment phase of the KeptnApp. The items of this list refer to the names of KeptnEvaluationDefinitions located in the same namespace as the KeptnApp, or in the Keptn namespace. |
| `postDeploymentEvaluations` _string array_ | PostDeploymentEvaluations is a list of all evaluations to be performed during the post-deployment phase of the KeptnApp. The items of this list refer to the names of KeptnEvaluationDefinitions located in the same namespace as the KeptnApp, or in the Keptn namespace. |


#### KeptnAppStatus



KeptnAppStatus defines the observed state of KeptnApp

_Appears in:_
- [KeptnApp](#keptnapp)

| Field | Description |
| --- | --- |
| `currentVersion` _string_ | CurrentVersion indicates the version that is currently deployed or being reconciled. |


#### KeptnAppVersion



KeptnAppVersion is the Schema for the keptnappversions API

_Appears in:_
- [KeptnAppVersionList](#keptnappversionlist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `lifecycle.keptn.sh/v1alpha3`
| `kind` _string_ | `KeptnAppVersion`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[KeptnAppVersionSpec](#keptnappversionspec)_ | Spec describes the desired state of the KeptnAppVersion. |
| `status` _[KeptnAppVersionStatus](#keptnappversionstatus)_ | Status describes the current state of the KeptnAppVersion. |


#### KeptnAppVersionList



KeptnAppVersionList contains a list of KeptnAppVersion



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `lifecycle.keptn.sh/v1alpha3`
| `kind` _string_ | `KeptnAppVersionList`
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[KeptnAppVersion](#keptnappversion) array_ |  |


#### KeptnAppVersionSpec



KeptnAppVersionSpec defines the desired state of KeptnAppVersion

_Appears in:_
- [KeptnAppVersion](#keptnappversion)

| Field | Description |
| --- | --- |
| `version` _string_ | Version defines the version of the application. For automatically created KeptnApps, the version is a function of all KeptnWorkloads that are part of the KeptnApp. |
| `revision` _integer_ | Revision can be modified to trigger another deployment of a KeptnApp of the same version. This can be used for restarting a KeptnApp which failed to deploy, e.g. due to a failed preDeploymentEvaluation/preDeploymentTask. |
| `workloads` _[KeptnWorkloadRef](#keptnworkloadref) array_ | Workloads is a list of all KeptnWorkloads that are part of the KeptnApp. |
| `preDeploymentTasks` _string array_ | PreDeploymentTasks is a list of all tasks to be performed during the pre-deployment phase of the KeptnApp. The items of this list refer to the names of KeptnTaskDefinitions located in the same namespace as the KeptnApp, or in the Keptn namespace. |
| `postDeploymentTasks` _string array_ | PostDeploymentTasks is a list of all tasks to be performed during the post-deployment phase of the KeptnApp. The items of this list refer to the names of KeptnTaskDefinitions located in the same namespace as the KeptnApp, or in the Keptn namespace. |
| `preDeploymentEvaluations` _string array_ | PreDeploymentEvaluations is a list of all evaluations to be performed during the pre-deployment phase of the KeptnApp. The items of this list refer to the names of KeptnEvaluationDefinitions located in the same namespace as the KeptnApp, or in the Keptn namespace. |
| `postDeploymentEvaluations` _string array_ | PostDeploymentEvaluations is a list of all evaluations to be performed during the post-deployment phase of the KeptnApp. The items of this list refer to the names of KeptnEvaluationDefinitions located in the same namespace as the KeptnApp, or in the Keptn namespace. |
| `appName` _string_ | AppName is the name of the KeptnApp. |
| `previousVersion` _string_ | PreviousVersion is the version of the KeptnApp that has been deployed prior to this version. |
| `traceId` _object (keys:string, values:string)_ | TraceId contains the OpenTelemetry trace ID. |


#### KeptnAppVersionStatus



KeptnAppVersionStatus defines the observed state of KeptnAppVersion

_Appears in:_
- [KeptnAppVersion](#keptnappversion)

| Field | Description |
| --- | --- |
| `preDeploymentStatus` _KeptnState_ | PreDeploymentStatus indicates the current status of the KeptnAppVersion's PreDeployment phase. |
| `postDeploymentStatus` _KeptnState_ | PostDeploymentStatus indicates the current status of the KeptnAppVersion's PostDeployment phase. |
| `preDeploymentEvaluationStatus` _KeptnState_ | PreDeploymentEvaluationStatus indicates the current status of the KeptnAppVersion's PreDeploymentEvaluation phase. |
| `postDeploymentEvaluationStatus` _KeptnState_ | PostDeploymentEvaluationStatus indicates the current status of the KeptnAppVersion's PostDeploymentEvaluation phase. |
| `workloadOverallStatus` _KeptnState_ | WorkloadOverallStatus indicates the current status of the KeptnAppVersion's Workload deployment phase. |
| `workloadStatus` _[WorkloadStatus](#workloadstatus) array_ | WorkloadStatus contains the current status of each KeptnWorkload that is part of the KeptnAppVersion. |
| `currentPhase` _string_ | CurrentPhase indicates the current phase of the KeptnAppVersion. |
| `preDeploymentTaskStatus` _[ItemStatus](#itemstatus) array_ | PreDeploymentTaskStatus indicates the current state of each preDeploymentTask of the KeptnAppVersion. |
| `postDeploymentTaskStatus` _[ItemStatus](#itemstatus) array_ | PostDeploymentTaskStatus indicates the current state of each postDeploymentTask of the KeptnAppVersion. |
| `preDeploymentEvaluationTaskStatus` _[ItemStatus](#itemstatus) array_ | PreDeploymentEvaluationTaskStatus indicates the current state of each preDeploymentEvaluation of the KeptnAppVersion. |
| `postDeploymentEvaluationTaskStatus` _[ItemStatus](#itemstatus) array_ | PostDeploymentEvaluationTaskStatus indicates the current state of each postDeploymentEvaluation of the KeptnAppVersion. |
| `phaseTraceIDs` _object (keys:string, values:object)_ | PhaseTraceIDs contains the trace IDs of the OpenTelemetry spans of each phase of the KeptnAppVersion. |
| `status` _KeptnState_ | Status represents the overall status of the KeptnAppVersion. |
| `startTime` _[Time](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#time-v1-meta)_ | StartTime represents the time at which the deployment of the KeptnAppVersion started. |
| `endTime` _[Time](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#time-v1-meta)_ | EndTime represents the time at which the deployment of the KeptnAppVersion finished. |


#### KeptnEvaluation



KeptnEvaluation is the Schema for the keptnevaluations API

_Appears in:_
- [KeptnEvaluationList](#keptnevaluationlist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `lifecycle.keptn.sh/v1alpha3`
| `kind` _string_ | `KeptnEvaluation`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[KeptnEvaluationSpec](#keptnevaluationspec)_ | Spec describes the desired state of the KeptnEvaluation. |
| `status` _[KeptnEvaluationStatus](#keptnevaluationstatus)_ | Status describes the current state of the KeptnEvaluation. |


#### KeptnEvaluationDefinition



KeptnEvaluationDefinition is the Schema for the keptnevaluationdefinitions API

_Appears in:_
- [KeptnEvaluationDefinitionList](#keptnevaluationdefinitionlist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `lifecycle.keptn.sh/v1alpha3`
| `kind` _string_ | `KeptnEvaluationDefinition`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[KeptnEvaluationDefinitionSpec](#keptnevaluationdefinitionspec)_ | Spec describes the desired state of the KeptnEvaluationDefinition. |
| `status` _string_ | unused field |


#### KeptnEvaluationDefinitionList



KeptnEvaluationDefinitionList contains a list of KeptnEvaluationDefinition



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `lifecycle.keptn.sh/v1alpha3`
| `kind` _string_ | `KeptnEvaluationDefinitionList`
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[KeptnEvaluationDefinition](#keptnevaluationdefinition) array_ |  |


#### KeptnEvaluationDefinitionSpec



KeptnEvaluationDefinitionSpec defines the desired state of KeptnEvaluationDefinition

_Appears in:_
- [KeptnEvaluationDefinition](#keptnevaluationdefinition)

| Field | Description |
| --- | --- |
| `objectives` _[Objective](#objective) array_ | Objectives is a list of objectives that have to be met for a KeptnEvaluation referencing this KeptnEvaluationDefinition to be successful. |


#### KeptnEvaluationList



KeptnEvaluationList contains a list of KeptnEvaluation



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `lifecycle.keptn.sh/v1alpha3`
| `kind` _string_ | `KeptnEvaluationList`
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[KeptnEvaluation](#keptnevaluation) array_ |  |




#### KeptnEvaluationSpec



KeptnEvaluationSpec defines the desired state of KeptnEvaluation

_Appears in:_
- [KeptnEvaluation](#keptnevaluation)

| Field | Description |
| --- | --- |
| `workload` _string_ | Workload defines the KeptnWorkload for which the KeptnEvaluation is done. |
| `workloadVersion` _string_ | WorkloadVersion defines the version of the KeptnWorkload for which the KeptnEvaluation is done. |
| `appName` _string_ | AppName defines the KeptnApp for which the KeptnEvaluation is done. |
| `appVersion` _string_ | AppVersion defines the version of the KeptnApp for which the KeptnEvaluation is done. |
| `evaluationDefinition` _string_ | EvaluationDefinition refers to the name of the KeptnEvaluationDefinition which includes the objectives for the KeptnEvaluation. The KeptnEvaluationDefinition can be located in the same namespace as the KeptnEvaluation, or in the Keptn namespace. |
| `retries` _integer_ | Retries indicates how many times the KeptnEvaluation can be attempted in the case of an error or missed evaluation objective, before considering the KeptnEvaluation to be failed. |
| `retryInterval` _[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#duration-v1-meta)_ | RetryInterval specifies the interval at which the KeptnEvaluation is retried in the case of an error or a missed objective. |
| `failAction` _string_ |  |
| `checkType` _CheckType_ | Type indicates whether the KeptnEvaluation is part of the pre- or postDeployment phase. |


#### KeptnEvaluationStatus



KeptnEvaluationStatus defines the observed state of KeptnEvaluation

_Appears in:_
- [KeptnEvaluation](#keptnevaluation)

| Field | Description |
| --- | --- |
| `retryCount` _integer_ | RetryCount indicates how many times the KeptnEvaluation has been attempted already. |
| `evaluationStatus` _object (keys:string, values:[EvaluationStatusItem](#evaluationstatusitem))_ | EvaluationStatus describes the status of each objective of the KeptnEvaluationDefinition referenced by the KeptnEvaluation. |
| `overallStatus` _KeptnState_ | OverallStatus describes the overall status of the KeptnEvaluation. The Overall status is derived from the status of the individual objectives of the KeptnEvaluationDefinition referenced by the KeptnEvaluation. |
| `startTime` _[Time](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#time-v1-meta)_ | StartTime represents the time at which the KeptnEvaluation started. |
| `endTime` _[Time](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#time-v1-meta)_ | EndTime represents the time at which the KeptnEvaluation finished. |


#### KeptnMetricReference





_Appears in:_
- [Objective](#objective)

| Field | Description |
| --- | --- |
| `name` _string_ | Name is the name of the referenced KeptnMetric. |
| `namespace` _string_ | Namespace is the namespace where the referenced KeptnMetric is located. |


#### KeptnTask



KeptnTask is the Schema for the keptntasks API

_Appears in:_
- [KeptnTaskList](#keptntasklist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `lifecycle.keptn.sh/v1alpha3`
| `kind` _string_ | `KeptnTask`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[KeptnTaskSpec](#keptntaskspec)_ | Spec describes the desired state of the KeptnTask. |
| `status` _[KeptnTaskStatus](#keptntaskstatus)_ | Status describes the current state of the KeptnTask. |


#### KeptnTaskDefinition



KeptnTaskDefinition is the Schema for the keptntaskdefinitions API

_Appears in:_
- [KeptnTaskDefinitionList](#keptntaskdefinitionlist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `lifecycle.keptn.sh/v1alpha3`
| `kind` _string_ | `KeptnTaskDefinition`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[KeptnTaskDefinitionSpec](#keptntaskdefinitionspec)_ | Spec describes the desired state of the KeptnTaskDefinition. |
| `status` _[KeptnTaskDefinitionStatus](#keptntaskdefinitionstatus)_ | Status describes the current state of the KeptnTaskDefinition. |


#### KeptnTaskDefinitionList



KeptnTaskDefinitionList contains a list of KeptnTaskDefinition



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `lifecycle.keptn.sh/v1alpha3`
| `kind` _string_ | `KeptnTaskDefinitionList`
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[KeptnTaskDefinition](#keptntaskdefinition) array_ |  |


#### KeptnTaskDefinitionSpec



KeptnTaskDefinitionSpec defines the desired state of KeptnTaskDefinition

_Appears in:_
- [KeptnTaskDefinition](#keptntaskdefinition)

| Field | Description |
| --- | --- |
| `function` _[RuntimeSpec](#runtimespec)_ | Deprecated Function contains the definition for the function that is to be executed in KeptnTasks based on the KeptnTaskDefinitions. |
| `python` _[RuntimeSpec](#runtimespec)_ | Python contains the definition for the python function that is to be executed in KeptnTasks based on the KeptnTaskDefinitions. |
| `deno` _[RuntimeSpec](#runtimespec)_ | Deno contains the definition for the Deno function that is to be executed in KeptnTasks based on the KeptnTaskDefinitions. |
| `container` _[ContainerSpec](#containerspec)_ | Container contains the definition for the container that is to be used in Job based on the KeptnTaskDefinitions. |
| `retries` _integer_ | Retries specifies how many times a job executing the KeptnTaskDefinition should be restarted in the case of an unsuccessful attempt. |
| `timeout` _[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#duration-v1-meta)_ | Timeout specifies the maximum time to wait for the task to be completed successfully. If the task does not complete successfully within this time frame, it will be considered to be failed. |
| `serviceAccount` _[ServiceAccountSpec](#serviceaccountspec)_ | ServiceAccount specifies the service account to be used in jobs to authenticate with the Kubernetes API and access cluster resources. |
| `automountServiceAccountToken` _[AutomountServiceAccountTokenSpec](#automountserviceaccounttokenspec)_ | AutomountServiceAccountToken allows to enable K8s to assign cluster API credentials to a pod, if set to false the pod will decline the service account |


#### KeptnTaskDefinitionStatus



KeptnTaskDefinitionStatus defines the observed state of KeptnTaskDefinition

_Appears in:_
- [KeptnTaskDefinition](#keptntaskdefinition)

| Field | Description |
| --- | --- |
| `function` _[FunctionStatus](#functionstatus)_ | Function contains status information of the function definition for the task. |


#### KeptnTaskList



KeptnTaskList contains a list of KeptnTask



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `lifecycle.keptn.sh/v1alpha3`
| `kind` _string_ | `KeptnTaskList`
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[KeptnTask](#keptntask) array_ |  |


#### KeptnTaskSpec



KeptnTaskSpec defines the desired state of KeptnTask

_Appears in:_
- [KeptnTask](#keptntask)

| Field | Description |
| --- | --- |
| `taskDefinition` _string_ | TaskDefinition refers to the name of the KeptnTaskDefinition which includes the specification for the task to be performed. The KeptnTaskDefinition can be located in the same namespace as the KeptnTask, or in the Keptn namespace. |
| `context` _[TaskContext](#taskcontext)_ | Context contains contextual information about the task execution. |
| `parameters` _[TaskParameters](#taskparameters)_ | Parameters contains parameters that will be passed to the job that executes the task. |
| `secureParameters` _[SecureParameters](#secureparameters)_ | SecureParameters contains secure parameters that will be passed to the job that executes the task. These will be stored and accessed as secrets in the cluster. |
| `checkType` _CheckType_ | Type indicates whether the KeptnTask is part of the pre- or postDeployment phase. |
| `retries` _integer_ | Retries indicates how many times the KeptnTask can be attempted in the case of an error before considering the KeptnTask to be failed. |
| `timeout` _[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#duration-v1-meta)_ | Timeout specifies the maximum time to wait for the task to be completed successfully. If the task does not complete successfully within this time frame, it will be considered to be failed. |


#### KeptnTaskStatus



KeptnTaskStatus defines the observed state of KeptnTask

_Appears in:_
- [KeptnTask](#keptntask)

| Field | Description |
| --- | --- |
| `jobName` _string_ | JobName is the name of the Job executing the Task. |
| `status` _KeptnState_ | Status represents the overall state of the KeptnTask. |
| `message` _string_ | Message contains information about unexpected errors encountered during the execution of the KeptnTask. |
| `startTime` _[Time](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#time-v1-meta)_ | StartTime represents the time at which the KeptnTask started. |
| `endTime` _[Time](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#time-v1-meta)_ | EndTime represents the time at which the KeptnTask finished. |
| `reason` _string_ | Reason contains more information about the reason for the last transition of the Job executing the KeptnTask. |


#### KeptnWorkload



KeptnWorkload is the Schema for the keptnworkloads API

_Appears in:_
- [KeptnWorkloadList](#keptnworkloadlist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `lifecycle.keptn.sh/v1alpha3`
| `kind` _string_ | `KeptnWorkload`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[KeptnWorkloadSpec](#keptnworkloadspec)_ | Spec describes the desired state of the KeptnWorkload. |
| `status` _[KeptnWorkloadStatus](#keptnworkloadstatus)_ | Status describes the current state of the KeptnWorkload. |


#### KeptnWorkloadInstance



KeptnWorkloadInstance is the Schema for the keptnworkloadinstances API

_Appears in:_
- [KeptnWorkloadInstanceList](#keptnworkloadinstancelist)

| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `lifecycle.keptn.sh/v1alpha3`
| `kind` _string_ | `KeptnWorkloadInstance`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[KeptnWorkloadInstanceSpec](#keptnworkloadinstancespec)_ | Spec describes the desired state of the KeptnWorkloadInstance. |
| `status` _[KeptnWorkloadInstanceStatus](#keptnworkloadinstancestatus)_ | Status describes the current state of the KeptnWorkloadInstance. |


#### KeptnWorkloadInstanceList



KeptnWorkloadInstanceList contains a list of KeptnWorkloadInstance



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `lifecycle.keptn.sh/v1alpha3`
| `kind` _string_ | `KeptnWorkloadInstanceList`
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[KeptnWorkloadInstance](#keptnworkloadinstance) array_ |  |


#### KeptnWorkloadInstanceSpec



KeptnWorkloadInstanceSpec defines the desired state of KeptnWorkloadInstance

_Appears in:_
- [KeptnWorkloadInstance](#keptnworkloadinstance)

| Field | Description |
| --- | --- |
| `app` _string_ | AppName is the name of the KeptnApp containing the KeptnWorkload. |
| `version` _string_ | Version defines the version of the KeptnWorkload. |
| `preDeploymentTasks` _string array_ | PreDeploymentTasks is a list of all tasks to be performed during the pre-deployment phase of the KeptnWorkload. The items of this list refer to the names of KeptnTaskDefinitions located in the same namespace as the KeptnApp, or in the Keptn namespace. |
| `postDeploymentTasks` _string array_ | PostDeploymentTasks is a list of all tasks to be performed during the post-deployment phase of the KeptnWorkload. The items of this list refer to the names of KeptnTaskDefinitions located in the same namespace as the KeptnWorkload, or in the Keptn namespace. |
| `preDeploymentEvaluations` _string array_ | PreDeploymentEvaluations is a list of all evaluations to be performed during the pre-deployment phase of the KeptnWorkload. The items of this list refer to the names of KeptnEvaluationDefinitions located in the same namespace as the KeptnWorkload, or in the Keptn namespace. |
| `postDeploymentEvaluations` _string array_ | PostDeploymentEvaluations is a list of all evaluations to be performed during the post-deployment phase of the KeptnWorkload. The items of this list refer to the names of KeptnEvaluationDefinitions located in the same namespace as the KeptnWorkload, or in the Keptn namespace. |
| `resourceReference` _[ResourceReference](#resourcereference)_ | ResourceReference is a reference to the Kubernetes resource (Deployment, DaemonSet, StatefulSet or ReplicaSet) the KeptnWorkload is representing. |
| `workloadName` _string_ | WorkloadName is the name of the KeptnWorkload. |
| `previousVersion` _string_ | PreviousVersion is the version of the KeptnWorkload that has been deployed prior to this version. |
| `traceId` _object (keys:string, values:string)_ | TraceId contains the OpenTelemetry trace ID. |


#### KeptnWorkloadInstanceStatus



KeptnWorkloadInstanceStatus defines the observed state of KeptnWorkloadInstance

_Appears in:_
- [KeptnWorkloadInstance](#keptnworkloadinstance)

| Field | Description |
| --- | --- |
| `preDeploymentStatus` _KeptnState_ | PreDeploymentStatus indicates the current status of the KeptnWorkloadInstance's PreDeployment phase. |
| `deploymentStatus` _KeptnState_ | DeploymentStatus indicates the current status of the KeptnWorkloadInstance's Deployment phase. |
| `preDeploymentEvaluationStatus` _KeptnState_ | PreDeploymentEvaluationStatus indicates the current status of the KeptnWorkloadInstance's PreDeploymentEvaluation phase. |
| `postDeploymentEvaluationStatus` _KeptnState_ | PostDeploymentEvaluationStatus indicates the current status of the KeptnWorkloadInstance's PostDeploymentEvaluation phase. |
| `postDeploymentStatus` _KeptnState_ | PostDeploymentStatus indicates the current status of the KeptnWorkloadInstance's PostDeployment phase. |
| `preDeploymentTaskStatus` _[ItemStatus](#itemstatus) array_ | PreDeploymentTaskStatus indicates the current state of each preDeploymentTask of the KeptnWorkloadInstance. |
| `postDeploymentTaskStatus` _[ItemStatus](#itemstatus) array_ | PostDeploymentTaskStatus indicates the current state of each postDeploymentTask of the KeptnWorkloadInstance. |
| `preDeploymentEvaluationTaskStatus` _[ItemStatus](#itemstatus) array_ | PreDeploymentEvaluationTaskStatus indicates the current state of each preDeploymentEvaluation of the KeptnWorkloadInstance. |
| `postDeploymentEvaluationTaskStatus` _[ItemStatus](#itemstatus) array_ | PostDeploymentEvaluationTaskStatus indicates the current state of each postDeploymentEvaluation of the KeptnWorkloadInstance. |
| `startTime` _[Time](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#time-v1-meta)_ | StartTime represents the time at which the deployment of the KeptnWorkloadInstance started. |
| `endTime` _[Time](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#time-v1-meta)_ | EndTime represents the time at which the deployment of the KeptnWorkloadInstance finished. |
| `currentPhase` _string_ | CurrentPhase indicates the current phase of the KeptnWorkloadInstance. This can be: - PreDeploymentTasks - PreDeploymentEvaluations - Deployment - PostDeploymentTasks - PostDeploymentEvaluations |
| `phaseTraceIDs` _object (keys:string, values:object)_ | PhaseTraceIDs contains the trace IDs of the OpenTelemetry spans of each phase of the KeptnWorkloadInstance |
| `status` _KeptnState_ | Status represents the overall status of the KeptnWorkloadInstance. |


#### KeptnWorkloadList



KeptnWorkloadList contains a list of KeptnWorkload



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `lifecycle.keptn.sh/v1alpha3`
| `kind` _string_ | `KeptnWorkloadList`
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `items` _[KeptnWorkload](#keptnworkload) array_ |  |


#### KeptnWorkloadRef



KeptnWorkloadRef refers to a KeptnWorkload that is part of a KeptnApp

_Appears in:_
- [KeptnAppSpec](#keptnappspec)
- [KeptnAppVersionSpec](#keptnappversionspec)
- [WorkloadStatus](#workloadstatus)

| Field | Description |
| --- | --- |
| `name` _string_ | Name is the name of the KeptnWorkload. |
| `version` _string_ | Version is the version of the KeptnWorkload. |


#### KeptnWorkloadSpec



KeptnWorkloadSpec defines the desired state of KeptnWorkload

_Appears in:_
- [KeptnWorkload](#keptnworkload)
- [KeptnWorkloadInstanceSpec](#keptnworkloadinstancespec)

| Field | Description |
| --- | --- |
| `app` _string_ | AppName is the name of the KeptnApp containing the KeptnWorkload. |
| `version` _string_ | Version defines the version of the KeptnWorkload. |
| `preDeploymentTasks` _string array_ | PreDeploymentTasks is a list of all tasks to be performed during the pre-deployment phase of the KeptnWorkload. The items of this list refer to the names of KeptnTaskDefinitions located in the same namespace as the KeptnApp, or in the Keptn namespace. |
| `postDeploymentTasks` _string array_ | PostDeploymentTasks is a list of all tasks to be performed during the post-deployment phase of the KeptnWorkload. The items of this list refer to the names of KeptnTaskDefinitions located in the same namespace as the KeptnWorkload, or in the Keptn namespace. |
| `preDeploymentEvaluations` _string array_ | PreDeploymentEvaluations is a list of all evaluations to be performed during the pre-deployment phase of the KeptnWorkload. The items of this list refer to the names of KeptnEvaluationDefinitions located in the same namespace as the KeptnWorkload, or in the Keptn namespace. |
| `postDeploymentEvaluations` _string array_ | PostDeploymentEvaluations is a list of all evaluations to be performed during the post-deployment phase of the KeptnWorkload. The items of this list refer to the names of KeptnEvaluationDefinitions located in the same namespace as the KeptnWorkload, or in the Keptn namespace. |
| `resourceReference` _[ResourceReference](#resourcereference)_ | ResourceReference is a reference to the Kubernetes resource (Deployment, DaemonSet, StatefulSet or ReplicaSet) the KeptnWorkload is representing. |


#### KeptnWorkloadStatus



KeptnWorkloadStatus defines the observed state of KeptnWorkload

_Appears in:_
- [KeptnWorkload](#keptnworkload)

| Field | Description |
| --- | --- |
| `currentVersion` _string_ | CurrentVersion indicates the version that is currently deployed or being reconciled. |


#### Objective





_Appears in:_
- [KeptnEvaluationDefinitionSpec](#keptnevaluationdefinitionspec)

| Field | Description |
| --- | --- |
| `keptnMetricRef` _[KeptnMetricReference](#keptnmetricreference)_ | KeptnMetricRef references the KeptnMetric that should be evaluated. |
| `evaluationTarget` _string_ | EvaluationTarget specifies the target value for the references KeptnMetric. Needs to start with either '<' or '>', followed by the target value (e.g. '<10'). |


#### ResourceReference



ResourceReference represents the parent resource of Workload

_Appears in:_
- [KeptnWorkloadInstanceSpec](#keptnworkloadinstancespec)
- [KeptnWorkloadSpec](#keptnworkloadspec)

| Field | Description |
| --- | --- |
| `uid` _UID_ |  |
| `kind` _string_ |  |
| `name` _string_ |  |


#### RuntimeSpec





_Appears in:_
- [KeptnTaskDefinitionSpec](#keptntaskdefinitionspec)

| Field | Description |
| --- | --- |
| `functionRef` _[FunctionReference](#functionreference)_ | FunctionReference allows to reference another KeptnTaskDefinition which contains the source code of the function to be executes for KeptnTasks based on this KeptnTaskDefinition. This can be useful when you have multiple KeptnTaskDefinitions that should execute the same logic, but each with different parameters. |
| `inline` _[Inline](#inline)_ | Inline allows to specify the code that should be executed directly in the KeptnTaskDefinition, as a multi-line string. |
| `httpRef` _[HttpReference](#httpreference)_ | HttpReference allows to point to an HTTP URL containing the code of the function. |
| `configMapRef` _[ConfigMapReference](#configmapreference)_ | ConfigMapReference allows to reference a ConfigMap containing the code of the function. When referencing a ConfigMap, the code of the function must be available as a value of the 'code' key of the referenced ConfigMap. |
| `parameters` _[TaskParameters](#taskparameters)_ | Parameters contains parameters that will be passed to the job that executes the task as env variables. |
| `secureParameters` _[SecureParameters](#secureparameters)_ | SecureParameters contains secure parameters that will be passed to the job that executes the task. These will be stored and accessed as secrets in the cluster. |
| `cmdParameters` _string_ | CmdParameters contains parameters that will be passed to the command |


#### SecureParameters





_Appears in:_
- [KeptnTaskSpec](#keptntaskspec)
- [RuntimeSpec](#runtimespec)

| Field | Description |
| --- | --- |
| `secret` _string_ | Secret contains the parameters that will be made available to the job executing the KeptnTask via the 'SECRET_DATA' environment variable. The 'SECRET_DATA'  environment variable's content will the same as value of the 'SECRET_DATA' key of the referenced secret. |


#### ServiceAccountSpec





_Appears in:_
- [KeptnTaskDefinitionSpec](#keptntaskdefinitionspec)

| Field | Description |
| --- | --- |
| `name` _string_ |  |


#### TaskContext





_Appears in:_
- [KeptnTaskSpec](#keptntaskspec)

| Field | Description |
| --- | --- |
| `workloadName` _string_ | WorkloadName the name of the KeptnWorkload the KeptnTask is being executed for. |
| `appName` _string_ | AppName the name of the KeptnApp the KeptnTask is being executed for. |
| `appVersion` _string_ | AppVersion the version of the KeptnApp the KeptnTask is being executed for. |
| `workloadVersion` _string_ | WorkloadVersion the version of the KeptnWorkload the KeptnTask is being executed for. |
| `taskType` _string_ | TaskType indicates whether the KeptnTask is part of the pre- or postDeployment phase. |
| `objectType` _string_ | ObjectType indicates whether the KeptnTask is being executed for a KeptnApp or KeptnWorkload. |


#### TaskParameters





_Appears in:_
- [KeptnTaskSpec](#keptntaskspec)
- [RuntimeSpec](#runtimespec)

| Field | Description |
| --- | --- |
| `map` _object (keys:string, values:string)_ | Inline contains the parameters that will be made available to the job executing the KeptnTask via the 'DATA' environment variable. The 'DATA'  environment variable's content will be a json encoded string containing all properties of the map provided. |


#### WorkloadStatus





_Appears in:_
- [KeptnAppVersionStatus](#keptnappversionstatus)

| Field | Description |
| --- | --- |
| `workload` _[KeptnWorkloadRef](#keptnworkloadref)_ | Workload refers to a KeptnWorkload that is part of the KeptnAppVersion. |
| `status` _KeptnState_ | Status indicates the current status of the KeptnWorkload. |


