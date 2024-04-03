---
comments: true
---

# KeptnTask

Keptn uses `KeptnTask` resources internally
to manage tasks (and their underlying Kubernetes Job resources)
that are run before and after deployment of your workloads  
(pre- and post-deployment tasks).
You do not need to create this resource yourself except in special situations,
for instance when using Keptn to manage workloads outside the k8s cluster.
See the [use case page](../../use-cases/non-k8s.md) for more details about this.

## Synopsis

```yaml
apiVersion: lifecycle.keptn.sh/v1
kind: KeptnTask
metadata:
  name: <name-of-this-run>
spec:
  taskDefinition: <name-of-KeptnTaskDefinition resource>
  context:
    appName: "<name-of-KeptnApp-resource>"
    appVersion: "KeptnApp-version"
    objectType: ""
    taskType: ""
    workloadName: "name-of-KeptnWorkload resource""
    workloadVersion: "version-of-KeptnWorkload resource"
    metadata:
      <custom-info1>: "<custom-info1-value>"
      <custom-info2>: "<custom-info2-value>"
      ...
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
          [KeptnApp](app.md) resource
          for which the `KeptnTask` is being executed.
        * **appVersion** (required) -- Version of the `KeptnApp` resource
          for which the `KeptnTask` is being executed.
        * **metadata** -- Additional key-value pairs with contextual information for the `KeptnTask`.
          Keptn populates this field based on the `spec.metadata` field of the
          [KeptnWorkloadVersion](../api-reference/lifecycle/v1/index.md#keptnworkloadversion)
          and
          [KeptnAppVersion](../api-reference/lifecycle/v1/index.md#keptnappversion)
          resources.

            For example, the following lines reference the commit ID and user ID:

            ```yaml
            commit-id: "1234"
            user-id: "person3"
            ```

            See [Context metadata](../../guides/metadata.md)
            for information about setting user defined values for those fields.

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
      See [Working with secrets](../../guides/tasks.md#working-with-secrets)
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

## Example

For a full example of how to create a `KeptnTask` resource
to use for a deployment being done outside of Kubernetes, see
[Keptn for Non-Kubernetes Applications](../../use-cases/non-k8s.md).

## Files

[API reference](../api-reference/lifecycle/v1/index.md#keptntaskspec)

## Differences between versions

The syntax of the `KeptnTask` resource changed significantly
in Keptn v0.8.0.

## See also

* [KeptnTaskDefinition](taskdefinition.md)
* [Keptn for Non-Kubernetes Applications](../../use-cases/non-k8s.md)
