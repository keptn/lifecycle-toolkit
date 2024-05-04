---
comments: true
---

# Day 2 Operations

After you have successfully rolled out your application by following
the instructions in the [integration guide](../guides/integrate.md),
Keptn also assists you with day 2 operations for your application.

Tasks that fall under this category include:

* Updating the version of one or more [workloads](https://kubernetes.io/docs/concepts/workloads/)
  that are part of the same application
* Adding a new [workload](https://kubernetes.io/docs/concepts/workloads/) to an existing application
* Monitoring the health of your application using `KeptnMetrics`, as described [here](../guides/evaluatemetrics.md)
* Optimizing the resource usage of your applications by integrating
  `KeptnMetrics` into a
  [HorizontalPodAutoscaler (HPA)](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/),
  as described [here](./hpa.md)

## Updating Workload Versions

After a first successful deployment of your application with Keptn,
you will eventually need to update the version of one or
more [workloads](https://kubernetes.io/docs/concepts/workloads/) that are part of the application.
This usually involves updating the image of a deployment
and changing the configuration of a deployment.
For example, using a different service account name for a pod
could be seen as an update.
Regardless of that, however, it is the user who decides what Keptn
sees as a version bump in the application by setting the value of
the `keptn.sh/version` or `app.kubernetes.io/version` labels/annotations
in their [workloads](https://kubernetes.io/docs/concepts/workloads/).

When this changes, Keptn interprets a change as a new version
and thus re-runs the pre- and post-tasks and evaluations for the application
as well as promotion tasks if the previous phases succeeded.

If the version label/annotation does not change, Keptn does not consider
a change of a [workload](https://kubernetes.io/docs/concepts/workloads/) configuration to be an update,
and therefore no pre- and post-tasks/evaluations are executed because they have already been
completed for the version set in the labels/annotations.

To illustrate the update of a workload,
let's assume the following example, including
a [workload](https://kubernetes.io/docs/concepts/workloads/) called `podtato-head-frontend` that includes a
pre-deployment task.

```yaml
{% include "./assets/day-2-operations/deployment-initial.yaml" %}
```

Now, let's assume that the configuration of that workload
needs to be changed.
In this example we assume that the image of that workload
should be updated, but a configuration change is not limited to that.
From here, you essentially have two options:

* **Only update the configuration *without* updating the `app.kubernetes.io/version`
  label:** This can be useful when the change in the configuration should happen regardless
  of the result of any task or evaluation, e.g., when the previously used image has a critical vulnerability
  and the image must be updated as quickly as possible.
  To do that, change `podtato-head-frontend` as follows:

    ```yaml
    {% include "./assets/day-2-operations/deployment-new-image.yaml" %}
    ```

* **Update the configuration *and* the version label:**
  Doing so causes the `KeptnWorkload` that is associated
  with the `podtato-head-frontend` deployment to be updated,
  and therefore the pre-task `my-task` and pre-evaluation `my-evaluation`
  are executed before the updated pods are scheduled.
  In this case, the deployment should be changed as follows:

    ```yaml
    {% include "./assets/day-2-operations/deployment-new-image-and-version.yaml" %}
    ```

Applying this causes the
[KeptnApp](../reference/crd-reference/app.md)
resource to be updated with a new
version, and a new
[KeptnAppVersion](../reference/api-reference/lifecycle/v1/index.md#keptnappversion)
resource to be created.
Due to this, all checks defined in the
[KeptnAppContext](../reference/api-reference/lifecycle/v1/index.md#keptnappcontext)
resource
as well as those defined in the deployment's `keptn.sh/pre-deployment-tasks`
label are executed again.

After applying the updated manifests, you can monitor the status
of the application and related workloads using the following commands:

```shell
$ kubectl get keptnworkloadversion -n podtato-kubectl

NAMESPACE   NAME                                             APPNAME         WORKLOADNAME                         WORKLOADVERSION      PHASE
podtato-kubectl   podtato-head-podtato-head-frontend-0.3.0   podtato-head    podtato-head-podtato-head-frontend   0.3.0                Completed
podtato-kubectl   podtato-head-podtato-head-hat-0.3.0        podtato-head    podtato-head-podtato-head-hat        0.3.0                Completed
podtato-kubectl   podtato-head-podtato-head-frontend-0.3.1   podtato-head    podtato-head-podtato-head-frontend   0.3.1                Completed
```

As can be seen in the output of the command, the
[KeptnWorkloadVersion](../reference/api-reference/lifecycle/v1/index.md#keptnworkloadversion)
resources from the previous deployment
are still here, but a new `KeptnWorkloadVersion` for the updated workload
has been added.
For the workload that
remained unchanged (`podtato-head-hat`), no new `KeptnWorkloadVersion` needed to be created.

Similarly, retrieving the list of `KeptnAppVersions` will reflect the update by
returning a newly created `KeptnAppVersion`.

```shell
$ kubectl get keptnappversion -n podtato-kubectl

NAMESPACE         NAME                               APPNAME        VERSION   PHASE
podtato-kubectl   podtato-head-f13dcb00ea-6b86b273   podtato-head   0.1.0     Completed
podtato-kubectl   podtato-head-1c40c739cf-d4735e3a   podtato-head   0.1.0     Completed
```

## Adding a new Workload to an Application

To add a new workload (e.g. a new deployment) to an existing app,
you must:

* Make sure the
  `keptn.sh/app`/`app.kubernetes.io/part-of` label/annotation is present
  on the new [workload](https://kubernetes.io/docs/concepts/workloads/)
* Add the new [workload](https://kubernetes.io/docs/concepts/workloads/) to the `KeptnApp`,
  if you have previously defined the `KeptnApp` resource manually.
  If the application has been discovered automatically, this step is not needed.

For example, to add the deployment `podtato-head-left-leg` to the
`podtato-head` application, the configuration for that new deployment
would look like this, with the required label being set:

```yaml
{% include "./assets/day-2-operations/new-deployment.yaml" %}
```

After applying the updated manifests, you can monitor the status
of the application and related workloads using the following commands:

```shell
$ kubectl get keptnworkloadversion -n podtato-kubectl

NAMESPACE   NAME                                             APPNAME         WORKLOADNAME                         WORKLOADVERSION      PHASE
podtato-kubectl   podtato-head-podtato-head-frontend-0.3.0   podtato-head    podtato-head-podtato-head-frontend   0.3.0                Completed
podtato-kubectl   podtato-head-podtato-head-frontend-0.3.1   podtato-head    podtato-head-podtato-head-frontend   0.3.1                Completed
podtato-kubectl   podtato-head-podtato-head-hat-0.3.0        podtato-head    podtato-head-podtato-head-hat        0.3.0                Completed
podtato-kubectl   podtato-head-podtato-head-left-leg-0.3.0   podtato-head    podtato-head-podtato-head-left-leg   0.3.0                Completed
```

As can be seen in the output of the command, in addition
to the previous `KeptnWorkloadVersions`, the newly created
`KeptnWorkloadVersion`, `podtato-head-podtato-head-left-leg-0.3.0` has been added
to the results.
