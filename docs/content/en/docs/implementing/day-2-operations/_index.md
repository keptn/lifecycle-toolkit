---
title: Day 2 Operations with Keptn
description: How to operate and maintain your Keptn Apps
weight: 500
---

After you have successfully rolled out your application by following
the instructions in the [integration guide](../integrate),
Keptn also assists you with day 2 operations for your application.

Tasks that fall under this category include:

* Updating the version of one or more workloads that are part of
the same application
* Adding a new workload to an existing application
* Monitoring the health of your application using `KeptnMetrics`
* Optimizing the resource usage of your applications by integrating
`KeptnMetrics` into a
[HorizontalPodAutoscaler (HPA)](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/)

## Updating Workload Versions

After a first successful deployment of your application with Keptn,
you will eventually need to update the version of one or
more workloads that are part of the application.
This usually involves updating the image of a deployment
and changing the configuration of a deployment.
For example, using a different service account name for a pod
could be seen as an update.
Regardless of that, however, it is the user who decides what Keptn
sees as a version bump in the application by setting the value of
the `keptn.sh/version` or `app.kubernetes.io/version` labels/annotations
in their workloads.

When this changes, Keptn interprets a change as a new version
and thus re-runs the pre- and post-tasks and evaluations for the application.

If the version label/annotation does not change, Keptn does not consider
a change of a workload configuration to be an update, and therefore no pre-
and post-tasks/evaluations are executed because they have already been
completed for the version set in the labels/annotations.

To illustrate the update of a workload,
let's assume the following example, including
a workload called `podtato-head-frontend` that includes a pre-task and
a pre-evaluation.

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: podtato-kubectl
  annotations:
    keptn.sh/lifecycle-toolkit: "enabled"

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: podtato-head-frontend
  namespace: podtato-kubectl
spec:
  template:
    metadata:
      labels:
        app.kubernetes.io/name: podtato-head-frontend
        app.kubernetes.io/part-of: podtato-head
        app.kubernetes.io/version: 0.1.0
        keptn.sh/pre-deployment-evaluations: my-evaluation
        keptn.sh/pre-deployment-tasks: my-task
    spec:
      containers:
        - name: podtato-head-frontend
          image: podtato-head-frontend:a
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: podtato-head-hat
  namespace: podtato-kubectl
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app.kubernetes.io/name: podtato-head-hat
        app.kubernetes.io/part-of: podtato-head
        app.kubernetes.io/version: 0.1.1
    spec:
      containers:
        - name: podtato-head-hat
          image: podtato-head-hat:a
```

Now, let's assume that the configuration of that workload needs to be changed.
In this example we assume that the image of that workload
should be updated, but a configuration change is not limited to that.
From here, you essentially have two options:

* **Only update the configuration *without* updating the `app.kubernetes.io/version`
label:** This can be useful when the change in the configuration should happen regardless
of the result of any task or evaluation, e.g., when the previously used image has a critical vulnerability
and the image must be updated as quickly as possible.
To do that, change `podtato-head-frontend` as follows:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: podtato-head-frontend
  namespace: podtato-kubectl
spec:
  template:
    metadata:
      labels:
        app.kubernetes.io/name: podtato-head-frontend
        app.kubernetes.io/part-of: podtato-head
        app.kubernetes.io/version: 0.1.0 # The version label stays the same and therefore no related KeptnWorkload will be affected
        keptn.sh/pre-deployment-evaluations: my-evaluation
        keptn.sh/pre-deployment-tasks: my-task
    spec:
      containers:
        - name: podtato-head-frontend
          image: podtato-head-frontend:b # Only the image tag has been updated from 'a' to 'b'
```

* **Update the configuration *and* the version label:**
   Doing so causes the `KeptnWorkload` that is associated
    with the `podtato-head-frontend` deployment to be updated,
    and therefore the pre-task `my-task` and pre-evaluation `my-evaluation`
    are executed before the updated pods are scheduled.
In this case, the deployment should be changed as follows:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: podtato-head-frontend
  namespace: podtato-kubectl
spec:
  template:
    metadata:
      labels:
        app.kubernetes.io/name: podtato-head-frontend
        app.kubernetes.io/part-of: podtato-head
        app.kubernetes.io/version: 0.2.0 # The version label is now updated as well
        keptn.sh/pre-deployment-evaluations: my-evaluation
        keptn.sh/pre-deployment-tasks: my-task
    spec:
      containers:
        - name: podtato-head-frontend
          image: podtato-head-frontend:b # The image tag has been updated from 'a' to 'b'
```

If you have defined the related `KeptnApp` resource yourself,
this must also be updated to refer to the updated `KeptnWorkload`.
This is a mandatory step, since the `KeptnWorkload` associated with
this updated deployment is not able to progress otherwise.
Therefore, make sure that the version of `podtato-head-frontend`
is updated accordingly:

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha3
kind: KeptnApp
metadata:
  name: podtato-head
  namespace: podtato-kubectl
spec:
  version: "0.1.0"
  preDeploymentTasks:
    - wait-for-prometheus
  postDeploymentTasks:
    - post-deployment-loadtests
  workloads:
  - name: podtato-head-frontend
    version: 0.2.0 # Updated version. Needs to be set to the value of the 'app.kubernetes.io/version' label of the updated deployment
  - name: podtato-head-hat
    version: 1.1.1
```

Updating the `KeptnApp` also causes all pre-/post-tasks/evaluations
of the `KeptnApp` to be executed again.
In this example, this means that the tasks `wait-for-prometheus`,
and `post-deployment-loadtests` will run again.

If you are using the [automatic app discovery](../integrate#use-keptn-automatic-app-discovery),
you do not need to update the `KeptnApp` resource.
Keptn will take care of that for you.

After applying the updated manifests, you can monitor the status
of the application and related workloads using the following commands:

```shell
$ kubectl get keptnworkloadversion -n podtato-kubectl

NAMESPACE   NAME                                             APPNAME         WORKLOADNAME                         WORKLOADVERSION      PHASE
podtato-kubectl   podtato-head-podtato-head-frontend-0.1.0   podtato-head    podtato-head-podtato-head-frontend   0.1.0                Completed
podtato-kubectl   podtato-head-podtato-head-hat-0.1.1        podtato-head    podtato-head-podtato-head-hat        0.1.1                Completed
podtato-kubectl   podtato-head-podtato-head-frontend-0.2.0   podtato-head    podtato-head-podtato-head-frontend   0.2.0                Completed
```

As can be seen in the output of the command, the `KeptnWorkloadVersions` from the previous deployment
are still here, but a new `KeptnWorkloadVersion` for the updated workload has been added.
For the workload that remained unchanged (`podtato-head-hat`), no new `KeptnWorkloadVersion`
needed to be created.

Similarly, retrieving the list of `KeptnAppVersions` will reflect the update by
returning a newly created `KeptnAppVersion`.

```shell
$ kubectl get keptnappversion -n podtato-kubectl

NAMESPACE         NAME                          APPNAME        VERSION   PHASE
podtato-kubectl   podtato-head-0.1.0-6bch3iak   podtato-head   0.1.0     Completed
podtato-kubectl   podtato-head-0.1.0-hf52kauz   podtato-head   0.1.0     Completed
```

## Adding a new Workload to an Application

To add a new workload (e.g. a new deployment) to an existing app,
you must:

* Make sure the
`keptn.sh/app`/`app.kubernetes.io/part-of` label/annotation is present
on the new workload
* Add the new workload to the `KeptnApp`,
if you have previously defined the `KeptnApp` resource manually.
If the application has been discovered automatically, this step is not needed.

For example, to add the deployment `podtato-head-left-leg` to the
`podtato-head` application, the configuration for that new deployment
would look like this, with the required label being set:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: podtato-head-left-leg
  namespace: podtato-kubectl
spec:
  template:
    metadata:
      labels:
        app.kubernetes.io/name: podtato-head-left-leg
        app.kubernetes.io/part-of: podtato-head
        app.kubernetes.io/version: 0.1.0
    spec:
      containers:
        - name: podtato-head-frontend
          image: podtato-head-left-leg:a
```

The `KeptnApp`, if defined by the user, should contain the
reference to the newly added workload.
This is mandatory, as the workload itself is not be able to
progress if it is not part of a `KeptnApp`.
For automatically discovered apps this is done
automatically.

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha3
kind: KeptnApp
metadata:
  name: podtato-head
  namespace: podtato-kubectl
spec:
  version: "0.1.0"
  preDeploymentTasks:
    - wait-for-prometheus
  postDeploymentTasks:
    - post-deployment-loadtests
  workloads:
  - name: podtato-head-frontend
    version: 0.1.0
  - name: podtato-head-hat
    version: 1.1.1
  - name: podtato-head-left-leg # The newly added workload
    version: 0.1.0
```

After applying the updated manifests, you can monitor the status
of the application and related workloads using the following commands:

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
