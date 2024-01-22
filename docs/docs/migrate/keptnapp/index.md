# Migrating KeptnApp to KeptnAppContext

Keptn has recently orphaned the concept of manually created KeptnApp
and instead uses the KeptnApp auto-discovery feature combined with
the newly introduced KeptnAppContext custom resource.

## Migration steps

To migrate an existing **manually created KeptnApp** to KeptnAppContext
you need to execute the following steps:

1. Create a KeptnAppContext custom resource with the same name
as your KeptnApp has
2. Move pre/post-deployment tasks and evaluation from KeptnApp
to KeptnAppContext
3. Introduce a `app.kubernetes.io/managed-by: keptn` annotation
to KeptnApp
4. Change the API version of KeptnApp to `lifecycle.keptn.sh/v1beta1`

**Warning:**
Please note that the migration is needed to be executed only for
**manually created KeptnApp** resources. KeptnApp resources created by
the automatic app-discovery feature do not need any actions.

**Note:**
Please make sure all of your application resources
(such as Pods, Deployments, StatefulSets or DaemonSets)
have the proper annotations/labels set.
These annotations/labels (especially the
`app.kubernetes.io/part-of` or `keptn.sh/app`)
are needed for the full migration to the
usage of automatic app-discovery feature.
More information about how to set up these annotations/labels
can be found [here](../../guides/integrate.md#basic-annotations).

## Example of migration

In the next subsection we are going to look at an example of the migration.
Let's say we have the following KeptnApp in our cluster:

```yaml
{% include "./asserts/keptnapp.yaml" %}
```

If we apply the migration steps from the previous subsection, we will ge the
following result:

```yaml
{% include "./asserts/keptnapp-migrated.yaml" %}
```

These modified resources can be now applied to your cluster.

**Note:**
Please make sure all of your application resources
(such as Pods, Deployments, StatefulSets or DaemonSets)
have the proper annotations/labels set.
More information about how to set up these annotations/labels
can be found [here](../../guides/integrate.md#basic-annotations).
