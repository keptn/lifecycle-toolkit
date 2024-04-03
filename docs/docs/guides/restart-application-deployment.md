---
comments: true
---

# Redeploy/Restart an Application

A [KeptnApp](../reference/crd-reference/app.md) can fail
because of an unsuccessful pre-deployment evaluation
or pre-deployment task.
For example, this happens if the target value of a
[KeptnEvaluationDefinition](../reference/crd-reference/evaluationdefinition.md)
resource is misconfigured
or a pre-deployment evaluation checks the wrong URL.

After you fix the configuration
that caused the pre-deployment evaluation or task to fail,
you can increment the `spec.revision` value
and apply the updated `KeptnApp` manifest
to create a new revision of the `KeptnApp`
without modifying the `version`.

Afterwards, all related `KeptnWorkloadVersions`
automatically refer to the newly created revision of the `KeptnAppVersion`
to determine whether they are allowed
to enter their respective deployment phases.

To illustrate this, consider the following example:

```yaml
{% include "./assets/restart-application-deployment/app-full.yaml" %}
```

Applying these resources causes a `KeptnApp`
to be created.
The `KeptnApp` contains all workloads that are part of
the application, and a version number that is
derived from the workloads:

```yaml
{% include "./assets/restart-application-deployment/generated-app.yaml" %}
```

Then, based on the `KeptnApp` and `KeptnAppContext`,
a `KeptnAppVersion` is created automatically and
the execution of the pre-deployment checks defined in
the `KeptnAppContext` starts.
In this example, the pre-deployment checks
contain a task which clearly fails.
Therefore, the `KeptnAppVersion` is not able
to proceed with the deployment.

You can inspect the status of the created `KeptnAppVersion`:

```shell
$ kubectl get keptnappversions.lifecycle.keptn.sh -n restartable-apps
NAME                   APPNAME        VERSION   PHASE
podtato-head-0.1.2-ab1223js   podtato-head   0.1.1     AppPreDeployTasks
```

Notice that the `KeptnAppVersion` stays
in the `AppPreDeployTasks` phase for a while,
due to the pre-check trying to run
until the failure threshold is reached.
Eventually, the `KeptnAppVersion`'s `PredeploymentPhase`
is in a `Failed` state, with the remaining phases being `Deprecated`.

<!-- markdownlint-disable MD013 -->
```shell
$ kubectl get keptnappversions.lifecycle.keptn.sh -n restartable-apps -owide
NAME                          APPNAME        VERSION   PHASE               PREDEPLOYMENTSTATUS   PREDEPLOYMENTEVALUATIONSTATUS   WORKLOADOVERALLSTATUS   POSTDEPLOYMENTSTATUS   POSTDEPLOYMENTEVALUATIONSTATUS
podtato-head-0.1.2-ab1223js   podtato-head   0.1.2     AppPreDeployTasks   Failed                Deprecated                      Deprecated              Deprecated             Deprecated
```
<!-- markdownlint-enable MD013 -->

To fix the deployment of this application,
we first need to fix the task that has failed earlier.
To do so, edit the `pre-deployment-check` `KeptnTaskDefinition`:

```shell
kubectl -n restartable-apps edit keptntaskdefinitions.lifecycle.keptn.sh pre-deployment-check
```

Modify the manifest to look like this:

```yaml
{% include "./assets/restart-application-deployment/fixed-task.yaml" %}
```

To restart the deployment of our `KeptnApplication`,
edit the manifest:

```shell
kubectl -n restartable-apps edit keptnapps.lifecycle.keptn.sh podtato-head
```

Increment the value of the `spec.revision` field by one:

```yaml
{% include "./assets/restart-application-deployment/generated-app-bumped-revision.yaml" %}
```

After those changes have been made,
you will notice a new revision of the `podtato-head` `KeptnAppVersion`:

```shell
$ kubectl get keptnappversions.lifecycle.keptn.sh -n restartable-apps       
NAME                          APPNAME        VERSION   PHASE
podtato-head-0.1.2-ab1223js   podtato-head   0.1.2     AppPreDeployTasks
podtato-head-0.1.2-xbhj9073   podtato-head   0.1.2     AppDeploy
```

See that the newly created revision `podtato-head-0.1.2-xbhj9073`
has made it beyond the pre-deployment check phase
and has reached its `AppDeployPhase`.

You can also verify the execution of the `pre-deployment-check`
by retrieving the list of `KeptnTasks` in
the `restartable-apps` namespace:

<!-- markdownlint-disable MD013 -->
```shell
$ kubectl get keptntasks.lifecycle.keptn.sh -n restartable-apps
NAME                             APPNAME        APPVERSION   WORKLOADNAME   WORKLOADVERSION   JOB NAME                              STATUS
pre-pre-deployment-check-49827   podtato-head   0.1.2                                         klc-pre-pre-deployment-check--77601   Failed
pre-pre-deployment-check-65056   podtato-head   0.1.2                                         klc-pre-pre-deployment-check--57313   Succeeded
```
<!-- markdownlint-enable MD013 -->

Notice that the previous failed instances are still available
for both `KeptnAppVersions` and `KeptnTasks`.
This may be useful historical data to keep track of
what went wrong during earlier deployment attempts.
