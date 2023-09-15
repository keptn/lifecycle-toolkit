---
title: KeptnTask
description: Define a run of a KeptnTaskDefinition
weight: 85
---

Keptn populates the `KeptnTask` resource
for tasks that deploy software on Kubernetes.
When using Keptn to deploy software built on a virtual machine,
you must manually populate the `KeptnTask` definition
and modify it manually for each new run.

## Synopsis

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha3
kind: KeptnTask
metadata:
  name: <name-of-this-run>
spec:
  workload: "my-workload"
  workloadVersion: "1.0.0"
  appVersion: "1.0.0"
  app: "my-app"
  taskDefinition: helloworldtask
  context:
    appName: "my-app"
    appVersion: "1.0.0"
    objectType: ""
    taskType: ""
    workloadName: "my-workload"
    workloadVersion: "1.0.0"
```

## Fields

* **apiVersion** -- API version being used.
`
* **kind** -- Resource type.
   Must be set to `KeptnTask`

* **metadata**
  * **name** -- Unique name of this run of the task.
    This name must be modified each time you run this `KeptnTask`,
    so a common practice is to add a number to the end of the string
    used for the `name` field of the `KeptnTaskDefinition` resource
    so you can increment the number for each run.
    Names must comply with the
    [Kubernetes Object Names and IDs](https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-subdomain-names)
    specification.
* **spec** - Defines the desired state of this `KeptnTask` resource
  * **workload** - Name of the
      [KeptnWorkload](../crd-ref/lifecycle/v1alpha3/#keptnworkload)
      resource for which the `KeptnTask` is being executed.

      TODO: API ref shows `workloadName` and `workloadVersion`
      in `TaskContext`, not directly in `spec`.
      And I don't see a plain `workload` field.

  * **workloadVersion** - `KeptnWorkload` version
      for which the `KeptnTask` is being executed.
  * **appName** - Name of the
      [KeptnApp](../yaml-crd-ref/app.md) resource
      for which the `KeptnTask` is being executed.
  * **appVersion** - Version of the `KeptnApp` resource
      for which the `KeptnTask` is being executed.

      TODO: So I can create a `KeptnApp` resource that lists multiple workloads
      to be processed?
       Do I need to annotate the Workloads
      and can I autocreate the `KeptnApp` or do I need to do it manually?

      TODO: API reference shows `appName` and `appVersion` in `TaskContext`

      TODO: API reference shows `parameters` and `secureParameters` here.
      Should they be added to this list?

      TODO: API reference also shows `CheckType`, `retries`, and `timeout` here.

  * **taskDefinition** - Name of the corresponding `KeptnTaskDefinition` resource.
  This `KeptnTaskDefinition` can be located in the same namespace
  or in the Keptn name space.
    * **context** - Contextual information about the task execution
      * **appName** - Name of the
          [KeptnApp](../yaml-crd-ref/app.md) resource
          for which the `KeptnTask` is being executed.
      * **appVersion** - Version of the `KeptnApp` resource
          for which the `KeptnTask` is being executed.

          TODO: So I can create a `KeptnApp` resource that lists multiple workloads
          to be processed?
           Do I need to annotate the Workloads
          and can I autocreate the `KeptnApp` or do I need to do it manually?
      * **objectType** - Indicates whether this `KeptnTask`
          is being executed for a `KeptnApp` or a `KeptnWorkload` resource.

          TODO: So why is this set to null (`""`) in the example?
      * **taskType** Indicates whether this `KeptnTask`
          is part of the pre- or post-deployment phase

          TODO: So why is this set to null (`""`) in the example?

      * **workloadName** - Name of the `KeptnWorkload`
          for which the `KeptnTask` is being executed
      * **workloadVersion** - Version of the `KeptnWorkload`
          for which the `KeptnTask` is being executed

          TODO: Why are these fields both here and directly in `spec` above?

## Usage

Applying this file causes Keptn to create a Job and a Pod
and run the associated `KeptnTaskDefinition`.

Use the following commands to show the current status of the jobs:

```shell
kubectl get keptntasks
kubectl get pods
```

Each time you want to rerun the `KeptnTask` resource,
you must update the follow fields
to create a unique `KeptnTask` name
the `KeptnTask` name needs to be unique, update the follow fields:

* `name`
* `spec.appVersion`
* `spec.workloadVersion`
* `spec.context.appVersion`
* `spec.context.workloadVersion`

A common practice is to just increment the value of each field.

## Examples

For this example, the `KeptnTaskDefinition` resource is:

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha3
kind: KeptnTaskDefinition 
metadata:
  name: helloworldtask
spec:
  retries: 0
  timeout: 30s
  container:
    name: cowsay
    image: rancher/cowsay:latest
    args:
      - 'hello world'
```

And the `KeptnTask` resource for the first run is:

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha3
kind: KeptnTask
metadata:
  name: runhelloworld1
spec:
  workload: "my-workload"
  workloadVersion: "1.0.0"
  appVersion: "1.0.0"
  app: "my-app"
  taskDefinition: helloworldtask
  context:
    appName: "my-app"
    appVersion: "1.0.0"
    objectType: ""
    taskType: ""
    workloadName: "my-workload"
    workloadVersion: "1.0.0"
```

The `KeptnTask resource for the second run is:

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha3
kind: KeptnTask
metadata:
  name: runhelloworld2
spec:
  workload: "my-workload"
  workloadVersion: "1.0.1"
  appVersion: "1.0.1"
  app: "my-app"
  taskDefinition: helloworldtask
  context:
    appName: "my-app"
    appVersion: "1.0.1"
    objectType: ""
    taskType: ""
    workloadName: "my-workload"
    workloadVersion: "1.0.1"
```

## Files

[API reference](../crd-ref/lifecycle/v1alpha3/#keptntaskspec)

## Differences between versions

## See also

* [KeptnTaskDefinition](taskdefinition.md)
* [Keptn for Non-Kubernetes Applications](../implementing/tasks-non-k8s-apps.md)
