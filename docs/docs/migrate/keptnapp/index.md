# Migrating KeptnApp to KeptnAppContext

The introduction of the `KeptnAppContext` resource
in Keptn v0.11.0
requires modifications to existing
[KeptnApp](../../reference/crd-reference/app.md)
resources that were generated manually.
This page gives instructions for doing that.

## Migration steps

> **Warning**
Migration is only required for
**manually created KeptnApp** resources.
`KeptnApp` resources created by
Keptn via the
[automatic app-discovery](../../guides/auto-app-discovery.md)
feature do not require modification
unless you edited them manually to add pre/post-deployment tasks
or evaluations.

To migrate an existing **manually created KeptnApp** to KeptnAppContext
you need to execute the following steps:

1. Create a `KeptnAppContext` custom resource that has the same name
as your `KeptnApp`.
2. Move the lists of pre/post-deployment tasks and evaluation from `KeptnApp`
to `KeptnAppContext`.
In other words, delete them from `KeptnApp.spec` and add them under the `KeptnAppContext.spec` field.
3. Add the `app.kubernetes.io/managed-by: keptn` annotation
to `KeptnApp`.

> **Note**
Please make sure all of your application resources
(such as Pods, Deployments, StatefulSets or DaemonSets)
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

Applying the migration steps from the previous subsection, we get the
following result:

```yaml
{% include "./assets/keptnapp-migrated.yaml" %}
```

These modified resources can be now applied to your cluster.
