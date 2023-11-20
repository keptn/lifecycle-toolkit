---
title: KeptnTask
description: Define a run of a KeptnTaskDefinition
weight: 85
---

Keptn uses KeptnTask resources internally
to manage tasks (and their underlying Kubernetes Job resources)
that are run before and after deployment of your workloads
(pre- and post-deployment tasks).
You do not need to create this resource yourself except in special situations,
like using Keptn to manage workloads that are outside of the k8s cluster.

## Synopsis

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha3
kind: KeptnTask
metadata:
  name: <name-of-this-run>
spec:
  taskDefinition: <name-of-KeptnTaskDefinition resource>
  context:
    appName: "<name-of-KeptnApp-resource>"
    appVersion: "1.0.0"
    objectType: ""
    taskType: ""
    workloadName: "my-workload"
    workloadVersion: "1.0.0"
  parameters: <parameters to pass to job>
  secureParameters: <secure parameters to pass to job>
  checkType: ""
  retries: <integer>
  timeout: <duration-in-seconds>
```

## Fields

* **apiVersion** -- API version being used.

* **kind** -- Resource type.
  Must be set to `KeptnTask`

* **metadata**
  * **name** -- Unique name of this run of the task.
    This name must be modified each time you run this `KeptnTask`,
    so a common practice is to add a number to the end of the string
    so you can increment the number for each run.
    Names must comply with the
    [Kubernetes Object Names and IDs](https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-subdomain-names)
    specification.
* **spec** - Defines the specification of this `KeptnTask` resource
  * **taskDefinition** (required) -- Name of the corresponding `KeptnTaskDefinition` resource.
    This `KeptnTaskDefinition` can be located in the same namespace
    or in the Keptn installation namespace.
  * **context** (required) -- Contextual information about the task execution
    * **appName** (required) -- Name of the
      [KeptnApp](../yaml-crd-ref/app.md) resource
      for which the `KeptnTask` is being executed.
    * **appVersion** (required) -- Version of the `KeptnApp` resource
      for which the `KeptnTask` is being executed.

    * **objectType** (required) -- Indicates whether this `KeptnTask`
      is being executed for a `KeptnApp` or a `KeptnWorkload` resource.
      When populating this resource manually
      to run a task for a non-Kubernetes deployment,
      set this value to `""`:
      Keptn populates this field based on annotations
      to the `KeptnWorkload` and `KeptnApp` resources.

    * **taskType** (required) -- Indicates whether this `KeptnTask`
      is part of the pre- or post-deployment phase.
      When populating this resource manually
      to run a task for a non-Kubernetes deployment,
      set this value to `""`:
      Keptn populates this field based on annotations
      to the `KeptnWorkload` and `KeptnApp` resources.

    * **workloadName** (required) -- Name of the `KeptnWorkload`
      for which the `KeptnTask` is being executed.
    * **workloadVersion** (required) -- Version of the `KeptnWorkload`
      for which the `KeptnTask` is being executed.
  * **parameters** -- Parameters that are passed to the job
    that executes the `KeptnTask`.
  * **secureParameters** -- Secure parameters that are passed
    to the job that executes the `KeptnTask`.
    These are stored and accessed as Kubernetes `Secrets` in the cluster.
    See [Working with secrets](../implementing/tasks.md#working-with-secrets)
    for more information.
  * **checkType** -- Defines whether task is part of pre- or post-deployment phase.
    Keptn populates this field based on annotations
    to the `KeptnWorkload` and `KeptnApp` resources.
  * **retries** -- If errors occur,
    this defines the number of attempts made
    before the `KeptnTask` is considered to be failed.
  * **timeout** -- Specifies the time, in seconds,
    to wait for the `KeptnTask` to complete successfully.
    If the `KeptnTask` does not complete successfully in this timeframe,
    it is considered to be failed.

## Usage

Applying this file causes Keptn to create a Job and a Pod
and run the associated `KeptnTaskDefinition`.

Use the following commands to show the current status of the jobs:

```shell
kubectl get keptntasks
kubectl get pods
```

Each time you want to rerun the `KeptnTask` resource,
you must update the value of the `metadata.name` field.
A common practice is to just increment the value incrementally,
so `helloworldtask-1` becomes `helloworldtask-2`, etc.

For a full example of how to create a `KeptnTask` resource
to use for a deployment being done outside of Kubernetes, see
[Keptn for Non-Kubernetes Applications](../implementing/tasks-non-k8s-apps.md).

## Files

[API reference](../crd-ref/lifecycle/v1alpha3/#keptntaskspec)

## Differences between versions

The syntax of the `KeptnTask` resource changed significantly
in Keptn v0.8.0.

## See also

* [KeptnTaskDefinition](taskdefinition.md)
* [Keptn for Non-Kubernetes Applications](../implementing/tasks-non-k8s-apps.md)
