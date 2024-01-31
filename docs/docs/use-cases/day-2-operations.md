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
and thus re-runs the pre- and post-tasks and evaluations for the application.

If the version label/annotation does not change, Keptn does not consider
a change of a [workload](https://kubernetes.io/docs/concepts/workloads/) configuration to be an update,
and therefore no pre- and post-tasks/evaluations are executed because they have already been
completed for the version set in the labels/annotations.

To illustrate the update of a [workload](https://kubernetes.io/docs/concepts/workloads/),
let's assume the following example, including
a [workload](https://kubernetes.io/docs/concepts/workloads/) called `podtato-head-frontend` that includes a
pre-deployment task.

```yaml
{% include "./assets/deployment-initial.yaml" %}
```

Now, let's assume that the configuration of that [workload](https://kubernetes.io/docs/concepts/workloads/) needs to be changed.
In this example we assume that the image of that [workload](https://kubernetes.io/docs/concepts/workloads/)
should be updated, but a configuration change is not limited to that.
From here, you essentially have two options:

* **Only update the configuration *without* updating the `app.kubernetes.io/version`
label:** This can be useful when the change in the configuration should happen regardless
of the result of any task or evaluation, e.g., when the previously used image has a critical vulnerability
and the image must be updated as quickly as possible.
To do that, change `podtato-head-frontend` as follows:

```yaml
{% include "./assets/deployment-new-image.yaml" %}
```

* **Update the configuration *and* the version label:**
Doing so causes the `KeptnWorkload` that is associated
with the `podtato-head-frontend` deployment to be updated,
and therefore the pre-task `my-task` and pre-evaluation `my-evaluation`
are executed before the updated pods are scheduled.
In this case, the deployment should be changed as follows:

```yaml
{% include "./assets/deployment-new-image-and-version.yaml" %}
```



             OBSOLETE If you have defined the related `KeptnApp` resource yourself,
             OBSOLETE this must also be updated to refer to the updated `KeptnWorkload`.
             OBSOLETE This is a mandatory step, since the `KeptnWorkload` associated with
             OBSOLETE this updated deployment is not able to progress otherwise.
             OBSOLETE Therefore, make sure that the version of `podtato-head-frontend`
             OBSOLETE is updated accordingly:
             OBSOLETE 
             OBSOLETE ```yaml
             OBSOLETE {% include "./assets/app-updated-version.yaml" %}
             OBSOLETE ```
             OBSOLETE 
             OBSOLETE Updating the `KeptnApp` also causes all pre-/post-tasks/evaluations
             OBSOLETE of the `KeptnApp` to be executed again.
             OBSOLETE In this example, this means that the tasks `wait-for-prometheus`,
             OBSOLETE and `post-deployment-loadtests` will run again.
             OBSOLETE 
             OBSOLETE If you are using the [automatic app discovery](../guides/auto-app-discovery.md),
             OBSOLETE you do not need to update the `KeptnApp` resource.
             OBSOLETE Keptn will take care of that for you.

After applying the updated manifests, you can monitor the status
of the application and related [workloads](https://kubernetes.io/docs/concepts/workloads/) using the following commands:

```shell
$ kubectl get keptnworkloadversion -n podtato-kubectl

NAMESPACE   NAME                                             APPNAME         WORKLOADNAME                         WORKLOADVERSION      PHASE
podtato-kubectl   podtato-head-podtato-head-frontend-0.1.0   podtato-head    podtato-head-podtato-head-frontend   0.1.0                Completed
podtato-kubectl   podtato-head-podtato-head-hat-0.1.1        podtato-head    podtato-head-podtato-head-hat        0.1.1                Completed
podtato-kubectl   podtato-head-podtato-head-frontend-0.2.0   podtato-head    podtato-head-podtato-head-frontend   0.2.0                Completed
```

As can be seen in the output of the command, the `KeptnWorkloadVersions` from the previous deployment
are still here, but a new `KeptnWorkloadVersion` for the updated [workload](https://kubernetes.io/docs/concepts/workloads/)
has been added.
For the [workload](https://kubernetes.io/docs/concepts/workloads/) that
remained unchanged (`podtato-head-hat`), no new `KeptnWorkloadVersion` needed to be created.

Similarly, retrieving the list of `KeptnAppVersions` will reflect the update by
returning a newly created `KeptnAppVersion`.

```shell
$ kubectl get keptnappversion -n podtato-kubectl

NAMESPACE         NAME                          APPNAME        VERSION   PHASE
podtato-kubectl   podtato-head-0.1.0-6bch3iak   podtato-head   0.1.0     Completed
podtato-kubectl   podtato-head-0.1.0-hf52kauz   podtato-head   0.1.0     Completed
```

## Adding a new Workload to an Application

To add a new [workload](https://kubernetes.io/docs/concepts/workloads/) (e.g. a new deployment) to an existing app,
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
{% include "./assets/new-deployment.yaml" %}
```

         OBSOLETE The `KeptnApp`, if defined by the user, should contain the
         OBSOLETE reference to the newly added [workload](https://kubernetes.io/docs/concepts/workloads/).
         OBSOLETE This is mandatory, as the [workload](https://kubernetes.io/docs/concepts/workloads/) itself is not able to
         OBSOLETE progress if it is not part of a `KeptnApp`.
         OBSOLETE For automatically discovered apps this is done
         OBSOLETE automatically.
         OBSOLETE 
         OBSOLETE ```yaml
         OBSOLETE {% include "./assets/app-with-new-workload.yaml" %}
         OBSOLETE ```

After applying the updated manifests, you can monitor the status
of the application and related [workloads](https://kubernetes.io/docs/concepts/workloads/) using the following commands:

```shell
$ kubectl get keptnworkloadversion -n podtato-kubectl

NAMESPACE   NAME                                             APPNAME         WORKLOADNAME                         WORKLOADVERSION      PHASE
podtato-kubectl   podtato-head-podtato-head-frontend-0.1.0   podtato-head    podtato-head-podtato-head-frontend   0.1.0                Completed
podtato-kubectl   podtato-head-podtato-head-hat-0.1.1        podtato-head    podtato-head-podtato-head-hat        0.1.1                Completed
podtato-kubectl   podtato-head-podtato-head-left-leg-0.1.0   podtato-head    podtato-head-podtato-head-left-leg   0.1.0                Completed
```

As can be seen in the output of the command, in addition
to the previous `KeptnWorkloadVersions`, the newly created
`KeptnWorkloadVersion`, `podtato-head-podtato-head-left-leg-0.1.0` has been added
to the results.
