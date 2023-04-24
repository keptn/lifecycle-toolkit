---
title: Flow of deployment
linktitle: Flow of deployment
description: Understand the execution flow of a deployment
weight: 25
---

A deployment is started by the following command:

```bash
kubectl apply -f deployment.yaml
```

This executes a series of
Kubernetes [events](https://kubernetes.io/docs/reference/kubernetes-api/cluster-resources/event-v1/).
Events are resources such as Deployments or Pods.

A user can view these events by executing:

```bash
kubectl get events -n <namespace> . 
```

The `kubectl apply` occurs at beginning of the deployment
but the created pods are blocked and in pending state
until all the required pre-deployment tasks/evaluation
defined on either the KeptnApp or Workload level pass.
Only then are the pods bound to a node and deployed.

The Keptn Lifecycle Toolkit implements a [Permit Scheduler Plugin](https://kubernetes.io/docs/concepts/scheduling-eviction/scheduling-framework/#permit)
that blocks the creation of the pods until all the pre-conditions are fulfilled.

## Summary of deployment flow

```bash
AppPreDeployTasks
  AppPreDeployTasksStarted
  AppPreDeployTasksSucceeded OR AppPreDeployTasksErrored
AppPreDeployEvaluations
  AppPreDeployEvaluationsStarted
  AppPreDeployEvaluationsSucceeded OR AppPreDeployEvaluationsErrored
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
AppDeploy
  AppDeploySucceeded OR AppDeployErrored
AppPostDeployTasks
  AppPostDeployTasksStarted
  AppPostDeployTasksSucceeded OR AppPostDeployTasksErrored
AppPostDeployEvaluations
  AppPostDeployEvaluationsStarted
  AppPostDeployEvaluationsSucceeded OR AppPostDeployEvaluationsErrored
```

## Pre-deployment phase

## Pre-deployment evaluation phase

## Deployment phase

The `AppDeploy` phase basically covers
the entire deployment and check phase of the workloads.
The KeptnApp just observes whether
all pre and post-deployment tasks/evaluation are successful
and that the pods are deployed successfully.
When all activities are successful,
the KeptnApp issues the `AppPostDeployEvaluationsSucceeded` event
and continues to the next phase.
If any of these activities fail,
the KeptnApp issues the `AppPostDeployEvaluationsErrored` event
and terminates the deployment.

## Post-deployment phase

## Post-deployment evaluation phase

## Completed phase

Additional phases/states exist,
such as those that describe what happens when something fails.

## Events that are not part of the deployment flow

Whenever something in the system happens (we create a new resource, etc.)
we genereta a Kubernetes event.
The following events are defined as part of the Keptn Lifecycle Toolkit
but they are not part of the deployment flow.
These include:

```bash
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

## Other phases

## More source (to be written up)

6 main phases (pre-deployments-tasks, pre-deployment-evaluation, deployment,
post-deployments-tasks, post-deployment-evaluation, completed)
