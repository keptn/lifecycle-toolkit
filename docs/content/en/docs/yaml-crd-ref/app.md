---
title: KeptnApp
description: Define all workloads and checks associated with an application
weight: 10
---

A `KeptnApp` resource lists all the [workloads](https://kubernetes.io/docs/concepts/workloads/)
that constitute a logical application.
It contains information about:

- All [workloads](https://kubernetes.io/docs/concepts/workloads/) and checks
  that are associated with a Keptn application
- A list of tasks and evaluations to be executed
  pre- and post-deployment.
- Tasks referenced by `KeptnApp` are defined in a
  [KeptnTaskDefinition](taskdefinition.md) resource.
  `KeptnApp` identifies the task by the value of the `metadata.name` field
  and does not need to understand what runner is used to define the task.

## Synopsis

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha3
kind: KeptnApp
metadata:
  name: <app-name>
  namespace: <app-namespace>
spec:
  version: "x.y"
  revision: x
  workloads:
  - name: <workload1-name>
    version: <version-string>
  - name: <workload2-name>
    version: <version-string>
  preDeploymentTasks:
  - <list of tasks>
  postDeploymentTasks:
  - <list of tasks>
  preDeploymentEvaluations:
  - <list of evaluations>
  postDeploymentEvaluations:
  - <list of evaluations>
```

## Fields

The first set of fields are created automatically
when the app discovery feature generates the `KeptnApp` resource:

- **apiVersion** -- API version being used.
- **kind** -- Resource type.
   Must be set to `KeptnApp`

- **metadata**
  - **name** -- Unique name of this application.
    Names must comply with the
    [Kubernetes Object Names and IDs](https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-subdomain-names)
    specification.

- **spec**
  - **version** (required) -- version of the Keptn application.
    Changing this version number causes a new execution
    of all application-level checks
  - **revision** -- revision of a `version`.
    The value is an integer that can be modified
    to trigger another deployment of a `KeptnApp` of the same version.
    For example, increment this number to restart a `KeptnApp` version
    that failed to deploy, perhaps because a
    `preDeploymentEvaluation` or `preDeploymentTask` failed.
    See
    [Restart an Application Deployment](../implementing/restart-application-deployment.md)
    for a longer discussion of this.
  - **workloads**
    - **name** (required) -- name of this Kubernetes
      [workload](https://kubernetes.io/docs/concepts/workloads/).
      Use the same naming rules listed above for the application name.
      Provide one entry for each [workload](https://kubernetes.io/docs/concepts/workloads/)
      associated with this Keptn application.
    - **version** (required) -- version number for this [workload](https://kubernetes.io/docs/concepts/workloads/).
      Changing this number causes a new execution
      of checks for this [workload](https://kubernetes.io/docs/concepts/workloads/) only,
      not the entire application.

The remaining fields are required only when implementing
the release lifecycle management feature.
If used, these fields must be populated manually:

- **preDeploymentTasks** -- list each task
    to be run as part of the pre-deployment stage.
    Task names must match the value of the `metadata.name` field
    for the associated [KeptnTaskDefinition](taskdefinition.md) resource.
- **postDeploymentTasks** -- list each task
    to be run as part of the post-deployment stage.
    Task names must match the value of the `metadata.name` field
    for the associated
    [KeptnTaskDefinition](taskdefinition.md)
    resource.
- **preDeploymentEvaluations** -- list each evaluation to be run
    as part of the pre-deployment stage.
    Evaluation names must match the value of the `metadata.name` field
    for the associated
    [KeptnEvaluationDefinition](evaluationdefinition.md)
    resource.
- **postDeploymentEvaluations** -- list each evaluation to be run
    as part of the post-deployment stage.
    Evaluation names must match the value of the `metadata.name` field
    for the associated [KeptnEvaluationDefinition](evaluationdefinition.md)
    resource.

## Usage

Kubernetes defines
[workloads](https://kubernetes.io/docs/concepts/workloads/)
but does not define applications.
Keptn adds the concept of applications
defined as a set of workloads that can be executed.
A `KeptnApp` resource is added
into the repository of the deployment engine
(ArgoCD, Flux, etc.)
and is then deployed by that deployment engine.

A `KeptnApp` resource is created automatically, using the
[automatic application discovery](../implementing/integrate.md#use-keptn-automatic-app-discovery)
feature to generate a `KeptnApp` resource
based on the
[basic annotations](../implementing/integrate.md#basic-annotations)
that are applied to any of the workload resources.
This allows you to use the Keptn observability features for existing resources
without manually populating any Keptn related resources.

The release lifecycle management feature
allows you to define pre- and post-deployment
evaluations and tasks to be run for the `KeptnApp` as a whole.
These must be added to the `KeptnApp` manifest manually.
Note that all evaluations or tasks for a specific stage
(such as `preDeploymentTasks`)
are executed in parallel.
If you have a series of tasks that should be executed sequentially,
you can code them all into a single `KeptnTaskDefinition`.

## Example

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha3
kind: KeptnApp
metadata:
  name: podtato-head
  namespace: podtato-kubectl
spec:
  version: "latest"
  workloads:
  - name: podtato-head-left-arm
    version: "my_vers12.5"
  - name: podtato-head-left-leg
    version: "my_v24"
  postDeploymentTasks:
  - post-deployment-hello
  preDeploymentEvaluations:
  - my-prometheus-definition
```

## Files

## Differences between versions

- The `spec.revision` field is introduced in v1alpha2.

## See also

- [KeptnTaskDefinition](taskdefinition.md)
- [KeptnEvaluationDefinition](evaluationdefinition.md)
- [Working with tasks](../implementing/tasks.md)
- [Architecture of KeptnWorkloads and KeptnTasks](../architecture/keptn-apps.md)
- [Pre- and post-deployment tasks](../implementing/integrate.md#pre--and-post-deployment-checks)
- [Orchestrate deployment checks](../intro/usecase-orchestrate.md)
- [Use Keptn automatic app discovery](../implementing/integrate.md#use-keptn-automatic-app-discovery)
- [Restart an Application Deployment](../implementing/restart-application-deployment.md)
