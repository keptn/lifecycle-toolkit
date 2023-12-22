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

{% include "./assets/restart-application-deployment.md_1.yaml" %}

In this example, the `KeptnApp` executes a pre-deployment check
which clearly fails because of the `pre-deployment-check` task,
and therefore is not able to proceed with the deployment.

After applying this manifest,
you can inspect the status of the created `KeptnAppVersion`:

```shell
$ kubectl get keptnappversions.lifecycle.keptn.sh -n restartable-apps
NAME                   APPNAME        VERSION   PHASE
podtato-head-0.1.1-1   podtato-head   0.1.1     AppPreDeployTasks
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
NAME                   APPNAME        VERSION   PHASE               PREDEPLOYMENTSTATUS   PREDEPLOYMENTEVALUATIONSTATUS   WORKLOADOVERALLSTATUS   POSTDEPLOYMENTSTATUS   POSTDEPLOYMENTEVALUATIONSTATUS
podtato-head-0.1.1-1   podtato-head   0.1.1     AppPreDeployTasks   Failed                Deprecated                      Deprecated              Deprecated             Deprecated
```
<!-- markdownlint-enable MD013 -->

To fix the deployment of this application,
we first need to fix the task that has failed earlier.
To do so, edit the `pre-deployment-check` `KeptnTaskDefinition`:

```shell
kubectl -n restartable-apps edit keptntaskdefinitions.lifecycle.keptn.sh pre-deployment-check
```

Modify the manifest to look like this:

{% include "./assets/restart-application-deployment.md_2.yaml" %}

To restart the deployment of our `KeptnApplication`,
edit the manifest:

```shell
kubectl -n restartable-apps edit keptnapps.lifecycle.keptn.sh podtato-head
```

Increment the value of the `spec.revision` field by one:

{% include "./assets/restart-application-deployment.md_3.yaml" %}

After those changes have been made,
you will notice a new revision of the `podtato-head` `KeptnAppVersion`:

```shell
$ kubectl get keptnappversions.lifecycle.keptn.sh -n restartable-apps       
NAME                   APPNAME        VERSION   PHASE
podtato-head-0.1.1-1   podtato-head   0.1.1     AppPreDeployTasks
podtato-head-0.1.1-2   podtato-head   0.1.1     AppDeploy
```

See that the newly created revision `podtato-head-0.1.1-2`
has made it beyond the pre-deployment check phase
and has reached its `AppDeployPhase`.

You can also verify the execution of the `pre-deployment-check`
by retrieving the list of `KeptnTasks` in
the `restartable-apps` namespace:

<!-- markdownlint-disable MD013 -->
```shell
$ kubectl get keptntasks.lifecycle.keptn.sh -n restartable-apps
NAME                             APPNAME        APPVERSION   WORKLOADNAME   WORKLOADVERSION   JOB NAME                              STATUS
pre-pre-deployment-check-49827   podtato-head   0.1.1                                         klc-pre-pre-deployment-check--77601   Failed
pre-pre-deployment-check-65056   podtato-head   0.1.1                                         klc-pre-pre-deployment-check--57313   Succeeded
```
<!-- markdownlint-enable MD013 -->

Notice that the previous failed instances are still available
for both `KeptnAppVersions` and `KeptnTasks`.
This may be useful historical data to keep track of
what went wrong during earlier deployment attempts.
