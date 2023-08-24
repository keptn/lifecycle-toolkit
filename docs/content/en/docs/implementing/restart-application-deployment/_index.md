---
title: Restart an Application Deployment
description: Learn how to restart an unsuccessful Keptn Application Deployment.
layout: quickstart
weight: 100
hidechildren: false # this flag hides all sub-pages in the sidebar-multicard.html
---

A [KeptnApp](../../yaml-crd-ref/app.md) can fail
because of an unsuccessful pre-deployment evaluation
or pre-deployment task.
For example, this happens if the target value of a
[KeptnEvaluationDefinition](../../yaml-crd-ref/evaluationdefinition.md)
resource is misconfigured
or a pre-deployment evaluation checks the wrong URL.

After you fix the configuration
that caused the pre-deployment evaluation or task to fail,
you can increment the `spec.revision` value
and apply the updated `KeptnApp` manifest
to create a new revision of the `KeptnApp`
without modifying the `version`.

Afterwards, all related `KeptnWorkloadInstances`
automatically refer to the newly created revision of the `KeptnAppVersion`
to determine whether they are allowed
to enter their respective deployment phases.

To illustrate this, consider the following example:

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: restartable-apps
  annotations:
    keptn.sh/lifecycle-toolkit: "enabled"
---
apiVersion: lifecycle.keptn.sh/v1alpha2
kind: KeptnApp
metadata:
  name: podtato-head
  namespace: restartable-apps
spec:
  version: "0.1.1"
  revision: 1
  workloads:
    - name: podtato-head-entry
      version: "0.1.2"
  preDeploymentTasks:
    - pre-deployment-check
---
apiVersion: lifecycle.keptn.sh/v1alpha2
kind: KeptnTaskDefinition
metadata:
  name: pre-deployment-check
  namespace: restartable-apps
spec:
  function:
    inline:
      code: |
        console.error("I failed")
        process.exit(1)
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: podtato-head-entry
  namespace: restartable-apps
  labels:
    app: podtato-head
spec:
  selector:
    matchLabels:
      component: podtato-head-entry
  template:
    metadata:
      labels:
        component: podtato-head-entry
        keptn.sh/workload: podtato-head-entry
        keptn.sh/app: podtato-head
        keptn.sh/version: "0.1.2"
    spec:
      terminationGracePeriodSeconds: 5
      containers:
        - name: server
          image: ghcr.io/podtato-head/entry:0.2.7
          imagePullPolicy: Always
          ports:
            - containerPort: 9000
          env:
            - name: PODTATO_PORT
              value: "9000"
```

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

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha2
kind: KeptnTaskDefinition
metadata:
  name: pre-deployment-check
  namespace: restartable-apps
spec:
  function:
    inline:
      code: |
        console.error("Success")
```

To restart the deployment of our `KeptnApplication`,
edit the manifest:

```shell
kubectl -n restartable-apps edit keptnapps.lifecycle.keptn.sh podtato-head
```

Increment the value of the `spec.revision` field by one:

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha2
kind: KeptnApp
metadata:
  name: podtato-head
  namespace: restartable-apps
spec:
  version: "0.1.1"
  revision: 2 # Increased this value from 1 to 2
  workloads:
    - name: podtato-head-entry
      version: "0.1.2"
  preDeploymentTasks:
    - pre-deployment-check
```

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
