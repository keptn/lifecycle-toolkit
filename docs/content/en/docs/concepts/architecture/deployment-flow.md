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

Within each phase,
the order of execution is determined
by the order in which evaluations and tasks are listed in the
[KeptnApp](../../../yaml-crd-ref/app/)
resource.

A [Kubernetes Event](https://kubernetes.io/docs/reference/kubernetes-api/cluster-resources/event-v1/)
is emitted at each phase
to provide additional Observability of the execution flow.

The Keptn Lifecycle Toolkit implements a
[Permit Scheduler Plugin](https://kubernetes.io/docs/concepts/scheduling-eviction/scheduling-framework/#permit)
that blocks the binding of the pods to a node
until all the pre-conditions are fulfilled.

A Kubernetes deployment is started by the following command:

```shell
kubectl apply -f deployment.yaml
```

TODO: Should this mention anything
about deployment engines like Flux?
Maybe they just call ``kubectl apply` under the hood?

The deployment is created
but the created pods are blocked and in pending state
until all the required pre-deployment tasks/evaluations
defined on either the KeptnApp or KeptnWorkload level pass.
Only then are the pods bound to a node and deployed.
If any evaluation or task fails,
the `KeptnApp` issues appropriate `*Errored` event
and terminates the deployment.
If all evaluations and tasks in a phase are successful,
the `KeptnApp` issues the appropriate `*Succeeded` event
and initiates the next phase.

## Summary of deployment flow

To view these events on your cluster, execute:

```shell
kubectl get events -n <namespace> . 
```

### Pre-deployment phase

TODO: Add intro statement about this phase.

```shell
AppPreDeployTasks
  AppPreDeployTasksStarted
  AppPreDeployTasksSucceeded OR AppPreDeployTasksErrored
```

### Pre-deployment evaluation phase

TODO: Add intro statement about this phase.

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
the `KeptnApp` issues the `AppDeploySucceeded` event
and continues to the next phase.
If any of these activities fail,
the `KeptnApp` issues the `AppDeployErrored` event
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
  WorkloadPostDeployTasks
    WorkloadPostDeployTasksStarted
    WorkloadPostDeployTasksSucceeded OR WorkloadPostDeployTasksErrored 
  WorkloadPostDeployEvaluations
    WorkloadPostDeployEvaluationsStarted
    WorkloadPostDeployEvaluationsSucceeded OR WorkloadPostDeployEvaluationsErrored
  AppDeploySucceeded OR AppDeployErrored
  ```

### Post-deployment phase

TODO: Add intro statement about this phase.

```shell
AppPostDeployTasks
  AppPostDeployTasksStarted
  AppPostDeployTasksSucceeded OR AppPostDeployTasksErrored
```

### Post-deployment evaluation phase

TODO: Add intro statement about this phase.

```shell
AppPostDeployEvaluations
  AppPostDeployEvaluationsStarted
  AppPostDeployEvaluationsSucceeded OR AppPostDeployEvaluationsErrored
```

### Completed phase

```shell
Completed
```

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

