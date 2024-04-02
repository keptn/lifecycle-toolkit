# Migrating KeptnApp to KeptnAppContext

The
[KeptnAppContext](../../reference/crd-reference/appcontext.md)
resource was introduced in the `v1beta1` API version.
Existing
[KeptnApp](../../reference/crd-reference/app.md)
resources that were generated manually must be migrated
because of this new feature.
This page gives instructions for doing that.
Versions `v1beta1` and `v1` are fully compatible.

> **Note**
> Manual migration is only required for:
>
> * `KeptnApp` resources with older versions than `v1beta1`
> * Manually created `KeptnApp` resources
> * Automatically created `KeptnApp` resources with
>   manually edited pre/post-deployment tasks or evaluations**
>
> `KeptnApp` resources created using the
> [automatic app-discovery](../../guides/auto-app-discovery.md)
> feature without any manual changes in the pre/post-deployment
> tasks/evaluations section do not require any action.

## Migration steps

You need the following steps to migrate an existing `KeptnApp` resource
to work with the `KeptnAppContext` resource
introduced in the `v1beta1` API version:

1. Create a `KeptnAppContext` custom resource
   that has the same name as your `KeptnApp`.
2. Move the lists of pre/post-deployment tasks and evaluations
   from `KeptnApp` to `KeptnAppContext`.
   In other words, delete them from `KeptnApp.spec`
   and add them under the `KeptnAppContext.spec` field.
3. Add the `app.kubernetes.io/managed-by: keptn` annotation
   to the `KeptnApp` resource if it is not already there.

You can migrate your KeptnApp manually or, if you have
[go](https://go.dev/)
installed
on your machine, use the script provided
[here](https://github.com/keptn/lifecycle-toolkit/tree/main/lifecycle-operator/converter).

```bash
    go run convert_app.go path_to_keptnapp_to_convert path_to_desired_output_file
```

For instance, to run the example file conversion, the command is:

```bash
    go run convert_app.go example_keptnapp.yaml example_output.yaml
```

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

After applying the migration steps from the previous subsection,
you see the following changes:

* The `app.kubernetes.io/managed-by` annotation
  has been added to the `metadata.labels` section of the revised `KeptnApp` resource.

     If your original `KeptnApp` resource was auto-generated,
     it already had this annotation.

* The list of tasks and evaluations
has been moved from the `KeptnApp` resource
to the `KeptnAppContext` resource.

```yaml
{% include "./assets/keptnapp-migrated.yaml" %}
```

These modified resources can be now applied to your cluster.

## What's next?

Making these modifications does not alter the behavior of Keptn.
However, you might want to enhance your traces, tasks, and evaluations
with the new functionality that is available:

* Add context metadata to your traces.
  This allows you to include information
  like the stage into which the application is deployed, a commit ID,
  or other information relevant to the deployment traces of
  the application and its workloads.
  For instructions, see
  [Metadata](../../guides/metadata.md).
* Add `KEPTN_CONTEXT` information to the `function` code in your
  [KeptnTaskDefinition](../../reference/crd-reference/taskdefinition.md)
  resource.
  This allows you to correlate a task to a specific application/workload,
  provide information about the phase in which the task is executed,
  and access any metadata that has been attached to the application/workload
  such as commit ID or user name.
  For instructions, see
  [Context](../../guides/tasks.md#context).
  