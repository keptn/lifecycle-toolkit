---
comments: true
---

# Flow of deployment

Keptn deploys a
[Kubernetes Workload](https://kubernetes.io/docs/concepts/workloads/)
by passing through a well-defined execution flow.

The execution flow goes through six main phases:

* Pre-deployments-tasks
* Pre-deployment-evaluation
* Deployment
* Post-deployment-tasks
* Post-deployment-evaluation
* Promotion
* Completed

Within each phase, all tasks and evaluations for each phase
are executed in parallel.
They are not affected by the order
in which evaluations and tasks are listed in the
[KeptnAppContext](../../reference/crd-reference/appcontext.md)
resource
or in the order of the pre/post-tasks, pre/post-evaluations, and promotion tasks
that are listed in the Workflow manifests.

## Kubernetes and Cloud Events

[Kubernetes Events](https://kubernetes.io/docs/reference/kubernetes-api/cluster-resources/event-v1/)
and [CloudEvents](https://cloudevents.io/)
are emitted at each phase
to provide additional Observability of the execution flow.

Keptn implements a
[Permit Scheduler Plugin](https://kubernetes.io/docs/concepts/scheduling-eviction/scheduling-framework/#permit)
that blocks the binding of the pods to a node
until all the pre-conditions are fulfilled.

A Kubernetes deployment is started by the deployment engine
that is implemented
(such as Flux or Argo)
or can be started by the following command:

```shell
kubectl apply -f deployment.yaml
```

Keptn does not care how a deployment manifest is applied to the cluster.
Both `kubectl` and Flux/Argo send the manifest to the Kubernetes API
so Keptn does not differentiate the actual deployment options.
This also means that one Keptn Application
can include services that are deployed with different methods.

The deployment is created
but the created pods are blocked and in pending state
until all the required pre-deployment tasks/evaluations
defined on either the `KeptnApp` or `KeptnWorkload` level pass.
Only then are the pods bound to a node and deployed.
If any pre-deployment evaluation or task fails,
the `KeptnApp` issues an appropriate `*Errored` event
and the deployment remains pending indefinitely,
until further changes or external intervention
either resolve the problem or terminate the execution.
If all evaluations and tasks in a phase are successful,
the `KeptnApp` issues the appropriate `*Succeeded` event
and initiates the next phase.

> **Note**
This behavior can be changed by configuring non-blocking deployment
functionality.
More information can be found in the
[Keptn non-blocking deployment section](./keptn-non-blocking.md).

## Summary of deployment flow

To view these events on your cluster, execute:

```shell
kubectl get events -n <namespace> . 
```

> **Note**
This only displays Kubernetes events, not Cloud Events.

### Pre-deployment phase

Pre-deployment tasks can perform any kind of action needed
to prepare for the deployment, including unit tests, load tests or other similar tests.

```shell
AppPreDeployTasks
  AppPreDeployTasksStarted
  AppPreDeployTasksSucceeded OR AppPreDeployTasksErrored
```

### Pre-deployment evaluation phase

Pre-deployment evaluation can be used to assert the status of the cluster
or of services the [workload](https://kubernetes.io/docs/concepts/workloads/) depends on,
to assure it is deployed only if the specified prerequisites are met.

```shell
AppPreDeployEvaluations
  AppPreDeployEvaluationsStarted
  AppPreDeployEvaluationsSucceeded OR AppPreDeployEvaluationsErrored
```

### Deployment phase

The `AppDeploy` phase basically covers
the entire deployment and check phase of the [workloads](https://kubernetes.io/docs/concepts/workloads/).
The `KeptnApp` just observes whether
all pre and post-deployment tasks/evaluation are successful
and that the pods are deployed successfully.
When all activities are successful,
the `KeptnApp` issues the `AppDeploySucceeded` event
and continues to the next phase.
If any of these activities fail,
the `KeptnApp` issues the `AppDeployErrored` event
and terminates the deployment.

> **Note**
By default Keptn observes the state of the Kubernetes workloads
for 5 minutes.
After this timeout is exceeded, the deployment phase (from Keptn
viewpoint) is considered as `Failed` and Keptn does not proceed
with post-deployment phases (tasks, evaluations or promotion phase).
This timeout can be modified for the cluster by changing the value
of the `observabilityTimeout` field in the
[KeptnConfig](../../reference/crd-reference/config.md)
resource.

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

The post-deployment phase is typically used
to run tests on the freshly deployed application,
such as gathering  performance data.

```shell
AppPostDeployTasks
  AppPostDeployTasksStarted
  AppPostDeployTasksSucceeded OR AppPostDeployTasksErrored
```

### Post-deployment evaluation phase

The post-deployment evaluation can be used
to analyze the cluster/application status after the new workload is deployed.
The result of this phase does not revert or influence the deployment
but can be used by other external tools, for instance, to react to a failure.

```shell
AppPostDeployEvaluations
  AppPostDeployEvaluationsStarted
  AppPostDeployEvaluationsSucceeded OR AppPostDeployEvaluationsErrored
```

### Promotion phase

The promotion phase is typically used
to run promotion tasks
(such as promoting the application to another stage)
on the freshly deployed application.

```shell
PromotionTasks
  PromotionTasksStarted
  PromotionTasksSucceeded OR PromotionTasksErrored
```

### Completed phase

```shell
Completed
```

## Events that are generated asynchronously

Additional phases/states exist,
such as those that describe what is currently happening in the system.
During the lifetime of the application, custom resources are created,
updated, deleted or reconciled.
Each reconciliation, or re-evaluation of the state of custom resources
by the controller, can cause the generation of events.
These include:

```shell
ReconcileEvaluation
ReconcileTask
ReconcileWorkload
CreateEvaluation
CreateTask
CreateApp
CreateAppVersion
CreateWorkload
CreateWorkloadVersion
CreateAppCreationRequest
UpdateWorkload
DeprecateAppVersion
AppCompleted
WorkloadCompleted
Deprecated
Completed
Cancelled
```
