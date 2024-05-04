---
comments: true
---

# KeptnAppContext

The `KeptnAppContext` custom resource works with the corresponding
[KeptnApp](app.md) resource
that has the same name and is located in the same namespace.
It allows you to

- Add metadata and links to traces for a specific application.
  This enables you to enrich your traces with additional information that
  you can use to better understand and analyze
  the performance of your applications.

- Define tasks and evaluations that run
  before or after the deployment of the `KeptnApp`.

## Synopsis

```yaml
apiVersion: lifecycle.keptn.sh/v1
kind: KeptnAppContext
metadata:
  name: <app-name>
  namespace: <app-namespace>
spec:
  metadata:
    <custom-attributes>
  spanLinks:
    - "<list of links>"
  preDeploymentTasks:
    - <list of tasks>
  postDeploymentTasks:
    - <list of tasks>
  preDeploymentEvaluations:
    - <list of evaluations>
  postDeploymentEvaluations:
    - <list of evaluations>
  promotionTasks:
    - <list of tasks>
```

## Fields

- **apiVersion** -- API version being used.
  Must be set to `lifecycle.keptn.sh/v1`
- **kind** -- Resource type.
  Must be set to `KeptnAppContext`

- **metadata**
    - **name** -- Unique name of this `KeptnAppContext` resource.
      Names must comply with the
      [Kubernetes Object Names and IDs](https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-subdomain-names)
      specification
      and match the `name` given to the associated `KeptnApp` resource.
    - **namespace** -- Namespace of this application.
      This must match the `namespace` given to the associated `KeptnApp` resource.
- **spec**
    - **metadata** -- list of key-value pairs
      that are propagated to the application trace as attributes.
      For example, the following lines add the `commit-id` and `author` information to the `KEPTN_CONTEXT`
      of the workload or `KeptnApp` where it is specified:

        ```yaml
        spec:
          metadata:
            commit-id: "1234"
            author: "myUser"
        ```

        For more information, see [Context metadata](../../guides/metadata.md).

    - **spanLinks** -- List of OpenTelemetry span links
      that connect multiple traces.
      For example, this can be used to connect deployments of the same application through different stages.
      You can retrieve the value to use from the JSON representation of the trace in Jaeger.
      The structure of this is:

        ```yaml
        00-<trace-id>-<span-id>-01
        ```

        After you add this field to your `KeptnAppContext` manifest, you must increment the `revision` number
        of the corresponding `KeptnApp` resource and apply the manifest to store the information in the traces.

        For more information, see
        [Advanced tracing configurations in Keptn Linking traces](../../guides/otel.md#advanced-tracing-configurations-in-keptn-linking-traces)

The remaining fields are required only when implementing the release lifecycle management feature.
If used, these fields must be populated manually:

- **spec**

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
    - **promotionTasks** -- list each task
      to be run as part of the promotion stage.
      Task names must match the value of the `metadata.name` field
      for the associated [KeptnTaskDefinition](taskdefinition.md) resource.

## Usage

`KeptnAppContext` lists the tasks and evaluations to be executed pre/post-deployment.
Tasks referenced by `KeptnAppContext` are defined in a [KeptnTaskDefinition](taskdefinition.md) resource.
`KeptnAppContext` identifies each task by the value of the `metadata.name` field
and does not need to understand what runner is used to define the task.
Similarly, evaluations referenced are defined in a [KeptnEvaluationDefinition](evaluationdefinition.md)
resource and identified by the value of the `metadata.name` field;
`KeptnAppContext` does not need to understand the data source or query being used for the evaluation.

## Example

```yaml
apiVersion: lifecycle.keptn.sh/v1
kind: KeptnAppContext
metadata:
  name: podtato-head
  namespace: podtato-kubectl
spec:
  preDeploymentTasks:
    - container-sleep
    - python-secret
```

## Files

[KeptnAppContext](../api-reference/lifecycle/v1/index.md#keptnappcontext)

## Differences between versions

The `KeptnAppContext` resource is new in the `v1beta1` version of the lifecycle operator.
Versions `v1beta1` and `v1` are fully compatible.

## See also

- [KeptnApp](app.md)
- [KeptnTaskDefinition](taskdefinition.md)
- [KeptnEvaluationDefinition](evaluationdefinition.md)
- [Deployment tasks](../../guides/tasks.md)
- [Architecture of KeptnWorkloads and KeptnTasks](../../components/lifecycle-operator/keptn-apps.md)
- Getting started with
  [Release Lifecycle Management](../../getting-started/lifecycle-management.md)
- [Use Keptn automatic app discovery](../../guides/auto-app-discovery.md)
- [Restart an Application Deployment](../../guides/restart-application-deployment.md)
- [Context metadata](../../guides/metadata.md)
- [Advanced tracing configurations in Keptn Linking traces](../../guides/otel.md#advanced-tracing-configurations-in-keptn-linking-traces)
