---
title: Flow of deployment
linktitle: Flow of deployment
description: Understand the execution flow of a deployment
weight: 25
---

The Keptn Lifecycle Toolkit (KLT) deploys a
[Kubernetes Workload](https://kubernetes.io/docs/concepts/workloads/)
by passing through a well-defined execution flow.

The execution flow goes through six main phases:

* Pre-deployments-tasks
* Pre-deployment-evaluation
* Deployment
* Post-deployment-tasks
* Post-deployment-evaluation
* Completed

A [Kubernetes Event](https://kubernetes.io/docs/reference/kubernetes-api/cluster-resources/event-v1/)
is emitted at each phase to provide additional Observability of the execution flow.

The Keptn Lifecycle Toolkit implements a
[Permit Scheduler Plugin](https://kubernetes.io/docs/concepts/scheduling-eviction/scheduling-framework/#permit)
that blocks the creation of the pods until all the pre-conditions are fulfilled.

A Kubernetes deployment is started by the following command:

```shell
kubectl apply -f deployment.yaml
```

The `kubectl apply` occurs at the beginning of the deployment
but the created pods are blocked and in pending state
until all the required pre-deployment tasks/evaluations
defined on either the KeptnApp or KeptnWorkload level pass.
Only then are the pods bound to a node and deployed.

## Summary of deployment flow

To view these events on your cluster, execute:

```shell
kubectl get events -n <namespace> . 
```

### Pre-deployment phase

```shell
AppPreDeployTasks
  AppPreDeployTasksStarted
  AppPreDeployTasksSucceeded OR AppPreDeployTasksErrored
```

### Pre-deployment evaluation phase

```shell
AppPreDeployEvaluations
  AppPreDeployEvaluationsStarted
  AppPreDeployEvaluationsSucceeded OR AppPreDeployEvaluationsErrored
```

### Deployment phase

The `AppDeploy` phase basically covers
the entire deployment and check phase of the workloads.
The `KeptnApp` just observes whether
all pre and post-deployment tasks/evaluation are successful
and that the pods are deployed successfully.
When all activities are successful,
the KeptnApp issues the `AppPostDeployEvaluationsSucceeded` event
and continues to the next phase.
If any of these activities fail,
the KeptnApp issues the `AppPostDeployEvaluationsErrored` event
and terminates the deployment.

```shell
AppDeploy
  AppDeployStarted
  WorkloadPreDeployTasks
    WorkloadPreDeployTasksStarted
    WorkloadPreDeployTasksSucceeded OR WorkloadPreDeployTasksErrored 
  WorkloadPreDeployEvaluations
    WorkloadPreDeployEvaluationsStarted
    WorkloadPreDeployEvaluationsSucceeded OR WorkloadPreDeployErrored
  WorkloadDeploy
    WorkloadDeployStarted
    WorkloadDeploySucceeded OR WorkloadDeployErrored
```

### Post-deployment phase

```shell
  WorkloadPostDeployTasks
    WorkloadPostDeployTasksStarted
    WorkloadPostDeployTasksSucceeded OR WorkloadPostDeployTasksErrored
  WorkloadPostDeployEvaluations
    WorkloadPostDeployEvaluationsStarted
    WorkloadPostDeployEvaluationsSucceeded OR WorkloadPostDeployEvaluationsErrored
AppDeploy
  AppDeploySucceeded OR AppDeployErrored
  
### Post-deployment evaluation phase

```shell
AppPostDeployTasks
  AppPostDeployTasksStarted
  AppPostDeployTasksSucceeded OR AppPostDeployTasksErrored
```

```shell
AppPostDeployEvaluations
  AppPostDeployEvaluationsStarted
  AppPostDeployEvaluationsSucceeded OR AppPostDeployEvaluationsErrored
```

### Completed phase

## Events that are not part of the deployment flow

Additional phases/states exist,
such as those that describe what happens when something fails.

Whenever something in the system happens (we create a new resource, etc.)
a Kubernetes event is generated.
The following events are defined as part of the Keptn Lifecycle Toolkit
but they are not part of the deployment flow.
These include:

```shell
CreateEvaluation
ReconcileEvaluation
ReconcileTask
CreateTask
CreateApp
CreateAppVersion
CreateWorkload
CreateWorkloadInstance
Completed
Deprecated
WorkloadDeployReconcile
  WorkloadDeployReconcileErrored
```
