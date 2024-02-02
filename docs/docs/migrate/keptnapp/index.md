# Migrating KeptnApp to KeptnAppContext

The
[KeptnAppContext](../../reference/crd-reference/appcontext.md)
resource is new in the `v1beta1` API included in Keptn v0.11.0.
Existing
[KeptnApp](../../reference/crd-reference/app.md)
resources that were generated manually must be migrated
because of this new feature..
This page gives instructions for doing that.

> **Note**
> Manual migration is only required for:
>
> * Manually created `KeptnApp` resources
> * Automatically created `KeptnApp` resources with
>   manually edited pre/post-deployment tasks or evaluations**
>
> `KeptnApp` resources created using the
> [automatic app-discovery](../../guides/auto-app-discovery.md)
> feature without any manual changes in the pre/post-deployment
> tasks/evaluations section do not require any action.

## Migration steps

The steps to migrate an existing `KeptnApp` resource
to work with the `KeptnAppContext` resource introduced in version 0.11.0,
you need to execute the following steps:

1. Create a `KeptnAppContext` custom resource
   that has the same name as your `KeptnApp`.
2. Move the lists of pre/post-deployment tasks and evaluations
   from `KeptnApp` to `KeptnAppContext`.
   In other words, delete them from `KeptnApp.spec`
   and add them under the `KeptnAppContext.spec` field.
3. If necessary, add the `app.kubernetes.io/managed-by: keptn` annotation
   to the `KeptnApp` resource.

> **Note**
Be sure that all of your application resources
(such as
[Pods](https://kubernetes.io/docs/concepts/workloads/pods/),
[Deployments](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/),
[StatefulSets](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/),
and
[DaemonSets](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/)
have the proper annotations/labels set.
These annotations/labels (especially the
`app.kubernetes.io/part-of` or `keptn.sh/app`)
are necessary for the migration to the
automatic app-discovery feature.
More information about how to set up these annotations/labels
can be found [here](../../guides/integrate.md#basic-annotations).

## Example of migration

Here, we provide an example of how to
migrate the `KeptnApp` definition to the `KeptnAppContext`.
Let's say we have the following `KeptnApp` in our cluster:

```yaml
{% include "./assets/keptnapp.yaml" %}
```

Applying the migration steps from the previous subsection,
we get the following result.
You see the following changes:

* The `app.kubernetes.io/managed-by` annotation
  has been added to the `labels` section of the revised `KeptnApp` resource.

     If your original `KeptnApp` resource was auto-generated,
     it already had this annotation.

* The list of tasks and evaluations
has been moved from the `KeptnApp` resource
to the `KeptnAppContext` resource

```yaml
{% include "./assets/keptnapp-migrated.yaml" %}
```

These modified resources can be now applied to your cluster.

## What's next?

Making these modifications allows your existing `KeptnApp` functionality
run as it did before.
However, new capabilities can be added:

* Add context metadata to your traces.
  This allows you to include information like `stage` in your workload traces
  and information such as commit ID or user name
  in your `KeptnApp` traces.
  For instructions, see
  [Metadata](../../guides/metadata.md).
* Add `KEPTN_CONTEXT` information to the `function` code in your
  [KeptnTaskDefinition](../../reference/crd-reference/taskdefinition.md)
  resource.
  This allows you to correlate a task to a specific application/workload,
  provide informaton about the phase in which the task is executed,
  and access any metadata that has been attached to the application/workload
  such as commit ID or user name.
  For instructions, see
  [Context](../../guides/tasks#context).
  